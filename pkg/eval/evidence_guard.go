package eval

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type runEvidenceSnapshot map[string]string

// snapshotRunEvidence identifies every pre-existing run-owned file outside the
// current cell's native directory. Arms may create evidence only inside that
// directory; all other run state is immutable for the duration of the call.
func snapshotRunEvidence(runDirectory, excludedDirectory string) (runEvidenceSnapshot, error) {
	runDirectory = filepath.Clean(runDirectory)
	excludedDirectory = filepath.Clean(excludedDirectory)
	snapshot := runEvidenceSnapshot{}
	err := filepath.WalkDir(runDirectory, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if filepath.Clean(path) == excludedDirectory {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(runDirectory, path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			snapshot[relative] = "symlink:" + target
			return nil
		}
		if !info.Mode().IsRegular() {
			return errors.Errorf("run evidence %q is not a regular file", relative)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		snapshot[relative] = digestBytes(data)
		return nil
	})
	return snapshot, errors.Wrap(err, "walk run evidence")
}

func verifyRunEvidence(runDirectory, excludedDirectory string, expected runEvidenceSnapshot) error {
	actual, err := snapshotRunEvidence(runDirectory, excludedDirectory)
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
