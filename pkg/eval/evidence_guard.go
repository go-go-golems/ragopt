package eval

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/go-go-golems/ragopt/pkg/runstore"
)

type runEvidenceSnapshot map[string]string

// snapshotRunEvidence hashes the fixed trust roots exposed to every arm. The
// set is bounded by configuration plus copied inputs; committed cells and
// native artifacts are audited once as a whole before run completion.
func snapshotRunEvidence(run *runstore.Run) (runEvidenceSnapshot, error) {
	paths := []string{"manifest.json", "config.json", "status.json", filepath.Join("inputs", "manifest.json")}
	for _, input := range run.Inputs() {
		paths = append(paths, input.CopiedPath)
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
