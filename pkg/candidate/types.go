// Package candidate loads and validates immutable system snapshots and
// exactly-one-mutation candidate bundles.
package candidate

const (
	SnapshotAPIVersion  = "ragopt-snapshot/v1"
	CandidateAPIVersion = "ragopt-candidate/v1"
)

// AssetRef identifies exact asset bytes within a candidate bundle.
type AssetRef struct {
	Name      string `json:"name" yaml:"name"`
	MediaType string `json:"media_type" yaml:"media_type"`
	Path      string `json:"path" yaml:"path"`
	SHA256    string `json:"sha256" yaml:"sha256"`
	SizeBytes int64  `json:"size_bytes" yaml:"size_bytes"`
}

// Snapshot is the complete semantic identity of one system configuration.
// Asset file contents are verified while loading but are not retained in this
// public projection.
type Snapshot struct {
	APIVersion    string            `json:"api_version" yaml:"api_version"`
	System        string            `json:"system" yaml:"system"`
	SnapshotID    string            `json:"snapshot_id" yaml:"snapshot_id"`
	LockedAssets  []AssetRef        `json:"locked_assets" yaml:"locked_assets"`
	MutableAssets []AssetRef        `json:"mutable_assets" yaml:"mutable_assets"`
	Dimensions    map[string]string `json:"dimensions" yaml:"dimensions"`
	ByteDigest    string            `json:"byte_digest,omitempty" yaml:"-"`

	assetBytes map[string][]byte
}

// Proposer records who authored a candidate. V1 does not change behavior by
// proposer kind.
type Proposer struct {
	Kind     string `json:"kind" yaml:"kind"`
	Identity string `json:"identity" yaml:"identity"`
}

// ExpectedImprovement declares the candidate's one target metric and optional
// case groups.
type ExpectedImprovement struct {
	Metric string   `json:"metric" yaml:"metric"`
	Groups []string `json:"groups,omitempty" yaml:"groups,omitempty"`
}

// MutationDeclaration records the candidate author's hypothesis. Validation
// computes the actual mutation independently.
type MutationDeclaration struct {
	Asset               string              `json:"asset" yaml:"asset"`
	Hypothesis          string              `json:"hypothesis" yaml:"hypothesis"`
	ExpectedImprovement ExpectedImprovement `json:"expected_improvement" yaml:"expected_improvement"`
	RegressionRisks     []string            `json:"regression_risks" yaml:"regression_risks"`
}

// Evidence links a proposal to the diagnostics that motivated it.
type Evidence struct {
	DiagnosticManifestDigest string   `json:"diagnostic_manifest_digest,omitempty" yaml:"diagnostic_manifest_digest,omitempty"`
	SelectedCaseIDs          []string `json:"selected_case_ids,omitempty" yaml:"selected_case_ids,omitempty"`
}

// CandidateManifest is the strict portable manifest stored as candidate.yaml.
type CandidateManifest struct {
	APIVersion        string              `json:"api_version" yaml:"api_version"`
	CandidateID       string              `json:"candidate_id" yaml:"candidate_id"`
	ParentSnapshot    string              `json:"parent_snapshot" yaml:"parent_snapshot"`
	CandidateSnapshot string              `json:"candidate_snapshot" yaml:"candidate_snapshot"`
	Proposer          Proposer            `json:"proposer" yaml:"proposer"`
	Mutation          MutationDeclaration `json:"mutation" yaml:"mutation"`
	Evidence          Evidence            `json:"evidence,omitempty" yaml:"evidence,omitempty"`
}

// Mutation is the independently verified single changed asset.
type Mutation struct {
	AssetName    string   `json:"asset_name"`
	Parent       AssetRef `json:"parent"`
	Candidate    AssetRef `json:"candidate"`
	ParentDigest string   `json:"parent_digest"`
	ChildDigest  string   `json:"child_digest"`
}

// Candidate is a fully loaded and validated immutable candidate view.
type Candidate struct {
	Manifest           CandidateManifest `json:"manifest"`
	Parent             Snapshot          `json:"parent"`
	Child              Snapshot          `json:"child"`
	Mutation           Mutation          `json:"mutation"`
	Digest             string            `json:"digest"`
	Root               string            `json:"root"`
	ManifestPath       string            `json:"manifest_path"`
	ManifestByteDigest string            `json:"manifest_byte_digest"`
}
