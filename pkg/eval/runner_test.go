package eval

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/go-go-golems/ragopt/pkg/candidate"
	"github.com/go-go-golems/ragopt/pkg/runstore"
)

func TestInterruptedRunResumesToUninterruptedCanonicalResult(t *testing.T) {
	fixture := newEvaluationFixture(t)
	ctx, cancel := context.WithCancel(t.Context())
	control := &scriptControl{cancel: cancel, cancelOnCall: 4}
	interruptedRequest := fixture.request(t.TempDir(),
		&scriptedArm{name: "incumbent", control: control},
		&scriptedArm{name: "challenger", control: control},
	)
	interrupted, err := Run(ctx, interruptedRequest)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected interruption, got %v", err)
	}
	if interrupted.Completed != 3 || interrupted.Expected != 16 {
		t.Fatalf("interrupted counts: %#v", interrupted)
	}
	reader, err := runstore.Open(interrupted.RunDirectory)
	mustNoError(t, err)
	if reader.Status().State != runstore.StateActive {
		t.Fatalf("interrupted run state: %q", reader.Status().State)
	}
	if got := len(readCells(t, interrupted.RunDirectory)); got != 3 {
		t.Fatalf("durable interrupted cells: got %d", got)
	}
	appendFile(t, filepath.Join(interrupted.RunDirectory, "results", "cells.jsonl"), []byte(`{"truncated"`))

	resumeControl := &scriptControl{}
	resumeRequest := fixture.request(filepath.Dir(interrupted.RunDirectory),
		&scriptedArm{name: "incumbent", control: resumeControl},
		&scriptedArm{name: "challenger", control: resumeControl},
	)
	resumed, err := Resume(t.Context(), interrupted.RunDirectory, resumeRequest)
	mustNoError(t, err)
	if resumed.Completed != 16 || resumed.Failures != 0 || !resumed.Resumed {
		t.Fatalf("resumed result: %#v", resumed)
	}
	if resumeControl.callCount() != 13 {
		t.Fatalf("resume reran completed cells: calls=%d", resumeControl.callCount())
	}
	if first := resumeControl.callsSnapshot()[0]; first != "challenger/cmp-1/1" {
		t.Fatalf("resume started at wrong cell: %s", first)
	}

	uninterruptedControl := &scriptControl{}
	uninterruptedRequest := fixture.request(t.TempDir(),
		&scriptedArm{name: "incumbent", control: uninterruptedControl},
		&scriptedArm{name: "challenger", control: uninterruptedControl},
	)
	uninterrupted, err := Run(t.Context(), uninterruptedRequest)
	mustNoError(t, err)
	if uninterrupted.Completed != 16 || uninterruptedControl.callCount() != 16 {
		t.Fatalf("uninterrupted result: %#v calls=%d", uninterrupted, uninterruptedControl.callCount())
	}
	loaded, err := LoadArtifactRun(t.Context(), uninterrupted.RunDirectory)
	mustNoError(t, err)
	if len(loaded.Cells) != 16 || loaded.Config.ChangedAsset != fixture.candidate.Mutation.AssetName || loaded.PolicyPath == "" {
		t.Fatalf("strictly loaded artifact run: cells=%d config=%#v policy=%q", len(loaded.Cells), loaded.Config, loaded.PolicyPath)
	}

	resumedCells := canonicalCells(t, readCells(t, resumed.RunDirectory))
	uninterruptedCells := canonicalCells(t, readCells(t, uninterrupted.RunDirectory))
	if !reflect.DeepEqual(resumedCells, uninterruptedCells) {
		t.Fatalf("resumed and uninterrupted cells differ\nresumed: %#v\nuninterrupted: %#v", resumedCells, uninterruptedCells)
	}
	for _, call := range append(control.viewsSnapshot(), resumeControl.viewsSnapshot()...) {
		if !strings.HasPrefix(call.assetPath, interrupted.RunDirectory+string(filepath.Separator)) {
			t.Fatalf("arm received non-copied asset path: %s", call.assetPath)
		}
	}

	// Original bundle edits cannot alter the completed run's copied inputs.
	writeFile(t, filepath.Join(fixture.candidate.Root, "candidate", "prompt.md"), []byte("edited after evaluation\n"))
	_, err = runstore.Open(resumed.RunDirectory)
	mustNoError(t, err)
}

func TestArmErrorRecordsFailedCellAndContinues(t *testing.T) {
	fixture := newEvaluationFixture(t)
	control := &scriptControl{failKey: "challenger/cmp-2/0"}
	request := fixture.request(t.TempDir(),
		&scriptedArm{name: "incumbent", control: control},
		&scriptedArm{name: "challenger", control: control},
	)
	request.Repeats = 1
	result, err := Run(t.Context(), request)
	mustNoError(t, err)
	if result.Completed != 8 || result.Failures != 1 {
		t.Fatalf("failed-cell accounting: %#v", result)
	}
	cells := readCells(t, result.RunDirectory)
	failed := 0
	for _, cell := range cells {
		if cell.Outcome.Failure == nil {
			continue
		}
		failed++
		if cell.Outcome.Failure.Class != "arm_error" || cell.Outcome.Completed || cell.Outcome.ContractValid {
			t.Fatalf("invalid failed outcome: %#v", cell.Outcome)
		}
		assertArtifact(t, result.RunDirectory, cell.Outcome.NativeArtifact)
	}
	if failed != 1 {
		t.Fatalf("failed cells: got %d", failed)
	}
}

func TestInvalidOutcomeIsCustodyFailure(t *testing.T) {
	fixture := newEvaluationFixture(t)
	control := &scriptControl{invalidArtifactKey: "incumbent/cmp-1/0"}
	request := fixture.request(t.TempDir(),
		&scriptedArm{name: "incumbent", control: control},
		&scriptedArm{name: "challenger", control: control},
	)
	result, err := Run(t.Context(), request)
	if err == nil || !strings.Contains(err.Error(), "outside its assigned cell directory") {
		t.Fatalf("expected native custody error, got %v", err)
	}
	reader, openErr := runstore.Open(result.RunDirectory)
	mustNoError(t, openErr)
	if reader.Status().State != runstore.StateFailed {
		t.Fatalf("custody failure state: %q", reader.Status().State)
	}
}

func TestNonFiniteMetricIsCustodyFailure(t *testing.T) {
	fixture := newEvaluationFixture(t)
	control := &scriptControl{nanKey: "challenger/cmp-1/0"}
	request := fixture.request(t.TempDir(),
		&scriptedArm{name: "incumbent", control: control},
		&scriptedArm{name: "challenger", control: control},
	)
	result, err := Run(t.Context(), request)
	if err == nil || !strings.Contains(err.Error(), "metric \"quality\" is not finite") {
		t.Fatalf("expected finite metric error, got %v", err)
	}
	reader, openErr := runstore.Open(result.RunDirectory)
	mustNoError(t, openErr)
	if reader.Status().State != runstore.StateFailed {
		t.Fatalf("metric custody failure state: %q", reader.Status().State)
	}
}

func TestResumeRejectsIdentityMismatchWithoutMutatingActiveRun(t *testing.T) {
	fixture := newEvaluationFixture(t)
	ctx, cancel := context.WithCancel(t.Context())
	control := &scriptControl{cancel: cancel, cancelOnCall: 2}
	request := fixture.request(t.TempDir(),
		&scriptedArm{name: "incumbent", control: control},
		&scriptedArm{name: "challenger", control: control},
	)
	result, err := Run(ctx, request)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected interruption, got %v", err)
	}
	request.Repeats++
	_, err = Resume(t.Context(), result.RunDirectory, request)
	if err == nil || !strings.Contains(err.Error(), "resume config digest mismatch") {
		t.Fatalf("expected resume identity mismatch, got %v", err)
	}
	reader, openErr := runstore.Open(result.RunDirectory)
	mustNoError(t, openErr)
	if reader.Status().State != runstore.StateActive {
		t.Fatalf("identity mismatch mutated active run: %q", reader.Status().State)
	}
}

func TestResumeRejectsMissingBoundInput(t *testing.T) {
	fixture := newEvaluationFixture(t)
	ctx, cancel := context.WithCancel(t.Context())
	control := &scriptControl{cancel: cancel, cancelOnCall: 2}
	request := fixture.request(t.TempDir(),
		&scriptedArm{name: "incumbent", control: control},
		&scriptedArm{name: "challenger", control: control},
	)
	result, err := Run(ctx, request)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected interruption, got %v", err)
	}
	reader, err := runstore.Open(result.RunDirectory)
	mustNoError(t, err)
	inputs := reader.Inputs()
	filtered := make([]runstore.InputRef, 0, len(inputs)-1)
	for _, input := range inputs {
		if input.Role == roleCandidateManifest {
			mustNoError(t, os.Remove(filepath.Join(result.RunDirectory, input.CopiedPath)))
			continue
		}
		filtered = append(filtered, input)
	}
	manifestData, err := json.MarshalIndent(filtered, "", "  ")
	mustNoError(t, err)
	writeFile(t, filepath.Join(result.RunDirectory, "inputs", "manifest.json"), append(manifestData, '\n'))

	resumeControl := &scriptControl{}
	request.Incumbent = &scriptedArm{name: "incumbent", control: resumeControl}
	request.Challenger = &scriptedArm{name: "challenger", control: resumeControl}
	_, err = Resume(t.Context(), result.RunDirectory, request)
	if err == nil || !strings.Contains(err.Error(), "bound input count mismatch") {
		t.Fatalf("expected complete input-set rejection, got %v", err)
	}
	reader, openErr := runstore.Open(result.RunDirectory)
	mustNoError(t, openErr)
	if reader.Status().State != runstore.StateFailed {
		t.Fatalf("missing-input resume state: %q", reader.Status().State)
	}
}

func TestResumeRejectsMalformedMiddleAndDuplicateCells(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(t *testing.T, runDirectory string)
		want   string
	}{
		{
			name: "malformed middle",
			mutate: func(t *testing.T, runDirectory string) {
				appendFile(t, filepath.Join(runDirectory, "results", "cells.jsonl"), []byte("not-json\n"))
			},
			want: "decode result cell line",
		},
		{
			name: "duplicate",
			mutate: func(t *testing.T, runDirectory string) {
				path := filepath.Join(runDirectory, "results", "cells.jsonl")
				data, err := os.ReadFile(path)
				mustNoError(t, err)
				appendFile(t, path, data)
			},
			want: "duplicate result cell key",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newEvaluationFixture(t)
			ctx, cancel := context.WithCancel(t.Context())
			control := &scriptControl{cancel: cancel, cancelOnCall: 2}
			request := fixture.request(t.TempDir(),
				&scriptedArm{name: "incumbent", control: control},
				&scriptedArm{name: "challenger", control: control},
			)
			result, err := Run(ctx, request)
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("expected interruption, got %v", err)
			}
			test.mutate(t, result.RunDirectory)
			resumeControl := &scriptControl{}
			request.Incumbent = &scriptedArm{name: "incumbent", control: resumeControl}
			request.Challenger = &scriptedArm{name: "challenger", control: resumeControl}
			_, err = Resume(t.Context(), result.RunDirectory, request)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected %q, got %v", test.want, err)
			}
			reader, openErr := runstore.Open(result.RunDirectory)
			mustNoError(t, openErr)
			if reader.Status().State != runstore.StateFailed {
				t.Fatalf("corrupt resume state: %q", reader.Status().State)
			}
		})
	}
}

func TestRunRequestValidation(t *testing.T) {
	fixture := newEvaluationFixture(t)
	valid := fixture.request(t.TempDir(), &scriptedArm{name: "same", control: &scriptControl{}}, &scriptedArm{name: "same", control: &scriptControl{}})
	_, err := Run(t.Context(), valid)
	if err == nil || !strings.Contains(err.Error(), "arm names must be unique") {
		t.Fatalf("expected arm identity rejection, got %v", err)
	}
	valid = fixture.request(t.TempDir(), &scriptedArm{name: "incumbent", control: &scriptControl{}}, &scriptedArm{name: "challenger", control: &scriptControl{}})
	valid.Repeats = 0
	_, err = Run(t.Context(), valid)
	if err == nil || !strings.Contains(err.Error(), "repeats") {
		t.Fatalf("expected repeat rejection, got %v", err)
	}
}

type evaluationFixture struct {
	suite     *SuiteDocument
	policy    string
	candidate *candidate.Candidate
}

func newEvaluationFixture(t *testing.T) *evaluationFixture {
	t.Helper()
	root := t.TempDir()
	suitePath := filepath.Join(root, "suite.json")
	writeFile(t, suitePath, []byte(`{
  "api_version":"ragopt-suite/v1",
  "name":"paired-fixture",
  "cases":[
    {"id":"cmp-1","groups":["comparison"],"input":{"question":"compare one"}},
    {"id":"cmp-2","groups":["comparison"],"input":{"question":"compare two"}},
    {"id":"single-1","groups":["single"],"input":{"question":"single"}},
    {"id":"abstain-1","groups":["abstention"],"input":{"question":"unknown"}}
  ]
}`))
	suite, err := LoadSuite(t.Context(), suitePath)
	mustNoError(t, err)
	policyPath := filepath.Join(root, "policy.json")
	writeFile(t, policyPath, []byte(`{"api_version":"fixture-policy/v1","target":"quality"}`))
	candidateValue := writeCandidateFixture(t, filepath.Join(root, "bundle"))
	return &evaluationFixture{suite: suite, policy: policyPath, candidate: candidateValue}
}

func (fixture *evaluationFixture) request(root string, incumbent, challenger Arm) RunRequest {
	return RunRequest{
		RunRoot:     root,
		Name:        "paired fixture",
		Description: "deterministic scripted evaluation",
		Suite:       fixture.suite,
		PolicyPath:  fixture.policy,
		Candidate:   fixture.candidate,
		Incumbent:   incumbent,
		Challenger:  challenger,
		Repeats:     2,
	}
}

type scriptControl struct {
	mu                 sync.Mutex
	calls              []string
	views              []viewObservation
	cancel             context.CancelFunc
	cancelOnCall       int
	failKey            string
	invalidArtifactKey string
	nanKey             string
}

type viewObservation struct {
	role      string
	assetPath string
}

func (control *scriptControl) begin(key string, view CandidateView) (int, bool, bool, bool) {
	control.mu.Lock()
	defer control.mu.Unlock()
	control.calls = append(control.calls, key)
	asset := view.Assets["prompt"]
	control.views = append(control.views, viewObservation{role: view.Role, assetPath: asset.Path})
	callNumber := len(control.calls)
	return callNumber, key == control.failKey, key == control.invalidArtifactKey, key == control.nanKey
}

func (control *scriptControl) callCount() int {
	control.mu.Lock()
	defer control.mu.Unlock()
	return len(control.calls)
}

func (control *scriptControl) callsSnapshot() []string {
	control.mu.Lock()
	defer control.mu.Unlock()
	return append([]string(nil), control.calls...)
}

func (control *scriptControl) viewsSnapshot() []viewObservation {
	control.mu.Lock()
	defer control.mu.Unlock()
	return append([]viewObservation(nil), control.views...)
}

type scriptedArm struct {
	name    string
	control *scriptControl
}

var _ Arm = (*scriptedArm)(nil)

func (arm *scriptedArm) Name() string { return arm.name }

func (arm *scriptedArm) Run(ctx context.Context, request Request) (Outcome, error) {
	key := fmt.Sprintf("%s/%s/%d", arm.name, request.Case.ID, request.RepeatIndex)
	callNumber, fail, invalid, nanMetric := arm.control.begin(key, request.Candidate)
	if arm.control.cancelOnCall == callNumber {
		arm.control.cancel()
		return Outcome{}, ctx.Err()
	}
	if fail {
		return Outcome{}, errors.New("scripted provider failure")
	}
	asset := request.Candidate.Assets["prompt"]
	assetBytes, err := os.ReadFile(asset.Path)
	if err != nil {
		return Outcome{}, err
	}
	payload := map[string]any{
		"arm":         arm.name,
		"case_id":     request.Case.ID,
		"repeat":      request.RepeatIndex,
		"role":        request.Candidate.Role,
		"prompt_text": string(assetBytes),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return Outcome{}, err
	}
	artifactPath := filepath.Join(request.NativeDirectory, "result.json")
	if err := os.WriteFile(artifactPath, data, 0o600); err != nil {
		return Outcome{}, err
	}
	if invalid {
		artifactPath = filepath.Join(request.RunDirectory, "config.json")
	}
	relative, err := filepath.Rel(request.RunDirectory, artifactPath)
	if err != nil {
		return Outcome{}, err
	}
	completed := true
	contractValid := true
	abstained := request.Case.ID == "abstain-1"
	quality := 0.5
	toolCalls := 1
	if request.Candidate.Role == "candidate" {
		quality = 0.8
		toolCalls = 2
	}
	if nanMetric {
		quality = math.NaN()
	}
	return Outcome{
		Completed:      completed,
		ContractValid:  contractValid,
		Abstained:      abstained,
		Metrics:        map[string]float64{"quality": quality},
		ProviderCalls:  1,
		ToolCalls:      toolCalls,
		InputTokens:    10,
		OutputTokens:   5,
		NativeArtifact: ArtifactRef{Path: relative},
	}, nil
}

func writeCandidateFixture(t *testing.T, root string) *candidate.Candidate {
	t.Helper()
	parentPrompt := []byte("parent prompt\n")
	childPrompt := []byte("candidate prompt\n")
	judge := []byte("judge prompt\n")
	writeFile(t, filepath.Join(root, "parent", "prompt.md"), parentPrompt)
	writeFile(t, filepath.Join(root, "candidate", "prompt.md"), childPrompt)
	writeFile(t, filepath.Join(root, "shared", "judge.md"), judge)
	locked := candidate.AssetRef{Name: "judge", MediaType: "text/markdown", Path: "shared/judge.md", SHA256: byteDigest(judge), SizeBytes: int64(len(judge))}
	parent := candidate.Snapshot{
		APIVersion:   candidate.SnapshotAPIVersion,
		System:       "fixture-rag",
		LockedAssets: []candidate.AssetRef{locked},
		MutableAssets: []candidate.AssetRef{{
			Name: "prompt", MediaType: "text/markdown", Path: "parent/prompt.md", SHA256: byteDigest(parentPrompt), SizeBytes: int64(len(parentPrompt)),
		}},
		Dimensions: map[string]string{"model": "fixture-model", "evaluator": "fixture-eval-v1"},
	}
	child := candidate.Snapshot{
		APIVersion:   candidate.SnapshotAPIVersion,
		System:       "fixture-rag",
		LockedAssets: []candidate.AssetRef{locked},
		MutableAssets: []candidate.AssetRef{{
			Name: "prompt", MediaType: "text/markdown", Path: "candidate/prompt.md", SHA256: byteDigest(childPrompt), SizeBytes: int64(len(childPrompt)),
		}},
		Dimensions: map[string]string{"model": "fixture-model", "evaluator": "fixture-eval-v1"},
	}
	var err error
	parent.SnapshotID, err = candidate.DigestSnapshot(parent)
	mustNoError(t, err)
	child.SnapshotID, err = candidate.DigestSnapshot(child)
	mustNoError(t, err)
	writeYAML(t, filepath.Join(root, "parent", "snapshot.yaml"), parent)
	writeYAML(t, filepath.Join(root, "candidate", "snapshot.yaml"), child)
	manifest := candidate.CandidateManifest{
		APIVersion:        candidate.CandidateAPIVersion,
		CandidateID:       "prompt-001",
		ParentSnapshot:    "parent/snapshot.yaml",
		CandidateSnapshot: "candidate/snapshot.yaml",
		Proposer:          candidate.Proposer{Kind: "human", Identity: "fixture"},
		Mutation: candidate.MutationDeclaration{
			Asset:      "prompt",
			Hypothesis: "The candidate prompt improves scripted quality.",
			ExpectedImprovement: candidate.ExpectedImprovement{
				Metric: "quality",
				Groups: []string{"comparison"},
			},
			RegressionRisks: []string{"one extra tool call"},
		},
	}
	writeYAML(t, filepath.Join(root, "candidate.yaml"), manifest)
	loaded, err := candidate.LoadCandidate(t.Context(), root, "candidate.yaml")
	mustNoError(t, err)
	return loaded
}

func writeYAML(t *testing.T, path string, value any) {
	t.Helper()
	data, err := yaml.Marshal(value)
	mustNoError(t, err)
	writeFile(t, path, data)
}

func byteDigest(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func appendFile(t *testing.T, path string, data []byte) {
	t.Helper()
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o600)
	mustNoError(t, err)
	_, err = file.Write(data)
	mustNoError(t, err)
	mustNoError(t, file.Sync())
	mustNoError(t, file.Close())
}

func readCells(t *testing.T, runDirectory string) []Cell {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(runDirectory, "results", "cells.jsonl"))
	mustNoError(t, err)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	result := make([]Cell, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		var cell Cell
		mustNoError(t, json.Unmarshal([]byte(line), &cell))
		result = append(result, cell)
	}
	return result
}

type canonicalCell struct {
	CaseID         string
	RepeatIndex    int
	Arm            string
	CandidateID    string
	SnapshotDigest string
	SuiteDigest    string
	PolicyDigest   string
	Outcome        Outcome
}

func canonicalCells(t *testing.T, cells []Cell) []canonicalCell {
	t.Helper()
	result := make([]canonicalCell, 0, len(cells))
	for _, cell := range cells {
		outcome := cell.Outcome
		outcome.Duration = 0
		outcome.NativeArtifact.Path = filepath.ToSlash(outcome.NativeArtifact.Path)
		result = append(result, canonicalCell{
			CaseID: cell.CaseID, RepeatIndex: cell.RepeatIndex, Arm: cell.Arm,
			CandidateID: cell.CandidateID, SnapshotDigest: cell.SnapshotDigest,
			SuiteDigest: cell.SuiteDigest, PolicyDigest: cell.PolicyDigest, Outcome: outcome,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		left := fmt.Sprintf("%s/%04d/%s", result[i].CaseID, result[i].RepeatIndex, result[i].Arm)
		right := fmt.Sprintf("%s/%04d/%s", result[j].CaseID, result[j].RepeatIndex, result[j].Arm)
		return left < right
	})
	return result
}

func assertArtifact(t *testing.T, runDirectory string, artifact ArtifactRef) {
	t.Helper()
	if artifact.Path == "" {
		return
	}
	data, err := os.ReadFile(filepath.Join(runDirectory, artifact.Path))
	mustNoError(t, err)
	if byteDigest(data) != artifact.SHA256 || int64(len(data)) != artifact.SizeBytes {
		t.Fatalf("artifact identity mismatch: %#v", artifact)
	}
}

func mustNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
