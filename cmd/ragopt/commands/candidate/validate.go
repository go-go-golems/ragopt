package candidate

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

	candidatelib "github.com/go-go-golems/ragopt/pkg/candidate"
)

// ValidateCommand validates one candidate bundle and emits its semantic
// identities as a structured row.
type ValidateCommand struct {
	*cmds.CommandDescription
}

type validateSettings struct {
	Bundle   string `glazed:"bundle"`
	Manifest string `glazed:"manifest"`
}

var _ cmds.GlazeCommand = (*ValidateCommand)(nil)

// NewValidateCommand constructs the Glazed candidate validate command.
func NewValidateCommand() (*ValidateCommand, error) {
	description := cmds.NewCommandDescription(
		"validate",
		cmds.WithShort("Validate an exactly-one-mutation candidate bundle"),
		cmds.WithLong(`Validate candidate.yaml, both snapshot manifests, and every declared asset.

The command rejects unknown YAML fields, path and symlink escapes, digest or
size drift, locked-input changes, semantic-dimension changes, and anything
other than exactly one mutable asset byte change. A successful command emits
one structured identity row and does not modify the bundle.

Examples:
  ragopt candidate validate --bundle ./candidate-prompt-001
  ragopt candidate validate --bundle ./candidate-prompt-001 --format json
`),
		cmds.WithFlags(
			fields.New(
				"bundle",
				fields.TypeString,
				fields.WithRequired(true),
				fields.WithHelp("Candidate bundle root directory"),
			),
			fields.New(
				"manifest",
				fields.TypeString,
				fields.WithDefault("candidate.yaml"),
				fields.WithHelp("Candidate manifest path relative to the bundle root"),
			),
		),
	)
	return &ValidateCommand{CommandDescription: description}, nil
}

// RunIntoGlazeProcessor validates the bundle and emits one result row.
func (command *ValidateCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	vals *values.Values,
	processor middlewares.Processor,
) error {
	settings := &validateSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return errors.Wrap(err, "decode candidate validation settings")
	}
	loaded, err := candidatelib.LoadCandidate(ctx, settings.Bundle, settings.Manifest)
	if err != nil {
		return errors.Wrap(err, "validate candidate bundle")
	}
	row := ValidationRow(loaded)
	if err := processor.AddRow(ctx, row); err != nil {
		return errors.Wrap(err, "emit candidate validation row")
	}
	return nil
}

// ValidationRow projects a validated candidate into the stable CLI row.
func ValidationRow(loaded *candidatelib.Candidate) types.Row {
	return types.NewRow(
		types.MRP("candidate_id", loaded.Manifest.CandidateID),
		types.MRP("candidate_digest", loaded.Digest),
		types.MRP("parent_snapshot", loaded.Parent.SnapshotID),
		types.MRP("child_snapshot", loaded.Child.SnapshotID),
		types.MRP("changed_asset", loaded.Mutation.AssetName),
		types.MRP("parent_asset_digest", loaded.Mutation.ParentDigest),
		types.MRP("child_asset_digest", loaded.Mutation.ChildDigest),
		types.MRP("bundle", filepath.Clean(loaded.Root)),
		types.MRP("valid", true),
	)
}
