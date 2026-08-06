package candidate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/types"
	"gopkg.in/yaml.v3"

	candidatelib "github.com/go-go-golems/ragopt/pkg/candidate"
)

func TestValidateCommandEmitsIdentityRow(t *testing.T) {
	bundle := writeCandidateBundle(t)
	command, err := NewValidateCommand()
	mustNoError(t, err)
	section, ok := command.GetDefaultSection()
	if !ok {
		t.Fatal("validate command has no default section")
	}
	sectionValues, err := values.NewSectionValues(
		section,
		values.WithFieldValue("bundle", bundle),
		values.WithFieldValue("manifest", "candidate.yaml"),
	)
	mustNoError(t, err)
	parsed := values.New()
	parsed.Set(schema.DefaultSlug, sectionValues)
	collector := &rowCollector{}
	mustNoError(t, command.RunIntoGlazeProcessor(t.Context(), parsed, collector))
	if len(collector.rows) != 1 {
		t.Fatalf("rows: got %d", len(collector.rows))
	}
	assertRowValue(t, collector.rows[0], "candidate_id", "prompt-001")
	assertRowValue(t, collector.rows[0], "changed_asset", "prompt")
	assertRowValue(t, collector.rows[0], "valid", true)
	value, ok := collector.rows[0].Get("candidate_digest")
	if !ok || !strings.HasPrefix(value.(string), "sha256:") {
		t.Fatalf("candidate digest: %v", value)
	}
}

func TestValidateCommandPropagatesValidationError(t *testing.T) {
	command, err := NewValidateCommand()
	mustNoError(t, err)
	section, ok := command.GetDefaultSection()
	if !ok {
		t.Fatal("validate command has no default section")
	}
	sectionValues, err := values.NewSectionValues(
		section,
		values.WithFieldValue("bundle", filepath.Join(t.TempDir(), "missing")),
		values.WithFieldValue("manifest", "candidate.yaml"),
	)
	mustNoError(t, err)
	parsed := values.New()
	parsed.Set(schema.DefaultSlug, sectionValues)
	err = command.RunIntoGlazeProcessor(t.Context(), parsed, &rowCollector{})
	if err == nil || !strings.Contains(err.Error(), "validate candidate bundle") {
		t.Fatalf("expected wrapped candidate error, got %v", err)
	}
}

type rowCollector struct {
	rows []types.Row
}

func (collector *rowCollector) AddRow(_ context.Context, row types.Row) error {
	collector.rows = append(collector.rows, row)
	return nil
}

func (collector *rowCollector) Close(context.Context) error { return nil }

func assertRowValue(t *testing.T, row types.Row, key string, expected any) {
	t.Helper()
	actual, ok := row.Get(key)
	if !ok || actual != expected {
		t.Fatalf("row %s: got %#v, expected %#v", key, actual, expected)
	}
}

func writeCandidateBundle(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	parentBytes := []byte("parent prompt\n")
	childBytes := []byte("candidate prompt\n")
	writeFile(t, root, "parent/prompt.md", parentBytes)
	writeFile(t, root, "candidate/prompt.md", childBytes)
	parent := candidatelib.Snapshot{
		APIVersion: candidatelib.SnapshotAPIVersion,
		System:     "cli-fixture",
		MutableAssets: []candidatelib.AssetRef{{
			Name:      "prompt",
			MediaType: "text/markdown",
			Path:      "parent/prompt.md",
			SHA256:    sha256Text(parentBytes),
			SizeBytes: int64(len(parentBytes)),
		}},
		Dimensions: map[string]string{"model": "fixture-model"},
	}
	child := candidatelib.Snapshot{
		APIVersion: candidatelib.SnapshotAPIVersion,
		System:     "cli-fixture",
		MutableAssets: []candidatelib.AssetRef{{
			Name:      "prompt",
			MediaType: "text/markdown",
			Path:      "candidate/prompt.md",
			SHA256:    sha256Text(childBytes),
			SizeBytes: int64(len(childBytes)),
		}},
		Dimensions: map[string]string{"model": "fixture-model"},
	}
	var err error
	parent.SnapshotID, err = candidatelib.DigestSnapshot(parent)
	mustNoError(t, err)
	child.SnapshotID, err = candidatelib.DigestSnapshot(child)
	mustNoError(t, err)
	writeYAML(t, root, "parent/snapshot.yaml", parent)
	writeYAML(t, root, "candidate/snapshot.yaml", child)
	manifest := candidatelib.CandidateManifest{
		APIVersion:        candidatelib.CandidateAPIVersion,
		CandidateID:       "prompt-001",
		ParentSnapshot:    "parent/snapshot.yaml",
		CandidateSnapshot: "candidate/snapshot.yaml",
		Proposer:          candidatelib.Proposer{Kind: "human", Identity: "test"},
		Mutation: candidatelib.MutationDeclaration{
			Asset:      "prompt",
			Hypothesis: "The candidate prompt is clearer.",
			ExpectedImprovement: candidatelib.ExpectedImprovement{
				Metric: "quality",
			},
			RegressionRisks: []string{"verbosity"},
		},
	}
	writeYAML(t, root, "candidate.yaml", manifest)
	return root
}

func sha256Text(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func writeYAML(t *testing.T, root, relative string, value any) {
	t.Helper()
	data, err := yaml.Marshal(value)
	mustNoError(t, err)
	writeFile(t, root, relative, data)
}

func writeFile(t *testing.T, root, relative string, data []byte) {
	t.Helper()
	path := filepath.Join(root, relative)
	mustNoError(t, os.MkdirAll(filepath.Dir(path), 0o700))
	mustNoError(t, os.WriteFile(path, data, 0o600))
}

func mustNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
