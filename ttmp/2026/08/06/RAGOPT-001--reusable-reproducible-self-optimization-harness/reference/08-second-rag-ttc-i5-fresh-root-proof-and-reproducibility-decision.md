---
Title: Second RAG-TTC I5 Fresh-Root Proof and Reproducibility Decision
Ticket: RAGOPT-001
Status: active
Topics:
    - rag
    - evaluation
    - reproducibility
    - rag-ttc
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: /home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/rag-ttc/cmd/rag-ttc/cmds/chat/tooleval/ragopt.go
      Note: Product-owned adapter exercised by both proof runs
    - Path: /home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/rag-ttc/configs/ragopt/i5-combined-comparison-v1/candidate.yaml
      Note: Strict one-mutation candidate reproduced and rejected
    - Path: repo://ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/reference/07-rag-ttc-i5-fresh-root-reproducibility-protocol.md
      Note: Frozen protocol and authorization envelope used for this run
ExternalSources: []
Summary: Fresh-root evidence proving stable RAGOPT identity, custody, pairing, and rejection semantics for the RAG-TTC I5 candidate.
LastUpdated: 2026-08-06T16:10:00-04:00
WhatFor: Decide whether the RAG-TTC adapter proof reproduces and whether GEC integration may begin.
WhenToUse: Read before checking Phase 5 task 67 or interpreting the two stochastic I5 runs.
---

# Second RAG-TTC I5 Fresh-Root Proof and Reproducibility Decision

## Executive decision

The RAG-TTC integration proof **passes**. The I5 combined-comparison candidate
**fails** and must not be promoted. These statements are compatible:

- The harness reproduced the same frozen system identities, copied-input
  custody, complete three-pair comparison, gate ordering, and final `fail`
  decision from a fresh run root.
- The stochastic answer runtime produced different cell-level behavior, but
  every difference was retained. Run two produced two fully judged pairs and
  the candidate lost faithfulness on both.
- Validation was not run because the feedback hard gates failed.

Task 67's purpose is to prove the integration and evidence chain, not to force
a candidate to pass. The evidence below satisfies that purpose. The next
sequential work is the GEC authorization and judge-accounting P0 repair, then a
frozen GEC candidate run twice.

## Quick reference

| Field | Result |
|---|---|
| Second run | `20260806T194517.057830400Z-ttc-i5-feedback-98a227142ee5` |
| Run state | `complete` |
| Cells | 6 expected; 6 recorded; 1 failed arm |
| Pairing | 3/3 complete |
| Candidate validity | 2/3 contract-valid |
| Decision | **FAIL — retain and do not promote** |
| Harness proof | **PASS — fresh-root rejection reproduced** |
| Answer calls | 22 of 24 maximum |
| Search calls | 17 of 18 maximum |
| Judge work calls | 8 of 12 maximum |
| Provider tokens | 108,548, below 1,000,000 |
| Validation | Not run |

## What was frozen

Both proof runs have byte-identical `config.json` files:

```text
sha256:99ca5018c34d85b1693894c8670f57f7a00044e5f62065376bcf699be61e627c
```

The important identity layers are:

| Identity | Digest |
|---|---|
| Suite semantic identity | `sha256:b009f9e913179007bf2da04cadb2e75200e22e1ab38689513015721ee76781b6` |
| Suite source bytes | `sha256:c22b6b186d5f4bd7c2a9177271bfa0005946474b67f78fc7ca85108b9106518d` |
| Policy semantic identity | `sha256:75c56ff46aab7bb2b281942514faf8a7de94fa1bb4c06e2e0df9548fb7ddb83f` |
| Policy source bytes | `sha256:7d3fba806e3824b26ff817b103fd6415cfe66117edc66ad3050c6ab6476b4d8e` |
| Candidate | `sha256:0d598d871694580d7e36af0238a3431cf02690249f6484ea4380b00179507106` |
| Parent snapshot | `sha256:4cc350f874416dd38e4288a090b95952a5d3d50de0b2ff8ff1fa0b26ee5b876a` |
| Child snapshot | `sha256:836466c1b61320034b9e88cd2ba5e1656555082f1196d6e500165653bc8e8eaa` |
| Parent search description | `sha256:5d0ff3a5780742be4d04be20e1df70dfccca17706ea034c47f946c3f84b50808` |
| Candidate search description | `sha256:359486f91306b2bb48d705ec1d3b442fe793e480fce882f16ab67cee29fd0e6d` |

The byte/semantic distinction matters. Source-byte digests prove exact copied
inputs. Semantic digests identify canonical decoded content. Pre-v0.1 schemas
currently overload the name `policy_digest`; Phase 6 must rename or split the
fields without a compatibility shim.

## Second-run cell evidence

| Case | Arm | Outcome | Calls / searches | Faithfulness / relevance |
|---|---|---|---:|---:|
| `ttc-y-005` | incumbent | valid answer | 3 / 2 | 1.000000 / 1.000000 |
| `ttc-y-005` | challenger | valid answer | 3 / 2 | 0.944444 / 1.000000 |
| `ttc-y-007` | incumbent | valid answer | 4 / 3 | 0.978723 / 1.000000 |
| `ttc-y-007` | challenger | valid answer | 4 / 3 | 0.950000 / 1.000000 |
| `ttc-expand-056` | incumbent | arm error at four-iteration ceiling | 4 / 4 requested | not judged |
| `ttc-expand-056` | challenger | invalid abstention | 4 / 3 | not judged |

The failed incumbent cell is retained as `arm-error.json`. Its final search
request reached the per-cell embedding-resource ceiling and the bounded answer
loop then reached four iterations. The challenger returned an abstention that
did not satisfy the answer contract. Neither cell is silently omitted.

## Gate result

| Phase | Gate | Result | Explanation |
|---|---|---:|---|
| identity | policy bytes | PASS | Copied policy bytes match the run identity |
| identity | complete pairing | PASS | 3 of 3 expected pairs exist |
| hard | require all cells | PASS | Every coordinate has a durable cell |
| hard | candidate completed | PASS | Candidate completed 3 of 3 |
| hard | candidate contract valid | **FAIL** | Candidate valid for 2 of 3 |
| hard | maximum failure rate | PASS | Candidate failure rate is zero |
| hard | faithfulness coverage/floor | **FAIL** | Only 2/3 pairs have metrics; missing is not dropped |

For the two judged pairs:

- Answer relevance is tied: incumbent 1.0, candidate 1.0, mean delta 0.0.
- Incumbent faithfulness mean is 0.989362.
- Candidate faithfulness mean is 0.947222.
- Candidate faithfulness mean delta is -0.042139, with zero wins, zero ties,
  and two losses.

Hard gates are lexicographic. Cost cannot rescue failed contract validity or
metric coverage. The generated promotion plan therefore has state
`review_required`, `human_apply_required: true`, and decision `fail`. It
describes the candidate but cannot apply it.

## Run-one versus run-two interpretation

| Property | Run one | Run two | Interpretation |
|---|---:|---:|---|
| Recorded cells | 6/6 | 6/6 | Stable custody |
| Complete pairs | 3/3 | 3/3 | Stable pairing |
| Candidate valid | 1/3 | 2/3 | Stochastic runtime changed |
| Fully judged pairs | 0 paired | 2 paired | Coverage improved in run two |
| Decision | fail | fail | Stable gate outcome |
| Validation | not run | not run | Spending gate preserved |

Canonical metric deltas cannot be byte-identical when model outputs differ.
That is not a harness mismatch. The reproducibility claim is narrower and
stronger: frozen identities reproduce the same coordinate set and decision
procedure; native evidence explains every differing outcome; invalid and
failed cells remain visible; the final decision is stable.

This also limits the conclusion. Six cells over three cases do not prove the
I5 description is universally harmful. They prove it did not earn promotion
under the frozen feedback policy. A larger exploratory study would be a new
experiment, not a reinterpretation of this gate.

## Custody and budget verification

The verification recomputed, rather than trusted:

1. Every copied input's SHA-256 and byte size from `inputs/manifest.json`.
2. Both run configurations' byte identity using `sha256sum` and `cmp -s`.
3. Every native `outcome.json` or `arm-error.json` reference against its
   recorded SHA-256 and size.
4. The exact six unique `(case, arm, repeat)` coordinates.
5. The existence of native session directories and answer traces.
6. Answer iterations, search tool calls, judge cache work calls, and provider
   token usage from native artifacts.

Observed external work remained within authorization:

```text
answer generations: 22 / 24
search embeddings:  17 / 18
judge generations:   8 / 12
answer tokens:       84,117
judge tokens:        24,431
combined tokens:     108,548 / 1,000,000
validation cells:    0
```

No source, profile, adapter, candidate, suite, policy, or budget was changed
between the two proof runs.

## Commands for reproduction review

The authorized product command was:

```bash
go run ./cmd/rag-ttc tool-loop ragopt \
  --split feedback \
  --profile ttc-live-luna-low \
  --cache-directory /home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/rag-ttc/.cache/rag-ttc \
  --format json \
  --log-level info
```

The comparison and report commands consume the terminal run and frozen policy:

```bash
go run /home/manuel/code/wesen/go-go-golems/ragopt/cmd/ragopt compare \
  --run experiments/ragopt-runs/20260806T194517.057830400Z-ttc-i5-feedback-98a227142ee5 \
  --policy configs/ragopt/i5-combined-comparison-v1/gate-policy.yaml

go run /home/manuel/code/wesen/go-go-golems/ragopt/cmd/ragopt report \
  --run experiments/ragopt-runs/20260806T194517.057830400Z-ttc-i5-feedback-98a227142ee5 \
  --policy configs/ragopt/i5-combined-comparison-v1/gate-policy.yaml \
  --output-directory experiments/ragopt-reports/20260806T194517.057830400Z-ttc-i5-feedback-98a227142ee5
```

## Phase decision and next work

Check task 67. Do not promote I5 and do not run its validation suite. Proceed
to GEC only in this order:

1. Fix the known authorization and judge-accounting P0 findings.
2. Freeze one human-authored candidate and its locked suite, policy, corpus,
   retrieval configuration, answer model, judge, prompts, and safety ceilings.
3. Run the candidate twice from fresh roots.
4. Preserve GEC-native evidence and project only common comparison fields into
   RAGOPT.
5. Document integration friction before adding a generic RAGOPT API.

The project should still not build a reflector, generic plugin/subprocess
protocol, workflow engine, automatic deployment path, transcript warehouse, or
multi-mutation optimizer. Those features are neither required nor proven by
this checkpoint.
