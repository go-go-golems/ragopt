package compare

import (
	"context"
	"math"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/go-go-golems/ragopt/pkg/eval"
)

type coordinate struct {
	caseID string
	repeat int
	arm    string
}

// Build strictly joins incumbent and candidate cells for every suite case and
// repeat. Missing pairs are retained explicitly and never receive synthetic
// metric values.
func Build(ctx context.Context, run *eval.ArtifactRun) (*Report, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if run == nil {
		return nil, errors.New("evaluation artifact run is required")
	}
	durable, err := run.DurableSnapshot(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "load durable evaluation evidence")
	}
	run = durable
	if run.Suite == nil {
		return nil, errors.New("evaluation artifact suite is required")
	}
	config := run.Config
	index := make(map[coordinate]eval.Cell, len(run.Cells))
	validCases := make(map[string]struct{}, len(run.Suite.Suite.Cases))
	for _, caseValue := range run.Suite.Suite.Cases {
		validCases[caseValue.ID] = struct{}{}
	}
	for _, cell := range run.Cells {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if _, exists := validCases[cell.CaseID]; !exists || cell.RepeatIndex < 0 || cell.RepeatIndex >= config.Repeats {
			return nil, errors.Errorf("cell has unexpected case/repeat coordinate %q/%d", cell.CaseID, cell.RepeatIndex)
		}
		if cell.Arm != config.IncumbentArm && cell.Arm != config.ChallengerArm {
			return nil, errors.Errorf("cell has unexpected arm %q", cell.Arm)
		}
		expectedSnapshot := config.ParentSnapshot
		if cell.Arm == config.ChallengerArm {
			expectedSnapshot = config.ChildSnapshot
		}
		if cell.RunID != run.Manifest.RunID || cell.SuiteDigest != config.SuiteDigest ||
			cell.PolicyDigest != config.PolicyDigest || cell.CandidateID != config.CandidateID ||
			cell.SnapshotDigest != expectedSnapshot {
			return nil, errors.Errorf("cell %q/%d/%s has cross-run identity", cell.CaseID, cell.RepeatIndex, cell.Arm)
		}
		for metric, value := range cell.Outcome.Metrics {
			if !finite(value) {
				return nil, errors.Errorf("cell %q/%d/%s metric %q is not finite", cell.CaseID, cell.RepeatIndex, cell.Arm, metric)
			}
		}
		key := coordinate{caseID: cell.CaseID, repeat: cell.RepeatIndex, arm: cell.Arm}
		if _, exists := index[key]; exists {
			return nil, errors.Errorf("duplicate comparison cell %q/%d/%s", cell.CaseID, cell.RepeatIndex, cell.Arm)
		}
		index[key] = cell
	}

	report := &Report{
		APIVersion:      ReportAPIVersion,
		RunID:           run.Manifest.RunID,
		RunState:        run.Status.State,
		SuiteDigest:     config.SuiteDigest,
		PolicyDigest:    config.PolicyDigest,
		CandidateID:     config.CandidateID,
		CandidateDigest: config.CandidateDigest,
		ParentSnapshot:  config.ParentSnapshot,
		ChildSnapshot:   config.ChildSnapshot,
		IncumbentArm:    config.IncumbentArm,
		ChallengerArm:   config.ChallengerArm,
		ExpectedPairs:   len(run.Suite.Suite.Cases) * config.Repeats,
	}
	for _, caseValue := range run.Suite.Suite.Cases {
		for repeat := 0; repeat < config.Repeats; repeat++ {
			incumbent, hasIncumbent := index[coordinate{caseID: caseValue.ID, repeat: repeat, arm: config.IncumbentArm}]
			candidateCell, hasCandidate := index[coordinate{caseID: caseValue.ID, repeat: repeat, arm: config.ChallengerArm}]
			key := PairKey{CaseID: caseValue.ID, RepeatIndex: repeat}
			groups := append([]string(nil), caseValue.Groups...)
			if !hasIncumbent || !hasCandidate {
				report.MissingPairs = append(report.MissingPairs, MissingPair{
					Key: key, Groups: groups,
					MissingIncumbent: !hasIncumbent, MissingCandidate: !hasCandidate,
				})
				continue
			}
			pair, err := buildPair(key, groups, incumbent, candidateCell)
			if err != nil {
				return nil, err
			}
			report.Pairs = append(report.Pairs, pair)
		}
	}
	report.CompletePairs = len(report.Pairs)
	report.Metrics, report.Groups, err = aggregate(report)
	if err != nil {
		return nil, err
	}
	return report, nil
}

func buildPair(key PairKey, groups []string, incumbent, candidateCell eval.Cell) (Pair, error) {
	metricNames := make(map[string]struct{}, len(incumbent.Outcome.Metrics)+len(candidateCell.Outcome.Metrics))
	for name := range incumbent.Outcome.Metrics {
		metricNames[name] = struct{}{}
	}
	for name := range candidateCell.Outcome.Metrics {
		metricNames[name] = struct{}{}
	}
	names := make([]string, 0, len(metricNames))
	for name := range metricNames {
		names = append(names, name)
	}
	sort.Strings(names)
	pair := Pair{Key: key, Groups: groups, Incumbent: incumbent, Candidate: candidateCell}
	for _, name := range names {
		incumbentValue, incumbentPresent := incumbent.Outcome.Metrics[name]
		candidateValue, candidatePresent := candidateCell.Outcome.Metrics[name]
		pair.MetricPresence = append(pair.MetricPresence, MetricPresence{
			Metric: name, IncumbentPresent: incumbentPresent, CandidatePresent: candidatePresent,
		})
		if incumbentPresent && candidatePresent {
			delta := candidateValue - incumbentValue
			if !finite(delta) {
				return Pair{}, errors.Errorf("metric %q delta for %q/%d is not finite", name, key.CaseID, key.RepeatIndex)
			}
			pair.Deltas = append(pair.Deltas, MetricDelta{
				Metric: name, Incumbent: incumbentValue, Candidate: candidateValue, Delta: delta,
			})
		}
	}
	pair.Costs = CostDelta{
		ProviderCalls: candidateCell.Outcome.ProviderCalls - incumbent.Outcome.ProviderCalls,
		ToolCalls:     candidateCell.Outcome.ToolCalls - incumbent.Outcome.ToolCalls,
		TotalTokens:   (candidateCell.Outcome.InputTokens + candidateCell.Outcome.OutputTokens) - (incumbent.Outcome.InputTokens + incumbent.Outcome.OutputTokens),
		DurationNanos: int64(candidateCell.Outcome.Duration - incumbent.Outcome.Duration),
	}
	return pair, nil
}

type metricAccumulator struct {
	aggregate MetricAggregate
	incumbent float64
	candidate float64
	delta     float64
}

type groupAccumulator struct {
	aggregate GroupAggregate
	provider  int64
	tool      int64
	tokens    int64
	duration  int64
}

func aggregate(report *Report) ([]MetricAggregate, []GroupAggregate, error) {
	expectedByGroup := map[string]int{"all": report.ExpectedPairs}
	// Expected group counts come from both complete and missing pairs.
	for _, pair := range report.Pairs {
		for _, group := range pair.Groups {
			expectedByGroup[group]++
		}
	}
	for _, missing := range report.MissingPairs {
		for _, group := range missing.Groups {
			expectedByGroup[group]++
		}
	}
	metrics := make(map[string]*metricAccumulator)
	groups := make(map[string]*groupAccumulator)
	for name, expected := range expectedByGroup {
		groups[name] = &groupAccumulator{aggregate: GroupAggregate{Group: name, ExpectedPairs: expected}}
	}
	for _, pair := range report.Pairs {
		pairGroups := append([]string{"all"}, pair.Groups...)
		for _, group := range pairGroups {
			accumulator := groups[group]
			accumulator.aggregate.CompletePairs++
			addOutcomeCounts(&accumulator.aggregate, pair)
			accumulator.provider += int64(pair.Costs.ProviderCalls)
			accumulator.tool += int64(pair.Costs.ToolCalls)
			accumulator.tokens += int64(pair.Costs.TotalTokens)
			accumulator.duration += pair.Costs.DurationNanos
			for _, delta := range pair.Deltas {
				key := strings.Join([]string{group, delta.Metric}, "\x00")
				metric := metrics[key]
				if metric == nil {
					metric = &metricAccumulator{aggregate: MetricAggregate{
						Group: group, Metric: delta.Metric, ExpectedPairs: expectedByGroup[group],
					}}
					metrics[key] = metric
				}
				metric.aggregate.PairsWithMetric++
				metric.incumbent += delta.Incumbent
				metric.candidate += delta.Candidate
				metric.delta += delta.Delta
				if !finite(metric.incumbent) || !finite(metric.candidate) || !finite(metric.delta) {
					return nil, nil, errors.Errorf("metric %q aggregate for group %q is not finite", delta.Metric, group)
				}
				switch {
				case delta.Delta > 0:
					metric.aggregate.Wins++
				case delta.Delta < 0:
					metric.aggregate.Losses++
				default:
					metric.aggregate.Ties++
				}
			}
		}
	}
	metricResults := make([]MetricAggregate, 0, len(metrics))
	for _, accumulator := range metrics {
		accumulator.aggregate.CompletePairs = groups[accumulator.aggregate.Group].aggregate.CompletePairs
		count := float64(accumulator.aggregate.PairsWithMetric)
		if count > 0 {
			accumulator.aggregate.MeanIncumbent = accumulator.incumbent / count
			accumulator.aggregate.MeanCandidate = accumulator.candidate / count
			accumulator.aggregate.MeanDelta = accumulator.delta / count
			if !finite(accumulator.aggregate.MeanIncumbent) || !finite(accumulator.aggregate.MeanCandidate) || !finite(accumulator.aggregate.MeanDelta) {
				return nil, nil, errors.Errorf("metric %q mean for group %q is not finite", accumulator.aggregate.Metric, accumulator.aggregate.Group)
			}
		}
		metricResults = append(metricResults, accumulator.aggregate)
	}
	sort.Slice(metricResults, func(i, j int) bool {
		if metricResults[i].Group != metricResults[j].Group {
			return metricResults[i].Group < metricResults[j].Group
		}
		return metricResults[i].Metric < metricResults[j].Metric
	})
	groupResults := make([]GroupAggregate, 0, len(groups))
	for _, accumulator := range groups {
		complete := float64(accumulator.aggregate.CompletePairs)
		if accumulator.aggregate.ExpectedPairs > 0 {
			accumulator.aggregate.CandidateFailureRate = float64(accumulator.aggregate.CandidateFailures) / float64(accumulator.aggregate.ExpectedPairs)
		}
		if complete > 0 {
			accumulator.aggregate.MeanProviderCallsDelta = float64(accumulator.provider) / complete
			accumulator.aggregate.MeanToolCallsDelta = float64(accumulator.tool) / complete
			accumulator.aggregate.MeanTotalTokensDelta = float64(accumulator.tokens) / complete
			accumulator.aggregate.MeanDurationNanosDelta = float64(accumulator.duration) / complete
		}
		groupResults = append(groupResults, accumulator.aggregate)
	}
	sort.Slice(groupResults, func(i, j int) bool { return groupResults[i].Group < groupResults[j].Group })
	return metricResults, groupResults, nil
}

func finite(value float64) bool { return !math.IsNaN(value) && !math.IsInf(value, 0) }

func addOutcomeCounts(aggregate *GroupAggregate, pair Pair) {
	if pair.Incumbent.Outcome.Completed {
		aggregate.IncumbentCompleted++
	}
	if pair.Candidate.Outcome.Completed {
		aggregate.CandidateCompleted++
	}
	if pair.Incumbent.Outcome.ContractValid {
		aggregate.IncumbentContractValid++
	}
	if pair.Candidate.Outcome.ContractValid {
		aggregate.CandidateContractValid++
	}
	if pair.Incumbent.Outcome.Failure != nil {
		aggregate.IncumbentFailures++
	}
	if pair.Candidate.Outcome.Failure != nil {
		aggregate.CandidateFailures++
	}
	if pair.Incumbent.Outcome.Abstained {
		aggregate.IncumbentAbstentions++
	}
	if pair.Candidate.Outcome.Abstained {
		aggregate.CandidateAbstentions++
	}
}
