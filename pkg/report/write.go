package report

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Write publishes the report and JSON plan atomically to explicit caller-owned
// paths. It never writes into or changes the evaluated run.
func Write(ctx context.Context, document *Document, markdownPath, planPath string) error {
	if document == nil {
		return errors.New("report document is required")
	}
	if markdownPath == "" || planPath == "" {
		return errors.New("markdown and plan output paths are required")
	}
	markdownAbsolute, err := resolveDestination(markdownPath)
	if err != nil {
		return errors.Wrap(err, "resolve markdown output path")
	}
	planAbsolute, err := resolveDestination(planPath)
	if err != nil {
		return errors.Wrap(err, "resolve plan output path")
	}
	if markdownAbsolute == planAbsolute {
		return errors.New("markdown and plan output paths must differ")
	}
	plan, err := json.MarshalIndent(document.Plan, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshal promotion plan")
	}
	plan = append(plan, '\n')
	outputs := []*stagedOutput{
		{path: markdownAbsolute, data: []byte(document.Markdown)},
		{path: planAbsolute, data: plan},
	}
	for _, output := range outputs {
		if err := output.stage(ctx); err != nil {
			cleanupStaged(outputs)
			return errors.Wrap(err, "stage report outputs")
		}
	}
	if err := publishOutputs(ctx, outputs); err != nil {
		cleanupStaged(outputs)
		return errors.Wrap(err, "publish report outputs")
	}
	return nil
}

// ValidateOutputsOutsideRun rejects report destinations that would overwrite
// any file in the immutable evaluated run.
func ValidateOutputsOutsideRun(runDirectory string, paths ...string) error {
	runAbsolute, err := resolveDestination(runDirectory)
	if err != nil {
		return errors.Wrap(err, "resolve evaluated run directory")
	}
	for _, path := range paths {
		outputAbsolute, err := resolveDestination(path)
		if err != nil {
			return errors.Wrap(err, "resolve report output path")
		}
		relative, err := filepath.Rel(runAbsolute, outputAbsolute)
		if err != nil {
			return errors.Wrap(err, "compare report output with evaluated run")
		}
		if relative == "." || (relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))) {
			return errors.Errorf("report output path %q is inside evaluated run %q", outputAbsolute, runAbsolute)
		}
	}
	return nil
}

// resolveDestination canonicalizes symlinks in every existing path component
// while preserving a not-yet-created suffix. This is the identity atomicWrite
// will actually address after it creates missing parent directories.
func resolveDestination(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	current := filepath.Clean(absolute)
	var suffix []string
	for {
		_, err := os.Lstat(current)
		if err == nil {
			resolved, err := filepath.EvalSymlinks(current)
			if err != nil {
				return "", err
			}
			for index := len(suffix) - 1; index >= 0; index-- {
				resolved = filepath.Join(resolved, suffix[index])
			}
			return filepath.Clean(resolved), nil
		}
		if !os.IsNotExist(err) {
			return "", err
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", err
		}
		suffix = append(suffix, filepath.Base(current))
		current = parent
	}
}

type stagedOutput struct {
	path      string
	data      []byte
	temporary string
	backup    string
	existed   bool
	published bool
}

func (o *stagedOutput) stage(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	directory := filepath.Dir(o.path)
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return errors.Wrap(err, "create output directory")
	}
	info, err := os.Lstat(o.path)
	if err == nil {
		if !info.Mode().IsRegular() {
			return errors.Errorf("output destination %q is not a regular file", o.path)
		}
		o.existed = true
	} else if !os.IsNotExist(err) {
		return errors.Wrap(err, "inspect output destination")
	}
	temporary, err := os.CreateTemp(directory, ".ragopt-report-*")
	if err != nil {
		return errors.Wrap(err, "create temporary output")
	}
	o.temporary = temporary.Name()
	defer func() { _ = temporary.Close() }()
	if err := temporary.Chmod(0o600); err != nil {
		return errors.Wrap(err, "set output permissions")
	}
	if _, err := temporary.Write(o.data); err != nil {
		return errors.Wrap(err, "write temporary output")
	}
	if err := temporary.Sync(); err != nil {
		return errors.Wrap(err, "sync temporary output")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return errors.Wrap(temporary.Close(), "close temporary output")
}

func publishOutputs(ctx context.Context, outputs []*stagedOutput) error {
	for _, output := range outputs {
		if !output.existed {
			continue
		}
		placeholder, err := os.CreateTemp(filepath.Dir(output.path), ".ragopt-report-backup-*")
		if err != nil {
			rollbackOutputs(outputs)
			return errors.Wrap(err, "reserve output backup")
		}
		backupPath := placeholder.Name()
		if err := placeholder.Close(); err != nil {
			_ = os.Remove(backupPath)
			rollbackOutputs(outputs)
			return errors.Wrap(err, "close output backup placeholder")
		}
		if err := os.Remove(backupPath); err != nil {
			rollbackOutputs(outputs)
			return errors.Wrap(err, "remove output backup placeholder")
		}
		if err := os.Rename(output.path, backupPath); err != nil {
			rollbackOutputs(outputs)
			return errors.Wrap(err, "backup existing output")
		}
		output.backup = backupPath
	}
	for _, output := range outputs {
		if err := ctx.Err(); err != nil {
			rollbackOutputs(outputs)
			return err
		}
		if err := os.Rename(output.temporary, output.path); err != nil {
			rollbackOutputs(outputs)
			return errors.Wrap(err, "publish staged output")
		}
		output.temporary = ""
		output.published = true
	}
	for _, directory := range outputDirectories(outputs) {
		if err := syncDirectory(directory); err != nil {
			rollbackOutputs(outputs)
			return err
		}
	}
	for _, output := range outputs {
		if output.backup != "" {
			if err := os.Remove(output.backup); err != nil {
				return errors.Wrap(err, "remove committed output backup")
			}
			output.backup = ""
		}
	}
	for _, directory := range outputDirectories(outputs) {
		if err := syncDirectory(directory); err != nil {
			return err
		}
	}
	return nil
}

func rollbackOutputs(outputs []*stagedOutput) {
	for _, output := range outputs {
		if output.published {
			_ = os.Remove(output.path)
			output.published = false
		}
	}
	for _, output := range outputs {
		if output.backup != "" {
			_ = os.Rename(output.backup, output.path)
			output.backup = ""
		}
	}
	for _, directory := range outputDirectories(outputs) {
		_ = syncDirectory(directory)
	}
}

func cleanupStaged(outputs []*stagedOutput) {
	for _, output := range outputs {
		if output.temporary != "" {
			_ = os.Remove(output.temporary)
			output.temporary = ""
		}
	}
}

func outputDirectories(outputs []*stagedOutput) []string {
	seen := make(map[string]struct{}, len(outputs))
	result := make([]string, 0, len(outputs))
	for _, output := range outputs {
		directory := filepath.Dir(output.path)
		if _, exists := seen[directory]; exists {
			continue
		}
		seen[directory] = struct{}{}
		result = append(result, directory)
	}
	return result
}

func syncDirectory(directory string) error {
	handle, err := os.Open(directory)
	if err != nil {
		return errors.Wrap(err, "open output directory")
	}
	if err := handle.Sync(); err != nil {
		_ = handle.Close()
		return errors.Wrap(err, "sync output directory")
	}
	return errors.Wrap(handle.Close(), "close output directory")
}
