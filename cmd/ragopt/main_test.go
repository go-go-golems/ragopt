package main

import (
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommandExposesNarrowCandidateValidateFlags(t *testing.T) {
	root, err := newRootCommand()
	if err != nil {
		t.Fatal(err)
	}
	validate, _, err := root.Find([]string{"candidate", "validate"})
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"bundle", "manifest", "format", "output-fields", "max-output-rows"} {
		if validate.Flags().Lookup(required) == nil {
			t.Errorf("missing validate flag --%s", required)
		}
	}
	for _, unwanted := range []string{"config-file", "print-parsed-fields", "print-schema", "print-yaml"} {
		if validate.Flags().Lookup(unwanted) != nil {
			t.Errorf("unexpected automatic flag --%s", unwanted)
		}
	}
	if root.PersistentFlags().Lookup("log-level") == nil {
		t.Fatal("root command is missing --log-level")
	}
}

func TestRootCandidateValidateRequiresBundle(t *testing.T) {
	root, err := newRootCommand()
	if err != nil {
		t.Fatal(err)
	}
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	root.SetArgs([]string{"candidate", "validate"})
	err = root.Execute()
	if err == nil || !strings.Contains(err.Error(), "bundle") {
		t.Fatalf("expected required --bundle error, got %v", err)
	}
}

func TestRootExposesNarrowCompareAndReportFlags(t *testing.T) {
	root, err := newRootCommand()
	if err != nil {
		t.Fatal(err)
	}
	compareCommand, _, err := root.Find([]string{"compare"})
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"run", "format", "output-fields", "max-output-rows"} {
		if compareCommand.Flags().Lookup(required) == nil {
			t.Errorf("compare missing --%s", required)
		}
	}
	reportCommand, _, err := root.Find([]string{"report"})
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"run", "output-path", "plan-path", "format"} {
		if reportCommand.Flags().Lookup(required) == nil {
			t.Errorf("report missing --%s", required)
		}
	}
	for _, command := range []*cobra.Command{compareCommand, reportCommand} {
		for _, unwanted := range []string{"config-file", "print-parsed-fields", "print-schema", "print-yaml"} {
			if command.Flags().Lookup(unwanted) != nil {
				t.Errorf("%s has unexpected automatic flag --%s", command.Name(), unwanted)
			}
		}
	}
}
