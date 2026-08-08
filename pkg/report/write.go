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
	markdownAbsolute, err := filepath.Abs(markdownPath)
	if err != nil {
		return errors.Wrap(err, "resolve markdown output path")
	}
	planAbsolute, err := filepath.Abs(planPath)
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
	if err := atomicWrite(ctx, markdownAbsolute, []byte(document.Markdown)); err != nil {
		return errors.Wrap(err, "write promotion report")
	}
	if err := atomicWrite(ctx, planAbsolute, plan); err != nil {
		return errors.Wrap(err, "write promotion plan")
	}
	return nil
}

// ValidateOutputsOutsideRun rejects report destinations that would overwrite
// any file in the immutable evaluated run.
func ValidateOutputsOutsideRun(runDirectory string, paths ...string) error {
	runAbsolute, err := filepath.Abs(runDirectory)
	if err != nil {
		return errors.Wrap(err, "resolve evaluated run directory")
	}
	for _, path := range paths {
		outputAbsolute, err := filepath.Abs(path)
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

func atomicWrite(ctx context.Context, path string, data []byte) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return errors.Wrap(err, "resolve output path")
	}
	directory := filepath.Dir(absolute)
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return errors.Wrap(err, "create output directory")
	}
	temporary, err := os.CreateTemp(directory, ".ragopt-report-*")
	if err != nil {
		return errors.Wrap(err, "create temporary output")
	}
	temporaryPath := temporary.Name()
	defer func() {
		_ = temporary.Close()
		_ = os.Remove(temporaryPath)
	}()
	if err := temporary.Chmod(0o600); err != nil {
		return errors.Wrap(err, "set output permissions")
	}
	if _, err := temporary.Write(data); err != nil {
		return errors.Wrap(err, "write temporary output")
	}
	if err := temporary.Sync(); err != nil {
		return errors.Wrap(err, "sync temporary output")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := temporary.Close(); err != nil {
		return errors.Wrap(err, "close temporary output")
	}
	if err := os.Rename(temporaryPath, absolute); err != nil {
		return errors.Wrap(err, "publish output")
	}
	directoryHandle, err := os.Open(directory)
	if err != nil {
		return errors.Wrap(err, "open output directory")
	}
	if err := directoryHandle.Sync(); err != nil {
		_ = directoryHandle.Close()
		return errors.Wrap(err, "sync output directory")
	}
	return errors.Wrap(directoryHandle.Close(), "close output directory")
}
