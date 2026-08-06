# Changelog

## 2026-08-06

- Initial workspace created


## 2026-08-06

Initialized ragopt from the committed Go template, validated the empty scaffold, audited proven RAG experiment mechanisms, and created the phased implementation ledger.

### Related Files

- /home/manuel/code/wesen/go-go-golems/ragopt/go.mod — Normalized project module baseline
- /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/tasks.md — Canonical phased implementation ledger

## 2026-08-06

Completed the intern-facing evidence-gated optimization harness design: evidence matrix, scope, schemas, APIs, algorithms, diagrams, integration boundaries, tests, decisions, and phase exit criteria.

### Related Files

- /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/design-doc/01-ragopt-intern-guide-to-a-reusable-evidence-gated-optimization-harness.md — Primary architecture and implementation guide
- /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/reference/01-implementation-diary.md — Chronological evidence and scope decisions

## 2026-08-06

Validated the complete ticket and delivered the default-layout four-document design bundle to reMarkable at /ai/2026/08/06/RAGOPT-001.

### Related Files

- /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/reference/01-implementation-diary.md — Validation and reMarkable delivery evidence

## 2026-08-06 - Phase 1 immutable run store

Implemented and exercised ragopt-run/v1 creation, copied inputs, atomic artifacts, synced JSONL, terminal states, strict read validation, and interruption recovery (commit 9f1ccb4).

### Related Files

- /home/manuel/code/wesen/go-go-golems/ragopt/pkg/runstore/read.go — Strict run integrity validation
- /home/manuel/code/wesen/go-go-golems/ragopt/pkg/runstore/run.go — Phase 1 lifecycle checkpoint
- /home/manuel/code/wesen/go-go-golems/ragopt/pkg/runstore/run_test.go — Recovery and invariant tests

## 2026-08-06 - Phase 2 snapshot and candidate validation

Implemented and exercised strict ragopt-snapshot/v1 identity, ragopt-candidate/v1 exactly-one-mutation validation, safe bundle resolution, and the narrow Glazed candidate validate command (commit d2329dd).

### Related Files

- /home/manuel/code/wesen/go-go-golems/ragopt/cmd/ragopt/commands/candidate/validate.go — Artifact CLI checkpoint
- /home/manuel/code/wesen/go-go-golems/ragopt/pkg/candidate/candidate.go — Phase 2 validation checkpoint
- /home/manuel/code/wesen/go-go-golems/ragopt/pkg/candidate/candidate_test.go — Mutation and corruption fixtures

## 2026-08-06 - Phase 3 resumable paired evaluation

Implemented and exercised ragopt-suite/v1, explicit incumbent/challenger arms, complete immutable input binding, synced ragopt-cell/v1 results, native artifact custody, failed-cell retention, and exact active-run resume (commit 0802670).

### Related Files

- /home/manuel/code/wesen/go-go-golems/ragopt/pkg/eval/resume.go — Resume and corruption semantics
- /home/manuel/code/wesen/go-go-golems/ragopt/pkg/eval/runner.go — Paired runner checkpoint
- /home/manuel/code/wesen/go-go-golems/ragopt/pkg/eval/runner_test.go — Interruption and equivalence proof

## 2026-08-06 - Phase 4 paired decisions and promotion reports

Implemented and exercised strict artifact loading, exact incumbent/candidate joins, transparent denominators, ragopt-gate-policy/v1 lexicographic decisions, deterministic Markdown reviews, non-applying promotion plans, and Glazed compare/report commands (commit 7fc1a98).

### Related Files

- cmd/ragopt/commands/report/report.go — Explicit non-applying report CLI
- pkg/compare/build.go — Exact pairing and aggregate accounting
- pkg/gate/evaluate.go — Lexicographic gate decision checkpoint
- pkg/report/render.go — Deterministic promotion review

## 2026-08-06 - Implementation bundle re-delivery and Phase 5 boundary

Replaced the earlier reMarkable design PDF with the expanded eight-document default-layout implementation bundle. Scoped the existing RAG-TTC per-query tool-loop, judge, and human-authored combined-comparison search-description seam; product integration awaits publication of the local ragopt commits.

### Related Files

- ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/reference/01-implementation-diary.md — Upload evidence and RAG-TTC dependency boundary

## 2026-08-06 - Phase 5 RAG-TTC proof-cycle inputs

Selected and strictly validated the human-authored I5 combined-comparison search-description candidate, 3-case feedback suite, disjoint 7-case validation suite, RAG-TTC gate policy, and locked model/corpus/index/judge/prompt/safety/evaluator identities. Committed the product-owned assets and drift tests in RAG-TTC as a6cfd93.

### Related Files

- abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/rag-ttc/configs/ragopt/i5-combined-comparison-v1/candidate.yaml — Validated one-mutation product candidate
- abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/rag-ttc/internal/ragoptassets/assets_test.go — Product asset and suite drift tests

## 2026-08-06 - Product-owned RAG-TTC adapter

Published ragopt main, pinned the portable module revision in RAG-TTC, and implemented the bounded native I5 proof-cycle command. The adapter executes the existing chat runtime and judge, preserves complete per-cell sessions, locks runtime/source/index/profile identities, supports exact active-run resume, and introduces no generic plugin or subprocess protocol. RAG-TTC commit: `56432fd`.

### Related Files

- abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/rag-ttc/cmd/rag-ttc/cmds/chat/tooleval/ragopt.go — Product-owned execution and native artifact adapter
- abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/rag-ttc/cmd/rag-ttc/cmds/chat/tooleval/ragopt_test.go — Locked environment and arm-binding tests
- ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/reference/01-implementation-diary.md — Publication, implementation, failure, and validation trail

## 2026-08-06 - Feedback proof run awaiting explicit provider approval

Verified the exact six-cell feedback command and tmux run procedure. The live start was rejected before execution because it would export evaluation questions and retrieved evidence to external answer/judge providers and may incur spend. No provider calls or run artifacts were created.

Refreshed the eight-document reMarkable bundle at `/ai/2026/08/06/RAGOPT-001`; the dry run explicitly confirmed `layout=default` before replacement.

## 2026-08-06 - First live feedback run rejected; budget bias corrected

Completed and retained the first six-cell feedback run, then rejected it after native evidence showed that a one-embedding budget biased execution toward the combined-query candidate. RAG-TTC commit `90485d8` locks sufficient 3/4/2 per-cell budgets and excludes judge overhead from product cost accounting. Focused tests and strict one-mutation validation pass; a fresh proof run is required.

## 2026-08-06 - Corrected RAG-TTC feedback proof formally rejected

Ran the corrected I5 candidate in an isolated worktree: 6/6 cells and all native artifact links validated, but only one candidate answer was contract-valid and the incumbent had one real iteration-ceiling failure. `ragopt compare` returned `fail`; the non-applying report says not to promote. Validation was intentionally skipped after feedback hard-gate failure.

### Related Files

- ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/reference/06-first-rag-ttc-i5-feedback-proof-rejection-report.md — Identities, six cell outcomes, gates, commands, and decision
- ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/reference/01-implementation-diary.md — Isolated setup, failures, commands, and lessons

Replaced the reMarkable bundle with the nine-document default-layout edition containing the formal rejection report.

## 2026-08-06

Step 15: clarified the product-proof boundary, recorded reproducible rejection as valid harness evidence, and preserved the feedback-before-validation spending gate; the fresh provider rerun awaits explicit approval (commit `15f5757`).

### Related Files

- /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/tasks.md — Phase 5 scope and task semantics

## 2026-08-06

Step 16: designed the shared production-refresh control plane for CoinVault/GEC, TTC Garden, and RAG-TTC Admin; added post-v0.1 implementation tasks and uploaded the standalone default-layout design to reMarkable (commit 8eaa35f).

### Related Files

- /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/design-doc/02-production-index-build-scheduling-resumability-and-ragopt-integration.md — Shared lifecycle, APIs, AWS/River decisions, product adapters, and implementation phases
- /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/tasks.md — Post-v0.1 shared refresh and product deployment ledger

## 2026-08-06

Step 17: completed no-provider readiness checks for the second RAG-TTC proof, distinguished policy/suite byte and semantic identities, and recorded the pre-v0.1 schema cleanup (commit 029c076); the six-cell provider run still awaits explicit approval.

### Related Files

- /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/reference/06-first-rag-ttc-i5-feedback-proof-rejection-report.md — Corrected first-run policy and suite identity layers
- /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/reference/07-rag-ttc-i5-fresh-root-reproducibility-protocol.md — Authoritative task 67 identities, readiness evidence, and comparison criteria

## 2026-08-06

Step 18: froze the six-cell provider consent envelope, including exact questions, payload classes, cache behavior, provider identities, 18/24/12 operation ceilings, local retention, and the validation exclusion (commit c71db2c); no provider run was launched.

### Related Files

- /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/reference/07-rag-ttc-i5-fresh-root-reproducibility-protocol.md — Auditable external execution authorization scope for task 67

## 2026-08-06

Step 19: recorded task 67 as externally blocked after three consecutive continuations without the explicitly bounded six-cell provider authorization; no second run, provider call, or validation cell exists (commit ed2f4a9).

### Related Files

- /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/index.md — Current blocker and exact resume condition
- /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/reference/01-implementation-diary.md — Strict blocker audit and preserved execution boundary

## 2026-08-06

Steps 20–23: ran the authorized six-cell RAG-TTC I5 feedback proof from a fresh root, verified frozen identities and complete artifact custody, reproduced the `fail` decision within 22/17/8 answer/search/judge calls and 108,548 provider tokens, and kept validation closed. Checked task 67, added the formal reproducibility report, and printed the RAGOPT brutalist mission board plus per-diary closing tickets.

### Related Files

- /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/reference/08-second-rag-ttc-i5-fresh-root-proof-and-reproducibility-decision.md — Fresh-root identities, cells, custody, budgets, stochastic interpretation, and phase decision
- /home/manuel/code/wesen/go-go-golems/ragopt/ttmp/2026/08/06/RAGOPT-001--reusable-reproducible-self-optimization-harness/sources/almanach/21-ragopt-mission-board.yaml — Reprintable physical project mission board
