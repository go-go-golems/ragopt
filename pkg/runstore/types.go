// Package runstore owns durable, self-describing experiment run directories.
//
// A run is append-only while active and immutable after it reaches a terminal
// state. Domain-specific evaluation data remains in caller-selected artifacts;
// runstore only owns their custody and the run's semantic identity.
package runstore

import "time"

const (
	// ManifestSchemaVersion identifies the on-disk manifest contract.
	ManifestSchemaVersion = "ragopt-run/v1"

	StateActive   = "active"
	StateComplete = "complete"
	StateFailed   = "failed"
)

// Options configures a new run directory.
type Options struct {
	Root        string
	Name        string
	Description string
	Dimensions  map[string]string
}

// InputRef records one immutable input copied into the run.
type InputRef struct {
	Role         string `json:"role"`
	OriginalPath string `json:"original_path"`
	CopiedPath   string `json:"copied_path"`
	SHA256       string `json:"sha256"`
	SizeBytes    int64  `json:"size_bytes"`
}

// Host captures the local execution environment.
type Host struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	CPUs     int    `json:"cpus"`
}

// Manifest describes one run and its immutable configuration identity.
type Manifest struct {
	SchemaVersion string            `json:"schema_version"`
	RunID         string            `json:"run_id"`
	Name          string            `json:"name"`
	Description   string            `json:"description,omitempty"`
	StartedAt     time.Time         `json:"started_at"`
	GoVersion     string            `json:"go_version"`
	ModulePath    string            `json:"module_path,omitempty"`
	ModuleVersion string            `json:"module_version,omitempty"`
	Host          Host              `json:"host"`
	ConfigDigest  string            `json:"config_digest"`
	Dimensions    map[string]string `json:"dimensions,omitempty"`
}

// Status is the current or terminal state of a run.
type Status struct {
	State      string     `json:"state"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	Error      string     `json:"error,omitempty"`
}

// Summary is the small common terminal summary. Domain-specific reports stay
// as ordinary files below results/.
type Summary struct {
	Message string         `json:"message,omitempty"`
	Metrics map[string]any `json:"metrics,omitempty"`
}
