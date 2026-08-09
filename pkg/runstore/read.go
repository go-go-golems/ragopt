package runstore

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Reader exposes a validated run without granting mutation methods.
type Reader struct {
	dir      string
	manifest Manifest
	status   Status
	inputs   []InputRef
}

// Open validates the common run contract and every declared copied input.
// It accepts active runs so callers can inspect safely interrupted work.
func Open(dir string) (*Reader, error) {
	if strings.TrimSpace(dir) == "" {
		return nil, errors.New("run directory is required")
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, errors.Wrap(err, "resolve run directory")
	}
	info, err := os.Stat(absDir)
	if err != nil {
		return nil, errors.Wrap(err, "stat run directory")
	}
	if !info.IsDir() {
		return nil, errors.New("run path is not a directory")
	}

	var manifest Manifest
	if err := readStrictJSON(filepath.Join(absDir, "manifest.json"), &manifest); err != nil {
		return nil, errors.Wrap(err, "read run manifest")
	}
	if manifest.SchemaVersion != ManifestSchemaVersion {
		return nil, errors.Errorf("unsupported run manifest schema %q", manifest.SchemaVersion)
	}
	if manifest.RunID == "" || filepath.Base(absDir) != manifest.RunID {
		return nil, errors.Errorf("manifest run ID %q does not match directory %q", manifest.RunID, filepath.Base(absDir))
	}
	if strings.TrimSpace(manifest.Name) == "" {
		return nil, errors.New("manifest run name is required")
	}
	if manifest.StartedAt.IsZero() {
		return nil, errors.New("manifest start time is required")
	}
	validatedDimensions, err := validateDimensions(manifest.Dimensions)
	if err != nil {
		return nil, errors.Wrap(err, "validate manifest dimensions")
	}
	manifest.Dimensions = validatedDimensions

	configData, err := readRegularFile(filepath.Join(absDir, "config.json"))
	if err != nil {
		return nil, errors.Wrap(err, "read run config")
	}
	canonicalConfig, err := canonicalizeJSON(configData)
	if err != nil {
		return nil, errors.Wrap(err, "validate run config")
	}
	if actual := digestBytes(canonicalConfig); actual != manifest.ConfigDigest {
		return nil, errors.Errorf("config digest mismatch: manifest=%s actual=%s", manifest.ConfigDigest, actual)
	}

	var status Status
	if err := readStrictJSON(filepath.Join(absDir, "status.json"), &status); err != nil {
		return nil, errors.Wrap(err, "read run status")
	}
	if !status.StartedAt.Equal(manifest.StartedAt) {
		return nil, errors.New("status start time does not match manifest")
	}
	if err := validateStatus(status); err != nil {
		return nil, err
	}
	if status.State == StateComplete {
		var summary Summary
		if err := readStrictJSON(filepath.Join(absDir, "results", "summary.json"), &summary); err != nil {
			return nil, errors.Wrap(err, "read completed run summary")
		}
	}

	inputs, err := readInputs(absDir)
	if err != nil {
		return nil, err
	}
	return &Reader{dir: absDir, manifest: manifest, status: status, inputs: inputs}, nil
}

// Dir returns the validated run directory.
func (reader *Reader) Dir() string {
	if reader == nil {
		return ""
	}
	return reader.dir
}

// Manifest returns a defensive copy of the validated manifest.
func (reader *Reader) Manifest() Manifest {
	if reader == nil {
		return Manifest{}
	}
	manifest := reader.manifest
	manifest.Dimensions = cloneStrings(reader.manifest.Dimensions)
	return manifest
}

// Status returns the validated state.
func (reader *Reader) Status() Status {
	if reader == nil {
		return Status{}
	}
	status := reader.status
	if reader.status.FinishedAt != nil {
		finishedAt := *reader.status.FinishedAt
		status.FinishedAt = &finishedAt
	}
	return status
}

// Inputs returns defensive copies of the validated copied-input records.
func (reader *Reader) Inputs() []InputRef {
	if reader == nil {
		return nil
	}
	return append([]InputRef(nil), reader.inputs...)
}

// Path returns a confined path inside the validated run.
func (reader *Reader) Path(relative string) (string, error) {
	if reader == nil {
		return "", errors.New("run reader is nil")
	}
	return joinWithin(reader.dir, relative)
}

func validateStatus(status Status) error {
	switch status.State {
	case StateActive:
		if status.FinishedAt != nil || status.Error != "" {
			return errors.New("active run cannot have a finish time or error")
		}
	case StateComplete:
		if status.FinishedAt == nil {
			return errors.New("complete run requires a finish time")
		}
		if status.Error != "" {
			return errors.New("complete run cannot have an error")
		}
	case StateFailed:
		if status.FinishedAt == nil {
			return errors.New("failed run requires a finish time")
		}
	default:
		return errors.Errorf("unsupported run state %q", status.State)
	}
	if status.FinishedAt != nil && status.FinishedAt.Before(status.StartedAt) {
		return errors.New("run finish time precedes start time")
	}
	return nil
}

func readInputs(root string) ([]InputRef, error) {
	manifestPath := filepath.Join(root, "inputs", "manifest.json")
	var inputs []InputRef
	err := readStrictJSON(manifestPath, &inputs)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			entries, readErr := os.ReadDir(filepath.Join(root, "inputs"))
			if readErr != nil {
				return nil, errors.Wrap(readErr, "read input directory")
			}
			for _, entry := range entries {
				if !strings.HasPrefix(entry.Name(), ".pending-") {
					return nil, errors.New("input directory has files but no input manifest")
				}
			}
			return nil, nil
		}
		return nil, errors.Wrap(err, "read input manifest")
	}
	roles := make(map[string]struct{}, len(inputs))
	paths := make(map[string]struct{}, len(inputs))
	for index, input := range inputs {
		if strings.TrimSpace(input.Role) == "" {
			return nil, errors.Errorf("input %d has no role", index)
		}
		if _, exists := roles[input.Role]; exists {
			return nil, errors.Errorf("duplicate input role %q", input.Role)
		}
		roles[input.Role] = struct{}{}
		clean := filepath.Clean(input.CopiedPath)
		if filepath.Dir(clean) != "inputs" || filepath.Base(clean) == "manifest.json" {
			return nil, errors.Errorf("input %q has invalid copied path %q", input.Role, input.CopiedPath)
		}
		if _, exists := paths[clean]; exists {
			return nil, errors.Errorf("duplicate input copied path %q", clean)
		}
		paths[clean] = struct{}{}
		path, pathErr := joinWithin(root, clean)
		if pathErr != nil {
			return nil, errors.Wrapf(pathErr, "validate input %q path", input.Role)
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil, errors.Wrapf(readErr, "read copied input %q", input.Role)
		}
		if int64(len(data)) != input.SizeBytes {
			return nil, errors.Errorf("input %q size mismatch: manifest=%d actual=%d", input.Role, input.SizeBytes, len(data))
		}
		if actual := digestBytes(data); actual != input.SHA256 {
			return nil, errors.Errorf("input %q digest mismatch: manifest=%s actual=%s", input.Role, input.SHA256, actual)
		}
	}
	entries, err := os.ReadDir(filepath.Join(root, "inputs"))
	if err != nil {
		return nil, errors.Wrap(err, "read input directory")
	}
	for _, entry := range entries {
		if entry.Name() == "manifest.json" || strings.HasPrefix(entry.Name(), ".pending-") {
			continue
		}
		if _, ok := paths[filepath.Join("inputs", entry.Name())]; !ok {
			return nil, errors.Errorf("untracked copied input %q", entry.Name())
		}
	}
	return inputs, nil
}

// recoverPendingInputs completes or discards the small transaction used by
// CopyInput. A manifest that names the final file commits the pending bytes;
// without that manifest entry, the pending bytes were never committed.
func recoverPendingInputs(root string) error {
	directory := filepath.Join(root, "inputs")
	info, err := os.Lstat(directory)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return errors.New("input directory is not a real directory")
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	var pending []os.DirEntry
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".pending-") {
			pending = append(pending, entry)
		}
	}
	if len(pending) == 0 {
		return nil
	}
	var manifest []InputRef
	manifestErr := readStrictJSON(filepath.Join(directory, "manifest.json"), &manifest)
	if manifestErr != nil && !errors.Is(manifestErr, os.ErrNotExist) {
		return errors.Wrap(manifestErr, "read input manifest during recovery")
	}
	committed := make(map[string]struct{}, len(manifest))
	for _, input := range manifest {
		committed[filepath.Base(input.CopiedPath)] = struct{}{}
	}
	for _, entry := range pending {
		pendingPath := filepath.Join(directory, entry.Name())
		finalName := strings.TrimPrefix(entry.Name(), ".pending-")
		finalPath := filepath.Join(directory, finalName)
		if _, ok := committed[finalName]; ok {
			if _, err := os.Lstat(finalPath); err == nil {
				if err := os.Remove(pendingPath); err != nil {
					return errors.Wrap(err, "discard redundant pending input")
				}
				continue
			} else if !os.IsNotExist(err) {
				return err
			}
			if err := os.Rename(pendingPath, finalPath); err != nil {
				return errors.Wrap(err, "finish pending input publication")
			}
			continue
		}
		if err := os.Remove(pendingPath); err != nil {
			return errors.Wrap(err, "discard uncommitted pending input")
		}
	}
	return syncDirectory(directory)
}

func readStrictJSON(path string, destination any) error {
	data, err := readRegularFile(path)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return errors.Wrap(err, "decode JSON")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("JSON contains multiple values")
		}
		return errors.Wrap(err, "check JSON trailing data")
	}
	return nil
}

func readRegularFile(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, errors.Errorf("%s is not a regular file", path)
	}
	return os.ReadFile(path)
}
