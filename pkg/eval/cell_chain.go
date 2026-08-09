package eval

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

func sealCell(cell *Cell, previous string) error {
	if cell == nil {
		return errors.New("cell is nil")
	}
	cell.PreviousDigest = previous
	cell.Digest = ""
	data, err := json.Marshal(cell)
	if err != nil {
		return errors.Wrap(err, "marshal cell identity")
	}
	if len(data) > maximumCellRecordBytes {
		return errors.Errorf("cell record exceeds %d bytes", maximumCellRecordBytes)
	}
	cell.Digest = digestBytes(data)
	return nil
}

func validateCellChain(cell Cell, previous string) error {
	if cell.PreviousDigest != previous {
		return errors.Errorf("previous digest mismatch: cell=%q expected=%q", cell.PreviousDigest, previous)
	}
	claimed := cell.Digest
	if !strings.HasPrefix(claimed, "sha256:") || len(claimed) != len("sha256:")+64 {
		return errors.New("cell digest is invalid")
	}
	cell.Digest = ""
	data, err := json.Marshal(cell)
	if err != nil {
		return errors.Wrap(err, "marshal cell identity")
	}
	if len(data) > maximumCellRecordBytes {
		return errors.Errorf("cell record exceeds %d bytes", maximumCellRecordBytes)
	}
	if actual := digestBytes(data); actual != claimed {
		return errors.Errorf("cell digest mismatch: stored=%s actual=%s", claimed, actual)
	}
	return nil
}
