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
	mutations := map[string]func(*eval.ArtifactRun, *compare.Report, *gate.PolicyDocument){
		"run": func(_ *eval.ArtifactRun, report *compare.Report, _ *gate.PolicyDocument) { report.RunID = "other-run" },
		"suite": func(_ *eval.ArtifactRun, report *compare.Report, _ *gate.PolicyDocument) {
			report.SuiteDigest = reportDigest('x')
		},
		"policy bytes": func(_ *eval.ArtifactRun, _ *compare.Report, policy *gate.PolicyDocument) {
			policy.ByteDigest = reportDigest('x')
		},
		"candidate ID": func(_ *eval.ArtifactRun, report *compare.Report, _ *gate.PolicyDocument) {
			report.CandidateID = "other"
		},
		"candidate digest": func(_ *eval.ArtifactRun, report *compare.Report, _ *gate.PolicyDocument) {
			report.CandidateDigest = reportDigest('x')
		},
		"parent snapshot": func(_ *eval.ArtifactRun, report *compare.Report, _ *gate.PolicyDocument) {
			report.ParentSnapshot = reportDigest('x')
		},
		"child snapshot": func(_ *eval.ArtifactRun, report *compare.Report, _ *gate.PolicyDocument) {
			report.ChildSnapshot = reportDigest('x')
		},
		"incumbent arm": func(_ *eval.ArtifactRun, report *compare.Report, _ *gate.PolicyDocument) {
			report.IncumbentArm = "other"
		},
		"challenger arm": func(_ *eval.ArtifactRun, report *compare.Report, _ *gate.PolicyDocument) {
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

func reportFixture() (*eval.ArtifactRun, *compare.Report, *gate.PolicyDocument, gate.Decision) {
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
	policy := &gate.PolicyDocument{ByteDigest: reportDigest('p'), Policy: gate.Policy{
		APIVersion: gate.PolicyAPIVersion, Name: "policy-1",
		HardGates: gate.HardGates{MaxFailureRate: &maximumFailure},
		Target:    gate.Target{Metric: "quality", Groups: []string{"comparison"}, MinimumMeanDelta: 0.05},
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
