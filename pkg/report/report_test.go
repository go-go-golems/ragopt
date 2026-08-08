package report

import (
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

func reportFixture() (*eval.ArtifactRun, *compare.Report, *gate.PolicyDocument, gate.Decision) {
	mutation := candidate.MutationDeclaration{
		Asset: "prompt", Hypothesis: "Make comparison behavior explicit.",
		ExpectedImprovement: candidate.ExpectedImprovement{Metric: "quality", Groups: []string{"comparison"}},
		RegressionRisks:     []string{"verbosity"},
	}
	run := &eval.ArtifactRun{
		Manifest: runstore.Manifest{RunID: "run-1"},
		Config: eval.RunConfig{
			SuiteDigest: reportDigest('s'), PolicyDigest: reportDigest('p'),
			CandidateID: "candidate-1", CandidateDigest: reportDigest('c'),
			ParentSnapshot: reportDigest('i'), ChildSnapshot: reportDigest('n'),
			IncumbentArm: "incumbent", ChallengerArm: "candidate",
			Mutation: mutation, ChangedAsset: "prompt", ParentAssetDigest: reportDigest('a'), ChildAssetDigest: reportDigest('b'),
		},
	}
	comparison := &compare.Report{
		APIVersion: compare.ReportAPIVersion, RunID: "run-1", SuiteDigest: reportDigest('s'), PolicyDigest: reportDigest('p'),
		CandidateID: "candidate-1", CandidateDigest: reportDigest('c'), ParentSnapshot: reportDigest('i'), ChildSnapshot: reportDigest('n'),
		IncumbentArm: "incumbent", ChallengerArm: "candidate",
		ExpectedPairs: 1, CompletePairs: 1,
		Groups:  []compare.GroupAggregate{{Group: "all", ExpectedPairs: 1, CompletePairs: 1, CandidateCompleted: 1, CandidateContractValid: 1}},
		Metrics: []compare.MetricAggregate{{Group: "all", Metric: "quality", ExpectedPairs: 1, CompletePairs: 1, PairsWithMetric: 1, MeanIncumbent: 0.5, MeanCandidate: 0.6, MeanDelta: 0.1, Wins: 1}},
		Pairs:   []compare.Pair{{Key: compare.PairKey{CaseID: "case-a"}, Groups: []string{"comparison"}, Candidate: eval.Cell{Outcome: eval.Outcome{Completed: true, ContractValid: true}}, Deltas: []compare.MetricDelta{{Metric: "quality", Incumbent: 0.5, Candidate: 0.6, Delta: 0.1}}, MetricPresence: []compare.MetricPresence{{Metric: "quality", IncumbentPresent: true, CandidatePresent: true}}}},
	}
	policy := &gate.PolicyDocument{Digest: reportDigest('d'), ByteDigest: reportDigest('p'), Policy: gate.Policy{Name: "policy-1"}}
	decision := gate.Decision{APIVersion: gate.DecisionAPIVersion, PolicyName: "policy-1", PolicyDigest: policy.Digest, Status: gate.DecisionPass, Checks: []gate.CheckResult{{Phase: "target", Name: "quality", Passed: true, Message: "quality improved"}}}
	return run, comparison, policy, decision
}

func reportDigest(value byte) string { return "sha256:" + strings.Repeat(string(value), 64) }
