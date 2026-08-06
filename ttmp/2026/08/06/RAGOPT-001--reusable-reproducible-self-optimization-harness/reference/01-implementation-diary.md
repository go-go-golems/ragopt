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
