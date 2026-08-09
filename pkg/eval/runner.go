package eval

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/go-go-golems/ragopt/pkg/candidate"
	"github.com/go-go-golems/ragopt/pkg/runstore"
)

const RunAPIVersion = "ragopt-eval-run/v1"

// RunConfig is the complete semantic identity required to resume and compare a
// paired evaluation run.
type RunConfig struct {
	APIVersion        string                        `json:"api_version"`
	SuiteDigest       string                        `json:"suite_digest"`
	PolicyDigest      string                        `json:"policy_digest"`
	CandidateID       string                        `json:"candidate_id"`
	CandidateDigest   string                        `json:"candidate_digest"`
	ParentSnapshot    string                        `json:"parent_snapshot"`
	ChildSnapshot     string                        `json:"child_snapshot"`
	IncumbentArm      string                        `json:"incumbent_arm"`
	ChallengerArm     string                        `json:"challenger_arm"`
	Repeats           int                           `json:"repeats"`
	InputDigests      map[string]string             `json:"input_digests"`
	Mutation          candidate.MutationDeclaration `json:"mutation"`
	ChangedAsset      string                        `json:"changed_asset"`
	ParentAssetDigest string                        `json:"parent_asset_digest"`
	ChildAssetDigest  string                        `json:"child_asset_digest"`
}

type preparedRequest struct {
	request    RunRequest
	suite      *SuiteDocument
	candidate  *candidate.Candidate
	policyPath string
	config     RunConfig
}

type scheduledCell struct {
	caseValue Case
	repeat    int
	arm       Arm
	armName   string
	view      CandidateView
}

// Run creates and executes one new deterministic paired run.
func Run(ctx context.Context, request RunRequest) (*RunResult, error) {
	prepared, err := prepareRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	run, err := runstore.Create(ctx, runstore.Options{
		Root:        prepared.request.RunRoot,
		Name:        prepared.request.Name,
		Description: prepared.request.Description,
		Dimensions: map[string]string{
			"suite_digest":     prepared.config.SuiteDigest,
			"policy_digest":    prepared.config.PolicyDigest,
			"candidate_digest": prepared.config.CandidateDigest,
		},
	}, prepared.config)
	if err != nil {
		return nil, errors.Wrap(err, "create paired evaluation run")
	}
	result := newRunResult(run, prepared, false)
	incumbentView, challengerView, err := bindInputs(ctx, run, prepared)
	if err != nil {
		return failRun(ctx, run, result, errors.Wrap(err, "bind immutable run inputs"))
	}
	return execute(ctx, run, prepared, incumbentView, challengerView, nil, result)
}

// Resume explicitly reopens an active run after validating the complete run
// config and copied inputs, then executes only missing cells.
func Resume(ctx context.Context, runDirectory string, request RunRequest) (*RunResult, error) {
	prepared, err := prepareRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	run, err := runstore.Resume(ctx, runDirectory, prepared.config)
	if err != nil {
		return nil, errors.Wrap(err, "resume paired evaluation run")
	}
	result := newRunResult(run, prepared, true)
	if err := verifyBoundInputs(run, prepared.config.InputDigests); err != nil {
		return failRun(ctx, run, result, errors.Wrap(err, "validate resumed immutable inputs"))
	}
	incumbentView, challengerView, err := viewsFromInputs(run, prepared.candidate)
	if err != nil {
		return failRun(ctx, run, result, errors.Wrap(err, "restore immutable candidate views"))
	}
	completed, err := loadCompletedCells(run, prepared, incumbentView, challengerView)
	if err != nil {
		return failRun(ctx, run, result, errors.Wrap(err, "load completed result cells"))
	}
	for _, cell := range completed {
		result.Completed++
		if cell.Outcome.Failure != nil {
			result.Failures++
		}
	}
	return execute(ctx, run, prepared, incumbentView, challengerView, completed, result)
}

func prepareRequest(ctx context.Context, request RunRequest) (*preparedRequest, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(request.RunRoot) == "" || strings.TrimSpace(request.Name) == "" {
		return nil, errors.New("run root and name are required")
	}
	absoluteRunRoot, err := filepath.Abs(request.RunRoot)
	if err != nil {
		return nil, errors.Wrap(err, "resolve run root")
	}
	request.RunRoot = absoluteRunRoot
	if request.Repeats <= 0 {
		return nil, errors.New("repeats must be greater than zero")
	}
	if request.Repeats > 10000 {
		return nil, errors.New("repeats exceeds the v1 limit of 10000")
	}
	if err := validateSuiteDocument(request.Suite); err != nil {
		return nil, errors.Wrap(err, "validate suite document")
	}
	reloadedSuite, err := LoadSuite(ctx, request.Suite.SourcePath)
	if err != nil {
		return nil, errors.Wrap(err, "reload suite source")
	}
	if reloadedSuite.Digest != request.Suite.Digest {
		return nil, errors.Errorf("suite source changed after load: expected=%s actual=%s", request.Suite.Digest, reloadedSuite.Digest)
	}
	if request.Candidate == nil {
		return nil, errors.New("candidate is required")
	}
	reloadedCandidate, err := candidate.LoadCandidate(ctx, request.Candidate.Root, request.Candidate.ManifestPath)
	if err != nil {
		return nil, errors.Wrap(err, "reload candidate bundle")
	}
	if reloadedCandidate.Digest != request.Candidate.Digest {
		return nil, errors.Errorf("candidate changed after load: expected=%s actual=%s", request.Candidate.Digest, reloadedCandidate.Digest)
	}
	if request.Incumbent == nil || request.Challenger == nil {
		return nil, errors.New("incumbent and challenger arms are required")
	}
	incumbentName := request.Incumbent.Name()
	challengerName := request.Challenger.Name()
	if !identifierPattern.MatchString(incumbentName) || !identifierPattern.MatchString(challengerName) {
		return nil, errors.New("arm names must be bounded identifiers")
	}
	if incumbentName == challengerName {
		return nil, errors.Errorf("arm names must be unique, both are %q", incumbentName)
	}
	policyPath, policyDigest, err := loadPolicyIdentity(request.PolicyPath)
	if err != nil {
		return nil, err
	}
	inputDigests, err := expectedInputDigests(reloadedSuite, policyDigest, reloadedCandidate)
	if err != nil {
		return nil, err
	}
	config := RunConfig{
		APIVersion:        RunAPIVersion,
		SuiteDigest:       reloadedSuite.Digest,
		PolicyDigest:      policyDigest,
		CandidateID:       reloadedCandidate.Manifest.CandidateID,
		CandidateDigest:   reloadedCandidate.Digest,
		ParentSnapshot:    reloadedCandidate.Parent.SnapshotID,
		ChildSnapshot:     reloadedCandidate.Child.SnapshotID,
		IncumbentArm:      incumbentName,
		ChallengerArm:     challengerName,
		Repeats:           request.Repeats,
		InputDigests:      inputDigests,
		Mutation:          reloadedCandidate.Manifest.Mutation,
		ChangedAsset:      reloadedCandidate.Mutation.AssetName,
		ParentAssetDigest: reloadedCandidate.Mutation.ParentDigest,
		ChildAssetDigest:  reloadedCandidate.Mutation.ChildDigest,
	}
	return &preparedRequest{
		request:    request,
		suite:      reloadedSuite,
		candidate:  reloadedCandidate,
		policyPath: policyPath,
		config:     config,
	}, nil
}

func loadPolicyIdentity(path string) (string, string, error) {
	if strings.TrimSpace(path) == "" {
		return "", "", errors.New("gate policy path is required")
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", "", errors.Wrap(err, "resolve gate policy path")
	}
	info, err := os.Stat(absolute)
	if err != nil {
		return "", "", errors.Wrap(err, "stat gate policy")
	}
	if !info.Mode().IsRegular() {
		return "", "", errors.New("gate policy path is not a regular file")
	}
	data, err := os.ReadFile(absolute)
	if err != nil {
		return "", "", errors.Wrap(err, "read gate policy")
	}
	return absolute, digestBytes(data), nil
}

func execute(
	ctx context.Context,
	run *runstore.Run,
	prepared *preparedRequest,
	incumbentView CandidateView,
	challengerView CandidateView,
	completed map[string]Cell,
	result *RunResult,
) (*RunResult, error) {
	if completed == nil {
		completed = map[string]Cell{}
	}
	schedule := buildSchedule(prepared, incumbentView, challengerView)
	for _, item := range schedule {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		key := expectedCellKey(prepared, item)
		if _, exists := completed[key]; exists {
			continue
		}
		cell, err := executeCell(ctx, run, prepared, item)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return result, err
			}
			return failRun(ctx, run, result, err)
		}
		if err := run.AppendJSONL(ctx, "results/cells.jsonl", cell); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return result, ctxErr
			}
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return result, err
			}
			return failRun(ctx, run, result, errors.Wrap(err, "append result cell"))
		}
		completed[key] = cell
		result.Completed++
		if cell.Outcome.Failure != nil {
			result.Failures++
		}
	}
	if result.Completed != result.Expected {
		return failRun(ctx, run, result, errors.Errorf("result cell count mismatch: completed=%d expected=%d", result.Completed, result.Expected))
	}
	if err := run.Complete(ctx, runstore.Summary{
		Message: "paired evaluation complete",
		Metrics: map[string]any{
			"expected_cells":  result.Expected,
			"completed_cells": result.Completed,
			"failed_cells":    result.Failures,
		},
	}); err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return result, ctxErr
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return result, err
		}
		return failRun(ctx, run, result, errors.Wrap(err, "complete paired evaluation run"))
	}
	return result, nil
}

func buildSchedule(prepared *preparedRequest, incumbentView, challengerView CandidateView) []scheduledCell {
	result := make([]scheduledCell, 0, len(prepared.suite.Suite.Cases)*prepared.config.Repeats*2)
	for _, caseValue := range prepared.suite.Suite.Cases {
		for repeat := 0; repeat < prepared.config.Repeats; repeat++ {
			result = append(result,
				scheduledCell{caseValue: cloneCase(caseValue), repeat: repeat, arm: prepared.request.Incumbent, armName: prepared.config.IncumbentArm, view: incumbentView},
				scheduledCell{caseValue: cloneCase(caseValue), repeat: repeat, arm: prepared.request.Challenger, armName: prepared.config.ChallengerArm, view: challengerView},
			)
		}
	}
	return result
}

func executeCell(ctx context.Context, run *runstore.Run, prepared *preparedRequest, item scheduledCell) (Cell, error) {
	nativeRelative := nativeCellPath(item.armName, item.caseValue.ID, item.repeat)
	nativeDirectory, err := run.Path(nativeRelative)
	if err != nil {
		return Cell{}, errors.Wrap(err, "resolve native artifact directory")
	}
	if err := os.RemoveAll(nativeDirectory); err != nil {
		return Cell{}, errors.Wrap(err, "clear uncommitted native artifact directory")
	}
	if err := os.MkdirAll(nativeDirectory, 0o700); err != nil {
		return Cell{}, errors.Wrap(err, "create native artifact directory")
	}
	protected, err := snapshotRunEvidence(run.Dir(), nativeDirectory)
	if err != nil {
		return Cell{}, errors.Wrap(err, "snapshot run evidence before arm")
	}
	startedAt := time.Now().UTC()
	outcome, armErr := item.arm.Run(ctx, Request{
		RunDirectory:    run.Dir(),
		NativeDirectory: nativeDirectory,
		Case:            cloneCase(item.caseValue),
		RepeatIndex:     item.repeat,
		Candidate:       cloneView(item.view),
	})
	finishedAt := time.Now().UTC()
	if err := verifyRunEvidence(run.Dir(), nativeDirectory, protected); err != nil {
		return Cell{}, errors.Wrap(err, "arm mutated run-owned evidence")
	}
	if armErr != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return Cell{}, ctxErr
		}
		if errors.Is(armErr, context.Canceled) || errors.Is(armErr, context.DeadlineExceeded) {
			return Cell{}, armErr
		}
		outcome, err = recordArmFailure(ctx, run, nativeRelative, armErr, finishedAt.Sub(startedAt))
		if err != nil {
			return Cell{}, err
		}
		if err := validateOutcome(run.Dir(), nativeDirectory, &outcome); err != nil {
			return Cell{}, errors.Wrap(err, "validate recorded arm failure")
		}
	} else {
		outcome.Duration = finishedAt.Sub(startedAt)
		if err := validateOutcome(run.Dir(), nativeDirectory, &outcome); err != nil {
			return Cell{}, errors.Wrap(err, "validate arm outcome")
		}
	}
	return Cell{
		APIVersion:     CellAPIVersion,
		RunID:          run.Manifest().RunID,
		CaseID:         item.caseValue.ID,
		RepeatIndex:    item.repeat,
		Arm:            item.armName,
		CandidateID:    prepared.config.CandidateID,
		SnapshotDigest: item.view.SnapshotDigest,
		SuiteDigest:    prepared.config.SuiteDigest,
		PolicyDigest:   prepared.config.PolicyDigest,
		StartedAt:      startedAt,
		FinishedAt:     finishedAt,
		Outcome:        outcome,
	}, nil
}

func nativeCellPath(armName, caseID string, repeat int) string {
	return filepath.Join(
		"native",
		nativeIdentityComponent("arm", armName),
		nativeIdentityComponent("case", caseID),
		fmt.Sprintf("%04d", repeat),
	)
}

func nativeIdentityComponent(kind, identity string) string {
	digest := sha256.Sum256([]byte(identity))
	return kind + "-" + hex.EncodeToString(digest[:])
}

func recordArmFailure(ctx context.Context, run *runstore.Run, nativeRelative string, armErr error, duration time.Duration) (Outcome, error) {
	relative := filepath.Join(nativeRelative, "arm-error.json")
	payload := map[string]any{
		"class":   "arm_error",
		"message": armErr.Error(),
	}
	if err := run.WriteJSON(ctx, relative, payload); err != nil {
		return Outcome{}, errors.Wrap(err, "record arm failure artifact")
	}
	path, err := run.Path(relative)
	if err != nil {
		return Outcome{}, err
	}
	artifact, err := identifyArtifact(run.Dir(), path)
	if err != nil {
		return Outcome{}, err
	}
	return Outcome{
		Completed:      false,
		ContractValid:  false,
		Failure:        &Failure{Class: "arm_error", Message: armErr.Error()},
		Duration:       duration,
		NativeArtifact: artifact,
	}, nil
}

func validateOutcome(runDirectory, nativeDirectory string, outcome *Outcome) error {
	if outcome == nil {
		return errors.New("outcome is nil")
	}
	if outcome.Failure == nil && !outcome.Completed {
		return errors.New("incomplete outcome requires a failure")
	}
	if outcome.Failure != nil {
		if !identifierPattern.MatchString(outcome.Failure.Class) || strings.TrimSpace(outcome.Failure.Message) == "" {
			return errors.New("outcome failure requires a valid class and message")
		}
		if outcome.Completed || outcome.ContractValid || outcome.Abstained {
			return errors.New("failed outcome cannot be completed, contract-valid, or abstained")
		}
	}
	if outcome.Abstained && !outcome.Completed {
		return errors.New("abstained outcome must be completed")
	}
	for name, value := range outcome.Metrics {
		if !identifierPattern.MatchString(name) {
			return errors.Errorf("invalid metric name %q", name)
		}
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return errors.Errorf("metric %q is not finite", name)
		}
	}
	if outcome.ProviderCalls < 0 || outcome.ToolCalls < 0 || outcome.InputTokens < 0 || outcome.OutputTokens < 0 || outcome.Duration < 0 {
		return errors.New("outcome counts and duration must be nonnegative")
	}
	if strings.TrimSpace(outcome.NativeArtifact.Path) == "" {
		return errors.New("outcome native artifact path is required")
	}
	path, err := resolveNativeArtifact(runDirectory, nativeDirectory, outcome.NativeArtifact.Path)
	if err != nil {
		return err
	}
	identified, err := identifyArtifact(runDirectory, path)
	if err != nil {
		return err
	}
	if outcome.NativeArtifact.SHA256 != "" && outcome.NativeArtifact.SHA256 != identified.SHA256 {
		return errors.Errorf("native artifact digest mismatch: outcome=%s actual=%s", outcome.NativeArtifact.SHA256, identified.SHA256)
	}
	if outcome.NativeArtifact.SizeBytes != 0 && outcome.NativeArtifact.SizeBytes != identified.SizeBytes {
		return errors.Errorf("native artifact size mismatch: outcome=%d actual=%d", outcome.NativeArtifact.SizeBytes, identified.SizeBytes)
	}
	outcome.NativeArtifact = identified
	return nil
}

func validateStoredOutcome(runDirectory, nativeDirectory string, outcome *Outcome) error {
	if outcome == nil {
		return errors.New("stored outcome is nil")
	}
	claimed := outcome.NativeArtifact
	if !strings.HasPrefix(claimed.SHA256, "sha256:") || len(claimed.SHA256) != len("sha256:")+64 {
		return errors.New("stored native artifact digest is required")
	}
	if claimed.SizeBytes < 0 {
		return errors.New("stored native artifact size is invalid")
	}
	if err := validateOutcome(runDirectory, nativeDirectory, outcome); err != nil {
		return err
	}
	if claimed.SHA256 != outcome.NativeArtifact.SHA256 || claimed.SizeBytes != outcome.NativeArtifact.SizeBytes {
		return errors.New("stored native artifact identity differs from file")
	}
	return nil
}

func resolveNativeArtifact(runDirectory, nativeDirectory, relative string) (string, error) {
	if filepath.IsAbs(relative) || filepath.Clean(relative) != relative {
		return "", errors.New("native artifact path must be canonical and run-relative")
	}
	path := filepath.Join(runDirectory, relative)
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", errors.Wrap(err, "resolve native artifact")
	}
	relativeToNative, err := filepath.Rel(nativeDirectory, resolved)
	if err != nil {
		return "", errors.Wrap(err, "compare native artifact directory")
	}
	if relativeToNative == ".." || strings.HasPrefix(relativeToNative, ".."+string(filepath.Separator)) {
		return "", errors.New("native artifact is outside its assigned cell directory")
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", errors.Wrap(err, "stat native artifact")
	}
	if !info.Mode().IsRegular() {
		return "", errors.New("native artifact is not a regular file")
	}
	if err := rejectExternalArtifactAlias(runDirectory, nativeDirectory, info); err != nil {
		return "", err
	}
	return resolved, nil
}

func rejectExternalArtifactAlias(runDirectory, nativeDirectory string, artifactInfo os.FileInfo) error {
	runDirectory = filepath.Clean(runDirectory)
	nativeDirectory = filepath.Clean(nativeDirectory)
	err := filepath.WalkDir(runDirectory, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if filepath.Clean(path) == nativeDirectory && entry.IsDir() {
			return filepath.SkipDir
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() && os.SameFile(artifactInfo, info) {
			relative, relErr := filepath.Rel(runDirectory, path)
			if relErr != nil {
				return relErr
			}
			return errors.Errorf("native artifact aliases protected run evidence %q", relative)
		}
		return nil
	})
	return errors.Wrap(err, "check native artifact aliases")
}

func identifyArtifact(runDirectory, path string) (ArtifactRef, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ArtifactRef{}, errors.Wrap(err, "read native artifact")
	}
	relative, err := filepath.Rel(runDirectory, path)
	if err != nil {
		return ArtifactRef{}, errors.Wrap(err, "make native artifact path relative")
	}
	return ArtifactRef{Path: relative, SHA256: digestBytes(data), SizeBytes: int64(len(data))}, nil
}

func newRunResult(run *runstore.Run, prepared *preparedRequest, resumed bool) *RunResult {
	return &RunResult{
		RunDirectory: run.Dir(),
		RunID:        run.Manifest().RunID,
		Expected:     len(prepared.suite.Suite.Cases) * prepared.config.Repeats * 2,
		Resumed:      resumed,
	}
}

func failRun(ctx context.Context, run *runstore.Run, result *RunResult, cause error) (*RunResult, error) {
	cleanupContext, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	if failErr := run.Fail(cleanupContext, cause); failErr != nil {
		return result, errors.Wrapf(cause, "also failed to mark run failed: %v", failErr)
	}
	return result, cause
}

func cloneView(view CandidateView) CandidateView {
	assets := make(map[string]ResolvedAsset, len(view.Assets))
	for name, asset := range view.Assets {
		assets[name] = asset
	}
	view.Assets = assets
	view.Dimensions = cloneStrings(view.Dimensions)
	return view
}

func expectedCellKey(prepared *preparedRequest, item scheduledCell) string {
	return strings.Join([]string{
		prepared.config.SuiteDigest,
		prepared.config.PolicyDigest,
		prepared.config.CandidateID,
		item.view.SnapshotDigest,
		item.caseValue.ID,
		fmt.Sprintf("%d", item.repeat),
		item.armName,
	}, "\x00")
}

func cellKey(cell Cell) string {
	return strings.Join([]string{
		cell.SuiteDigest,
		cell.PolicyDigest,
		cell.CandidateID,
		cell.SnapshotDigest,
		cell.CaseID,
		fmt.Sprintf("%d", cell.RepeatIndex),
		cell.Arm,
	}, "\x00")
}
