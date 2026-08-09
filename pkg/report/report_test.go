package report

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/go-go-golems/ragopt/pkg/candidate"
	"github.com/go-go-golems/ragopt/pkg/compare"
	"github.com/go-go-golems/ragopt/pkg/eval"
	"github.com/go-go-golems/ragopt/pkg/gate"
	"github.com/go-go-golems/ragopt/pkg/policy"
	"github.com/go-go-golems/ragopt/pkg/runstore"
)

func TestBuildAndWritePromotionEvidenceWithoutApplying(t *testing.T) {
	run, comparison, policy, decision := reportFixture()
	first, err := Build(t.Context(), run, comparison, policy, decision)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Build(t.Context(), run, comparison, policy, decision)
	if err != nil {
		t.Fatal(err)
	}
	if first.Markdown != second.Markdown {
		t.Fatal("report rendering is not deterministic")
	}
	for _, section := range []string{"Evidence identity", "Proposed mutation", "Gate decision", "Group outcomes", "Metric aggregates", "Missing pairs", "Paired results", "Review action"} {
		if !strings.Contains(first.Markdown, "## "+section) {
			t.Errorf("report is missing %q section", section)
		}
	}
	if first.Plan.State != "review_required" || !first.Plan.HumanApplyRequired || first.Plan.Decision != gate.DecisionPass {
		t.Fatalf("unsafe or incorrect plan: %#v", first.Plan)
	}
	directory := t.TempDir()
	markdownPath := filepath.Join(directory, "review.md")
	planPath := filepath.Join(directory, "plan.json")
	if err := Write(t.Context(), first, markdownPath, planPath); err != nil {
		t.Fatal(err)
	}
	markdown, err := os.ReadFile(markdownPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(markdown) != first.Markdown {
		t.Fatal("published Markdown differs from rendered report")
	}
	planData, err := os.ReadFile(planPath)
	if err != nil {
		t.Fatal(err)
	}
	var plan PromotionPlan
	if err := json.Unmarshal(planData, &plan); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(plan, first.Plan) {
		t.Fatalf("published plan differs: %#v", plan)
	}
}

func TestBuildRejectsInconsistentIdentities(t *testing.T) {
	mutations := map[string]func(*eval.ArtifactRun, *compare.Report, *policy.Document){
		"run": func(_ *eval.ArtifactRun, report *compare.Report, _ *policy.Document) { report.RunID = "other-run" },
		"suite": func(_ *eval.ArtifactRun, report *compare.Report, _ *policy.Document) {
			report.SuiteDigest = reportDigest('x')
		},
		"policy bytes": func(_ *eval.ArtifactRun, _ *compare.Report, policy *policy.Document) {
			policy.ByteDigest = reportDigest('x')
		},
		"candidate ID": func(_ *eval.ArtifactRun, report *compare.Report, _ *policy.Document) {
			report.CandidateID = "other"
		},
		"candidate digest": func(_ *eval.ArtifactRun, report *compare.Report, _ *policy.Document) {
			report.CandidateDigest = reportDigest('x')
		},
		"parent snapshot": func(_ *eval.ArtifactRun, report *compare.Report, _ *policy.Document) {
			report.ParentSnapshot = reportDigest('x')
		},
		"child snapshot": func(_ *eval.ArtifactRun, report *compare.Report, _ *policy.Document) {
			report.ChildSnapshot = reportDigest('x')
		},
		"incumbent arm": func(_ *eval.ArtifactRun, report *compare.Report, _ *policy.Document) {
			report.IncumbentArm = "other"
		},
		"challenger arm": func(_ *eval.ArtifactRun, report *compare.Report, _ *policy.Document) {
			report.ChallengerArm = "other"
		},
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			run, comparison, policy, decision := reportFixture()
			mutate(run, comparison, policy)
			if _, err := Build(t.Context(), run, comparison, policy, decision); err == nil {
				t.Fatal("expected inconsistent report identities to fail")
			}
		})
	}
}

func TestBuildRejectsDecisionFromDifferentComparison(t *testing.T) {
	run, comparison, policy, decision := reportFixture()
	comparison.Pairs[0].Deltas[0].Delta = -1
	if _, err := Build(t.Context(), run, comparison, policy, decision); err == nil {
		t.Fatal("expected stale gate decision to be rejected")
	}
}

func TestBuildRejectsComparisonAndDecisionNotBackedByRun(t *testing.T) {
	run, comparison, policy, _ := reportFixture()
	comparison.Pairs[0].Deltas[0].Candidate = 2
	comparison.Pairs[0].Deltas[0].Delta = 1.5
	comparison.Metrics[0].MeanCandidate = 2
	comparison.Metrics[0].MeanDelta = 1.5
	decision, err := gate.Evaluate(t.Context(), policy, comparison)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Build(t.Context(), run, comparison, policy, decision); err == nil || !strings.Contains(err.Error(), "differs from artifact run evidence") {
		t.Fatalf("expected fabricated comparison rejection, got %v", err)
	}
}

func TestBuildRejectsIncompleteDirectoryBackedRun(t *testing.T) {
	run, comparison, policy, decision := reportFixture()
	run.Config.APIVersion = eval.RunAPIVersion
	run.Config.InputDigests = map[string]string{"suite": reportDigest('s')}
	directory := t.TempDir()
	data, err := json.Marshal(run.Config)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "config.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(data)
	run.Directory = directory
	run.Manifest.ConfigDigest = "sha256:" + hex.EncodeToString(digest[:])
	if _, err := Build(t.Context(), run, comparison, policy, decision); err == nil || !strings.Contains(err.Error(), "load durable evaluation evidence") {
		t.Fatalf("incomplete directory-backed run was accepted: %v", err)
	}
}

func TestWriteRejectsSameResolvedOutputPath(t *testing.T) {
	run, comparison, policy, decision := reportFixture()
	document, err := Build(t.Context(), run, comparison, policy, decision)
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	t.Chdir(directory)
	if err := Write(t.Context(), document, "review.out", filepath.Join(directory, "review.out")); err == nil {
		t.Fatal("expected relative and absolute aliases to be rejected")
	}
}

func TestReportDestinationAliasUsesFilesystemCaseSemantics(t *testing.T) {
	directory := t.TempDir()
	caseInsensitive, err := filesystemCaseInsensitive(directory)
	if err != nil {
		t.Fatal(err)
	}
	alias, err := destinationsAlias(
		filepath.Join(directory, "review.out"),
		filepath.Join(directory, "Review.out"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if alias != caseInsensitive {
		t.Fatalf("case-folded alias = %v, filesystem case-insensitive = %v", alias, caseInsensitive)
	}
}

func TestWriteDoesNotPublishEitherOutputWhenOneDestinationIsInvalid(t *testing.T) {
	directory := t.TempDir()
	markdownPath := filepath.Join(directory, "review.md")
	planPath := filepath.Join(directory, "plan.json")
	if err := os.WriteFile(markdownPath, []byte("old markdown"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(planPath, 0o700); err != nil {
		t.Fatal(err)
	}
	err := Write(t.Context(), &Document{Markdown: "new markdown"}, markdownPath, planPath)
	if err == nil {
		t.Fatal("expected invalid plan destination to reject transaction")
	}
	data, readErr := os.ReadFile(markdownPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(data) != "old markdown" {
		t.Fatalf("markdown was partially published: %q", data)
	}
}

func TestPublishOutputsRollsBackFirstOutputWhenSecondCommitFails(t *testing.T) {
	directory := t.TempDir()
	markdownPath := filepath.Join(directory, "review.md")
	planPath := filepath.Join(directory, "plan.json")
	if err := os.WriteFile(markdownPath, []byte("old markdown"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(planPath, []byte("old plan"), 0o600); err != nil {
		t.Fatal(err)
	}
	outputs := []*stagedOutput{{path: markdownPath, data: []byte("new markdown")}, {path: planPath, data: []byte("new plan")}}
	for _, output := range outputs {
		if err := output.stage(t.Context()); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Remove(outputs[1].temporary); err != nil {
		t.Fatal(err)
	}
	if err := publishOutputs(t.Context(), outputs); err == nil {
		t.Fatal("expected second publication to fail")
	}
	cleanupStaged(outputs)
	for path, want := range map[string]string{markdownPath: "old markdown", planPath: "old plan"} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != want {
			t.Fatalf("rollback for %s = %q, want %q", path, data, want)
		}
	}
}

func TestValidateOutputsOutsideRun(t *testing.T) {
	root := t.TempDir()
	runDirectory := filepath.Join(root, "run")
	if err := os.MkdirAll(runDirectory, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := ValidateOutputsOutsideRun(runDirectory, filepath.Join(root, "review.md")); err != nil {
		t.Fatal(err)
	}
	if err := ValidateOutputsOutsideRun(runDirectory, filepath.Join(runDirectory, "config.json")); err == nil {
		t.Fatal("expected output inside the evaluated run to be rejected")
	}
}

func TestValidateOutputsOutsideRunUsesFilesystemCaseSemantics(t *testing.T) {
	root := t.TempDir()
	caseInsensitive, err := filesystemCaseInsensitive(root)
	if err != nil {
		t.Fatal(err)
	}
	if !caseInsensitive {
		t.Skip("filesystem is case-sensitive")
	}
	runDirectory := filepath.Join(root, "RunID")
	if err := os.MkdirAll(runDirectory, 0o700); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(root, "runid", "CONFIG.JSON")
	if err := ValidateOutputsOutsideRun(runDirectory, output); err == nil {
		t.Fatal("case-varied output inside evaluated run was accepted")
	}
}

func TestReportPathsResolveExistingParentSymlinks(t *testing.T) {
	root := t.TempDir()
	runDirectory := filepath.Join(root, "run")
	outputs := filepath.Join(root, "outputs")
	if err := os.MkdirAll(runDirectory, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outputs, 0o700); err != nil {
		t.Fatal(err)
	}
	runLink := filepath.Join(root, "run-link")
	outputLink := filepath.Join(root, "output-link")
	if err := os.Symlink(runDirectory, runLink); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if err := os.Symlink(outputs, outputLink); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if err := ValidateOutputsOutsideRun(runDirectory, filepath.Join(runLink, "config.json")); err == nil {
		t.Fatal("expected symlinked run destination to be rejected")
	}
	run, comparison, policy, decision := reportFixture()
	document, err := Build(t.Context(), run, comparison, policy, decision)
	if err != nil {
		t.Fatal(err)
	}
	if err := Write(t.Context(), document, filepath.Join(outputLink, "same.out"), filepath.Join(outputs, "same.out")); err == nil {
		t.Fatal("expected symlink aliases of one output to be rejected")
	}
}

func reportFixture() (*eval.ArtifactRun, *compare.Report, *policy.Document, gate.Decision) {
	mutation := candidate.MutationDeclaration{
		Asset: "prompt", Hypothesis: "Make comparison behavior explicit.",
		ExpectedImprovement: candidate.ExpectedImprovement{Metric: "quality", Groups: []string{"comparison"}},
		RegressionRisks:     []string{"verbosity"},
	}
	run := &eval.ArtifactRun{
		Manifest: runstore.Manifest{RunID: "run-1"},
		Status:   runstore.Status{State: runstore.StateComplete},
		Config: eval.RunConfig{
			SuiteDigest: reportDigest('s'), PolicyDigest: reportDigest('p'),
			CandidateID: "candidate-1", CandidateDigest: reportDigest('c'),
			ParentSnapshot: reportDigest('i'), ChildSnapshot: reportDigest('n'),
			IncumbentArm: "incumbent", ChallengerArm: "candidate", Repeats: 1,
			Mutation: mutation, ChangedAsset: "prompt", ParentAssetDigest: reportDigest('a'), ChildAssetDigest: reportDigest('b'),
		},
	}
	run.Suite = &eval.SuiteDocument{Digest: run.Config.SuiteDigest, Suite: eval.Suite{
		APIVersion: eval.SuiteAPIVersion, Name: "report-fixture",
		Cases: []eval.Case{{ID: "case-a", Groups: []string{"comparison"}, Input: []byte(`{"id":"a"}`)}},
	}}
	run.Cells = []eval.Cell{
		reportCell(run, run.Config.IncumbentArm, run.Config.ParentSnapshot, 0.5),
		reportCell(run, run.Config.ChallengerArm, run.Config.ChildSnapshot, 0.6),
	}
	comparison, err := compare.Build(context.Background(), run)
	if err != nil {
		panic(err)
	}
	maximumFailure := 1.0
	policy := &policy.Document{ByteDigest: reportDigest('p'), Policy: policy.Policy{
		APIVersion: policy.PolicyAPIVersion, Name: "policy-1",
		HardGates: policy.HardGates{MaxFailureRate: &maximumFailure},
		Target:    policy.Target{Metric: "quality", Groups: []string{"comparison"}, MinimumMeanDelta: 0.05},
	}}
	semantic, err := json.Marshal(policy.Policy)
	if err != nil {
		panic(err)
	}
	sum := sha256.Sum256(semantic)
	policy.Digest = "sha256:" + hex.EncodeToString(sum[:])
	decision, err := gate.Evaluate(context.Background(), policy, comparison)
	if err != nil {
		panic(err)
	}
	return run, comparison, policy, decision
}

func reportCell(run *eval.ArtifactRun, arm, snapshot string, quality float64) eval.Cell {
	return eval.Cell{
		APIVersion: eval.CellAPIVersion, RunID: run.Manifest.RunID,
		CaseID: "case-a", Arm: arm, CandidateID: run.Config.CandidateID,
		SnapshotDigest: snapshot, SuiteDigest: run.Config.SuiteDigest, PolicyDigest: run.Config.PolicyDigest,
		Outcome: eval.Outcome{Completed: true, ContractValid: true, Metrics: map[string]float64{"quality": quality}},
	}
}

func reportDigest(value byte) string { return "sha256:" + strings.Repeat(string(value), 64) }
