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
gates, and human-owned promotion reports. The initial ticket is design and
project initialization only; it intentionally adds no optimizer behavior.

## Key Links

- [Intern design and implementation guide](design-doc/01-ragopt-intern-guide-to-a-reusable-evidence-gated-optimization-harness.md)
- [Implementation diary](reference/01-implementation-diary.md)
- [Phased task ledger](tasks.md)
- [Changelog](changelog.md)

## Status

Current status: **active — design review**

The repository scaffold is complete. Implementation begins only after Phase 0
scope review is accepted.

## Scope Boundary

V1 owns artifact identity, custody, pairing, comparison, gates, and reports. It
does not own retrieval, chat, judging, transcript warehousing, candidate
proposal, deployment, or online self-modification.

## Structure

- `design-doc/` — architecture, APIs, pseudocode, decisions, and phase details.
- `reference/` — chronological implementation diary.
- `tasks.md` — canonical phase and task ledger.
- `changelog.md` — concise ticket-level changes.
