package eval

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/go-go-golems/ragopt/pkg/runstore"
)

func loadCompletedCells(
	run *runstore.Run,
	prepared *preparedRequest,
	incumbentView CandidateView,
	challengerView CandidateView,
) (map[string]Cell, error) {
	path, err := run.Path(filepath.Join("results", "cells.jsonl"))
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]Cell{}, nil
		}
		return nil, errors.Wrap(err, "read result cells")
	}
	if len(data) > 0 && data[len(data)-1] != '\n' {
		lastNewline := bytes.LastIndexByte(data, '\n')
		keep := 0
		if lastNewline >= 0 {
			keep = lastNewline + 1
		}
		if err := truncateAndSync(path, int64(keep)); err != nil {
			return nil, errors.Wrap(err, "discard uncommitted truncated JSONL tail")
		}
		data = data[:keep]
	}

	expected := make(map[string]scheduledCell)
	for _, item := range buildSchedule(prepared, incumbentView, challengerView) {
		expected[expectedCellKey(prepared, item)] = item
	}
	completed := make(map[string]Cell)
	lines := bytes.Split(data, []byte{'\n'})
	for index, line := range lines {
		if len(line) == 0 {
			if index != len(lines)-1 {
				return nil, errors.Errorf("result cell line %d is blank", index+1)
			}
			continue
		}
		var cell Cell
		if err := decodeStrictJSON(line, &cell); err != nil {
			return nil, errors.Wrapf(err, "decode result cell line %d", index+1)
		}
		key := cellKey(cell)
		item, exists := expected[key]
		if !exists {
			return nil, errors.Errorf("result cell line %d has an unexpected identity", index+1)
		}
		if _, exists := completed[key]; exists {
			return nil, errors.Errorf("duplicate result cell key on line %d", index+1)
		}
		if err := validateStoredCell(run, prepared, item, &cell); err != nil {
			return nil, errors.Wrapf(err, "validate result cell line %d", index+1)
		}
		completed[key] = cell
	}
	return completed, nil
}

func validateStoredCell(run *runstore.Run, prepared *preparedRequest, item scheduledCell, cell *Cell) error {
	if cell.APIVersion != CellAPIVersion {
		return errors.Errorf("unsupported cell API version %q", cell.APIVersion)
	}
	if cell.RunID != run.Manifest().RunID {
		return errors.Errorf("cell run ID %q does not match run %q", cell.RunID, run.Manifest().RunID)
	}
	if cell.StartedAt.IsZero() || cell.FinishedAt.IsZero() || cell.FinishedAt.Before(cell.StartedAt) {
		return errors.New("cell has invalid timestamps")
	}
	if cell.CaseID != item.caseValue.ID || cell.RepeatIndex != item.repeat || cell.Arm != item.armName {
		return errors.New("cell schedule coordinates do not match expected work")
	}
	if cell.CandidateID != prepared.config.CandidateID || cell.SnapshotDigest != item.view.SnapshotDigest ||
		cell.SuiteDigest != prepared.config.SuiteDigest || cell.PolicyDigest != prepared.config.PolicyDigest {
		return errors.New("cell semantic identity does not match resumed run")
	}
	nativeRelative := filepath.Join("native", item.armName, item.caseValue.ID, formatRepeat(item.repeat))
	nativeDirectory, err := run.Path(nativeRelative)
	if err != nil {
		return err
	}
	return validateStoredOutcome(run.Dir(), nativeDirectory, &cell.Outcome)
}

func truncateAndSync(path string, size int64) error {
	file, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return errors.Wrap(err, "open JSONL artifact for recovery")
	}
	if err := file.Truncate(size); err != nil {
		_ = file.Close()
		return errors.Wrap(err, "truncate JSONL artifact")
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return errors.Wrap(err, "sync recovered JSONL artifact")
	}
	if err := file.Close(); err != nil {
		return errors.Wrap(err, "close recovered JSONL artifact")
	}
	directory, err := os.Open(filepath.Dir(path))
	if err != nil {
		return errors.Wrap(err, "open JSONL directory for recovery sync")
	}
	if err := directory.Sync(); err != nil {
		_ = directory.Close()
		return errors.Wrap(err, "sync JSONL directory after recovery")
	}
	return errors.Wrap(directory.Close(), "close JSONL recovery directory")
}

func formatRepeat(repeat int) string {
	return fmt.Sprintf("%04d", repeat)
}
