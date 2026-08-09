package eval

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/go-go-golems/ragopt/pkg/runstore"
)

type runEvidenceSnapshot map[string]string

// snapshotRunEvidence hashes the fixed trust roots exposed to every arm plus
// the committed cell journal when present. Each cell authenticates its native
// artifact digest, so protecting the bounded journal identity prevents an arm
// from rewriting both an earlier cell and its artifact; artifact-only changes
// are rejected by the final audit.
func snapshotRunEvidence(run *runstore.Run) (runEvidenceSnapshot, error) {
	paths := []string{"manifest.json", "config.json", "status.json", filepath.Join("inputs", "manifest.json")}
	for _, input := range run.Inputs() {
		paths = append(paths, input.CopiedPath)
	}
	cellsPath := filepath.Join("results", "cells.jsonl")
	absoluteCells, err := run.Path(cellsPath)
	if err != nil {
		return nil, err
	}
	if _, err := os.Lstat(absoluteCells); err == nil {
		paths = append(paths, cellsPath)
	} else if !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "inspect committed cell journal")
	}
	snapshot := make(runEvidenceSnapshot, len(paths))
	for _, relative := range paths {
		path, err := run.Path(relative)
		if err != nil {
			return nil, err
		}
		info, err := os.Lstat(path)
		if err != nil {
			return nil, errors.Wrapf(err, "inspect run evidence %q", relative)
		}
		if !info.Mode().IsRegular() {
			return nil, errors.Errorf("run evidence %q is not a regular file", relative)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "read run evidence %q", relative)
		}
		snapshot[relative] = digestBytes(data)
	}
	return snapshot, nil
}

func verifyRunEvidence(run *runstore.Run, expected runEvidenceSnapshot) error {
	actual, err := snapshotRunEvidence(run)
	if err != nil {
		return err
	}
	if len(actual) != len(expected) {
		return errors.Errorf("run evidence file count changed: before=%d after=%d", len(expected), len(actual))
	}
	for path, digest := range expected {
		actualDigest, exists := actual[path]
		if !exists {
			return errors.Errorf("run evidence %q was removed", path)
		}
		if actualDigest != digest {
			return errors.Errorf("run evidence %q changed", path)
		}
	}
	return nil
}
