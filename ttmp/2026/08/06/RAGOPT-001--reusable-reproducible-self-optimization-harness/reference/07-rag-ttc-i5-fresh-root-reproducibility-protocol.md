---
Title: RAG-TTC I5 Fresh-Root Reproducibility Protocol
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
    - Path: abs:///tmp/rag-ttc-ragopt-proof/cmd/rag-ttc/cmds/chat/tooleval/ragopt.go
      Note: Locked product adapter, budgets, index default, and environment validation for task 67
    - Path: abs:///tmp/rag-ttc-ragopt-proof/configs/ragopt/i5-combined-comparison-v1/shared/gate-policy.yaml
      Note: Exact frozen policy bytes for both proof runs
    - Path: repo://pkg/eval/runner.go
      Note: Evaluation run policy byte-digest custody and semantic run identity
    - Path: repo://pkg/gate/evaluate.go
      Note: Byte identity check and semantic decision identity
    - Path: repo://pkg/gate/policy.go
      Note: Distinct policy byte and semantic digest computation
ExternalSources: []
Summary: Precommitted identities, custody checks, and decision criteria for the second corrected RAG-TTC I5 feedback run.
LastUpdated: 2026-08-06T15:07:51.868218493-04:00
WhatFor: Prevent post-hoc definitions of reproducibility when completing RAGOPT-001 task 67.
WhenToUse: Read before authorizing, running, comparing, or reviewing the second provider-backed I5 feedback proof.
---


# RAG-TTC I5 Fresh-Root Reproducibility Protocol

## Goal

Freeze the exact invariants that the second corrected RAG-TTC I5 feedback run
must satisfy. The second run is a fresh provider execution, so answer text,
token counts, durations, tool trajectories, native artifact digests, and
quality scores may vary. Reproducibility means semantic identity, complete
custody, correct paired accounting, and the same gate semantics—not byte-for-
byte stochastic transcripts.

## Context

The first corrected run is:

```text
20260806T180004.824520651Z-ttc-i5-feedback-a1781d11159f
```

It completed six expected cells, retained one incumbent arm failure, and
rejected I5 because only one of three candidate cells was contract-valid. The
validation split was correctly not run. Task 67 requests one repetition from a
fresh run root before GEC integration begins.

The second run sends the same three evaluation questions and retrieved evidence
to configured external answer and judge providers. It requires explicit
approval for the additional export and spend.

## Quick Reference

### Identities that must match exactly

| Identity | Required value |
|---|---|
| candidate | `sha256:0d598d871694580d7e36af0238a3431cf02690249f6484ea4380b00179507106` |
| parent snapshot | `sha256:4cc350f874416dd38e4288a090b95952a5d3d50de0b2ff8ff1fa0b26ee5b876a` |
| child snapshot | `sha256:836466c1b61320034b9e88cd2ba5e1656555082f1196d6e500165653bc8e8eaa` |
| feedback suite semantic digest | `sha256:b009f9e913179007bf2da04cadb2e75200e22e1ab38689513015721ee76781b6` |
| feedback suite byte digest | `sha256:c22b6b186d5f4bd7c2a9177271bfa0005946474b67f78fc7ca85108b9106518d` |
| gate policy semantic digest | `sha256:75c56ff46aab7bb2b281942514faf8a7de94fa1bb4c06e2e0df9548fb7ddb83f` |
| gate policy byte digest | `sha256:7d3fba806e3824b26ff817b103fd6415cfe66117edc66ad3050c6ab6476b4d8e` |
| parent search asset | `sha256:5d0ff3a5780742be4d04be20e1df70dfccca17706ea034c47f946c3f84b50808` |
| child search asset | `sha256:359486f91306b2bb48d705ec1d3b442fe793e480fce882f16ab67cee29fd0e6d` |
| index manifest bytes | `sha256:b2e92c7cd9964b6ab36c0b3dc47f96b1cdc72f1966c3ceeaf6772abf82f194dd` |
| product commit | `90485d8539515173e88d2a9703f3a7ad4aed74bc` |

The run config must also contain the same complete `input_digests` map. A new
run ID and timestamps are required and intentionally excluded from semantic
identity.

#### Policy identity naming warning

The current pre-v0.1 schemas use the JSON name `policy_digest` for two related
but different identities:

- `ragopt-eval-run/v1` config, cells, and comparison report contain the policy
  **byte digest** (`7d3f…`) used to prove that the exact copied YAML is unchanged;
- `ragopt-gate-decision/v1` and `ragopt-promotion-plan/v1` contain the policy
  **semantic digest** (`75c5…`) computed from the strictly decoded policy.

That API naming is integration friction, not evidence drift. For task 67, both
values must match the first corrected run at their respective artifact layers.
Do not change the schemas between the two proof executions. Before v0.1, rename
or explicitly split the fields so a caller cannot compare the two by accident;
no backwards-compatibility layer is required.

The suite already separates these layers more clearly: `config.suite_digest`
is the semantic digest (`b009…`), while `input_digests["evaluation-suite"]` is
the copied file's byte digest (`c22b…`).

### Coordinates that must exist exactly once

```text
incumbent  / ttc-y-005      / repeat 0
challenger / ttc-y-005      / repeat 0
incumbent  / ttc-y-007      / repeat 0
challenger / ttc-y-007      / repeat 0
incumbent  / ttc-expand-056 / repeat 0
challenger / ttc-expand-056 / repeat 0
```

### Custody requirements

- run state is terminal `complete`;
- summary says `expected_cells=6` and `completed_cells=6`;
- comparison says `expected_pairs=3` and `complete_pairs=3`;
- every cell references a run-relative native artifact;
- every native artifact exists and matches its recorded SHA-256 and size;
- run input copies match the config's complete digest map;
- no result is silently omitted because it failed, abstained, or was invalid;
- no validation cell exists.

### Decision requirements

The exact stochastic cell statuses are observations, not preconditions. Apply
the same frozen policy to the second run. Reproduction succeeds if the harness:

1. passes identity and complete-pairing checks;
2. accounts for all invalid/failed cells in denominators;
3. evaluates every hard gate before target/cost tie-breakers;
4. produces a deterministic decision from that run's evidence;
5. emits a non-applying, human-reviewed promotion plan;
6. makes any decision difference explainable from paired native outcomes.

The expected decision class is `fail`, because the same candidate previously
had poor contract-valid coverage. If the second decision differs, do not call
the harness non-reproducible automatically and do not promote. Inspect the
paired native evidence and judge coverage, document stochastic instability,
and decide whether the three-case feedback suite is sufficiently stable.

## Usage Examples

### No-provider readiness evidence

The following checks completed on 2026-08-06 without a new provider call:

- detached proof worktree HEAD is exactly `90485d8539515173e88d2a9703f3a7ad4aed74bc`;
- the proof worktree's `profiles.yaml` and the source workspace copy both have
  SHA-256 `72c197d530ccc5575b545ad890163f8cc6415602306977bf03b13838cf24c596`;
- the default index symlink
  `.cache/rag-ttc/indexes/ttc-056cbd53e148922e847ceabab1f7c4ef`
  resolves to the established source-worktree index and its manifest byte
  digest is `b2e92c…`;
- focused adapter, root-command, and asset tests pass;
- strict candidate validation reproduces candidate, parent, child, and changed
  asset identities;
- read-only comparison of the first corrected run still yields `fail`, three
  expected pairs, and three complete pairs;
- no RAGOPT proof tmux session is running; only the unrelated
  `ttc-garden-chat-luna-low` session exists;
- the proof run and report roots currently contain only the first corrected run,
  so a second launch must create a new run ID rather than resume the old one.

Run only after explicit provider approval:

```bash
cd /tmp/rag-ttc-ragopt-proof
go run ./cmd/rag-ttc tool-loop ragopt \
  --split feedback \
  --profile ttc-live-luna-low \
  --cache-directory /home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/rag-ttc/.cache/rag-ttc \
  --format json \
  --log-level info
```

Then run the existing read-only commands against the new run root:

```bash
go run github.com/go-go-golems/ragopt/cmd/ragopt compare \
  --run experiments/ragopt-runs/<new-run-id> --format json

go run github.com/go-go-golems/ragopt/cmd/ragopt report \
  --run experiments/ragopt-runs/<new-run-id> \
  --output-path experiments/ragopt-reports/<new-run-id>/review.md \
  --plan-path experiments/ragopt-reports/<new-run-id>/promotion-plan.json \
  --format json
```

The first baseline comparison is `fail`, with three complete pairs, candidate
completion 3/3, candidate contract validity 1/3, zero candidate arm failures,
and faithfulness coverage for only one pair.

## Related

- [First corrected rejection report](./06-first-rag-ttc-i5-feedback-proof-rejection-report.md)
- [Implementation diary](./01-implementation-diary.md)
- [Phase 5 task ledger](../tasks.md)
