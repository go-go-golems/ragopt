---
Title: Implementation Diary
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
    - Path: abs:///home/manuel/code/gec/2026-03-16--gec-rag/ttmp/2026/08/05/GEC-RAG-OPT-001--retrieval-optimization-reranker-eval-growth-and-benchmarked-retrieval-experiments/design-doc/02-deep-review-from-ad-hoc-retrieval-tuning-to-reproducible-self-optimization.md
      Note: Gap analysis that motivated the ragopt scope
    - Path: abs:///home/manuel/code/wesen/go-go-golems/go-template/go.mod
      Note: Committed template module used for the initial scaffold
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/rag-ttc/pkg/experiment/run.go
      Note: Implemented immutable run and synced JSONL evidence audited in Step 2
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/rag-ttc/pkg/rag/tooleval/runner.go
      Note: Implemented small Arm and Outcome comparison boundary audited in Step 2
    - Path: repo://go.mod
      Note: Normalized ragopt module and dependency baseline
ExternalSources: []
Summary: Chronological record of the ragopt repository bootstrap, evidence audit, scope decisions, documentation, validation, and delivery.
LastUpdated: 2026-08-06T10:05:00-04:00
WhatFor: Preserve the exact reasoning, commands, failures, and review guidance behind the initial ragopt design.
WhenToUse: Read when reviewing the design or resuming implementation from the phased task ledger.
---


# Implementation Diary

## Goal

Record the creation of the `ragopt` repository and the evidence-backed design
for a reusable optimization harness. The diary distinguishes work completed in
this ticket from mechanisms merely proposed for later implementation.

## Step 1: Bootstrap the repository from the committed Go template

The target repository already existed on GitHub with a clean initial commit,
its own remote, and a docmgr seed. I imported the committed `HEAD` of the local
`go-template`, deliberately excluding the template's `.git` directory and its
unrelated dirty `AGENT.md` working-tree change. I then replaced template
identifiers with `ragopt`, renamed the command directory, and kept the command
and package bodies empty.

This checkpoint establishes only proven project plumbing: Go module metadata,
Glazed and zerolog/logcopter dependencies, CI, linting, release configuration,
and docmgr. No optimization behavior was added.

**Commit (code):** `fc3a703` — "Initialize ragopt from go-template"

### Prompt Context

**User prompt (verbatim):**

> ok, created it in /home/manuel/code/wesen/go-go-golems/ragopt . Actually initialize it with the tempalte in ~/code/wesen/go-go-golems/go-template which already sets up the glazed framework and such.
>
> Don't add new features, try not to overengineer, keep it to what we already built and proved to work.
>
> Create a new docmgr ticket in the ragopt repo ttmp, and Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.
>
> Create phases of implementation and for each phase, a precise list of tasks so that we can keep track of our progress. Keep a detailed diary according to the diary skill format as you work, commit at appropriate intervals.

**Assistant interpretation:** Initialize the existing `ragopt` Git repository
from the standard committed Go template, add no optimizer implementation in
this pass, and create an evidence-backed design package with an actionable task
ledger and detailed diary.

**Inferred intent:** Prevent another design-only or placeholder
"self-optimization" effort by giving a new repository a small, testable scope
derived from mechanisms that already produced trustworthy RAG experiments.

### What I did

- Inspected the target's status, remote, history, and worktree list.
- Fetched `origin` and confirmed `main...origin/main` was `0 0`.
- Inspected the template's status and found an unrelated local `AGENT.md` diff.
- Archived the template's committed `HEAD` and extracted it into `ragopt`.
- Preserved `ragopt/.git`, `.ttmp.yaml`, and `ttmp/`.
- Renamed `cmd/XXX` to `cmd/ragopt` and normalized module, release, Makefile,
  logcopter, workflow, and README identifiers.
- Removed inherited trailing whitespace detected by `git diff --cached --check`.
- Ran generation, formatting, tests, and build validation.

### Why

- Copying the template's `.git` would have replaced the new repository's
  identity and remote.
- Copying its dirty working tree would have imported an unrelated policy
  change that conflicted with the instructions supplied for this task.
- A separate scaffold commit makes later documentation and implementation
  reviewable independently.

### What worked

- `GOWORK=off go generate ./...` completed with no output or errors.
- `GOWORK=off go fmt ./...` completed successfully.
- `GOWORK=off go test ./...` passed for the root, command, and package stubs.
- `GOWORK=off go build ./...` passed.
- The first commit contains only normalized template files.

### What didn't work

- The first staged whitespace check reported trailing spaces in three inherited
  workflow files, including:
  `github/codeql-action/analyze@v4 ` and several whitespace-only lines. I
  removed those spaces and reran `git diff --cached --check` successfully.

### What I learned

- The local template is intentionally minimal; it supplies project plumbing
  and dependencies, not a prebuilt Glazed command tree.
- `ragopt` already had a correctly initialized docmgr workspace, so running a
  destructive or redundant initialization was unnecessary.

### What was tricky to build

The important ownership boundary was repository identity. The template's files
were desired, while its Git metadata, remote, docmgr store, and dirty local
change were not. Using the committed archive made that boundary explicit.

### What warrants a second pair of eyes

- Review `.goreleaser.yaml` and the disabled docs publication block before the
  first release; they are normalized template policy, not exercised release
  evidence.
- Confirm the inherited Go/toolchain versions remain the organization default.

### What should be done in the future

- Replace the empty command/package stubs only as each implementation phase
  begins; do not pre-scaffold speculative packages.

### Code review instructions

- Review commit `fc3a703` from the repository root.
- Run `GOWORK=off go test ./...` and `GOWORK=off go build ./...`.
- Search for unresolved template identifiers with
  `rg --hidden 'XXX|go-template' --glob '!.git/**' --glob '!ttmp/**'`.

### Technical details

```text
Target remote: git@github.com:go-go-golems/ragopt.git
Template commit: a7d7972 (Harden template glazed-lint target)
Scaffold commit: fc3a703
Module: github.com/go-go-golems/ragopt
Command path: ./cmd/ragopt
```

## Step 2: Audit the proven experiment mechanisms and the failed framework boundary

This step reviewed the GEC-RAG handoff and deep review, the RAG-TTC
optimization designs and diaries, the project knowledge map, and the concrete
RAG-TTC packages that actually ran experiments. The evidence supports a much
narrower repository than the earlier GEPA design: `ragopt` should package the
measurement and promotion spine that worked, while leaving product execution,
retrieval, judging, diagnostics, and candidate proposal with their owners.

The critical finding is negative but useful. RAG-TTC implemented immutable run
custody, typed evaluation arms, strict semantic config, durable native
artifacts, and deterministic diagnostic packets. The GEPA optimizer ticket
designed a reflector, transcript warehouse, candidate system, and promotion
loop, but its own diary records that the work was documentation and design
only. GEC-RAG then implemented candidate experiments as command branches
without immutable run identity, first-class candidates, paired repeats, or a
safe mutable tool-description asset. `ragopt` must not blur those categories.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Reconcile the earlier RAG-TTC and GEC-RAG work
into the minimum standalone package whose implementation can be demonstrated,
not merely described.

**Inferred intent:** Make experimental discipline the default path so a future
colleague cannot skip custody, comparison, and gates while still calling a
command an optimizer.

### What I did

- Read the GEC-RAG optimizer handoff and its adversarial review guide.
- Re-read the deep review and pragmatic chatbot blueprint, focusing on the
  immutable-run, candidate, pairing, gate, and reflector findings.
- Read the RAG-TTC GEPA optimizer guide, task ledger, and complete diary.
- Inspected `pkg/experiment`, `pkg/rag/tooleval`, `pkg/rag/diagnostic`, and
  `pkg/rag/toolconfig` implementations and tests.
- Read the `rag-ttc` project knowledge map and located the historical
  researchctl/RAG project reports.
- Defined an evidence classification: implemented and exercised; implemented
  but product-specific; designed only; explicitly deferred.

### Why

- A reusable repository should extract stable mechanisms, not copy an entire
  product or repeat a speculative architecture.
- The term self-optimization is unsafe unless the evaluator and promotion
  path can prove what changed, what stayed fixed, and why a candidate passed.

### What worked

- The implemented run store already demonstrates atomic JSON artifacts,
  append-and-sync JSONL, copied input digests, terminal status, and path
  confinement.
- `tooleval.Arm` and its small `Outcome` prove that a shared evaluator can
  compare native systems without replacing their native artifact formats.
- `diagnostic.ValidateSources` proves that cross-artifact identity and count
  mismatches can be rejected instead of silently normalized.
- `toolconfig.Load` proves strict YAML, repository-relative path safety,
  resolved asset digests, and compiled ceilings outside experiment control.

### What didn't work

- No reusable optimizer package exists in RAG-TTC; searches found the proposed
  design ticket but no `pkg/optimizer` implementation.
- The earlier GEPA task ledger leaves every implementation phase unchecked.
- GEC-RAG's current command-oriented experiments do not create immutable
  evaluation runs or first-class candidates and do not pair repeated answer
  outcomes.

### What I learned

- The honest v1 product is an **evidence-gated optimization harness**, not an
  autonomous optimizer.
- A transcript warehouse is useful but is not required for candidate custody,
  paired execution, comparison, or promotion; it is therefore outside v1.
- A reflector is a proposer, not the optimization framework. It must remain
  optional until a human-authored candidate traverses the same path twice.
- Native application artifacts should remain authoritative. The common
  outcome is a comparison projection with a link to the native evidence.

### What was tricky to build

The largest design trap is making a standalone binary responsible for running
arbitrary product code. A subprocess or plugin protocol would be a new,
unproved compatibility surface. The pragmatic boundary is a Go library that
consumer commands call for execution, plus a standalone CLI for validating,
inspecting, comparing, and reporting portable artifacts.

### What warrants a second pair of eyes

- Confirm the package remains domain-neutral enough for GEC-RAG and RAG-TTC
  without becoming a generic workflow engine.
- Review the proposed snapshot and candidate contracts: they must be strict
  while allowing each product to declare its required locked dimensions.
- Review which aggregate gate policies can be generic and which must remain
  product-owned.

### What should be done in the future

- Implement phases in order and stop after each exit criterion.
- Demonstrate two byte-identifiable runs of one human-authored candidate before
  accepting any model-generated proposer work.

### Code review instructions

- Read the primary design guide's evidence matrix before its proposed APIs.
- Cross-check claims against the absolute paths in the References section.
- Verify the RAG-TTC optimizer ticket's unchecked task ledger before treating
  any reflector or warehouse behavior as implemented.

### Technical details

```text
Implemented and reusable:
  immutable run custody -> rag-ttc/pkg/experiment
  small arm projection  -> rag-ttc/pkg/rag/tooleval
  source consistency    -> rag-ttc/pkg/rag/diagnostic
  semantic config       -> rag-ttc/pkg/rag/toolconfig

Designed, not proved:
  transcript warehouse
  candidate store
  reflector
  promotion automation

ragopt v1 boundary:
  identity + custody + resume + paired compare + gates + promotion report
```

## Step 3: Write the reusable harness design and phased implementation guide

This step converted the evidence audit into the primary intern-facing design.
The guide treats the missing experimental control plane as the product: system
snapshots, one-mutation candidates, immutable runs, synced result cells,
strict pairing, gate evaluation, and promotion reports. It also records the
specific boundaries that keep the repository pragmatic.

The design does not add runtime code. API blocks and schemas are proposed
contracts for phased review. Package directories must not be created until
their phase begins and the preceding exit criterion passes.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Produce enough technical detail that a new intern
can implement the package without reconstructing earlier projects or confusing
designed-only reflector work with proven mechanisms.

**Inferred intent:** Turn lessons from the earlier struggle into enforceable
defaults while resisting a broad optimizer platform.

### What I did

- Wrote the evidence classification and mechanism matrix.
- Compared `ragopt` with Ragkit, RAG-TTC, GEC-RAG, and researchctl ownership.
- Defined the v1 package layout and product dependency direction.
- Defined proposed asset, snapshot, candidate, suite, outcome, cell, policy,
  run-directory, decision, and promotion-plan contracts.
- Added lifecycle, component, and state diagrams.
- Added validation, scheduling/resume, pairing, and gate pseudocode.
- Defined the artifact-only Glazed CLI boundary.
- Defined separate RAG-TTC and GEC-RAG integration sequences.
- Added failure, security, privacy, testing, review, and operating guidance.
- Recorded architecture decisions and rejected alternatives.
- Expanded implementation phases and kept `tasks.md` as the canonical ledger.

### Why

- The earlier work had strong principles but no reusable implementation path.
- A new intern needs to know which code can be ported, which code must remain
  product-owned, and which appealing features are explicitly deferred.
- Schemas and negative tests reduce the chance that implementation silently
  weakens identity or failure accounting.

### What worked

- The design maps every v1 mechanism to implemented RAG-TTC or researchctl
  evidence, or labels it as a proposed composition of those mechanisms.
- The artifact CLI plus in-process library boundary avoids a new process/plugin
  compatibility layer.
- The phased plan can produce useful checkpoints before any live LLM run.

### What didn't work

- N/A. No runtime behavior was attempted in this step.

### What I learned

- “Batteries included” should mean the custody and gate path is complete, not
  that the package owns every diagnosis, proposal, execution, and deployment
  concern.
- The product's native artifact should remain authoritative; a generic outcome
  only needs comparable fields and a digest-linked pointer.
- A standalone optimizer command cannot safely execute arbitrary products
  without inventing an unproved protocol. The library boundary is smaller.

### What was tricky to build

Snapshot dimensions must be generic enough for different products without
becoming untyped, undigested metadata. The design allows bounded string
dimensions but requires each integration to declare the keys it considers
mandatory. The first two integrations should determine the smallest final API.

### What warrants a second pair of eyes

- Choose one public execution API after the first RAG-TTC spike: explicit arms
  or one evaluator receiving candidate views. Do not ship both.
- Decide whether v0.1 requires native artifacts inside the run directory.
- Review cancellation and active-run recovery semantics before Phase 1 freezes
  the status schema.
- Confirm candidate bundles can copy complete parent assets without excessive
  size before adding content-addressed indirection.

### What should be done in the future

- Obtain Phase 0 scope acceptance.
- Implement only Phase 1 next and prove interruption recovery with a scripted
  fixture before creating candidate or comparison packages.

### Code review instructions

- Read Sections 1–4 first to verify evidence strength.
- Review Sections 7–9 for schema and algorithm invariants.
- Review every architecture decision against the non-goals.
- Cross-check implementation sequence with `tasks.md`.

### Technical details

```text
Primary guide:
ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/design-doc/01-ragopt-intern-guide-to-a-reusable-evidence-gated-optimization-harness.md

Proposed v1 packages:
pkg/runstore, pkg/candidate, pkg/eval, pkg/compare, pkg/gate

Standalone CLI responsibility:
validate + inspect + compare + report

Product responsibility:
execute real arms + score/judge + diagnose + apply reviewed promotion
```
