---
Title: Paired Evaluation V1 API, Cell, and Resume Contract
Ticket: RAGOPT-001
Status: active
Topics:
    - rag
    - evaluation
    - experiments
    - reproducibility
    - architecture
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://pkg/eval/bind.go
      Note: Immutable input binding and copied candidate views
    - Path: repo://pkg/eval/resume.go
      Note: Committed-cell validation and active-tail recovery
    - Path: repo://pkg/eval/runner.go
      Note: New-run preparation, schedule, cell execution, and terminal state
    - Path: repo://pkg/eval/runner_test.go
      Note: Executable paired-run and resume contract
    - Path: repo://pkg/eval/suite.go
      Note: Strict suite normalization and semantic digest
    - Path: repo://pkg/eval/suite_test.go
      Note: Suite schema and identity tests
    - Path: repo://pkg/eval/types.go
      Note: Public suite, arm, outcome, cell, view, and request API
    - Path: repo://pkg/runstore/run.go
      Note: Explicit single-writer active resume primitive
ExternalSources: []
Summary: Implemented explicit two-arm paired execution, immutable input binding, synced cells, native artifact verification, and exact active-run resume.
LastUpdated: 2026-08-06T11:39:16.041933846-04:00
WhatFor: Explain how product code runs or resumes a deterministic incumbent/candidate evaluation while preserving exact inputs, native evidence, failures, and cell identities.
WhenToUse: Read before implementing a product Arm, starting or resuming a run, inspecting cells.jsonl, or building comparison and gate logic.
---


# Paired Evaluation V1 API, Cell, and Resume Contract

## Goal

`pkg/eval` is the execution bridge between a validated candidate and durable
comparison evidence. It deliberately knows nothing about retrieval, chat
messages, model providers, judges, or product answer schemas. A product supplies
two in-process `Arm` implementations. The runner supplies deterministic work
coordinates, copied candidate assets, native artifact directories, result-cell
custody, and exact resume.

The Phase 3 guarantee is stronger than “the loop can restart.” The runner proves
that every resumed cell belongs to the same suite, policy bytes, candidate,
snapshot, arm, case, and repeat; that every native artifact still exists with
the recorded bytes; and that only missing work executes.

## Context

The public design chose one of the two earlier arm alternatives: **two explicit
in-process arms**, one incumbent and one challenger. This matches the exercised
RAG-TTC evaluator boundary and avoids shipping both an arm-per-system API and a
second evaluator-plus-view compatibility API.

```text
validated candidate bundle
       |
       | bind exact manifests, snapshots, policy, suite, and assets
       v
+---------------- immutable run ----------------+
| inputs/                                         |
| config.json + manifest.json + active status     |
|                                                  |
| case 0 / repeat 0: incumbent -> sync cell        |
| case 0 / repeat 0: challenger -> sync cell       |
| case 0 / repeat 1: incumbent -> sync cell        |
| ...                                              |
|                                                  |
| native/<arm>/<case>/<repeat>/...                 |
| results/cells.jsonl                              |
+--------------------------------------------------+
       | complete                     | interruption
       v                              v
 terminal summary              active, explicitly resumable
```

The standalone `ragopt` binary does not execute product arms. A consumer
repository imports `pkg/eval` in its own experiment command, where it can build
its real runtime safely. This avoids subprocess, RPC, Go plugin, and dynamic
adapter protocols that have not been proved necessary.

## Quick Reference

### Public execution API

```go
type Arm interface {
    Name() string
    Run(ctx context.Context, request eval.Request) (eval.Outcome, error)
}

func eval.LoadSuite(ctx context.Context, path string) (*eval.SuiteDocument, error)

func eval.Run(ctx context.Context, request eval.RunRequest) (*eval.RunResult, error)

func eval.Resume(
    ctx context.Context,
    runDirectory string,
    request eval.RunRequest,
) (*eval.RunResult, error)
```

Every concrete consumer arm should assert its contract:

```go
var _ eval.Arm = (*ToolQAArm)(nil)
```

### Run request

```go
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
```

Rules:

- `RunRoot` and `Name` are required; the root becomes absolute before run
  creation.
- `Repeats` is in `[1, 10000]`.
- the suite must come from strict `LoadSuite` and its source is reloaded;
- the candidate must come from strict `LoadCandidate` and is reloaded;
- the policy is an opaque required regular file in Phase 3;
- arm names are distinct bounded identifiers;
- arm names are captured once before execution, not queried dynamically per
  cell;
- v1 executes sequentially and has no worker count or scheduler flags.

Phase 4 replaces the opaque policy interpretation with a strict gate-policy
schema. Phase 3 already locks the policy's exact bytes and digest in every
cell, so comparison cannot accidentally use another policy.

### Suite schema and identity

```json
{
  "api_version": "ragopt-suite/v1",
  "name": "retrieval-validation-v3",
  "cases": [
    {
      "id": "cmp-03",
      "groups": ["comparison", "multi-subject"],
      "input": {
        "question": "Compare A and B",
        "product_case_ref": "questions/cmp-03"
      }
    }
  ]
}
```

The loader:

- rejects unknown fields and multiple top-level JSON values;
- requires `ragopt-suite/v1`, a bounded name, and at least one case;
- requires unique bounded case IDs;
- accepts any valid JSON value as opaque input, including `null`;
- requires unique bounded group names and sorts them canonically;
- canonicalizes each opaque input with `json.Number` preservation;
- preserves declared case order in the suite digest;
- stores the absolute source path so the raw file can be copied into the run.

The suite digest is SHA-256 over normalized suite JSON. Group and JSON object
ordering noise is removed; case order remains semantic because it determines
execution order.

### Immutable input binding

Before either arm runs, `eval.Run` copies all of these through `runstore`:

- raw evaluation suite;
- raw gate policy;
- candidate manifest;
- parent snapshot manifest;
- candidate snapshot manifest;
- every parent locked and mutable asset;
- every candidate locked and mutable asset.

The run config contains the complete expected role-to-byte-digest map. After
copying, and again on resume, the runner requires:

```text
count(run inputs) == count(expected inputs)
for every expected role:
  role exists exactly once
  copied SHA-256 == expected source SHA-256
```

This is the candidate `validated -> running` transition. An arm never receives
paths into the original source bundle. It receives absolute paths to copies
inside the run. Editing the original bundle after execution cannot change the
run reader's valid inputs or historical cells.

Asset roles include side, locked/mutable class, escaped logical name, and a
digest suffix. This prevents role collisions without adding a content store.
Parent and candidate copies are separate even when a locked file has the same
bytes, keeping ownership literal for v1.

### Candidate view passed to an arm

```go
type CandidateView struct {
    Role            string // incumbent | candidate
    CandidateID     string
    CandidateDigest string
    SnapshotDigest  string
    System          string
    Assets          map[string]ResolvedAsset
    Dimensions      map[string]string
}

type ResolvedAsset struct {
    Ref    candidate.AssetRef
    Path   string // absolute copied path inside run/inputs
    Locked bool
}
```

The runner clones the asset and dimension maps for each invocation. The view
contains only configuration and copied paths. It does not expose a writable
candidate store, runstore handle, promotion API, or safety-policy override.

### Per-cell arm request

```go
type Request struct {
    RunDirectory    string
    NativeDirectory string
    Case            Case
    RepeatIndex     int
    Candidate       CandidateView
}
```

`NativeDirectory` is unique for one arm/case/repeat:

```text
native/<captured-arm-name>/<case-id>/<zero-padded-repeat>/
```

The arm owns files below that directory. It returns one run-relative path to
its authoritative native artifact. The runner resolves symlinks, requires the
target to remain inside the assigned cell directory, requires a regular file,
and recomputes its digest and size.

### Outcome contract

```go
type Outcome struct {
    Completed      bool
    ContractValid  bool
    Abstained      bool
    Failure        *Failure
    Metrics        map[string]float64
    ProviderCalls  int
    ToolCalls      int
    InputTokens    int
    OutputTokens   int
    Duration       time.Duration
    NativeArtifact ArtifactRef
}
```

Validation rules:

- missing native artifact is a custody failure;
- artifact path is canonical, run-relative, and inside the assigned directory;
- provided artifact digest/size, when set, must match recomputed values;
- metric names are bounded identifiers and values must be finite;
- missing metric and metric value zero remain distinct;
- counts and duration cannot be negative;
- the runner replaces `Duration` with measured wall duration;
- incomplete outcome requires a failure;
- failure requires a bounded class and nonblank message;
- failed outcome cannot also be completed, contract-valid, or abstained;
- abstention must be a completed outcome.

`ContractValid=false` on a completed nonfailed outcome remains representable.
That allows a product to record an answer artifact that completed but violated
an answer or citation contract. Phase 4 accounts for it as a separate hard
dimension instead of dropping it.

### Ordinary arm errors

If `Arm.Run` returns an ordinary error, the runner does not abort the matrix.
It writes `arm-error.json` in the assigned native directory, computes its
identity, and appends a failed cell:

```json
{
  "completed": false,
  "contract_valid": false,
  "failure": {
    "class": "arm_error",
    "message": "scripted provider failure"
  },
  "native_artifact": {
    "path": "native/challenger/cmp-03/0000/arm-error.json",
    "sha256": "sha256:...",
    "size_bytes": 91
  }
}
```

This keeps failures in denominators and lets later cells run. Cancellation and
deadline errors are different: they leave the run active without fabricating a
failed cell, so explicit resume can execute the interrupted coordinate again.

### Cell schema

Each durable line in `results/cells.jsonl` is `ragopt-cell/v1`:

```json
{
  "api_version": "ragopt-cell/v1",
  "run_id": "20260806T...",
  "case_id": "cmp-03",
  "repeat_index": 0,
  "arm": "challenger",
  "candidate_id": "prompt-clarity-001",
  "snapshot_digest": "sha256:child...",
  "suite_digest": "sha256:suite...",
  "policy_digest": "sha256:policy-bytes...",
  "started_at": "2026-08-06T15:00:00Z",
  "finished_at": "2026-08-06T15:00:01Z",
  "outcome": {}
}
```

The full key is:

```text
(suite digest,
 policy digest,
 candidate ID,
 snapshot digest,
 case ID,
 repeat index,
 captured arm name)
```

Candidate ID is the proposal identity shared by both arms. Snapshot digest and
arm distinguish parent from child. The run ID is validated separately so a
cell copied from another run cannot be accepted even if its semantic key
matches.

### Deterministic schedule

V1 always uses:

```text
for case in suite.cases, declared order:
  for repeat = 0; repeat < repeats; repeat++:
    execute incumbent
    append + fsync incumbent cell
    execute challenger
    append + fsync challenger cell
```

Interleaving reduces temporal drift and produces natural strict pairs. There is
no randomized arm order or seed. There is no goroutine or worker pool. If live
measurements later show sequential execution is unacceptable, concurrency must
preserve deterministic cell identity and custody and receive its own design.

### Run config identity

`runstore` hashes a `ragopt-eval-run/v1` config containing:

- suite semantic digest;
- policy byte digest;
- candidate ID and semantic digest;
- parent and child snapshot digests;
- captured incumbent and challenger names;
- repeat count;
- complete expected input role/digest map.

Resume supplies the same config to `runstore.Resume`. A single mismatch refuses
write access before any result file changes. The run store also requires
`status=active`. Complete and failed runs cannot be reopened.

### Explicit resume algorithm

```text
prepare request again
  reload suite, candidate, policy, and expected input bytes
  recompute complete eval run config

runstore.Resume(runDir, expectedConfig)
  strict Open common run artifacts
  require active status
  require config digest equality

verify exact copied input role/digest set
rebuild incumbent and candidate views from copied inputs
read cells.jsonl

if final bytes do not end in newline:
  truncate to last committed newline
  fsync file and containing directory

for each nonempty committed line:
  strict-decode one cell
  require key exists in deterministic schedule
  reject any duplicate key
  validate run ID, coordinates, timestamps, semantic identities
  revalidate native artifact path, bytes, digest, and size

execute only schedule keys not already present
require completed count == expected count
write summary and complete status
```

The final truncated fragment is discarded because `AppendJSONL` writes a
newline before `fsync`; a no-newline fragment never crossed the successful
record boundary. A malformed middle line, blank middle line, duplicate cell,
or unexpected identity is corruption and marks the run failed. It is never
silently skipped.

### Resume authority and single-writer rule

`runstore.Resume` is explicit and requires the exact config. It does not use an
OS lock and cannot prove another process has stopped. The operator or product
command must ensure there is exactly one writer. This is an honest local v1
constraint, not an invitation to add a distributed lease system.

If identity mismatch occurs before a writer is acquired, the active run remains
active. If copied-input or result corruption is discovered after acquisition,
the runner marks the run failed while preserving existing artifacts.

### Terminal summary and result

On success the runner writes:

```json
{
  "message": "paired evaluation complete",
  "metrics": {
    "expected_cells": 16,
    "completed_cells": 16,
    "failed_cells": 1
  }
}
```

`RunResult` reports the directory, run ID, expected/completed/failed cells, and
whether the call was a resume. `failed_cells` counts recorded outcomes with a
`Failure`; completed but contract-invalid outcomes remain separately visible in
cells for Phase 4.

### Error classes and state effects

| Condition | Cell recorded? | Run state |
|---|---:|---|
| successful arm outcome | yes | active until all work completes |
| ordinary arm error | failed cell | continue |
| context cancellation/deadline | no for interrupted cell | active |
| candidate/suite/policy identity mismatch before create/resume | no | no run mutation |
| config mismatch in `runstore.Resume` | no | remains active |
| missing copied input after writer acquisition | no | failed |
| malformed/duplicate/unexpected cell | no | failed |
| native artifact escape/missing/digest drift | no | failed |
| NaN/infinite metric | no | failed |
| result append or terminal publication failure | no | best-effort failed status |

Custody errors stop because later evidence would not be trustworthy. Product
arm errors continue because they are experimental outcomes.

## Usage Examples

### Implement a product arm

```go
type ToolQAArm struct {
    runtimeFactory RuntimeFactory
}

var _ eval.Arm = (*ToolQAArm)(nil)

func (arm *ToolQAArm) Name() string { return "rag-ttc-tool-qa" }

func (arm *ToolQAArm) Run(ctx context.Context, request eval.Request) (eval.Outcome, error) {
    prompt, ok := request.Candidate.Assets["orchestration_prompt"]
    if !ok {
        return eval.Outcome{}, errors.New("orchestration_prompt asset missing")
    }

    runtime, err := arm.runtimeFactory.Build(ctx, RuntimeConfig{
        PromptPath: prompt.Path,
        Dimensions: request.Candidate.Dimensions,
    })
    if err != nil {
        return eval.Outcome{}, errors.Wrap(err, "build product runtime")
    }

    nativePath := filepath.Join(request.NativeDirectory, "answer.json")
    result, err := runtime.Answer(ctx, request.Case.Input, nativePath)
    if err != nil {
        return eval.Outcome{}, err // runner records an attributable failed cell
    }
    relative, err := filepath.Rel(request.RunDirectory, nativePath)
    if err != nil {
        return eval.Outcome{}, errors.Wrap(err, "make artifact path relative")
    }

    return eval.Outcome{
        Completed:      true,
        ContractValid:  result.ContractValid,
        Abstained:      result.Abstained,
        Metrics:        result.Metrics,
        ProviderCalls:  result.ProviderCalls,
        ToolCalls:      result.ToolCalls,
        InputTokens:    result.InputTokens,
        OutputTokens:   result.OutputTokens,
        NativeArtifact: eval.ArtifactRef{Path: relative},
    }, nil
}
```

The complete product result stays in `answer.json`. `Outcome` is only the
common comparison projection.

### Start a run

```go
suite, err := eval.LoadSuite(ctx, "fixtures/validation-suite.json")
if err != nil {
    return errors.Wrap(err, "load validation suite")
}
proposal, err := candidate.LoadCandidate(ctx, candidateRoot, "candidate.yaml")
if err != nil {
    return errors.Wrap(err, "load candidate")
}

request := eval.RunRequest{
    RunRoot:     "artifacts/evaluations",
    Name:        "prompt clarity validation",
    Description: proposal.Manifest.Mutation.Hypothesis,
    Suite:       suite,
    PolicyPath:  "policies/promotion-v1.yaml",
    Candidate:   proposal,
    Incumbent:   incumbentArm,
    Challenger:  challengerArm,
    Repeats:     2,
}

result, err := eval.Run(ctx, request)
if err != nil {
    // Preserve result.RunDirectory when result is nonnil for diagnosis/resume.
    return errors.Wrap(err, "run paired validation")
}
```

### Resume after cancellation

```go
result, err := eval.Resume(ctx, interruptedRunDirectory, request)
if err != nil {
    return errors.Wrap(err, "resume paired validation")
}
```

Use the same semantic request. `RunRoot`, name, and description do not identify
completed cells, but the suite, policy, candidate, snapshots, arm names,
repeats, and bound inputs must match exactly.

### Review and validation commands

```bash
go test ./pkg/eval ./pkg/runstore ./pkg/candidate -count=1
go test ./pkg/eval ./pkg/runstore ./pkg/candidate -race -count=1
go test ./...
go build ./...
```

Start review with these tests:

- `TestInterruptedRunResumesToUninterruptedCanonicalResult`;
- `TestArmErrorRecordsFailedCellAndContinues`;
- `TestResumeRejectsIdentityMismatchWithoutMutatingActiveRun`;
- `TestResumeRejectsMissingBoundInput`;
- `TestResumeRejectsMalformedMiddleAndDuplicateCells`;
- `TestInvalidOutcomeIsCustodyFailure`;
- `TestNonFiniteMetricIsCustodyFailure`.

## Related

- System design:
  `design-doc/01-ragopt-intern-guide-to-a-reusable-evidence-gated-optimization-harness.md`
- Run store contract:
  `reference/02-runstore-v1-on-disk-contract-and-recovery-guarantees.md`
- Snapshot/candidate contract:
  `reference/03-snapshot-and-candidate-v1-bundle-contract.md`
- Diary: `reference/01-implementation-diary.md`
- Ledger: `tasks.md`
- Code checkpoint: commit `0802670`
