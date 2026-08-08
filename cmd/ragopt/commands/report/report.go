package report

import (
	"context"
	"path/filepath"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/pkg/errors"

	"github.com/go-go-golems/ragopt/pkg/compare"
	"github.com/go-go-golems/ragopt/pkg/eval"
	"github.com/go-go-golems/ragopt/pkg/gate"
	reportlib "github.com/go-go-golems/ragopt/pkg/report"
)

type Command struct{ *cmds.CommandDescription }

type settings struct {
	Run        string `glazed:"run"`
	OutputPath string `glazed:"output-path"`
	PlanPath   string `glazed:"plan-path"`
}

var _ cmds.GlazeCommand = (*Command)(nil)

func NewCommand() (*Command, error) {
	return &Command{CommandDescription: cmds.NewCommandDescription(
		"report",
		cmds.WithShort("Render a promotion review and non-applying plan"),
		cmds.WithLong(`Strictly compare and gate an existing run, then atomically write a Markdown review and JSON promotion plan to explicit paths.

The plan is always marked review_required and human_apply_required. This
command has no operation that changes a product or evaluated run.

Example:
  ragopt report --run ./runs/20260806T... --output-path ./review.md --plan-path ./promotion-plan.json
`),
		cmds.WithFlags(
			fields.New("run", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Existing paired evaluation run directory")),
			fields.New("output-path", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Markdown promotion review output path")),
			fields.New("plan-path", fields.TypeString, fields.WithRequired(true), fields.WithHelp("JSON promotion plan output path")),
		),
	)}, nil
}

func (command *Command) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, processor middlewares.Processor) error {
	configuration := &settings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, configuration); err != nil {
		return errors.Wrap(err, "decode report settings")
	}
	run, err := eval.LoadArtifactRun(ctx, configuration.Run)
	if err != nil {
		return errors.Wrap(err, "load evaluation artifacts")
	}
	policy, err := gate.LoadPolicy(ctx, run.PolicyPath)
	if err != nil {
		return errors.Wrap(err, "load copied gate policy")
	}
	comparison, err := compare.Build(ctx, run)
	if err != nil {
		return errors.Wrap(err, "compare paired outcomes")
	}
	decision, err := gate.Evaluate(ctx, policy, comparison)
	if err != nil {
		return errors.Wrap(err, "evaluate gate policy")
	}
	document, err := reportlib.Build(ctx, run, comparison, policy, decision)
	if err != nil {
		return errors.Wrap(err, "build promotion report")
	}
	if err := reportlib.ValidateOutputsOutsideRun(run.Directory, configuration.OutputPath, configuration.PlanPath); err != nil {
		return err
	}
	if err := reportlib.Write(ctx, document, configuration.OutputPath, configuration.PlanPath); err != nil {
		return err
	}
	row := types.NewRow(
		types.MRP("run_id", comparison.RunID), types.MRP("candidate_id", comparison.CandidateID),
		types.MRP("decision", string(decision.Status)), types.MRP("report_path", cleanPath(configuration.OutputPath)),
		types.MRP("plan_path", cleanPath(configuration.PlanPath)), types.MRP("human_apply_required", true),
	)
	return errors.Wrap(processor.AddRow(ctx, row), "emit report row")
}

func cleanPath(path string) string {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}
	return absolute
}
