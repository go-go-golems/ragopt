---
Title: First RAG-TTC I5 Feedback Proof Rejection Report
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
    - Path: /home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/rag-ttc/cmd/rag-ttc/cmds/chat/tooleval/ragopt.go
    - Path: /home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/rag-ttc/configs/ragopt/i5-combined-comparison-v1/candidate.yaml
ExternalSources: []
Summary: Evidence, cell outcomes, gate decision, and operational lessons from the corrected first product proof.
LastUpdated: 2026-08-06T14:03:18.131169279-04:00
WhatFor: ""
WhenToUse: ""
---

# First RAG-TTC I5 Feedback Proof Rejection Report

## Goal

Record the first corrected, real-product use of ragopt against RAG-TTC and the
evidence-based decision not to promote the I5 combined-comparison search
description. This is both a candidate review and a proof that the harness runs
the genuine product, preserves native evidence, and fails closed.

## Context

The candidate changes one text asset: the TTC search tool description. Its
hypothesis is that asking for both comparison subjects and requested
attributes in the first search will preserve answer quality while reducing
follow-up work. The parent, candidate, three-case feedback suite, judge,
answer model, embedding model, corpus, index, prompts, schema, tool safety,
budgets, evaluator source, and gate policy are immutable inputs.

An earlier run is explicitly excluded from evaluation because its adapter
allowed only one query embedding. That budget caused five artificial
abstentions and favored the combined-query candidate. Commit `90485d8` corrected
the budget to the product's real maximum of three searches and separated judge
overhead from product cost. The report below covers only the fresh corrected
run.

## Quick Reference

| Field | Value |
|---|---|
| Decision | **FAIL — do not promote** |
| Run | `20260806T180004.824520651Z-ttc-i5-feedback-a1781d11159f` |
| Run state | `complete` |
| Cells | 6/6 committed; one arm failure |
| Candidate | `sha256:0d598d871694580d7e36af0238a3431cf02690249f6484ea4380b00179507106` |
| Parent | `sha256:4cc350f874416dd38e4288a090b95952a5d3d50de0b2ff8ff1fa0b26ee5b876a` |
| Child | `sha256:836466c1b61320034b9e88cd2ba5e1656555082f1196d6e500165653bc8e8eaa` |
| Suite | `sha256:b009f9e913179007bf2da04cadb2e75200e22e1ab38689513015721ee76781b6` |
| Policy | `sha256:75c56ff46aab7bb2b281942514faf8a7de94fa1bb4c06e2e0df9548fb7ddb83f` |
| Changed asset | `search_description` only |

### Cell evidence

| Case | Arm | Product result | Calls / searches | Judge |
|---|---|---|---:|---|
| `ttc-y-005` | incumbent | Failed at four-iteration ceiling | unavailable after arm-error projection | not run |
| `ttc-y-005` | challenger | Contract-valid answer | 4 / 3 | faithfulness 1.0; relevance 1.0 |
| `ttc-y-007` | incumbent | Contract-invalid abstention | 4 / 3 | invalid; not judged |
| `ttc-y-007` | challenger | Contract-invalid abstention | 3 / 2 | invalid; not judged |
| `ttc-expand-056` | incumbent | Contract-invalid abstention | 3 / 2 | invalid; not judged |
| `ttc-expand-056` | challenger | Contract-invalid abstention | 4 / 3 | invalid; not judged |

The failed incumbent cell has a durable ragopt `arm-error.json`. Successful
runtime invocations have a native `outcome.json` which links the complete TTC
session archive, runtime projection, judge cell, judge configuration, cache
outcomes, and usage. ragopt validated every referenced artifact digest and
size while loading the run for comparison.

### Gate decision

| Phase | Check | Result | Reason |
|---|---|---:|---|
| identity | policy bytes | PASS | Copied policy matches the run identity |
| identity | complete pairing | PASS | Three expected pairs are present |
| hard | all cells | PASS | Three candidate cells paired with three incumbent cells |
| hard | candidate completed | PASS | Candidate completed 3/3 |
| hard | candidate contract-valid | **FAIL** | Candidate valid for only 1/3 cases |
| hard | candidate failure rate | PASS | Candidate failures 0/3 |
| hard | faithfulness floor | **FAIL** | Score is 1.0, but exists for only one of three pairs |

No target improvement or cost tie-breaker can override a hard-gate failure.
There are no defensible paired metric aggregates because the incumbent has no
judged cell and the candidate has metrics for only one case.

## Usage Examples

The proof was run in an isolated worktree at corrected commit `90485d8`:

```bash
go run ./cmd/rag-ttc tool-loop ragopt \
  --split feedback \
  --profile ttc-live-luna-low \
  --cache-directory /home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/rag-ttc/.cache/rag-ttc \
  --format json \
  --log-level info
```

The immutable run was strictly loaded and gated:

```bash
go run github.com/go-go-golems/ragopt/cmd/ragopt compare \
  --run experiments/ragopt-runs/20260806T180004.824520651Z-ttc-i5-feedback-a1781d11159f \
  --format json

go run github.com/go-go-golems/ragopt/cmd/ragopt report \
  --run experiments/ragopt-runs/20260806T180004.824520651Z-ttc-i5-feedback-a1781d11159f \
  --output-path experiments/ragopt-reports/20260806T180004.824520651Z-ttc-i5-feedback-a1781d11159f/review.md \
  --plan-path experiments/ragopt-reports/20260806T180004.824520651Z-ttc-i5-feedback-a1781d11159f/promotion-plan.json \
  --format json
```

The promotion plan is `review_required` and `human_apply_required`. It contains
no mechanism that modifies RAG-TTC.

### Decision and next action

Do not run the 28-cell validation split for this candidate. The staged policy
uses feedback as a cheap hard-gate screen; validation spend is justified only
after feedback passes. Retain both the rejected candidate and corrected run.

The product result suggests a future diagnostic question—why comparison cases
often exhaust the four-call loop or abstain—but that is RAG-TTC product work,
not a reason to weaken ragopt gates, increase safety ceilings, or mutate a
second surface in this candidate.

## Related

- [Primary intern guide](../design-doc/01-ragopt-intern-guide-to-a-reusable-evidence-gated-optimization-harness.md)
- [Implementation diary](01-implementation-diary.md)
- [Paired comparison and gate contract](05-paired-comparison-gate-and-promotion-report-v1-contract.md)
- RAG-TTC adapter commit `56432fd`
- RAG-TTC correction commit `90485d8`
