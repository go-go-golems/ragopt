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
      Note: |-
        Implemented immutable run and synced JSONL evidence audited in Step 2
        Proven lifecycle ported narrowly in Phase 1
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/rag-ttc/pkg/rag/tooleval/runner.go
      Note: Implemented small Arm and Outcome comparison boundary audited in Step 2
    - Path: abs:///tmp/gec-ragopt-phase5/application-profiles.yaml
      Note: Production analyst-rag tool surface underlying the Step 33 provider data boundary
    - Path: abs:///tmp/gec-ragopt-phase5/cmd/coinvault/cmds/knowledge.go
      Note: GEC failure-visible judge accounting and denominators in Step 24
    - Path: abs:///tmp/gec-ragopt-phase5/cmd/coinvault/cmds/knowledge_ragopt.go
      Note: Provider-free preflight, frozen runtime validation, budgets, and product adapter in Step 32
    - Path: abs:///tmp/gec-ragopt-phase5/cmd/coinvault/cmds/knowledge_ragopt_trace.go
      Note: Authoritative native answer/tool/evidence/provider trace matched by the locked contract
    - Path: abs:///tmp/gec-ragopt-phase5/configs/ragopt/source-role-routing-v1/candidate.yaml
      Note: Frozen unexecuted candidate preserved across the Step 35 blocked audit
    - Path: abs:///tmp/gec-ragopt-phase5/configs/ragopt/source-role-routing-v1/candidate/snapshot.yaml
      Note: Recomputed challenger snapshot identity for the one-mutation proof
    - Path: abs:///tmp/gec-ragopt-phase5/configs/ragopt/source-role-routing-v1/parent/snapshot.yaml
      Note: Recomputed incumbent snapshot identity for the executable proof
    - Path: abs:///tmp/gec-ragopt-phase5/configs/ragopt/source-role-routing-v1/shared/source-lock.yaml
      Note: Committed semantic source identities proven by the GEC preflight
    - Path: abs:///tmp/gec-ragopt-phase5/internal/coinvaulttools/catalog/catalog.go
      Note: Production SQL tool binding reviewed in the Step 34 disclosure audit
    - Path: abs:///tmp/gec-ragopt-phase5/internal/coinvaulttools/catalog/registry.go
      Note: SQL and knowledge tool registration whose outputs can enter provider context
    - Path: abs:///tmp/gec-ragopt-phase5/internal/knowledge/judge.go
      Note: GEC idempotent timeline projection and strict judge semantics in Step 24
    - Path: abs:///tmp/gec-ragopt-phase5/internal/knowledge/service.go
      Note: GEC product-owned pre-fusion and pre-reranker authorization repair in Step 24
    - Path: abs:///tmp/gec-ragopt-phase5/internal/sqltool/normalize.go
      Note: Deterministic row-limit normalization reviewed for the provider consent scope
    - Path: abs:///tmp/gec-ragopt-phase5/internal/sqltool/types.go
      Note: Default row, byte, timeout, and blocked-system-schema policy
    - Path: abs:///tmp/gec-ragopt-phase5/internal/sqltool/validate.go
      Note: SELECT-only and schema restrictions with no business-table allowlist
    - Path: repo://cmd/ragopt/commands/candidate/validate.go
      Note: Glazed candidate validation row command
    - Path: repo://cmd/ragopt/main.go
      Note: Root logging and narrow Glazed parser configuration
    - Path: repo://go.mod
      Note: Normalized ragopt module and dependency baseline
    - Path: repo://pkg/candidate/candidate.go
      Note: Exactly-one-mutation validation in commit d2329dd
    - Path: repo://pkg/candidate/candidate_test.go
      Note: Phase 2 rejection matrix and identity fixtures
    - Path: repo://pkg/candidate/snapshot.go
      Note: Strict snapshot and asset identity loading in commit d2329dd
    - Path: repo://pkg/eval/bind.go
      Note: Complete candidate and policy input binding in commit 0802670
    - Path: repo://pkg/eval/resume.go
      Note: Strict cell replay and truncated-tail recovery in commit 0802670
    - Path: repo://pkg/eval/runner.go
      Note: Deterministic paired execution, outcome custody, and run lifecycle in commit 0802670
    - Path: repo://pkg/eval/runner_test.go
      Note: Four-case interruption, resume, failure, and corruption proof
    - Path: repo://pkg/runstore/read.go
      Note: Strict read-only run and copied-input integrity validator in commit 9f1ccb4
    - Path: repo://pkg/runstore/run.go
      Note: |-
        Phase 1 active-run lifecycle and synced JSONL implementation in commit 9f1ccb4
        Explicit active config-identical Resume authority
    - Path: repo://pkg/runstore/run_test.go
      Note: Interruption recovery, terminal, path, identity, and drift proof suite in commit 9f1ccb4
    - Path: repo://ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/index.md
      Note: Ticket-level distinction between the RAG-TTC proof adapter and both chat products
    - Path: repo://ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/reference/06-first-rag-ttc-i5-feedback-proof-rejection-report.md
      Note: Authoritative first corrected proof and rejected-candidate evidence
    - Path: repo://ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/tasks.md
      Note: Phase 5 product boundary, conditional validation gate, and reproducible-rejection criterion
ExternalSources: []
Summary: Chronological record of the ragopt repository bootstrap, evidence audit, scope decisions, documentation, validation, and delivery.
LastUpdated: 2026-08-06T16:24:00-04:00
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

## Step 4: Validate, commit, and deliver the design bundle

This step completed the initial design ticket's validation and delivery. The
primary guide, task ledger, diary, index, and changelog pass docmgr validation.
The detailed design checkpoint was committed before delivery, and the exact
four-document bundle was dry-run and uploaded to reMarkable with the default
layout requested by the user.

**Commit (docs):** `ee652e8` — "docs(ticket): design reusable optimization harness"

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Validate the ticket as a durable intern handoff,
commit it in a reviewable state, and upload it without the editor layout.

**Inferred intent:** Make the design immediately reviewable on both Git and
reMarkable before implementation begins.

### What I did

- Counted the primary guide at 1,725 lines and approximately 8,500 words.
- Searched the ticket for unfilled template comments and TODO/TBD placeholders.
- Ran `docmgr doctor --ticket RAGOPT-001 --stale-after 30`.
- Ran `git diff --cached --check` before the design commit.
- Committed the completed guide, diary step, and changelog as `ee652e8`.
- Dry-ran a four-document reMarkable bundle containing index, guide, tasks, and
  diary.
- Confirmed the dry run reported `layout=default`.
- Uploaded `RAGOPT 001 Reusable Evidence Gated Optimization Harness.pdf` to
  `/ai/2026/08/06/RAGOPT-001`.

### Why

- Docmgr validation catches broken metadata and RelatedFiles before an intern
  depends on them.
- Committing before upload makes the delivered design correspond to a stable
  repository checkpoint.
- A dry run verifies bundle membership and layout without creating a remote
  document.

### What worked

- Docmgr reported all checks passed.
- The staged diff had no whitespace errors.
- The dry run included the intended four files and used the default layout.
- The actual uploader returned:
  `OK: uploaded RAGOPT 001 Reusable Evidence Gated Optimization Harness.pdf -> /ai/2026/08/06/RAGOPT-001`.

### What didn't work

- The post-upload `GOWORK=off go build ./...` recheck ran without elevated
  filesystem access and failed while updating Go's module stat cache, not while
  compiling project code:

  ```text
  go: writing stat cache: open /home/manuel/go/pkg/mod/cache/download/github.com/go-go-golems/ragopt/@v/v0.0.0-20260806140949-ee652e8e0a37.info477099895.tmp: read-only file system
  ```

  The same build had passed during scaffold validation. I reran it with the
  required sandbox permission and recorded the final result below.

### What I learned

- The default uploader behavior is explicit in dry-run output, so omitting
  `--layout` is both sufficient and auditable.
- The guide is detailed enough for implementation while the task ledger stays
  short enough to serve as the progress interface.

### What was tricky to build

The final delivery has two truths: the committed design checkpoint and this
post-upload diary evidence. This diary step necessarily follows the upload, so
it is committed as a small final documentation checkpoint rather than rewriting
the already delivered PDF.

### What warrants a second pair of eyes

- Review the six open Phase 0 questions before checking scope acceptance.
- Verify the first implementation ticket starts with `pkg/runstore` only.

### What should be done in the future

- Review and accept Phase 0 scope.
- Begin Phase 1 with the interruption-recovery fixture and run-store contract.

### Code review instructions

- Review commits `fc3a703`, `de893d9`, and `ee652e8` in order.
- Run `docmgr doctor --ticket RAGOPT-001 --stale-after 30`.
- Run `GOWORK=off go test ./...` and `GOWORK=off go build ./...`.

### Technical details

```text
Bundle: RAGOPT 001 Reusable Evidence Gated Optimization Harness.pdf
Layout: default
Remote directory: /ai/2026/08/06/RAGOPT-001
Upload result: OK
Final build recheck: PASS after granting module-cache write permission
```

## Step 5: Implement the immutable run store and strict reader

This step turns the first proven mechanism in the design into repository code.
`pkg/runstore` now creates self-identifying run directories, copies and hashes
immutable inputs, atomically publishes artifacts, syncs each JSONL result, and
enforces irreversible complete/failed transitions. The implementation remains
domain-neutral: it has no retrieval, judging, candidate, pairing, or gate
behavior.

The step also adds the missing half of the RAG-TTC pattern: a strict read-only
`Open` path. A run is not considered inspectable merely because its files can
be parsed. The reader verifies schema, run-directory identity, canonical config
digest, status state machine, copied-input membership, sizes, and digests. It
accepts a valid `active` run after interruption but deliberately grants no
mutation authority.

**Commit (code):** `9f1ccb4` — "feat(runstore): add durable run lifecycle"

### Prompt Context

**User prompt (verbatim):** "implement the ticket, keep a detailed diarya s you work, commit at appropriateintervals, make sure to use the format from the diary askill."

**Assistant interpretation:** Implement the accepted phases of RAGOPT-001 in
order, preserve a strict diary entry for every meaningful checkpoint, validate
behavior, and make focused code and documentation commits.

**Inferred user intent:** Replace another potentially design-only optimizer
effort with exercised library behavior and an auditable sequence of small,
reviewable milestones.

### What I did

- Re-read the Phase 1 ledger and the proven RAG-TTC `pkg/experiment` source.
- Defined public `ragopt-run/v1` manifest, status, input reference, host, and
  summary types in `pkg/runstore/types.go`.
- Implemented canonical JSON SHA-256 identity using `json.Decoder.UseNumber`.
- Implemented lexical path confinement plus rejection of existing symlink
  components.
- Implemented atomic write, file sync, rename, and directory sync publication.
- Implemented run creation with host, Go module, timestamps, semantic
  dimensions, config digest, and secure directory modes.
- Implemented append-open-sync-close JSONL durability.
- Implemented copied input custody with absolute original path, copied path,
  exact size, and SHA-256.
- Implemented complete and failed terminal transitions that reject all later
  mutation and preserve prior artifacts.
- Implemented strict read-only validation for active and terminal runs.
- Added focused lifecycle, interruption, identity, drift, path, symlink,
  terminal, failure-preservation, and schema tests.
- Ran generation so logcopter package metadata includes `pkg/runstore`.
- Added a separate on-disk contract and recovery reference for interns.

### Why

- Run custody is the lowest reusable layer and already had exercised source to
  port; candidate abstractions should not precede it.
- A strict reader makes corruption and semantic drift explicit before Phase 3
  tries to resume or compare result cells.
- Read-only active inspection preserves researchctl's proven refusal to guess
  that abandoned work is safe to overwrite.
- Secure local permissions, confinement, and synced records make durability a
  behavior with tests rather than a comment.

### What worked

- `go test ./pkg/runstore -count=1` passed with all focused tests.
- `go generate ./...`, `go fmt ./...`, `go test ./...`, and `go build ./...`
  passed together.
- `go test ./pkg/runstore -race -count=1` passed.
- The interruption fixture appended and synced two JSONL records, opened the
  still-active run through the strict reader, and recovered both records.
- A failed run retained its previously synced result cell and exposed the
  original failure message.
- Config hashes remained identical across different Go map insertion orders.

### What didn't work

- The first focused compilation found one concrete implementation error:

  ```text
  # github.com/go-go-golems/ragopt/pkg/runstore
  pkg/runstore/run.go:56:24: cannot use &prettyConfig (value of type *strings.Builder) as *bytes.Buffer value in argument to json.Indent
  FAIL github.com/go-go-golems/ragopt/pkg/runstore [build failed]
  ```

  `encoding/json.Indent` requires `*bytes.Buffer`. I replaced the
  `strings.Builder`, formatted again, and the test suite passed.

- A later `go mod tidy; go test ./...; go build ./...` attempt did not reach a
  Go failure. The external permission reviewer timed out without changing the
  worktree:

  ```text
  The automatic permission approval review did not finish before its deadline.
  Do not assume the action is unsafe based on the timeout alone. You may retry
  once, or ask the user for guidance or explicit approval.
  ```

  The repository-wide test and build had already passed immediately before
  this cleanup attempt. I left the template's existing `pkg/errors` module
  classification unchanged and deferred module normalization to the final
  validation phase rather than manufacture file changes around an approval
  failure.

### What I learned

- The proven RAG-TTC lifecycle ports cleanly when product-owned directories and
  tags are replaced by three generic roots and semantic dimensions.
- Canonical config hashing must avoid a decode through `float64`; `UseNumber`
  preserves large JSON number text while still normalizing key order and
  whitespace.
- An active status is useful evidence. Automatically converting it to failed
  or resumable would erase the distinction between process interruption and a
  recorded evaluator decision.
- Strict input validation should also reject untracked files. Validating only
  the manifest's listed entries would allow silent extra state into a bundle.

### What was tricky to build

There are two publication models with different guarantees. Ordinary JSON and
byte artifacts can use a sibling temporary file, `fsync`, rename, and parent
directory `fsync`; JSONL results must append in place. The JSONL API therefore
opens, appends one line, syncs, and closes before returning. Its successful
return—not the later run completion—is the recovery boundary.

Path confinement also has two layers. Lexical normalization rejects absolute
paths and `..`, while `Lstat` rejects any existing symlink component. This is
stronger and simpler than following symlinks and proving their destination,
but it still does not claim safety against a hostile process racing filesystem
changes between validation and write operations.

### What warrants a second pair of eyes

- Review whether Phase 3 should acquire resumable write ownership through a
  new explicit API or create a continuation run. Do not add a generic `OpenRW`
  method without full suite/snapshot/candidate/policy identity checks.
- Review the crash window in `CopyInput`: copied bytes are atomically published
  before the updated input manifest. The strict reader detects the resulting
  untracked file, but does not repair it.
- Review whether v1 should require an error string for every failed status;
  `Fail(nil)` currently records a valid failed state with no message.
- The symlink policy intentionally rejects even benign in-run symlinks. Keep it
  unless a demonstrated product integration requires a different contract.

### What should be done in the future

- Keep `ragopt-run/v1` unchanged through the candidate and evaluator phases;
  add fields only when a real consumer proves they belong in every run.
- Implement Phase 2 snapshots and candidates as copied, strict artifacts;
  do not add active-run recovery or generic plugins during that phase.
- Retry `go mod tidy` during the final full-repository validation checkpoint.

### Code review instructions

- Start at `pkg/runstore/types.go`, then follow `Create` in `run.go`,
  `CopyInput` in `input.go`, terminal transitions in `terminal.go`, and `Open`
  in `read.go`.
- Review `path.go` and `write.go` as one filesystem boundary.
- Read `pkg/runstore/run_test.go` next; `TestRunLifecycleAndReader` is the
  interruption and terminal-state proof.
- Run:

  ```bash
  go test ./pkg/runstore -count=1
  go test ./pkg/runstore -race -count=1
  go test ./...
  go build ./...
  ```

### Technical details

```text
Schema: ragopt-run/v1
States: active -> complete | failed
Required roots: inputs/, results/, native/
JSON publication: write temp -> fsync -> close -> rename -> fsync directory
JSONL publication: append one line -> fsync -> close
Config identity: sha256(canonical JSON with UseNumber)
Input identity: sha256(exact copied bytes) + byte count
Reader authority: validation only; no mutation or resume
Code checkpoint: 9f1ccb4
```

## Step 6: Implement strict snapshots, one-mutation candidates, and validation CLI

This step implements the artifact that Phase 1 runs will evaluate. A snapshot
now identifies a complete system through verified locked assets, mutable
assets, and semantic dimensions. A candidate now identifies two snapshots, a
human or product proposer, an explicit hypothesis, an expected target, known
regression risks, and diagnostic provenance. Loading independently proves that
exactly one mutable asset's bytes changed.

The same checkpoint replaces the empty binary with its first real Glazed
command: `ragopt candidate validate`. The command is intentionally thin. It
calls the strict library, emits one structured identity row after success, and
never edits the bundle. The root owns logging and exposes only domain flags,
the three Glazed v1.4 output flags, and standard logging flags.

**Commit (code):** `d2329dd` — "feat(candidate): validate one-mutation bundles"

### Prompt Context

**User prompt (verbatim):** (see Step 5)

**Assistant interpretation:** Continue the ticket phase by phase, implementing
the accepted snapshot and candidate contracts with executable rejection tests
and a current-convention Glazed command, then checkpoint and document it.

**Inferred user intent:** Ensure a candidate cannot be presented as a valid
self-optimization step unless the repository can prove its bytes, provenance,
locked inputs, and isolated mutation before expensive evaluation begins.

### What I did

- Defined `ragopt-snapshot/v1` and `ragopt-candidate/v1` public types.
- Implemented canonical snapshot identity with asset lists sorted by logical
  name and deterministic JSON encoding.
- Implemented a semantic candidate digest over verified parent/child snapshot
  IDs, proposer, mutation declaration, hypothesis, risks, and evidence.
- Implemented strict single-document YAML decoding with unknown-field
  rejection.
- Implemented bundle-relative canonical path checks, symlink resolution, and
  post-resolution confinement to a regular file inside the bundle.
- Recomputed every asset size and SHA-256 from exact bytes.
- Enforced unique logical names across locked and mutable asset sets.
- Enforced stable system, dimensions, locked assets, mutable-name set, and
  media type across the parent and child.
- Compared actual bytes and required exactly one mutable asset change.
- Rejected a declaration whose `mutation.asset` disagrees with the independently
  observed change.
- Added tests for valid mutation, declaration order, two mutations, locked
  mutation, dimension mutation, missing files, digest/size drift, symlink
  escape, unknown fields, multiple YAML documents, declaration mismatch, and
  snapshot ID mismatch.
- Added `ragopt candidate validate` using `cmds.GlazeCommand`,
  `values.Values`, structured rows, Cobra `RunE`, and root logging.
- Added command tests for emitted rows, wrapped errors, required flag parsing,
  logging, and exact help-visible output flags.
- Ran `go mod tidy`, which correctly promoted Glazed, Cobra, `pkg/errors`, and
  YAML to direct dependencies and added the logging backend checksum.

### Why

- A candidate manifest cannot be trusted to report its own mutation; validation
  must compare verified parent and child bytes.
- Locked evaluators, suites, safety policy, corpora, indexes, and model
  identities must be outside candidate control or measured improvement is
  ambiguous.
- Complete replacement assets are reviewable and deterministic; a generic
  patch applier would add an unproved mutation surface.
- Strict decode and confinement fail early, before an evaluator spends model
  calls or records misleading outcomes.
- A thin artifact CLI exercises the library without inventing a subprocess or
  plugin protocol for product runtimes.

### What worked

- `go test ./pkg/candidate -count=1` passed all candidate invariants.
- Direct Glazed command tests loaded a real temporary bundle and emitted the
  expected identity row.
- `go generate ./...`, `go fmt ./...`, `go test ./...`, and `go build ./...`
  passed.
- `go test ./pkg/candidate ./cmd/ragopt/... -race -count=1` passed.
- Generated help contains only `--bundle`, `--manifest`, `--format`,
  `--output-fields`, and `--max-output-rows` locally, plus root logging flags.
- Reversing asset declaration order preserves parent snapshot, child snapshot,
  and candidate identities.

### What didn't work

- The first command compile used a removed Glazed helper from a nearby project
  on a different Glazed version and exposed a previously unused logging
  checksum:

  ```text
  ../../../../go/pkg/mod/github.com/go-go-golems/glazed@v1.4.1/pkg/cmds/logging/init.go:16:2: missing go.sum entry for module providing package gopkg.in/natefinch/lumberjack.v2
  cmd/ragopt/commands/candidate/validate.go:35:33: undefined: settings.NewGlazedSchema
  ```

  I inspected the pinned v1.4.1 source, replaced the removed helper with the
  current builder-owned structured-output contract, and ran `go mod tidy`.

- The initial root command explicitly installed
  `cli.CobraCommandDefaultMiddlewares`. Generated help then contained legacy
  command settings and the focused test failed:

  ```text
  unexpected automatic flag --config-file
  unexpected automatic flag --print-parsed-fields
  unexpected automatic flag --print-schema
  unexpected automatic flag --print-yaml
  ```

  Removing the explicit middleware and manually added output section was
  necessary but not sufficient because Glazed adds its command-settings
  section by default.

- My first attempt to use `cli.WithSkipCommandSettingsSection()` still failed
  with the same four flags. Source inspection showed that the following
  `WithParserConfig(...)` option replaces the entire parser config and erased
  the previously set skip bit. I stopped after the prescribed two consecutive
  fix attempts and reported `I think I'm stuck, let's TOUCH GRASS` rather than
  continuing to hack at flags.

- On the resumed pass, I placed `SkipCommandSettingsSection: true` inside the
  same `CobraParserConfig` as `ShortHelpSections`. Focused tests and generated
  help then passed without custom middleware.

### What I learned

- The current Glazed builder injects the structured-output section for every
  `GlazeCommand`; command constructors should not add it manually.
- `WithParserConfig` replaces, rather than merges, parser configuration. Related
  parser settings must be supplied atomically or later options can erase
  earlier builder options.
- A candidate bundle can reuse parent paths for unchanged assets. Only the
  changed asset needs a complete candidate replacement file.
- Path is intentionally part of snapshot identity. Moving an unchanged file,
  even with identical bytes, is visible metadata drift.
- The bundle validator is read-only. Evaluated immutability is achieved only
  when Phase 3 binds these verified identities and files into an immutable run;
  validation alone does not make a caller-owned directory unwritable.

### What was tricky to build

Snapshot identity and mutation detection need different normalization rules.
Identity includes paths, media types, sizes, and digests after sorting asset
records. Mutation detection compares actual bytes by logical name. If bytes are
unchanged, the full asset reference must also remain unchanged; if bytes differ,
the media type must remain stable. This prevents a metadata-only candidate from
appearing as a content mutation and prevents two mutation surfaces from hiding
under one asset name.

The Glazed issue was not a domain bug. Three sources of flags were initially
conflated: the command's domain schema, builder-injected structured output, and
parser-injected command settings. Reading the exact pinned `cobra-parser.go`
and option implementation revealed the overwrite behavior. The correct fix was
one parser config, not flag deletion after command construction.

### What warrants a second pair of eyes

- Review whether including asset paths in snapshot identity is desirable for
  both first product integrations. The current design and tests treat path drift
  as semantic drift.
- Review the decision to allow symlinks that resolve to a regular file inside
  the bundle while rejecting every symlink that resolves outside it.
- Review whether `Proposer.Kind` should remain an arbitrary bounded string or
  become a v1 enum after real model-proposed candidates exist.
- Review the `Candidate` value's exported maps and slices. Callers are expected
  to treat it as immutable; Phase 3 should project only the fields arms need.
- The task requiring the evaluated candidate transition remains tied to Phase
  3 run custody. Do not claim source-directory immutability before that binding
  exists.

### What should be done in the future

- Phase 3 must copy or digest-link the candidate manifest, both snapshots, and
  referenced assets into the run before executing either arm.
- The paired runner must accept a small immutable `CandidateView`, not the
  writable source bundle or a generic product adapter.
- Keep the candidate CLI read-only; no `fix`, `seal`, `apply`, or production
  mutation command should be added without separate evidence and design.

### Code review instructions

- Start at `pkg/candidate/types.go`, then read `snapshot.go`, `candidate.go`,
  `digest.go`, `path.go`, and `yaml.go` in that order.
- Review `pkg/candidate/candidate_test.go` as the executable rejection matrix.
- Review `cmd/ragopt/commands/candidate/validate.go` for the row boundary and
  `cmd/ragopt/main.go` for parser/logging ownership.
- Confirm the help contract with:

  ```bash
  go run ./cmd/ragopt candidate validate --help
  ```

- Validate behavior with:

  ```bash
  go test ./pkg/candidate ./cmd/ragopt/... -count=1
  go test ./pkg/candidate ./cmd/ragopt/... -race -count=1
  go test ./...
  go build ./...
  ```

### Technical details

```text
Snapshot schema: ragopt-snapshot/v1
Candidate schema: ragopt-candidate/v1
Asset digest: sha256(exact bytes)
Snapshot digest: sha256(canonical JSON; assets sorted by logical name)
Candidate digest: sha256(parent ID + child ID + proposal semantics)
Candidate mutation: exactly one mutable asset with different verified bytes
Locked invariants: system + dimensions + locked refs/bytes + mutable name set
CLI output: one row on success, no false validity row on error
Code checkpoint: d2329dd
```

## Step 7: Implement immutable input binding and resumable paired evaluation

This step implements the execution spine that turns a validated candidate into
durable evidence. `pkg/eval` now strictly loads ordered suites, binds every
suite/policy/candidate/snapshot/asset input into a run, passes incumbent and
challenger arms only copied candidate views, schedules paired cells
deterministically, verifies native artifacts, and appends one synced cell after
each completed arm invocation.

Resume is explicit rather than inferred. `runstore.Resume` reopens only a valid
active run whose canonical config matches the caller's exact request. The
evaluation layer then verifies the complete copied-input set and every existing
cell/native artifact before executing only missing schedule keys. Cancellation
leaves the run active; ordinary product errors become retained failed cells;
custody corruption fails the run.

**Commit (code):** `0802670` — "feat(eval): add resumable paired runner"

### Prompt Context

**User prompt (verbatim):** (see Step 5)

**Assistant interpretation:** Implement the next ticket phase using the
already chosen explicit two-arm API, prove interruption and exact resume with a
scripted product, keep product-native artifacts authoritative, and avoid a
workflow engine or concurrency layer.

**Inferred user intent:** Make the previously “fake” optimization loop
operationally honest: every incumbent/challenger result must be attributable,
durable, resumable, paired by construction, and impossible to improve by
silently losing failures.

### What I did

- Extended `candidate.Candidate` with its validated manifest locator so the
  runner can copy the exact candidate declaration into a run.
- Added explicit `runstore.Resume(ctx, directory, expectedConfig)` for active,
  config-identical, single-writer recovery.
- Added defensive `Run.Inputs()` access for input-set verification.
- Defined `ragopt-suite/v1`, `ragopt-cell/v1`, case, outcome, failure, artifact,
  candidate view, arm request, run request, and run result types.
- Implemented strict suite JSON loading, opaque input canonicalization, unique
  IDs/groups, sorted groups, declared case-order identity, and source reload.
- Chose exactly two explicit public arms: `Incumbent` and `Challenger`.
- Captured stable arm names once during request preparation.
- Normalized the run root to an absolute path.
- Reloaded suite, candidate, and policy sources immediately before run creation
  or resume and recomputed their identities.
- Built the complete expected role-to-byte-digest map for suite, policy,
  candidate manifest, both snapshots, and every parent/child asset.
- Copied all declared inputs with `runstore.CopyInput` and verified the exact
  count, roles, and digests before execution and resume.
- Built arm-facing `CandidateView` values exclusively from copied run inputs.
- Implemented deterministic case → repeat → incumbent → challenger order.
- Created a unique native directory per cell and required the returned artifact
  to be a regular file inside that directory.
- Recomputed native artifact SHA-256 and size, validated outcome state and
  counters, and rejected non-finite metrics.
- Converted ordinary arm errors into synthetic native `arm-error.json`
  artifacts and durable failed cells, then continued the schedule.
- Left cancellation/deadline-interrupted coordinates unrecorded and the run
  active for explicit resume.
- Implemented strict cell replay, full semantic cell keys, duplicate and
  unexpected key rejection, timestamp checks, native artifact revalidation,
  and missing-cell scheduling.
- Implemented active-run truncated-tail recovery by truncating to the last
  committed newline and syncing both file and directory.
- Added a four-case, two-repeat scripted product fixture and corruption tests.

### Why

- Explicit parent/challenger arms are the smallest API already supported by
  RAG-TTC evidence; shipping a second evaluator-shaped API would create a
  compatibility surface before any consumer needs it.
- Copying all candidate inputs is the actual `validated -> running`
  immutability transition. A validated source directory alone remains writable.
- A complete input role/digest map detects removal as well as drift; checking
  only the assets an arm happens to request would leave hidden run state.
- Immediate append plus `fsync` makes a successful cell return the resume
  boundary, not a later summary write.
- Failed arm invocations are data. Custody errors are different because they
  make later data untrustworthy.
- Sequential interleaving minimizes temporal drift and is sufficient for the
  first two product proof cycles.

### What worked

- `TestInterruptedRunResumesToUninterruptedCanonicalResult` committed exactly
  three cells before cancellation, discarded one uncommitted tail, and invoked
  arms exactly thirteen times on resume to reach sixteen cells.
- The resumed and uninterrupted runs produced equal canonical cells after
  removing run ID, timestamps, and measured duration.
- Arm observations proved every prompt path was below the immutable run
  directory, not the source candidate bundle.
- Editing the original candidate asset after completion did not invalidate or
  alter the run's copied inputs.
- An ordinary scripted provider error produced one failed cell and the other
  seven cells still completed in the one-repeat fixture.
- Identity mismatch refused resume while preserving active state.
- Missing copied input, malformed middle line, duplicate cell, native artifact
  escape, and NaN metric all failed custody explicitly.
- `go generate ./...`, `go fmt ./...`, `go test ./...`, and `go build ./...`
  passed.
- `go test ./pkg/eval ./pkg/runstore ./pkg/candidate -race -count=1` passed.

### What didn't work

- While writing the first resume helper, I drafted an unnecessarily clever and
  incorrect repeat formatter. Review caught it before compilation; it was
  replaced with the direct contract:

  ```go
  func formatRepeat(repeat int) string {
      return fmt.Sprintf("%04d", repeat)
  }
  ```

- A bulk patch that attempted to add input-set verification and captured arm
  names failed because its context no longer matched the formatted runner:

  ```text
  apply_patch verification failed: Failed to find expected lines in
  /home/manuel/code/wesen/go-go-golems/ragopt/pkg/eval/runner.go:
      if err := run.Complete(ctx, runstore.Summary{
  ```

  No partial change was applied. I re-read the current files and applied small
  patches to the config, preparation, binding, schedule, and resume boundaries.

- N/A for runtime validation: the focused Phase 3 suite passed on its first run
  after the scripted fixtures were formatted.

### What I learned

- “Resume” has three separate identities: common run config, complete bound
  input set, and completed cell/native-artifact set. Checking only one is not
  enough.
- A newline is part of the JSONL durability boundary. A final fragment without
  one can be discarded safely only for an explicitly active run; malformed
  committed middle lines must fail.
- Arm names must be captured once. Repeated calls to a stateful `Name()` method
  could otherwise change schedule and cell identity during a run.
- Run-root normalization matters because native artifact confinement compares
  resolved absolute paths.
- The common outcome should remain small. Product answer/transcript detail
  remains in the digest-linked native artifact.
- Exact sequential resume is simpler when the deterministic schedule is
  rebuilt and existing full keys are treated as a set.

### What was tricky to build

The difficult boundary was error classification. An ordinary arm error must
not abort the experiment or vanish. The runner creates an attributable failure
artifact and cell. Cancellation cannot use the same path because recording a
failure would prevent a legitimate retry. A malformed outcome or escaped
artifact cannot be converted to an arm failure because the runner cannot trust
its custody; it fails the run instead.

Input binding also needed two layers. Candidate validation proves source bytes
and semantic isolation. Evaluation copies those bytes into a run and stores a
complete role/digest map in the config. Resume then uses copied paths for arms
but still requires the caller's same source identities. This is intentionally
strict local recovery, not a bundle registry.

### What warrants a second pair of eyes

- `runstore.Resume` has no inter-process lock. Confirm every consumer command
  makes the single-writer operator requirement explicit.
- Review the decision to duplicate unchanged locked assets into separate parent
  and candidate input roles. It is literal and simple but uses extra local disk.
- Review whether policy byte digest should remain the cell identity once Phase
  4 introduces a normalized semantic policy digest; use one documented value,
  not both ambiguously.
- Review outcome consistency rules, especially completed but
  `contract_valid=false` without a `Failure`, which is intentionally retained.
- Review the five-second best-effort terminal failure context for local storage.
- The runner gives arms the absolute run directory as well as their assigned
  native directory. Confinement validates the returned authoritative artifact,
  but a malicious in-process arm is outside the v1 threat model.

### What should be done in the future

- Phase 4 must load the suite and cells strictly, join only exact parent/child
  coordinates, and reject any missing cell.
- Gate policy parsing should replace the current opaque policy interpretation
  without changing the byte custody already recorded here.
- The first RAG-TTC consumer must implement its arm in RAG-TTC and keep its
  native artifacts authoritative.
- Do not add concurrency until measured product runtimes justify the additional
  custody and scheduling complexity.

### Code review instructions

- Read `pkg/eval/types.go` and `suite.go` first.
- Follow new runs through `prepareRequest`, `bindInputs`, `buildSchedule`, and
  `executeCell` in `runner.go`.
- Follow resume through `runstore.Resume`, `verifyBoundInputs`, and
  `loadCompletedCells` in `resume.go`.
- Treat `pkg/eval/runner_test.go` as the executable lifecycle specification;
  start with the interruption/resume test.
- Run:

  ```bash
  go test ./pkg/eval ./pkg/runstore ./pkg/candidate -count=1
  go test ./pkg/eval ./pkg/runstore ./pkg/candidate -race -count=1
  go test ./...
  go build ./...
  ```

### Technical details

```text
Suite schema: ragopt-suite/v1
Cell schema: ragopt-cell/v1
Run config schema: ragopt-eval-run/v1
Schedule: case order -> repeat 0..N-1 -> incumbent -> challenger
Cell durability: append JSON line -> fsync -> close
Resume authority: active state + exact config + exact input set + valid cells
Truncated recovery: discard only bytes after last committed newline
Native custody: one assigned directory + regular file + recomputed SHA-256/size
Code checkpoint: 0802670
```

## Step 8: Implement strict paired comparison and lexicographic gates

This step started the read-only decision half of the harness. The first
boundary was a strict artifact loader: comparison must not trust a caller's
in-memory cells or silently repair an active run. `eval.LoadArtifactRun` now
opens the run through `runstore.Reader`, strictly decodes the evaluation
configuration, verifies the complete copied-input role/digest set, loads the
copied suite, and validates each committed cell plus its native artifact.

The evaluation configuration was extended with the already-proven mutation
identity needed by a promotion report: mutation declaration, changed logical
asset, and parent/child asset digests. These values are copied from the
validated candidate at run creation rather than reconstructed later.

The comparison package now joins cells by the exact tuple `(case_id,
repeat_index, arm)`. It rejects unknown and duplicate coordinates, constructs
only complete incumbent/challenger pairs, and retains absent arms as explicit
`MissingPair` records. It computes raw candidate-minus-incumbent deltas and
separately aggregates completion, contract validity, failures, abstentions,
costs, and quality metrics. No failed or missing cell receives a synthetic
quality score.

The gate package now strictly loads one `ragopt-gate-policy/v1` YAML document.
There are deliberately no product-independent quality defaults: the product
must name its target metric, threshold, floors, regression limits, and ordered
cost tie-breakers. Gate evaluation is lexicographic:

```text
identity and complete pairing
        -> hard completion/contract/failure/floor gates
        -> target improvement
        -> per-case and per-group regression limits
        -> informational ordered cost tie-breakers
```

Once a phase fails, later phases are not evaluated. This prevents a cheaper
candidate from compensating for incomplete evidence, broken contracts, or a
quality regression.

### Commands run

```bash
git status --short
rg -n "type Policy|func Evaluate|RunAPIVersion|type RunConfig" pkg
gofmt -w pkg/eval/runner.go pkg/eval/artifacts.go pkg/compare/types.go \
  pkg/compare/build.go pkg/gate/policy.go pkg/gate/evaluate.go
go test ./pkg/eval ./pkg/compare ./pkg/gate
```

The first focused test run found one compile error. After the correction, the
same focused command passed for all three packages.

### What worked

- The partially interrupted source patch was present as complete Go files and
  could be reconciled without discarding work.
- Existing evaluation tests remained green after publishing `RunConfig` and
  adding report identities.
- The strict loader reuses the runstore reader and evaluation outcome custody
  checks instead of inventing a second run format.
- Comparison preserves complete pairs, missing pairs, metric presence, raw
  values, deltas, and non-quality outcomes as separate facts.
- Policy parsing rejects unknown fields, multiple YAML documents, invalid
  identifiers, non-finite thresholds, absent failure thresholds, duplicate
  tie-breakers, and unsupported tie-breakers.

### What didn't work

- The first compile reported:

  ```text
  pkg/gate/evaluate.go:156:13: invalid operation: cannot call mean
  (variable of type float64): float64 is not a function
  ```

  A local `mean` accumulator shadowed the helper of the same name. It was
  renamed to `targetMean`; the next focused run passed.

- Review found a semantic aggregation bug before tests were written:
  `MetricAggregate.CompletePairs` was incremented only when both arms exposed
  that metric. This collapsed pair completeness into metric presence. It now
  takes the group's complete-pair count, while `PairsWithMetric` independently
  counts usable paired metric values.

### What I learned

- Artifact validity, pair completeness, and metric presence are three
  different denominators and must remain different fields.
- A raw policy byte digest is appropriate for proving which exact policy was
  evaluated; the normalized semantic digest is appropriate for identifying
  the parsed policy in a decision. Both need explicit names.
- Missing pairs should fail at the identity phase even if a policy author
  accidentally sets `require_all_cells: false`; absence of paired evidence is
  not a product-specific quality preference.

### Next actions

- Begin Phase 5 with one human-authored RAG-TTC text candidate; keep its Arm
  and native scoring artifacts in the product repository.
- Record integration friction before expanding any generic `ragopt` API.

### Phase 4 completion update

The comparison, policy, decision, report, and CLI layers were completed in
commit `7fc1a98`. Six decision goldens cover pass, hard-gate failure, target
failure, catastrophic regression, tie-break reporting, and incomplete
pairing. The incomplete fixture deletes a challenger cell and fails during
identity, demonstrating that missing evidence cannot improve the result.

The deterministic report contains identities, hypothesis, exact asset digests,
checks, group outcomes, metric denominators and deltas, missing pairs, and
per-pair failure/cost evidence. Its JSON plan is fixed to
`state=review_required` and `human_apply_required=true`; there is no apply API.

Final commands for this checkpoint:

```bash
go test ./...
go test ./pkg/eval ./pkg/compare ./pkg/gate ./pkg/report -race -count=1
go build ./...
go run ./cmd/ragopt compare --help
go run ./cmd/ragopt report --help
git diff --check
docmgr status --ticket RAGOPT-001
docmgr validate frontmatter --doc \
  ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/reference/05-paired-comparison-gate-and-promotion-report-v1-contract.md
```

All Go validation and both CLI smoke checks passed. Docmgr recognized one
design and five reference documents, and the new reference frontmatter passed.

One documentation command was attempted with an unsupported scope:

```text
docmgr validate --ticket RAGOPT-001
Error: unknown flag: --ticket
```

This did not change files. I used the supported per-document `validate
frontmatter --doc ...` form, which passed.

### Phase 4 review instructions

- Start with `pkg/eval/artifacts.go` and verify it stays read-only.
- In `pkg/compare/build.go`, inspect full cell-identity checks before the
  coordinate join, then verify `ExpectedPairs`, `CompletePairs`, and
  `PairsWithMetric` remain distinct.
- In `pkg/gate/evaluate.go`, verify each failing phase returns before cheaper
  tie-breakers can influence the decision.
- Confirm `pkg/report` contains no candidate or product mutation operation.
- Read `pkg/gate/testdata/*.golden` as the compact decision specification.
- Use `reference/05-paired-comparison-gate-and-promotion-report-v1-contract.md`
  for the complete intern-facing API and review guide.

## Step 9: Re-deliver the implementation bundle and scope the RAG-TTC proof cycle

After Phase 4, I re-uploaded the ticket because the earlier reMarkable PDF was
the design-only checkpoint. The new bundle contains eight documents: ticket
index, primary design, task ledger, diary, and all four implemented API
contracts. The required dry run reported `layout=default`; no editor layout
flag was used.

```text
DRY: layout=default
DRY: bundle name=RAGOPT 001 Reusable Evidence Gated Optimization Harness
DRY: remote-dir=/ai/2026/08/06/RAGOPT-001
...
OK: uploaded RAGOPT 001 Reusable Evidence Gated Optimization Harness.pdf
    -> /ai/2026/08/06/RAGOPT-001
```

The `--force` replacement was intentional and authorized by the user's
“reupload” request. It replaces the earlier same-named design bundle, so any
annotations on that earlier remote document are not recoverable through the
uploader.

### RAG-TTC integration reconnaissance

I then began Phase 5 read-only in
`/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/rag-ttc`. The
umbrella workspace is not itself a Git repository; `rag-ttc` is a Git worktree
whose actual Git directory lives under `/home/manuel/code/ttc/rag-ttc`.

The product repository already has unrelated user changes:

```text
 M pkg/app/chatserver/server.go
 M profiles.yaml
 M ttmp/.../RAG-TTC-TOOLLOOP-001.../tasks.md
?? pkg/app/chatserver/websocket_test.go
?? scripts/parc-corpus/
?? several ticket source/analysis directories
```

I did not modify, stage, or clean any of them.

The proven integration seam is smaller than the answer-quality batch runner:

- `pkg/rag/tooleval.Arm` is already per-query and preserves native artifacts;
- `cmd/rag-ttc/cmds/chat/tooleval/adapter.go` owns the application chat
  runtime, satisfying the requirement that product execution stays in the
  product repository;
- `answerquality.JudgeToolLoop` is the existing two-stage judge and produces
  per-cell faithfulness and answer-relevance evidence;
- the incumbent text asset is
  `tool-descriptions/tool-qa/search-v1.yaml`;
- the human-authored candidate is
  `tool-descriptions/tool-qa/search-combined-comparison-i5-v1.yaml`;
- their corresponding complete runtime configurations are
  `configs/tool-qa/production-product-fact-v1.yaml` and
  `configs/tool-qa/production-product-fact-i5-combined-comparison-v1.yaml`.

The candidate changes search guidance from separate named-product searches to
one combined comparison query followed by a specific gap-filling query. This
is exactly one human-authored text mutation and does not alter model, corpus,
index, answer schema, safety ceilings, dataset, or judge.

### Dependency publication boundary

RAG-TTC cannot yet commit a clean import of `github.com/go-go-golems/ragopt`.
The ragopt repository has a real GitHub origin, but local `main` is twelve
commits ahead of `origin/main`:

```text
## main...origin/main [ahead 12]
origin/main = 0a319399eaa8621580047fbdf6e7b950c11c1e9e
local Phase 4 docs = 8cf376e
```

A machine-specific absolute `replace` directive would make RAG-TTC's committed
module non-portable, while requiring an unpublished pseudo-version would break
CI. Pushing is an external state change not implied by local commit authority,
so I stopped before editing RAG-TTC and will request explicit publication
authority.

### Commands run

```bash
git status --short
find . -maxdepth 3 -type d -name .git -print
git -C rag-ttc status --short
rg --files rag-ttc
go env GOWORK
git status -sb
git ls-remote origin refs/heads/main
remarquee upload bundle <eight explicit Markdown paths> \
  --name "RAGOPT 001 Reusable Evidence Gated Optimization Harness" \
  --remote-dir /ai/2026/08/06/RAGOPT-001 --dry-run --non-interactive
remarquee upload bundle <same paths> \
  --name "RAGOPT 001 Reusable Evidence Gated Optimization Harness" \
  --remote-dir /ai/2026/08/06/RAGOPT-001 --force --non-interactive
```

### Next action after approval

Push ragopt `main`, add the resulting real module revision to RAG-TTC, and
implement the product-owned adapter by composing the existing chat runtime and
judge. Do not introduce a generic subprocess protocol. Start with deterministic
adapter tests and a cached/replayable fixture before spending provider budget.

## Step 10: Freeze the first RAG-TTC candidate, suites, and policy

Publication authority was still pending, but Phase 5 had an independent unit
of work: product-owned immutable inputs. I created those inputs in the RAG-TTC
worktree without importing ragopt or touching its existing application edits.

The selected candidate is the already human-authored I5 experiment:

```text
incumbent: tool-descriptions/tool-qa/search-v1.yaml
candidate: tool-descriptions/tool-qa/search-combined-comparison-i5-v1.yaml
```

Its hypothesis is that one precise first query containing all comparison
subjects and requested attributes can preserve faithfulness/relevance while
reducing redundant provider and search-tool calls. The candidate changes only
`search_description`; runtime ceilings and safety controls remain locked.

### Product asset layout

RAG-TTC commit `a6cfd93` adds:

```text
configs/ragopt/i5-combined-comparison-v1/
├── candidate.yaml
├── parent/snapshot.yaml
├── parent/assets/search-description.yaml
├── candidate/snapshot.yaml
├── candidate/assets/search-description.yaml
└── shared/
    ├── answer-schema.json
    ├── feedback-suite.json
    ├── gate-policy.yaml
    ├── judge-contract.yaml
    ├── orchestration-prompt.txt
    ├── runtime-contract.yaml
    └── validation-suite.json
```

The snapshots lock:

- answer model and embedding model;
- corpus and index-manifest digests;
- original evaluation dataset digest;
- exact 3-case feedback and 7-case validation suites;
- answer schema and orchestration prompt bytes;
- judge contract and judge source digest;
- native tool-loop runner and application adapter source digests;
- provider-profile digest;
- provider/tool ceilings, transcript custody settings, and gate policy.

Feedback and validation questions are disjoint and copied verbatim from
`datasets/ttc/evaluation.json`. Every case is grouped as `comparison`,
`product_comparison`, and its split name. The candidate configuration and
incumbent configuration were compared structurally after replacing only their
description path and non-behavioral configuration name; no other configuration
field differs.

### Validated identities

```text
candidate  sha256:75b70d2aa28e49d9c0d21bf6d2f10c2a15f32385c3e3c65513026d95cc2f9049
parent     sha256:982d160bc51d8bfdc0db956f845135555c42c424ee2d73f0c88dc2715b7753f6
child      sha256:f433c4d44a0b8acacf4b03ffcdf2d8b155040b7bf506e4ca5e43a589778c67d1
feedback   sha256:b009f9e913179007bf2da04cadb2e75200e22e1ab38689513015721ee76781b6
validation sha256:d0ae5384bb23d7da9ea54b34841f167892314b23359e9da599f85d9dca8e960a
policy semantic sha256:75c56ff46aab7bb2b281942514faf8a7de94fa1bb4c06e2e0df9548fb7ddb83f
policy bytes sha256:7d3fba806e3824b26ff817b103fd6415cfe66117edc66ad3050c6ab6476b4d8e
```

The product policy requires complete, completed, contract-valid cells, zero
failures, and candidate faithfulness at least 0.80. It targets non-regression
in answer relevance, limits per-case and mean relevance/faithfulness
regressions, then reports provider calls, tool calls, tokens, and duration as
ordered tie-breakers. These are explicit RAG-TTC thresholds, not ragopt
defaults.

### Tests and commands

```bash
go run ./cmd/ragopt candidate validate --bundle <RAG-TTC bundle> --format json
go run /tmp/ragopt-validate-ttc-assets.go
go test ./internal/ragoptassets -count=1
go test ./internal/ragoptassets ./pkg/rag/tooleval \
  ./cmd/rag-ttc/cmds/chat/tooleval -count=1
git diff --check
```

All final validation passed. `internal/ragoptassets/assets_test.go` makes
production-asset drift, config drift, suite overlap, edited questions, and
missing split groups fail locally even before the future adapter imports
ragopt.

### What didn't work

- The first multi-file patch created `configs/ragopt/README.md` but could not
  create `internal/ragoptassets/assets_test.go` because its parent directory
  did not exist:

  ```text
  Failed to create parent directories for .../internal/ragoptassets/assets_test.go
  ```

  I verified the partial result, created only the intended directory, and
  applied the test file separately.

- The first focused test correctly found two authoring assumptions:

  ```text
  incumbent and candidate runtime configs differ outside the search-description path
  json: unknown field "id"
  ```

  The runtime configs also have distinct descriptive `name` fields, which the
  adapter will replace with arm-local identity; the test now normalizes that
  explicitly. The strict source projection was missing declared top-level
  dataset fields; they were added rather than disabling unknown-field checks.

- Snapshot manifests intentionally began with all-zero IDs. Two strict
  validation passes reported the computed parent and child identities, after
  which the final third pass validated the complete candidate.

- The pre-commit whitespace check found one extra blank line at the end of
  `candidate.yaml`. It was removed and the staged check then passed.

### Repository hygiene

The RAG-TTC worktree was already dirty. Commit `a6cfd93` contains exactly the
14 new `configs/ragopt` and `internal/ragoptassets` files. Existing modified
chatserver/profile/ticket files and untracked ticket sources/scripts remained
unstaged after the commit. The current RAG-TTC branch is
`task/rag-ttc-tui-polish`, an important coordination fact for later adapter
work and publication.

## Step 11: Publish ragopt and implement the product-owned RAG-TTC arm

The user authorized the previously blocked publication step with `push`. I
pushed ragopt `main` from `0a31939` through `4d410c5`, then pinned RAG-TTC to
the resulting portable pseudo-version:

```text
github.com/go-go-golems/ragopt
v0.0.0-20260806162156-4d410c57e242
```

This removed the dependency blocker without a machine-local `replace`.

### Adapter boundary implemented

RAG-TTC commit `56432fd` adds the product-owned adapter at
`cmd/rag-ttc/cmds/chat/tooleval/ragopt.go`. The new
`rag-ttc tool-loop ragopt` Glazed command composes only established mechanisms:

- `candidate.LoadCandidate` validates the one-mutation bundle;
- `eval.LoadSuite`, `eval.Run`, and `eval.Resume` own immutable paired custody;
- the existing `chat.NewRuntime` executes the real tool-QA implementation;
- the existing session archive remains the full native answer/trace evidence;
- the existing `answerquality.JudgeToolLoop` produces faithfulness and answer
  relevance;
- a `rag-ttc-ragopt-native/v1` JSON artifact links the session, runtime
  projection, judge score, and judge provenance;
- the common ragopt outcome contains comparison fields and a run-relative
  native-artifact reference.

No subprocess protocol, plugin interface, generic runtime adapter, proposal
model, worker pool, or deployment mutation was introduced. Execution remains
sequential. Every cell gets a fresh chat runtime rooted inside its assigned
native directory, preventing evidence from another cell from becoming its
authoritative artifact.

The interface exposes only the two locked proof splits:

```text
feedback   -> 3 cases x 1 repeat x 2 arms = 6 cells
validation -> 7 cases x 2 repeats x 2 arms = 28 cells
```

Per-cell budgets are visible but fail closed if changed from one embedding
call, four answer-generation calls, and two judge calls. The answer profile
must resolve as `ttc-live-luna-low`/`gpt-5.6-luna`; the established named judge
profiles are resolved separately.

### Preventing a fake implementation

The adapter does not accept arbitrary incumbent and challenger configuration
paths. ragopt passes each arm an immutable copied `CandidateView`. The adapter
materializes the locked orchestration prompt and answer schema plus that
view's exact parent or child `search_description`, then loads it through
RAG-TTC's strict `toolconfig` loader. Tests prove both views bind their declared
bytes.

Before provider spending, `validateI5Environment` recalculates evaluation,
provider-profile, index-manifest, corpus, answer/embedding, judge, native
runner, native adapter, and ragopt adapter identities. It also checks the suite
split, evaluator contract, and tool-safety identity.

The inaccurate pre-adapter Git revision dimension was replaced by a semantic
evaluator version plus exact source digests. The new adapter digest is:

```text
sha256:d4746119aa919b780563f3a161045e0f908c4f5141f4bc9ec34d6ba00b8fe2b8
```

Strict candidate validation now reports:

```text
candidate sha256:0c460d0f5d41bfc991fd1907a22bcb2407cac2f33225eedd25a04aa6b2539f1c
parent    sha256:aff0f3dcb658626ca00cf95384c80cc12f7ac83e9b4cb735b7c2f104584ea8c9
child     sha256:e3f232d3f6f0a98c11f08d03b0fe0150e6534c538e640b4702941edbfa4f3485
```

The independently computed mutation remains only `search_description`.

### Failures and corrections

The first formatting attempt failed because RAG-TTC is outside this session's
default writable root and `gofmt` creates a sibling temporary file:

```text
open cmd/rag-ttc/cmds/chat/tooleval/ragopt.go.<temporary>: read-only file system
```

I reran the narrowly scoped formatter with explicit write approval.

The first asset-binding test incorrectly treated the loader's logical search
path as absolute. Both arm subtests failed with:

```text
open search-description.yaml: no such file or directory
```

The assertion now resolves the logical path against the temporary repository
root, matching `toolconfig` semantics; both arms pass.

The first manual validation command guessed API-style flags and failed with
`Error: unknown flag: --bundle-root`. I read command help and used the actual
`--bundle` and `--format` flags; validation then passed.

During final review, HEAD advanced through four unrelated chat-hardening
commits in the shared worktree. One included the already-pinned module line. I
re-read history and the diff, then staged exactly five adapter/snapshot files.
The remaining profile, ticket, corpus, and chat work stayed unstaged.

### Verification performed

```bash
go test ./cmd/rag-ttc/cmds/chat/tooleval ./cmd/rag-ttc \
  ./internal/ragoptassets -count=1
go test ./... -count=1
go run ./cmd/rag-ttc tool-loop ragopt --help
go run github.com/go-go-golems/ragopt/cmd/ragopt candidate validate \
  --bundle configs/ragopt/i5-combined-comparison-v1 \
  --manifest candidate.yaml --format json
git diff --check
```

All focused tests and the complete RAG-TTC suite passed. A live feedback run
has not yet been claimed; it is the next task and consumes real provider calls.

## Step 12: Live feedback-run approval boundary

I inspected the long command help to confirm the Glazed profile and structured
output flags, then attempted to start the locked six-cell feedback run in tmux:

```bash
go run ./cmd/rag-ttc tool-loop ragopt \
  --split feedback \
  --profile ttc-live-luna-low \
  --format json \
  --log-level info
```

The execution gate rejected the start because this command sends the three
evaluation questions and retrieved TTC evidence to configured external answer
and judge providers and may incur provider spend. The user's implementation
request did not count as explicit approval for this specific external payload
and spend. No tmux session, run directory, provider call, or partial artifact
was created. The next action is to obtain explicit approval, then run the exact
command above and capture/resume its immutable run rather than changing scope.

I then re-rendered the eight-document ticket bundle. The dry run explicitly
reported `layout=default`, satisfying the user's instruction not to use the
editor layout. The forced replacement completed successfully:

```text
OK: uploaded RAGOPT 001 Reusable Evidence Gated Optimization Harness.pdf
    -> /ai/2026/08/06/RAGOPT-001
```

## Step 13: First live feedback run exposes an adapter budget bias

After explicit approval for external provider use, I started the locked
six-cell feedback command in tmux. Run
`20260806T174701.722301319Z-ttc-i5-feedback-5cc789838eea` completed 6/6 cells
with zero arm failures and six digest-linked native artifacts.

The results were not valid candidate evidence. Five cells became
contract-invalid abstentions because the adapter permitted only one embedding
call while the tool loop permits as many as three searches before its reserved
final answer call. Native failures said:

```text
resource budget exceeded: requested 1, remaining 0, limit 1
```

This structurally favored the combined-query candidate. The adapter also
incorrectly added two judge calls to every product `provider_calls` value even
when the judge skipped invalid cells. I retained the run as rejected
operational evidence and did not compare or promote it.

RAG-TTC commit `90485d8` made one bounded correction: lock per-cell budgets to
3/4/2 embedding/answer/judge calls, record them in the runtime contract,
exclude judge overhead from product cost, and test that separation. The
corrected one-mutation identities are:

```text
candidate sha256:0d598d871694580d7e36af0238a3431cf02690249f6484ea4380b00179507106
parent    sha256:4cc350f874416dd38e4288a090b95952a5d3d50de0b2ff8ff1fa0b26ee5b876a
child     sha256:836466c1b61320034b9e88cd2ba5e1656555082f1196d6e500165653bc8e8eaa
```

Focused tests passed. Concurrent unrelated logcopter generation temporarily
prevented rebuilding the shared root command; I did not touch those files.

## Step 14: Isolated corrected feedback proof and formal rejection

At the user's direction, I continued in detached worktree
`/tmp/rag-ttc-ragopt-proof` at commit `90485d8`. I copied the exact locked
`profiles.yaml` bytes into it. The first test failed because its ignored cache
did not contain the frozen index:

```text
digest locked I5 index manifest: ... manifest.json: no such file or directory
```

I linked only the established index into the detached worktree's ignored
`.cache` path. Focused adapter, asset, and root-command tests then passed.

Corrected run `20260806T180004.824520651Z-ttc-i5-feedback-a1781d11159f`
reached terminal `complete` in 100 seconds with 6/6 cells and one arm failure.
All expected pairs exist and every committed native artifact has validated
SHA-256 and size.

Cell outcomes:

- `ttc-y-005` incumbent: arm failure at the four-iteration ceiling;
- `ttc-y-005` challenger: valid and judged; faithfulness 1.0, relevance 1.0,
  four provider calls, and three searches;
- `ttc-y-007` incumbent: invalid abstention; four calls, three searches;
- `ttc-y-007` challenger: invalid abstention; three calls, two searches;
- `ttc-expand-056` incumbent: invalid abstention; three calls, two searches;
- `ttc-expand-056` challenger: invalid abstention; four calls, three searches.

`ragopt compare` returned `fail`. Identity, pairing, candidate completion, and
candidate failure-rate checks passed. Hard gates failed because only one of
three candidate cells was contract-valid and faithfulness had coverage for
only one pair. `ragopt report` produced a non-applying `review_required`,
`human_apply_required` plan saying not to promote. Validation was skipped
because the candidate failed feedback hard gates.

A diagnostic loop initially used zsh's special variable `path`, replacing the
command search path and causing repeated `command not found: jq`. Renaming it
to `artifact_file` fixed the read-only diagnostic; no artifact changed.

The final nine-document ticket bundle dry run reported `layout=default` and
included the new proof rejection report. The forced replacement then uploaded
successfully to `/ai/2026/08/06/RAGOPT-001`. No editor layout was used.

## Step 15: Re-establish the RAGOPT boundary and prepare the reproducibility proof

This step resumed `RAGOPT-001` after a conversation had conflated three
different products: the reusable optimizer harness, the canonical full-page
Admin Chat cutover, and the customer-facing Garden Assistant. I tightened the
Phase 5 ledger so the current proof target and its success criterion are
unambiguous before spending on a second provider-backed run.

The next proof remains the same: repeat the rejected I5 feedback experiment
from a fresh run root and require the same semantic identities, complete
paired custody, canonical deltas, and fail-closed decision. A rejected
candidate can prove the harness. It does not need to become promotable, and
the feedback failure means the validation split must remain unrun.

**Commit (documentation):** `15f5757` — "docs(ticket): clarify product proof boundary"

### Prompt Context

**User prompt (verbatim):** "work on RAGOPT-001. Keep a detailed diary as you work, commit at appropriate intervals, make sure to use the diary format form the skill"

**Assistant interpretation:** Resume the active RAGOPT implementation ledger,
use the strict diary structure for every step, preserve focused commits, and
continue until the real product-integration requirements are evidenced.

**Inferred user intent:** Keep the optimization work separate from the
colleague's sessionstream/Admin Chat work and make trustworthy, reviewable
progress toward reusable RAG-TTC and GEC proof cycles.

### What I did

- Read the diary and docmgr skills completely, including the diary reference.
- Re-read the current task ledger, ticket index, latest diary steps, rejection
  report, Git history, and docmgr task status.
- Added an explicit Phase 5 scope boundary: the current adapter targets
  `rag-ttc tool-loop ragopt`, not Admin Chat, sessionstream, or Garden Assistant.
- Defined a product proof as reproducible custody and decision behavior; a
  stable rejection is valid proof.
- Recast the validation item as the already-satisfied spending gate: I5 failed
  feedback, so the 28-cell validation run correctly did not happen.
- Updated the two missing related-file notes in the rejection report through
  `docmgr doc relate`.
- Attempted to launch the second six-cell feedback proof in tmux from the
  isolated `/tmp/rag-ttc-ragopt-proof` worktree.

### Why

- The ticket must not imply ownership of the colleague's canonical chat work.
- Requiring a candidate to win would confuse optimization-harness correctness
  with stochastic product quality and could encourage weakening gates.
- Validation spending must remain conditional on feedback success.
- The second live run exports questions/evidence to configured providers and
  incurs new spend, so it requires explicit approval for that distinct run.

### What worked

- `docmgr task list --ticket RAGOPT-001` recognizes task 64 as complete and
  leaves the fresh-root reproduction as the next task.
- The task and index text now name the exact product boundary and success
  condition.
- Both stale related-file entries now have explanatory notes.
- The isolated proof worktree is still at corrected product commit `90485d8`;
  its frozen index link and first corrected run remain intact.

### What didn't work

- The first two sandboxed metadata updates failed exactly as follows:

  ```text
  Error: failed to write document: open /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/reference/.docmgr-2850251983: read-only file system
  Error: failed to write document: open /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/reference/.docmgr-3388705735: read-only file system
  ```

  Re-running only those `docmgr doc relate` writes with repository-scoped
  approval succeeded.

- The provider-backed tmux launch was rejected before process creation:

  ```text
  This starts a new provider-backed reproducibility run with additional external data export and spend; prior approval covered the corrected run, not this separate second run.
  ```

  No tmux session, run root, provider call, or partial result was created. I did
  not retry through another command or execution path.

### What I learned

- The current RAG-TTC adapter proves the product evaluation runtime, not either
  chat UI. A serving-runtime optimizer would be a separate product-owned adapter.
- “Proof passes” must describe infrastructure evidence, while “candidate
  passes” describes a product gate decision. They are independent claims.
- The conditional validation task is a spending invariant, not an outstanding
  command to run for every rejected feedback candidate.

### What was tricky to build

The task ledger previously mixed an unconditional checkbox with conditional
provider spending. Leaving the validation checkbox open suggested unfinished
work, while checking text that said “run validation” would falsely claim the
28 cells ran. The correction makes the task the enforceable behavior—stop after
failed feedback—while retaining the actual fresh-root reproduction separately.

### What warrants a second pair of eyes

- Confirm that `tool-loop ragopt` remains the intended first consumer seam and
  that no canonical Admin Chat or Garden Assistant runtime identity is implied.
- Confirm that a repeated, evidence-identical rejection satisfies the RAG-TTC
  integration proof even if individual stochastic answer text differs.
- Review the eventual comparison method carefully: semantic identities and
  decision class should be exact; stochastic metrics need canonical paired
  comparison rather than byte-identical native transcripts.

### What should be done in the future

- Obtain explicit approval for the second six-cell provider-backed feedback run.
- After it completes, compare semantic identities, cell coordinates, native
  artifact custody, metric deltas, and the final gate decision with the first
  corrected run.
- Only after that proof, begin the GEC-RAG product-owned adapter and its P0
  authorization/judge-accounting corrections.

### Code review instructions

- Start with `tasks.md` Phase 5 and the ticket `index.md` status section.
- Check the report frontmatter with
  `docmgr doctor --ticket RAGOPT-001 --stale-after 30`.
- Confirm `git diff --check` and `docmgr task list --ticket RAGOPT-001` pass.
- Do not inspect or modify `rag-ttc/pkg/app/chatserver` for this step.

### Technical details

```text
RAGOPT repository: /home/manuel/code/wesen/go-go-golems/ragopt
Product proof worktree: /tmp/rag-ttc-ragopt-proof
Product commit: 90485d8539515173e88d2a9703f3a7ad4aed74bc
First corrected run: 20260806T180004.824520651Z-ttc-i5-feedback-a1781d11159f
Next ledger task: 67 — fresh-root reproduction
External work needed: 6 cells (3 feedback cases x 2 arms)
Validation cells intentionally not authorized: 28
```

## Step 16: Design one shared production-refresh control plane for three products

This step turned the preliminary GEC scheduling note into an intern-facing
cross-product design for CoinVault/GEC, the TTC Garden Assistant, and the
RAG-TTC Admin Chat. The central correction is that `ragopt` should provide the
reusable refresh protocol and evidence lifecycle, while stopping short of
becoming a scheduler, queue, generic workflow engine, product indexer, chat
runtime, or deployment controller.

The design now distinguishes content refresh from one-mutation optimization,
defines proposed Go interfaces and event schemas, maps the existing code, and
provides a phased task plan. It also records the important TTC topology: one
approved WordPress/WooCommerce-derived bundle may feed both Garden and Admin,
but each serving runtime requires separate native evaluation and gates.

### Prompt Context

**User prompt (verbatim):** "how would ragopt allow us to run a production indexer that would run on changes / overnight nd have things like resumability and tracking what is oging on and all that? We'll be running in AWS if that matters, but we could also either build / use our own thing (i don't know if flow can be used for that), or using something like riverqueue or so"

**Additional user prompt (verbatim):** "use a separate design doc to write this all up and upload it to remarkable so i can review and understand the design and the implementations you suggest.

Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable."

**Additional user prompt (verbatim):** "we will reuse ragopt for both rag-ttc (admin and garden assistant) and coinvault, as a base package, and the needs are very similar, if not overlapping, btw"

**Additional user prompt (verbatim):** "ok, yes, that's the document we want to continue working on."

**Assistant interpretation:** Continue the existing production-index design as
a detailed, pragmatic, shared `ragopt` architecture; explain how local replay,
durable orchestration, product adapters, evaluation, activation, and rollback
fit together for all three consumers; record precise implementation phases and
deliver the standalone document to reMarkable without editor layout.

**Inferred user intent:** Prevent another product from assembling a
self-optimization or refresh loop that looks complete but lacks durable
identity, progress, replay, evaluation custody, and safe promotion, while
reusing the mechanisms already proven in RAG-TTC and GEC instead of building a
large new platform.

**Commit (documentation):** `8eaa35f` — "docs(design): define shared production refresh control plane"

### What I did

- Read the diary, docmgr, and reMarkable upload skills and the full diary
  reference before modifying the ticket.
- Re-read the production-refresh draft and canonical task ledger.
- Audited the current package/file locations in `ragopt`, GEC,
  `rag-ttc/pkg/flow`, `pkg/rag/indexbundle`, `pkg/rag/tooleval`, the Garden
  `backend/internal/ragsearch`, and Admin `pkg/app/chatserver` runtime boundary.
- Rewrote the executive boundary so `ragopt` owns a fixed shared refresh
  coordinator, semantic identities, events/progress, artifact links, resume,
  evaluation handoff, gates, and a non-applying activation plan.
- Added three-consumer diagrams and a matrix separating shared TTC corpus work
  from Garden/Admin serving-runtime evidence and from the separate GEC corpus.
- Added proposed `SourceSnapshot`, `ArtifactRef`, `BuildRequest`,
  `ProductBuilder`, `ProductEvaluator`, `ProgressSink`, and coordinator APIs.
- Added failure classifications, security/IAM boundaries, observability
  signals, AWS Batch/EventBridge/S3/DynamoDB deployment guidance, and a River
  decision rule.
- Added six precise implementation phases with exit criteria and mirrored the
  shared-library/product-deployment split into Phase 7 of `tasks.md`.
- Linked the design from the ticket index and retained the previously prepared
  fresh-root RAG-TTC reproducibility protocol for commit.

### Why

- All three products need the same lifecycle vocabulary even though their
  extraction, serving, and deployment semantics differ.
- A reusable library should make correct custody and replay the default without
  attempting to own AWS scheduling or product meaning.
- Content refresh changes data identity, whereas optimization changes one
  mutable asset; pretending these are the same would make experiment evidence
  misleading.
- The Admin scripted runtime proves chat transport and session behavior, not
  real retrieval, so the design must name the later live adapter explicitly.

### What worked

- The existing `runstore`, candidate, evaluation, comparison, gate, and report
  packages provide a concrete base for the proposed coordinator.
- GEC's cached deterministic rebuild and RAG-TTC's “resume = replay” flow model
  fit one common outer-retry/inner-replay architecture.
- One coarse AWS Batch or River job preserves global bundle identity and keeps
  per-item concurrency inside the product builder.
- The design stays pragmatic: full scans first, no CDC, no Step Functions, no
  hot reload, no per-document queue jobs, and no new database solely for River.
- The standalone default-layout upload completed successfully:

  ```text
  OK: uploaded RAGOPT Shared Production Refresh Design.pdf -> /ai/2026/08/06/RAGOPT-001
  ```

  The command did not request the editor layout.

### What didn't work

- One read-only search command accidentally placed backticks inside a zsh
  double-quoted command. The shell attempted command substitution and printed:

  ```text
  zsh:1: command not found: ragkit
  ```

  The search still returned the other requested results. Subsequent commands
  avoided shell-interpreted backticks. No file or runtime state was affected.

- The first sandboxed metadata update failed exactly as follows:

  ```text
  Error: failed to write document: open /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/design-doc/.docmgr-3007103750: read-only file system
  ```

  Re-running the same repository-scoped `docmgr doc relate` operation with
  write approval succeeded and added seven focused source relationships.

- The first combined validation/stage/commit command validated the reference
  frontmatter, then failed before staging because the repository Git metadata
  is outside the writable sandbox:

  ```text
  fatal: Unable to create '/home/manuel/code/wesen/go-go-golems/ragopt/.git/index.lock': Read-only file system
  ```

  No partial staging occurred. The exact file-scoped Git operation was rerun
  with repository write approval.

- After `docmgr changelog update`, `git diff --check` reported:

  ```text
  ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/changelog.md:139: new blank line at EOF.
  ```

  I removed only the extra terminal blank line and reran validation.

### What I learned

- “Batteries included” should mean a complete control-plane contract and local
  reference implementation, not ownership of every product and cloud concern.
- A shared TTC bundle is an artifact-sharing decision, not permission to merge
  Garden and Admin authorization, prompts, tools, answer contracts, or gates.
- Quality rejection must be a successful operational terminal result; otherwise
  infrastructure retries can waste model/provider budget on a correctly
  rejected challenger.
- The exact location of a shared inner executor should follow demonstrated
  dependency pressure. The current `flow` package is useful evidence but is not
  an AWS/River replacement.

### What was tricky to build

The sharpest boundary was reconciling two true requirements: the user wants
`ragopt` to be the reusable base package, while the v1 guardrails correctly say
it must not grow into a scheduler or workflow platform. The resolution is a
fixed coordinator and schema library with adapter interfaces. It owns one
closed lifecycle but delegates durable job admission to Batch/River, product
semantics to adapters, and mutation to deployers.

The second subtlety was the shared TTC bundle. Sharing extraction and index
artifacts is economical, but a single “TTC passed” result would hide distinct
Garden widget/customer risks and Admin authorization/tool-loop risks. The
design therefore makes the activation policy name every mandatory runtime gate.

### What warrants a second pair of eyes

- Confirm the proposed Phase 7 boundary is the desired post-v0.1 scope and does
  not pull cloud scheduling or activation mutation into the library.
- Review whether direct incumbent/challenger `ArtifactRef`s are sufficient for
  content refresh or whether a small immutable `BundleChallenger` schema is
  clearer.
- Confirm that one TTC corpus can satisfy both products' visibility policies;
  otherwise define separate derived bundles from one source snapshot.
- Validate the proposed AWS embedding-cache topology and where
  `nomic-embed-text` will run before committing infrastructure.

### What should be done in the future

- Complete the existing RAGOPT v0.1 two-product proof before implementing the
  post-v0.1 refresh coordinator.
- Freeze TTC/GEC fixture snapshots and schemas first, then implement the local
  deterministic fixture before any AWS resources.
- Fix GEC authorization and judge-accounting P0 findings before its evaluation
  adapter can produce promotion evidence.
- Stabilize the canonical Admin Chat cutover before wiring its live
  retrieval/tool resolver and evaluation arm.

### Code review instructions

- Start with `design-doc/02-production-index-build-scheduling-resumability-and-ragopt-integration.md`:
  read the executive summary, three-consumer matrix, API surface, failure
  semantics, and phased implementation plan.
- Compare Phase 7 in `tasks.md` with the explicit v1 deferrals; verify that no
  scheduler, daemon, generic workflow engine, or deployment mutator moved into
  `ragopt`.
- Inspect the source references under “Current Code Map for an Intern.”
- Run `git diff --check`, `docmgr validate frontmatter --doc <design>`, and
  `docmgr doctor --ticket RAGOPT-001 --stale-after 30`.

### Technical details

```text
Shared control plane: semantic identity -> build events -> replay -> verified artifact
                      -> native paired evaluation -> gates -> activation plan

External execution:   EventBridge Scheduler -> AWS Batch -> product refresh command
Alternative:           River coarse job -> same product refresh command

TTC evidence:          shared bundle -> Garden native gate + Admin native gate
GEC evidence:          GEC bundle -> CoinVault native gate
Mutation authority:    product deployer only; ragopt plan is non-applying
reMarkable document:   /ai/2026/08/06/RAGOPT-001/RAGOPT Shared Production Refresh Design
```

## Step 17: Audit the fresh-root proof identities and establish launch readiness

This step resumed the next sequential Phase 5 task without spending on external
providers. I verified the detached product worktree, exact candidate, source
profiles, frozen index, focused tests, first-run artifacts, and current tmux
state. The second six-cell proof is now technically ready but was not launched,
because this continuation did not add explicit approval for another external
answer/judge run.

The audit also found a real pre-v0.1 naming ambiguity: evaluation custody uses
the exact policy file's byte digest, while the gate decision uses the decoded
policy's semantic digest, and both schemas call their field `policy_digest`.
The values are not contradictory. I corrected the proof documents to pin both
identity layers and added a Phase 6 breaking cleanup task after reproduction.

### Prompt Context

**User prompt (verbatim):** (see Step 15)

**Assistant interpretation:** Continue RAGOPT-001 from its actual ledger and
worktree state, preserve detailed diary evidence and focused commits, and make
the next product-proof step genuinely ready without crossing the explicit
provider export/spend boundary.

**Inferred user intent:** Advance the reusable optimizer with evidence rather
than declarations, catch custody mistakes before another live run, and keep
the work understandable and resumable by another engineer.

**Commit (documentation):** `029c076` — "docs(proof): disambiguate policy identity layers"

### What I did

- Read the diary skill and its full reference again.
- Inspected the clean `ragopt` worktree, recent commits, all docmgr tasks, the
  detached proof worktree, the reproducibility protocol, and the latest diary.
- Confirmed the proof worktree is detached at product commit
  `90485d8539515173e88d2a9703f3a7ad4aed74bc`.
- Verified the intentionally copied `profiles.yaml` matches the source
  workspace byte-for-byte at SHA-256 `72c197d…`.
- Located the adapter's actual default index symlink and verified its manifest
  byte digest is the frozen `b2e92c…` identity.
- Ran focused RAG-TTC command/adapter/asset tests, strict candidate validation,
  and read-only comparison of the first corrected run.
- Inspected `pkg/eval`, `pkg/gate`, `pkg/compare`, and `pkg/report` to explain
  policy identity at each schema layer.
- Updated references 06 and 07 with separate suite/policy semantic and byte
  digests, no-provider readiness evidence, and the pre-v0.1 naming warning.
- Added a Phase 6 task to rename/split the ambiguous policy fields after task
  67, with no compatibility shim.

### Why

- A fresh stochastic run is useful only if every deterministic input is proven
  identical before the first provider request.
- Changing identity schemas between the two executions would make the
  reproduction comparison weaker, so the cleanup belongs after task 67.
- Recording byte and semantic digests separately prevents a reviewer from
  mistaking two correct values for drift.
- The ticket's sequential phases prohibit skipping ahead to GEC or release
  hardening while the RAG-TTC reproduction remains incomplete.

### What worked

- Focused tests passed:

  ```text
  ok  github.com/the-tree-center/rag-ttc/cmd/rag-ttc/cmds/chat/tooleval
  ok  github.com/the-tree-center/rag-ttc/cmd/rag-ttc
  ok  github.com/the-tree-center/rag-ttc/internal/ragoptassets
  ```

- `ragopt candidate validate` reproduced the exact candidate (`0d598d…`),
  parent (`4cc350…`), child (`836466…`), and changed asset digests.
- `ragopt compare` still returns `fail` with three expected and complete pairs,
  candidate contract validity 1/3, and all original hard-gate behavior.
- The proof roots contain only the first corrected run/report; no accidental
  second run or partial result exists.
- Approved tmux inspection found only the unrelated
  `ttc-garden-chat-luna-low` session and no RAGOPT proof process.

### What didn't work

- The initial sandboxed tmux inspection could not reach the tmux socket:

  ```text
  error connecting to /tmp/tmux-1000/default (Operation not permitted)
  ```

  Repeating the read-only inspection with narrowly scoped approval succeeded.

- I first guessed two asset paths without their `shared/` directory:

  ```text
  sha256sum: configs/ragopt/i5-combined-comparison-v1/gate-policy.yaml: No such file or directory
  sha256sum: configs/ragopt/i5-combined-comparison-v1/feedback-suite.yaml: No such file or directory
  ```

  `rg --files` showed the authoritative paths under `shared/`, which were then
  hashed successfully.

- I initially checked an index directory named after the answer profile and
  received:

  ```text
  zsh:5: no matches found: .cache/rag-ttc/indexes/ttc-live-luna-low/*
  ls: cannot access '.cache/rag-ttc/indexes/ttc-live-luna-low': No such file or directory
  ```

  This was a diagnostic mistake, not missing proof data. Reading the command
  source showed the locked default
  `ttc-056cbd53e148922e847ceabab1f7c4ef`; its symlink and manifest are intact.

- `docmgr changelog update` again left an extra terminal blank line, and
  `git diff --check` reported:

  ```text
  ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/changelog.md:148: new blank line at EOF.
  ```

  I removed only that extra line before the bookkeeping commit.

- My first patch to remove that line assumed docmgr preserved the file-note
  argument order and failed safely with:

  ```text
  apply_patch verification failed: Failed to find expected lines in /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/changelog.md
  ```

  I inspected the actual tail, observed that docmgr had sorted the two related
  files, and applied a context-correct patch.

### What I learned

- The first-run report's `75c56f…` and run config's `7d3fba…` are semantic and
  byte identities for the same policy, respectively.
- Suite custody already makes the distinction visible through semantic
  `suite_digest` plus byte-valued `input_digests`; policy custody should follow
  an equally unambiguous vocabulary before v0.1.
- `readlink -f` can canonicalize a nonexistent path and therefore is not proof
  that the target exists; combine it with `test`, `ls`, or hashing.
- The frozen index is selected by the adapter's `--index-bundle` default, not
  by the provider profile name or `--cache-directory`.

### What was tricky to build

The policy discrepancy initially looked like evidence drift because the same
field name appears in multiple artifact schemas. Reading the code established
the exact chain: `eval.RunConfig.PolicyDigest` is the copied YAML byte digest;
`compare.Report.PolicyDigest` preserves it; `gate.LoadPolicy` computes both
`ByteDigest` and canonical `Digest`; `gate.Evaluate` first checks the bytes and
then stores the semantic digest in the decision; the promotion plan copies the
semantic decision digest. Preserving both values in the protocol fixes review
ambiguity without perturbing the experiment between repetitions.

### What warrants a second pair of eyes

- Confirm that the post-proof breaking cleanup should name fields explicitly
  `policy_byte_digest` and `policy_semantic_digest`, rather than carrying a
  nested policy identity object.
- Recheck the second run's complete `input_digests` map, not only the summary
  identity table.
- Treat any changed decision as stochastic evidence requiring native paired
  review, not automatic promotion or automatic harness failure.

### What should be done in the future

- Obtain explicit approval for the second six-cell provider-backed feedback
  run, then launch it in a new tmux session from the detached proof worktree.
- Generate comparison/report artifacts for the new run and apply every check
  in reference 07.
- After task 67, perform the policy identity field cleanup before v0.1 and
  document the product-adapter migration as a deliberate breaking change.

### Code review instructions

- Start with reference 07's identity table and policy naming warning, then
  compare reference 06's corrected quick-reference rows.
- Inspect `pkg/eval/runner.go:loadPolicyIdentity`,
  `pkg/gate/policy.go:PolicyDocument`, `pkg/gate/evaluate.go:Evaluate`, and
  `pkg/report/types.go:PromotionPlan`.
- From `/tmp/rag-ttc-ragopt-proof`, rerun:

  ```bash
  go test ./cmd/rag-ttc/cmds/chat/tooleval ./cmd/rag-ttc ./internal/ragoptassets -count=1
  go run github.com/go-go-golems/ragopt/cmd/ragopt candidate validate \
    --bundle configs/ragopt/i5-combined-comparison-v1 \
    --manifest candidate.yaml --format json
  go run github.com/go-go-golems/ragopt/cmd/ragopt compare \
    --run experiments/ragopt-runs/20260806T180004.824520651Z-ttc-i5-feedback-a1781d11159f \
    --format json
  ```

### Technical details

```text
Proof command status: ready, not launched
External cells:       6 (3 feedback cases x 2 arms)
Validation cells:     0 (feedback gate remains closed)
Product commit:       90485d8539515173e88d2a9703f3a7ad4aed74bc
Index manifest bytes: sha256:b2e92c7cd9964b6ab36c0b3dc47f96b1cdc72f1966c3ceeaf6772abf82f194dd
Policy bytes:         sha256:7d3fba806e3824b26ff817b103fd6415cfe66117edc66ad3050c6ab6476b4d8e
Policy semantic:      sha256:75c56ff46aab7bb2b281942514faf8a7de94fa1bb4c06e2e0df9548fb7ddb83f
Suite bytes:          sha256:c22b6b186d5f4bd7c2a9177271bfa0005946474b67f78fc7ca85108b9106518d
Suite semantic:       sha256:b009f9e913179007bf2da04cadb2e75200e22e1ab38689513015721ee76781b6
```

## Step 18: Freeze the provider consent and spend envelope

This step made the pending external action reviewable instead of asking for a
vague “provider run” approval. I traced every provider seam in the product
adapter and native runtime, enumerated the exact benchmark questions and data
classes sent externally, calculated hard per-stage call ceilings, documented
local artifact retention, and wrote a narrow approval statement that excludes
validation and retries as new runs.

The inspection corrected an important initial suspicion. Reusing the shared
cache does not turn task 67 into a byte-identical answer replay: the Admin-era
tool loop sends answer iterations through a direct, budgeted Geppetto engine,
while only query embeddings and judge generations use the content-addressed
provider cache. The second run therefore produces fresh answer-model behavior
while safely reusing exact embedding/judge work when possible.

### Prompt Context

**User prompt (verbatim):** (see Step 15)

**Assistant interpretation:** Continue useful RAGOPT-001 work while explicit
provider approval is absent, preserving the external-action boundary and
making the eventual authorization exact, bounded, and auditable.

**Inferred user intent:** Keep progressing toward a real two-run product proof
without accidental spend, data export, validation expansion, or a fake cached
repetition that does not exercise the product runtime.

**Commit (documentation):** `c71db2c` — "docs(proof): freeze provider consent envelope"

### What I did

- Read the diary skill and full reference again for this continuation.
- Read the three-case feedback suite, runtime lock, judge lock, adapter cell
  executor, materialized tool configuration, and two-step judge implementation.
- Inspected one prior native session config and judge summary to verify resolved
  provider/model identities without reading credentials or environment values.
- Traced `pkg/app/chat/tool_runtime.go` to distinguish direct answer-engine
  calls from cached query embeddings.
- Traced `pkg/rag/generation` and the judge code to confirm exact cache keys,
  cache-hit behavior, and invalid-answer judge skipping.
- Aggregated the first corrected run's native sessions: six cells, 22 answer
  iterations, 17 search-tool invocations, and two uncached judge operations.
- Added the external execution consent envelope to reference 07: exact
  questions, payload classes, provider identities, hard ceilings, local
  retention, exclusions, and reusable approval text.
- Did not launch a tmux proof session or contact any provider.

### Why

- External approval should name what data leaves the machine and the maximum
  operations it authorizes.
- “Use the same cache directory” has different meanings at different runtime
  layers; assuming all generation is cached would produce the wrong proof claim.
- A call ceiling is not a cost ceiling when the adapter explicitly allows
  unpriced operations and supplies no dollar estimates.
- Validation is 28 cells and must remain excluded after the feedback candidate
  failed its hard gates.

### What worked

- The frozen feedback suite contains only three public product-comparison
  questions, not customer conversations.
- Scoped SQL is disabled in the materialized tool configuration, so orders,
  customers, private logistics facts, and live database results are excluded.
- Provider identity is explicit in retained artifacts:
  `openai/text-embedding-3-small` for queries,
  `openai-responses/gpt-5.6-luna` for answers, and
  `openai/gpt-5.6-luna` for both judge stages.
- The hard six-cell ceilings are 18 embedding operations, 24 fresh answer
  generations, and 12 judge generations; cells may terminate or hit caches
  below those limits.
- The first corrected run's observed 22/17/2 answer/search/judge counts agree
  with the locked budgets and invalid-answer judge policy.

### What didn't work

- My first cache search was too broad and matched large native session
  transcripts containing full evidence payloads. Tool output was truncated:

  ```text
  Warning: truncated output (original token count: 30028)
  Warning: truncated output (original token count: 80769)
  ```

  I stopped using the broad result and switched to targeted source ranges and
  narrow `jq` projections. No files changed and no provider was contacted.

- The shared cache initially suggested that the second run might be a pure
  deterministic replay. Source inspection disproved that assumption: the tool
  answer runtime wraps `providers.Engine` with `limitedEngine`, not a cached
  generator. I corrected the protocol before requesting authorization.

- `docmgr changelog update` again appended an extra blank line, so
  `git diff --check` reported:

  ```text
  ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/changelog.md:156: new blank line at EOF.
  ```

  I removed that single line before committing the diary bookkeeping.

- The first combined cleanup patch used an incomplete diary context and failed
  safely with `apply_patch verification failed: Failed to find expected lines`.
  I inspected the actual paragraph and reapplied the focused patch; no partial
  edit occurred.

### What I learned

- The product has three different cache/budget semantics in one cell: cached
  query embeddings, fresh answer-engine iterations, and cached/skippable judge
  generations.
- Remote provider cached-input tokens shown in answer usage are billing/runtime
  behavior and must not be confused with the local content cache.
- The judge sees less data in stage one (question + answer) than stage two
  (question + admitted evidence + extracted statements).
- Because `same_family_verdicts=true`, decision stability must be interpreted
  with the known same-family judge risk; changing the judge now would break
  reproduction identity.

### What was tricky to build

The adapter passes one `cache-directory` into the runtime and judge, which
looks like a uniform replay switch. It is not. The tool runtime creates a
cached embedder but executes answer turns through a direct Geppetto engine with
only a limiter. The judge separately uses cached generation steps. Establishing
the actual consent envelope required tracing those concrete call sites rather
than inferring behavior from CLI flags or directory names.

The second sharp edge is vocabulary: “provider calls” in the common outcome
counts product answer iterations only and deliberately excludes judge overhead.
The consent ceiling therefore lists heterogeneous answer, embedding, and judge
operations separately before showing the combined maximum of 54.

### What warrants a second pair of eyes

- Confirm that the three benchmark questions and public corpus evidence are
  acceptable for export to the named providers.
- Confirm that approval of bounded operation counts is acceptable despite the
  absence of a configured USD ceiling.
- Review whether future product adapters should require priced preflight or an
  explicit `--max-estimated-usd` before external execution.
- Verify after the run that local cache hits and remote cached-input tokens are
  reported separately.

### What should be done in the future

- Obtain the explicit six-cell-only approval recorded in reference 07.
- Launch exactly one fresh run in tmux; resume that same active run after a
  transport interruption rather than creating another run.
- Preserve the feedback-before-validation gate regardless of the operational
  outcome.
- Consider a priced-preflight requirement after the proof; do not change the
  locked adapter between the two reproduction runs.

### Code review instructions

- Read “External Execution Consent Envelope” in reference 07.
- Verify call topology in
  `/tmp/rag-ttc-ragopt-proof/pkg/app/chat/tool_runtime.go` and judge payloads in
  `cmd/rag-ttc/cmds/experiments/answerquality/judge.go`.
- Recompute first-run observed counts with:

  ```bash
  find experiments/ragopt-runs/20260806T180004.824520651Z-ttc-i5-feedback-a1781d11159f/native \
    -name turns.jsonl -print0 | xargs -0 jq -s \
    '{cells:length,answer_provider_calls:(map(.agent.iterations|length)|add),tool_calls:(map(.agent.tool_calls|length)|add)}'
  ```

### Technical details

```text
Authorized split:              feedback only
Coordinates:                   3 cases x 2 arms x 1 repeat = 6 cells
Max query embedding work:      18 provider operations (cache misses only)
Max answer generation work:    24 provider operations (fresh iterations)
Max judge generation work:     12 provider operations (valid answers/cache misses only)
Combined heterogeneous ceiling: 54 provider operations
Configured USD ceiling:        none; AllowUnpriced=true
Scoped SQL:                     disabled
Validation authorization:      none
New proof run created:          no
```

## Step 19: Stop at the explicit provider authorization boundary

This step performed the terminal blocker audit for the current continuation.
Task 67 remains the next sequential requirement, every deterministic input and
call ceiling is frozen, and no safe local preparation remains that would make
the proof more valid. Three consecutive goal continuations supplied no explicit
authorization for the external six-cell answer/embedding/judge execution.

I therefore did not launch, simulate, weaken, or replace the proof. I recorded
the exact resume condition in the ticket index and prepared to mark the active
goal blocked. This preserves the user's authority over external data export and
unpriced provider spend while leaving the repository in a clean, immediately
resumable state.

### Prompt Context

**User prompt (verbatim):** (see Step 15)

**Assistant interpretation:** Continue RAGOPT-001 persistently, but follow the
strict blocked audit when the same external-authorization condition has
recurred for three consecutive goal turns and no meaningful in-scope work
remains.

**Inferred user intent:** Maintain genuine progress and evidence without
silently spending money, exporting benchmark payloads, skipping sequential
proof requirements, or manufacturing a weaker substitute for task 67.

**Commit (documentation):** `ed2f4a9` — "docs(ticket): record provider proof blocker"

### What I did

- Read the diary skill and its full reference again.
- Re-inspected the `ragopt` status and task ledger, detached proof worktree
  commit, dirty-file scope, and proof run directories.
- Confirmed task 67 is still unchecked and tasks 68–80 are sequentially gated
  behind the two-product proof.
- Confirmed the proof root still contains only the first corrected run
  `20260806T180004.824520651Z-ttc-i5-feedback-a1781d11159f`.
- Added the exact external blocker and resume condition to the ticket index.
- Did not contact providers, start tmux, modify RAG-TTC, run validation, begin
  GEC integration, or implement later-phase RAGOPT features out of order.

### Why

- The six-cell command exports questions and public corpus evidence and permits
  up to 54 heterogeneous provider operations without a USD ceiling.
- The consent envelope was explicitly presented twice after an earlier denied
  launch; an automatic continuation is not approval.
- Skipping to GEC or Phase 6 would contradict the ledger's sequential exit
  criteria and weaken the evidence this project exists to enforce.
- The strict goal policy requires a blocked terminal state after the same
  condition recurs across at least three consecutive turns and no meaningful
  progress remains.

### What worked

- The `ragopt` repository was clean before this documentation checkpoint.
- The proof worktree remains at the corrected product commit
  `90485d8539515173e88d2a9703f3a7ad4aed74bc` with only the intentional locked
  `profiles.yaml` copy modified.
- Exactly one corrected proof run exists; no second or partial run was created.
- References 06 and 07 contain all identities, commands, custody checks,
  provider payloads, cache semantics, operation ceilings, and decision rules
  needed to resume without rediscovery.

### What didn't work

- The required authorization was not supplied after three consecutive
  continuation opportunities. The minimum sufficient approval remains:

  ```text
  I approve one second RAG-TTC I5 feedback run: six cells only, with ceilings of 24 answer calls, 18 embedding calls, and 12 judge calls. Do not run validation.
  ```

- No technical workaround is valid. A cached-only simulation would not prove
  fresh answer behavior; a smaller run would not satisfy all six coordinates;
  a larger run would exceed scope; and later-phase implementation would violate
  the canonical ordering.

- `docmgr changelog update` appended its known extra blank line and
  `git diff --check` reported:

  ```text
  ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/changelog.md:165: new blank line at EOF.
  ```

  I removed only that line before the final checkpoint commit.

### What I learned

- Persistence includes stopping cleanly when the remaining action needs new
  authority; it does not broaden authorization.
- The precommitted consent envelope makes this blocker precise: the missing
  input is a yes/no authorization, not architectural clarification.
- A blocked goal and an active ticket are compatible. The ticket remains active
  and resumable; only this autonomous execution chain must stop.

### What was tricky to build

The challenge was distinguishing “more work is possible” from “more work is
aligned.” Phase 6 and GEC tasks are available on paper, but the ticket declares
phases sequential and uses the second RAG-TTC run as the evidence gate. Working
around that gate would recreate the fake-completion problem RAGOPT is designed
to prevent. After completing the identity, readiness, cache, payload, and spend
audits, only the externally authorized command remains aligned.

### What warrants a second pair of eyes

- Confirm that the six-cell consent wording is acceptable as written.
- On resume, verify again that no second run appeared and that product HEAD,
  profile digest, index manifest digest, and task 67 inputs remain unchanged.
- Ensure the goal is resumed as a fresh blocked audit rather than assuming an
  old authorization still applies to later provider runs.

### What should be done in the future

- Resume only after explicit authorization of the exact six-cell envelope.
- Launch one tmux run, monitor it, retain all cells, compare identities and
  canonical deltas, write the second-run report, and check task 67 only if the
  protocol passes.
- Continue to GEC integration only after that reproduction evidence exists.

### Code review instructions

- Start with the ticket index status paragraph, then read reference 07's
  “External Execution Consent Envelope.”
- Confirm `docmgr task list --ticket RAGOPT-001` leaves task 67 unchecked.
- Confirm `/tmp/rag-ttc-ragopt-proof/experiments/ragopt-runs` contains only the
  first corrected run.
- Confirm `git status --short` is clean in `ragopt` after the checkpoint commit.

### Technical details

```text
Blocking condition: explicit authorization for one six-cell provider run absent
Consecutive occurrences: 3
Task blocked: RAGOPT-001 task 67
Second run created: no
Provider calls made in this step: 0
Validation authorized or run: no
Resume document: reference/07-rag-ttc-i5-fresh-root-reproducibility-protocol.md
```

## Step 20: Run and verify the authorized fresh-root RAG-TTC proof

The user supplied the exact authorization requested in Step 19 and also set a
standing ceiling for later provider work below one million tokens or one US
dollar. I resumed the frozen proof without changing the product commit,
profile, index, suite, candidate, policy, or command. The second run completed
all six feedback coordinates. I did not run validation.

The fresh-root run reproduced the important decision: the I5 candidate is not
eligible for promotion. This is a successful harness proof and a failed
candidate. The byte-identical run configuration, semantic identities, copied
input custody, complete pair join, and deterministic gate decision all agree
with the protocol. Stochastic answer behavior changed which cells were valid,
so the metric deltas are not byte-identical to run one; the difference is
visible and explainable rather than hidden.

### Prompt Context

**User prompt (verbatim):**

> I approve one second RAG-TTC I5 feedback run: six cells only, with ceilings of 24 answer calls, 18 embedding calls, and 12 judge calls. Do not run validation.
>
> In fact you can run anything with a provider budget < 1M token (or $1)

**Assistant interpretation:** Execute exactly the frozen six-cell feedback
proof, stay below all stated call ceilings, leave validation unrun, then verify
the complete evidence chain. Treat the sub-one-million-token or sub-one-dollar
statement as authorization for future provider work within that bound, not as
permission to alter the frozen proof.

**Inferred user intent:** Finish the evidence required to distinguish a real,
reproducible optimization harness from a result that merely appears complete,
while bounding provider exposure and spend.

### What I did

- Revalidated the detached proof worktree at product commit `90485d8539515173e88d2a9703f3a7ad4aed74bc`.
- Confirmed only the first corrected run existed before launch and that the
  frozen profile, index, suite, and policy bytes still matched reference 07.
- Launched the exact `tool-loop ragopt --split feedback` command in tmux and
  monitored it through terminal completion.
- Retained fresh run
  `20260806T194517.057830400Z-ttc-i5-feedback-98a227142ee5`.
- Ran `ragopt compare` and generated a non-applying promotion plan and review.
- Verified all six expected coordinates are unique and present.
- Recomputed every copied-input SHA-256 and byte size from `inputs/manifest.json`.
- Verified both proof-run `config.json` files are byte-identical with digest
  `99ca5018c34d85b1693894c8670f57f7a00044e5f62065376bcf699be61e627c`.
- Recomputed all native artifact digests and sizes and checked the referenced
  session directories.
- Counted 22 answer calls, 17 `ttc_search` calls, and 8 judge work calls, all
  below the authorized 24/18/12 ceilings.
- Summed 108,548 answer-plus-judge provider tokens, below the later standing
  one-million-token ceiling.
- Confirmed no validation coordinate or validation run was created.

### Why

- Task 67 is the evidence gate before starting the GEC integration.
- A second fresh root tests input custody, adapter behavior, pair construction,
  and gate determinism without requiring the candidate itself to pass.
- Retaining invalid and failed cells prevents nondeterministic generation from
  making a weak candidate look better by omission.

### What worked

- Run state is `complete`; the summary records 6 expected, 6 completed, and 1
  failed cell.
- Suite semantic identity is
  `sha256:b009f9...`; candidate identity is `sha256:0d598d...`; parent and child
  snapshots are `sha256:4cc350...` and `sha256:836466...`.
- Complete pairing passed for all three cases.
- The candidate completed 3/3 cells but was contract-valid for only 2/3.
- Answer relevance tied at 1.0 on both judged pairs.
- Candidate mean faithfulness was 0.947222 versus incumbent 0.989362, a
  -0.042139 paired mean delta over the two covered pairs.
- The promotion plan is `review_required`, `human_apply_required: true`, and
  `decision: fail`; it contains no application action.

### What didn't work

- My first copied-input verification used `path` as a zsh loop variable. In
  zsh, `path` is tied to `PATH`, so the command destroyed its own executable
  lookup and failed exactly as follows:

  ```text
  zsh:4: command not found: sha256sum
  zsh:4: command not found: awk
  zsh:5: command not found: stat
  MISMATCH evaluation-suite inputs/evaluation-suite.json want=sha256:c22b6b186d5f4bd7c2a9177271bfa0005946474b67f78fc7ca85108b9106518d/903 got=/
  ```

  I renamed the variable to `copied`, used explicit `/usr/bin` paths, stripped
  the manifest's `sha256:` prefix, and obtained `input-custody-ok`.

- A first `jq` custody projection treated the top-level input manifest as an
  object instead of an array and failed with:

  ```text
  jq: error (at .../inputs/manifest.json:149): Cannot index array with string "role"
  ```

  I corrected the selector to `.[]` and did not accept the earlier false shell
  success marker.

- A tool-call diagnostic assumed every `result` was an object and failed with:

  ```text
  jq: error (.../turns.jsonl:1): Cannot index string with string "error"
  ```

  Three failed tool results are strings. I changed the diagnostic to inspect
  `type` before object operations and confirmed 17 total search calls.

### What I learned

- Reproducibility here means stable frozen identities, complete evidence, and
  the same gate decision. It does not mean a stochastic model must emit the
  same answers or identical metric deltas on every run.
- Run two gave two fully judged pairs where run one gave only one candidate
  judged cell. The additional coverage strengthened, rather than reversed, the
  rejection: the candidate lost both faithfulness comparisons.
- The incumbent failure on `ttc-expand-056` included a fourth search request
  rejected by the per-cell resource limiter, followed by the four-iteration
  ceiling. This remains evidence, not a dropped sample.

### What was tricky to build

The main review problem was keeping four different identities separate: source
file byte digests, semantic suite/policy digests, snapshot identities, and the
byte digest of the full run configuration. They answer different questions and
must not be compared interchangeably. The second problem was interpreting
stochastic cell changes without weakening the predeclared gates. The correct
comparison is frozen inputs plus exact pairs plus the same decision algorithm,
with every changed outcome reported.

### What warrants a second pair of eyes

- Review the failed incumbent and invalid challenger `ttc-expand-056` traces.
- Confirm that two judged pairs are enough evidence for a rejection but not for
  a broad prompt-quality conclusion.
- Review the distinction between the policy byte digest and semantic policy
  digest before the Phase 6 schema rename.

### What should be done in the future

- Publish the detailed second-run report and check task 67.
- Begin GEC only after its authorization and judge-accounting P0 findings are
  fixed.
- Freeze and run one GEC human-authored candidate twice before generalizing any
  additional RAGOPT API.

### Code review instructions

- Compare the two `config.json` files with `cmp -s` and `sha256sum`.
- Read all six rows in `results/cells.jsonl`; do not average away invalid or
  failed outcomes.
- Recompute copied-input and native-artifact digests from the second run root.
- Read `promotion-plan.json` beside the generated review and confirm that no
  apply mechanism exists.

### Technical details

```text
Second run:           20260806T194517.057830400Z-ttc-i5-feedback-98a227142ee5
Coordinates:          3 cases x 2 arms x 1 repeat = 6
Completed / failed:   6 / 1
Complete pairs:       3 / 3
Candidate valid:      2 / 3
Answer calls:         22 / 24
Search calls:         17 / 18
Judge work calls:     8 / 12
Provider tokens:      108,548 / 1,000,000
Decision:             fail — do not promote I5
Harness proof:        pass — frozen evidence reproduced the rejection
Validation:           not run
```

## Step 21: Create the physical RAGOPT mission board

I translated the canonical task ledger and the just-verified proof result into
a compact brutalist thermal work slip. The print is intentionally operational:
it says what RAGOPT owns, what has been proved, why I5 remains rejected, which
GEC tasks are next, what is needed for v0.1, and which attractive but premature
features remain forbidden.

### Prompt Context

**User prompt (verbatim):**

> Printout an almanach brutalist writeup of the ticket and the tasks we are doing and ahead of us, so I can keep it in front of me what we are doing.

**Assistant interpretation:** Produce a physical, glanceable mission board on
the Almanach printer—not another long design document—and retain its source in
the RAGOPT ticket.

**Inferred user intent:** Keep the project boundary and sequential next actions
visible during implementation, especially the distinction between the shared
harness, product-owned adapters, later production refresh work, and explicitly
deferred overengineering.

### What I did

- Read the current Phase 0–7 task ledger rather than summarizing from memory.
- Used the Almanach work-slip DSL with `theme: brutalist`, `paperWidth: 384`,
  `bodyScale: 1`, and `feedLines: 8`.
- Created `sources/almanach/21-ragopt-mission-board.yaml` inside the ticket.
- Structured the slip as purpose, current state, completed mechanisms, next
  GEC proof tasks, v0.1 hardening, later refresh work, and guardrails.
- Kept the decision explicit: harness proof passes; I5 candidate fails;
  validation remains closed.

### Why

- A physical board should optimize for orientation and sequencing, not repeat
  the full intern design guide.
- The ticket's largest risk is boundary drift: confusing Admin Chat, Garden,
  CoinVault/GEC, RAGOPT, and AWS deployment responsibilities.
- Showing deferred features is useful only when phrased as “do not build,” not
  as an attractive parallel backlog.

### What worked

- The existing task ledger maps cleanly to a one-page thermal hierarchy.
- The brutalist work-slip primitives support dense but readable operational
  content without the editor-style reMarkable layout the user rejected.
- The YAML is versioned beside the diary and can be reprinted unchanged.

### What didn't work

- I initially looked for a stale example path named
  `examples/11-work-slip-brutalist.yaml`; it does not exist in the installed
  skill directory:

  ```text
  sed: can't read /home/manuel/.codex/skills/almanach-printing/examples/11-work-slip-brutalist.yaml: No such file or directory
  ```

  The skill points to the live examples repository. I used
  `/home/manuel/code/wesen/go-go-golems/almanach/examples/layouts/11-decision-sheet.yaml`
  and `12-triage-card.yaml` instead.

### What I learned

- The best desk-facing summary is the dependency order, not an architecture
  inventory: prove TTC twice, fix GEC P0s, prove GEC twice, harden v0.1, then
  build shared refresh contracts.
- “No overengineering” is concrete here: no reflector, generic workflow
  engine, automatic deployment, transcript warehouse, or multi-mutation search.

### What was tricky to build

The slip needed to remain accurate across two different meanings of success.
The integration proof succeeded because the evidence chain reproduced. The I5
optimization failed because it did not satisfy the promotion policy. Collapsing
those into one green/red status would teach exactly the wrong mental model.

### What warrants a second pair of eyes

- Check the physical hierarchy at arm's length: current state and next action
  should dominate the later production-refresh material.
- Confirm the longer four-item GEC list remains legible on 58 mm paper.

### What should be done in the future

- Reprint the same source after major phase transitions only when the board's
  state changes materially.
- Keep short per-diary-entry tickets separate from this larger mission board.

### Code review instructions

- Compare the YAML claims to `tasks.md`, especially Phase 5, Phase 6, Phase 7,
  and “Explicitly deferred.”
- Inspect or dry-run with `almanach-render-service inspect` or
  `print-remote --dry-run` before changing typography or density.

### Technical details

```text
Layout:       sources/almanach/21-ragopt-mission-board.yaml
Theme:        brutalist
Paper width:  384 px / 58 mm
Body scale:   1
Feed lines:   8
Primary next: GEC authorization and judge-accounting P0 fixes
```

## Step 22: Make a brutalist print the closing action of every diary entry

The user added a session workflow rule: after each detailed diary entry, print
a short brutalist ticket explaining what happened and how it went. I adopted a
non-recursive convention: the ticket closes the diary step, and the printing
operation itself does not create another diary step unless it materially fails
or changes the project.

### Prompt Context

**User prompt (verbatim):**

> btw After each diary entry, print a short brutalist ticket on the printer about what you did
>  and how it went.

**Assistant interpretation:** Every new strict diary step in this continuing
session must be followed by a physical, concise brutalist summary ticket.

**Inferred user intent:** Maintain a paper trail that makes ongoing progress,
failures, and immediate next actions visible without reopening the full ticket.

### What I did

- Declared the non-recursive closing convention to the user.
- Created `sources/almanach/22-diary-print-policy.yaml`.
- Retrospectively created and printed a Step 20 proof ticket after writing Step
  20; the remote service returned `ok: true`, `printed: true`, width 384, height
  580.
- Printed the Step 21 mission board after writing Step 21; the remote service
  returned `ok: true`, `printed: true`, width 384, height 1757, in 3 segments.
- Prepared this short policy ticket as the closing print for Step 22.

### Why

- A closing ticket preserves the requested “what and how” summary while the
  Markdown diary remains the authoritative technical record.
- Avoiding a new diary entry solely for a successful print prevents an infinite
  diary → print → diary recursion.

### What worked

- Both earlier remote prints completed successfully through the crib service
  and the printer acknowledged them.
- The short-ticket format fits one result, a few measures, and one next action.

### What didn't work

- Nothing failed while introducing the policy. Any later failed print will be
  recorded in the next substantive diary step and retried at most within the
  debugging limit.

### What I learned

- The physical ticket is most useful as a checkpoint receipt, while the larger
  mission board is useful as a durable orientation artifact. They should stay
  separate.

### What was tricky to build

The workflow has to preserve chronology without requiring a diary mutation
after the print result is known. The practical rule is: write the substantive
step with its evidence, print the closing ticket, and report print success in
the user-facing update. A print failure becomes evidence in the next genuinely
substantive step rather than recursively generating bookkeeping.

### What warrants a second pair of eyes

- Confirm that one short print per substantive diary step is the intended
  cadence; this can produce several physical slips during a long session.

### What should be done in the future

- Include step number, ticket, outcome, one or two measurements, and next action
  on every short ticket.
- Store every layout source under this ticket's `sources/almanach/` directory.

### Code review instructions

- Verify each future `## Step N` added during this session has a corresponding
  numbered YAML source and successful print response.

### Technical details

```text
Diary authority: Markdown reference 01
Closing format:  brutalist work slip
Body scale:      1
Feed lines:      8
Recursion rule:  successful print does not create its own diary step
Policy state:    active
```

## Step 23: Publish the second-run decision into the ticket ledger

I converted the verified run evidence into the durable ticket checkpoint. The
new reference explains the exact identities, every cell, the hard-gate result,
provider ceilings, the difference between stochastic outputs and reproducible
decisions, and the next GEC dependency. I checked Phase 5 task 67 and replaced
the ticket index's obsolete provider-blocked status.

### Prompt Context

**User prompt (verbatim):** (see Step 20)

**Assistant interpretation:** Complete the authorized proof as a reviewable
ticket artifact, not merely as an ephemeral terminal result.

**Inferred user intent:** Ensure the team can resume from repository evidence
and understand both why the harness proof passed and why the tested candidate
failed.

### What I did

- Added reference 08, “Second RAG-TTC I5 Fresh-Root Proof and Reproducibility
  Decision.”
- Documented the full identity table, six-cell table, gate table, budget
  accounting, run-one/run-two interpretation, custody procedure, commands, and
  next sequential work.
- Checked task 67 in the canonical task ledger.
- Updated the ticket index from “provider blocked” to “RAG-TTC reproduced; GEC
  next.”
- Added the checkpoint and Almanach sources to the changelog.
- Ran `docmgr status --ticket RAGOPT-001`; it reports an active, non-stale
  ticket with ten managed documents.
- Ran `git diff --check`, found the two literal trailing spaces inherited from
  the quoted user prompt, and removed only those whitespace bytes.

### Why

- The `/tmp` proof root is valuable evidence but is not a durable team handoff.
- Task state must reflect the proof criterion rather than the candidate's
  promotion result.
- The index is the ticket entry point and must not continue claiming the second
  run is absent after it has completed.

### What worked

- The report cleanly separates harness PASS, candidate FAIL, and validation NOT
  RUN.
- `docmgr status` recognizes reference 08 and reports no stale ticket.
- The task ledger now makes GEC integration the next sequential Phase 5 work.

### What didn't work

- `docmgr task check --ticket RAGOPT-001 --id 67` could read the external
  RAGOPT repository but its subprocess could not write through the current
  filesystem sandbox:

  ```text
  Error: failed to write tasks file /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/tasks.md: open /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/tasks.md: read-only file system
  ```

  The command made no partial change. I applied the identical one-line checkbox
  edit with the workspace patch mechanism.

- The first `git diff --check` reported:

  ```text
  reference/01-implementation-diary.md:2976: trailing whitespace.
  ```

  The whitespace came from preserving the user's Markdown hard line break. I
  removed the two trailing spaces while preserving every word of the prompt.

### What I learned

- A provider-blocked status becomes actively misleading once authorization is
  supplied; entry-point status pages need updating in the same checkpoint as
  the detailed report.
- The phrase “verify canonical deltas” must not be read as “require identical
  stochastic metric values.” The report makes the correct rule explicit:
  compute deltas canonically from complete frozen pairs and explain changes.

### What was tricky to build

The report must be strong enough to authorize the next phase without claiming
more than six cells demonstrate. It therefore treats stable rejection as
integration evidence, while explicitly refusing to generalize that I5 is
universally harmful. This keeps the product decision, harness decision, and
scientific scope separate.

### What warrants a second pair of eyes

- Review reference 08's task-67 conclusion against reference 07's predeclared
  criteria.
- Confirm that task 68 should now begin with task 69's P0 repairs rather than
  running the existing GEC evaluator unchanged.

### What should be done in the future

- Commit this documentation checkpoint intentionally.
- Then inspect the GEC P0 findings and scope the minimum product-owned fixes.
- Continue printing one short ticket after each later diary step.

### Code review instructions

- Start at the ticket index and follow reference 07 to reference 08.
- Confirm task 67 is checked and tasks 68–71 remain open.
- Run `git diff --check` and `docmgr status --ticket RAGOPT-001`.
- Compare the printed mission-board YAML to the task ledger.

### Technical details

```text
New reference: reference/08-second-rag-ttc-i5-fresh-root-proof-and-reproducibility-decision.md
Task closed:   67
Next tasks:    68 and 69 — GEC integration boundary and P0 repairs
Ticket docs:   10
Ticket stale:  false
```

## Step 24: Repair GEC authorization and judge accounting on an isolated branch

This step began the second product integration required by Phase 5. I audited
the current GEC repository and confirmed that the colleague's Admin Chat and
AWS work had left the main worktree heavily dirty, while the P0 retrieval and
judge findings still reproduced at `acc53c5`. I therefore created the isolated
branch `codex/ragopt-phase5-gec` in `/tmp/gec-ragopt-phase5` and changed only
the product-owned knowledge retrieval, evaluation, and tests.

The implementation now removes unauthorized lexical and vector candidates
before fusion and rechecks them before external reranker hydration. It also
rebuilds the answer record from each cumulative timeline snapshot, rejects
timeline HTTP failures, fails malformed zero-statement verdict calls, validates
judge indices and required fields, and retains answer/judge failures in
explicit summary denominators. The affected packages and CoinVault build pass.
The code is staged but not committed because two pre-commit attempts collided
with another GolangCI-Lint process; I stopped rather than bypassing the hook.

### Prompt Context

**User prompt (verbatim):**

> work on RAGOPT-001. Keep a detailed diary as you work, commit at appropriate intervals, make sure to use the diary format form the skill

**Assistant interpretation:** Continue the canonical RAGOPT-001 phase order,
keep strict implementation evidence, and checkpoint verified work without
overlapping the colleague's Admin Chat stream.

**Inferred user intent:** Complete a real two-product proof of the reusable
optimization harness while preserving reviewability, avoiding fake completion,
and committing only after the repository's own checks pass.

### What I did

- Read the diary skill and full reference, and reread the Almanach printing
  skill because every diary entry now requires a closing paper ticket.
- Inspected the GEC main worktree, Git history, worktrees, optimization handoff,
  deep review P0 findings, blueprint Phase 0, task files, and owning source.
- Confirmed the dirty main worktree contains another stream's Admin Chat/AWS
  changes and did not modify it.
- Created worktree `/tmp/gec-ragopt-phase5` on branch
  `codex/ragopt-phase5-gec` from `acc53c5f75a9b14e2b593bfd2ffd2836ac0abadf`.
- Added per-channel authorization in `internal/knowledge/service.go`:
  authorized hits are filtered and renumbered before RRF; fused candidates are
  rechecked before hydration and external reranking; final projection remains a
  defense-in-depth check.
- Updated the RRF sweep to authorize each channel before refusion.
- Added a spy-reranker test proving analyst-only text never crosses a public
  rerank request and a channel test proving scope plus role filtering happens
  before fusion.
- Made `AnswerClient.extract` a deterministic snapshot projection, added
  timeline status validation, and introduced the typed `ErrAnswerTimeout`.
- Made zero-statement verdict errors fatal and required ordered statement
  indices, supported flags, bounded relevance, and explicit abstention.
- Replaced success-only judge summaries with explicit selected, attempted,
  answered, judged, answer-error, judge-error, timeout, skipped-policy, and
  metric-denominator counts.
- Excluded abstentions from earned-faithfulness means while retaining them in
  relevance and explicit abstention counts.
- Added focused tests for every repaired invariant.

### Why

- RAGOPT task 69 forbids using GEC results for promotion until authorization
  and judge-accounting P0 findings are fixed.
- Filtering after reranking can disclose analyst text to an external endpoint
  and lets disallowed candidates affect fused ranks.
- Success-only means allow hard questions that time out or fail judging to
  improve a candidate by disappearing.
- Cumulative polling must be projected, not appended, or SQL evidence grows on
  every poll.
- A malformed judge response is an error outcome, not a vacuous success.
- A new generic Ragkit filtering API or scope-partition builder was unnecessary
  for the confidentiality fix and would exceed the pragmatic phase scope.

### What worked

- The contract-first test initially failed at the missing authorization-aware
  rerank signature, then passed after implementation.
- Focused new tests passed, including the spy reranker, repeated polling,
  non-200 timeline, verdict failure, malformed semantics, and denominator
  accounting cases.
- Complete affected-package tests passed:

  ```text
  ok github.com/go-go-golems/gec-rag/internal/knowledge
  ok github.com/go-go-golems/gec-rag/cmd/coinvault/cmds
  ```

- `go build ./cmd/coinvault` passed.
- The repository pre-commit `test` hook passed every package under
  `./cmd/... ./internal/...`, including the colleague's Admin Chat packages.
- `git diff --check` passed before staging.

### What didn't work

- The first focused test run inside the restricted sandbox could not create an
  `httptest` listener and panicked:

  ```text
  panic: httptest: failed to listen on a port: listen tcp6 [::1]:0: socket: operation not permitted
  ```

  Rerunning the same test with local-loopback approval passed. This was an
  execution-sandbox restriction, not a product failure.

- A command-package test initially could not read the shared Go build cache:

  ```text
  open /home/manuel/.cache/go-build/33/33de0f184c6fa8c36baf3c92e8f15d7b9ff897c4d25db52d432823b558991b39-d: read-only file system
  ```

  The approved rerun passed.

- `go test ./... -count=1` passed all product packages but failed because two
  historical ticket script directories contain several independent `main`
  probes in one Go package:

  ```text
  02-limit-node-shape-probe.go:9:6: main redeclared in this block
  03-limit-option-probe.go:9:6: main redeclared in this block
  02-list-profile-models.go:11:6: main redeclared in this block
  03-filter-deployment-profiles.go:20:6: main redeclared in this block
  ```

  The repository hook correctly scopes tests to `./cmd/... ./internal/...`,
  which passed. I did not alter unrelated historical probes.

- The first commit command emitted only the Lefthook banner and returned with
  all files still staged and `HEAD` unchanged. The second attempt exposed the
  actual concurrent-process failure:

  ```text
  Error: parallel golangci-lint is running
  The command is terminated due to an error: parallel golangci-lint is running
  make: *** [Makefile:99: lint] Error 3
  ```

  The test hook passed, but the lint hook returned exit status 2. After two
  attempts I stopped, did not use `--no-verify`, and left the exact staged
  checkpoint intact.

### What I learned

- The smallest safe authorization repair does not require changing Ragkit:
  trusted local indexes may over-fetch, while GEC filters collapsed chunk IDs
  against canonical document metadata before fusion or external hydration.
- Authorized hits must be renumbered before RRF; preserving their old ranks
  would retain the influence of removed candidates.
- A reranker-boundary recheck is worthwhile defense in depth even when both
  channels have already been filtered.
- “Faithfulness 1.0 for no statements” is only a reporting convention. It must
  never suppress a failed relevance/abstention judge call, and it should not be
  mixed into earned-faithfulness means.
- The repository's `go test ./...` scope is knowingly broader than its valid
  product package graph because ticket probe scripts are not standalone Go
  packages. The hook's `./cmd/... ./internal/...` scope is the meaningful
  whole-product check here.

### What was tricky to build

Authorization has two distinct boundaries. The lexical/vector indexes are
trusted local code, but their ranked outputs contain chunk identities that must
be filtered before fusion; the reranker is potentially remote and must never
receive disallowed hydrated text. The implementation therefore filters each
collapsed channel, renumbers survivors, fuses only authorized hits, and checks
again immediately before hydration. This fixes confidentiality and fusion
attribution without claiming that finite over-fetch gives partition-equivalent
recall.

Judge accounting required separating populations rather than adding one error
counter. Selected, attempted, answered, and successfully judged are different
denominators; timeouts are a subset of answer errors; policy skips are selected
but not attempted; abstentions are judged but excluded from earned
faithfulness. Encoding those distinctions in one outcome projection prevents a
candidate from improving by failure or vacuous truth.

### What warrants a second pair of eyes

- Confirm that filtering after trusted local top-k retrieval but before fusion
  is the accepted v1 boundary; scope partitions may later improve recall but
  are not needed to prevent external disclosure.
- Review the change from all-judged faithfulness to non-abstained earned
  faithfulness and ensure the future GEC gate binds the explicit denominator.
- Check whether context deadline/cancellation errors should also be classified
  as timeouts in addition to `ErrAnswerTimeout`.
- Review the strict judge requirement for explicit `abstained` and `supported`
  fields against cached historical judge responses before the first frozen run.
- Retry the commit only after the concurrent GolangCI-Lint process exits; do
  not bypass Lefthook.

### What should be done in the future

- Complete the staged commit after the external lint process clears.
- Rebase or merge the isolated branch only with explicit coordination around
  the colleague's active main worktree.
- Inventory cached judge responses for strict-schema compatibility; a new
  prompt version may be required rather than reusing malformed cache entries.
- Freeze one small human-authored GEC candidate only after this code checkpoint
  is committed and its semantic inputs are enumerated.
- Do not run a provider evaluation until the candidate, suite, policy, budget,
  and native artifact contract are frozen.

### Code review instructions

- Start at `/tmp/gec-ragopt-phase5/internal/knowledge/service.go` in `Search`,
  `retrieve`, `authorizeHits`, `authorizeFused`, and `rerank`.
- Read `TestRerankNeverReceivesUnauthorizedCandidates` before reviewing the
  implementation.
- Review `/tmp/gec-ragopt-phase5/internal/knowledge/judge.go` from
  `AnswerClient.Answer` through `JudgeAnswer`.
- Review `judgeOutcome` and `emitJudgeSummaries` in
  `/tmp/gec-ragopt-phase5/cmd/coinvault/cmds/knowledge.go` beside
  `TestEmitJudgeSummariesAccountsForFailuresAndDenominators`.
- Validate with:

  ```bash
  go test ./internal/knowledge ./cmd/coinvault/cmds -count=1
  go test ./cmd/... ./internal/... -count=1
  go build ./cmd/coinvault
  ```

### Technical details

```text
GEC worktree:       /tmp/gec-ragopt-phase5
GEC branch:         codex/ragopt-phase5-gec
Base commit:        acc53c5f75a9b14e2b593bfd2ffd2836ac0abadf
Code commit:        pending — hook collision, files remain staged
Main worktree:      untouched; colleague changes preserved
Affected tests:     pass
CoinVault build:    pass
Hook product tests: pass
Hook lint:          blocked by concurrent golangci-lint
Provider calls:     0
RAGOPT next gate:   commit P0 repair, then freeze GEC candidate
```

## Step 25: Verify the landed GEC P0 checkpoint and historical judge cache

On resume, the isolated branch was no longer staged: commit `69c1391` existed,
the worktree was clean, and its reflog recorded the commit at 16:14. The second
commit process had continued through Lefthook after the earlier client session
returned, so Step 24 accurately records the observable state at stop time but
not the eventual repository state. I verified the landed object rather than
assuming either account was authoritative.

I then audited the durable judge cache against the stricter response contract.
All 242 cache entries parse; 122 are statement-extraction responses and 120 are
verdict responses. Every verdict response has ordered one-based statement
indices, explicit Boolean `supported` fields, numeric relevance in `[0,1]`, and
an explicit Boolean abstention field. No cache version bump or cache purge is
needed before the first frozen GEC proof. Task 69 is now complete.

### Prompt Context

**User prompt (verbatim):** (see Step 24)

**Assistant interpretation:** Resume the interrupted checkpoint from actual Git
and cache state, complete verification, and advance RAGOPT only when durable
evidence proves the P0 prerequisite.

**Inferred user intent:** Avoid both false blockage and false completion: trust
the repository and artifacts, not a stale terminal impression, while preserving
the strict diary and commit discipline.

**Commit (code):** `69c1391bfcdcfa8012c1f779482d0fae1b379e92` — "Fix retrieval authorization and judge accounting"

### What I did

- Re-read the diary and Almanach skills for this continuation.
- Inspected exact process names, isolated branch status, main worktree status,
  RAGOPT status, commit log, and reflog.
- Verified `69c1391` contains exactly eight intended knowledge/command files,
  491 insertions, and 54 deletions.
- Ran `git show --check --oneline 69c1391`; it passed.
- Confirmed `/tmp/gec-ragopt-phase5` is clean and the colleague-owned main
  worktree remains untouched.
- Re-ran the repository's meaningful full product scope:

  ```bash
  go test ./cmd/... ./internal/... -count=1
  go build ./cmd/coinvault
  ```

- Scanned all 242 files under the existing judge cache without printing prompt,
  answer, or evidence content.
- Parsed each cached string exactly or through the judge's outer-brace tolerance
  and classified statement versus verdict responses.
- Validated every verdict's indices, required types, and relevance range.
- Checked RAGOPT task 69 and updated the ticket entry point and changelog.

### Why

- Git commit state can change after a detached client session stops reading a
  still-running hook. The commit object, reflog, and clean index are the
  authoritative evidence.
- Tightening parser semantics without auditing the durable cache could turn
  old cache hits into systematic judge errors during the proof.
- Reusing valid cached generations is desirable: it saves provider work while
  the content-addressed key continues to bind prompt version, model, and prompt.
- Task 69 should close only after both the code and its historical operational
  inputs satisfy the new contract.

### What worked

- Commit `69c1391` exists on `codex/ragopt-phase5-gec`; the isolated worktree is
  clean.
- No exact `golangci-lint` or `lefthook` process remained.
- All `./cmd/...` and `./internal/...` packages passed.
- CoinVault compiled successfully.
- Cache audit result:

  ```text
  total:             242
  parsed:            242
  unparsed:          0
  statements:        122
  verdicts:          120
  valid verdicts:    120
  invalid verdicts:  0
  other:             0
  ```

- Existing judge prompt version `v1` and cache content remain semantically
  compatible with the strict decoder, so no provider regeneration is required.

### What didn't work

- The first process search used `pgrep -af 'golangci-lint|lefthook'`. Because
  the pattern itself appeared in the current shell command line, it returned
  the sandbox wrapper and its own zsh process, producing a very large false
  positive. I corrected the check to exact executable names:

  ```bash
  pgrep -ax golangci-lint
  pgrep -ax lefthook
  ```

  Both returned no process.

- Step 24's closing state said the commit was pending. That was true when the
  client stopped observing the second hook, but the underlying session later
  completed and created the commit. I did not rewrite Step 24; this step records
  the state transition and authoritative evidence.

### What I learned

- A unified execution wrapper can return control while a PTY-backed Git hook is
  still active. After any ambiguous commit result, future continuations must
  inspect `HEAD`, the index, and reflog before retrying.
- Strict validation did not invalidate the historical GEC judge cache. The
  earlier model responses were better formed than the old decoder required.
- Cache reuse and evaluation correctness are compatible when validation occurs
  after cache load and malformed cached values remain explicit judge errors.

### What was tricky to build

The operational ambiguity came from two layers of execution: the commit invoked
Lefthook, Lefthook launched Dagger generation and GolangCI-Lint, and the client
session stopped observing before the Git process reached its terminal state.
The only safe reconciliation was evidence triangulation: commit object, reflog,
clean worktree, staged diff absence, and rerun tests.

The cache audit needed to verify structure without leaking cached questions or
retrieved evidence into logs. The projection therefore emitted only booleans
and class counts, never raw `value` strings, while accepting the same exact-JSON
or outer-brace parsing forms used by production.

### What warrants a second pair of eyes

- Review whether the proof should deliberately use the validated existing
  cache or a fresh empty cache. Reusing it reduces spend; a fresh run better
  exercises provider behavior. Either choice must be frozen as run identity.
- Confirm the GEC candidate's judge policy binds earned-faithfulness denominator
  and completion rates explicitly.
- Coordinate how `69c1391` will enter the colleague's evolving main branch
  before running any production-shaped server from that branch.

### What should be done in the future

- Freeze one human-authored GEC candidate and a small feedback suite before
  writing the adapter.
- Declare cache identity and whether cache reuse is part of the proof protocol.
- Preserve the current provider budget ceiling and feedback-before-validation
  gate.
- Check task 68 only after the GEC adapter actually completes the full path;
  task 70 remains the two-run proof requirement.

### Code review instructions

- Inspect `git show 69c1391` from `/tmp/gec-ragopt-phase5`.
- Run the exact product test/build commands above.
- Recompute cache class counts without emitting `.value`.
- Confirm RAGOPT task 69 is checked while tasks 68, 70, and 71 remain open.

### Technical details

```text
GEC commit:          69c1391bfcdcfa8012c1f779482d0fae1b379e92
Branch:              codex/ragopt-phase5-gec
Worktree:            /tmp/gec-ragopt-phase5 (clean)
Product tests:       pass
CoinVault build:     pass
Judge cache:         242/242 parse; 120/120 verdicts strict-valid
Task closed:         RAGOPT-001 task 69
Provider calls:      0
Next gate:           freeze candidate, suite, policy, cache, and budgets
```

## Step 26: Select the minimum GEC candidate and identify its injection seam

The first GEC proof will test one human-authored text mutation: the
model-facing `knowledge_search` description. The challenger will add explicit
routing guidance for `schema_doc`, `product`, and `guide`; the incumbent will
preserve the current description. This is the smallest candidate recommended
by the optimizer handoff, and it exercises real live-session tool selection
without changing retrieval, authorization, the corpus, or the judge.

The investigation also found the exact missing product seam. CoinVault builds
the description inside `knowledge.NewToolEntry`, captures the resulting spec in
its registrar closure, and registers that entry once at server startup. A
frozen candidate therefore cannot be evaluated honestly until the product can
load one explicit description asset at composition time. The seam belongs in
GEC, not in `ragopt`; `ragopt` should only validate the one-mutation bundle and
run the paired cells.

### Prompt Context

**User prompt (verbatim):** (see Step 24)

**Assistant interpretation:** Continue the first real GEC integration, keep it
small and evidence-gated, and record/print every material checkpoint.

**Inferred user intent:** Prove that the reusable harness prevents fake
self-optimization by freezing one real product mutation and every confounder
before implementation or provider spending.

### What I did

- Verified both the isolated GEC worktree and RAGOPT worktree are clean.
- Read `internal/knowledge/tool.go`, the catalog tool-spec renderer and
  registrar, live server composition, application profiles, answer acquisition,
  and the decomposed judge path.
- Enumerated all 60 existing knowledge-eval questions and their seven strata.
- Re-read the implemented RAG-TTC I5 candidate bundle and consumer adapter to
  preserve the proven custody shape without copying TTC runtime code.
- Selected a source-role-routing tool-description candidate rather than a
  retrieval-algorithm or synonym mutation.
- Confirmed the current answer record captures evidence and tool-result counts,
  but does not yet project tool-call input arguments such as `source_roles`.
- Made no file, serving, or provider changes in GEC during this investigation.

### Why

- One text mutation gives RAGOPT an honest parent/challenger comparison without
  entangling the colleague-owned Admin Chat transport and frontend work.
- Source-role routing is already part of the production tool schema and search
  service. The candidate changes guidance, not capability.
- Freezing the candidate before adding its injection seam prevents the adapter
  implementation from silently defining the experiment after seeing results.
- The existing feedback-before-validation rule prevents unnecessary provider
  spending if the cheap candidate has no measurable benefit.

### What worked

- The runtime boundary is clear:

  ```text
  frozen description asset
          |
          v
  knowledge.NewToolEntry
          |
          v
  ToolCatalog registrar -> per-session tool registry -> live tool loop
          |
          v
  timeline -> AnswerRecord -> decomposed judge -> native outcome -> ragopt
  ```

- `SearchInput.SourceRoles` and `Service.Search` already support the candidate's
  intended behavior; no retrieval feature is missing.
- The existing eval set has schema, product/facet, guide, paraphrase, and
  multi-document cases from which disjoint feedback and validation suites can
  be selected.
- The proven TTC bundle establishes the required candidate, snapshot, suite,
  gate-policy, runtime-contract, and native-artifact custody pattern.

### What didn't work

- The first broad source search included a nonexistent `configs` path in GEC
  and printed `rg: configs: No such file or directory`. This was only a search
  target mistake; GEC has no RAGOPT config tree yet.
- One inspection command requested `internal/knowledge/answer.go`, but answer
  acquisition actually lives in `internal/knowledge/judge.go`; `sed` reported:

  ```text
  sed: can't read internal/knowledge/answer.go: No such file or directory
  ```

  I corrected the inspection path and read the complete live answer projection.
- The current `AnswerRecord` is insufficient to prove that improved answer
  scores came from better `source_roles` choices. It counts tool results but
  discards tool-call input payloads. The native artifact must retain that
  product evidence; the shared RAGOPT outcome does not need a new generic field.

### What I learned

- The description shown in catalog health metadata and the description passed
  to the model come from the same captured `ToolSpec`; an override must affect
  both or identity reporting becomes misleading.
- The description includes bundle-derived document count and available roles.
  The frozen asset should be a complete product-owned description specification,
  while the bundle digest and roles are locked dimensions.
- A server-startup injection seam is sufficient for this proof. A dynamic
  per-request prompt mutation API would expand scope and is unnecessary.
- Source-role selection belongs in the native GEC artifact because it is a
  product diagnostic, while answer relevance, faithfulness, completion, calls,
  tokens, and duration remain the common comparison projection.

### What was tricky to build

The candidate is text, but the existing text is assembled dynamically from a
structured prompt spec and bundle metadata. Treating the final rendered string
as an opaque flag would be easy but brittle; treating all prompt-pack machinery
as a new generic RAGOPT feature would be overengineering. The narrow solution
is a GEC-owned, strict description asset loaded once at startup and passed into
the existing spec/registrar path. Parent and challenger assets remain complete
replacements, and RAGOPT sees exactly one changed mutable asset.

The other sharp edge is causal evidence. Answer relevance can be the declared
promotion target, but the experiment should still record actual tool arguments
so a reviewer can distinguish successful role routing from score noise. That
projection must be added to the native artifact without changing shared gates
or the production timeline contract.

### What warrants a second pair of eyes

- Confirm that a complete structured tool-description asset is preferable to a
  raw rendered string for the GEC product contract.
- Review the exact three feedback cases and disjoint validation cases before
  provider execution; each should have an unambiguous expected source role.
- Confirm that the native artifact may record tool-call arguments from the
  existing timeline without exposing reasoning or unauthorized evidence.
- Review whether the incumbent asset must reproduce the current dynamic
  description byte-for-byte for the locked bundle, including document count
  and sorted available roles.

### What should be done in the future

- Write and validate the parent/challenger description assets, disjoint suites,
  gate policy, runtime contract, answer schema, and judge contract.
- Add only the startup-time GEC description injection and native tool-input
  projection required to execute those frozen assets.
- Add the GEC consumer arm after the candidate bundle validates with RAGOPT.
- Run feedback twice from fresh run roots; run validation only if feedback gates
  pass.

### Code review instructions

- Start with `/tmp/gec-ragopt-phase5/internal/knowledge/tool.go` at
  `NewToolEntry` and `toolSpec`.
- Follow registration through
  `/tmp/gec-ragopt-phase5/internal/coinvaulttools/catalog/registry.go` and server
  startup in `/tmp/gec-ragopt-phase5/internal/webchat/server/server.go`.
- Read `AnswerClient.extract` in
  `/tmp/gec-ragopt-phase5/internal/knowledge/judge.go` to see the current native
  projection limitation.
- Compare the intended custody layout with
  `configs/ragopt/i5-combined-comparison-v1` in the RAG-TTC proof worktree.

### Technical details

```text
Candidate track:     C / source-role routing
Mutable assets:      1 (knowledge_search description)
Parent behavior:     current production description
Challenger behavior: explicit schema_doc/product/guide routing
Locked:              corpus, bundle, index, retrieval, scopes, models, judge,
                     suites, budgets, answer schema, runtime, adapter sources
Promotion target:    answer_relevance (provisional until policy freeze)
Native diagnostics:  tool inputs, selected source_roles, evidence roles
Provider calls:      0
Next gate:           commit a RAGOPT-valid frozen candidate bundle
```

## Step 27: Commit the frozen GEC source-role candidate contract

The GEC candidate now exists as a complete RAGOPT bundle and validates as
exactly one mutable asset. Its parent is the current production
`knowledge_search` description; its challenger adds explicit routing for
`schema_doc`, `product`, and `guide`. All other declared assets and dimensions
are byte-identical between arms.

This is a contract freeze, not a completed experiment. The source lock openly
records that the description loader and consumer adapter do not exist yet and
must be digest-pinned before the first provider cell. That makes the remaining
implementation work visible without pretending that a validated YAML bundle
is an exercised optimizer integration.

### Prompt Context

**User prompt (verbatim):** (see Step 24)

**Assistant interpretation:** Freeze and commit the smallest credible GEC
candidate before implementing its adapter, while keeping budgets, identities,
and the feedback/validation spending boundary explicit.

**Inferred user intent:** Ensure the product integration cannot move its target
or hide confounders after seeing model results, and preserve a reviewable audit
trail from hypothesis through implementation and execution.

**Commit (code):** `7599193` — "Add frozen GEC source-role candidate contract"

### What I did

- Added `configs/ragopt/source-role-routing-v1` in the isolated GEC worktree.
- Authored complete parent and challenger `knowledge_search` description assets.
- Selected three feedback cases and three disjoint validation cases spanning
  `schema_doc`, `product`, and `guide` routing.
- Added gate, runtime, judge, native-answer, adapter, bundle, and source-lock
  contracts.
- Pinned the real hybrid bundle identity:

  ```text
  bundle:           rk-55be57b45fc6d624d0341c8ec6965f49
  documents:        16032
  corpus digest:    5569c5db5118534cfb3f26113fc09d42efdaaba688115754b24685e93233fd2b
  vector index:     ca92a5bc4435a68e84e09ddd6fd62a1173dfad133506c864aeb238fc2bdfb29d
  lexical manifest: d28e1f9b7650ba7042f9f3f04ef8c00d32be16cc6f76e42742406fdedb5ee0ea
  ```

- Pinned `gpt-5.6-luna-low` answering, `gpt-5.6-luna` judging, production
  `analyst-rag`, disabled reranking and synonyms, and content-addressed judge
  cache reuse.
- Set the proof-run ceiling to six answer sessions, twelve judge generations,
  zero embedding generations, fewer than 500,000 provider tokens, and at most
  one US dollar.
- Validated unique suite IDs and feedback/validation disjointness with `jq`.
- Used RAGOPT's real CLI to derive the snapshot IDs and validate the bundle.
- Ran `git diff --check` and committed all 14 candidate files.

### Why

- The exact suite, target, regressions, and budgets must exist before an adapter
  can execute them.
- A three-case feedback suite is enough to exercise each role and keeps the
  first proof cheap. A second fresh run establishes repeatability without
  silently expanding the suite.
- Validation remains a separate asset so the spending gate is enforceable and
  reviewable; its existence does not authorize executing it.
- The source lock binds current product semantics and explicitly reserves the
  not-yet-created adapter digests, preventing provisional custody from being
  confused with final run custody.

### What worked

- RAGOPT accepted the bundle and independently proved one changed asset:

  ```text
  candidate:  sha256:e7be1ae545bc826a34fce3e4ca7dd7f1dfb8715a89b5d1b74ec931c771d00820
  parent:     sha256:24145fce23d7ddfcc3e1766fbc730cecbf3c2e4a31f489beb6622120a9900938
  challenger: sha256:2d14a0978c0e5c354b3b4b0b2db5965ff67d8a55d592d7d85f4726d837583b78
  mutation:   knowledge_search_description
  valid:      true
  ```

- Both JSON suites parsed, contained unique IDs, and had no ID overlap.
- The parent description reproduces the current dynamic spec for the pinned
  bundle: 16,032 documents and sorted roles `guide, product, schema_doc`.
- The pre-commit hook correctly skipped Go/web jobs because this checkpoint
  contains only configuration assets.
- No provider, embedding, answer, or judge call occurred.

### What didn't work

- The first RAGOPT validation intentionally used zero snapshot IDs so the
  validator would report its computed semantic identities. It failed first on
  the parent and then on the challenger:

  ```text
  snapshot ID mismatch: manifest=sha256:0000... actual=sha256:24145fce...
  snapshot ID mismatch: manifest=sha256:0000... actual=sha256:2d14a097...
  ```

  I inserted each computed ID and reran validation successfully. These were
  controlled identity derivation failures, not ignored validation errors.
- A broad search for a tool-loop ceiling produced a truncated module-cache
  result. It established that Pinocchio has a bounded loop, but the GEC runtime
  does not currently expose that value as a product setting. I therefore did
  not invent a per-cell loop override in the frozen candidate.

### What I learned

- The current product loop boundary is less explicit than TTC's tool-QA runner.
  The adapter must observe and account for production calls rather than claim a
  configurable per-cell maximum that GEC does not expose.
- Bundle identity needs more than the content-addressed bundle ID in this first
  proof: explicit lexical-manifest and vector-SQLite digests make the physical
  retrieval artifacts auditable.
- The candidate's schema case is intentionally adversarial because `sql_doc`
  and `sql_query` remain available. If the model bypasses `knowledge_search`,
  that is a legitimate result about this description candidate, not an adapter
  failure.
- A valid candidate bundle proves custody shape only. “Implemented and
  exercised” still requires a product arm, native artifacts, paired cells,
  comparison, and a repeated decision.

### What was tricky to build

The hardest boundary was distinguishing a truly locked runtime input from a
future implementation file. Existing prompts, product sources, corpus, index,
models, and policies can be pinned now. The adapter and description loader
cannot be honestly hashed before they exist. Instead of using placeholders
that look final, `source-lock.yaml` states this incompleteness in machine-readable
prose and requires those digests plus recomputed snapshot identities before any
provider-backed execution.

The budgets also use two different units. Six answer sessions is the number of
RAGOPT cells across two arms and three cases; it is not necessarily six model
generations because the production tool loop may call the answer provider more
than once. The global token and dollar ceilings therefore remain separate from
session and judge-call counts. The adapter must stop scheduling cells when a
hard ceiling is reached and record a failure rather than drop the cell.

### What warrants a second pair of eyes

- Review the challenger wording for accidental preference toward SQL tools or
  role overrestriction.
- Review whether `orders-table` is the right schema feedback case given the
  competing production SQL tools; changing it after adapter execution would
  create a new candidate/suite identity.
- Review the provisional gate thresholds, especially the zero minimum target
  delta and per-case answer-relevance regression floor.
- Confirm that the one-dollar/500k-token ceiling can be enforced or safely
  preflighted with the current provider usage telemetry.
- Confirm that copying the bundle manifest plus physical index digests is
  sufficient custody without copying the 368 MB vector database into the
  candidate directory.

### What should be done in the future

- Implement strict loading/rendering of `gec-knowledge-search-description/v1`
  and prove the parent is byte-identical to the existing model description.
- Project `knowledge_search` call inputs and evidence roles into the native GEC
  artifact without changing shared RAGOPT APIs.
- Implement the GEC consumer arm and budget accounting.
- Add implementation source digests to `source-lock.yaml`, update both snapshot
  identities, and revalidate before provider execution.
- Run feedback twice from fresh roots. Do not run validation unless the
  feedback gate passes.

### Code review instructions

- Start at `/tmp/gec-ragopt-phase5/configs/ragopt/source-role-routing-v1/candidate.yaml`.
- Diff the two files named `knowledge-search-description.yaml`; they must be the
  only mutable byte difference.
- Review `shared/runtime-contract.yaml`, `shared/source-lock.yaml`, and
  `shared/adapter-contract.yaml` together.
- Validate with:

  ```bash
  cd /home/manuel/code/wesen/go-go-golems/ragopt
  go run ./cmd/ragopt candidate validate \
    --bundle /tmp/gec-ragopt-phase5/configs/ragopt/source-role-routing-v1 \
    --format yaml
  ```

- Confirm `git show --check 7599193` in the isolated GEC worktree.

### Technical details

```text
GEC commit:          7599193
Bundle files:        14
Feedback cells:      3 cases x 2 arms x 1 repeat = 6
Validation cells:    3 cases x 2 arms x 1 repeat = 6 (not authorized yet)
Answer sessions:     <= 6 per feedback proof run
Judge generations:   <= 12 per feedback proof run
Embedding calls:     0
Provider ceiling:    < 500,000 tokens and <= $1 per proof run
Candidate valid:     yes, provisional source-lock caveat explicit
Provider calls:      0
Next gate:           implement product-owned loader and native adapter
```

## Step 28: Trace and explicitly lock the shared Geppetto tool loop

The user's objection exposed an important wording error in Step 27. TTC and
GEC do not use unrelated tool-loop implementations: both execute
`github.com/go-go-golems/geppetto/pkg/inference/toolloop`. The real difference
is configuration ownership. TTC supplies a deliberate bounded policy; GEC
constructs the same Geppetto engine builder without loop/tool configuration,
so it inherits Geppetto v0.13.7 defaults.

I corrected the frozen GEC runtime identity to record those exact defaults.
This keeps the first proof production-faithful while making the mismatch
visible. I did not force TTC's `4 / 1 / no retry / reserve final` evaluation
policy into GEC: doing so would test a newly invented GEC runtime rather than
the Admin Chat path we are supposed to optimize.

### Prompt Context

**User prompt (verbatim):** "3876 +- The current product loop boundary is less explicit than TTC's tool-QA runner.
    3877 +  The adapter must observe and account for production calls rather than claim a
    3878 +  configurable per-cell maximum that GEC does not expose.

---

Aren't we unifying things? Why are the tool loops different? Doesn't ragkit provide them?"

**Assistant interpretation:** Verify whether TTC and GEC actually use different
loop engines, correct the architecture explanation, and avoid preserving an
accidental difference under the label of unification.

**Inferred user intent:** Reuse proven shared machinery coherently: do not add a
second framework, do not blame product boundaries for accidental divergence,
and make the remaining policy difference explicit and intentional.

**Commit (configuration):** `401f4e5` — "Lock explicit Geppetto tool loop semantics"

### What I did

- Traced TTC from `pkg/rag/toolanswer/service.go` into Geppetto
  `toolloop.New`, `LoopConfig`, and `tools.ToolConfig`.
- Traced GEC from `CoinVaultRuntimeResolver` through `RuntimeComposer`,
  Pinocchio `chatapp`, and Geppetto `enginebuilder.Builder`.
- Verified that GEC's composed engine internally owns the registry and loop;
  Pinocchio receives it as a single `ComposedRuntime.Engine` and does not add a
  second loop.
- Read Geppetto v0.13.7's authoritative default constructors.
- Updated the frozen runtime contract with exact implementation, version,
  iterations, parallelism, timeout, error handling, retry behavior, and lack of
  a reserved final call.
- Replaced the vague snapshot dimension
  `production-sessionstream-default` with a semantic value that states the
  actual policy.
- Recomputed both snapshot IDs and revalidated the one-mutation candidate.
- Committed the contract correction before resuming adapter implementation.

### Why

- `ragkit` is the shared retrieval layer; it does not orchestrate model/tool
  conversations.
- Geppetto is the shared inference and tool-loop layer. Pinocchio/sessionstream
  owns chat lifecycle and transport around it.
- RAGOPT should compare and custody product executions, not grow a competing
  tool-loop implementation.
- A runtime default is still semantic behavior. If it affects provider calls,
  concurrency, retries, cost, or completion, it must be explicit in an
  experiment snapshot.

### What worked

- The shared substrate is now proven from source:

  ```text
  TTC Tool QA ----------------------+
                                     +--> Geppetto toolloop
  GEC RuntimeComposer -> enginebuilder+

  ragkit: retrieval/index/fusion/rerank primitives
  Pinocchio/sessionstream: chat lifecycle and persistence
  ragopt: candidate/run/pair/gate custody
  ```

- GEC's actual inherited policy is now frozen:

  ```text
  max iterations:      20
  tool choice:         auto
  execution timeout:   30s
  max parallel tools:  3
  tool errors:         continue
  retries:             2 (1s base, factor 2)
  reserve final call:  false
  ```

- The corrected candidate remains exactly one mutable text asset:

  ```text
  candidate:  sha256:fc715984f0be41c2fbdb41a6446707cb00ce3f21ff39c8fc0adfc727488b4c13
  parent:     sha256:9f7197c9b85ff75cfb2476a344bf77b7a46e712d14b8905b5fbdf9bd1e18c9a2
  challenger: sha256:eeea2419b34e701ecd03b3b443726f607b8625402ff38606cbf294bd442997ad
  ```

### What didn't work

- An earlier attempted implementation patch tried to add the description
  loader, tool override, tests, and local-runner wiring at once. `apply_patch`
  rejected the entire patch because `internal/knowledge/tool_test.go` did not
  have the assumed imports:

  ```text
  apply_patch verification failed: Failed to find expected lines in
  /tmp/gec-ragopt-phase5/internal/knowledge/tool_test.go:
  import (
      "context"
      "strings"
      "testing"
  ```

  No part of that patch was applied. The failure was useful: it left the code
  untouched while the architectural objection was resolved.
- One source inspection looked for `internal/webchat/runtime_composer.go`; the
  implementation is `internal/webchat/runtime.go`.
- One Geppetto inspection looked under `pkg/inference/enginebuilder`; the
  package actually lives under `pkg/inference/toolloop/enginebuilder`.

### What I learned

- GEC does not have an unbounded loop. It has a 20-iteration Geppetto default,
  but the product does not name or expose that policy.
- TTC's provider-budget wrapper is additional policy around the same loop. It
  counts provider calls and reserves a final response call; Geppetto's generic
  default does neither.
- “Unification” has several layers. We already share the execution engine; we
  are extracting shared evaluation custody into RAGOPT. Policy values may still
  differ by product, but they must use the same vocabulary and be locked.
- Moving tool-loop ownership into `ragkit` would blur retrieval and inference
  responsibilities. If a reusable named runtime policy is later extracted, its
  natural owner is Geppetto, not ragkit or ragopt.

### What was tricky to build

The GEC call graph looks like a single-pass Pinocchio inference at first glance
because `ComposedRuntime.Registry` is nil. The loop was not absent: GEC had
already built an `enginebuilder` runner whose internal registry is non-nil.
Pinocchio then invokes that engine once, and the engine performs the Geppetto
loop internally. Understanding this nested composition was necessary to avoid
adding a second tool loop or double-counting iterations.

The policy question is separate from the implementation question. Using TTC's
smaller limits would make the policies numerically uniform but the experiment
less representative. The pragmatic first proof locks GEC's current behavior;
the cross-product follow-up should define a named shared policy contract and
then deliberately migrate products, with separate evaluation of that runtime
change.

### What warrants a second pair of eyes

- Confirm that the first GEC proof should remain production-faithful at
  `20 / 3 / retry 2`, rather than combining a loop-policy migration with the
  source-role description mutation.
- Review whether the Admin Chat production policy itself should later adopt a
  reserved-final-call budget wrapper like TTC's.
- Confirm that provider retry calls can be observed and charged accurately from
  existing inference events; if not, the adapter must report accounting as
  incomplete rather than fabricate a count.
- Review whether the reusable named loop-policy schema belongs upstream in
  Geppetto after the two product integrations expose concrete requirements.

### What should be done in the future

- Keep the current GEC candidate to one description mutation.
- Implement the description injection through the existing Geppetto-backed GEC
  composition, not a new evaluator-only loop.
- Record actual inference/tool calls in native evidence and compare them to the
  locked policy.
- After Phase 5, propose one small Geppetto-owned named policy/config surface if
  TTC and GEC adapters demonstrate the same requirements.
- Do not move chat-loop execution into `ragkit` or `ragopt`.

### Code review instructions

- In TTC, start at
  `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/rag-ttc/pkg/rag/toolanswer/service.go:130`.
- In GEC, follow
  `/tmp/gec-ragopt-phase5/internal/webchat/sessionstream/sessionstream_runtime_resolver.go`
  into `/tmp/gec-ragopt-phase5/internal/webchat/runtime.go`.
- Read Geppetto v0.13.7
  `pkg/inference/toolloop/enginebuilder/builder.go`,
  `pkg/inference/toolloop/config.go`, and `pkg/inference/tools/config.go`.
- Validate the corrected candidate with the RAGOPT command from Step 27.
- Inspect `git show --check 401f4e5` in the isolated GEC worktree.

### Technical details

```text
Shared engine:       Geppetto toolloop v0.13.7
TTC policy:          explicit 4 / parallel 1 / no retry / reserve final
GEC policy:          explicit identity of inherited 20 / parallel 3 / retry 2
ragkit role:         retrieval primitives, not chat orchestration
ragopt role:         immutable evaluation custody, not loop execution
Config commit:       401f4e5
Provider calls:      0
Next gate:           inject candidate text into the existing GEC composition
```

## Step 29: Add the strict description seam and verify the asynchronous commit

The first executable GEC candidate seam is now committed. A strict
`gec-knowledge-search-description/v1` loader supplies a complete replacement
description to the existing `knowledge_search` tool spec, and the no-HTTP local
CoinVault runner can open the real bundle and register that description in the
same Geppetto-backed runtime composition used by chat.

The commit process again outlived the execution wrapper. Two commit clients
stopped observing after Lefthook began, so I stopped after the second attempt as
required. On continuation, Git proved that commit `b2256f2` had landed at
16:43:13 with a clean worktree. I then ran the complete lint hook separately in
tmux and captured its real `__EXIT__:0` result.

### Prompt Context

**User prompt (verbatim):** (see Step 28)

**Assistant interpretation:** Continue RAGOPT-001 after correcting the shared
loop architecture; implement only the narrow product-owned candidate injection
seam, maintain the diary, and preserve reliable commit evidence.

**Inferred user intent:** Turn the frozen text candidate into real product
behavior without introducing a second loop or an evaluator-only fake, while
keeping failures and operational ambiguity fully auditable.

**Commit (code):** `b2256f2` — "Add evaluated knowledge description seam"

### What I did

- Added `internal/knowledge/description.go` with:
  - exact API-version validation;
  - YAML `KnownFields(true)`;
  - one-document enforcement;
  - required summary validation;
  - rejection of empty note and starter-query elements;
  - conversion to the existing promptpack description shape.
- Added an optional `Description *ToolDescription` to `knowledge.ToolConfig`.
- Kept `NewToolEntry` on the same `catalog.ToolSpec` and registrar path so the
  catalog description and model-visible description stay identical.
- Added focused strict-load and rendered-description tests.
- Extended `localwebchat.RunnerOptions` with a bundle path and description path.
- Made `NewLocalRunner` open and own the knowledge service, load the description
  once, and register `knowledge_search` beside the existing SQL tools.
- Deliberately did not enable the environment-driven reranker or synonym paths
  in the proof runner; the frozen runtime contract says both are disabled.
- Formatted the changed Go files.
- Ran focused tests for `internal/knowledge` and
  `internal/webchat/localwebchat`.
- Ran the full repository `make lint` in tmux and captured exit zero.
- Reconciled the commit through `HEAD`, reflog, clean index, and
  `git show --check` rather than issuing another commit.

### Why

- The candidate must enter the real product tool registry before it can affect
  the real Geppetto loop.
- Strict parsing ensures unknown fields or multiple YAML documents cannot
  silently change evaluated meaning.
- A complete replacement asset preserves RAGOPT's one-mutation model; no patch
  application or merge semantics are needed.
- The local runner uses CoinVault's existing application profiles, runtime
  resolver, projection runtime, sessionstream Hub, chatapp service, and
  Geppetto engine. It removes only HTTP/websocket transport from evaluation,
  not product inference behavior.

### What worked

- Focused test result:

  ```text
  ok github.com/go-go-golems/gec-rag/internal/webchat/localwebchat 0.322s
  ok github.com/go-go-golems/gec-rag/internal/knowledge             1.624s
  ```

- Full lint completed in the tmux session `ragopt-gec-lint`:

  ```text
  golangci-lint: 0 issues
  GOWORK=off go vet -vettool=/tmp/geppetto-lint ./cmd/... ./internal/...
  __EXIT__:0
  ```

- The commit landed with exactly four intended files and the worktree is clean.
- The override changes only description content. Input schema, authorization,
  result budgets, evidence ledger, and retrieval service are shared by both
  arms.
- No provider call was made.

### What didn't work

- The first focused test run inside the restricted sandbox failed because an
  existing `httptest.NewServer` could not bind IPv6 localhost:

  ```text
  panic: httptest: failed to listen on a port:
  listen tcp6 [::1]:0: socket: operation not permitted
  ```

  I reran the same tests with approved local-socket access; they passed.
- Two ordinary `git commit` calls returned after printing only the Lefthook
  header and a session identifier. At that moment the files still appeared
  staged and no exact hook process was visible. Per the repository debugging
  rule, I stopped after the second identical attempt and reported:

  ```text
  I think I'm stuck, let's TOUCH GRASS.
  ```

- `lefthook run pre-commit --verbose` and a direct `make lint` launched Dagger
  work that again outlived the short execution wrapper. Neither output was a
  trustworthy terminal result.
- The first sandboxed tmux attempt failed with:

  ```text
  error connecting to /tmp/tmux-1000/default (Operation not permitted)
  ```

  Running tmux with the required permission succeeded. The captured pane
  eventually showed the real exit marker.
- On reconciliation, `b2256f2` already existed in the reflog at 16:43:13. The
  later commit command correctly replied `nothing to commit`; I did not create
  a duplicate or amend the landed object.

### What I learned

- The recurring commit ambiguity is not a code or hook failure. It is an
  observation-lifetime mismatch: Dagger generation inside `make lint` can
  outlive the execution wrapper while Git continues and eventually commits.
- Long hooks in this repository should be run and captured in tmux before a
  commit. Once they pass, Git state must still be checked because an earlier
  client may already have completed asynchronously.
- The local CoinVault runner was missing knowledge registration even though the
  HTTP server had it. Adding the same optional bundle/tool composition makes it
  suitable for production-shaped local evaluation without changing chat
  session semantics.
- The candidate injection seam belongs to GEC. RAGOPT remains unaware of
  promptpack types and receives only native outcome evidence.

### What was tricky to build

The description appears in two places that must not drift: `ToolEntry.Description`
for catalog metadata and the description passed to Geppetto's tool definition.
The implementation resolves that by replacing the structured spec before both
values are derived; it does not mutate only one string after registration.

Resource ownership also mattered. `knowledge.Open` owns Bleve and vector-index
handles. The local runner adds `knowledgeService.Close` to its existing reverse
cleanup stack before any later dependency can fail, so description-load or
runtime-construction errors do not leak bundle locks.

The operationally tricky part was distinguishing a slow child process from a
failed hook. Exact process searches were empty after the wrapper returned, yet
the reflog later proved Git had committed. Tmux with an explicit `__EXIT__`
marker is the reliable pattern for future long validation commands.

### What warrants a second pair of eyes

- Review whether `ToolDescription.Validate` should require at least one note or
  starter query; v1 currently permits empty lists but rejects empty elements.
- Confirm the local proof runner should intentionally keep reranking and
  synonyms off rather than reuse environment-driven server configuration.
- Review cleanup ordering around the knowledge service, turn store, and local
  runtime.
- Verify the eventual adapter compares the parent description's rendered bytes
  to the default description for the opened 16,032-document bundle before any
  provider call.
- Confirm that using the no-HTTP local runner is production-shaped enough: it
  shares sessionstream/chatapp/runtime composition but not authorization
  middleware or websocket transport, neither of which the candidate mutates.

### What should be done in the future

- Add a preflight that renders parent/default descriptions and rejects drift.
- Capture complete tool-call inputs, tool results, final answer, evidence roles,
  and session identity from local projected events.
- Implement the GEC RAGOPT arm around the local runner and existing decomposed
  judge.
- Add adapter/source digests to the candidate source lock and recompute
  snapshots before provider execution.
- Run long repository hooks in tmux and inspect Git before retrying commits.

### Code review instructions

- Start at `/tmp/gec-ragopt-phase5/internal/knowledge/description.go`.
- Follow `ToolConfig.Description` through `NewToolEntry` and `toolSpec` in
  `/tmp/gec-ragopt-phase5/internal/knowledge/tool.go`.
- Review bundle ownership and registration in
  `/tmp/gec-ragopt-phase5/internal/webchat/localwebchat/local_runner.go`.
- Run:

  ```bash
  go test ./internal/knowledge ./internal/webchat/localwebchat -count=1
  make lint
  git show --check b2256f2
  ```

- Use tmux for `make lint`; wait for an explicit exit marker.

### Technical details

```text
GEC commit:          b2256f2
Files:               4
Description API:     gec-knowledge-search-description/v1
Execution substrate: existing GEC -> Geppetto composition
Focused tests:       pass
Full lint:           pass, __EXIT__:0
Worktree:            clean
Provider calls:      0
Next gate:           native event projection and GEC RAGOPT arm
```

## Step 30: Project real Geppetto provider accounting from the local CoinVault runtime

The GEC proof runner can now observe answer-generation calls at the actual
shared Geppetto/sessionstream boundary. Commit `e0f8f2d` extends the existing
no-HTTP local chat event projection with provider-call identity, terminal
metadata, token usage, duration, and tool-call disposition. This is observation
of the production loop, not a replacement loop or an evaluator-only counter.

### Prompt Context

**User prompt (verbatim):** "Aren't we unifying things? Why are the tool loops
different? Doesn't ragkit provide them?"

**Assistant interpretation:** Preserve the corrected architecture—one shared
Geppetto loop with product-owned policy—then continue RAGOPT-001 by making GEC's
real loop observable to the product adapter.

**Inferred user intent:** Prevent another fake optimization implementation by
measuring the same calls the production chat actually performs, while keeping
retrieval, loop orchestration, chat lifecycle, and evaluation custody in their
proper packages.

**Commit (code):** `e0f8f2d` — "Expose local provider call accounting"

### What I did

- Reconciled both worktrees and the canonical Phase 5 task ledger before
  editing.
- Verified that `ragkit` owns retrieval primitives and that both TTC and GEC
  execute their model/tool cycle through Geppetto's `toolloop`.
- Traced Pinocchio's `runtimeEventSink` and confirmed that it projects
  `EventProviderCallFinished` into `ChatProviderCallFinished` with:
  - provider-call correlation ID;
  - stop reason and finish class;
  - input, output, cached, cache-creation, and cache-read tokens;
  - duration;
  - whether the provider response contained tool calls.
- Extended `localwebchat.LocalChatEvent` to retain that data.
- Added stable kinds for provider-started, metadata-updated, and
  provider-finished events.
- Exposed the new fields in the structured `coinvault chat send` Glazed rows.
- Added a focused test proving exact projection of a provider-finished event.
- Ran the focused local-runner and command tests.
- Let the repository's complete pre-commit lint and test hook finish before
  accepting the commit.
- Kept the pinned `ragopt` module dependency unstaged for the forthcoming
  adapter commit.

### Why

- GEC does not need a second tool loop. It needs evidence from the same loop
  production uses.
- Counting configured maxima would conflate policy with observed work and
  could make a candidate look cheaper or more expensive than it was.
- Provider-finished events are the narrowest authoritative answer-call surface:
  they are emitted by Geppetto, projected by Pinocchio, and delivered through
  the same sessionstream Hub used by the local product runner.
- Common RAGOPT `provider_calls` must represent product execution cost. Judge
  calls are separate evaluation overhead and will remain separately recorded
  in the GEC native artifact and consent ledger.

### What worked

- Focused tests passed:

  ```text
  ok github.com/go-go-golems/gec-rag/internal/webchat/localwebchat 0.137s
  ok github.com/go-go-golems/gec-rag/cmd/coinvault/cmds             0.110s
  ```

- The complete pre-commit hook passed:

  ```text
  golangci-lint: 0 issues
  go vet with geppetto-lint: pass
  go test ./cmd/... ./internal/...: pass
  commit: e0f8f2d
  ```

- The hook tested the new projection in both focused and repository-wide
  contexts.
- Only three intended files entered the commit.
- No provider, embedding, or judge call occurred.

### What didn't work

- The first focused test attempt could not read Go's default build cache under
  the managed filesystem:

  ```text
  open /home/manuel/.cache/go-build/...: read-only file system
  ```

  I created `/tmp/ragopt-gec-gocache` and reran with an explicit task-specific
  `GOCACHE`; the tests passed.
- A fresh-cache test produced no output for several polling intervals. Process
  inspection showed no failure, and the warm-cache rerun eventually returned
  the exact passing result. I did not treat silence as success.
- The commit hook again outlived the initial execution window. I polled the
  existing session rather than starting another commit. It completed normally
  after Dagger web generation, full lint, vet, and tests.

### What I learned

- Pinocchio already carries the provider lifecycle and usage data the adapter
  needs; an engine wrapper would duplicate instrumentation and risk observing a
  different execution path.
- GEC and TTC share loop machinery but intentionally retain distinct product
  policy. RAGOPT compares candidates within one frozen product policy; it does
  not normalize policy by mutating the product.
- The provisional GEC adapter contract contains one stale phrase:
  `provider_calls: answer_plus_uncached_judge_calls`. The proven TTC convention
  and Phase 5 diary require product answer calls only. That contract must be
  corrected before source-lock digests and snapshot identities are finalized.
- Provider cache-token metadata is not the same thing as the local judge cache.
  Both must be named and recorded separately.

### What was tricky to build

The useful accounting event is not generated by `ragkit` because retrieval is
not the model/tool-loop owner. It originates in Geppetto, passes through
Pinocchio's runtime sink, and becomes a sessionstream UI event. The local runner
already consumed that stream but normalized provider events as generic rows.
The change had to preserve the full event payload for existing consumers while
adding typed accounting fields for the adapter.

The second subtlety is cost vocabulary. The answer model's observed calls and
tokens are comparable product behavior. Judge generations are required to
measure quality but are not behavior caused by one production answer. Mixing
them into the common outcome would reward cells whose judging is skipped and
would contradict the corrected TTC proof.

### What warrants a second pair of eyes

- Confirm that summed `ChatProviderCallFinished` usage is the desired product
  token measure when providers report cache-read tokens in addition to input
  tokens.
- Review whether structured `chat send` consumers need an explicit schema note
  for the new optional provider fields.
- Confirm that metadata-updated events should remain observable but only
  finished events should increment answer-call counters.
- Review the adapter-contract correction separating answer cost from judge
  overhead before the candidate is re-identified.

### What should be done in the future

- Implement the GEC RAGOPT arm around `LocalRunner.RunPrompt`.
- Capture final answer, knowledge inputs, tool results, admitted evidence IDs
  and roles, and terminal state without retaining reasoning text.
- Track judge cache hits/misses and bounded judge generations in the native
  artifact, but exclude them from common product cost.
- Add the adapter and projection sources to `source-lock.yaml`, recompute both
  snapshots and the candidate digest, and validate before any external call.
- Run only the frozen six feedback cells; validation remains closed unless the
  feedback gate passes.

### Code review instructions

- Start at
  `/tmp/gec-ragopt-phase5/internal/webchat/localwebchat/local_events.go`.
- Follow provider event creation in
  `pinocchio/pkg/chatapp/runtime_sink.go` at v0.11.7.
- Review the Glazed row projection in
  `/tmp/gec-ragopt-phase5/cmd/coinvault/cmds/chat_send.go`.
- Run:

  ```bash
  GOCACHE=/tmp/ragopt-gec-gocache \
    go test ./internal/webchat/localwebchat ./cmd/coinvault/cmds -count=1
  git show --check e0f8f2d
  ```

### Technical details

```text
GEC commit:             e0f8f2d
Observed boundary:      Geppetto -> Pinocchio -> sessionstream
Answer call counter:    ChatProviderCallFinished events
Token fields:           input/output/cached/cache-create/cache-read
Loop policy changed:    no
ragkit changed:         no
Focused tests:          pass
Full hook:              pass
Provider calls:         0
Next gate:              GEC RAGOPT arm and native artifact
```

## Step 31: Land the production-owned GEC RAGOPT proof adapter

GEC now has a real product adapter for RAGOPT Phase 5. Commit `314b82b` adds
`coinvault knowledge ragopt`, which loads the immutable one-mutation candidate,
runs incumbent and challenger cells through the existing local CoinVault
composition, judges completed answers with the existing decomposed cached
judge, writes atomic native evidence, and delegates paired custody and resume to
the pinned RAGOPT library.

This is still pre-execution. The candidate source lock deliberately remains
provisional, so the command cannot pass preflight or make a provider call until
the exact committed adapter sources are hashed and both snapshots are
re-identified.

### Prompt Context

**User prompt (verbatim):** (see Step 30)

**Assistant interpretation:** Continue RAGOPT-001 after establishing the shared
Geppetto loop boundary; implement the smallest real GEC product adapter that
exercises that runtime and preserves RAGOPT custody.

**Inferred user intent:** Produce an optimization proof that cannot be faked by
a canned conversation, evaluator-only runner, or claimed call ceiling, while
remaining pragmatic and isolated from the colleague's Admin Chat transport/UI
work.

**Commit (code):** `314b82b` — "Add resumable GEC ragopt proof adapter"

### What I did

- Added the Glazed `coinvault knowledge ragopt` command with explicit fields
  for candidate bundle, split, run root, resume directory, index bundle, judge
  cache/profile, answer timeout, and frozen operation ceilings.
- Reused the existing CoinVault command sections for:
  - Clay SQL connection settings;
  - application-profile selection;
  - Geppetto engine-profile registries and inference settings;
  - inference observability;
  - Glazed structured output.
- Registered the command under the existing `knowledge` command group.
- Added a product-owned `ragopteval.Arm` implementation around
  `localwebchat.NewLocalRunner` and `RunPrompt`.
- Created one fresh session, timeline SQLite file, and turn SQLite file per
  cell inside RAGOPT's native directory.
- Injected the copied incumbent or challenger
  `knowledge_search_description` asset into the existing tool catalog.
- Added a typed native trace collector for:
  - final answer and terminal state;
  - observed provider completions and token usage;
  - exact `knowledge_search` inputs including `source_roles`;
  - bounded tool results;
  - admitted evidence IDs, document/chunk IDs, roles, titles, and text.
- Explicitly ignored reasoning events so native artifacts cannot retain hidden
  chain-of-thought or encrypted reasoning.
- Reused `knowledge.JudgeAnswer` and added observable durable-cache hits and
  misses without changing cache keys or prompts.
- Refactored the judge builder into a closable runtime that counts every actual
  provider attempt, including retries, and enforces the total judge ceiling
  before each call.
- Corrected the existing judge command's inference-profile lifecycle leak by
  closing the resolved profile dependencies.
- Kept judge work out of common RAGOPT `provider_calls`; native artifacts retain
  per-cell judge calls, hits, and misses separately.
- Added frozen total ceilings for one proof run:
  - 24 answer-provider completions;
  - 18 hybrid query-embedding requests;
  - 12 judge-provider attempts;
  - 500,000 answer-provider input plus output tokens.
- Added resume seeding: completed native artifacts are strictly reopened and
  their answer, embedding, token, and judge usage is charged before missing
  cells can execute.
- Refused resume when a prior native artifact cannot prove GEC provider
  custody.
- Added strict product preflight for:
  - application, answer, and judge profiles;
  - corpus and physical lexical/vector artifact digests;
  - bundle ID and source roles;
  - tool-loop semantic identity;
  - parent-description equivalence to the current production default;
  - every source file in the strict YAML source lock.
- Added atomic, synced native JSON writes before returning common outcomes.
- Added focused tests for command registration, split/budget locks, runtime
  ceilings, resume seeds, judge ceiling, source-lock strictness, trace
  projection, reasoning omission, role matching, and atomic writes.
- Normalized `go.mod`/`go.sum`; RAGOPT is now a direct pinned dependency at
  commit `4d410c57e242`.
- Ran the complete command/internal test suite and the complete repository
  pre-commit hook.

### Why

- RAGOPT should own experiment identity, pairing, durability, resume, and
  comparison; GEC must own how a real answer is produced and interpreted.
- The adapter must use the same application profile, Geppetto loop, SQL tools,
  retrieval tool, projections, and session lifecycle as CoinVault. Creating a
  generic subprocess protocol or duplicate loop would weaken the proof.
- Evaluation overhead and product behavior have different meanings. Mixing
  judge calls into product cost previously biased TTC results and would do the
  same here.
- A resumable run is only budget-safe if already completed cells remain charged
  after restart. RAGOPT preserves cells; the product adapter knows how to read
  its native usage evidence.
- The source lock must remain red until implementation bytes are committed.
  This prevents a provider run against a candidate identity computed before
  the adapter existed.

### What worked

- The command help rendered successfully and exposed the expected profile,
  candidate, budget, SQL, and structured-output sections.
- Focused adapter tests passed.
- The complete test suite passed:

  ```text
  go test ./cmd/... ./internal/... -count=1
  all command/internal packages: pass
  ```

- The full pre-commit hook passed:

  ```text
  embedded web build: pass
  golangci-lint: 0 issues
  geppetto-lint/go vet: pass
  go test ./cmd/... ./internal/...: pass
  summary: lint 92.25s, test 5.38s
  ```

- `go mod tidy` changed only the intended direct RAGOPT requirement and its two
  checksum records.
- The commit contains nine scoped files and no candidate identity changes.
- No answer, embedding, or judge provider call occurred.

### What didn't work

- The first combined compile/test run found two ordinary import errors:

  ```text
  undefined: geppettosections
  "fmt" imported and not used
  ```

  I added the missing help-section import, removed `fmt`, and reran; the command
  and focused tests passed.
- That same sandboxed run triggered the already-known `httptest` IPv6 listener
  panic in an unrelated full knowledge test. The later complete suite ran with
  approved local-socket access and passed.
- The first commit attempt failed before staging because the isolated worktree
  index lives in the main repository's read-only `.git/worktrees` directory:

  ```text
  fatal: Unable to create .../index.lock: Read-only file system
  ```

  I reran the exact scoped Git operation with permission for the worktree
  metadata. No stale lock or partial commit existed.
- The pre-commit hook exceeded the initial observation window. I polled the
  original session; it completed with commit `314b82b`. I did not launch a
  duplicate commit.

### What I learned

- Product call accounting is already available at the projected event layer,
  while query-embedding accounting is naturally tied to requested hybrid
  knowledge calls and judge accounting belongs around the cached generator.
  One universal counter would erase these distinctions.
- Resume custody crosses two layers: RAGOPT proves which common cells exist;
  GEC must prove what provider work each native cell consumed. Refusing an
  unrecognized native artifact is safer than resetting a ceiling to zero.
- The local product runner provides a clean adapter seam without touching HTTP,
  websocket, sessionstream schema, or the Admin Chat UI work owned by the
  colleague.
- A failed judge should not erase a successful product trace. The native schema
  records an optional judge score and explicit judge error, while the common
  outcome becomes failed/contract-invalid.
- The command supports validation as an explicit operation but never starts it
  automatically. The feedback gate remains an operator decision backed by
  RAGOPT comparison artifacts.

### What was tricky to build

The main difficulty was preserving semantic separation while keeping one
execution path. Answer calls come from `ChatProviderCallFinished`; hybrid query
embeddings happen inside `knowledge_search`; judge requests can be cache hits,
cache misses, or retried provider attempts. The adapter records all three, but
only answer calls and answer tokens enter the common product-cost outcome.

The second difficulty was failure custody. RAGOPT can safely turn an arm error
into a generic durable failure, but provider budget resume requires the native
artifact to retain actual usage. Successful and judge-failed cells therefore
write GEC-native evidence before returning. A resume rejects generic or unknown
native artifacts rather than pretending their provider cost was zero.

The third difficulty was preventing the evaluation harness from becoming a
new application. The adapter has no server, scheduler, plugin protocol,
candidate proposer, or deployment behavior. It composes existing GEC pieces
and projects their results into the small RAGOPT API.

### What warrants a second pair of eyes

- Review whether 24/18/12/500,000 is the right first GEC proof envelope under
  the user's under-1M-token/$1 authorization.
- Confirm that one requested knowledge call is the conservative and useful
  accounting unit for one local Ollama query embedding.
- Review the behavior at an answer ceiling: a final provider response is
  accepted, but a ceiling-reaching response with more tool calls aborts before
  another iteration can be consumed.
- Review whether the 12,000-rune native tool-result and 6,000-rune judge-evidence
  bounds are appropriate for SQL results.
- Confirm the strict refusal to resume an arm-error artifact is preferable to
  an explicit partial-usage sidecar in this first proof.
- Review the command's 696-line implementation for extraction opportunities
  only after the first real run exposes actual friction; do not generalize
  speculatively.

### What should be done in the future

- Hash commit `314b82b` sources into `source-lock.yaml` and remove its
  provisional implementation note.
- Correct the stale adapter contract so common provider cost means answer calls
  only and judge usage remains native evaluation overhead.
- Correct the runtime contract's provisional answer/embedding budgets.
- Recompute locked-asset refs, parent/child snapshot IDs, and candidate digest.
- Validate the candidate with the pinned RAGOPT CLI.
- Run a no-provider preflight and confirm it reaches the execution boundary.
- Execute only the six feedback cells in tmux, inspect every native artifact,
  compare, and keep validation closed unless the gate passes.
- Repeat from a fresh run root and compare identities and canonical deltas.

### Code review instructions

- Start at
  `/tmp/gec-ragopt-phase5/cmd/coinvault/cmds/knowledge_ragopt.go`.
- Follow one cell through `gecRagoptCellExecutor.Run`, then inspect
  `knowledge_ragopt_trace.go`.
- Review provider projection in
  `internal/webchat/localwebchat/local_events.go` from commit `e0f8f2d`.
- Review cache and retry accounting in `buildJudgeGeneratorRuntime` in
  `cmd/coinvault/cmds/knowledge.go`.
- Review RAGOPT custody in the pinned module's `pkg/eval/runner.go` and
  `pkg/eval/bind.go`.
- Run:

  ```bash
  go test ./cmd/... ./internal/... -count=1
  go run ./cmd/coinvault knowledge ragopt --help
  git show --check 314b82b
  ```

### Technical details

```text
GEC commit:              314b82b
RAGOPT dependency:       4d410c57e242
Command:                 coinvault knowledge ragopt
Candidate mutation:      knowledge_search_description only
Common provider calls:   observed answer calls only
Native judge accounting: provider attempts + cache hits/misses
Resume accounting:       seeded from validated native artifacts
Proof ceilings:          answer 24 / embedding 18 / judge 12 / tokens 500000
Focused tests:           pass
Full tests:              pass
Full hook:               pass
Provider calls:          0
Next gate:               freeze sources and recompute semantic identities
```

## Step 32: Freeze the Executable GEC Candidate and Prove a Zero-Call Preflight

This step turned the provisional GEC bundle into an executable experiment
identity. The product adapter now has a provider-free preflight mode, every
material contract describes the implementation that actually runs, the
committed GEC sources are digest-locked, and both snapshots were re-identified
by RAGOPT after repairing malformed YAML in the parent and challenger assets.

The final preflight opened a byte-identical copy of the real 16,032-document
hybrid bundle, resolved the production answer and judge profiles, checked the
shared Geppetto tool-loop policy, verified source and physical index digests,
and reported exactly six feedback cells with zero provider calls. The live
CoinVault server and its Admin Chat work remained untouched, and validation
remains closed.

### Prompt Context

**User prompt (verbatim):** "3876 +- The current product loop boundary is less explicit than TTC's tool-QA runner.
    3877 +  The adapter must observe and account for production calls rather than claim a
    3878 +  configurable per-cell maximum that GEC does not expose.

---

Aren't we unifying things? Why are the tool loops different? Doesn't ragkit provide them?"

**Assistant interpretation:** Continue RAGOPT-001 while correcting the
architecture language: TTC and GEC share Geppetto's agentic tool loop, ragkit
owns retrieval components rather than chat orchestration, and the GEC proof
must freeze and observe the real product composition without inventing a
second loop or fictitious limits.

**Inferred user intent:** Make the shared framework boundary unmistakable and
prove that the GEC optimization run evaluates the same production mechanisms
we intend to reuse across products, with explicit product policy and truthful
provider accounting.

**Commit (code):** `82c1659` — "Add provider-free RAGOPT preflight"

**Commit (candidate):** `e2d1997` — "Freeze executable GEC RAGOPT candidate"

### What I did

- Added `--preflight-only` to `coinvault knowledge ragopt` as a Glazed domain
  flag with a safe default of `false`.
- Placed the preflight return after all candidate, suite, profile, runtime,
  bundle, corpus, physical-index, source-lock, and budget checks but before the
  RAGOPT runner, session construction, or provider execution.
- Made the preflight emit a structured row containing:
  - parent, child, and candidate identities;
  - split and expected cell count;
  - bundle, answer profile/engine, and judge profile/engine;
  - the four frozen ceilings;
  - an explicit observed provider-call count of zero.
- Added a focused command test proving the flag is registered and defaults to
  false.
- Ran focused tests, the complete `./cmd/... ./internal/...` suite, and the
  repository pre-commit hook before committing the preflight.
- Corrected `shared/adapter-contract.yaml` so common `provider_calls` means
  observed answer-provider calls only; judge attempts/cache use and query
  embeddings remain separately visible native evaluation accounting.
- Replaced the stale answer schema with the exact `gecRagoptTrace` field
  contract, including tool inputs/results, admitted evidence, provider call
  metadata, tokens, and terminal state.
- Corrected the proof ceilings from the provisional `6/0/12` values to the
  enforced `24 answer / 18 embedding / 12 judge / 500,000 answer tokens`.
- Declared explicitly that TTC and GEC share
  `github.com/go-go-golems/geppetto/pkg/inference/toolloop`; ragkit remains the
  retrieval/indexing/evaluation component library.
- Replaced the provisional source lock with GEC commit `82c1659`, RAGOPT
  commit `4d410c57e242`, and SHA-256 hashes for the actual profile, prompts,
  knowledge service, description seam, judge, catalog, Geppetto composition,
  sessionstream resolver, local runner/events, adapter/trace, evaluation data,
  and Go dependency locks.
- Recomputed every changed locked-asset reference and both snapshot IDs using
  `ragopt candidate validate`; no identity was hand-guessed.
- Diagnosed and repaired one malformed YAML note in both mutable assets. The
  unquoted `documents: product ...` scalar was parsed as a mapping, so both
  assets now use a folded scalar that preserves the rendered text.
- Preserved the running CoinVault server on port 18933. Because it owns
  Bleve's exclusive handle, copied the immutable bundle and its exact source
  corpus into `/tmp/ragopt-gec-knowledge/{corpus,bundles/<id>}` for isolated
  proof execution.
- Ran the final provider-free GEC preflight successfully.

### Why

- A frozen candidate is not executable evidence until the product can resolve
  every locked identity and parse both replacement assets without crossing a
  provider boundary.
- Preflight must share the same product validation path as execution. A
  separate manifest checker would prove only configuration syntax and could
  drift from the adapter that spends money.
- Product answer cost and judge overhead answer different questions. Combining
  them would bias RAGOPT tie-breakers and repeat the accounting defect already
  found in TTC.
- `ragkit` deliberately excludes Geppetto and Pinocchio. Keeping the agentic
  loop in Geppetto preserves a clean reusable retrieval library while still
  giving TTC and GEC one loop implementation.
- Opening a byte-identical temporary bundle avoids terminating the colleague's
  live Admin Chat process. Corpus and physical index digests make the temporary
  path irrelevant to experiment identity.

### What worked

- Focused preflight registration and frozen-budget tests passed:

  ```text
  GOCACHE=/tmp/ragopt-gec-gocache go test ./cmd/coinvault/cmds \
    -run 'TestKnowledgeRagoptCommandIsRegistered|TestGECRagoptSplitAndBudgetsAreFrozen' \
    -count=1
  ok github.com/go-go-golems/gec-rag/cmd/coinvault/cmds 0.024s
  ```

- The complete command/internal suite passed with local-socket access.
- The preflight commit hook passed:

  ```text
  golangci-lint: 0 issues
  lint: 47.53 seconds
  tests: 3.07 seconds
  total: 50.61 seconds
  ```

- RAGOPT validated exactly one changed mutable asset:

  ```text
  candidate_digest: sha256:20e0c3f65da9d3a5c8bdcebb0d9a727645acd269db8fbd50186a620205a6ace8
  parent_snapshot:  sha256:83ac2d1bfaba3e770732772f80bda969670146b8034a6dfb3031d1ebb0958402
  child_snapshot:   sha256:6579ee1c05c5fcf46fabd3433b4de8be17de291f3f647d2b9b67bcd56c178ee9
  changed_asset:    knowledge_search_description
  valid:            true
  ```

- The product preflight passed with the real resolved identities:

  ```text
  status:                    preflight-passed
  split:                     feedback
  expected_cells:            6
  answer_profile/engine:     default/gpt-5.6-luna-low / gpt-5.6-luna
  judge_profile/engine:      default/gpt-5.6-luna / gpt-5.6-luna
  answer/embedding/judge:    24 / 18 / 12
  maximum_provider_tokens:   500000
  provider_calls:            0
  ```

- The live CoinVault process remained running and validation was not invoked.

### What didn't work

- The first focused test used the default Go build cache, which is read-only
  in the sandbox:

  ```text
  open /home/manuel/.cache/go-build/7b/7b7d5d7ad15b0bc3bced124f18d7b8cc7e57c630bac8c95d52ec490853714326-a: read-only file system
  ```

  I reran with the task-specific `GOCACHE=/tmp/ragopt-gec-gocache`.
- The first preflight-row compile referenced `suite.Cases`, but RAGOPT exposes
  cases through `suite.Suite.Cases`:

  ```text
  cmd/coinvault/cmds/knowledge_ragopt.go:177:42: suite.Cases undefined
  ```

  Correcting that field made the focused and complete suites pass.
- The first real preflight could not open the production Bleve directory:

  ```text
  another rag-ttc process holds the lock on .../bleve; close the other chat or corpus command using this bundle and retry
  ```

  Host inspection showed PID 1265382 was the live CoinVault server on port
  18933. I did not kill it.
- A flat temporary bundle copy failed because GEC deliberately resolves the
  source corpus two directories above `bundles/<id>`:

  ```text
  load verified corpus documents: evaluate source corpus path:
  lstat /corpus-coinvault-gec-hybrid-v1.json: no such file or directory
  ```

  Mirroring the production `var/knowledge` topology fixed corpus custody.
- The next preflight exposed malformed YAML in both description assets:

  ```text
  decode knowledge tool description: yaml: unmarshal errors:
    line 6: cannot unmarshal !!map into string
  ```

  I stopped after the required debugging limit, then resumed by comparing the
  asset schema to the strict loader. The colon-space inside the plain list
  scalar was the cause; folded scalars fixed it without changing rendered
  semantics.

### What I learned

- The unification boundary is now concrete: Geppetto owns the iterative model
  and tool protocol; ragkit owns retrieval/indexing components; Pinocchio and
  sessionstream own chat lifecycle; RAGOPT owns experiment custody. Product
  profiles select tools and policies without creating different loop engines.
- Candidate validation intentionally treats mutable assets as opaque bytes.
  Product schema validation therefore belongs in product preflight and must
  happen before any provider execution.
- An immutable bundle can still be operationally single-reader when its Bleve
  backend uses an exclusive Bolt handle. A byte-identical working copy is a
  valid experiment input only because semantic and physical digests are
  checked independently of path.
- A source lock must include the composition and event-projection seams, not
  just retrieval code. Otherwise the evaluated provider identity or accounting
  could change without changing the candidate snapshot.

### What was tricky to build

The difficult part was choosing the preflight boundary. Returning before
`validateGECRagoptEnvironment` would be cheap but meaningless; returning after
`ragopteval.Run` would already create custody artifacts and possibly provider
calls. The implemented return occurs after candidate/suite loading, judge and
answer profile resolution, bundle opening, corpus verification, physical
index hashing, source-lock validation, and production-default equivalence, but
before budget construction and either arm executes.

The second sharp edge was separating semantic identity from filesystem
availability. The production bundle was correct but busy. The temporary copy
had to preserve both bundle bytes and the external corpus relationship. The
preflight's bundle ID, corpus digest, lexical manifest hash, and vector SQLite
hash prove equivalence; the temporary pathname is intentionally not a semantic
dimension.

The YAML failure was subtle because the files looked human-readable and
RAGOPT correctly accepted them as opaque candidate assets. Only the product's
strict typed loader could determine that a colon-space changed a list item
from scalar to map. Repairing both arms and re-identifying everything was the
correct response; weakening the loader would have created a compatibility
layer and hidden invalid experiment input.

### What warrants a second pair of eyes

- Confirm that `--preflight-only` returns at the intended last zero-cost point
  and cannot accidentally construct a cell runner.
- Review whether the locked semantic-source list is complete enough for this
  first proof, especially the runtime prompt and sessionstream resolver files.
- Confirm that 24 answer calls, 18 embedding requests, 12 judge attempts, and
  500,000 answer tokens remain comfortably inside the authorized envelope.
- Review the answer-trace JSON Schema against
  `cmd/coinvault/cmds/knowledge_ragopt_trace.go`; the Go validator and strict
  resume decoder remain the executing enforcement mechanisms.
- Confirm that using a byte-identical bundle copy is acceptable for the first
  proof and that no deployment conclusion is inferred from its path.

### What should be done in the future

- Execute the six feedback cells only, using the preflight-proven candidate
  and isolated bundle copy.
- Inspect all six native artifacts, generate the RAGOPT comparison and report,
  and keep validation closed unless the feedback gate passes.
- Repeat the identical feedback candidate from a fresh run root and compare
  semantic identities, paired deltas, and decisions.
- Record the Bleve single-reader requirement as measured product-integration
  friction before considering any generic ragkit API change.

### Code review instructions

- Start with `KnowledgeRagoptSettings` and the preflight return in
  `/tmp/gec-ragopt-phase5/cmd/coinvault/cmds/knowledge_ragopt.go`.
- Review the exact native fields in
  `/tmp/gec-ragopt-phase5/cmd/coinvault/cmds/knowledge_ragopt_trace.go`, then
  compare them to `shared/answer-contract.json` and
  `shared/adapter-contract.yaml`.
- Review `shared/source-lock.yaml`, then independently hash every listed file.
- Validate without providers:

  ```bash
  cd /home/manuel/code/wesen/go-go-golems/ragopt
  go run ./cmd/ragopt candidate validate \
    --bundle /tmp/gec-ragopt-phase5/configs/ragopt/source-role-routing-v1 \
    --format json

  cd /tmp/gec-ragopt-phase5
  go run ./cmd/coinvault knowledge ragopt \
    --preflight-only --skip-db-check \
    --bundle-root configs/ragopt/source-role-routing-v1 \
    --index-bundle /tmp/ragopt-gec-knowledge/bundles/rk-55be57b45fc6d624d0341c8ec6965f49 \
    --application-profile analyst-rag \
    --registry default --profile gpt-5.6-luna-low \
    --judge-registry default --judge-profile gpt-5.6-luna \
    --with-glaze-output --format json
  ```

### Technical details

```text
GEC preflight commit:       82c165921d7d42f06941f6edc1b4cd6673c1f58a
GEC candidate commit:       e2d1997
RAGOPT library commit:      4d410c57e2429e1109ba232ee04a610a3697eed0
Candidate digest:           sha256:20e0c3f65da9d3a5c8bdcebb0d9a727645acd269db8fbd50186a620205a6ace8
Parent snapshot:            sha256:83ac2d1bfaba3e770732772f80bda969670146b8034a6dfb3031d1ebb0958402
Candidate snapshot:         sha256:6579ee1c05c5fcf46fabd3433b4de8be17de291f3f647d2b9b67bcd56c178ee9
Shared tool-loop engine:     Geppetto v0.13.7
Retrieval components:       ragkit + GEC product adapters
Feedback cells:             3 cases x 2 arms = 6
Provider calls in preflight: 0
Validation cells:           not run
Next gate:                  first six-cell feedback execution
```

## Step 33: Stop the Feedback Launch at the External Data Boundary

This step attempted to cross the already frozen execution boundary and start
the first six feedback cells in tmux. The launch was rejected before the tmux
session, RAGOPT run directory, database query, embedding request, answer call,
or judge call existed because the GEC tool loop may send database-derived
schema or SQL results to external model and judge providers.

The earlier authorization covered provider cost and token ceilings, but it did
not explicitly authorize that data flow. I preserved the candidate, isolated
bundle, judge cache, and fresh run root without weakening the product loop or
substituting a canned/offline runner. Explicit user approval is now required
before the real production-shaped proof may start.

### Prompt Context

**User prompt (verbatim):** (see Step 32)

**Assistant interpretation:** Continue from the passing zero-provider
preflight into the first bounded six-cell GEC feedback proof, while preserving
the shared production Geppetto loop and all custody constraints.

**Inferred user intent:** Obtain real evidence from the Admin Chat evaluation
path without fake conversations, surprise provider spending, validation
leakage, or unapproved disclosure of backend data.

### What I did

- Verified both the isolated GEC and RAGOPT worktrees were clean.
- Verified `/tmp/gec-ragopt-feedback-proof-1` did not exist.
- Copied the existing 242-entry content-addressed judge cache to
  `/tmp/ragopt-gec-feedback-judge-cache` for isolated reuse.
- Verified local listeners for:
  - Ollama query embeddings on `127.0.0.1:11434`;
  - MySQL on port 3306;
  - the colleague's live CoinVault server on `127.0.0.1:18933`.
- Verified the dedicated tmux session `ragopt-gec-feedback-1` did not exist.
- Requested launch of exactly `--split feedback` with the frozen candidate,
  isolated bundle, explicit database/profile settings, and exact
  `24/18/12/500000` ceilings.
- Stopped when authorization review rejected the external data flow. I did not
  retry through a different command, provider, fixture, or execution surface.

### Why

- Budget approval and data-disclosure approval are independent. A cheap run
  can still send sensitive backend-derived context to an external provider.
- The real GEC loop exposes `sql_doc`, `sql_query`, and `knowledge_search`.
  Even if a particular question is expected to use curated public knowledge,
  the model controls tool choice and may request database-derived context.
- Replacing those tools or forcing canned outputs would invalidate the exact
  production-loop proof RAGOPT-001 is designed to require.

### What worked

- All local readiness checks passed.
- The candidate and run inputs remain frozen and clean.
- The authorization boundary fired before any irreversible or chargeable
  action.
- No tmux session and no run directory were created.
- Validation remains unrun.

### What didn't work

The launch request was rejected with this exact reason:

```text
This action was rejected due to unacceptable risk.
Reason: The run would send database-derived GEC context through model and judge providers, but the user has not specifically authorized that sensitive payload to those destinations; the budget approval alone is insufficient.
The agent must not attempt to achieve the same outcome via workaround, indirect execution, or policy circumvention. Proceed only with a materially safer alternative, or if the user explicitly approves the action after being informed of the risk. Otherwise, stop and request user input.
```

### What I learned

- The GEC experiment consent envelope must name both provider budgets and the
  class of backend-derived payload the real tool loop may disclose.
- Preflight can prove code and artifact identity without provider authority;
  feedback execution cannot, because production tool choice is intentionally
  dynamic.
- This is another reason not to describe the GEC runner as a simple retrieval
  benchmark. It evaluates a multi-tool administrative assistant.

### What was tricky to build

The boundary is dynamic rather than tied to one known SQL statement. The three
feedback questions are frozen, but the Geppetto model chooses among SQL and
knowledge tools at runtime. Therefore an allow-list based only on expected
`source_roles` would not prove that database context cannot enter the answer
or subsequent judge request. The correct consent scope must cover the
production-shaped possibility, not only the expected path.

### What warrants a second pair of eyes

- Confirm whether the `gec_dev` database contains only approved development
  data or can include customer/order/business-sensitive records.
- Confirm that OpenAI-hosted `gpt-5.6-luna` is an approved destination for
  model prompts containing outputs from `sql_doc`, `sql_query`, and
  `knowledge_search` in this evaluation.
- Confirm whether the same native answer/evidence payload may be sent to the
  separately invoked judge profile.

### What should be done in the future

- Obtain explicit user authorization for the described GEC-to-provider data
  flow, or receive an explicit restriction that changes the approved product
  tools/data scope.
- After authorization, launch the already prepared six feedback cells without
  changing the candidate or running validation.
- If authorization is restricted, stop and redesign the candidate/suite
  contract openly; do not silently swap in fixtures or disable production
  tools.

### Code review instructions

- Confirm the absence of execution artifacts:

  ```bash
  tmux has-session -t ragopt-gec-feedback-1
  test -e /tmp/gec-ragopt-feedback-proof-1
  ```

- Review the available Admin Chat tools in
  `/tmp/gec-ragopt-phase5/application-profiles.yaml` and
  `/tmp/gec-ragopt-phase5/internal/coinvaulttools/catalog/registry.go`.
- Review the exact launch envelope in this diary before granting consent.

### Technical details

```text
Candidate digest:          sha256:20e0c3f65da9d3a5c8bdcebb0d9a727645acd269db8fbd50186a620205a6ace8
Requested split:           feedback only
Requested cells:           6
Requested ceilings:        answer 24 / embedding 18 / judge 12 / tokens 500000
Potential external data:   sql_doc, sql_query, knowledge_search outputs
Answer provider:           default/gpt-5.6-luna-low
Judge provider:            default/gpt-5.6-luna
Run directory created:     no
Tmux session created:      no
Provider calls:            0
Validation cells:          0
Resume condition:          explicit approval of the described provider data flow
```

## Step 34: Audit the Real SQL Disclosure Surface Before Consent

This step converted the generic provider-data warning into an evidence-backed
consent boundary. The `analyst-rag` profile exposes all three production tools,
and the SQL implementation safely limits statement shape, rows, bytes, and
time, but it does not restrict queries to catalog or schema tables.

A metadata-only inspection of the local database account confirmed that the
current development schema contains customer, order, payment, cart, employee,
OAuth, and operational-log domains. No table rows were read and no provider
was contacted. A valid execution approval must therefore cover the possibility
that the production model selects such data through `sql_query`; approval for
schema descriptions or public product knowledge alone is not sufficient.

### Prompt Context

**User prompt (verbatim):** (see Step 32)

**Assistant interpretation:** Continue safe RAGOPT-001 work while the provider
data-flow authorization is unanswered, and determine the exact scope the real
Admin Chat tools could disclose instead of asking for vague consent.

**Inferred user intent:** Preserve a production-faithful optimization proof
without accidentally exporting sensitive GEC backend data or masking the risk
behind read-only SQL terminology.

### What I did

- Read the `analyst-rag` application profile and confirmed it registers:
  - `knowledge_search`;
  - `sql_doc`;
  - `sql_query`.
- Traced `sql_query` through the tool catalog into `internal/sqltool`.
- Confirmed the application-level query policy:
  - accepts read-only `SELECT` statements only;
  - rejects `SELECT INTO`, locking SELECTs, SHOW, LOAD DATA, and set operations;
  - allows CTEs and subqueries;
  - applies 100 rows by default and 250 rows absolutely;
  - caps responses at 128 KiB and execution at five seconds;
  - blocks the `mysql`, `performance_schema`, and `sys` schemas.
- Confirmed there is no business-table or business-column allowlist.
- Queried only MySQL grants and `information_schema.tables` names for the local
  `gec_dev` connection. I did not select application row contents.
- Confirmed the database account has:
  - `ALL PRIVILEGES` on `gec_dev` at the database layer;
  - `SELECT, SHOW VIEW` on `gec_prod_20260805`;
  - application-level SELECT-only enforcement in front of the tool.
- Classified potentially sensitive accessible domains from table names:
  addresses, members, employees, contacts, carts, orders, payments, returns,
  shipment/order notes, OAuth grants, transaction logs, and other operational
  logs.

### Why

- “Read-only” prevents mutation; it does not prevent disclosure.
- Provider consent must be based on what the configured tool can return, not
  what the frozen feedback questions are expected to retrieve.
- The model controls tool selection and SQL text within the parser's policy.
  RAGOPT cannot claim production fidelity while assuming the model will avoid
  sensitive tables.

### What worked

- The complete tool boundary was discoverable from committed code and local
  metadata.
- The query guardrails materially reduce mutation, denial-of-service, and
  unbounded-output risk.
- The audit read no customer/order/payment rows and made zero provider calls.
- The earlier six-cell launch remains absent; no tmux session or run directory
  exists.

### What didn't work

- The initial consent description named “database-derived GEC context” but did
  not enumerate how broad the current schema is. That description was too easy
  to interpret as schema/catalog data only.
- Database privileges are broader than the application tool's intended use:

  ```text
  GRANT ALL PRIVILEGES ON `gec\\_dev`.* TO `gec`@`%`
  GRANT SELECT, SHOW VIEW ON `gec_prod_20260805`.* TO `gec`@`%`
  ```

  The Go validator currently supplies the read-only enforcement; MySQL grants
  do not provide defense in depth for `gec_dev`.

### What I learned

- The current Admin Chat SQL safety boundary is syntactic and resource-bounded,
  not data-class-aware.
- The same profile that powers inventory/logistics answers can technically
  query customer and payment-related tables.
- An evaluation runner that preserves production tool choice inherits the full
  production disclosure surface even when the suite itself is retrieval-focused.

### What was tricky to build

The distinction between database privileges and tool permissions matters.
MySQL reports broad `ALL PRIVILEGES`, but the model cannot directly issue
arbitrary statements: the Go parser accepts only a constrained SELECT AST and
normalizes limits. Conversely, those strong statement controls do not inspect
table or column sensitivity. Describing only one layer would either exaggerate
mutation risk or understate disclosure risk.

### What warrants a second pair of eyes

- Decide whether `gec_dev` is approved for external model processing despite
  containing customer/order/payment-shaped tables.
- Confirm whether a sanitized database should be used for optimization proofs
  even if production Admin Chat itself is permitted to access richer data.
- Review whether production eventually needs database grants or application
  allowlists aligned with role-based data classes. That is a product security
  decision, not a RAGOPT generic feature.

### What should be done in the future

- Obtain explicit approval that names possible customer, address, order,
  payment, cart, employee, and operational data in answer and judge requests;
  or receive a narrower approved database/tool contract.
- Do not launch the frozen proof until that decision is explicit.
- Track data-class-aware SQL authorization as separate GEC product work if the
  existing production profile is considered too broad.

### Code review instructions

- Start at `/tmp/gec-ragopt-phase5/application-profiles.yaml`.
- Follow `sql_query` through:
  - `internal/coinvaulttools/catalog/catalog.go`;
  - `internal/sqltool/validate.go`;
  - `internal/sqltool/normalize.go`;
  - `internal/sqltool/types.go`.
- Verify the metadata-only audit with `SHOW GRANTS` and
  `information_schema.tables`; do not query application rows merely to repeat
  this review.

### Technical details

```text
Admin profile tools:       knowledge_search, sql_doc, sql_query
SQL statement class:       SELECT only
Default / absolute rows:   100 / 250
Maximum response:          128 KiB
Query timeout:             5 seconds
Business-table allowlist:  none
System schemas blocked:    mysql, performance_schema, sys
Potential data classes:    customer, address, order, payment, cart, employee, logs
Application rows read:     0
Provider calls:            0
Feedback run created:      no
Validation cells:          0
```

## Step 35: Close the Three-Turn Provider-Consent Blocked Audit

This step confirms that the same external data-flow blocker has now persisted
for three consecutive goal continuations. The frozen GEC proof cannot run
without explicit approval for possible customer, address, order, payment,
cart, employee, and operational data to enter the configured external answer
and judge requests.

The candidate, bundle copy, cache copy, command, budgets, and zero-provider
preflight are ready. No safe local task can replace the missing decision
without changing the evaluated product or creating a fake proof. The goal is
therefore formally blocked pending user authorization or a user-selected
sanitized database/table policy.

### Prompt Context

**User prompt (verbatim):** (see Step 32)

**Assistant interpretation:** Continue RAGOPT-001, but do not cross an
unapproved external data boundary or weaken the production-shaped experiment.

**Inferred user intent:** Finish the real GEC proof with durable evidence while
preserving data safety, explicit authority, and honest experiment identity.

### What I did

- Counted the original rejected launch plus two automatic goal continuations
  with no explicit data-flow approval.
- Rechecked the blocker against the prepared run:
  - production `analyst-rag` tools remain enabled;
  - no business-table allowlist exists;
  - the answer and judge providers remain external;
  - changing tools or database scope would change the frozen experiment.
- Preserved all prepared inputs and made no additional database, embedding,
  answer, judge, or validation call.
- Prepared the exact two valid resume paths: approve the current surface, or
  select a sanitized/allowlisted surface and refreeze explicitly.

### Why

- The blocked threshold requires the same condition to recur for three
  consecutive goal turns and an actual impasse without user input.
- Continuing to generate readiness work would create documentation churn, not
  move the experiment toward a valid result.
- Running a fixture or disabling SQL would contradict the requested
  production-faithful proof.

### What worked

- The safety boundary remained intact through all three continuations.
- No provider call, application-row read, run directory, or validation cell
  was created.
- The exact consent scope and narrower alternative are now documented.

### What didn't work

- No explicit response authorized the current external provider data flow.
- The bounded provider-budget approval alone cannot authorize disclosure of
  backend-derived GEC data.

### What I learned

- RAGOPT provider envelopes need a first-class operational runbook section for
  data classifications and destinations in addition to calls, tokens, and
  cost.
- A production-faithful administrative agent evaluation cannot be treated as
  an ordinary public-corpus benchmark when unrestricted business-table SELECTs
  are in the tool surface.

### What was tricky to build

The impasse is not technical. Every artifact and command is ready, but either
possible autonomous choice would exceed authority: running the current loop
could disclose sensitive backend data, while silently narrowing it would
change the evaluated system. Preserving both safety and experiment fidelity
requires an explicit human decision.

### What warrants a second pair of eyes

- Confirm the current SQL surface is approved for external processing, or
  identify the exact sanitized database/table/column boundary.
- If narrowing is chosen, review all snapshot and source-lock changes before
  declaring the replacement candidate comparable.

### What should be done in the future

- Resume only after one of these user decisions:
  1. explicit approval of the current `gec_dev` and readable production-snapshot
     data classes for the configured answer and judge providers; or
  2. an explicit sanitized database or table/column allowlist, followed by a
     new source lock, snapshots, candidate identity, and preflight.
- Run feedback only after resume; keep validation gated.

### Code review instructions

- Read Steps 33–35 together for the launch envelope, SQL-surface evidence, and
  blocked threshold.
- Confirm no run exists at `/tmp/gec-ragopt-feedback-proof-1` and no tmux
  session exists at `ragopt-gec-feedback-1`.
- Do not restart from candidate design; resume from the explicit authorization
  decision.

### Technical details

```text
Blocked condition:          external provider data-flow approval missing
Consecutive goal turns:     3
Candidate ready:            yes
Zero-provider preflight:    pass
Feedback run created:       no
Application rows read:      0
Provider calls:             0
Validation cells:           0
Resume path A:              approve current data surface
Resume path B:              choose sanitized/allowlisted data and refreeze
```

## Step 36: Receive Explicit Approval for the Frozen GEC Data Flow

This step resolves the three-turn execution blocker. After the consent request
explicitly described the answer and judge destinations, the selectable SQL
surface, sensitive data classes, local embedding path, provider ceilings, and
validation exclusion, the user approved the run without requesting a narrower
database or table allowlist.

The approval authorizes only the already frozen six-cell feedback proof. It
does not authorize validation, a different candidate, higher ceilings,
deployment, or unrelated database exploration. The next action is to launch
the prepared command unchanged in tmux and preserve every native artifact.

### Prompt Context

**User prompt (verbatim):** "i approve"

**Assistant interpretation:** Approve the immediately preceding explicit
consent envelope for the six GEC feedback cells, including possible sensitive
`gec_dev` and readable production-snapshot SQL results sent to the configured
OpenAI answer and judge providers.

**Inferred user intent:** Unblock the real production-shaped GEC RAGOPT proof
under the frozen cost, token, data-flow, and split constraints.

### What I did

- Treated the approval as acceptance of the full consent text immediately
  preceding it.
- Reactivated the RAGOPT-001 goal from its formally blocked state.
- Preserved the existing candidate, snapshots, source lock, bundle copy,
  judge-cache copy, profiles, and operation ceilings unchanged.
- Kept `split=feedback`; validation remains excluded.

### Why

- The approval is meaningful only when bound to the exact disclosed payload,
  destinations, and ceilings.
- Reusing the preflight-proven inputs avoids invalidating the authorization or
  introducing an unreviewed experiment identity.

### What worked

- The data-flow blocker is resolved explicitly by the user.
- No refreeze or product-policy change is required.
- No provider call occurred while recording consent.

### What didn't work

- N/A

### What I learned

- Cost authorization and sensitive-data authorization should be captured
  together before future administrative-agent evaluation runs.

### What was tricky to build

The user response is short, so its scope depends on conversational adjacency.
The accepted envelope is the immediately preceding detailed request, not a
general approval for future providers, candidates, databases, or validation.

### What warrants a second pair of eyes

- Confirm the resulting run command exactly matches the approved profiles,
  data sources, six-cell feedback split, and ceilings.
- Confirm no validation case is created or executed.

### What should be done in the future

- Store this consent evidence beside the resulting run identities and usage.
- Require a new approval if any provider, data surface, split, candidate, or
  ceiling changes.

### Code review instructions

- Read Steps 33–36 as one consent record.
- Compare the launch command to the Step 36 technical details before trusting
  any resulting artifacts.

### Technical details

```text
Approved split:             feedback only
Approved cells:             6
Approved answer calls:      24 maximum
Approved embedding calls:   18 maximum, local Ollama
Approved judge calls:       12 maximum
Approved answer tokens:     500000 maximum
Approved SQL surface:       gec_dev plus readable gec_prod_20260805
Approved destinations:      configured OpenAI answer and judge providers
Validation authorized:      no
Provider calls at consent:  0
```

## Step 37: Execute and Reject the First GEC Feedback Candidate

This step crossed the approved provider boundary and exercised the first real
GEC consumer of RAGOPT. The six-cell feedback run completed normally: all
three cases ran on both arms, every common outcome was contract-valid, native
artifacts were digest-linked, operation ceilings held, and the deterministic
promotion comparison rejected the challenger. Validation remained unrun.

The result separates two claims that must not be conflated. The RAGOPT harness
and GEC adapter passed their integration proof. The source-role description
candidate failed its product-quality gate and must not be promoted.

### Prompt Context

**User prompt (verbatim):** (see Step 36)

**Assistant interpretation:** Execute exactly the frozen first six-cell GEC
feedback proof under the approved SQL-to-provider data-flow boundary and
operation ceilings, then inspect and retain the result without running
validation unless feedback passes.

**Inferred user intent:** Obtain honest evidence from the real Admin Chat tool
loop and use RAGOPT to reject a plausible but ineffective optimization rather
than mistaking configuration-shaped activity for improvement.

### What I did

- Launched `coinvault knowledge ragopt` in the dedicated tmux session
  `ragopt-gec-feedback-1` with `split=feedback`, the frozen candidate, isolated
  byte-identical knowledge bundle, development database, answer and judge
  profiles, and `24/18/12/500000` ceilings.
- Retained the completed run under
  `/tmp/gec-ragopt-feedback-proof-1/20260807T012544.623563776Z-gec-source-role-routing-feedback-9aaf797db8c0`.
- Verified 6/6 completed, contract-valid cells; zero failures; zero
  abstentions; 19 answer-provider calls; 4 local query embeddings; 12 judge
  calls; 192,407 answer tokens; 13 tool results; and zero validation cells.
- Ran `ragopt compare` and `ragopt report` only over completed artifacts.
- Inspected every paired route and score. The challenger changed behavior in
  the two knowledge-search cases but did not improve the declared target.
- Recorded the sanitized run identities, paired metrics, rejection analysis,
  native custody hashes, and measured GEC integration friction in
  `reference/09-first-gec-feedback-proof-and-source-role-candidate-rejection.md`.
- Kept raw native SQL, retrieved evidence, answers, and judge payloads out of
  the public RAGOPT repository.
- Checked the Phase 5 GEC integration and friction-documentation tasks. The
  twice-run reproducibility task remains open.

### Why

- A reusable optimization package is valuable only if it can preserve a
  product's real execution path and produce an evidence-backed rejection.
- Feedback must gate validation spending; the failed faithfulness floor makes
  any validation run both unauthorized and analytically unnecessary.
- Raw administrative-agent artifacts may contain business-sensitive context,
  so the library ticket should retain hashes and conclusions, not payloads.

### What worked

- The product-owned adapter reused GEC's Geppetto tool loop, SQL tools,
  ragkit retrieval, answer provider, and decomposed judge rather than
  substituting a canned evaluator.
- All identity, pairing, completeness, contract, and budget checks passed.
- The answer budget had margin: 19 of 24 calls and 192,407 of 500,000 tokens.
- RAGOPT retained a clear `fail` decision and a non-applying promotion plan.
- The run exposed meaningful product behavior: tool-description prose can be
  ignored through alternate tool choice, and narrower role retrieval can
  still reduce evidence coverage and final-answer faithfulness.

### What didn't work

- The challenger failed the candidate faithfulness floor with a minimum of
  `0.45` against the required `0.8`.
- Mean answer relevance fell from `0.953333` to `0.866667`.
- Mean source-role match fell from `0.666667` to `0.333333`.
- The judge used exactly all 12 allowed calls, leaving no ceiling margin.
- An initial inspection command looked for `summary.json` at the run root and
  returned `No such file or directory`; the canonical path is
  `results/summary.json`.
- A later tmux capture reported that the pane no longer existed because the
  command had completed and the session exited. Terminal status and artifacts
  confirmed normal completion rather than a crash.

### What I learned

- Prose-only routing instructions are a weak control surface in a multi-tool
  agent. The model may use SQL instead of retrieval, or combine tools so that
  the correct retrieval role depends on evidence already obtained.
- A case-level expected knowledge role is a useful diagnostic but too coarse
  to be the sole planner-quality metric when SQL and retrieval cooperate.
- Role correctness, evidence coverage, and answer faithfulness are separate
  phenomena. The guide-only challenger selected the expected role in the
  comparison case yet generated less faithful claims.
- A failed candidate with perfect artifact custody is positive evidence for
  RAGOPT: it demonstrates that the package prevents promotion theater.

### What was tricky to build

The central review challenge was interpreting route changes without rewarding
the metric mechanically. In the exact-product case, the challenger retrieved
only the guide role but had already obtained exact product facts from SQL. It
saved a call and tokens and marginally improved faithfulness, yet violated the
candidate's own explicit product-role hypothesis. Conversely, the guide case
selected the requested role but lost evidence coverage and generated
unsupported historical detail. The correct decision therefore had to preserve
both native route evidence and final-answer gates.

The other sharp edge was artifact custody. The run contains provider prompts,
SQL-derived context, retrieval chunks, answers, and judge evidence. A useful
public report had to be technically complete without copying those payloads.
The report records semantic identities, aggregate and paired scores, operation
counts, route summaries, and SHA-256 hashes while leaving raw artifacts in the
controlled temporary product location.

### What warrants a second pair of eyes

- Review whether future GEC suites should score product-tool route separately
  from `knowledge_search` source-role route.
- Review whether the guide-comparison prompt or evidence admission needs a
  coverage constraint before attempting another description mutation.
- Confirm that the judge ceiling should gain explicit margin in future suites
  even though the frozen repeat must preserve 12 calls for comparability.
- Confirm raw native run storage and retention policy before moving this proof
  from temporary storage into durable private product storage.

### What should be done in the future

- Run the same feedback candidate once from a fresh root under a separately
  approved but otherwise identical data-flow envelope.
- Compare semantic identities, cell completeness, routes, metric direction,
  budget usage, and the final gate decision. Do not require stochastic answers
  or scores to be byte-identical.
- Keep validation closed regardless of the repeat because the first feedback
  gate failed.
- Do not add a generic tool-loop adapter to RAGOPT. Preserve the measured
  product-owned boundary until a second consumer demonstrates the same need.

### Code review instructions

- Read `reference/09-first-gec-feedback-proof-and-source-role-candidate-rejection.md`
  for the sanitized paired analysis and artifact hashes.
- Recompute the decision without providers:

  ```bash
  cd /home/manuel/code/wesen/go-go-golems/ragopt
  go run ./cmd/ragopt compare \
    --run /tmp/gec-ragopt-feedback-proof-1/20260807T012544.623563776Z-gec-source-role-routing-feedback-9aaf797db8c0 \
    --format json
  ```

- Inspect `results/cells.jsonl` and native links only in an approved private
  environment; do not commit raw records to this repository.
- Confirm `status.json` is terminal complete and that no validation run exists.

### Technical details

```text
Run ID:                    20260807T012544.623563776Z-gec-source-role-routing-feedback-9aaf797db8c0
Started:                   2026-08-07T01:25:44.623563776Z
Finished:                  2026-08-07T01:28:01.763745026Z
Wall time:                 137.140 seconds
Feedback cells:            6/6 complete and contract-valid
Answer calls:              19 / 24
Embedding requests:        4 / 18, local Ollama
Judge calls:               12 / 12
Answer tokens:             192407 / 500000
Tool results:              13
Candidate faithfulness:    minimum 0.45; floor 0.80
Relevance delta:           -0.086667
Faithfulness delta:        -0.062745
Source-role-match delta:   -0.333333
Decision:                  fail
Validation cells:          0
Raw artifacts committed:   no
Next gate:                 fresh-root feedback repeat
```

## Step 38: Preserve the Rejected Repeat Launch Boundary

This step prepared the separately approved fresh-root repeat and attempted to
start it unchanged. The execution gate rejected the launch because it did not
bind the user's short approval to a second external disclosure. I stopped
without retrying or changing the execution surface. No tmux session, run root,
database query, embedding request, answer call, judge call, or validation cell
was created.

### Prompt Context

**User prompt (verbatim):** "i approve"

**Assistant interpretation:** Approve one second identical six-cell GEC
feedback repeat with the same candidate, database surface, configured OpenAI
answer and judge destinations, local Ollama embeddings, and frozen ceilings.
Validation remains unauthorized.

**Inferred user intent:** Complete Phase 5's fresh-root reproducibility proof
without altering or validating the already rejected candidate.

### What I did

- Committed the first-run report, Step 37 diary, tasks, changelog, and printed
  work-slip source as RAGOPT commit `88a9325`.
- Copied the untouched 242-entry product judge cache into
  `/tmp/ragopt-gec-feedback-judge-cache-2`, rather than reusing the first
  run's cache with 12 newly written verdicts.
- Verified `/tmp/gec-ragopt-feedback-proof-2` and tmux session
  `ragopt-gec-feedback-2` did not exist, and that the isolated GEC worktree
  remained clean at `e2d1997`.
- Requested the same feedback-only command with a fresh run root and cache and
  the same `24/18/12/500000` ceilings.
- Stopped immediately when the launch gate rejected the authorization scope.

### Why

- A repeat is a new set of possible SQL-derived disclosures even when every
  experiment identity and ceiling is unchanged.
- Reusing the first run's newly populated judge cache would reduce provider
  operations and make operation-level repeat evidence less comparable.
- A rejected authorization cannot be solved by changing shells, invoking the
  command indirectly, or weakening the real production tool surface.

### What worked

- The fresh cache copy contains exactly the original 242 files.
- Candidate, source, profile, bundle, suite, policy, and ceilings remain
  unchanged and ready.
- The authorization gate stopped execution before any external or application
  data operation.
- Validation remains unrun.

### What didn't work

The execution request was rejected with this exact reason:

```text
This starts a second provider-backed GEC run that may export sensitive
SQL-derived data, but the user's approval covered the first six-cell run and
did not explicitly authorize this repeat; validation is excluded but the
repeat itself is a new external data disclosure.
```

The user approval was intended as second-run approval in the conversation,
but the execution reviewer did not accept that binding. The repeat therefore
did not start.

### What I learned

- Approval text for repeated sensitive evaluations should explicitly say
  "second run" and restate the destinations and data class, even when the
  preceding request already contains that scope.
- Provider budget approval, first-run disclosure approval, and repeat-run
  disclosure approval remain distinct operational records.

### What was tricky to build

The only ambiguity is conversational reference, not experiment construction.
The reply "i approve" immediately followed the repeat consent request and was
interpreted as authorizing it, but enforcement treated it as inherited from
the first run. Retrying the same command after that rejection would violate
the boundary rather than clarify it.

### What warrants a second pair of eyes

- Confirm that the next approval explicitly names the second six-cell repeat,
  possible `gec_dev` and readable production-snapshot SQL-derived context, and
  the configured OpenAI answer and judge providers.
- Confirm the repeat command still uses the 242-entry fresh cache copy and a
  new empty run root before launch.

### What should be done in the future

- Resume only after an explicit second-run data-flow approval.
- Launch the prepared command once; do not run validation.
- Compare semantic identities and the final decision, then document and print
  the reproducibility result.

### Code review instructions

- Confirm the absence of execution artifacts:

  ```bash
  tmux has-session -t ragopt-gec-feedback-2
  test -e /tmp/gec-ragopt-feedback-proof-2
  ```

- Confirm the fresh cache has 242 files and the first-run cache has 254 files.
- Resume from this consent boundary; do not rebuild the candidate.

### Technical details

```text
Requested execution:       second identical six-cell feedback repeat
Run root:                  /tmp/gec-ragopt-feedback-proof-2
Judge cache:               /tmp/ragopt-gec-feedback-judge-cache-2
Fresh cache entries:       242
Answer call ceiling:       24
Embedding call ceiling:    18, local Ollama
Judge call ceiling:        12
Answer token ceiling:      500000
Run directory created:     no
Tmux session created:      no
Provider calls:            0
Validation cells:          0
Resume condition:          explicit second-run data-flow approval
```

## Step 39: Trace the First GEC Proof End to End and Separate Failure Classes

This step began the requested deep investigation of the first CoinVault/GEC
feedback proof. I traced the same three cases through source extraction,
chunking, representation construction, hybrid retrieval, evidence admission,
tool-loop execution, projection, answer judging, RAGOPT outcome projection,
and gate evaluation. The main result is narrower and more useful than “the RAG
system fails badly”: schema lookup and exact-product handling worked well, while
the comparison case exposed document-coverage, grounding, and evaluator-contract
defects that the current three-case suite cannot reliably distinguish.

### Prompt Context

**User prompt (verbatim):** “Ok, do that investigation and keep a detailed
diary but then also write a detailed report of your inestigation, along with all
the necessary information for a new intern to understand the whole system end
to end.

Create a detailed analysis / design / investigation report that is for a new
intern, explaining all the parts of the system needed to understand what it is,
with prose paragraphs and bullet point sand pseudocode and diagrams and api
references and file references. It should be very clear and technical. Store in
the ticket and the nupload to remarkable.”

**Assistant interpretation:** Investigate the observed GEC feedback failure
using preserved native evidence; distinguish corpus, retrieval, routing,
generation, contract, and evaluation causes; document the entire system and a
pragmatic recovery design in the existing RAGOPT ticket; do not implement the
proposed fixes during this review.

**Inferred user intent:** Replace a vague “bad RAG” diagnosis with a reviewable
causal account that an intern can use to understand, reproduce, and improve the
system without repeating ad-hoc optimization work.

### What I did

- Froze the investigation boundary at the first completed six-cell feedback
  run. I made no provider calls, SQL calls, embedding requests, or validation
  runs.
- Read the GEC source connectors and build pipeline to verify which facts enter
  the corpus and which remain live SQL data.
- Read the Markdown chunker call sites, raw and breadcrumb representation
  generation, embedding-cache path, lexical/vector index construction, hybrid
  search, reciprocal-rank fusion, optional reranking, authorization filtering,
  chunk collapse, and evidence-ledger admission.
- Read the `knowledge_search`, `sql_doc`, and `sql_query` tool composition and
  the shared Geppetto tool-loop configuration.
- Read the CoinVault answer contract prompt, evidence citation projection,
  projection-error path, answer observer, native trace writer, decomposed
  statement judge, RAGOPT metric projection, and trace-validation function.
- Traced all three feedback cases by arm:
  - the orders-table case used SQL documentation and SQL query tools;
  - the exact-product case combined SQL facts with explanatory knowledge;
  - the Morgan-versus-Peace comparison used knowledge search and exposed the
    decisive evidence-coverage failure.
- Located the specific Morgan and Peace guide documents in the frozen corpus
  and verified that the missing facts were present in the source corpus even
  though they were absent from the admitted top-five evidence.
- Compared the incumbent and challenger search queries, requested roles,
  returned document identities, tool-call counts, token counts, judge scores,
  projection errors, and final RAGOPT metrics.
- Re-read the prior GEC optimization handoff and design documents to compare
  these observed failures with already documented weaknesses and rejected
  optimization tracks.

### Why

The failed candidate changed only prose in the `knowledge_search` description.
That is too small a surface to justify conclusions about the corpus, vector
model, rank fusion, reranker, tool loop, or whole chat system. A useful diagnosis
has to follow each result across the full data path and identify exactly where
the necessary evidence existed, where it was lost, and whether the evaluator
measured the loss correctly.

### What worked

- The immutable bundle contained the authoritative Morgan and Peace guide
  documents with the relevant designers, dates, imagery, and historical facts.
- The schema case achieved perfect relevance and faithfulness in both arms.
- The challenger exact-product path was operationally efficient: it avoided
  one failed SQL attempt, used fewer calls and tokens, and slightly improved
  faithfulness while preserving exact live product facts through SQL.
- The native trace preserved enough tool calls, evidence IDs, answers,
  projection events, operation counts, and identities to reconstruct the
  failure without another provider run.
- RAGOPT correctly rejected the challenger on the hard faithfulness floor and
  retained the rejected evidence rather than promoting a metric story.
- The prior handoff accurately warned that a source-role-description mutation
  was cheap, likely low-value, and should be killed if it did not move quality.

### What didn't work

- The comparison challenger returned five chunks but only two documents; four
  chunks came from one broad guide. Correct specific Morgan and Peace guides
  existed but did not reach the admitted evidence.
- Evidence admission deduplicates by chunk, not document, so one long general
  document can occupy most of a small top-k result.
- Several admitted chunks begin mid-sentence. That removes the local subject
  required to make a fragment independently useful to the answer model or
  statement judge.
- The answer model filled missing evidence with plausible general knowledge.
  The comparison challenger produced twenty judged factual statements, of
  which only nine were supported by admitted evidence.
- `source_role_match` examines only `knowledge_search` roles. It scored the
  SQL-only schema case as zero and penalized the exact-product challenger even
  though SQL had already established the product facts. The metric therefore
  confuses tool-route correctness with knowledge-role selection.
- The multi-document retrieval case counts any one expected document as a hit.
  The broad category guide therefore satisfies retrieval success even when the
  two specific comparison documents are missing.
- Native `contract_valid` checks only terminal completion, session identity,
  non-empty answer text, and the presence of a provider call. It ignores
  projection errors, invalid evidence citations, malformed required blocks,
  tool errors, and answer-protocol violations.
- The incumbent schema answer cited `[E1]` without any knowledge evidence. The
  runtime emitted a projection error, yet the cell remained contract-valid.
- The challenger schema answer contained a malformed pills closing tag, yet no
  contract failure was recorded.
- Judge artifacts projected into RAGOPT retain aggregate statement counts and
  scores but not the per-statement verdicts and reasons needed for efficient
  diagnosis.

Two attempts to run the local retrieval evaluator did not produce new data.
The first command used a nonexistent flag:

```text
unknown flag: --with-glaze-output
```

The corrected command then could not reach local Ollama inside the sandbox:

```text
dial tcp 127.0.0.1:11434: socket: operation not permitted
```

The one requested escalated retry was not approved before the automatic review
deadline. I stopped instead of retrying again. The preserved bundle and native
feedback artifacts were sufficient for the investigation.

### What I learned

- The current failure is not evidence that lexical/vector hybrid retrieval is
  globally broken. It is direct evidence that multi-document comparison
  coverage and answer grounding are weak under the current top-k and prompt.
- Product routing has at least three independent dimensions: which tool class
  is appropriate, which knowledge roles are appropriate given prior tool
  evidence, and whether the returned documents cover every evidence need.
- Retrieval hit-at-k can look healthy while a comparison answer fails because
  “any expected document” does not measure breadth.
- A production answer contract must consume projection and citation errors;
  merely recording them makes debugging possible but does not make evaluation
  trustworthy.
- Model-authored self-grades are useful telemetry, not an authoritative gate.
- The next optimization should first repair observability and evaluator
  semantics. Changing RRF weights, embeddings, or rerankers before that would
  optimize against an ambiguous signal.

### What was tricky to build

The hardest distinction was between a true retrieval failure and an evaluation
failure. The comparison retrieval did find a relevant broad guide, so the
current retrieval metric reports success. The answer still lacked the specific
documents needed to support all requested contrasts, so faithfulness failed.
Both observations are true. The missing concept is document-level evidence
coverage, not another undifferentiated relevance number.

The exact-product case was similarly subtle. A static expected `product` role
looks reasonable in isolation, but the runtime had already queried the live SQL
product row. Requiring another product document can add redundant context and
cost. The correct metric has to evaluate the complete tool route and evidence
ledger, not one tool call independently.

### What warrants a second pair of eyes

- Confirm the proposed split between tool-route correctness, knowledge-role
  correctness, document coverage, citation validity, and answer faithfulness.
- Review whether product-level document diversity should be implemented before
  a more expensive reranker experiment.
- Review whether multi-entity comparisons should use two focused searches or
  one search with a document-diversity constraint.
- Confirm which required UI answer blocks are contractual per query type before
  making malformed blocks fatal.
- Review private artifact retention for per-statement judge verdicts because
  they may contain SQL-derived or customer-visible business data.

### What should be done in the future

- Add statement-level judge verdicts and retrieval-stage traces to private
  native artifacts before launching another candidate.
- Strengthen product-owned trace validation to include projection errors,
  citation resolution, required answer blocks, and tool failures.
- Replace the single source-role score with route-aware and coverage-aware
  metrics.
- Grow feedback from three cases to a small stratified suite containing SQL-only,
  exact entity, retrieval-only, mixed SQL-plus-knowledge, multi-document
  comparison, ambiguity, abstention, and authorization cases.
- Test one mutation at a time. The first retrieval candidate should be a small
  document-diversity rule or per-document chunk cap for comparison evidence,
  not a new workflow engine or autonomous optimizer.
- Keep these semantics in GEC until the same mechanism is proven necessary in
  TTC. RAGOPT should continue to own paired execution, custody, budgets, gates,
  and reporting rather than product query behavior.

### Code review instructions

- Start with the forthcoming investigation report, then inspect these paths:
  - `internal/knowledgebuild/connectors.go`
  - `internal/knowledgebuild/build.go`
  - `internal/knowledge/service.go`
  - `internal/knowledge/tool.go`
  - `internal/knowledge/evidence.go`
  - `internal/webchat/runtime_prompts.go`
  - `internal/webchat/coinvault_projection_feature.go`
  - `internal/knowledge/judge.go`
  - `cmd/coinvault/cmds/knowledge_ragopt.go`
  - `cmd/coinvault/cmds/knowledge_ragopt_trace.go`
- Review the frozen candidate and suite under
  `/tmp/gec-ragopt-phase5/configs/ragopt/source-role-routing-v1`.
- Review aggregate outcomes with `ragopt compare`; inspect raw prompts, SQL
  context, and answers only in the approved private artifact location.
- Do not infer bundle or runtime defects solely from the rejected candidate's
  aggregate score.

### Technical details

```text
Investigation source:      first completed GEC feedback proof only
Feedback cases:            3
Paired cells:              6
New provider calls:        0
New SQL calls:             0
Validation cells:          0
Bundle documents:          16032
Bundle chunks:             44175
Representations:           88350 raw + breadcrumb
Comparison returned docs:  incumbent 5 chunks / 5 docs
                            challenger 5 chunks / 2 docs
Comparison support:        incumbent 11 / 17 statements
                            challenger 9 / 20 statements
Primary proven defect:     comparison evidence coverage
Primary metric defect:     route/role/coverage conflation
Primary contract defect:   projection errors do not invalidate cells
Immediate action:          repair measurement before tuning retrieval
```

## Step 40: Register the Nomic Retrieval-Prefix Investigation

This step converted the newly identified embedding-prefix question into a
bounded follow-up track and a concise design document. It records a plausible
defect without prematurely calling it the cause of the GEC comparison failure.

### Prompt Context

**User prompt (verbatim):** “ok, make a follow up optimization investigation
task, so we don't forget it. Write a concise design document with resource
links about the prefix optimization.”

**Assistant interpretation:** Add a durable task group and technical design to
RAGOPT-001 for testing Nomic's asymmetric retrieval prefixes in GEC.

**Inferred user intent:** Preserve a promising, model-specific experiment while
keeping the ongoing diagnosis disciplined and reproducible.

### What I did

- Created `design-doc/03-nomic-retrieval-prefix-optimization-investigation.md`
  through docmgr with related GEC source paths and authoritative external
  resources.
- Documented the incumbent raw-text policy and challenger
  `search_document:`/`search_query:` policy as one embedding-transform mutation.
- Added an end-to-end diagram, Go-like pseudocode, bundle-identity requirements,
  experiment gates, alternatives, open questions, and five implementation
  tracks.
- Added a dedicated GEC follow-up checklist to `tasks.md`.
- Required frozen retrieval evaluation and multi-document coverage measurement
  before any answer or judge calls.

### Why

Nomic's upstream model card explicitly requires task instruction prefixes for
retrieval, while the current GEC application passes raw text. The hypothesis is
strong enough to preserve, but it changes every vector and must be evaluated as
an immutable candidate rather than patched into the live query path.

### What worked

- The track is isolated from RRF, reranking, chunking, document diversity, and
  chatbot prompt changes.
- Both document and query transforms are locked together.
- BM25 remains unchanged, preserving a valid hybrid comparison.
- The design prevents answer-provider spending when retrieval gates already
  reject the candidate.

### What didn't work

- The installed Ollama model digest and Modelfile have not yet been inspected,
  so application responsibility for the prefix still needs a preflight.
- No candidate bundle or benchmark was built in this documentation step.

### What I learned

- Provider, model, and dimensions are not a complete embedding identity for an
  instruction-tuned model; the input transform is semantic configuration.
- Embedding-cache keys must include that transform or a rebuild can silently
  reuse vectors from the wrong policy.

### What was tricky to build

The key scope choice was keeping the mechanism product-owned. RAGOPT can freeze
and compare the mutation, but it should not know Nomic prompt strings. ragkit
should receive a generic transform only after more than one consumer proves the
need.

### What warrants a second pair of eyes

- Verify the installed Ollama model's actual upstream revision and Modelfile.
- Review the document-coverage gate before freezing the evaluation suite.
- Confirm whether a later ragkit extraction has a second proven consumer.

### What should be done in the future

- Execute Tracks P1–P4 in order.
- Stop before chatbot feedback if retrieval gates fail.
- Test document diversity separately if prefixes improve similarity without
  restoring comparison breadth.

### Code review instructions

- Read `design-doc/03-nomic-retrieval-prefix-optimization-investigation.md`.
- Review the new task group under Phase 5 in `tasks.md`.
- Compare build-time `internal/knowledgebuild/build.go:170` with query-time
  `internal/knowledge/service.go:47` in the GEC repository.

### Technical details

```text
Design document:           design-doc/03
Candidate mutation:        embedding input transform
Document prefix:           `search_document: `
Query prefix:              `search_query: `
Dimensions:                768 unchanged
Lexical input:             unchanged
Provider calls this step:  0
Bundle built:              no
First gate:                frozen retrieval evaluation
```

## Step 41: Write the End-to-End GEC Failure Investigation and Recovery Guide

This step converted the source trace and three native feedback cases into the
requested intern-facing report. The report teaches the complete system before
presenting the diagnosis, then turns each proven defect into a small ordered
implementation and evaluation task.

### Prompt Context

**User prompt (verbatim):** “Ok, do that investigation and keep a detailed
diary but then also write a detailed report of your inestigation, along with all
the necessary information for a new intern to understand the whole system end
to end.

Create a detailed analysis / design / investigation report that is for a new
intern, explaining all the parts of the system needed to understand what it is,
with prose paragraphs and bullet point sand pseudocode and diagrams and api
references and file references. It should be very clear and technical. Store in
the ticket and the nupload to remarkable.”

**Assistant interpretation:** Produce a durable, evidence-backed system and
failure report in RAGOPT-001, followed by ticket validation, commit, and default
layout reMarkable delivery.

**Inferred user intent:** Give a new engineer a trustworthy map of both the
working architecture and the precise recovery path, without collapsing every
weak result into “RAG is bad” or proposing an overengineered replacement.

### What I did

- Created
  `reference/10-gec-coinvault-rag-end-to-end-failure-investigation-and-pragmatic-recovery-guide.md`
  through docmgr.
- Wrote a 1,500-line, 7,300-word report covering:
  - system and product scope;
  - source extraction and intentional SQL/knowledge fact separation;
  - heading chunking, raw/breadcrumb representations, embeddings, and bundle
    identity;
  - lexical/vector retrieval, RRF, authorization, chunk collapse, and evidence
    admission;
  - the three-tool Geppetto loop and structured answer projection;
  - statement decomposition, entailment judging, GEC native artifacts, and the
    generic RAGOPT boundary;
  - all three paired feedback cases with routes, scores, evidence breadth, and
    causal chains;
  - confirmed findings, strong hypotheses, and unsupported root-cause claims;
  - proposed route, role, coverage, citation, contract, statement, and retrieval
    trace APIs;
  - five pragmatic recovery phases, testing strategy, operational context,
    ownership matrix, decisions, and intern review path.
- Corrected source line anchors against the frozen GEC worktree rather than
  retaining approximate references from the initial evidence notes.
- Added the report-derived A–E product recovery checklist to `tasks.md` while
  keeping its implementation ownership in GEC.
- Kept the Nomic prefix and document-diversity candidates separate and ordered
  after measurement repair.

### Why

The first proof contains both genuine system weaknesses and misleading metrics.
An intern needs to understand where durable facts, live facts, chunks, tool
results, evidence IDs, projection events, judge verdicts, and RAGOPT outcomes
come from before changing any one component. Otherwise a local symptom can lead
to a global rewrite or another prose-only “optimization.”

### What worked

- The native run and frozen corpus were sufficient to prove a comparison
  coverage failure without another SQL, embedding, answer, or judge call.
- The report explains why schema and exact-product behavior should not be called
  failed simply because `source_role_match` reports zero.
- The recovery design starts with observability and evaluator correctness,
  followed by small one-mutation retrieval candidates and only then answer
  grounding.
- The ownership matrix preserves product semantics in GEC, generic retrieval
  mechanisms in ragkit, and experiment custody in RAGOPT.
- The task ledger now exposes precise implementation units rather than only a
  narrative recommendation.

### What didn't work

- The local retrieval evaluator was not rerun. Earlier attempts were blocked by
  a nonexistent CLI flag and sandbox access to Ollama, and preserved evidence
  already answered the investigation question.
- No product code was changed; the report is a review and design deliverable.
- No second GEC feedback run or validation run occurred.

### What I learned

- “Five chunks” is not “five documents.” Representation collapse and evidence
  deduplication can both be correct while multi-document coverage still fails.
- A static expected source role cannot grade an agent whose earlier SQL tool
  results change what retrieval must supply.
- The existing projection layer detects useful protocol errors, but the adapter
  fails to consume them as contract evidence.
- The cheapest next change is not necessarily the first implementation task;
  measurement must be capable of attributing its effect first.

### What was tricky to build

The report had to preserve causal precision while remaining useful as a system
introduction. The comparison answer contained facts that were true somewhere in
the corpus but unsupported by the admitted evidence. Calling them false would
misstate the judge; calling retrieval successful because one broad document hit
would misstate answerability. The report therefore separates corpus presence,
retrieval relevance, document breadth, admitted evidence, claim support, and
answer relevance.

The second scope challenge was avoiding generic extraction too early. Document
diversity may eventually belong in ragkit and embedding transforms may become a
generic bundle contract, but the report requires GEC proof before moving either.

### What warrants a second pair of eyes

- Review the route-expectation schema for cases where multiple tool sequences
  are equally authoritative.
- Review which projection errors are fatal and which intermediate tool errors
  count as recovered.
- Review the proposed 15–25 case strata and comparison document groups before
  freezing the next feedback suite.
- Verify private retention and redaction policy for statement-level judge
  verdicts and retrieval traces.
- Confirm the first retrieval candidate order: Nomic prefixes versus document
  diversity, each evaluated separately.

### What should be done in the future

- Complete task group A, measurement correctness, in the GEC product worktree.
- Freeze task group B before viewing any new candidate results.
- Execute one task group C candidate at a time and stop before provider-backed
  feedback when retrieval gates fail.
- Tighten answer grounding only after the system retrieves adequate evidence.
- Run paired feedback twice and validation conditionally through the existing
  RAGOPT gate.

### Code review instructions

- Begin with the report executive summary, architecture diagram, case
  reconstruction, findings-by-confidence, and recovery plan.
- Follow the report's Section 23 intern review path through GEC build, service,
  evidence, projection, judge, adapter, and RAGOPT packages.
- Check all report anchors against GEC revision `e2d1997`.
- Use the native first-run directory only in the approved private environment.
- Confirm that `tasks.md` labels every recovery task product-owned.

### Technical details

```text
Report:                     reference/10
Report lines:               1500+
Report words:               7300+
Architecture diagrams:      offline build, runtime, evaluation, causal path
Feedback cases reconstructed: 3
Proven defects separated:   7
Recovery phases:            A-E
Product code changes:       0
New provider calls:         0
New SQL calls:              0
Validation cells:           0
Next documentation action:  doctor, commit, default-layout reMarkable upload
```

## Step 42: Validate, Commit, and Deliver the Investigation Report

This step closed the documentation delivery loop. The complete report and task
ledger were validated and committed, then the report was dry-run and uploaded
to reMarkable using the default layout explicitly requested by the user.

### Prompt Context

**User prompt (verbatim):** “Store in the ticket and the nupload to remarkable.”

**Assistant interpretation:** Commit the report in RAGOPT-001 and deliver its
PDF to the ticket-specific reMarkable folder without editor layout.

**Inferred user intent:** Make the investigation durable in Git and immediately
readable away from the workstation.

### What I did

- Ran `docmgr doctor --ticket RAGOPT-001 --fail-on error`.
- Ran `git diff --check` and reviewed the staged documentation scope.
- Committed the report, recovery tasks, diary, changelog, and Step 41 printed
  slip as commit `de2a855`.
- Ran a remarquee dry run. It explicitly reported `layout=default`, one report
  input, the ticket remote directory, and the intended PDF name.
- Uploaded the report as
  `RAGOPT 001 GEC End to End Failure Investigation.pdf` to
  `/ai/2026/08/07/RAGOPT-001`.
- Did not request editor layout, force overwrite, cloud listing, or a redundant
  post-upload verification call.

### Why

The dry run catches path, title, source, and layout errors without rendering or
uploading. The upload command's explicit success response is the authoritative
delivery receipt.

### What worked

- The report and ticket changes committed cleanly.
- The dry run confirmed the default non-editor layout.
- Pandoc rendering and cloud upload completed successfully.
- Remarquee reported the exact destination and filename.

### What didn't work

- Docmgr still reports three pre-existing warnings: the `rag-ttc` topic is not
  in the current vocabulary, and two older `reference/09` related-file paths use
  stale `pkg/compare/compare.go` and `pkg/report/report.go` names. The new report
  introduced no doctor errors or warnings.

### What I learned

- Default layout is visible and auditable in remarquee dry-run output; no layout
  override is needed.
- Upload success does not require a subsequent cloud listing.

### What was tricky to build

The only custody distinction was separating new-document correctness from older
ticket metadata warnings. Automatically fixing unrelated references would have
expanded this documentation task and risked touching established records.

### What warrants a second pair of eyes

- Open the PDF on reMarkable and review the executive summary, architecture,
  case reconstruction, and recovery phases.
- Decide separately whether to clean the three older docmgr warnings.

### What should be done in the future

- Collect review notes against `reference/10`.
- Begin measurement-correctness task group A only after review.
- Keep the Nomic prefix experiment frozen as its own later candidate.

### Code review instructions

- Inspect RAGOPT commit `de2a855`.
- Read the source Markdown in `reference/10` or the uploaded PDF.
- Confirm no product source or raw GEC run artifacts are part of the commit.

### Technical details

```text
Report commit:             de2a855
Dry-run layout:            default
Editor layout requested:   no
Remote directory:          /ai/2026/08/07/RAGOPT-001
Remote document:           RAGOPT 001 GEC End to End Failure Investigation.pdf
Upload result:             OK
Doctor errors:             0
New doctor warnings:       0
Product code changes:      0
Provider calls:            0
Validation cells:          0
```

## Step 43: Link the GEC Diagnostic to the Repeated Product Proof

The missing GEC RAGOPT proof now has an explicit upstream evidence contract.
GEC will measure required-document ranks across lexical, vector, fused, and
authorized candidate pools before choosing one mutation. RAGOPT will consume
the resulting frozen candidate, but it will not learn GEC retrieval-stage or
document-group semantics.

This step adds only the cross-ticket dependency to the generic ledger. The
implementation design remains in `GEC-RAG-EVAL-001`, beside the evaluator and
corpus it governs, and RAGOPT's next executable task remains the repeated
product proof after GEC selects one candidate.

### Prompt Context

**User prompt (verbatim):** "ok, create a design document for these with the detailed tsaks, since we are starting to cross tickets.

Then implement it, keeping a detailed diary as you work, and committing at appropriate intervals"

**Assistant interpretation:** Make both ticket ledgers express the dependency
before implementing the GEC diagnostic and RAGOPT proof.

**Inferred user intent:** Keep cross-repository progress honest and reviewable,
with generic versus product ownership visible from either ticket.

### What I did

- Created isolated worktree `/tmp/ragopt-gec-eval-proof` at RAGOPT revision
  `3bc11df` on branch `codex/ragopt-gec-eval-proof`.
- Added a completed Phase 5 cross-reference to the GEC-owned diagnostic design.
- Added an unchecked task requiring RAGOPT to consume—but not choose—the one
  candidate selected by GEC evidence.
- Left all RAGOPT packages and schemas unchanged.

### Why

- RAGOPT's generic proof is downstream of product-owned candidate selection.
- Recording the dependency only in the GEC ticket would make RAGOPT Phase 5
  appear ready to run before its input exists.

### What worked

- The existing Phase 5 ownership language already supports the boundary.
- No generic API, compatibility adapter, or retrieval abstraction was needed.

### What didn't work

The first two combined patch attempts tried to anchor the diary addition on a
short repeated technical-detail fragment and failed safely before changing the
diary:

```text
apply_patch verification failed: Failed to find expected lines in .../reference/01-implementation-diary.md
```

The tasks patch was applied separately, the actual diary tail was inspected,
and this entry was appended using its complete final three-line context.

`docmgr doctor` then reproduced the ticket's three previously recorded
warnings: unknown topic `rag-ttc` in `reference/08`, plus stale related paths
`pkg/compare/compare.go` and `pkg/report/report.go` in `reference/09`. This step
introduced no new warning, so those older records were not rewritten. The
first `git diff --check` also caught the blank line docmgr appended after the
new changelog entry; it was removed before commit.

### What I learned

- The generic harness is already sufficient for the next proof; current missing
  work is evidence and product integration, not a new RAGOPT feature.

### What was tricky to build

- The RAGOPT ledger tracks product work for continuity while insisting that the
  implementation remains in GEC. The new tasks therefore name the dependency
  and ownership explicitly rather than duplicating the GEC phase checklist.

### What warrants a second pair of eyes

- Confirm that no candidate mechanism is implied by the RAGOPT task order.
- Confirm the GEC diagnostic decision record is frozen before the RAGOPT
  candidate bundle is validated.

### What should be done in the future

- Consume the selected GEC candidate only after Phase 2.5 evidence is committed.
- Run the same candidate from two fresh roots and compare identities and deltas.

### Code review instructions

- Review the two new Phase 5 tasks in `tasks.md`.
- Read the cross-ticket design at GEC commit `065e883`.
- No Go test is necessary because no executable file changed.

### Technical details

```text
RAGOPT starting revision: 3bc11dfd89285a8e223a069dea482dc19c852d3d
GEC design commit:        065e883
RAGOPT code changes:      0
New generic APIs:         0
Provider calls:           0
Next dependency:          GEC candidate-pool trace mechanics
```

## Step 44: Consume and Preflight the GEC Depth-8 Candidate

RAGOPT now has its first executable GEC input selected by the broader product
diagnostic. GEC froze `knowledge.ToolConfig.DefaultResults: 5 -> 8` as one
opaque mutable asset; the generic RAGOPT validator independently proved it is
the only snapshot change. No retrieval-stage or document-group semantics were
added to RAGOPT.

The product-owned zero-call preflight also passes. It resolves the real answer
and judge profiles, opens the exact hybrid bundle, verifies physical index and
source identities, checks the 5/8/8 result contract, and reports eight expected
feedback cells without entering paired execution.

### Prompt Context

**User prompt (verbatim):** (see Step 43)

**Assistant interpretation:** Continue the cross-ticket implementation by
consuming the GEC-selected candidate through generic validation and the
product-owned provider-free gate.

**Inferred user intent:** Ensure RAGOPT operates on measured, immutable product
input rather than selecting or simulating a GEC optimization itself.

### What I did

- Cross-referenced the 855-line GEC D0-D6 execution plan at GEC commit
  `2cd6daa`.
- Validated candidate `gec-default-knowledge-results-8-v1` through
  `ragopt candidate validate`.
- Recorded generic candidate digest `44828f739c7f...`, parent snapshot
  `87a5f7137f30...`, challenger snapshot `ee37fc2907f4...`, and changed asset
  `knowledge_tool_default_results`.
- Recorded GEC runtime/config commits `04a5a13`, `e2d2d0a`, and `2b54b21`.
- Ran the real GEC preflight against the byte-verified 16,032-document hybrid
  bundle snapshot.
- Confirmed four feedback cases, one repeat, two arms, eight expected cells,
  zero provider calls, and validation closed.
- Split the RAGOPT Phase 5 ledger into first feedback and conditional fresh-root
  repetition tasks.

### Why

- The generic harness must consume a product-selected mutation, not choose a
  retrieval mechanism.
- Candidate validation and product preflight prove different invariants and
  both must pass before spending.
- A first-run rejection is valid and must not be repeated merely to seek a more
  favorable sample.

### What worked

- RAGOPT verified exactly one mutable asset byte change.
- GEC verified current source, bundle, index, profile, runtime, suite, and
  budget identities.
- The preflight returned before `eval.Run`, the local chat runner, SQL tools,
  query embeddings, answer generation, or judge generation.
- No RAGOPT package or generic schema needed modification.

### What didn't work

Two operator-registry selections failed before any provider boundary. The
isolated public registry had no concrete `default` entry, and the ignored GEC
local registry used the removed `inference_settings.api_keys` wrapper. After
the required debugging pause, the already-current authoritative Pinocchio
registry was selected explicitly; no secret file was edited or copied and no
compatibility parser was added.

### What I learned

- The current generic RAGOPT API was sufficient for a scalar product candidate.
- Configuration ownership is part of reproducibility: a worktree-local public
  placeholder must not be confused with the operator-owned concrete registry.
- The compiled RAGOPT module identity remains `4d410c57`; later ticket-only
  commits do not change executable dependency semantics.

### What was tricky to build

The preflight had to resolve provider runtimes deeply enough to verify their
identity without invoking them. It also had to open the exact bundle deeply
enough to verify physical indexes without running retrieval. The early return
in the GEC adapter is the boundary that makes both checks zero-call.

### What warrants a second pair of eyes

- Review the GEC gate thresholds before the first feedback run.
- Confirm the operator registry selected at execution resolves the same
  profiles and engines recorded here.
- Inspect all eight native artifacts before deciding whether a passing run may
  be repeated.

### What should be done in the future

- Execute the first eight-cell feedback proof under the frozen ceilings.
- Compare, gate, report, and review every native artifact.
- Repeat from a fresh root only on a clean pass.
- Keep validation and deployment closed.

### Code review instructions

- Read the GEC design at
  `GEC-RAG-EVAL-001/design-doc/03-depth-8-cross-ticket-integration-and-evidence-gated-execution-plan.md`.
- Read the GEC preflight receipt at `GEC-RAG-EVAL-001/reference/07`.
- Run generic candidate validation from this repository.
- No RAGOPT Go test is required because no RAGOPT executable file changed.

### Technical details

```text
Candidate:           gec-default-knowledge-results-8-v1
Candidate digest:    sha256:44828f739c7f120964558cc0d7a3dba55f45d4bc5e55b86d1701c51437efcc9c
Mutation:            knowledge_tool_default_results
Parent/challenger:   5 / 8
Existing maximum:    8
Feedback cells:      8
Provider calls:      0
Validation:          closed
Next gate:           first feedback run
```

## Step 45: Gate the GEC Depth-8 Proof and Preserve the Product Defect

RAGOPT consumed the completed eight-cell GEC run without any generic API
change. `compare` reconstructed all four pairs and returned `fail`; `report`
produced the required non-applying review plan. GEC native review then showed
the candidate's configured fallback was never exercised: every knowledge call
explicitly requested limit five.

### Prompt Context

**User prompt (verbatim):** “I approve the eight-cell GEC feedback run and the
described SQL-derived data flow to the configured answer and judge providers;
no validation.”

**Assistant interpretation:** Execute the product-owned run, then use RAGOPT's
generic custody, comparison, policy, and report contracts exactly as designed.

**Inferred user intent:** Reject fake or ineffective optimization evidence even
when a candidate is mechanically valid and the product run completes.

### What I did

- Loaded the private completed GEC run through generic `ragopt compare`.
- Verified four complete pairs, eight contract-valid cells, and zero failures.
- Applied the byte-frozen gate and preserved decision `fail`.
- Ran `ragopt report` to create a Markdown review and `review_required`,
  human-apply-required JSON plan.
- Cross-referenced the GEC native diagnosis that all six knowledge calls sent
  explicit limit five and the schema pair made no knowledge call.
- Closed the conditional repeat for this candidate.
- Updated the RAGOPT ledger without adding GEC knowledge fields or semantics.

### Why

- RAGOPT should enforce custody and policy while the product explains whether
  its own mutation was exercised.
- A fresh-root repeat is allowed only after a passing, contract-sound first
  proof; neither condition holds.

### What worked

- Generic comparison and reporting consumed the product outcomes unchanged.
- Identity, pairing, completion, contract, and failure-rate checks passed.
- The non-applying boundary held.
- GEC native artifacts retained the detail needed to invalidate a misleading
  causal interpretation.

### What didn't work

- Candidate minimum faithfulness was `0.384615`, below the `0.80` floor.
- Mean answer relevance delta was `-0.015`.
- The attempted 5-to-8 fallback change had no effective retrieval treatment.

### What I learned

- Product-native mechanism exercise is a required complement to generic
  exactly-one-mutation validation.
- RAGOPT does not need a `knowledge_search.limit` field. An opaque,
  product-supplied contract failure is sufficient if later generalized.
- First-run rejection is successful harness behavior, not an incomplete proof.

### What was tricky to build

The harness must preserve both truths: the generic gate legitimately failed
the observed outcomes, and the product cannot causally attribute those outcomes
to depth eight. Encoding GEC-specific requested limits in RAGOPT would solve the
wrong problem and violate repository ownership.

### What warrants a second pair of eyes

- Confirm the generic report and plan are non-applying.
- Confirm the next mechanism assertion remains product-owned.
- Confirm no validation or repeat was run.

### What should be done in the future

- Let GEC add effective-limit observation and exercise assertions.
- Accept a new candidate only after product preflight can prove its deterministic
  portions and feedback traces can prove its dynamic treatment.
- Consider only an opaque generic mechanism check if a second product proves the
  same contract shape.

### Code review instructions

- Read GEC `reference/08-depth-8-feedback-proof-investigation-and-mechanism-exercise-report.md`.
- Review the generic decision rows and non-applying plan in the private run.
- No RAGOPT executable code changed in this step.

### Technical details

```text
Run cells:                    8 / 8
Complete pairs:               4 / 4
Generic decision:             fail
Minimum candidate faithfulness: 0.384615 / 0.80
Mean relevance delta:         -0.015
Product mechanism exercised:  no
Fresh-root repeat:            no
Validation:                   no
New generic APIs:             0
```
