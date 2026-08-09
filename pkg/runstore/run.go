package runstore

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// Run owns one active run directory.
type Run struct {
	mu       sync.Mutex
	dir      string
	manifest Manifest
	status   Status
	inputs   []InputRef
	terminal bool
}

// Resume explicitly reopens a validated active run for a single writer. The
// caller must supply the exact semantic config used by Create. Resume does not
// acquire an inter-process lock; callers are responsible for ensuring that the
// original writer has stopped and that only one resumed writer exists.
func Resume(ctx context.Context, dir string, expectedConfig any) (*Run, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, errors.Wrap(err, "resolve run directory")
	}
	if err := recoverPendingInputs(absDir); err != nil {
		return nil, errors.Wrap(err, "recover pending copied inputs")
	}
	reader, err := Open(absDir)
	if err != nil {
		return nil, errors.Wrap(err, "validate run before resume")
	}
	if reader.status.State != StateActive {
		return nil, errors.Errorf("cannot resume run in state %q", reader.status.State)
	}
	canonicalConfig, err := canonicalJSON(expectedConfig)
	if err != nil {
		return nil, errors.Wrap(err, "canonicalize expected resume config")
	}
	if actual := digestBytes(canonicalConfig); actual != reader.manifest.ConfigDigest {
		return nil, errors.Errorf("resume config digest mismatch: run=%s requested=%s", reader.manifest.ConfigDigest, actual)
	}
	return &Run{
		dir:      reader.dir,
		manifest: reader.manifest,
		status:   reader.status,
		inputs:   append([]InputRef(nil), reader.inputs...),
	}, nil
}

// Create initializes an active run and durably writes its configuration,
// manifest, and status.
func Create(ctx context.Context, options Options, config any) (*Run, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(options.Root) == "" {
		return nil, errors.New("run root is required")
	}
	root, err := filepath.Abs(options.Root)
	if err != nil {
		return nil, errors.Wrap(err, "resolve run root")
	}
	name, err := safeName("run name", options.Name)
	if err != nil {
		return nil, err
	}
	dimensions, err := validateDimensions(options.Dimensions)
	if err != nil {
		return nil, err
	}
	rawConfig, err := json.Marshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "marshal run config")
	}
	canonicalConfig, err := canonicalizeJSON(rawConfig)
	if err != nil {
		return nil, errors.Wrap(err, "canonicalize run config")
	}
	var prettyConfig bytes.Buffer
	if err := json.Indent(&prettyConfig, rawConfig, "", "  "); err != nil {
		return nil, errors.Wrap(err, "format run config")
	}

	startedAt := time.Now().UTC()
	runID, err := newRunID(startedAt, name)
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(root, runID)
	for _, child := range []string{"inputs", "results", "native"} {
		if err := os.MkdirAll(filepath.Join(dir, child), 0o700); err != nil {
			return nil, errors.Wrap(err, "create run directory")
		}
	}

	hostname, _ := os.Hostname()
	manifest := Manifest{
		SchemaVersion: ManifestSchemaVersion,
		RunID:         runID,
		Name:          options.Name,
		Description:   options.Description,
		StartedAt:     startedAt,
		GoVersion:     runtime.Version(),
		Host: Host{
			Hostname: hostname,
			OS:       runtime.GOOS,
			Arch:     runtime.GOARCH,
			CPUs:     runtime.NumCPU(),
		},
		ConfigDigest: digestBytes(canonicalConfig),
		Dimensions:   dimensions,
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		manifest.ModulePath = info.Main.Path
		manifest.ModuleVersion = info.Main.Version
	}
	run := &Run{
		dir:      dir,
		manifest: manifest,
		status:   Status{State: StateActive, StartedAt: startedAt},
	}
	if err := run.writeBytes(ctx, "config.json", append(prettyConfig.Bytes(), '\n')); err != nil {
		return nil, err
	}
	if err := run.writeJSON(ctx, "manifest.json", manifest); err != nil {
		return nil, err
	}
	if err := run.writeJSON(ctx, "status.json", run.status); err != nil {
		return nil, err
	}
	return run, nil
}

// Dir returns the run directory.
func (run *Run) Dir() string {
	if run == nil {
		return ""
	}
	return run.dir
}

// Manifest returns a defensive copy of the run manifest.
func (run *Run) Manifest() Manifest {
	if run == nil {
		return Manifest{}
	}
	run.mu.Lock()
	defer run.mu.Unlock()
	manifest := run.manifest
	manifest.Dimensions = cloneStrings(run.manifest.Dimensions)
	return manifest
}

// Inputs returns defensive copies of all copied-input records currently bound
// to the run.
func (run *Run) Inputs() []InputRef {
	if run == nil {
		return nil
	}
	run.mu.Lock()
	defer run.mu.Unlock()
	return append([]InputRef(nil), run.inputs...)
}

// Path returns a validated path within the run without creating it.
func (run *Run) Path(relative string) (string, error) {
	if run == nil {
		return "", errors.New("run is nil")
	}
	return joinWithin(run.dir, relative)
}

// WriteJSON atomically writes a JSON artifact inside an active run.
func (run *Run) WriteJSON(ctx context.Context, relative string, value any) error {
	if run == nil {
		return errors.New("run is nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	run.mu.Lock()
	defer run.mu.Unlock()
	if run.terminal {
		return errors.New("run is terminal")
	}
	return run.writeJSON(ctx, relative, value)
}

// WriteBytes atomically writes an artifact inside an active run.
func (run *Run) WriteBytes(ctx context.Context, relative string, data []byte) error {
	if run == nil {
		return errors.New("run is nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	run.mu.Lock()
	defer run.mu.Unlock()
	if run.terminal {
		return errors.New("run is terminal")
	}
	return run.writeBytes(ctx, relative, data)
}

// AppendJSONL appends and fsyncs one JSON record. Each successful return is a
// durability boundary suitable for interruption and later resume.
func (run *Run) AppendJSONL(ctx context.Context, relative string, value any) error {
	if run == nil {
		return errors.New("run is nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	path, err := joinWithin(run.dir, relative)
	if err != nil {
		return err
	}
	data, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "marshal JSONL record")
	}

	run.mu.Lock()
	defer run.mu.Unlock()
	if run.terminal {
		return errors.New("run is terminal")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return errors.Wrap(err, "create JSONL parent directory")
	}
	created := false
	if _, err := os.Lstat(path); os.IsNotExist(err) {
		created = true
	} else if err != nil {
		return errors.Wrap(err, "inspect JSONL artifact")
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return errors.Wrap(err, "open JSONL artifact")
	}
	if _, err := file.Write(append(data, '\n')); err != nil {
		_ = file.Close()
		return errors.Wrap(err, "append JSONL artifact")
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return errors.Wrap(err, "sync JSONL artifact")
	}
	if err := file.Close(); err != nil {
		return errors.Wrap(err, "close JSONL artifact")
	}
	if created {
		return errors.Wrap(syncDirectory(filepath.Dir(path)), "sync new JSONL directory entry")
	}
	return nil
}

func newRunID(startedAt time.Time, name string) (string, error) {
	random := make([]byte, 6)
	if _, err := rand.Read(random); err != nil {
		return "", errors.Wrap(err, "generate run ID")
	}
	prefix := startedAt.Format("20060102T150405.000000000Z") + "-"
	suffix := "-" + hex.EncodeToString(random)
	name = boundedSafeName(name, maximumFileComponentBytes-len(prefix)-len(suffix))
	return prefix + name + suffix, nil
}

func validateDimensions(dimensions map[string]string) (map[string]string, error) {
	if dimensions == nil {
		return nil, nil
	}
	keys := make([]string, 0, len(dimensions))
	for key := range dimensions {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make(map[string]string, len(dimensions))
	for _, key := range keys {
		trimmedKey := strings.TrimSpace(key)
		if trimmedKey == "" {
			return nil, errors.New("semantic dimension key is required")
		}
		if trimmedKey != key {
			return nil, errors.Errorf("semantic dimension key %q has surrounding whitespace", key)
		}
		value := strings.TrimSpace(dimensions[key])
		if value == "" {
			return nil, errors.Errorf("semantic dimension %q has an empty value", key)
		}
		result[key] = value
	}
	return result, nil
}

func cloneStrings(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}
	result := make(map[string]string, len(values))
	for key, value := range values {
		result[key] = value
	}
	return result
}
