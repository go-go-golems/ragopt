---
Title: First GEC Feedback Proof and Source-Role Candidate Rejection
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
    - Path: abs:///tmp/gec-ragopt-phase5/cmd/coinvault/cmds/knowledge_ragopt.go
      Note: Product-owned RAGOPT composition, preflight, budgets, and arm execution
    - Path: abs:///tmp/gec-ragopt-phase5/cmd/coinvault/cmds/knowledge_ragopt_trace.go
      Note: Native trace projection into common outcomes
    - Path: abs:///tmp/gec-ragopt-phase5/configs/ragopt/source-role-routing-v1/candidate.yaml
      Note: Frozen one-mutation candidate identity and hypothesis
    - Path: abs:///tmp/gec-ragopt-phase5/configs/ragopt/source-role-routing-v1/shared/gate-policy.yaml
      Note: Feedback promotion gates and metric thresholds
    - Path: abs:///tmp/gec-ragopt-phase5/configs/ragopt/source-role-routing-v1/shared/runtime-contract.yaml
      Note: Frozen runtime and operation ceilings
    - Path: repo://pkg/compare/compare.go
      Note: Strict paired comparison and decision input
    - Path: repo://pkg/eval/runner.go
      Note: Common paired durable execution lifecycle
    - Path: repo://pkg/report/report.go
      Note: Non-applying promotion report projection
ExternalSources: []
Summary: The first production-shaped GEC feedback proof completed six valid cells within budget and rejected the source-role description candidate for lower relevance, lower role match, and a faithfulness-floor failure.
LastUpdated: 2026-08-06T21:30:23.110970378-04:00
WhatFor: Review the exact first-run identities, routes, metrics, custody evidence, rejection reason, and measured GEC integration friction.
WhenToUse: Read before repeating the GEC candidate, considering validation, changing the candidate, or generalizing any product adapter API into ragopt.
---


# First GEC Feedback Proof and Source-Role Candidate Rejection

## Goal

Document the first real GEC-RAG integration proof for RAGOPT-001 without
committing sensitive native answers or SQL tool results to the RAGOPT
repository. This report establishes what ran, what the candidate changed, how
the real Admin Chat loop behaved, why the candidate failed, and what remains
to prove on a fresh repeat.

The result is a successful harness proof and a failed product candidate. Those
are not contradictory: RAGOPT preserved identities, ran paired production
arms, stayed within budgets, retained native evidence, and produced a
deterministic rejection. The human-authored description mutation did not
improve the system.

## Context

### The evaluated system

Both arms used the same product composition:

```text
frozen feedback case
        |
        v
GEC local CoinVault session
        |
        v
Geppetto v0.13.7 tool loop
        |
        +--> sql_doc --------------------+
        +--> sql_query -> gec_dev -------+--> answer provider
        +--> knowledge_search -----------+        |
                         |                         v
                         +--> ragkit hybrid      final answer
                              retrieval             |
                                                    v
                                      decomposed cached judge
                                                    |
                                                    v
                                      GEC native outcome.json
                                                    |
                                                    v
                                      RAGOPT common cell + gate
```

`ragkit` supplied the hybrid retrieval/index abstractions. It did not supply
the chat loop. TTC and GEC share the Geppetto tool-loop implementation; GEC's
application profile selected the product-specific SQL and knowledge tools.
Pinocchio/sessionstream components supplied the production-shaped local chat
lifecycle. RAGOPT owned immutable inputs, pairing, append-and-sync result
custody, comparison, and the promotion decision.

### The only mutation

The challenger replaced the complete model-facing `knowledge_search`
description. It added explicit routing advice:

- use `schema_doc` for curated schema semantics;
- use `product` for exact catalog items and facets;
- use `guide` for explanatory, historical, terminology, and comparison work;
- omit roles only for genuinely cross-role questions.

It did not change the tool schema, retrieval code, index, corpus, model,
prompt, SQL policy, judge, suite, gate, or loop configuration.

### Frozen identities

| Identity | Value |
|---|---|
| Run ID | `20260807T012544.623563776Z-gec-source-role-routing-feedback-9aaf797db8c0` |
| Candidate ID | `gec-source-role-routing-v1` |
| Candidate digest | `sha256:20e0c3f65da9d3a5c8bdcebb0d9a727645acd269db8fbd50186a620205a6ace8` |
| Parent snapshot | `sha256:83ac2d1bfaba3e770732772f80bda969670146b8034a6dfb3031d1ebb0958402` |
| Child snapshot | `sha256:6579ee1c05c5fcf46fabd3433b4de8be17de291f3f647d2b9b67bcd56c178ee9` |
| Parent asset | `sha256:c64fba98eaff6111d558e267e85d48f40aa81ed710648a198b1d7e0d282c00ff` |
| Challenger asset | `sha256:ea73739ed4ff57215f278fdf2b52924c1157d72e76327821a521fd07a7ed8bb8` |
| Suite semantic digest | `sha256:555cfedd4bf2f8a4c7d645aebc1af77b6c8d7c391a684b098a30fb22416cd71c` |
| Policy copied-byte digest | `sha256:04a5b4262d552610fa025e68fac3c5a0a071a1e5409aec4715a62caae694dd66` |
| Policy semantic digest | `sha256:fda4fd08552badce0e2d1156e24d1e99d234914c44c2d06bae25b548f5492103` |
| RAGOPT revision | `4d410c57e2429e1109ba232ee04a610a3697eed0` |
| GEC source revision | `82c165921d7d42f06941f6edc1b4cd6673c1f58a` |
| GEC candidate commit | `e2d1997` |

The two policy digests name different identity layers: copied YAML bytes and
the normalized semantic policy object. RAGOPT's comparison explicitly passed
the `policy_bytes` check. Renaming these fields to eliminate ambiguity remains
a pre-v0.1 Phase 6 task; no compatibility shim should be added.

## Quick Reference

### Decision

> **FAIL — do not promote the source-role description candidate.**

The first failed hard gate was:

```text
candidate faithfulness minimum 0.450000 across 3 pairs; floor 0.800000
```

All custody and structural checks passed before that quality failure:

- policy bytes matched;
- all 3 expected pairs existed;
- all 6 cells completed;
- all 6 cells were contract-valid;
- no cell failed;
- no cell abstained;
- no validation case ran.

### Provider and operation usage

| Resource | Observed | Frozen ceiling | Result |
|---|---:|---:|---:|
| Answer-provider calls | 19 | 24 | PASS |
| Query-embedding requests | 4 | 18 | PASS |
| Judge-provider calls | 12 | 12 | PASS, exact ceiling |
| Answer input tokens | 188,157 | — | — |
| Answer output tokens | 4,250 | — | — |
| Answer tokens total | 192,407 | 500,000 | PASS |
| Tool results | 13 | observed, not budgeted | — |
| Judge cache hits / misses | 0 / 12 | observed | — |
| Wall time | 137.140 seconds | six cells | — |

Judge calls are evaluation overhead and are not included in the common RAGOPT
`provider_calls` cost field. Query embeddings ran through local Ollama.

### Aggregate outcome

| Metric | Incumbent mean | Challenger mean | Delta | Wins / ties / losses |
|---|---:|---:|---:|---:|
| Answer relevance | 0.953333 | 0.866667 | -0.086667 | 0 / 1 / 2 |
| Faithfulness | 0.862745 | 0.800000 | -0.062745 | 1 / 1 / 1 |
| Source-role match | 0.666667 | 0.333333 | -0.333333 | 0 / 2 / 1 |

### Paired case analysis

| Case | Expected role | Incumbent route | Challenger route | Relevance delta | Faithfulness delta | Role delta |
|---|---|---|---|---:|---:|---:|
| `orders-table` | `schema_doc` | `sql_doc` + `sql_query`; no knowledge call | same | 0.00 | 0.00 | 0 |
| `facet-sf-eagle-2011` | `product` | SQL plus `knowledge_search([guide,product])` | SQL plus `knowledge_search([guide])` | -0.05 | +0.008824 | -1 |
| `compare-morgan-peace` | `guide` | `knowledge_search([guide,product])` | `knowledge_search([guide])` | -0.21 | -0.197059 | 0 |

#### Schema case

Both arms correctly answered the table-column question through `sql_doc` and
live `information_schema` SQL. Neither called `knowledge_search`, so changing
its description could not affect the executed path. Both answers scored 1.0
for faithfulness and relevance and 0 for the suite's expected
`schema_doc`-through-knowledge metric.

This is not evidence that either answer was bad. It is evidence that the case
does not isolate the mutated tool description in a multi-tool production
profile. The candidate itself declared this risk. Future diagnostic suites
should distinguish "correct product tool route" from "expected
knowledge-search role" rather than treating SQL use as role failure.

#### Exact-product case

The incumbent searched both guide and product representations and retrieved
the exact product document. The challenger used SQL to establish the exact
live catalog item, then searched only `guide` for explanatory grade context.
That behavior is understandable in a multi-tool agent, but it directly
contradicts the challenger's instruction to use `product` for exact catalog
items and failed the expected-role metric.

The challenger saved one answer call, one tool call, and 13,274 tokens, and
slightly improved faithfulness. It nevertheless lost 0.05 relevance and did
not satisfy the declared routing hypothesis. This is a useful example of why
cost improvements cannot rescue a candidate that misses its quality target.

#### Guide-comparison case

The challenger followed the new description exactly and narrowed retrieval
from `[guide,product]` to `[guide]`. The narrower evidence did not constrain
the generated answer: the model added dates, designers, iconography, and
historical claims not fully supported by the admitted chunks. Faithfulness
fell from 0.647059 to 0.45 and relevance fell from 0.86 to 0.65.

The result disproves the simple hypothesis that narrowing retrieval roles is
enough to improve grounded answer quality. Retrieval scope, evidence coverage,
and generation behavior must be evaluated separately.

### Why the candidate failed

The mutation tried to influence a probabilistic planning decision through
description prose. The real loop demonstrated three limitations:

1. The model can choose a different tool entirely, making the mutated
   description irrelevant.
2. In a multi-tool answer, the semantically appropriate knowledge role depends
   on what SQL already established. A case-level expected role can be too
   coarse.
3. Correct role selection does not guarantee sufficient evidence coverage or
   faithful generation.

This does not imply that source roles are useless. It means the optimization
surface and metric must correspond to the behavior being changed. If the goal
is deterministic tool routing, prose-only tool descriptions are a weak
control. If the goal is final-answer quality, routing metrics are diagnostic,
not substitutes for faithfulness and relevance.

### Harness result versus product result

```text
RAGOPT harness: PASS
  identities frozen
  one mutation proven
  production adapter used
  all cells durable
  budgets enforced
  comparison deterministic
  rejection retained

GEC candidate: FAIL
  target relevance decreased
  faithfulness floor violated
  role match decreased
  no promotion
  no validation spending
```

### Measured integration friction

- Bleve's current Bolt-backed index open is exclusive. The live CoinVault
  server held the production bundle, so the proof used a byte-identical copy
  with the required `knowledge/{corpus,bundles/<id>}` topology. Locked corpus,
  lexical-manifest, and vector-SQLite digests proved equivalence.
- Candidate assets are intentionally opaque to RAGOPT. GEC product preflight
  caught malformed YAML that generic candidate validation could not interpret.
- The multi-tool product makes expected knowledge roles conditional on prior
  SQL results. A retrieval-only role metric cannot fully score planner quality.
- Judge budget reached exactly 12 calls because all six cells missed both
  decomposed cache entries. The ceiling was sufficient but had no margin.
- Raw and semantic policy digests remain separately useful but confusingly
  named in the current pre-v0.1 schemas.
- Native artifacts contain potentially sensitive answer and SQL context. They
  must remain in controlled product storage rather than being copied into the
  public RAGOPT repository.

## Usage Examples

### Reproduce the comparison without providers

```bash
cd /home/manuel/code/wesen/go-go-golems/ragopt

go run ./cmd/ragopt compare \
  --run /tmp/gec-ragopt-feedback-proof-1/20260807T012544.623563776Z-gec-source-role-routing-feedback-9aaf797db8c0 \
  --format json

go run ./cmd/ragopt report \
  --run /tmp/gec-ragopt-feedback-proof-1/20260807T012544.623563776Z-gec-source-role-routing-feedback-9aaf797db8c0 \
  --output-path /tmp/gec-ragopt-feedback-analysis-1/promotion-review.md \
  --plan-path /tmp/gec-ragopt-feedback-analysis-1/promotion-plan.json \
  --format json
```

These commands read completed artifacts only. They make no provider,
embedding, SQL, or validation call.

### Artifact integrity checks

| Artifact | SHA-256 |
|---|---|
| Run manifest | `86cbf4edf65a5f73c00e2b640e519c7a81b6709fd5aa667d2152635ca5efe886` |
| Run config | `a86e091c38e41d32cdbe0479a1996449cc891696ba51f9e24a23bf6a21652b0f` |
| Terminal status | `23b068c8f6ccd72cb29b4792c08e931e4da682cf90178ed299c9e6560650f54e` |
| Common cells JSONL | `33f5756f7a539e4481202b0fe561ac562044f6ec9d17e3461c64846c00e6d19b` |
| Run summary | `d53b0cfd211f25fc9d4b9c15bf773f336c3143d6d115cbe8dabcf05b16859774` |
| Promotion review | `f7ab025fe32486a50597b18d57b60c7e20b34ed4bbcfcabcf9ba4d998c74004b` |
| Promotion plan | `fbb4fbd7f1e6f5be4f846a706325af4b0d8e9d74dbdd1937b9ffea631a156dd1` |

### Fresh-repeat protocol

The first result is enough to reject promotion and keep validation closed. It
is not enough to finish Phase 5's reproducibility claim. A second execution
must:

1. use a new empty run root;
2. reuse the same candidate, suite, policy, snapshots, source lock, profiles,
   and physical index digests;
3. run feedback only;
4. preserve the same ceilings;
5. compare semantic identities before comparing stochastic metric values;
6. explain any decision change rather than requiring byte-identical model
   answers;
7. leave validation unrun regardless of the repeat because the first feedback
   gate already failed.

The Step 36 consent was recorded as specific to the first six-cell run. The
user subsequently approved one second identical six-cell run under the same
data-flow boundary and ceilings. That approval does not authorize validation.

## Related

- Candidate and preflight implementation:
  `/tmp/gec-ragopt-phase5/cmd/coinvault/cmds/knowledge_ragopt.go`
- Native trace projection:
  `/tmp/gec-ragopt-phase5/cmd/coinvault/cmds/knowledge_ragopt_trace.go`
- Frozen candidate:
  `/tmp/gec-ragopt-phase5/configs/ragopt/source-role-routing-v1/candidate.yaml`
- Runtime and budget contract:
  `/tmp/gec-ragopt-phase5/configs/ragopt/source-role-routing-v1/shared/runtime-contract.yaml`
- Gate policy:
  `/tmp/gec-ragopt-phase5/configs/ragopt/source-role-routing-v1/shared/gate-policy.yaml`
- RAGOPT pairing and comparison:
  `pkg/eval/runner.go`, `pkg/compare/compare.go`, and `pkg/report/report.go`
- Consent, launch, and implementation narrative:
  `reference/01-implementation-diary.md`, Steps 30–37
