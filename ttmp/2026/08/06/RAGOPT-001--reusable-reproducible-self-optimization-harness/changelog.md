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
