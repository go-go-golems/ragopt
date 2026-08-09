---
Title: Runstore V1 On-Disk Contract and Recovery Guarantees
Ticket: RAGOPT-001
Status: active
Topics:
    - rag
    - evaluation
    - experiments
    - reproducibility
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://pkg/runstore/input.go
      Note: Copied immutable input custody
    - Path: repo://pkg/runstore/path.go
      Note: Artifact path and symlink confinement
    - Path: repo://pkg/runstore/read.go
      Note: Strict Open validation contract
    - Path: repo://pkg/runstore/run.go
      Note: Create and active artifact API implementation
    - Path: repo://pkg/runstore/run_test.go
      Note: Executable examples of recovery guarantees
    - Path: repo://pkg/runstore/terminal.go
      Note: Complete and failed state transitions
    - Path: repo://pkg/runstore/types.go
      Note: Public ragopt-run/v1 schema definitions
    - Path: repo://pkg/runstore/write.go
      Note: Canonical digest and atomic publication primitives
ExternalSources: []
Summary: Implemented ragopt-run/v1 schemas, artifact layout, durability boundaries, read validation, and interruption recovery contract.
LastUpdated: 2026-08-06T11:02:10.95879371-04:00
WhatFor: Explain the implemented Phase 1 storage API precisely enough to create, inspect, test, and review a run without inferring behavior from source alone.
WhenToUse: Read before writing an evaluator, resume implementation, artifact inspection command, or migration for ragopt-run/v1.
---


# Runstore V1 On-Disk Contract and Recovery Guarantees

## Goal

`pkg/runstore` is the durable custody layer for `ragopt`. It answers a narrow
question: after an experiment process exits, crashes, or is cancelled, which
configuration and copied inputs defined the run, which artifacts were
successfully published, and whether the run was active, complete, or failed?

It does not interpret product results. An evaluator chooses its JSONL record
schema and native artifacts. Runstore provides identity, containment,
publication, and terminal-state rules that those higher layers can rely on.

## Context

The implementation is a narrow port of the run lifecycle exercised in
RAG-TTC. The port deliberately changes the schema name and generic vocabulary,
adds a strict reader, strengthens existing-symlink rejection, and uses caller
supplied semantic dimensions instead of product tags. It does not add a
database, file locks, retry engine, or active-run mutation API.

The central lifecycle is:

```text
Create
  |
  v
+-------------------------------+
| active run                    |
| - copy immutable inputs       |
| - atomically write artifacts  |
| - append + fsync JSONL cells  |
+-------------------------------+
  | Complete(summary)       | Fail(error)
  v                         v
+----------------+      +----------------+
| complete       |      | failed         |
| summary exists |      | prior files stay|
+----------------+      +----------------+
  \___________________________/
               |
               v
       all mutation rejected
```

`Open` is deliberately read-only. It validates active and terminal runs but
does not turn an abandoned active directory back into a writer. Phase 3 must
design explicit resume ownership using full evaluation identities; it must not
infer that any arbitrary active directory is safe to overwrite.

## Quick Reference

### Package API

The implemented entry points are:

```go
func Create(ctx context.Context, options runstore.Options, config any) (*runstore.Run, error)
func Open(dir string) (*runstore.Reader, error)

func (r *runstore.Run) Dir() string
func (r *runstore.Run) Manifest() runstore.Manifest
func (r *runstore.Run) Path(relative string) (string, error)
func (r *runstore.Run) CopyInput(ctx context.Context, role, source string) (runstore.InputRef, error)
func (r *runstore.Run) WriteJSON(ctx context.Context, relative string, value any) error
func (r *runstore.Run) WriteBytes(ctx context.Context, relative string, data []byte) error
func (r *runstore.Run) AppendJSONL(ctx context.Context, relative string, value any) error
func (r *runstore.Run) Complete(ctx context.Context, summary runstore.Summary) error
func (r *runstore.Run) Fail(ctx context.Context, cause error) error

func (r *runstore.Reader) Dir() string
func (r *runstore.Reader) Manifest() runstore.Manifest
func (r *runstore.Reader) Status() runstore.Status
func (r *runstore.Reader) Inputs() []runstore.InputRef
func (r *runstore.Reader) Path(relative string) (string, error)
```

The source of truth is:

- `pkg/runstore/types.go` for public schemas;
- `pkg/runstore/run.go` for active-run writes and JSONL durability;
- `pkg/runstore/input.go` for copied inputs;
- `pkg/runstore/terminal.go` for terminal transitions;
- `pkg/runstore/read.go` for read validation;
- `pkg/runstore/path.go` and `write.go` for filesystem invariants.

### Directory contract

Every created run has this minimum shape:

```text
<root>/
└── <timestamp>-<safe-name>-<random>/
    ├── manifest.json
    ├── config.json
    ├── status.json
    ├── inputs/
    │   ├── manifest.json       # appears after the first CopyInput
    │   └── <safe-role>.<ext>   # exact copied bytes
    ├── results/
    │   ├── ...                 # caller-owned common result artifacts
    │   └── summary.json        # appears only on Complete
    └── native/
        └── ...                 # caller-owned native evidence
```

The directory name is the run ID. `Open` rejects a copied or renamed run whose
directory basename no longer equals `manifest.run_id`. This makes accidental
directory relabeling visible.

### Manifest schema

The schema identifier is exactly `ragopt-run/v1`:

```json
{
  "schema_version": "ragopt-run/v1",
  "run_id": "20260806T150000.000000000Z-paired-evaluation-a1b2c3d4e5f6",
  "name": "Paired Evaluation",
  "description": "optional operator context",
  "started_at": "2026-08-06T15:00:00Z",
  "go_version": "go1.26.5",
  "module_path": "github.com/example/product",
  "module_version": "(devel)",
  "host": {
    "hostname": "worker-1",
    "os": "linux",
    "arch": "amd64",
    "cpus": 16
  },
  "config_digest": "sha256:<64 lowercase hexadecimal characters>",
  "dimensions": {
    "suite": "retrieval-validation-v3",
    "model": "provider/model-version"
  }
}
```

Semantic dimensions are bounded string facts chosen by the caller. Keys must
be non-empty and must not have surrounding whitespace. Values are trimmed and
must remain non-empty. The map is cloned on input and output so caller mutation
cannot silently change the in-memory manifest.

The config digest is SHA-256 over canonical JSON. Object keys are sorted by Go
JSON encoding, whitespace is irrelevant, and numbers are decoded with
`json.Number` before re-encoding so large integer text is not routed through a
lossy float. `config.json` remains indented for people; its semantic digest is
independent of formatting.

### Copied input schema

`inputs/manifest.json` is an array of:

```json
{
  "role": "evaluation suite",
  "original_path": "/absolute/path/to/suite.json",
  "copied_path": "inputs/evaluation-suite.json",
  "sha256": "sha256:<digest of exact copied bytes>",
  "size_bytes": 1234
}
```

Roles and copied paths must be unique. The role is preserved for display, while
its safe lowercase form determines the destination filename. Extensions are
preserved. An input is copied by bytes; a later change to the original file
does not alter the run. `Open` recomputes the copied file's size and digest and
rejects untracked files under `inputs/`.

### Status state machine

Valid combinations are:

| State | `finished_at` | `error` | `results/summary.json` |
|---|---:|---:|---:|
| `active` | absent | empty | not required |
| `complete` | required | empty | required and valid |
| `failed` | required | optional | not required |

The status start time must equal the manifest start time, and a finish time
cannot precede it. Unknown JSON fields and multiple JSON values are rejected.
There is no compatibility fallback for unknown schema versions or states.

### Durability boundaries

`WriteJSON` and `WriteBytes` publish atomically:

```text
validate relative path
mkdir parents with 0700
create sibling temporary file
chmod 0600
write all bytes
fsync temporary file
close temporary file
rename over destination
fsync containing directory
```

Therefore a reader should see the previous complete artifact or the new
complete artifact, not a partially written JSON file.

`AppendJSONL` uses a different contract because append streams cannot be
published by replacement after every record:

```text
marshal one value to one compact JSON line
open destination with O_CREATE | O_APPEND | O_WRONLY
append line + newline
fsync file
close file
return success
```

Each successful return is the custody boundary. If the process is interrupted
after N successful calls, those N lines must remain parseable. The caller must
not batch unrecorded in-memory results and then claim resumability.

### Path and permissions rules

- Artifact paths must be non-empty and relative.
- `.` and any lexical `..` escape are rejected.
- Existing symlink components are rejected, even when they point back inside
  the run; this keeps the rule simple and auditable.
- Created directories use mode `0700`; created files use mode `0600`.
- The containment check assumes a normal single-owner workspace. It does not
  claim protection against another process racing symlink changes between
  validation and filesystem operations.
- The root itself may be relative. `Open` normalizes the reader's directory to
  an absolute path.

### Failure and interruption guarantees

`Fail` writes only terminal status. It does not delete or truncate prior
artifacts. A failed arm may therefore leave useful native evidence and synced
cells for diagnosis.

An interrupted process does not execute `Complete` or `Fail`, so the on-disk
state remains `active`. This is intentional: the library refuses to fabricate
a terminal interpretation. `Open` validates that active state and exposes the
synced artifacts through confined paths. A future resume layer must compare
the suite, snapshots, candidate, policy, repeat coordinates, and result cell
keys before acquiring write ownership.

### Reader validation checklist

`Open` rejects the run if any common invariant fails:

- missing or non-directory run path;
- unknown manifest field or schema;
- run ID/directory mismatch;
- missing name or start time;
- invalid semantic dimension;
- malformed config or config digest drift;
- malformed status, invalid state combination, or timestamp mismatch;
- missing summary for a complete run;
- malformed input manifest;
- duplicate input roles or copied paths;
- unsafe copied path or symlink component;
- missing, size-drifted, or digest-drifted copied input;
- untracked copied file in `inputs/`.

It does not validate caller-defined files in `results/` or `native/`. Their
schemas and digests belong to Phase 3's evaluation contract.

## Usage Examples

### Create, record, and complete a run

```go
ctx := context.Background()

run, err := runstore.Create(ctx, runstore.Options{
    Root:        "artifacts/runs",
    Name:        "retrieval candidate validation",
    Description: "candidate changes only the retrieval prompt",
    Dimensions: map[string]string{
        "suite": "retrieval-validation-v3",
        "model": "provider/model-version",
    },
}, map[string]any{
    "repeats": 2,
    "policy":  "promotion-policy-v1",
})
if err != nil {
    return errors.Wrap(err, "create evaluation run")
}

if _, err := run.CopyInput(ctx, "suite", "fixtures/suite.json"); err != nil {
    _ = run.Fail(context.Background(), err)
    return errors.Wrap(err, "copy evaluation suite")
}

for _, cell := range completedCells {
    if err := run.AppendJSONL(ctx, "results/cells.jsonl", cell); err != nil {
        _ = run.Fail(context.Background(), err)
        return errors.Wrap(err, "record result cell")
    }
}

if err := run.Complete(ctx, runstore.Summary{
    Message: "paired evaluation complete",
    Metrics: map[string]any{"cells": len(completedCells)},
}); err != nil {
    return errors.Wrap(err, "complete evaluation run")
}
```

The cleanup calls use a fresh context only to make a best effort to record a
terminal failure after the evaluation context is cancelled. Whether an
evaluator should mark cancellation failed or leave the run active for explicit
resume is a Phase 3 policy decision; runstore does not choose it.

### Inspect an interrupted run

```go
reader, err := runstore.Open(runDirectory)
if err != nil {
    return errors.Wrap(err, "validate run")
}

status := reader.Status()
if status.State == runstore.StateActive {
    cellsPath, err := reader.Path("results/cells.jsonl")
    if err != nil {
        return err
    }
    // Parse complete JSONL records and compare their full cell identities.
    _ = cellsPath
}
```

Do not construct `Run` directly and do not rename an interrupted directory to
pretend it is a new run. `Reader` deliberately has no mutation methods.

### Review and validation commands

```bash
go test ./pkg/runstore -count=1
go test ./pkg/runstore -race -count=1
go test ./...
go build ./...
```

The focused tests in `pkg/runstore/run_test.go` exercise:

- active reopen after two synced JSONL cells;
- completion and all terminal write rejection;
- duplicate input role and copied-path rejection;
- stable canonical config digest under map insertion-order changes;
- lexical path escape and symlink escape rejection;
- failed-run artifact preservation;
- config and copied-input drift detection;
- unknown fields and invalid semantic dimensions.

## Related

- Primary architecture:
  `design-doc/01-ragopt-intern-guide-to-a-reusable-evidence-gated-optimization-harness.md`
- Canonical phase ledger: `tasks.md`
- Implementation history: `reference/01-implementation-diary.md`
- Code checkpoint: commit `9f1ccb4` (`feat(runstore): add durable run lifecycle`)
