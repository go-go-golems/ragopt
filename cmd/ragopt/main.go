package main

import (
	"fmt"
	"os"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	candidatecmd "github.com/go-go-golems/ragopt/cmd/ragopt/commands/candidate"
)

var version = "dev"

func newRootCommand() (*cobra.Command, error) {
	root := &cobra.Command{
		Use:          "ragopt",
		Short:        "Evidence-gated optimization artifact tools",
		Version:      version,
		SilenceUsage: true,
		PersistentPreRunE: func(command *cobra.Command, _ []string) error {
			return logging.InitLoggerFromCobra(command)
		},
	}
	if err := logging.AddLoggingSectionToRootCommand(root, "ragopt"); err != nil {
		return nil, errors.Wrap(err, "add logging flags")
	}

	candidateParent := &cobra.Command{
		Use:   "candidate",
		Short: "Validate and inspect candidate bundles",
	}
	validateCommand, err := candidatecmd.NewValidateCommand()
	if err != nil {
		return nil, err
	}
	validateCobra, err := cli.BuildCobraCommandFromCommand(
		validateCommand,
		cli.WithParserConfig(cli.CobraParserConfig{
			ShortHelpSections:          []string{schema.DefaultSlug},
			SkipCommandSettingsSection: true,
		}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "build candidate validate command")
	}
	candidateParent.AddCommand(validateCobra)
	root.AddCommand(candidateParent)
	return root, nil
}

func main() {
	root, err := newRootCommand()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
