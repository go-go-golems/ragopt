package gate

import (
	"context"
	"fmt"
	"sort"

	"github.com/pkg/errors"

	"github.com/go-go-golems/ragopt/pkg/compare"
)

const DecisionAPIVersion = "ragopt-gate-decision/v1"

type DecisionStatus string

const (
	DecisionPass DecisionStatus = "pass"
	DecisionFail DecisionStatus = "fail"
)

type CheckResult struct {
	Phase   string         `json:"phase"`
	Name    string         `json:"name"`
	Passed  bool           `json:"passed"`
	Message string         `json:"message"`
	Values  map[string]any `json:"values,omitempty"`
}

type Decision struct {
	APIVersion   string         `json:"api_version"`
	PolicyName   string         `json:"policy_name"`
	PolicyDigest string         `json:"policy_digest"`
	Status       DecisionStatus `json:"status"`
	Checks       []CheckResult  `json:"checks"`
	Reasons      []string       `json:"reasons,omitempty"`
}

// Evaluate applies identity, hard, target, regression, and tie-break phases in
// lexicographic order. It is pure and performs no I/O.
func Evaluate(ctx context.Context, policy *PolicyDocument, report *compare.Report) (Decision, error) {
	if err := ctx.Err(); err != nil {
		return Decision{}, err
	}
	if policy == nil || report == nil {
		return Decision{}, errors.New("gate policy and comparison report are required")
	}
	semanticDigest, err := policyDigest(policy.Policy)
	if err != nil {
		return Decision{}, errors.Wrap(err, "validate gate policy semantics")
	}
	if semanticDigest != policy.Digest {
		return Decision{}, errors.Errorf("gate policy semantic digest mismatch: document=%s actual=%s", policy.Digest, semanticDigest)
	}
	decision := Decision{
		APIVersion: DecisionAPIVersion, PolicyName: policy.Policy.Name,
		PolicyDigest: policy.Digest, Status: DecisionPass,
	}
	identity := []CheckResult{
		check("identity", "policy_bytes", policy.ByteDigest == report.PolicyDigest,
			"run policy bytes match the loaded gate policy", map[string]any{"run": report.PolicyDigest, "loaded": policy.ByteDigest}),
		check("identity", "complete_pairing", report.CompletePairs == report.ExpectedPairs && len(report.MissingPairs) == 0,
			fmt.Sprintf("complete pairs %d of %d", report.CompletePairs, report.ExpectedPairs), nil),
	}
	if stopAfter(&decision, identity) {
		return decision, nil
	}

	all := findGroup(report, "all")
	hard := make([]CheckResult, 0)
	if policy.Policy.HardGates.RequireAllCells {
		hard = append(hard, check("hard", "require_all_cells", report.CompletePairs == report.ExpectedPairs,
			fmt.Sprintf("complete pairs %d of %d", report.CompletePairs, report.ExpectedPairs), nil))
	}
	if policy.Policy.HardGates.RequireCompleted {
		hard = append(hard, check("hard", "require_completed", all != nil && all.CandidateCompleted == all.ExpectedPairs,
			fmt.Sprintf("candidate completed %d of %d", valueInt(all, func(group *compare.GroupAggregate) int { return group.CandidateCompleted }), report.ExpectedPairs), nil))
	}
	if policy.Policy.HardGates.RequireContractValid {
		hard = append(hard, check("hard", "require_contract_valid", all != nil && all.CandidateContractValid == all.ExpectedPairs,
			fmt.Sprintf("candidate contract-valid %d of %d", valueInt(all, func(group *compare.GroupAggregate) int { return group.CandidateContractValid }), report.ExpectedPairs), nil))
	}
	failureRate := 1.0
	if all != nil {
		failureRate = all.CandidateFailureRate
	}
	maxFailure := *policy.Policy.HardGates.MaxFailureRate
	hard = append(hard, check("hard", "max_failure_rate", failureRate <= maxFailure,
		fmt.Sprintf("candidate failure rate %.6f <= %.6f", failureRate, maxFailure), map[string]any{"actual": failureRate, "maximum": maxFailure}))
	floorMetrics := sortedMetricNames(policy.Policy.HardGates.MetricFloors)
	for _, metric := range floorMetrics {
		floor := policy.Policy.HardGates.MetricFloors[metric]
		passed, minimum, present := candidateMetricFloor(report.Pairs, metric, floor)
		hard = append(hard, check("hard", "metric_floor:"+metric, passed,
			fmt.Sprintf("candidate %s minimum %.6f across %d pairs; floor %.6f", metric, minimum, present, floor),
			map[string]any{"minimum": minimum, "pairs_with_metric": present, "expected_pairs": report.ExpectedPairs, "floor": floor}))
	}
	if stopAfter(&decision, hard) {
		return decision, nil
	}

	targetChecks := evaluateTarget(policy.Policy.Target, report)
	if stopAfter(&decision, targetChecks) {
		return decision, nil
	}
	regressionChecks := evaluateRegressions(policy.Policy.Regressions, report)
	if stopAfter(&decision, regressionChecks) {
		return decision, nil
	}
	decision.Checks = append(decision.Checks, evaluateTieBreakers(policy.Policy.TieBreakers, all)...)
	return decision, nil
}

func evaluateTarget(target Target, report *compare.Report) []CheckResult {
	selectedGroups := groupSet(target.Groups)
	deltas := make([]struct {
		repeat int
		value  float64
	}, 0)
	selectedPairs := 0
	missingMetric := 0
	for _, pair := range report.Pairs {
		if !containsGroup(pair.Groups, selectedGroups) {
			continue
		}
		selectedPairs++
		delta, ok := pairDelta(pair, target.Metric)
		if !ok {
			missingMetric++
			continue
		}
		deltas = append(deltas, struct {
			repeat int
			value  float64
		}{repeat: pair.Key.RepeatIndex, value: delta})
	}
	targetMean := 0.0
	for _, delta := range deltas {
		targetMean += delta.value
	}
	if len(deltas) > 0 {
		targetMean /= float64(len(deltas))
	}
	checks := []CheckResult{
		check("target", "metric_presence:"+target.Metric, selectedPairs > 0 && missingMetric == 0 && len(deltas) == selectedPairs,
			fmt.Sprintf("target metric present in %d of %d selected pairs", len(deltas), selectedPairs), nil),
		check("target", "minimum_mean_delta:"+target.Metric, selectedPairs > 0 && missingMetric == 0 && targetMean >= target.MinimumMeanDelta,
			fmt.Sprintf("target mean delta %.6f >= %.6f", targetMean, target.MinimumMeanDelta), map[string]any{"actual": targetMean, "minimum": target.MinimumMeanDelta}),
	}
	if target.RequirePositiveEachRepeat {
		byRepeat := map[int][]float64{}
		for _, delta := range deltas {
			byRepeat[delta.repeat] = append(byRepeat[delta.repeat], delta.value)
		}
		repeats := make([]int, 0, len(byRepeat))
		for repeat := range byRepeat {
			repeats = append(repeats, repeat)
		}
		sort.Ints(repeats)
		passed := selectedPairs > 0 && missingMetric == 0
		means := map[string]float64{}
		for _, repeat := range repeats {
			value := mean(byRepeat[repeat])
			means[fmt.Sprintf("%d", repeat)] = value
			if value <= 0 {
				passed = false
			}
		}
		checks = append(checks, check("target", "positive_each_repeat:"+target.Metric, passed,
			"target mean delta must be positive in every represented repeat", map[string]any{"repeat_means": means}))
	}
	return checks
}

func evaluateRegressions(regressions Regressions, report *compare.Report) []CheckResult {
	checks := make([]CheckResult, 0)
	for _, metric := range sortedMetricNames(regressions.MaximumCaseDelta) {
		minimum := regressions.MaximumCaseDelta[metric]
		passed := true
		present := 0
		worst := 0.0
		initialized := false
		for _, pair := range report.Pairs {
			delta, ok := pairDelta(pair, metric)
			if !ok {
				passed = false
				continue
			}
			present++
			if !initialized || delta < worst {
				worst, initialized = delta, true
			}
			if delta < minimum {
				passed = false
			}
		}
		passed = passed && present == report.ExpectedPairs
		checks = append(checks, check("regression", "maximum_case_delta:"+metric, passed,
			fmt.Sprintf("worst case delta %.6f >= %.6f with metric in %d of %d pairs", worst, minimum, present, report.ExpectedPairs), nil))
	}
	groups := make([]string, 0, len(regressions.MaximumMeanDelta))
	for group := range regressions.MaximumMeanDelta {
		groups = append(groups, group)
	}
	sort.Strings(groups)
	for _, group := range groups {
		for _, metric := range sortedMetricNames(regressions.MaximumMeanDelta[group]) {
			minimum := regressions.MaximumMeanDelta[group][metric]
			aggregate := findMetric(report, group, metric)
			passed := aggregate != nil && aggregate.PairsWithMetric == aggregate.ExpectedPairs && aggregate.MeanDelta >= minimum
			actual, present, expected := 0.0, 0, 0
			if aggregate != nil {
				actual, present, expected = aggregate.MeanDelta, aggregate.PairsWithMetric, aggregate.ExpectedPairs
			}
			checks = append(checks, check("regression", "maximum_mean_delta:"+group+":"+metric, passed,
				fmt.Sprintf("%s mean delta %.6f >= %.6f with metric in %d of %d pairs", group, actual, minimum, present, expected), nil))
		}
	}
	return checks
}

func evaluateTieBreakers(tieBreakers []string, all *compare.GroupAggregate) []CheckResult {
	checks := make([]CheckResult, 0, len(tieBreakers))
	for _, name := range tieBreakers {
		value := 0.0
		if all != nil {
			switch name {
			case "provider_calls":
				value = all.MeanProviderCallsDelta
			case "tool_calls":
				value = all.MeanToolCallsDelta
			case "total_tokens":
				value = all.MeanTotalTokensDelta
			case "duration":
				value = all.MeanDurationNanosDelta
			}
		}
		checks = append(checks, check("tie_break", name, true,
			fmt.Sprintf("candidate minus incumbent mean %s delta %.6f; lower is preferred only after quality gates", name, value), map[string]any{"delta": value}))
	}
	return checks
}

func stopAfter(decision *Decision, checks []CheckResult) bool {
	decision.Checks = append(decision.Checks, checks...)
	failed := false
	for _, checkResult := range checks {
		if !checkResult.Passed {
			failed = true
			decision.Reasons = append(decision.Reasons, checkResult.Message)
		}
	}
	if failed {
		decision.Status = DecisionFail
	}
	return failed
}

func check(phase, name string, passed bool, message string, values map[string]any) CheckResult {
	return CheckResult{Phase: phase, Name: name, Passed: passed, Message: message, Values: values}
}

func findGroup(report *compare.Report, group string) *compare.GroupAggregate {
	for index := range report.Groups {
		if report.Groups[index].Group == group {
			return &report.Groups[index]
		}
	}
	return nil
}

func findMetric(report *compare.Report, group, metric string) *compare.MetricAggregate {
	for index := range report.Metrics {
		if report.Metrics[index].Group == group && report.Metrics[index].Metric == metric {
			return &report.Metrics[index]
		}
	}
	return nil
}

func pairDelta(pair compare.Pair, metric string) (float64, bool) {
	for _, delta := range pair.Deltas {
		if delta.Metric == metric {
			return delta.Delta, true
		}
	}
	return 0, false
}

func candidateMetricFloor(pairs []compare.Pair, metric string, floor float64) (bool, float64, int) {
	minimum := 0.0
	present := 0
	passed := true
	initialized := false
	for _, pair := range pairs {
		value, ok := pair.Candidate.Outcome.Metrics[metric]
		if !ok {
			passed = false
			continue
		}
		present++
		if !initialized || value < minimum {
			minimum, initialized = value, true
		}
		if value < floor {
			passed = false
		}
	}
	return passed && present == len(pairs) && len(pairs) > 0, minimum, present
}

func sortedMetricNames(values map[string]float64) []string {
	result := make([]string, 0, len(values))
	for name := range values {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value
	}
	return total / float64(len(values))
}

func valueInt(group *compare.GroupAggregate, fn func(*compare.GroupAggregate) int) int {
	if group == nil {
		return 0
	}
	return fn(group)
}
