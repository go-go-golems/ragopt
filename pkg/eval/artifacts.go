package eval

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/go-go-golems/ragopt/pkg/runstore"
)

// ArtifactRun is a strictly loaded read-only paired evaluation run.
type ArtifactRun struct {
	Directory  string
	Manifest   runstore.Manifest
	Status     runstore.Status
	Config     RunConfig
	Suite      *SuiteDocument
	PolicyPath string
	Cells      []Cell
	Inputs     []runstore.InputRef
}

// LoadArtifactRun validates a run's common artifacts, evaluation config,
// complete copied-input set, suite, committed cells, and native artifacts.
// It never repairs or mutates the run.
func LoadArtifactRun(ctx context.Context, directory string) (*ArtifactRun, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	reader, err := runstore.Open(directory)
	if err != nil {
		return nil, errors.Wrap(err, "open evaluation run")
	}
	configPath, err := reader.Path("config.json")
	if err != nil {
		return nil, err
	}
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "read evaluation config")
	}
	var config RunConfig
	if err := decodeStrictJSON(configData, &config); err != nil {
		return nil, errors.Wrap(err, "decode evaluation config")
	}
	if err := validateRunConfig(config); err != nil {
		return nil, err
	}
	inputs, inputByRole, err := validateArtifactInputs(reader, config.InputDigests)
	if err != nil {
		return nil, err
	}
	suiteInput, exists := inputByRole[roleSuite]
	if !exists {
		return nil, errors.New("evaluation run is missing suite input")
	}
	suitePath, err := reader.Path(suiteInput.CopiedPath)
	if err != nil {
		return nil, err
	}
	suite, err := LoadSuite(ctx, suitePath)
	if err != nil {
		return nil, errors.Wrap(err, "load copied evaluation suite")
	}
	if suite.Digest != config.SuiteDigest {
		return nil, errors.Errorf("copied suite digest mismatch: config=%s actual=%s", config.SuiteDigest, suite.Digest)
	}
	policyInput, exists := inputByRole[rolePolicy]
	if !exists {
		return nil, errors.New("evaluation run is missing gate policy input")
	}
	if policyInput.SHA256 != config.PolicyDigest {
		return nil, errors.Errorf("copied gate policy digest mismatch: config=%s actual=%s", config.PolicyDigest, policyInput.SHA256)
	}
	policyPath, err := reader.Path(policyInput.CopiedPath)
	if err != nil {
		return nil, err
	}
	cells, err := loadArtifactCells(reader, config, suite)
	if err != nil {
		return nil, err
	}
	return &ArtifactRun{
		Directory:  reader.Dir(),
		Manifest:   reader.Manifest(),
		Status:     reader.Status(),
		Config:     config,
		Suite:      suite,
		PolicyPath: policyPath,
		Cells:      cells,
		Inputs:     inputs,
	}, nil
}

func validateRunConfig(config RunConfig) error {
	if config.APIVersion != RunAPIVersion {
		return errors.Errorf("unsupported evaluation run API version %q", config.APIVersion)
	}
	for label, digest := range map[string]string{
		"suite": config.SuiteDigest, "policy": config.PolicyDigest,
		"candidate": config.CandidateDigest, "parent snapshot": config.ParentSnapshot,
		"child snapshot": config.ChildSnapshot, "parent asset": config.ParentAssetDigest,
		"child asset": config.ChildAssetDigest,
	} {
		if !strings.HasPrefix(digest, "sha256:") || len(digest) != len("sha256:")+64 {
			return errors.Errorf("%s digest is invalid", label)
		}
	}
	if !identifierPattern.MatchString(config.CandidateID) || !identifierPattern.MatchString(config.ChangedAsset) {
		return errors.New("evaluation candidate or changed-asset identity is invalid")
	}
	if !identifierPattern.MatchString(config.IncumbentArm) || !identifierPattern.MatchString(config.ChallengerArm) || config.IncumbentArm == config.ChallengerArm {
		return errors.New("evaluation arm identities are invalid")
	}
	if config.Repeats <= 0 || config.Repeats > 10000 {
		return errors.New("evaluation repeat count is invalid")
	}
	if len(config.InputDigests) == 0 {
		return errors.New("evaluation input digest map is required")
	}
	return nil
}

func validateArtifactInputs(reader *runstore.Reader, expected map[string]string) ([]runstore.InputRef, map[string]runstore.InputRef, error) {
	inputs := reader.Inputs()
	if len(inputs) != len(expected) {
		return nil, nil, errors.Errorf("evaluation input count mismatch: run=%d config=%d", len(inputs), len(expected))
	}
	byRole := make(map[string]runstore.InputRef, len(inputs))
	for _, input := range inputs {
		if _, exists := byRole[input.Role]; exists {
			return nil, nil, errors.Errorf("duplicate evaluation input role %q", input.Role)
		}
		byRole[input.Role] = input
	}
	for role, digest := range expected {
		input, exists := byRole[role]
		if !exists {
			return nil, nil, errors.Errorf("evaluation input role %q is missing", role)
		}
		if input.SHA256 != digest {
			return nil, nil, errors.Errorf("evaluation input role %q digest mismatch", role)
		}
	}
	return inputs, byRole, nil
}

type artifactExpectedCell struct {
	caseValue Case
	repeat    int
	arm       string
	snapshot  string
}

func loadArtifactCells(reader *runstore.Reader, config RunConfig, suite *SuiteDocument) ([]Cell, error) {
	path, err := reader.Path(filepath.Join("results", "cells.jsonl"))
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "read evaluation cells")
	}
	if len(data) > 0 && data[len(data)-1] != '\n' {
		return nil, errors.New("evaluation cells have a truncated final line")
	}
	expected := make(map[string]artifactExpectedCell)
	for _, caseValue := range suite.Suite.Cases {
		for repeat := 0; repeat < config.Repeats; repeat++ {
			for _, arm := range []struct {
				name     string
				snapshot string
			}{
				{name: config.IncumbentArm, snapshot: config.ParentSnapshot},
				{name: config.ChallengerArm, snapshot: config.ChildSnapshot},
			} {
				cell := Cell{
					SuiteDigest: config.SuiteDigest, PolicyDigest: config.PolicyDigest,
					CandidateID: config.CandidateID, SnapshotDigest: arm.snapshot,
					CaseID: caseValue.ID, RepeatIndex: repeat, Arm: arm.name,
				}
				expected[cellKey(cell)] = artifactExpectedCell{caseValue: caseValue, repeat: repeat, arm: arm.name, snapshot: arm.snapshot}
			}
		}
	}
	seen := make(map[string]struct{})
	lines := bytes.Split(data, []byte{'\n'})
	cells := make([]Cell, 0, len(lines)-1)
	for index, line := range lines {
		if len(line) == 0 {
			if index != len(lines)-1 {
				return nil, errors.Errorf("evaluation cell line %d is blank", index+1)
			}
			continue
		}
		var cell Cell
		if err := decodeStrictJSON(line, &cell); err != nil {
			return nil, errors.Wrapf(err, "decode evaluation cell line %d", index+1)
		}
		key := cellKey(cell)
		expectedCell, exists := expected[key]
		if !exists {
			return nil, errors.Errorf("evaluation cell line %d has unexpected identity", index+1)
		}
		if _, exists := seen[key]; exists {
			return nil, errors.Errorf("duplicate evaluation cell on line %d", index+1)
		}
		seen[key] = struct{}{}
		if cell.APIVersion != CellAPIVersion || cell.RunID != reader.Manifest().RunID {
			return nil, errors.Errorf("evaluation cell line %d has invalid schema or run ID", index+1)
		}
		if cell.StartedAt.IsZero() || cell.FinishedAt.IsZero() || cell.FinishedAt.Before(cell.StartedAt) {
			return nil, errors.Errorf("evaluation cell line %d has invalid timestamps", index+1)
		}
		if cell.SnapshotDigest != expectedCell.snapshot {
			return nil, errors.Errorf("evaluation cell line %d has invalid snapshot", index+1)
		}
		nativeDirectory, err := reader.Path(filepath.Join("native", expectedCell.arm, expectedCell.caseValue.ID, formatRepeat(expectedCell.repeat)))
		if err != nil {
			return nil, err
		}
		if err := validateOutcome(reader.Dir(), nativeDirectory, &cell.Outcome); err != nil {
			return nil, errors.Wrapf(err, "validate evaluation cell line %d outcome", index+1)
		}
		cells = append(cells, cell)
	}
	return cells, nil
}
