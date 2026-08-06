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
    - Path: abs:///tmp/rag-ttc-ragopt-proof/cmd/rag-ttc/cmds/experiments/answerquality/judge.go
      Note: Exact two-stage judge payloads, cache behavior, and invalid-answer skipping
    - Path: abs:///tmp/rag-ttc-ragopt-proof/configs/ragopt/i5-combined-comparison-v1/shared/gate-policy.yaml
      Note: Exact frozen policy bytes for both proof runs
    - Path: abs:///tmp/rag-ttc-ragopt-proof/pkg/app/chat/tool_runtime.go
      Note: Direct answer-engine calls, cached query embeddings, budgets, and disabled SQL runtime behavior
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
approval for the additional export and spend. Answer-generation calls are
intentionally fresh; query embeddings and judge calls use the shared
content-addressed cache and call providers only on misses.

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

## External Execution Consent Envelope

Approval for task 67 is deliberately narrow. It authorizes **only** one new
feedback run: three cases, two arms, one repeat, six cells. It does not authorize
the seven-case validation split, its two repeats, a different candidate, a
different model, retries after a terminal run, or any production mutation.

### Questions exported

| Case | Question |
|---|---|
| `ttc-y-005` | Compare Blue Ice and Carolina Sapphire Arizona Cypress for site requirements. |
| `ttc-y-007` | How do Thuja Green Giant and Leyland Cypress differ as privacy screens? |
| `ttc-expand-056` | Compare Pink Pearl Black Diamond and Dynamite Crape Myrtle by flower color and mature size. |

These are product benchmark questions, not customer conversations. The search
runtime reads only the frozen TTC customer-facing corpus. Scoped SQL is disabled,
so the run does not send orders, customer records, private logistics facts,
credentials, or live database results.

### Provider payloads

| Stage | Resolved provider/model | What leaves the machine | Cache behavior |
|---|---|---|---|
| Query embedding | `openai` / `text-embedding-3-small` / 1536 | model-generated search query strings | content-addressed; exact hits cause no provider call |
| Answer/tool loop | `openai-responses` / `gpt-5.6-luna` through profile `ttc-live-luna-low` | question, orchestration/tool instructions, conversation state, and admitted public corpus evidence/tool results | intentionally not served from the generation cache; each admitted iteration is a fresh provider call |
| Judge statement extraction | `openai` / `gpt-5.6-luna` through `ttc-judge-statements` | question and completed answer; no evidence | content-addressed; skipped for invalid answers |
| Judge verdict | `openai` / `gpt-5.6-luna` through `ttc-judge-verdicts` | question, admitted evidence passages, and extracted answer statements | content-addressed; skipped for invalid answers |

The verdict judge is the same model family as the answer model; the native
judge record preserves `same_family_verdicts=true`. That is an evaluation risk
already locked into both runs, not a reason to change the judge between them.

### Hard call ceiling

| Operation | Per cell | Six-cell maximum | When provider work is skipped |
|---|---:|---:|---|
| query embeddings | 3 | 18 | exact content-cache hit |
| answer generations | 4 | 24 | never by local content cache; a cell may terminate earlier |
| judge generations | 2 | 12 | invalid/abstained answer or exact content-cache hit |
| **combined heterogeneous operations** | **9** | **54** | according to the rules above |

These are operation ceilings, not a USD ceiling. The adapter sets
`AllowUnpriced: true` and does not configure per-operation dollar estimates or
`MaxEstimatedUSD`, so the command cannot honestly promise a monetary maximum.
Do not infer “free” from a cache hit: the answer loop remains fresh, and remote
provider-side cached-input billing is separate from the local content cache.

For scale, the first corrected run executed 22 fresh answer-model iterations,
recorded 17 search-tool invocations, and performed two uncached judge calls.
Only one of six cells reached the judge because four answers were invalid and
one arm failed. Those observations do not lower the second run's hard ceiling.

### Local retention

The run writes a new immutable directory under `experiments/ragopt-runs/` with:

- copied candidate, suite, policy, prompt, schema, runtime, and judge inputs;
- one cell for every arm/case coordinate, including failures;
- full native chat sessions containing answer text, reasoning text/summary,
  tool payloads, admitted evidence, citations, provider usage, and errors;
- judge statements, verdicts, metrics, cache outcomes, and provider identities;
- artifact digests, terminal status, comparison evidence, and a non-applying
  promotion plan.

Encrypted reasoning is omitted. Transcript limits remain the frozen contract:
12,000 runes per tool payload, 24,000 reasoning runes, at most ten distinct
evidence chunks, and at most 18,000 evidence runes per search configuration.

### Approval text

An unambiguous authorization can be as short as:

> I approve one second RAG-TTC I5 feedback run: six cells only, with ceilings
> of 24 answer calls, 18 embedding calls, and 12 judge calls. Do not run
> validation.

After that approval, launch exactly one new run. A transport failure may be
resumed through the same active run directory; do not silently create another
fresh run or expand the cell count.

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
