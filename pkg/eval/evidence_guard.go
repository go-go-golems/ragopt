package eval

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/go-go-golems/ragopt/pkg/runstore"
)

const maximumCellRecordBytes = 16 << 20

type runEvidenceSnapshot struct {
	files         map[string]string
	journalDigest string
	journalInfo   os.FileInfo
}

// snapshotRunEvidence hashes the fixed trust roots exposed to every arm and
// authenticates the bounded tail of the chained cell journal. The complete
// chain is read once during resume/final audit, not once per arm.
func snapshotRunEvidence(run *runstore.Run, expectedJournalDigest string) (runEvidenceSnapshot, error) {
	paths := []string{"manifest.json", "config.json", "status.json", filepath.Join("inputs", "manifest.json")}
	for _, input := range run.Inputs() {
		paths = append(paths, input.CopiedPath)
	}
	snapshot := runEvidenceSnapshot{
		files:         make(map[string]string, len(paths)),
		journalDigest: expectedJournalDigest,
	}
	for _, relative := range paths {
		path, err := run.Path(relative)
		if err != nil {
			return runEvidenceSnapshot{}, err
		}
		info, err := os.Lstat(path)
		if err != nil {
			return runEvidenceSnapshot{}, errors.Wrapf(err, "inspect run evidence %q", relative)
		}
		if !info.Mode().IsRegular() {
			return runEvidenceSnapshot{}, errors.Errorf("run evidence %q is not a regular file", relative)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return runEvidenceSnapshot{}, errors.Wrapf(err, "read run evidence %q", relative)
		}
		snapshot.files[relative] = digestBytes(data)
	}
	actualJournalDigest, journalInfo, err := journalTailState(run)
	if err != nil {
		return runEvidenceSnapshot{}, err
	}
	if actualJournalDigest != expectedJournalDigest {
		return runEvidenceSnapshot{}, errors.Errorf("cell journal head changed: actual=%q expected=%q", actualJournalDigest, expectedJournalDigest)
	}
	snapshot.journalInfo = journalInfo
	return snapshot, nil
}

func verifyRunEvidence(run *runstore.Run, expected runEvidenceSnapshot) error {
	actual, err := snapshotRunEvidence(run, expected.journalDigest)
	if err != nil {
		return err
	}
	if !sameJournalState(expected.journalInfo, actual.journalInfo) {
		return errors.New("committed cell journal changed during arm execution")
	}
	if len(actual.files) != len(expected.files) {
		return errors.Errorf("run evidence file count changed: before=%d after=%d", len(expected.files), len(actual.files))
	}
	for path, digest := range expected.files {
		actualDigest, exists := actual.files[path]
		if !exists {
			return errors.Errorf("run evidence %q was removed", path)
		}
		if actualDigest != digest {
			return errors.Errorf("run evidence %q changed", path)
		}
	}
	return nil
}

func journalTailState(run *runstore.Run) (string, os.FileInfo, error) {
	path, err := run.Path(filepath.Join("results", "cells.jsonl"))
	if err != nil {
		return "", nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil, nil
		}
		return "", nil, errors.Wrap(err, "open committed cell journal")
	}
	defer func() { _ = file.Close() }()
	info, err := file.Stat()
	if err != nil {
		return "", nil, errors.Wrap(err, "stat committed cell journal")
	}
	if !info.Mode().IsRegular() {
		return "", nil, errors.New("committed cell journal is not a regular file")
	}
	if info.Size() == 0 {
		return "", info, nil
	}
	offset := max(info.Size()-maximumCellRecordBytes, 0)
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return "", nil, errors.Wrap(err, "seek committed cell journal tail")
	}
	data, err := io.ReadAll(io.LimitReader(file, maximumCellRecordBytes+1))
	if err != nil {
		return "", nil, errors.Wrap(err, "read committed cell journal tail")
	}
	if len(data) == 0 || data[len(data)-1] != '\n' {
		return "", nil, errors.New("committed cell journal has a truncated tail")
	}
	data = bytes.TrimSuffix(data, []byte{'\n'})
	lastNewline := bytes.LastIndexByte(data, '\n')
	if lastNewline >= 0 {
		data = data[lastNewline+1:]
	} else if offset > 0 {
		return "", nil, errors.Errorf("cell journal record exceeds %d bytes", maximumCellRecordBytes)
	}
	var cell Cell
	if err := decodeStrictJSON(data, &cell); err != nil {
		return "", nil, errors.Wrap(err, "decode committed cell journal tail")
	}
	return cell.Digest, info, nil
}

func sameJournalState(left, right os.FileInfo) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return os.SameFile(left, right) && left.Size() == right.Size() && left.ModTime().Equal(right.ModTime())
}
