// Package eval executes deterministic, resumable incumbent/candidate pairs and
// records one durable cell per case, repeat, and arm.
package eval

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-go-golems/ragopt/pkg/candidate"
)

const (
	SuiteAPIVersion = "ragopt-suite/v1"
	CellAPIVersion  = "ragopt-cell/v1"
)

// Suite is the ordered set of opaque product cases evaluated together.
type Suite struct {
	APIVersion string `json:"api_version"`
	Name       string `json:"name"`
	Cases      []Case `json:"cases"`
}

// Case contains a stable ID, product-defined groups, and an opaque JSON input.
type Case struct {
	ID     string          `json:"id"`
	Groups []string        `json:"groups,omitempty"`
	Input  json.RawMessage `json:"input"`
}

// SuiteDocument is a validated suite plus its semantic identity and source
// file, which the runner copies into the immutable run.
type SuiteDocument struct {
	Suite      Suite  `json:"suite"`
	Digest     string `json:"digest"`
	ByteDigest string `json:"byte_digest,omitempty"`
	SourcePath string `json:"source_path"`
}

// ArtifactRef identifies one verified native artifact inside a run.
type ArtifactRef struct {
	Path      string `json:"path"`
	SHA256    string `json:"sha256"`
	SizeBytes int64  `json:"size_bytes"`
}

// Failure records an attributable arm-level failure without dropping the cell.
type Failure struct {
	Class   string `json:"class"`
	Message string `json:"message"`
}

// Outcome is the small comparison projection shared across product arms.
type Outcome struct {
	Completed      bool               `json:"completed"`
	ContractValid  bool               `json:"contract_valid"`
	Abstained      bool               `json:"abstained"`
	Failure        *Failure           `json:"failure,omitempty"`
	Metrics        map[string]float64 `json:"metrics,omitempty"`
	ProviderCalls  int                `json:"provider_calls,omitempty"`
	ToolCalls      int                `json:"tool_calls,omitempty"`
	InputTokens    int                `json:"input_tokens,omitempty"`
	OutputTokens   int                `json:"output_tokens,omitempty"`
	Duration       time.Duration      `json:"duration"`
	NativeArtifact ArtifactRef        `json:"native_artifact"`
}

// Cell is one durable point in the paired evaluation matrix.
type Cell struct {
	APIVersion     string    `json:"api_version"`
	RunID          string    `json:"run_id"`
	CaseID         string    `json:"case_id"`
	RepeatIndex    int       `json:"repeat_index"`
	Arm            string    `json:"arm"`
	CandidateID    string    `json:"candidate_id"`
	SnapshotDigest string    `json:"snapshot_digest"`
	SuiteDigest    string    `json:"suite_digest"`
	PolicyDigest   string    `json:"policy_digest"`
	StartedAt      time.Time `json:"started_at"`
	FinishedAt     time.Time `json:"finished_at"`
	Outcome        Outcome   `json:"outcome"`
	PreviousDigest string    `json:"previous_digest,omitempty"`
	Digest         string    `json:"digest"`
}

// ResolvedAsset is an immutable copied run input exposed to an arm.
type ResolvedAsset struct {
	Ref    candidate.AssetRef `json:"ref"`
	Path   string             `json:"path"`
	Locked bool               `json:"locked"`
}

// CandidateView is the bounded system view passed to one arm execution.
type CandidateView struct {
	Role            string                   `json:"role"`
	CandidateID     string                   `json:"candidate_id"`
	CandidateDigest string                   `json:"candidate_digest"`
	SnapshotDigest  string                   `json:"snapshot_digest"`
	System          string                   `json:"system"`
	Assets          map[string]ResolvedAsset `json:"assets"`
	Dimensions      map[string]string        `json:"dimensions"`
}

// Request is one arm invocation. NativeDirectory is a unique directory inside
// RunDirectory; the arm returns a run-relative native artifact path.
type Request struct {
	RunDirectory    string
	NativeDirectory string
	Case            Case
	RepeatIndex     int
	Candidate       CandidateView
}

// Arm executes one product-owned system implementation.
type Arm interface {
	Name() string
	Run(ctx context.Context, request Request) (Outcome, error)
}

// RunRequest defines one exact paired evaluation.
type RunRequest struct {
	RunRoot     string
	Name        string
	Description string
	Suite       *SuiteDocument
	PolicyPath  string
	Candidate   *candidate.Candidate
	Incumbent   Arm
	Challenger  Arm
	Repeats     int
}

// RunResult reports the durable run location and cell counts.
type RunResult struct {
	RunDirectory string `json:"run_directory"`
	RunID        string `json:"run_id"`
	Expected     int    `json:"expected_cells"`
	Completed    int    `json:"completed_cells"`
	Failures     int    `json:"failed_cells"`
	Resumed      bool   `json:"resumed"`
}
