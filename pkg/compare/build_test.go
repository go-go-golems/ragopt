package compare

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/go-go-golems/ragopt/pkg/eval"
	"github.com/go-go-golems/ragopt/pkg/runstore"
)

func TestBuildPairsAndAggregatesWithoutHidingMissingMetrics(t *testing.T) {
	run := comparisonFixture()
	// Quality is absent from one candidate cell. The pair remains complete,
	// while metric presence remains three of four.
	delete(run.Cells[7].Outcome.Metrics, "quality")
	report, err := Build(t.Context(), run)
	if err != nil {
		t.Fatal(err)
	}
	if report.ExpectedPairs != 4 || report.CompletePairs != 4 || len(report.MissingPairs) != 0 {
		t.Fatalf("pair counts: expected=%d complete=%d missing=%d", report.ExpectedPairs, report.CompletePairs, len(report.MissingPairs))
	}
	if report.RunState != runstore.StateComplete {
		t.Fatalf("run state = %q, want complete", report.RunState)
	}
	quality := findMetricForTest(t, report, "all", "quality")
	if quality.CompletePairs != 4 || quality.PairsWithMetric != 3 {
		t.Fatalf("quality denominators: %#v", quality)
	}
	if quality.Wins != 2 || quality.Ties != 0 || quality.Losses != 1 {
		t.Fatalf("quality outcomes: %#v", quality)
	}
	comparisonGroup := findGroupForTest(t, report, "comparison")
	if comparisonGroup.ExpectedPairs != 2 || comparisonGroup.CompletePairs != 2 {
		t.Fatalf("comparison group: %#v", comparisonGroup)
	}
	if got := report.Pairs[0].Deltas[0].Delta; math.Abs(got-0.1) > 1e-9 {
		t.Fatalf("raw candidate-incumbent delta: got %.6f", got)
	}
}

func TestBuildRetainsMissingPair(t *testing.T) {
	run := comparisonFixture()
	run.Cells = run.Cells[:len(run.Cells)-1]
	report, err := Build(t.Context(), run)
	if err != nil {
		t.Fatal(err)
	}
	if report.CompletePairs != 3 || len(report.MissingPairs) != 1 || !report.MissingPairs[0].MissingCandidate {
		t.Fatalf("missing pairing: %#v", report.MissingPairs)
	}
}

func TestBuildRejectsDuplicateAndCrossIdentityCells(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*eval.ArtifactRun)
		want   string
	}{
		{name: "duplicate", mutate: func(run *eval.ArtifactRun) { run.Cells = append(run.Cells, run.Cells[0]) }, want: "duplicate comparison cell"},
		{name: "suite", mutate: func(run *eval.ArtifactRun) { run.Cells[0].SuiteDigest = digest('x') }, want: "cross-run identity"},
		{name: "policy", mutate: func(run *eval.ArtifactRun) { run.Cells[0].PolicyDigest = digest('x') }, want: "cross-run identity"},
		{name: "snapshot", mutate: func(run *eval.ArtifactRun) { run.Cells[0].SnapshotDigest = digest('x') }, want: "cross-run identity"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			run := comparisonFixture()
			test.mutate(run)
			_, err := Build(t.Context(), run)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("got %v, want error containing %q", err, test.want)
			}
		})
	}
}

func TestBuildRejectsNonFiniteDerivedArithmetic(t *testing.T) {
	t.Run("delta overflow", func(t *testing.T) {
		run := comparisonFixture()
		run.Cells[0].Outcome.Metrics["quality"] = -math.MaxFloat64
		run.Cells[1].Outcome.Metrics["quality"] = math.MaxFloat64
		_, err := Build(t.Context(), run)
		if err == nil || !strings.Contains(err.Error(), "delta") {
			t.Fatalf("expected non-finite delta rejection, got %v", err)
		}
	})
	t.Run("aggregate overflow", func(t *testing.T) {
		run := comparisonFixture()
		for index := range run.Cells {
			run.Cells[index].Outcome.Metrics["quality"] = math.MaxFloat64
		}
		_, err := Build(t.Context(), run)
		if err == nil || !strings.Contains(err.Error(), "aggregate") {
			t.Fatalf("expected non-finite aggregate rejection, got %v", err)
		}
	})
}

func comparisonFixture() *eval.ArtifactRun {
	config := eval.RunConfig{
		APIVersion: eval.RunAPIVersion, SuiteDigest: digest('s'), PolicyDigest: digest('p'),
		CandidateID: "candidate-1", CandidateDigest: digest('c'), ParentSnapshot: digest('i'), ChildSnapshot: digest('n'),
		IncumbentArm: "incumbent", ChallengerArm: "candidate", Repeats: 2,
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
	values := []struct{ incumbent, candidate float64 }{{0.5, 0.6}, {0.7, 0.8}, {0.8, 0.7}, {0.5, 0.6}}
	index := 0
	for _, caseValue := range suite.Suite.Cases {
		for repeat := 0; repeat < config.Repeats; repeat++ {
			value := values[index]
			index++
			run.Cells = append(run.Cells,
				comparisonCell(config, caseValue.ID, repeat, config.IncumbentArm, config.ParentSnapshot, value.incumbent),
				comparisonCell(config, caseValue.ID, repeat, config.ChallengerArm, config.ChildSnapshot, value.candidate),
			)
		}
	}
	return run
}

func comparisonCell(config eval.RunConfig, caseID string, repeat int, arm, snapshot string, quality float64) eval.Cell {
	return eval.Cell{
		APIVersion: eval.CellAPIVersion, RunID: "run-1", CaseID: caseID, RepeatIndex: repeat, Arm: arm,
		CandidateID: config.CandidateID, SnapshotDigest: snapshot, SuiteDigest: config.SuiteDigest, PolicyDigest: config.PolicyDigest,
		StartedAt: time.Unix(1, 0), FinishedAt: time.Unix(2, 0),
		Outcome: eval.Outcome{Completed: true, ContractValid: true, Metrics: map[string]float64{"quality": quality, "safety": 0.9}, ProviderCalls: 1, InputTokens: 10, Duration: time.Second},
	}
}

func digest(value byte) string { return "sha256:" + strings.Repeat(string(value), 64) }

func findMetricForTest(t *testing.T, report *Report, group, metric string) MetricAggregate {
	t.Helper()
	for _, value := range report.Metrics {
		if value.Group == group && value.Metric == metric {
			return value
		}
	}
	t.Fatalf("missing metric %s/%s", group, metric)
	return MetricAggregate{}
}

func findGroupForTest(t *testing.T, report *Report, group string) GroupAggregate {
	t.Helper()
	for _, value := range report.Groups {
		if value.Group == group {
			return value
		}
	}
	t.Fatalf("missing group %s", group)
	return GroupAggregate{}
}
