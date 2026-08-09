package report

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/go-go-golems/ragopt/pkg/compare"
	"github.com/go-go-golems/ragopt/pkg/eval"
	"github.com/go-go-golems/ragopt/pkg/gate"
)

// Build creates deterministic report content without writing any files.
func Build(ctx context.Context, run *eval.ArtifactRun, comparison *compare.Report, policy *gate.PolicyDocument, decision gate.Decision) (*Document, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if run == nil || comparison == nil || policy == nil {
		return nil, errors.New("artifact run, comparison, and policy are required")
	}
	recomputed, err := gate.Evaluate(ctx, policy, comparison)
	if err != nil {
		return nil, errors.Wrap(err, "recompute gate decision")
	}
	if !reflect.DeepEqual(decision, recomputed) {
		return nil, errors.New("supplied gate decision differs from recomputed comparison decision")
	}
	if decision.APIVersion != gate.DecisionAPIVersion ||
		decision.PolicyDigest != policy.Digest ||
		policy.ByteDigest != run.Config.PolicyDigest ||
		comparison.RunID != run.Manifest.RunID ||
		comparison.SuiteDigest != run.Config.SuiteDigest ||
		comparison.PolicyDigest != run.Config.PolicyDigest ||
		comparison.CandidateID != run.Config.CandidateID ||
		comparison.CandidateDigest != run.Config.CandidateDigest ||
		comparison.ParentSnapshot != run.Config.ParentSnapshot ||
		comparison.ChildSnapshot != run.Config.ChildSnapshot ||
		comparison.IncumbentArm != run.Config.IncumbentArm ||
		comparison.ChallengerArm != run.Config.ChallengerArm {
		return nil, errors.New("report inputs have inconsistent identities")
	}
	plan := PromotionPlan{
		APIVersion: PromotionPlanAPIVersion, State: "review_required", HumanApplyRequired: true,
		RunID: comparison.RunID, CandidateID: comparison.CandidateID, CandidateDigest: comparison.CandidateDigest,
		ParentSnapshot: comparison.ParentSnapshot, ChildSnapshot: comparison.ChildSnapshot,
		Mutation: run.Config.Mutation, ChangedAsset: run.Config.ChangedAsset,
		ParentAssetDigest: run.Config.ParentAssetDigest, ChildAssetDigest: run.Config.ChildAssetDigest,
		PolicyName: policy.Policy.Name, PolicyDigest: policy.Digest,
		Decision: decision.Status, Reasons: append([]string(nil), decision.Reasons...),
	}
	return &Document{Markdown: renderMarkdown(run, comparison, policy, decision), Plan: plan}, nil
}

func renderMarkdown(run *eval.ArtifactRun, comparison *compare.Report, policy *gate.PolicyDocument, decision gate.Decision) string {
	var output strings.Builder
	fmt.Fprintf(&output, "# Promotion review: %s\n\n", comparison.CandidateID)
	fmt.Fprintf(&output, "> Decision: **%s**. This report does not apply the candidate; human review is required.\n\n", strings.ToUpper(string(decision.Status)))
	output.WriteString("## Evidence identity\n\n")
	fmt.Fprintf(&output, "- Run: `%s`\n- Suite: `%s`\n- Policy: `%s` (`%s`)\n", comparison.RunID, comparison.SuiteDigest, policy.Policy.Name, policy.Digest)
	fmt.Fprintf(&output, "- Candidate: `%s` (`%s`)\n- Parent snapshot: `%s`\n- Child snapshot: `%s`\n\n", comparison.CandidateID, comparison.CandidateDigest, comparison.ParentSnapshot, comparison.ChildSnapshot)
	output.WriteString("## Proposed mutation\n\n")
	fmt.Fprintf(&output, "- Asset: `%s`\n- Parent bytes: `%s`\n- Candidate bytes: `%s`\n", run.Config.ChangedAsset, run.Config.ParentAssetDigest, run.Config.ChildAssetDigest)
	fmt.Fprintf(&output, "- Hypothesis: %s\n- Expected improvement: `%s`", clean(run.Config.Mutation.Hypothesis), run.Config.Mutation.ExpectedImprovement.Metric)
	if len(run.Config.Mutation.ExpectedImprovement.Groups) > 0 {
		fmt.Fprintf(&output, " in `%s`", strings.Join(run.Config.Mutation.ExpectedImprovement.Groups, "`, `"))
	}
	output.WriteString("\n")
	if len(run.Config.Mutation.RegressionRisks) > 0 {
		fmt.Fprintf(&output, "- Declared regression risks: %s\n", clean(strings.Join(run.Config.Mutation.RegressionRisks, "; ")))
	}
	output.WriteString("\n## Gate decision\n\n")
	output.WriteString("| Phase | Check | Result | Detail |\n|---|---|---:|---|\n")
	for _, check := range decision.Checks {
		result := "FAIL"
		if check.Passed {
			result = "PASS"
		}
		fmt.Fprintf(&output, "| %s | %s | %s | %s |\n", table(check.Phase), table(check.Name), result, table(check.Message))
	}
	output.WriteString("\n## Group outcomes\n\n")
	output.WriteString("| Group | Pairs | Candidate completed | Contract valid | Failures | Abstentions | Failure rate |\n|---|---:|---:|---:|---:|---:|---:|\n")
	for _, group := range comparison.Groups {
		fmt.Fprintf(&output, "| %s | %d/%d | %d | %d | %d | %d | %.4f |\n", table(group.Group), group.CompletePairs, group.ExpectedPairs, group.CandidateCompleted, group.CandidateContractValid, group.CandidateFailures, group.CandidateAbstentions, group.CandidateFailureRate)
	}
	output.WriteString("\n## Metric aggregates\n\n")
	output.WriteString("| Group | Metric | Present | Incumbent mean | Candidate mean | Mean delta | Wins / ties / losses |\n|---|---|---:|---:|---:|---:|---:|\n")
	for _, metric := range comparison.Metrics {
		fmt.Fprintf(&output, "| %s | %s | %d/%d | %.6f | %.6f | %+.6f | %d / %d / %d |\n", table(metric.Group), table(metric.Metric), metric.PairsWithMetric, metric.ExpectedPairs, metric.MeanIncumbent, metric.MeanCandidate, metric.MeanDelta, metric.Wins, metric.Ties, metric.Losses)
	}
	output.WriteString("\n## Missing pairs\n\n")
	if len(comparison.MissingPairs) == 0 {
		output.WriteString("None.\n")
	} else {
		output.WriteString("| Case | Repeat | Missing incumbent | Missing candidate |\n|---|---:|---:|---:|\n")
		for _, missing := range comparison.MissingPairs {
			fmt.Fprintf(&output, "| %s | %d | %t | %t |\n", table(missing.Key.CaseID), missing.Key.RepeatIndex, missing.MissingIncumbent, missing.MissingCandidate)
		}
	}
	output.WriteString("\n## Paired results\n\n")
	metricNames := pairedMetricNames(comparison.Pairs)
	output.WriteString("| Case | Repeat | Groups | Candidate completion | Candidate contract | Candidate failure | Metric deltas | Cost deltas |\n|---|---:|---|---:|---:|---|---|---|\n")
	for _, pair := range comparison.Pairs {
		deltas := make([]string, 0, len(metricNames))
		for _, name := range metricNames {
			if delta, ok := deltaFor(pair, name); ok {
				deltas = append(deltas, fmt.Sprintf("%s=%+.6f", name, delta))
			} else {
				deltas = append(deltas, name+"=missing")
			}
		}
		failure := ""
		if pair.Candidate.Outcome.Failure != nil {
			failure = pair.Candidate.Outcome.Failure.Class + ": " + pair.Candidate.Outcome.Failure.Message
		}
		costs := fmt.Sprintf("provider=%+d; tool=%+d; tokens=%+d; duration_ns=%+d", pair.Costs.ProviderCalls, pair.Costs.ToolCalls, pair.Costs.TotalTokens, pair.Costs.DurationNanos)
		fmt.Fprintf(&output, "| %s | %d | %s | %t | %t | %s | %s | %s |\n", table(pair.Key.CaseID), pair.Key.RepeatIndex, table(strings.Join(pair.Groups, ", ")), pair.Candidate.Outcome.Completed, pair.Candidate.Outcome.ContractValid, table(failure), table(strings.Join(deltas, "; ")), table(costs))
	}
	output.WriteString("\n## Review action\n\n")
	if decision.Status == gate.DecisionPass {
		output.WriteString("The evidence gates passed. A human may review the native artifacts and apply the described asset replacement in the product repository.\n")
	} else {
		output.WriteString("Do not promote this candidate. Retain the run and candidate as rejection evidence.\n")
	}
	return output.String()
}

func pairedMetricNames(pairs []compare.Pair) []string {
	set := map[string]struct{}{}
	for _, pair := range pairs {
		for _, presence := range pair.MetricPresence {
			set[presence.Metric] = struct{}{}
		}
	}
	result := make([]string, 0, len(set))
	for name := range set {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

func deltaFor(pair compare.Pair, metric string) (float64, bool) {
	for _, delta := range pair.Deltas {
		if delta.Metric == metric {
			return delta.Delta, true
		}
	}
	return 0, false
}

func clean(value string) string { return strings.Join(strings.Fields(value), " ") }

func table(value string) string { return strings.ReplaceAll(clean(value), "|", "\\|") }
