package compare

import (
	"context"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/pkg/errors"

	comparelib "github.com/go-go-golems/ragopt/pkg/compare"
	"github.com/go-go-golems/ragopt/pkg/eval"
	"github.com/go-go-golems/ragopt/pkg/gate"
	"github.com/go-go-golems/ragopt/pkg/policy"
)

type Command struct{ *cmds.CommandDescription }

type settings struct {
	Run string `glazed:"run"`
}

var _ cmds.GlazeCommand = (*Command)(nil)

func NewCommand() (*Command, error) {
	return &Command{CommandDescription: cmds.NewCommandDescription(
		"compare",
		cmds.WithShort("Compare and gate one existing paired evaluation run"),
		cmds.WithLong(`Strictly load an immutable evaluation run, pair incumbent and candidate cells, and evaluate its copied gate policy.

The command never repairs the run or writes reports. It emits decision, check,
metric, group, and missing-pair records through Glazed structured output.

Example:
  ragopt compare --run ./runs/20260806T... --format json
`),
		cmds.WithFlags(fields.New("run", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Existing paired evaluation run directory"))),
	)}, nil
}

func (command *Command) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, processor middlewares.Processor) error {
	configuration := &settings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, configuration); err != nil {
		return errors.Wrap(err, "decode compare settings")
	}
	run, err := eval.LoadArtifactRun(ctx, configuration.Run)
	if err != nil {
		return errors.Wrap(err, "load evaluation artifacts")
	}
	policyDocument, err := policy.Load(ctx, run.PolicyPath)
	if err != nil {
		return errors.Wrap(err, "load copied gate policy")
	}
	comparison, err := comparelib.Build(ctx, run)
	if err != nil {
		return errors.Wrap(err, "compare paired outcomes")
	}
	decision, err := gate.Evaluate(ctx, policyDocument, comparison)
	if err != nil {
		return errors.Wrap(err, "evaluate gate policy")
	}
	for _, row := range Rows(comparison, decision) {
		if err := processor.AddRow(ctx, row); err != nil {
			return errors.Wrap(err, "emit comparison row")
		}
	}
	return nil
}

func Rows(comparison *comparelib.Report, decision gate.Decision) []types.Row {
	rows := []types.Row{types.NewRow(
		types.MRP("record_type", "decision"), types.MRP("run_id", comparison.RunID),
		types.MRP("candidate_id", comparison.CandidateID), types.MRP("decision", string(decision.Status)),
		types.MRP("expected_pairs", comparison.ExpectedPairs), types.MRP("complete_pairs", comparison.CompletePairs),
	)}
	for _, check := range decision.Checks {
		rows = append(rows, types.NewRow(
			types.MRP("record_type", "check"), types.MRP("phase", check.Phase), types.MRP("name", check.Name),
			types.MRP("passed", check.Passed), types.MRP("message", check.Message),
		))
	}
	for _, group := range comparison.Groups {
		rows = append(rows, types.NewRow(
			types.MRP("record_type", "group"), types.MRP("group", group.Group),
			types.MRP("expected_pairs", group.ExpectedPairs), types.MRP("complete_pairs", group.CompletePairs),
			types.MRP("candidate_completed", group.CandidateCompleted), types.MRP("candidate_contract_valid", group.CandidateContractValid),
			types.MRP("candidate_failures", group.CandidateFailures), types.MRP("candidate_failure_rate", group.CandidateFailureRate),
		))
	}
	for _, metric := range comparison.Metrics {
		rows = append(rows, types.NewRow(
			types.MRP("record_type", "metric"), types.MRP("group", metric.Group), types.MRP("metric", metric.Metric),
			types.MRP("pairs_with_metric", metric.PairsWithMetric), types.MRP("mean_incumbent", metric.MeanIncumbent),
			types.MRP("mean_candidate", metric.MeanCandidate), types.MRP("mean_delta", metric.MeanDelta),
			types.MRP("wins", metric.Wins), types.MRP("ties", metric.Ties), types.MRP("losses", metric.Losses),
		))
	}
	for _, missing := range comparison.MissingPairs {
		rows = append(rows, types.NewRow(
			types.MRP("record_type", "missing_pair"), types.MRP("case_id", missing.Key.CaseID), types.MRP("repeat_index", missing.Key.RepeatIndex),
			types.MRP("missing_incumbent", missing.MissingIncumbent), types.MRP("missing_candidate", missing.MissingCandidate),
		))
	}
	return rows
}
