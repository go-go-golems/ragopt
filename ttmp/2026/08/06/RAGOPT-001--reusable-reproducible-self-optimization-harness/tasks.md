# Tasks

This is the canonical implementation ledger. Phases are sequential. A later
phase does not start until the preceding exit criterion is demonstrated.

## Phase 0: Repository and contract design

- [x] Initialize `ragopt` from committed `go-template` content without copying template Git metadata.
- [x] Normalize module, command, logging, Makefile, CI, and release identifiers.
- [x] Run generation, formatting, complete tests, and complete build.
- [x] Create ticket `RAGOPT-001`, primary design guide, diary, and task ledger.
- [x] Audit GEC-RAG handoff, review, design, diaries, and implementation boundaries.
- [x] Audit RAG-TTC experiment, tool-eval, diagnostic, and semantic-config implementations.
- [x] Separate implemented evidence from proposed GEPA/warehouse/reflector work.
- [x] Review and accept the v1 scope and public terminology.

Exit criterion: reviewers agree that v1 is an evidence-gated optimization
harness and that autonomous proposal, transcript warehousing, and deployment
are excluded.

## Phase 1: Immutable run store and semantic identity

- [x] Define `ragopt-run/v1` manifest, status, input-reference, and summary schemas.
- [x] Port the proven `pkg/experiment` lifecycle into a narrowly named `pkg/runstore` package.
- [x] Replace RAG-TTC-specific module and directory assumptions with documented generic names only.
- [x] Preserve path confinement for every artifact write.
- [x] Preserve atomic JSON writes and append-plus-`fsync` JSONL writes.
- [x] Copy declared inputs into `inputs/` and record SHA-256, size, role, and original path.
- [x] Record host, Go module, start time, config digest, and caller-supplied semantic dimensions.
- [x] Reject writes after a run becomes terminal.
- [x] Preserve all prior artifacts when a run fails.
- [x] Add tests for path escape, interrupted JSONL, terminal writes, duplicate completion, and stable config digests.
- [x] Add a run reader that validates schema and exposes status without mutating artifacts.
- [x] Document the run-directory contract and recovery guarantees.

Exit criterion: a fixture process can be interrupted after N results, and all N
synced records plus copied inputs remain valid and inspectable.

## Phase 2: System snapshots and one-mutation candidate bundles

- [x] Define `ragopt-snapshot/v1` with system name, locked assets, mutable assets, and semantic dimensions.
- [x] Define one canonical `AssetRef` containing logical name, media type, copied path, SHA-256, and size.
- [x] Canonically sort logical maps before computing snapshot identity.
- [x] Define `ragopt-candidate/v1` with ID, parent snapshot, proposer identity, hypothesis, expected improvement, and regression risks.
- [x] Store complete replacement assets; do not support model-generated patch application.
- [x] Resolve every candidate path inside its bundle root and reject symlink/path escapes.
- [x] Strictly reject unknown YAML fields, multiple YAML documents, missing files, and digest mismatches.
- [x] Compare parent and candidate snapshots and require exactly one changed mutable asset.
- [x] Reject changes to locked assets or semantic dimensions.
- [x] Make evaluated candidate directories immutable by contract and status transition.
- [x] Add fixtures for valid text mutation, two mutations, locked mutation, missing asset, digest mismatch, and unsafe path.
- [x] Implement `ragopt candidate validate` as a Glazed structured-output command.

Exit criterion: a human-authored asset replacement produces one deterministic
candidate digest, and every multi-surface or locked-input mutation is rejected.

## Phase 3: Resumable paired evaluation runner

- [x] Define `ragopt-suite/v1` with stable case IDs, groups, opaque JSON input, and suite digest.
- [x] Define the small `Arm` interface that executes a case and returns a common `Outcome` projection.
- [x] Keep native artifacts application-owned and require each outcome to reference its native artifact path and digest.
- [x] Define outcome fields for completion, contract validity, abstention, failure, metrics, calls, tokens, duration, and semantic identity.
- [x] Require unique arm names, case IDs, candidate IDs, and repeat indices.
- [x] Execute the incumbent and challenger over the same ordered cases and repeat numbers.
- [x] Append and sync one result cell immediately after completion.
- [x] Resume by loading completed cell keys and rejecting identity mismatches.
- [x] Run sequentially in v1; do not add worker pools until measured runtime requires them.
- [x] Continue after an arm error by recording a failed cell; abort only on custody or identity errors.
- [x] Record the exact suite, snapshots, candidates, policy, and caller config as run inputs.
- [x] Add deterministic scripted-arm tests for success, failure, interruption, resume, and duplicate cells.
- [x] Add an integration fixture proving resumed output equals uninterrupted output after canonical sorting.

Exit criterion: an interrupted paired run resumes without re-running completed
cells and produces the same canonical results as an uninterrupted fixture run.

## Phase 4: Paired comparison, gates, and promotion report

- [x] Join incumbent and candidate outcomes strictly by case ID and repeat index.
- [x] Reject missing, duplicate, cross-suite, cross-policy, or cross-snapshot cells.
- [x] Compute per-cell metric deltas without replacing native values.
- [x] Aggregate wins, ties, losses, failures, and deltas by metric and case group.
- [x] Keep completion, contract validity, and failure rates separate from quality means.
- [x] Define `ragopt-gate-policy/v1` with hard conditions, one declared target metric, regression limits, and ordered tie-breakers.
- [x] Require products to choose metrics and thresholds; ship no universal RAG quality threshold.
- [x] Evaluate hard gates before target improvement and cost tie-breakers.
- [x] Treat missing or invalid outcomes as failures, never as dropped samples.
- [x] Retain every rejected candidate and machine-readable rejection reason.
- [x] Generate a Markdown promotion report containing identities, hypothesis, asset diff, paired table, group aggregates, all regressions, and gate decisions.
- [x] Generate a machine-readable promotion plan that describes but does not apply the product change.
- [x] Add golden tests for pass, hard-gate fail, target fail, catastrophic regression, tie-break, and incomplete pairing.
- [x] Implement `ragopt compare` and `ragopt report` Glazed commands over existing artifacts.

Exit criterion: the fixture candidate yields a deterministic decision and a
reviewable report; changing a failed cell into a missing cell cannot improve
the decision.

## Phase 5: First real integration and proof cycle

Phase 5 validates product-owned adapters around retrieval and evaluation
runtimes. It does not own either chat product: the RAG-TTC proof targets the
existing `tool-loop ragopt` evaluation command, not the canonical Admin Chat
cutover, `sessionstream`, or the customer-facing Garden Assistant. A product
proof passes when identities, custody, paired outcomes, gates, and the final
decision reproduce from a fresh run root. The evaluated candidate may be
correctly rejected; candidate promotion is not a phase-exit requirement.

- [x] Select one existing RAG-TTC human-authored text candidate with no safety-policy mutation.
- [x] Implement the RAG-TTC consumer arm in the RAG-TTC repository, not in `ragopt`.
- [x] Declare RAG-TTC's required locked dimensions: corpus, index, suite, answer model, judge, prompts, tool safety, and evaluator.
- [x] Run incumbent and candidate once on feedback with the corrected locked runtime.
- [x] Enforce the feedback-before-validation spending gate; I5 failed, so validation was correctly left unrun.
- [x] Verify native RAG-TTC artifacts remain authoritative and are digest-linked from common outcomes.
- [x] Review every feedback case outcome and record the reject decision.
- [x] Repeat the same candidate from a fresh run root and verify semantic identities and canonical deltas.
- [x] Integrate GEC-RAG only after the RAG-TTC integration proof reproduces; the I5 candidate does not need to pass promotion gates.
- [x] In GEC-RAG, fix authorization and judge-accounting P0 findings before using results for promotion.
- [x] Cross-reference the GEC-owned candidate-pool diagnosis and candidate
  selection contract in `GEC-RAG-EVAL-001` without adding retrieval-stage or
  document-group semantics to RAGOPT.
- [ ] Consume the single candidate selected by the frozen GEC diagnostic; do
  not preselect prefixes, diversity, fusion, or reranking in RAGOPT.
- [ ] Run one GEC-RAG human-authored candidate twice with a frozen suite and policy.
- [x] Document integration friction before adding any generic API to `ragopt`.

### GEC follow-up optimization investigation: Nomic retrieval prefixes

This is a product-owned investigation recorded here so it is not lost after
the first source-role candidate rejection. It does not block RAGOPT v0.1 unless
the GEC proof chooses it as the next candidate.

- [x] Document the Nomic `search_document:` / `search_query:` hypothesis and
  authoritative resources in
  `design-doc/03-nomic-retrieval-prefix-optimization-investigation.md`.
- [ ] Record the installed Ollama `nomic-embed-text` digest, upstream revision,
  and Modelfile; verify whether the server already transforms raw input.
- [ ] Add document-coverage expectations for frozen multi-document retrieval
  cases before inspecting candidate results.
- [ ] Freeze one embedding-transform candidate with both prefixes; do not
  combine it with RRF, reranking, chunking, or document-diversity changes.
- [ ] Include transform version and exact document/query prefixes in the GEC
  bundle's semantic identity and embedding-cache namespace.
- [ ] Apply the document prefix only on embedding input and the query prefix
  only on vector-query input; keep lexical BM25 input unchanged.
- [ ] Build and verify fresh incumbent and challenger vector bundles.
- [ ] Run frozen lexical/vector/fused retrieval evaluation and compare overall,
  group, and multi-document coverage metrics plus latency.
- [ ] Reject without answer/judge calls if retrieval gates fail.
- [ ] If retrieval gates pass, run bounded paired feedback twice; keep
  validation closed until feedback passes.
- [ ] Record whether the mechanism remains GEC-owned or has evidence for later
  extraction into ragkit; do not generalize it into RAGOPT.

### GEC recovery track after the first feedback investigation

These are product-owned follow-up tasks derived from
`reference/10-gec-coinvault-rag-end-to-end-failure-investigation-and-pragmatic-recovery-guide.md`.
They are recorded in this integration ticket for continuity; implementation
belongs in GEC unless a mechanism is independently proven generic.

#### A. Measurement correctness

- [ ] Preserve statement-level judge verdicts, evidence IDs, and reasons in the
  private GEC native artifact.
- [ ] Record lexical, vector, fused, authorized, reranked, and returned hit IDs,
  ranks, document IDs, roles, and scores.
- [ ] Make unresolved citation and projection errors fail the final answer
  contract.
- [ ] Validate required final-answer blocks by query type.
- [ ] Distinguish recovered intermediate tool errors from unresolved final
  failures.
- [ ] Replace `source_role_match` with tool-route, knowledge-role,
  document-coverage, and citation metrics.

#### B. Feedback-suite coverage

- [ ] Freeze 15–25 stratified feedback cases covering SQL-only, exact entity,
  retrieval-only, mixed SQL/knowledge, comparison, jargon, ambiguity,
  abstention, and authorization.
- [ ] Declare evidence needs and allowed tools per case.
- [ ] Add required document groups and complete coverage@k for comparisons.
- [ ] Freeze a separate held-out validation suite before candidate results are
  inspected.

#### C. Targeted retrieval candidates

- [ ] Execute the Nomic prefix track above as its own one-mutation candidate.
- [ ] If evidence still concentrates by document, freeze a separate simple
  per-document chunk-cap/diversity candidate.
- [ ] Protect exact-product, schema, and single-guide groups from regression.
- [ ] Do not combine prefix, diversity, chunking, RRF, reranker, or synonym
  changes in one candidate.

#### D. Answer grounding

- [ ] After retrieval evidence is adequate, freeze one prompt candidate that
  forbids filling missing material facts from model memory.
- [ ] Require adjacent citations for material knowledge claims and explicit
  qualification when evidence is insufficient.
- [ ] Measure unsupported-claim rate and citation coverage separately from
  relevance.

#### E. Proof and release boundary

- [ ] Run a passing candidate twice from fresh feedback roots with declared
  provider and token ceilings.
- [ ] Inspect all failures and protected-group regressions.
- [ ] Keep validation closed unless feedback passes.
- [ ] Produce a non-applying promotion plan and leave bundle activation and
  rollback to the GEC deployment path.

Exit criterion: two product repositories use the library without a generic
subprocess/plugin protocol, and one human candidate completes the path twice in
each repository with explainable identities and paired outcomes. A stable,
evidence-backed rejection is a valid completed decision; validation spending
remains conditional on feedback gates.

## Phase 6: CLI hardening and release

- [ ] Rename or split policy byte-digest and semantic-digest fields across pre-v0.1 schemas, product adapters, tests, and docs; add no compatibility shim.
- [ ] Add `ragopt run inspect` for manifest, input, result-count, and terminal-state inspection.
- [ ] Add structured table/JSON/YAML output to all artifact commands through Glazed.
- [ ] Add `--log-level` through Glazed fields and zerolog; do not read environment variables.
- [ ] Add schema/API reference documentation and one scripted fixture example.
- [ ] Add corruption diagnostics for malformed JSONL tails, missing inputs, and digest mismatches.
- [ ] Run `go test ./...`, `go build ./...`, `make lint`, `make logcopter-check`, and `make gosec`.
- [ ] Verify release metadata and disabled docs publishing configuration.
- [ ] Tag v0.1 only after the real integration exit criterion passes.

Exit criterion: v0.1 is useful as a library and artifact CLI, with passing CI
and evidence from at least two consumers.

## Phase 7: Shared production refresh contracts (post-v0.1)

Phase 7 implements the fixed control-plane lifecycle designed in
`design-doc/02-production-index-build-scheduling-resumability-and-ragopt-integration.md`.
It does not add a scheduler, daemon, generic workflow engine, indexer, or
deployment mutator to `ragopt`.

- [ ] Freeze TTC and GEC fixture snapshots and the cross-product responsibility matrix.
- [ ] Define generic build-run, build-event, artifact-reference, and terminal-result v1 schemas.
- [ ] Define semantic build identity and canonical serialization rules.
- [ ] Define `ProductBuilder`, `ProductEvaluator`, `ProgressSink`, registry, and activation-plan interfaces.
- [ ] Specify content-refresh challenger custody without weakening the one-mutation optimization candidate.
- [ ] Implement a local event log and projection by reusing `pkg/runstore` durability mechanics.
- [ ] Implement the fixed snapshot → build → verify → evaluate → gate → plan coordinator.
- [ ] Implement local acquire/resume, heartbeats, verified artifact reuse, and cancellation.
- [ ] Add fixture and filesystem registries plus conformance tests for production registry adapters.
- [ ] Add Glazed `ragopt refresh run` and `ragopt refresh inspect` commands.
- [ ] Prove interruption/replay equivalence with a deterministic fixture adapter.
- [ ] Implement the CoinVault/GEC build and evaluation adapter after its P0 evaluation fixes.
- [ ] Prove a local interrupted/resumed GEC build and a no-change cached rebuild.
- [ ] Implement the shared TTC index adapter around the already benchmarked pipeline.
- [ ] Implement and prove the Garden native evaluation adapter.
- [ ] Implement the Admin live retrieval/tool evaluation adapter after its chat cutover stabilizes.
- [ ] Require separate Garden and Admin gates when both consume one TTC bundle.
- [ ] Document measured integration friction before moving or generalizing `flow`.

Exit criterion: CoinVault/GEC and TTC use the same refresh lifecycle and event
contracts while retaining native builders, query runtimes, evidence, and
deployment ownership; interrupted refreshes replay safely and produce stable
decisions.

## Product deployment tracks (outside `ragopt`)

- [ ] Package one digest-pinned refresh image per product.
- [ ] Select AWS Batch by default, or River only where operational Postgres and an always-on Go worker already exist.
- [ ] Configure EventBridge Scheduler, coarse retries, DLQ, leases, and least-privilege IAM.
- [ ] Store private builds, verified bundles, native evidence, and channel pointers durably.
- [ ] Implement conditional activation, health-checked rollout, and previous-bundle rollback in each product.
- [ ] Add change-event reconciliation only after nightly full refresh is stable.
- [ ] Run worker-death, corrupt-artifact, concurrent-activation, and rollback drills.

These tasks belong to the consuming product/deployment repositories. They are
tracked here to make the end-to-end dependency explicit, not to move AWS or
deployment code into the `ragopt` module.

## Explicitly deferred and not part of v1

- [ ] Do not implement a reflector or model proposer before Phase 5 evidence exists.
- [ ] Do not implement a transcript warehouse in `ragopt`.
- [ ] Do not implement a generic workflow engine, scheduler, daemon, web UI, or plugin system.
- [ ] Do not implement automatic production mutation or deployment.
- [ ] Do not allow candidates to mutate evaluators, judges, suites, safety ceilings, authorization, or secrets.
- [ ] Do not add population search, Pareto selection, multi-component mutation, or automatic SQL generation.

These checkboxes are guardrails, not a backlog to complete for v1.
