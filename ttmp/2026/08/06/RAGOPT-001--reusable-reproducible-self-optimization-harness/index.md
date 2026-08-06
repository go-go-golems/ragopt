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
LastUpdated: 2026-08-06T16:15:00-04:00
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
- [Shared production refresh, scheduling, and resumability design](design-doc/02-production-index-build-scheduling-resumability-and-ragopt-integration.md)
- [Implementation diary](reference/01-implementation-diary.md)
- [Runstore v1 contract](reference/02-runstore-v1-on-disk-contract-and-recovery-guarantees.md)
- [Snapshot and candidate v1 contract](reference/03-snapshot-and-candidate-v1-bundle-contract.md)
- [Paired evaluation v1 contract](reference/04-paired-evaluation-v1-api-cell-and-resume-contract.md)
- [Comparison, gate, and promotion report v1 contract](reference/05-paired-comparison-gate-and-promotion-report-v1-contract.md)
- [First corrected RAG-TTC I5 rejection](reference/06-first-rag-ttc-i5-feedback-proof-rejection-report.md)
- [Fresh-root reproduction protocol](reference/07-rag-ttc-i5-fresh-root-reproducibility-protocol.md)
- [Second fresh-root proof and decision](reference/08-second-rag-ttc-i5-fresh-root-proof-and-reproducibility-decision.md)
- [Phased task ledger](tasks.md)
- [Changelog](changelog.md)

## Status

Current status: **active — RAG-TTC proof reproduced; GEC integration next**

The reusable custody, candidate, paired-run, comparison, gate, and report path
is implemented and exercised with deterministic fixtures. Phase 5 has now
proved the boundary twice in RAG-TTC; GEC-RAG remains required before v0.1.

The current RAG-TTC proof exercises the product-owned
`rag-ttc tool-loop ragopt` evaluation adapter. It does not modify or evaluate
the canonical Admin Chat/sessionstream cutover, and it does not own the
customer-facing TTC Garden Assistant. Those products may add their own
product-owned adapters later if their exact serving runtimes need optimization.

Phase 5 task 67 is complete. The authorized six-cell fresh-root execution
reproduced all identities, copied-input custody, complete pairing, and the
`fail` decision. The I5 candidate remains rejected; validation was correctly
left unrun. Reference 08 records the exact outcomes, budgets, stochastic
differences, and proof decision. Tasks 68–71 now move to GEC-RAG, beginning
with the known authorization and judge-accounting P0 findings.

The post-v0.1 production-refresh design makes `ragopt` the shared fixed
control-plane protocol for CoinVault/GEC, the TTC Garden Assistant, and the
RAG-TTC Admin Chat. Product repositories retain extraction, indexing, serving,
native evaluation, scheduling, and deployment. This design is documented and
tasked, not implemented.

## Scope Boundary

V1 owns artifact identity, custody, pairing, comparison, gates, and reports. It
does not own retrieval, chat, judging, transcript warehousing, candidate
proposal, deployment, or online self-modification.

## Structure

- `design-doc/` — architecture, APIs, pseudocode, decisions, and phase details.
- `reference/` — chronological diary and implemented API contracts.
- `tasks.md` — canonical phase and task ledger.
- `changelog.md` — concise ticket-level changes.
