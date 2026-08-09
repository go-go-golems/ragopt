package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-go-golems/ragopt/pkg/compare"
	"github.com/go-go-golems/ragopt/pkg/eval"
	"github.com/go-go-golems/ragopt/pkg/runstore"
)

func TestGateDecisionGoldens(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*eval.ArtifactRun, *PolicyDocument)
	}{
		{name: "pass"},
		{name: "hard-fail", mutate: func(run *eval.ArtifactRun, _ *PolicyDocument) {
			run.Cells[1].Outcome.Completed = false
			run.Cells[1].Outcome.Failure = &eval.Failure{Class: "provider", Message: "unavailable"}
		}},
		{name: "target-fail", mutate: func(run *eval.ArtifactRun, _ *PolicyDocument) {
			run.Cells[1].Outcome.Metrics["quality"] = 0.5
			run.Cells[3].Outcome.Metrics["quality"] = 0.6
		}},
		{name: "catastrophic-regression", mutate: func(run *eval.ArtifactRun, _ *PolicyDocument) {
			run.Cells[1].Outcome.Metrics["safety"] = 0.65
		}},
		{name: "tie-break", mutate: func(run *eval.ArtifactRun, policy *PolicyDocument) {
			run.Cells[1].Outcome.Metrics["quality"] = 0.5
			run.Cells[3].Outcome.Metrics["quality"] = 0.6
			policy.Policy.Target.MinimumMeanDelta = 0
			policy.Policy.Target.RequirePositiveEachRepeat = false
		}},
		{name: "incomplete-pairing", mutate: func(run *eval.ArtifactRun, _ *PolicyDocument) {
			run.Cells = run.Cells[:len(run.Cells)-1]
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			run := gateFixture()
			policy := gatePolicy(run.Config.PolicyDigest)
			if test.mutate != nil {
				test.mutate(run, policy)
			}
			policy.Digest = mustPolicyDigest(policy.Policy)
			report, err := compare.Build(t.Context(), run)
			if err != nil {
				t.Fatal(err)
			}
			decision, err := Evaluate(t.Context(), policy, report)
			if err != nil {
				t.Fatal(err)
			}
			actual := decisionSummary(decision)
			expected, err := os.ReadFile(filepath.Join("testdata", test.name+".golden"))
			if err != nil {
				t.Fatal(err)
			}
			if actual != string(expected) {
				t.Fatalf("decision mismatch\n--- actual ---\n%s--- expected ---\n%s", actual, expected)
			}
		})
	}
}

func TestEvaluateRejectsMutatedPolicyDocument(t *testing.T) {
	run := gateFixture()
	policy := gatePolicy(run.Config.PolicyDigest)
	report, err := compare.Build(t.Context(), run)
	if err != nil {
		t.Fatal(err)
	}
	policy.Policy.Target.MinimumMeanDelta = -100
	if _, err := Evaluate(t.Context(), policy, report); err == nil || !strings.Contains(err.Error(), "semantic digest mismatch") {
		t.Fatalf("expected mutated policy rejection, got %v", err)
	}
}

func TestTargetAllSelectsEveryPair(t *testing.T) {
	run := gateFixture()
	policy := gatePolicy(run.Config.PolicyDigest)
	policy.Policy.Target.Groups = []string{"all"}
	policy.Digest = mustPolicyDigest(policy.Policy)
	report, err := compare.Build(t.Context(), run)
	if err != nil {
		t.Fatal(err)
	}
	decision, err := Evaluate(t.Context(), policy, report)
	if err != nil {
		t.Fatal(err)
	}
	if decision.Status != DecisionPass {
		t.Fatalf("target all did not select the complete report: %#v", decision)
	}
}

func TestEvaluateRejectsTargetMeanOverflow(t *testing.T) {
	run := gateFixture()
	policy := gatePolicy(run.Config.PolicyDigest)
	policy.Policy.Target.Groups = []string{"selected"}
	policy.Digest = mustPolicyDigest(policy.Policy)
	report, err := compare.Build(t.Context(), run)
	if err != nil {
		t.Fatal(err)
	}
	for index := range report.Pairs {
		report.Pairs[index].Groups = []string{"selected"}
		for deltaIndex := range report.Pairs[index].Deltas {
			if report.Pairs[index].Deltas[deltaIndex].Metric == "quality" {
				report.Pairs[index].Deltas[deltaIndex].Delta = 1.0e308
			}
		}
	}
	if _, err := Evaluate(t.Context(), policy, report); err == nil || !strings.Contains(err.Error(), "mean accumulation is not finite") {
		t.Fatalf("expected target mean overflow rejection, got %v", err)
	}
}

func TestEvaluateRejectsIncompleteRunStateBeforePromotionChecks(t *testing.T) {
	for _, state := range []string{runstore.StateActive, runstore.StateFailed} {
		t.Run(string(state), func(t *testing.T) {
			run := gateFixture()
			run.Status.State = state
			policy := gatePolicy(run.Config.PolicyDigest)
			report, err := compare.Build(t.Context(), run)
			if err != nil {
				t.Fatal(err)
			}
			decision, err := Evaluate(t.Context(), policy, report)
			if err != nil {
				t.Fatal(err)
			}
			if decision.Status != DecisionFail || len(decision.Checks) < 2 || decision.Checks[1].Name != "run_complete" || decision.Checks[1].Passed {
				t.Fatalf("incomplete run was not rejected: %#v", decision)
			}
		})
	}
}

func TestLoadPolicyIsStrictAndHasNoThresholdDefaults(t *testing.T) {
	directory := t.TempDir()
	valid := `api_version: ragopt-gate-policy/v1
name: fixture
hard_gates:
  require_all_cells: true
  require_completed: true
  require_contract_valid: true
  max_failure_rate: 0
target:
  metric: quality
  minimum_mean_delta: 0.1
  require_positive_each_repeat: true
regressions: {}
`
	path := filepath.Join(directory, "policy.yaml")
	if err := os.WriteFile(path, []byte(valid), 0o600); err != nil {
		t.Fatal(err)
	}
	document, err := LoadPolicy(t.Context(), path)
	if err != nil {
		t.Fatal(err)
	}
	if document.Policy.Target.Metric != "quality" || document.Policy.HardGates.MaxFailureRate == nil {
		t.Fatalf("loaded policy: %#v", document.Policy)
	}

	invalid := strings.Replace(valid, "  max_failure_rate: 0\n", "", 1)
	if err := os.WriteFile(path, []byte(invalid), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadPolicy(t.Context(), path); err == nil || !strings.Contains(err.Error(), "max_failure_rate is required") {
		t.Fatalf("missing threshold error: %v", err)
	}
	missingTarget := strings.Replace(valid, "  minimum_mean_delta: 0.1\n", "", 1)
	if err := os.WriteFile(path, []byte(missingTarget), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadPolicy(t.Context(), path); err == nil || !strings.Contains(err.Error(), "minimum_mean_delta is required") {
		t.Fatalf("missing target threshold error: %v", err)
	}
	unknown := valid + "unknown: true\n"
	if err := os.WriteFile(path, []byte(unknown), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadPolicy(t.Context(), path); err == nil || !strings.Contains(err.Error(), "field unknown not found") {
		t.Fatalf("unknown-field error: %v", err)
	}
}

func gateFixture() *eval.ArtifactRun {
	config := eval.RunConfig{
		APIVersion: eval.RunAPIVersion, SuiteDigest: testDigest('s'), PolicyDigest: testDigest('p'),
		CandidateID: "candidate-1", CandidateDigest: testDigest('c'), ParentSnapshot: testDigest('i'), ChildSnapshot: testDigest('n'),
		IncumbentArm: "incumbent", ChallengerArm: "candidate", Repeats: 1,
	}
	suite := &eval.SuiteDocument{Digest: config.SuiteDigest, Suite: eval.Suite{
		APIVersion: eval.SuiteAPIVersion, Name: "fixture", Cases: []eval.Case{
			{ID: "case-a", Groups: []string{"comparison"}, Input: []byte(`{"id":"a"}`)},
			{ID: "case-b", Groups: []string{"factual"}, Input: []byte(`{"id":"b"}`)},
		},
	}}
	run := &eval.ArtifactRun{
		Manifest: runstore.Manifest{RunID: "run-1"},
		Status:   runstore.Status{State: runstore.StateComplete},
		Config:   config, Suite: suite,
	}
	for index, caseValue := range suite.Suite.Cases {
		incumbentQuality := 0.5 + float64(index)*0.1
		candidateQuality := incumbentQuality + 0.1
		run.Cells = append(run.Cells,
			gateCell(config, caseValue.ID, config.IncumbentArm, config.ParentSnapshot, incumbentQuality, 2, 20),
			gateCell(config, caseValue.ID, config.ChallengerArm, config.ChildSnapshot, candidateQuality, 1, 15),
		)
	}
	return run
}

func gateCell(config eval.RunConfig, caseID, arm, snapshot string, quality float64, calls, tokens int) eval.Cell {
	return eval.Cell{
		APIVersion: eval.CellAPIVersion, RunID: "run-1", CaseID: caseID, Arm: arm,
		CandidateID: config.CandidateID, SnapshotDigest: snapshot, SuiteDigest: config.SuiteDigest, PolicyDigest: config.PolicyDigest,
		StartedAt: time.Unix(1, 0), FinishedAt: time.Unix(2, 0),
		Outcome: eval.Outcome{Completed: true, ContractValid: true, Metrics: map[string]float64{"quality": quality, "safety": 0.9}, ProviderCalls: calls, InputTokens: tokens, Duration: time.Second},
	}
}

func gatePolicy(byteDigest string) *PolicyDocument {
	maximumFailure := 0.0
	document := &PolicyDocument{
		ByteDigest: byteDigest,
		Policy: Policy{
			APIVersion: PolicyAPIVersion, Name: "fixture-policy",
			HardGates:   HardGates{RequireAllCells: true, RequireCompleted: true, RequireContractValid: true, MaxFailureRate: &maximumFailure, MetricFloors: map[string]float64{"safety": 0.6}},
			Target:      Target{Metric: "quality", MinimumMeanDelta: 0.05, RequirePositiveEachRepeat: true},
			Regressions: Regressions{MaximumCaseDelta: map[string]float64{"safety": -0.2}, MaximumMeanDelta: map[string]map[string]float64{"all": {"quality": -0.2}}},
			TieBreakers: []string{"provider_calls", "total_tokens"},
		},
	}
	document.Digest = mustPolicyDigest(document.Policy)
	return document
}

func mustPolicyDigest(policy Policy) string {
	digest, err := policyDigest(policy)
	if err != nil {
		panic(err)
	}
	return digest
}

func decisionSummary(decision Decision) string {
	var result strings.Builder
	result.WriteString("status=" + string(decision.Status) + "\n")
	for _, check := range decision.Checks {
		status := "fail"
		if check.Passed {
			status = "pass"
		}
		result.WriteString(check.Phase + "/" + check.Name + "=" + status + "\n")
	}
	return result.String()
}

func testDigest(value byte) string { return "sha256:" + strings.Repeat(string(value), 64) }
