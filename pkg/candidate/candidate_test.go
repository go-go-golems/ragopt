package candidate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoadCandidateValidExactlyOneMutation(t *testing.T) {
	fixture := newFixture(t)
	loaded, err := LoadCandidate(t.Context(), fixture.root, "candidate.yaml")
	mustNoError(t, err)
	if loaded.Manifest.CandidateID != "prompt-clarity-001" {
		t.Fatalf("candidate ID: got %q", loaded.Manifest.CandidateID)
	}
	if loaded.Mutation.AssetName != "prompt" {
		t.Fatalf("mutation asset: got %q", loaded.Mutation.AssetName)
	}
	if loaded.Mutation.ParentDigest == loaded.Mutation.ChildDigest {
		t.Fatal("changed asset digests unexpectedly match")
	}
	if loaded.Digest == "" || loaded.Parent.SnapshotID == loaded.Child.SnapshotID {
		t.Fatalf("invalid identities: %#v", loaded)
	}
	if !filepath.IsAbs(loaded.Root) {
		t.Fatalf("bundle root is not absolute: %q", loaded.Root)
	}
}

func TestCandidateIdentityIgnoresAssetDeclarationOrder(t *testing.T) {
	fixture := newFixture(t)
	first, err := LoadCandidate(t.Context(), fixture.root, "candidate.yaml")
	mustNoError(t, err)

	parent := fixture.readSnapshot(t, fixture.manifest.ParentSnapshot)
	child := fixture.readSnapshot(t, fixture.manifest.CandidateSnapshot)
	reverseAssets(parent.LockedAssets)
	reverseAssets(parent.MutableAssets)
	reverseAssets(child.LockedAssets)
	reverseAssets(child.MutableAssets)
	fixture.writeSnapshot(t, fixture.manifest.ParentSnapshot, parent)
	fixture.writeSnapshot(t, fixture.manifest.CandidateSnapshot, child)

	second, err := LoadCandidate(t.Context(), fixture.root, "candidate.yaml")
	mustNoError(t, err)
	if first.Parent.SnapshotID != second.Parent.SnapshotID || first.Child.SnapshotID != second.Child.SnapshotID || first.Digest != second.Digest {
		t.Fatal("semantic identity changed with asset declaration order")
	}
}

func TestLoadCandidateRejectsMultipleMutableChanges(t *testing.T) {
	fixture := newFixture(t)
	child := fixture.readSnapshot(t, fixture.manifest.CandidateSnapshot)
	data := []byte("description changed too\n")
	path := "candidate/assets/description.md"
	fixture.writeFile(t, path, data)
	replaceAsset(child.MutableAssets, "description", assetRef("description", "text/markdown", path, data))
	fixture.refreshSnapshotID(t, &child)
	fixture.writeSnapshot(t, fixture.manifest.CandidateSnapshot, child)

	_, err := LoadCandidate(t.Context(), fixture.root, "candidate.yaml")
	mustErrorContains(t, err, "exactly one mutable asset, changed=2")
}

func TestLoadCandidateRejectsLockedMutation(t *testing.T) {
	fixture := newFixture(t)
	child := fixture.readSnapshot(t, fixture.manifest.CandidateSnapshot)
	data := []byte("changed judge\n")
	path := "candidate/assets/judge.md"
	fixture.writeFile(t, path, data)
	replaceAsset(child.LockedAssets, "judge_prompt", assetRef("judge_prompt", "text/markdown", path, data))
	fixture.refreshSnapshotID(t, &child)
	fixture.writeSnapshot(t, fixture.manifest.CandidateSnapshot, child)

	_, err := LoadCandidate(t.Context(), fixture.root, "candidate.yaml")
	mustErrorContains(t, err, "changes locked asset \"judge_prompt\"")
}

func TestLoadCandidateRejectsMissingAndDriftedAssets(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		fixture := newFixture(t)
		mustNoError(t, os.Remove(filepath.Join(fixture.root, "candidate/assets/prompt.md")))
		_, err := LoadCandidate(t.Context(), fixture.root, "candidate.yaml")
		mustErrorContains(t, err, "resolve asset \"prompt\"")
	})
	t.Run("digest mismatch", func(t *testing.T) {
		fixture := newFixture(t)
		fixture.writeFile(t, "candidate/assets/prompt.md", []byte("drifted with same manifest"))
		_, err := LoadCandidate(t.Context(), fixture.root, "candidate.yaml")
		if err == nil || (!strings.Contains(err.Error(), "size mismatch") && !strings.Contains(err.Error(), "digest mismatch")) {
			t.Fatalf("expected asset integrity error, got %v", err)
		}
	})
}

func TestLoadCandidateRejectsPathEscapeThroughSymlink(t *testing.T) {
	fixture := newFixture(t)
	outside := t.TempDir()
	mustNoError(t, os.WriteFile(filepath.Join(outside, "secret.md"), []byte("secret"), 0o600))
	link := filepath.Join(fixture.root, "escape")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	child := fixture.readSnapshot(t, fixture.manifest.CandidateSnapshot)
	data := []byte("secret")
	replaceAsset(child.MutableAssets, "prompt", assetRef("prompt", "text/markdown", "escape/secret.md", data))
	fixture.refreshSnapshotID(t, &child)
	fixture.writeSnapshot(t, fixture.manifest.CandidateSnapshot, child)

	_, err := LoadCandidate(t.Context(), fixture.root, "candidate.yaml")
	mustErrorContains(t, err, "resolves outside bundle root")
}

func TestLoadCandidateRejectsStrictYAMLFailures(t *testing.T) {
	t.Run("unknown field", func(t *testing.T) {
		fixture := newFixture(t)
		path := filepath.Join(fixture.root, "candidate.yaml")
		file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o600)
		mustNoError(t, err)
		_, err = file.WriteString("unknown_field: true\n")
		mustNoError(t, err)
		mustNoError(t, file.Close())
		_, err = LoadCandidate(t.Context(), fixture.root, "candidate.yaml")
		mustErrorContains(t, err, "field unknown_field not found")
	})
	t.Run("multiple documents", func(t *testing.T) {
		fixture := newFixture(t)
		path := filepath.Join(fixture.root, "candidate.yaml")
		file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o600)
		mustNoError(t, err)
		_, err = file.WriteString("---\napi_version: second\n")
		mustNoError(t, err)
		mustNoError(t, file.Close())
		_, err = LoadCandidate(t.Context(), fixture.root, "candidate.yaml")
		mustErrorContains(t, err, "multiple documents")
	})
}

func TestLoadCandidateRejectsDeclaredAndActualMutationMismatch(t *testing.T) {
	fixture := newFixture(t)
	fixture.manifest.Mutation.Asset = "description"
	fixture.writeYAML(t, "candidate.yaml", fixture.manifest)
	_, err := LoadCandidate(t.Context(), fixture.root, "candidate.yaml")
	mustErrorContains(t, err, "does not match changed asset \"prompt\"")
}

func TestLoadCandidateRejectsLockedDimensionChange(t *testing.T) {
	fixture := newFixture(t)
	child := fixture.readSnapshot(t, fixture.manifest.CandidateSnapshot)
	child.Dimensions["model"] = "model-b"
	fixture.refreshSnapshotID(t, &child)
	fixture.writeSnapshot(t, fixture.manifest.CandidateSnapshot, child)
	_, err := LoadCandidate(t.Context(), fixture.root, "candidate.yaml")
	mustErrorContains(t, err, "changes locked snapshot dimensions")
}

func TestLoadSnapshotRejectsSnapshotIDMismatch(t *testing.T) {
	fixture := newFixture(t)
	parent := fixture.readSnapshot(t, fixture.manifest.ParentSnapshot)
	parent.SnapshotID = "sha256:" + strings.Repeat("0", 64)
	fixture.writeSnapshot(t, fixture.manifest.ParentSnapshot, parent)
	_, err := LoadSnapshot(t.Context(), fixture.root, fixture.manifest.ParentSnapshot)
	mustErrorContains(t, err, "snapshot ID mismatch")
}

type candidateFixture struct {
	root     string
	manifest CandidateManifest
}

func newFixture(t *testing.T) *candidateFixture {
	t.Helper()
	fixture := &candidateFixture{root: t.TempDir()}
	judge := []byte("judge prompt v1\n")
	promptParent := []byte("answer with evidence\n")
	promptChild := []byte("answer every named subject with evidence\n")
	description := []byte("search the corpus\n")
	fixture.writeFile(t, "assets/judge.md", judge)
	fixture.writeFile(t, "parent/assets/prompt.md", promptParent)
	fixture.writeFile(t, "candidate/assets/prompt.md", promptChild)
	fixture.writeFile(t, "parent/assets/description.md", description)

	locked := assetRef("judge_prompt", "text/markdown", "assets/judge.md", judge)
	descriptionRef := assetRef("description", "text/markdown", "parent/assets/description.md", description)
	parent := Snapshot{
		APIVersion:   SnapshotAPIVersion,
		System:       "fixture-rag",
		LockedAssets: []AssetRef{locked},
		MutableAssets: []AssetRef{
			assetRef("prompt", "text/markdown", "parent/assets/prompt.md", promptParent),
			descriptionRef,
		},
		Dimensions: map[string]string{"model": "model-a", "suite": "suite-v1"},
	}
	child := Snapshot{
		APIVersion:   SnapshotAPIVersion,
		System:       "fixture-rag",
		LockedAssets: []AssetRef{locked},
		MutableAssets: []AssetRef{
			assetRef("prompt", "text/markdown", "candidate/assets/prompt.md", promptChild),
			descriptionRef,
		},
		Dimensions: map[string]string{"suite": "suite-v1", "model": "model-a"},
	}
	fixture.refreshSnapshotID(t, &parent)
	fixture.refreshSnapshotID(t, &child)
	fixture.writeSnapshot(t, "parent/snapshot.yaml", parent)
	fixture.writeSnapshot(t, "candidate/snapshot.yaml", child)
	fixture.manifest = CandidateManifest{
		APIVersion:        CandidateAPIVersion,
		CandidateID:       "prompt-clarity-001",
		ParentSnapshot:    "parent/snapshot.yaml",
		CandidateSnapshot: "candidate/snapshot.yaml",
		Proposer:          Proposer{Kind: "human", Identity: "intern"},
		Mutation: MutationDeclaration{
			Asset:      "prompt",
			Hypothesis: "The explicit named-subject instruction improves comparison completeness.",
			ExpectedImprovement: ExpectedImprovement{
				Metric: "comparison_completeness",
				Groups: []string{"comparison"},
			},
			RegressionRisks: []string{"more tokens on simple questions"},
		},
		Evidence: Evidence{SelectedCaseIDs: []string{"cmp-03"}},
	}
	fixture.writeYAML(t, "candidate.yaml", fixture.manifest)
	return fixture
}

func assetRef(name, mediaType, path string, data []byte) AssetRef {
	return AssetRef{Name: name, MediaType: mediaType, Path: path, SHA256: digestBytes(data), SizeBytes: int64(len(data))}
}

func (fixture *candidateFixture) refreshSnapshotID(t *testing.T, snapshot *Snapshot) {
	t.Helper()
	digest, err := DigestSnapshot(*snapshot)
	mustNoError(t, err)
	snapshot.SnapshotID = digest
}

func (fixture *candidateFixture) readSnapshot(t *testing.T, relative string) Snapshot {
	t.Helper()
	var snapshot Snapshot
	data, err := os.ReadFile(filepath.Join(fixture.root, relative))
	mustNoError(t, err)
	mustNoError(t, yaml.Unmarshal(data, &snapshot))
	return snapshot
}

func (fixture *candidateFixture) writeSnapshot(t *testing.T, relative string, snapshot Snapshot) {
	t.Helper()
	fixture.writeYAML(t, relative, snapshot)
}

func (fixture *candidateFixture) writeYAML(t *testing.T, relative string, value any) {
	t.Helper()
	data, err := yaml.Marshal(value)
	mustNoError(t, err)
	fixture.writeFile(t, relative, data)
}

func (fixture *candidateFixture) writeFile(t *testing.T, relative string, data []byte) {
	t.Helper()
	path := filepath.Join(fixture.root, relative)
	mustNoError(t, os.MkdirAll(filepath.Dir(path), 0o700))
	mustNoError(t, os.WriteFile(path, data, 0o600))
}

func replaceAsset(assets []AssetRef, name string, replacement AssetRef) {
	for index := range assets {
		if assets[index].Name == name {
			assets[index] = replacement
			return
		}
	}
}

func reverseAssets(assets []AssetRef) {
	for left, right := 0, len(assets)-1; left < right; left, right = left+1, right-1 {
		assets[left], assets[right] = assets[right], assets[left]
	}
}

func mustNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func mustErrorContains(t *testing.T, err error, expected string) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), expected) {
		t.Fatalf("expected error containing %q, got %v", expected, err)
	}
}
