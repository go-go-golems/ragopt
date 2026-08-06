---
Title: Reusable Reproducible Self-Optimization Harness
Ticket: RAGOPT-001
Status: active
Topics:
    - rag
    - evaluation
    - experiments
    - reproducibility
    - architecture
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Design and implementation ledger for a small reusable Go harness that makes candidate identity, immutable evaluation, paired comparison, gates, and human promotion evidence mandatory.
LastUpdated: 2026-08-06T10:05:00-04:00
WhatFor: Prevent ad hoc experiment commands or unimplemented reflectors from being mistaken for a working self-optimization system.
WhenToUse: Start here before implementing or integrating ragopt; then read the design guide and follow tasks.md in phase order.
---

# Reusable Reproducible Self-Optimization Harness

## Overview

`ragopt` packages the experiment-control spine already demonstrated in
RAG-TTC and found missing in GEC-RAG: semantic identity, immutable run
artifacts, one-mutation candidates, resumable paired execution, explicit
gates, and human-owned promotion reports. Phases 1–4 implement that generic
artifact path; product execution and candidate proposal remain outside it.

## Key Links

- [Intern design and implementation guide](design-doc/01-ragopt-intern-guide-to-a-reusable-evidence-gated-optimization-harness.md)
- [Implementation diary](reference/01-implementation-diary.md)
- [Runstore v1 contract](reference/02-runstore-v1-on-disk-contract-and-recovery-guarantees.md)
- [Snapshot and candidate v1 contract](reference/03-snapshot-and-candidate-v1-bundle-contract.md)
- [Paired evaluation v1 contract](reference/04-paired-evaluation-v1-api-cell-and-resume-contract.md)
- [Comparison, gate, and promotion report v1 contract](reference/05-paired-comparison-gate-and-promotion-report-v1-contract.md)
- [Phased task ledger](tasks.md)
- [Changelog](changelog.md)

## Status

Current status: **active — Phase 4 complete; product proof cycle pending**

The reusable custody, candidate, paired-run, comparison, gate, and report path
is implemented and exercised with deterministic fixtures. Phase 5 must prove
the boundary in RAG-TTC and GEC-RAG before v0.1 is released.

The current RAG-TTC proof exercises the product-owned
`rag-ttc tool-loop ragopt` evaluation adapter. It does not modify or evaluate
the canonical Admin Chat/sessionstream cutover, and it does not own the
customer-facing TTC Garden Assistant. Those products may add their own
product-owned adapters later if their exact serving runtimes need optimization.

## Scope Boundary

V1 owns artifact identity, custody, pairing, comparison, gates, and reports. It
does not own retrieval, chat, judging, transcript warehousing, candidate
proposal, deployment, or online self-modification.

## Structure

- `design-doc/` — architecture, APIs, pseudocode, decisions, and phase details.
- `reference/` — chronological diary and implemented API contracts.
- `tasks.md` — canonical phase and task ledger.
- `changelog.md` — concise ticket-level changes.
