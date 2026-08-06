---
Title: Snapshot and Candidate V1 Bundle Contract
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
    - Path: repo://cmd/ragopt/commands/candidate/validate.go
      Note: Structured CLI projection
    - Path: repo://pkg/candidate/candidate.go
      Note: Candidate declaration and one-mutation validation
    - Path: repo://pkg/candidate/candidate_test.go
      Note: Executable bundle rejection matrix
    - Path: repo://pkg/candidate/digest.go
      Note: Canonical snapshot and candidate identity algorithms
    - Path: repo://pkg/candidate/path.go
      Note: Resolved bundle path confinement
    - Path: repo://pkg/candidate/snapshot.go
      Note: Snapshot identity and exact asset verification
    - Path: repo://pkg/candidate/types.go
      Note: Public v1 snapshot, candidate, asset, and mutation schemas
    - Path: repo://pkg/candidate/yaml.go
      Note: Strict single-document YAML decoder
ExternalSources: []
Summary: Implemented strict snapshot identity, exactly-one-mutation candidate validation, safe bundle paths, and the candidate validate CLI.
LastUpdated: 2026-08-06T11:19:05.014209565-04:00
WhatFor: Explain how to author, validate, review, and consume a ragopt-snapshot/v1 and ragopt-candidate/v1 bundle without weakening its identity or one-mutation guarantees.
WhenToUse: Read before creating candidate fixtures, integrating a product snapshot, invoking candidate validate, or passing a validated candidate into the paired runner.
---


# Snapshot and Candidate V1 Bundle Contract

## Goal

`pkg/candidate` turns a directory of human-authored YAML and asset files into a
validated, content-identified proposal. Its most important guarantee is not
that a manifest says one component changed. The package independently reads
both snapshots and all asset bytes and proves that exactly one mutable asset
changed while locked assets, semantic dimensions, system name, and the mutable
asset set stayed fixed.

The package is intentionally a validator, not an editor. It does not apply
patches, create candidate text, rewrite digests, repair malformed bundles, or
promote anything. `LoadCandidate` returns a candidate only after every common
invariant passes.

## Context

The candidate layer sits between proposal and evaluation:

```text
human/product proposer
        |
        v
candidate bundle on disk
        |
        | strict load: schemas + paths + bytes + identities + diff
        v
validated Candidate value
        |
        | Phase 3 binds exact identities to immutable run inputs
        v
paired incumbent/candidate evaluation
```

This boundary exists because an imperative `--try-reranker` branch cannot
prove what else changed. A candidate is useful only when a reviewer can answer
all of these questions from artifacts:

- What complete parent system was the incumbent?
- Which assets and scalar dimensions were locked?
- Which one mutable asset changed?
- Do declared SHA-256 values match the actual files?
- Was the mutation hypothesis written before evaluation?
- Which diagnostics motivated it and which regressions were anticipated?

The v1 package uses YAML for portable manifests, SHA-256 for content identity,
and local directory bundles. There is no content-addressed store, registry,
database, patch language, or backward-compatibility decoder.

## Quick Reference

### Public API

```go
func candidate.LoadSnapshot(
    ctx context.Context,
    bundleRoot string,
    manifestPath string,
) (*candidate.Snapshot, error)

func candidate.DigestSnapshot(snapshot candidate.Snapshot) (string, error)

func candidate.LoadCandidate(
    ctx context.Context,
    bundleRoot string,
    manifestPath string,
) (*candidate.Candidate, error)
```

`LoadSnapshot` and `LoadCandidate` resolve all file paths relative to the bundle
root. `manifestPath` is also bundle-relative. `DigestSnapshot` computes the
value that must appear in `snapshot_id`; it does not read or verify files, so
authoring code must still call `LoadSnapshot` or `LoadCandidate` afterward.

There is intentionally no separate public `ValidateCandidate` that accepts an
arbitrarily constructed, partially loaded object. The one public load path
means callers cannot accidentally treat a decoded but unverified manifest as
evaluation-ready.

### Bundle layout

A pragmatic bundle reuses parent paths for unchanged assets and contains full
replacement bytes for the changed asset:

```text
candidate-prompt-001/
├── candidate.yaml
├── shared/
│   ├── evaluation-suite.yaml
│   └── judge-prompt.md
├── parent/
│   ├── snapshot.yaml
│   └── assets/
│       ├── orchestration-prompt.md
│       └── search-description.md
├── candidate/
│   ├── snapshot.yaml
│   └── assets/
│       └── orchestration-prompt.md
└── provenance/
    └── selected-cases.json
```

The candidate snapshot may point at `shared/*` and unchanged `parent/assets/*`
paths. The changed asset points at a complete replacement file under
`candidate/assets/*`. It is not a patch and no production file is overwritten.

### Asset reference

```yaml
- name: orchestration_prompt
  media_type: text/markdown
  path: parent/assets/orchestration-prompt.md
  sha256: sha256:0123456789abcdef...64-lowercase-hex-total...
  size_bytes: 2810
```

Rules enforced during load:

- `name` matches `[a-z0-9][a-z0-9._-]{0,127}`;
- logical names are unique across locked and mutable assets;
- `media_type` is nonempty and no longer than 255 bytes;
- `path` is nonempty, relative, lexically canonical, and confined;
- the resolved target is a regular file;
- a symlink may resolve only to a regular file still inside the bundle;
- `sha256` is exactly `sha256:` plus 64 lowercase hexadecimal digits;
- `size_bytes` is nonnegative and equals the exact file length;
- the digest is recomputed from the exact file bytes and must match.

The path is part of snapshot identity. Repointing an unchanged asset to another
path, even one with identical bytes, is metadata drift and candidate validation
rejects it. This keeps snapshot manifests literal and reviewable.

### Snapshot schema

```yaml
api_version: ragopt-snapshot/v1
system: rag-ttc-tool-qa
snapshot_id: sha256:<semantic-snapshot-digest>

locked_assets:
  - name: evaluation_suite
    media_type: application/yaml
    path: shared/evaluation-suite.yaml
    sha256: sha256:<byte-digest>
    size_bytes: 18432
  - name: judge_prompt
    media_type: text/markdown
    path: shared/judge-prompt.md
    sha256: sha256:<byte-digest>
    size_bytes: 4380

mutable_assets:
  - name: orchestration_prompt
    media_type: text/markdown
    path: parent/assets/orchestration-prompt.md
    sha256: sha256:<byte-digest>
    size_bytes: 2810
  - name: search_description
    media_type: text/markdown
    path: parent/assets/search-description.md
    sha256: sha256:<byte-digest>
    size_bytes: 1420

dimensions:
  answer_model: provider/model-version
  judge_model: provider/judge-version
  corpus_digest: sha256:<digest>
  index_digest: sha256:<digest>
  evaluator: rag-ttc-answer-quality/v3
  tool_safety: rag-ttc-tool-safety/v1
```

The system name follows the logical-name pattern. At least one mutable asset is
required because a snapshot with nothing eligible for isolated mutation cannot
participate in this candidate contract. Locked assets may be empty for small
fixtures, although real integrations should enumerate every suite, judge,
corpus, index, safety policy, and evaluator input that affects outcomes.

Dimension keys match `[A-Za-z][A-Za-z0-9_.-]{0,127}`. Values must be nonempty,
must not have surrounding whitespace, and are bounded to 1024 bytes. A product
integration is responsible for requiring its domain-specific keys.

### Snapshot identity

The snapshot ID is:

```text
sha256(JSON({
  api_version,
  system,
  locked_assets: sort_by_name(locked_assets),
  mutable_assets: sort_by_name(mutable_assets),
  dimensions
}))
```

`snapshot_id` itself is excluded from the payload. JSON object keys and the
dimension map are deterministically ordered by the Go JSON encoder. Asset
declaration order is explicitly sorted away. Tests reverse asset declarations
and prove the identity remains stable.

Paths, media types, declared byte digests, and sizes remain in the identity.
The loader then checks those declarations against actual bytes, so the digest
cannot be made truthful merely by editing YAML.

### Candidate schema

```yaml
api_version: ragopt-candidate/v1
candidate_id: prompt-clarity-001
parent_snapshot: parent/snapshot.yaml
candidate_snapshot: candidate/snapshot.yaml

proposer:
  kind: human
  identity: intern-name

mutation:
  asset: orchestration_prompt
  hypothesis: >-
    Comparison failures occur because each named subject is not explicitly
    required to receive separately acquired evidence.
  expected_improvement:
    metric: comparison_completeness
    groups:
      - comparison
  regression_risks:
    - additional searches on simple questions
    - increased answer latency

evidence:
  diagnostic_manifest_digest: sha256:<digest>
  selected_case_ids:
    - cmp-03
    - cmp-08
```

Candidate IDs follow the same bounded logical-name pattern. Proposer kind and
identity are required bounded strings; v1 does not behave differently for a
human or model label. The hypothesis and target metric are required. At least
one unique, nonblank regression risk is required. Groups and selected case IDs
are optional but, when present, must be nonblank and unique. A diagnostic
manifest digest is optional but must use the canonical SHA-256 form when set.

### Candidate validation algorithm

The implementation follows this order:

```text
resolve bundle root
resolve + strict-decode candidate.yaml
validate candidate declaration

load parent snapshot:
  strict-decode one YAML document
  validate schema and scalar fields
  resolve every asset inside bundle
  read bytes; verify size and SHA-256
  recompute snapshot identity

load child snapshot using the same path

require parent.system == child.system
require parent.dimensions == child.dimensions
require parent.locked_assets == child.locked_assets, including bytes
require names(parent.mutable_assets) == names(child.mutable_assets)

for each mutable asset name:
  if bytes equal:
    require full AssetRef equal
  else:
    require media_type unchanged
    append independently detected mutation

require detected mutation count == 1
require declared mutation.asset == detected asset name
require parent snapshot ID != child snapshot ID
compute candidate semantic digest
return validated Candidate
```

Missing, malformed, duplicate, or extra state is an error. The loader uses
`yaml.Decoder.KnownFields(true)` and requires EOF after the first document, so
unknown fields and multi-document YAML cannot hide configuration.

### Candidate digest

The returned `Candidate.Digest` is SHA-256 over canonical JSON containing:

- candidate API version and human-readable candidate ID;
- verified parent and child snapshot IDs;
- proposer;
- mutation declaration and hypothesis;
- expected improvement and regression risks;
- evidence links.

The filesystem root and manifest locator paths are not part of this digest;
the verified snapshot identities replace them. Moving the same complete bundle
does not change candidate identity. Changing provenance or hypothesis does.

### Mutation result

The loader returns the independently observed mutation:

```go
type Mutation struct {
    AssetName    string
    Parent       AssetRef
    Candidate    AssetRef
    ParentDigest string
    ChildDigest  string
}
```

Evaluation code should use this result for reports. It must not trust only
`manifest.mutation.asset`.

### Immutability boundary

The validator never writes the bundle. A returned candidate is a validated
in-memory view of the bytes read during that call. If the directory changes,
the next load will recompute every identity and either produce a different
candidate or reject drift.

This is not yet evaluated-run custody. Phase 3 must copy or digest-link the
candidate manifest, both snapshots, and all referenced inputs into the
immutable run before executing an arm. That `validated -> running` binding is
the status transition after which the run's candidate identity cannot be
changed by editing the original bundle. There is deliberately no writable
`candidate validate` status file and no claim that filesystem permissions make
an arbitrary source directory immutable.

### CLI contract

```text
ragopt candidate validate \
  --bundle ./candidate-prompt-001 \
  --manifest candidate.yaml \
  --format table|json|jsonl|csv|tsv|yaml
```

The default manifest is `candidate.yaml`. The command emits one row only after
successful validation:

| Field | Meaning |
|---|---|
| `candidate_id` | Human-readable stable ID |
| `candidate_digest` | Full semantic candidate digest |
| `parent_snapshot` | Verified parent snapshot ID |
| `child_snapshot` | Verified child snapshot ID |
| `changed_asset` | Independently detected asset name |
| `parent_asset_digest` | Parent asset byte digest |
| `child_asset_digest` | Child asset byte digest |
| `bundle` | Resolved absolute bundle root |
| `valid` | Always `true` on an emitted row |

Validation failures return an error and emit no misleading `valid=false` row.
The command uses the three Glazed v1.4 serialization flags only:
`--format`, `--output-fields`, and `--max-output-rows`. `--log-level` and the
other standard logging flags are root flags. No environment values are read.

## Usage Examples

### Author a snapshot ID in Go

```go
snapshot := candidate.Snapshot{
    APIVersion: candidate.SnapshotAPIVersion,
    System:     "product-rag",
    MutableAssets: []candidate.AssetRef{
        promptRef,
    },
    LockedAssets: []candidate.AssetRef{
        suiteRef,
        judgeRef,
    },
    Dimensions: map[string]string{
        "model":     "provider/model-version",
        "evaluator": "product-answer-quality/v1",
    },
}

snapshot.SnapshotID, err = candidate.DigestSnapshot(snapshot)
if err != nil {
    return errors.Wrap(err, "compute snapshot identity")
}

// Marshal snapshot, then prove paths and byte declarations by loading it.
loaded, err := candidate.LoadSnapshot(ctx, bundleRoot, "parent/snapshot.yaml")
if err != nil {
    return errors.Wrap(err, "validate authored snapshot")
}
_ = loaded
```

The authoring process must calculate every `AssetRef.SHA256` from bytes. Do not
copy a digest from an earlier file based on filename or intent.

### Validate from a product command

```go
validated, err := candidate.LoadCandidate(ctx, bundleRoot, "candidate.yaml")
if err != nil {
    return errors.Wrap(err, "load evaluation candidate")
}

log.Info().
    Str("candidate_id", validated.Manifest.CandidateID).
    Str("candidate_digest", validated.Digest).
    Str("changed_asset", validated.Mutation.AssetName).
    Msg("candidate validated")
```

Do not mutate `validated.Parent`, `validated.Child`, or their maps and slices.
They are value projections for execution and reporting, not builder objects.

### Expected rejections

| Change | Result |
|---|---|
| prompt bytes differ | valid if it is the only mutable change |
| prompt patch file without complete replacement | missing/invalid asset |
| two mutable byte changes | `changed=2` rejection |
| locked judge or suite change | locked asset rejection |
| model/evaluator dimension change | locked dimension rejection |
| changed declaration names wrong asset | declaration/detected mismatch |
| unknown YAML field | strict schema rejection |
| second YAML document | multiple-document rejection |
| file content edited without digest update | size or digest mismatch |
| symlink resolves outside bundle | confinement rejection |
| unchanged bytes moved to another path | metadata/path drift rejection |

### Validation commands

```bash
go test ./pkg/candidate -count=1
go test ./cmd/ragopt/... -count=1
go test ./pkg/candidate ./cmd/ragopt/... -race -count=1
go run ./cmd/ragopt candidate validate --help
go test ./...
go build ./...
```

## Related

- Primary system design:
  `design-doc/01-ragopt-intern-guide-to-a-reusable-evidence-gated-optimization-harness.md`
- Run custody contract:
  `reference/02-runstore-v1-on-disk-contract-and-recovery-guarantees.md`
- Implementation diary: `reference/01-implementation-diary.md`
- Canonical ledger: `tasks.md`
- Code checkpoint: commit `d2329dd`
