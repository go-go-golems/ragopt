---
Title: Paired Comparison, Gate, and Promotion Report v1 Contract
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
    - Path: repo://pkg/eval/artifacts.go
      Note: Strict read-only evaluation artifact loading
    - Path: repo://pkg/compare/build.go
      Note: Exact pairing and transparent aggregation
    - Path: repo://pkg/gate/policy.go
      Note: Strict product-owned gate policy schema
    - Path: repo://pkg/gate/evaluate.go
      Note: Lexicographic decision evaluation
    - Path: repo://pkg/report/render.go
      Note: Deterministic human promotion report
    - Path: repo://pkg/report/types.go
      Note: Non-applying machine-readable promotion plan
    - Path: repo://cmd/ragopt/commands/compare/compare.go
      Note: Glazed comparison command
    - Path: repo://cmd/ragopt/commands/report/report.go
      Note: Glazed report command
ExternalSources: []
Summary: Implemented strict paired comparison, explicit denominators, product-authored lexicographic gates, deterministic promotion reports, and non-applying review plans.
LastUpdated: 2026-08-06T12:04:02-04:00
WhatFor: Explain how immutable evaluation cells become an evidence-backed pass or fail decision and a human-owned promotion review without automatic mutation.
WhenToUse: Read after the paired-evaluation contract and before implementing a product gate policy or reviewing a candidate decision.
---

# Paired Comparison, Gate, and Promotion Report v1 Contract

## Purpose and trust boundary

Phase 4 turns immutable evaluation cells into a decision without pretending
that a mean score is sufficient evidence. The implementation owns four narrow
mechanisms:

- strict read-only loading of an existing run;
- exact incumbent/candidate pairing and transparent aggregation;
- product-authored, lexicographically evaluated gates;
- a review document and a non-applying promotion plan.

It does not score native answers, choose universal quality metrics, generate a
candidate, apply a changed asset, or deploy a product. Product integrations
remain responsible for those operations.

```text
immutable run directory
        |
        v
eval.LoadArtifactRun       validates custody and identities
        |
        v
compare.Build              exact pairs + raw deltas + denominators
        |
        v
gate.Evaluate              identity -> hard -> target -> regression -> cost
        |
        v
report.Build / Write       Markdown review + review_required JSON plan
```

The important design rule is that each arrow consumes a validated value. A
caller cannot use the reporting package to launder arbitrary cells into a
promotion decision.

## Source map

Read the implementation in this order:

- `pkg/eval/artifacts.go`: strict read-only artifact loading;
- `pkg/eval/runner.go`: `RunConfig` and identities captured at execution time;
- `pkg/compare/types.go`: paired and aggregate schemas;
- `pkg/compare/build.go`: strict join and aggregation;
- `pkg/gate/policy.go`: policy schema, parser, and validation;
- `pkg/gate/evaluate.go`: lexicographic decision algorithm;
- `pkg/report/types.go`: machine-readable promotion plan;
- `pkg/report/render.go`: deterministic Markdown;
- `pkg/report/write.go`: explicit atomic output publication;
- `cmd/ragopt/commands/compare/compare.go`: structured inspection command;
- `cmd/ragopt/commands/report/report.go`: explicit report-output command.

Executable specifications live in:

- `pkg/compare/build_test.go`;
- `pkg/gate/evaluate_test.go` and `pkg/gate/testdata/*.golden`;
- `pkg/report/report_test.go`;
- `pkg/eval/runner_test.go`;
- `cmd/ragopt/main_test.go`.

## Strict artifact loading

`eval.LoadArtifactRun(ctx, directory)` is the only supported disk-to-Phase-4
entry point. It opens the run through `runstore.Open`, which already verifies
the manifest, status, canonical configuration digest, copied-input manifest,
every copied input's bytes and size, path confinement, and terminal summary.

It then performs evaluation-specific checks:

1. strictly decode `config.json` as `ragopt-eval-run/v1`;
2. validate all candidate, snapshot, asset, suite, and policy identities;
3. require the exact configured input-role set—neither missing nor extra;
4. load the copied suite and compare its semantic digest;
5. locate the copied policy and compare its byte digest;
6. strictly parse every newline-committed cell;
7. reject blank middle lines, a truncated tail, duplicate cells, and unknown
   coordinates;
8. recompute and compare every native artifact digest and size.

The loader accepts active, complete, and failed runs for inspection. An active
run can therefore produce an explicit incomplete-pairing failure, but the
loader never truncates, resumes, or repairs it.

`RunConfig` captures the report identity at run creation:

```go
type RunConfig struct {
    SuiteDigest       string
    PolicyDigest      string // exact copied policy bytes
    CandidateID       string
    CandidateDigest   string
    ParentSnapshot    string
    ChildSnapshot     string
    IncumbentArm      string
    ChallengerArm     string
    Repeats           int
    InputDigests      map[string]string
    Mutation          candidate.MutationDeclaration
    ChangedAsset      string
    ParentAssetDigest string
    ChildAssetDigest  string
}
```

The hypothesis and asset diff are captured from a strictly validated
one-mutation candidate. Reporting does not reconstruct intent from filenames
or compare mutable source files after the run.

## Exact pairing

The complete cell identity is:

```text
(run ID, suite digest, policy byte digest, candidate ID, snapshot digest,
 case ID, repeat index, arm)
```

The comparison coordinate is the projection:

```text
(case ID, repeat index)
```

For each coordinate, exactly one incumbent and one challenger cell form a
`Pair`. An absent side creates a `MissingPair`; it never creates a zero score.
A duplicate side or cross-run identity is an error because the correct value
would be ambiguous.

Pseudocode:

```text
index = map[(case, repeat, arm)]cell

for cell in committed cells:
    require cell identity matches run config
    require coordinate belongs to suite and repeat range
    require coordinate not already in index
    index[coordinate] = cell

for case in suite order:
    for repeat in 0..repeats-1:
        incumbent = index[(case, repeat, incumbentArm)]
        candidate = index[(case, repeat, candidateArm)]
        if either absent:
            append explicit MissingPair
        else:
            append Pair with original cells and raw deltas
```

`Pair` retains the two complete cells. `MetricDelta` adds, but never replaces,
the native values:

```text
delta = candidate metric - incumbent metric
```

Positive is therefore always “candidate higher.” Whether higher is desirable
is encoded by the product's metric and thresholds, not guessed by `ragopt`.

## Denominators and aggregates

Three denominators must not be conflated:

- `ExpectedPairs`: suite cases multiplied by repeats;
- `CompletePairs`: coordinates with both arms;
- `PairsWithMetric`: complete pairs where both arms expose one metric.

Group `all` is implicit and reserved. Product case groups are copied from the
suite. For each group, the comparison reports completion, contract validity,
failures, abstentions, and mean cost deltas separately from metric means.

For each metric/group combination it reports:

- incumbent and candidate means over paired metric values;
- mean raw delta;
- wins, ties, and losses by sign of raw delta;
- all three relevant pair counts.

A failed cell remains a complete pair if both arms produced committed cells.
Its failure count increases, but no absent quality metric is synthesized. A
missing cell reduces complete pairs and creates a missing-pair record. Both
conditions fail appropriate gates without being silently excluded.

## Gate-policy API

Policies are product-owned YAML files. `gate.LoadPolicy` uses known-field
decoding and permits exactly one YAML document. It computes two identities:

- `ByteDigest`: the exact bytes copied into the evaluation run;
- `Digest`: canonical JSON of the validated semantic policy.

The byte digest proves the run used this exact file. The semantic digest names
the parsed policy in a decision.

```yaml
api_version: ragopt-gate-policy/v1
name: product-policy-v1
hard_gates:
  require_all_cells: true
  require_completed: true
  require_contract_valid: true
  max_failure_rate: 0
  metric_floors:
    product_safety: 0.95
target:
  metric: product_quality
  groups: [comparison]
  minimum_mean_delta: 0.05
  require_positive_each_repeat: true
regressions:
  maximum_case_delta:
    product_safety: -0.10
  maximum_mean_delta:
    all:
      product_relevance: -0.01
tie_breakers: [provider_calls, tool_calls, total_tokens, duration]
```

This is schema shape, not a recommended threshold. The repository contains no
universal faithfulness, relevance, safety, latency, or token default. Every
product must define and review its own metric semantics.

## Lexicographic decision

`gate.Evaluate` is pure: it performs no I/O and changes no candidate or run.
It evaluates phases in order and stops after the first failing phase.

```text
1 identity
  exact policy bytes and complete pairing

2 hard
  completion, contract, failure rate, and absolute candidate floors

3 target
  presence, minimum mean improvement, optional positivity per repeat

4 regression
  minimum allowed per-case and group-mean deltas

5 tie_break
  ordered informational candidate-minus-incumbent cost deltas
```

Tie-breakers cannot rescue a quality failure. V1 records them after all quality
gates pass; it does not invent cost thresholds or collapse cost and quality
into a weighted score.

Every evaluated check contains phase, name, pass/fail, explanatory message,
and optional machine-readable values. `Decision.Reasons` retains messages from
the first failing phase. A rejected candidate and its run remain immutable
evidence.

## Reports and promotion plans

`report.Build` requires mutually consistent run, comparison, policy, and
decision identities. It deterministically renders:

- run, suite, policy, candidate, and snapshot identities;
- hypothesis, expected improvement, regression risks, and exact asset digests;
- every evaluated gate check;
- group outcome denominators;
- metric means, deltas, and win/tie/loss counts;
- every missing pair;
- a per-pair table with failures, metric deltas, and cost deltas;
- an explicit human review action.

The JSON plan has schema `ragopt-promotion-plan/v1` and fixed safety fields:

```json
{
  "api_version": "ragopt-promotion-plan/v1",
  "state": "review_required",
  "human_apply_required": true,
  "decision": "pass"
}
```

Even a passing decision is not an automatic promotion. The plan describes the
asset replacement and evidence identities so product-owned tooling or a human
can review it. There is deliberately no `Apply` API.

`report.Write` requires distinct explicit Markdown and JSON paths. Each file
is written through a same-directory temporary file, synced, renamed, and its
directory synced. It does not write into the immutable evaluation run unless a
caller explicitly chooses such a path; the CLI documentation directs callers
to separate review outputs.

## CLI API

Read-only structured comparison:

```bash
go run ./cmd/ragopt compare \
  --run ./runs/<run-id> \
  --format json
```

Rows have `record_type` values `decision`, `check`, `group`, `metric`, and
`missing_pair`. Table, JSON, JSONL, CSV, TSV, and YAML serialization are
provided by Glazed.

Explicit report publication:

```bash
go run ./cmd/ragopt report \
  --run ./runs/<run-id> \
  --output-path ./reviews/candidate-1.md \
  --plan-path ./reviews/candidate-1.plan.json \
  --format json
```

The result row contains run ID, candidate ID, decision, absolute output paths,
and `human_apply_required=true`. Both commands expose the shared Glazed output
flags and root logging flags; neither reads environment variables.

## Failure examples

- Candidate arm failed but committed a valid failure cell: pairing completes;
  hard completion/failure gates fail.
- Candidate cell is deleted: identity `complete_pairing` fails before hard or
  quality checks. Deletion cannot improve the decision.
- Metric missing on one complete pair: metric presence/floor/regression check
  fails with the actual denominator.
- Policy edited after evaluation: loaded policy byte digest differs from the
  run; identity fails.
- Cell copied from another suite, policy, snapshot, candidate, or run:
  comparison rejects it as cross-run identity.
- Native artifact changed after execution: artifact loading fails before
  comparison.

## Validation commands

```bash
go test ./...
go test ./pkg/eval ./pkg/compare ./pkg/gate ./pkg/report -race -count=1
go build ./...
go run ./cmd/ragopt compare --help
go run ./cmd/ragopt report --help
```

The six decision goldens cover pass, hard-gate failure, target failure,
catastrophic regression, tie-break reporting, and incomplete pairing. Report
tests prove deterministic rendering, atomic publication, fixed human-review
state, and identity mismatch rejection.

## Review checklist

- Confirm product metrics have documented direction and range.
- Confirm the copied policy digest equals the run policy digest.
- Check expected, complete, and metric-present denominators independently.
- Read failures and contract-validity counts before quality means.
- Inspect every per-case regression and native artifact named by the product.
- Treat cost tie-breakers only after all earlier phases pass.
- Confirm the plan says `review_required` and `human_apply_required=true`.
- Apply an accepted asset only through the product's reviewed workflow.

## Deliberately deferred

Phase 4 does not add a weighted multi-objective optimizer, Pareto frontier,
policy editor, report web UI, database, candidate registry, plugin protocol,
automatic apply, or deployment integration. Those would be new features, not
requirements for the proven evidence path.
