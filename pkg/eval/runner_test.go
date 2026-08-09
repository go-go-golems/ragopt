package eval

import (
	"bytes"
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

func TestDurableSnapshotReloadsAllMutableRunEvidence(t *testing.T) {
	fixture := newEvaluationFixture(t)
	control := &scriptControl{}
	request := fixture.request(t.TempDir(), &scriptedArm{name: "incumbent", control: control}, &scriptedArm{name: "challenger", control: control})
	request.Repeats = 1
	result, err := Run(t.Context(), request)
	mustNoError(t, err)
	loaded, err := LoadArtifactRun(t.Context(), result.RunDirectory)
	mustNoError(t, err)
	wantRunID := loaded.Manifest.RunID
	wantMetric := loaded.Cells[0].Outcome.Metrics["quality"]
	wantGroup := loaded.Suite.Suite.Cases[0].Groups[0]
	loaded.Manifest.RunID = "mutated"
	loaded.Status.State = runstore.StateFailed
	loaded.Cells[0].Outcome.Metrics["quality"] = -100
	loaded.Suite.Suite.Cases[0].Groups[0] = "mutated"
	durable, err := loaded.DurableSnapshot(t.Context())
	mustNoError(t, err)
	if durable.Manifest.RunID != wantRunID || durable.Status.State != runstore.StateComplete || durable.Cells[0].Outcome.Metrics["quality"] != wantMetric || durable.Suite.Suite.Cases[0].Groups[0] != wantGroup {
		t.Fatalf("durable snapshot trusted mutable fields: %#v", durable)
	}
}

func TestFinalAuditCancellationLeavesRunActive(t *testing.T) {
	run := mustCreateEvidenceRun(t)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	result := &RunResult{RunDirectory: run.Dir(), RunID: run.Manifest().RunID}
	_, err := finalizeRun(ctx, run, result)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
	reader, openErr := runstore.Open(run.Dir())
	mustNoError(t, openErr)
	if reader.Status().State != runstore.StateActive {
		t.Fatalf("canceled final audit made run terminal: %q", reader.Status().State)
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

func TestAssetRoleEncodingIsInjective(t *testing.T) {
	digest := byteDigest([]byte("same"))
	dotted := assetRole("parent", candidate.AssetRef{Name: "a.b", SHA256: digest}, false)
	escapedWord := assetRole("parent", candidate.AssetRef{Name: "a-dot-b", SHA256: digest}, false)
	if dotted == escapedWord {
		t.Fatalf("asset roles collide: %q", dotted)
	}
}

func TestArmCannotMutateBoundInputs(t *testing.T) {
	fixture := newEvaluationFixture(t)
	request := fixture.request(t.TempDir(), &mutatingInputArm{name: "incumbent"}, &scriptedArm{name: "challenger", control: &scriptControl{}})
	result, err := Run(t.Context(), request)
	if err == nil || !strings.Contains(err.Error(), "arm mutated run-owned evidence") {
		t.Fatalf("expected bound-input mutation rejection, got %v", err)
	}
	statusData, readErr := os.ReadFile(filepath.Join(result.RunDirectory, "status.json"))
	mustNoError(t, readErr)
	var status runstore.Status
	mustNoError(t, json.Unmarshal(statusData, &status))
	if status.State != runstore.StateFailed {
		t.Fatalf("input mutation did not fail run: %q", status.State)
	}
}

func TestArmCannotMutateOtherRunEvidence(t *testing.T) {
	fixture := newEvaluationFixture(t)
	control := &scriptControl{}
	arm := &mutatingRunArm{delegate: &scriptedArm{name: "incumbent", control: control}}
	request := fixture.request(t.TempDir(), arm, &scriptedArm{name: "challenger", control: control})
	result, err := Run(t.Context(), request)
	if err == nil || !strings.Contains(err.Error(), `run evidence "config.json" changed`) {
		t.Fatalf("expected run-evidence mutation rejection, got %v", err)
	}
	_, openErr := runstore.Open(result.RunDirectory)
	if openErr == nil || !strings.Contains(openErr.Error(), "config digest mismatch") {
		t.Fatalf("expected tampered config to remain detectable, got %v", openErr)
	}
}

func TestCancellationReturnsContextErrorFromGenericArmError(t *testing.T) {
	fixture := newEvaluationFixture(t)
	ctx, cancel := context.WithCancel(t.Context())
	request := fixture.request(t.TempDir(), &cancelingErrorArm{name: "incumbent", cancel: cancel}, &scriptedArm{name: "challenger", control: &scriptControl{}})
	result, err := Run(ctx, request)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation, got %v", err)
	}
	reader, openErr := runstore.Open(result.RunDirectory)
	mustNoError(t, openErr)
	if reader.Status().State != runstore.StateActive {
		t.Fatalf("canceled run became terminal: %q", reader.Status().State)
	}
}

func TestCancellationAfterSuccessfulArmLeavesRunResumable(t *testing.T) {
	fixture := newEvaluationFixture(t)
	ctx, cancel := context.WithCancel(t.Context())
	control := &scriptControl{}
	arm := &cancelingSuccessArm{delegate: &scriptedArm{name: "incumbent", control: control}, cancel: cancel}
	request := fixture.request(t.TempDir(), arm, &scriptedArm{name: "challenger", control: control})
	result, err := Run(ctx, request)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation, got %v", err)
	}
	reader, openErr := runstore.Open(result.RunDirectory)
	mustNoError(t, openErr)
	if reader.Status().State != runstore.StateActive {
		t.Fatalf("canceled run became terminal: %q", reader.Status().State)
	}
}

func TestNativeCellPathsAreCaseIndependent(t *testing.T) {
	first := nativeCellPath("Arm", "Case", 0)
	second := nativeCellPath("arm", "case", 0)
	if strings.EqualFold(first, second) {
		t.Fatalf("case-distinct coordinates collide: %q and %q", first, second)
	}
}

func TestNativeCellPathComponentsStayWithinFilesystemLimits(t *testing.T) {
	identity := strings.Repeat("a", 128)
	path := nativeCellPath(identity, identity, 9999)
	for _, component := range strings.Split(path, string(filepath.Separator)) {
		if len(component) > 255 {
			t.Fatalf("native path component has %d bytes: %q", len(component), component)
		}
	}
}

func TestSuiteBytesRemainBoundToLoadedSemantics(t *testing.T) {
	fixture := newEvaluationFixture(t)
	control := &scriptControl{}
	request := fixture.request(t.TempDir(), &scriptedArm{name: "incumbent", control: control}, &scriptedArm{name: "challenger", control: control})
	prepared, err := prepareRequest(t.Context(), request)
	mustNoError(t, err)
	data, err := os.ReadFile(prepared.suite.SourcePath)
	mustNoError(t, err)
	writeFile(t, prepared.suite.SourcePath, append(data, '\n'))
	run, err := runstore.Create(t.Context(), runstore.Options{Root: request.RunRoot, Name: request.Name}, prepared.config)
	mustNoError(t, err)
	_, _, err = bindInputs(t.Context(), run, prepared)
	if err == nil || !strings.Contains(err.Error(), "evaluation suite changed during binding") {
		t.Fatalf("expected suite byte drift rejection, got %v", err)
	}
}

func TestCandidateManifestBytesRemainBoundToLoadedSemantics(t *testing.T) {
	tests := []struct {
		name string
		path func(*candidate.Candidate) string
		want string
	}{
		{name: "candidate", path: func(value *candidate.Candidate) string { return value.ManifestPath }, want: "candidate manifest changed during binding"},
		{name: "parent snapshot", path: func(value *candidate.Candidate) string { return value.Manifest.ParentSnapshot }, want: "parent snapshot changed during binding"},
		{name: "candidate snapshot", path: func(value *candidate.Candidate) string { return value.Manifest.CandidateSnapshot }, want: "candidate snapshot changed during binding"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newEvaluationFixture(t)
			control := &scriptControl{}
			request := fixture.request(t.TempDir(), &scriptedArm{name: "incumbent", control: control}, &scriptedArm{name: "challenger", control: control})
			prepared, err := prepareRequest(t.Context(), request)
			mustNoError(t, err)
			manifestPath := filepath.Join(prepared.candidate.Root, test.path(prepared.candidate))
			data, err := os.ReadFile(manifestPath)
			mustNoError(t, err)
			writeFile(t, manifestPath, append(data, '\n'))
			run, err := runstore.Create(t.Context(), runstore.Options{Root: request.RunRoot, Name: request.Name}, prepared.config)
			mustNoError(t, err)
			_, _, err = bindInputs(t.Context(), run, prepared)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected %s, got %v", test.want, err)
			}
		})
	}
}

func TestNativeArtifactCannotHardLinkProtectedRunEvidence(t *testing.T) {
	fixture := newEvaluationFixture(t)
	control := &scriptControl{}
	request := fixture.request(t.TempDir(), &scriptedArm{name: "incumbent", control: control}, &hardLinkArtifactArm{name: "challenger"})
	_, err := Run(t.Context(), request)
	if err == nil || !strings.Contains(err.Error(), "aliases protected run evidence") {
		t.Fatalf("expected hard-link alias rejection, got %v", err)
	}
}

func TestNativeArtifactHardLinkedOutsideRunBecomesOwnedCopy(t *testing.T) {
	fixture := newEvaluationFixture(t)
	external := filepath.Join(t.TempDir(), "external.json")
	writeFile(t, external, []byte(`{"external":true}`))
	control := &scriptControl{}
	request := fixture.request(t.TempDir(), &scriptedArm{name: "incumbent", control: control}, &externalHardLinkArm{name: "challenger", source: external})
	request.Repeats = 1
	result, err := Run(t.Context(), request)
	mustNoError(t, err)
	cells := readCells(t, result.RunDirectory)
	var artifact ArtifactRef
	for _, cell := range cells {
		if cell.Arm == "challenger" {
			artifact = cell.Outcome.NativeArtifact
			break
		}
	}
	externalInfo, err := os.Stat(external)
	mustNoError(t, err)
	artifactInfo, err := os.Stat(filepath.Join(result.RunDirectory, artifact.Path))
	mustNoError(t, err)
	if os.SameFile(externalInfo, artifactInfo) {
		t.Fatal("native artifact retained external hard-link identity")
	}
	writeFile(t, external, []byte(`{"mutated":true}`))
	_, err = LoadArtifactRun(t.Context(), result.RunDirectory)
	mustNoError(t, err)
}

func TestEvidenceSnapshotDoesNotGrowWithCommittedResults(t *testing.T) {
	run := mustCreateEvidenceRun(t)
	before, err := snapshotRunEvidence(run)
	mustNoError(t, err)
	mustNoError(t, run.WriteBytes(t.Context(), "native/old.json", []byte(`{}`)))
	mustNoError(t, run.AppendJSONL(t.Context(), "results/cells.jsonl", map[string]string{"cell": "old"}))
	after, err := snapshotRunEvidence(run)
	mustNoError(t, err)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("bounded evidence snapshot grew with results: before=%v after=%v", before, after)
	}
}

func TestResumeClearsInterruptedNativeDirectory(t *testing.T) {
	fixture := newEvaluationFixture(t)
	ctx, cancel := context.WithCancel(t.Context())
	request := fixture.request(t.TempDir(), &partialCancelArm{name: "incumbent", cancel: cancel}, &scriptedArm{name: "challenger", control: &scriptControl{}})
	result, err := Run(ctx, request)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected interruption, got %v", err)
	}
	control := &scriptControl{}
	checker := &emptyDirectoryArm{delegate: &scriptedArm{name: "incumbent", control: control}}
	request.Incumbent = checker
	request.Challenger = &scriptedArm{name: "challenger", control: control}
	resumed, err := Resume(t.Context(), result.RunDirectory, request)
	if err != nil {
		t.Fatal(err)
	}
	if !checker.called || resumed.Completed != resumed.Expected {
		t.Fatalf("resume did not execute cleanly: called=%t result=%#v", checker.called, resumed)
	}
}

func TestLoadRejectsStoredCellWithoutArtifactIdentity(t *testing.T) {
	fixture := newEvaluationFixture(t)
	control := &scriptControl{}
	request := fixture.request(t.TempDir(), &scriptedArm{name: "incumbent", control: control}, &scriptedArm{name: "challenger", control: control})
	result, err := Run(t.Context(), request)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(result.RunDirectory, "results", "cells.jsonl")
	data, err := os.ReadFile(path)
	mustNoError(t, err)
	lines := bytes.Split(data, []byte{'\n'})
	var cell Cell
	mustNoError(t, json.Unmarshal(lines[0], &cell))
	cell.Outcome.NativeArtifact.SHA256 = ""
	cell.Outcome.NativeArtifact.SizeBytes = 0
	lines[0], err = json.Marshal(cell)
	mustNoError(t, err)
	writeFile(t, path, bytes.Join(lines, []byte{'\n'}))
	_, err = LoadArtifactRun(t.Context(), result.RunDirectory)
	if err == nil || !strings.Contains(err.Error(), "artifact digest is required") {
		t.Fatalf("expected missing stored identity rejection, got %v", err)
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

type mutatingInputArm struct{ name string }

func (arm *mutatingInputArm) Name() string { return arm.name }
func (arm *mutatingInputArm) Run(_ context.Context, request Request) (Outcome, error) {
	return Outcome{}, os.WriteFile(request.Candidate.Assets["prompt"].Path, []byte("tampered"), 0o600)
}

type mutatingRunArm struct{ delegate Arm }

func (arm *mutatingRunArm) Name() string { return arm.delegate.Name() }
func (arm *mutatingRunArm) Run(ctx context.Context, request Request) (Outcome, error) {
	outcome, err := arm.delegate.Run(ctx, request)
	if err != nil {
		return Outcome{}, err
	}
	if err := os.WriteFile(filepath.Join(request.RunDirectory, "config.json"), []byte(`{"tampered":true}`), 0o600); err != nil {
		return Outcome{}, err
	}
	return outcome, nil
}

type hardLinkArtifactArm struct{ name string }

func (arm *hardLinkArtifactArm) Name() string { return arm.name }
func (arm *hardLinkArtifactArm) Run(_ context.Context, request Request) (Outcome, error) {
	source := filepath.Join(request.RunDirectory, "results", "cells.jsonl")
	artifactPath := filepath.Join(request.NativeDirectory, "result.json")
	if err := os.Link(source, artifactPath); err != nil {
		return Outcome{}, err
	}
	relative, err := filepath.Rel(request.RunDirectory, artifactPath)
	if err != nil {
		return Outcome{}, err
	}
	return Outcome{
		Completed: true, ContractValid: true,
		Metrics:        map[string]float64{"quality": 0.8},
		NativeArtifact: ArtifactRef{Path: relative},
	}, nil
}

type externalHardLinkArm struct {
	name   string
	source string
}

func (arm *externalHardLinkArm) Name() string { return arm.name }
func (arm *externalHardLinkArm) Run(_ context.Context, request Request) (Outcome, error) {
	artifactPath := filepath.Join(request.NativeDirectory, "result.json")
	if err := os.Link(arm.source, artifactPath); err != nil {
		return Outcome{}, err
	}
	relative, err := filepath.Rel(request.RunDirectory, artifactPath)
	if err != nil {
		return Outcome{}, err
	}
	return Outcome{Completed: true, ContractValid: true, Metrics: map[string]float64{"quality": 0.8}, NativeArtifact: ArtifactRef{Path: relative}}, nil
}

type cancelingSuccessArm struct {
	delegate Arm
	cancel   context.CancelFunc
}

func (arm *cancelingSuccessArm) Name() string { return arm.delegate.Name() }
func (arm *cancelingSuccessArm) Run(ctx context.Context, request Request) (Outcome, error) {
	outcome, err := arm.delegate.Run(ctx, request)
	arm.cancel()
	return outcome, err
}

type cancelingErrorArm struct {
	name   string
	cancel context.CancelFunc
}

func (arm *cancelingErrorArm) Name() string { return arm.name }
func (arm *cancelingErrorArm) Run(_ context.Context, _ Request) (Outcome, error) {
	arm.cancel()
	return Outcome{}, errors.New("request aborted")
}

type partialCancelArm struct {
	name   string
	cancel context.CancelFunc
}

func (arm *partialCancelArm) Name() string { return arm.name }
func (arm *partialCancelArm) Run(ctx context.Context, request Request) (Outcome, error) {
	if err := os.WriteFile(filepath.Join(request.NativeDirectory, "partial.tmp"), []byte("partial"), 0o600); err != nil {
		return Outcome{}, err
	}
	arm.cancel()
	return Outcome{}, ctx.Err()
}

type emptyDirectoryArm struct {
	delegate Arm
	called   bool
}

func (arm *emptyDirectoryArm) Name() string { return arm.delegate.Name() }
func (arm *emptyDirectoryArm) Run(ctx context.Context, request Request) (Outcome, error) {
	entries, err := os.ReadDir(request.NativeDirectory)
	if err != nil {
		return Outcome{}, err
	}
	if len(entries) != 0 {
		return Outcome{}, fmt.Errorf("native directory was not cleared: %v", entries)
	}
	arm.called = true
	return arm.delegate.Run(ctx, request)
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

func mustCreateEvidenceRun(t *testing.T) *runstore.Run {
	t.Helper()
	run, err := runstore.Create(t.Context(), runstore.Options{Root: t.TempDir(), Name: "evidence"}, map[string]string{"config": "fixed"})
	mustNoError(t, err)
	source := filepath.Join(t.TempDir(), "suite.json")
	writeFile(t, source, []byte(`{}`))
	_, err = run.CopyInput(t.Context(), "suite", source)
	mustNoError(t, err)
	return run
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
