---
Title: Nomic Retrieval Prefix Optimization Investigation
Ticket: RAGOPT-001
Status: active
Topics:
    - rag
    - evaluation
    - experiments
    - reproducibility
    - architecture
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: /home/manuel/code/gec/2026-03-16--gec-rag/data/knowledge-manifest-hybrid.yaml
      Note: Current GEC embedding provider, model, dimensions, and build endpoint
    - Path: /home/manuel/code/gec/2026-03-16--gec-rag/internal/knowledgebuild/build.go
      Note: Build-time raw and breadcrumb representation embedding path
    - Path: /home/manuel/code/gec/2026-03-16--gec-rag/internal/knowledgebuild/embed.go
      Note: Ollama embedder adapter and query-side reconstruction
    - Path: /home/manuel/code/gec/2026-03-16--gec-rag/internal/knowledge/service.go
      Note: Runtime lexical/vector retrieval and fusion path
ExternalSources:
    - https://huggingface.co/nomic-ai/nomic-embed-text-v1.5
    - https://arxiv.org/abs/2402.01613
    - https://ollama.com/library/nomic-embed-text
Summary: Evidence-gated design for testing search_document and search_query task prefixes without changing lexical retrieval or the chatbot.
LastUpdated: 2026-08-07T09:12:52.601520767-04:00
WhatFor: ""
WhenToUse: ""
---

# Nomic Retrieval Prefix Optimization Investigation

## Executive Summary

The frozen GEC hybrid index uses Ollama's `nomic-embed-text` at 768 dimensions
for both indexed representations and runtime queries. The application currently
passes raw text on both paths. Nomic's upstream model card for
`nomic-embed-text-v1.5` says retrieval inputs require asymmetric task prefixes:

```text
document = "search_document: " + passage
query    = "search_query: " + user_question
```

This is a credible retrieval-quality defect, but it is not yet a demonstrated
cause of the GEC comparison failure. We should test it as one immutable
embedding-policy mutation. The test must rebuild the vector index, leave BM25
and the chat runtime unchanged, and pass retrieval gates before consuming
answer or judge calls.

## Problem Statement

The current indexing path sends raw and breadcrumb representations directly to
the embedder in `internal/knowledgebuild/build.go`. The query path reconstructs
the same model identity from the bundle manifest, then sends the raw search
query through `internal/knowledge/service.go` and ragkit's vector index.

This makes document and query vectors dimensionally compatible, but it does not
give the instruction-tuned model their different retrieval roles. Nomic's model
card calls the task instruction prefix important and explicitly prescribes the
two prefixes above for RAG.

Three uncertainties prevent treating this as an automatic fix:

- The `nomic-embed-text` Ollama package must be verified against the upstream
  version and Modelfile; the server might apply a transform not visible in the
  application call site.
- Existing GEC retrieval benchmarks were run without application-side prefixes
  and already showed strong aggregate hybrid scores. The change could improve,
  preserve, or regress this corpus.
- The observed Morgan-versus-Peace failure also involved top-k document
  concentration. Better vector similarity does not guarantee evidence breadth.

The investigation question is therefore:

> With every other semantic input locked, does the Nomic-recommended asymmetric
> prefix policy improve frozen GEC retrieval, especially multi-document
> comparison coverage, without unacceptable group regressions?

## Proposed Solution

### Candidate boundary

Treat the prefix policy as a single replacement asset:

```yaml
embedding_transform:
  version: nomic-retrieval-prefix/v1
  document_prefix: "search_document: "
  query_prefix: "search_query: "
```

The incumbent policy has empty prefixes. The challenger supplies both. Corpus,
chunker, representations, model, dimensions, vector backend, lexical index,
fusion parameters, reranker setting, suite, and evaluator remain locked.

### Runtime behavior

```text
INDEX BUILD
  normalized document
    -> chunk
    -> raw representation ---------+
    -> breadcrumb representation --+-> prefix as document -> embed -> vector

QUERY
  user/tool query
    +-> unchanged text -> BM25
    +-> prefix as query -> embed -> vector
    -> unchanged weighted RRF -> unchanged result admission
```

Pseudocode:

```go
type EmbeddingTransform struct {
    Version        string
    DocumentPrefix string
    QueryPrefix    string
}

func (t EmbeddingTransform) Document(text string) string {
    return t.DocumentPrefix + text
}

func (t EmbeddingTransform) Query(text string) string {
    return t.QueryPrefix + text
}

// Build path: apply to every raw and breadcrumb representation.
vectors := EmbedAll(model, map(representations, transform.Document))

// Search path: never prefix lexical text.
lexicalHits := BM25.Search(query)
vectorHits  := Vector.Search(transform.Query(query))
hits        := WeightedRRF(lexicalHits, vectorHits)
```

The transform identity and exact strings must be recorded in the immutable
bundle manifest. A service must refuse to open a vector bundle if it cannot
reconstruct the declared query transform. This prevents prefixed documents
from being queried without a prefix, or vice versa.

### Evaluation sequence

1. Inspect the installed Ollama model metadata and perform a small embedding
   preflight to establish that the application is responsible for prefixes.
2. Freeze the existing retrieval suite and add no cases after seeing candidate
   results.
3. Add explicit multi-document coverage metrics to comparison cases. The
   current any-expected-document hit is insufficient.
4. Build incumbent and challenger bundles from the same frozen corpus. Reuse
   source extraction and chunk artifacts, but do not reuse old embedding-cache
   entries across different transform identities.
5. Run lexical, vector, and fused retrieval evaluation for both bundles.
6. Compare overall hit@k/MRR plus group and case deltas, document coverage,
   latency, embedding failures, and bundle identity.
7. Reject on a material protected-group regression. If retrieval does not
   improve, stop; do not run chatbot feedback.
8. Only after retrieval passes, run a small paired answer feedback suite under
   the existing provider budgets. Validation remains separately gated.

The comparison-case coverage metric should be explicit:

```text
required documents = {Morgan guide, Peace guide}
returned documents = unique(document_id for hit in top_k)
coverage@k         = |required ∩ returned| / |required|
complete@k         = coverage@k == 1.0
```

## Design Decisions

- **Prefix both sides.** Query-only prefixing would compare against vectors
  created under a different task transform and is not an admissible candidate.
- **Do not prefix BM25.** The string is a vector-model instruction, not corpus
  content or a lexical expansion term.
- **Rebuild the vector bundle.** Every representation vector changes. The old
  bundle cannot be mutated in place.
- **Identity includes transformation.** Provider/model/dimensions alone are no
  longer a sufficient semantic embedding identity.
- **Retrieval first.** This question can be answered cheaply and clearly before
  chat-planner and answer-generation variance enter the experiment.
- **No generic RAGOPT behavior.** GEC owns this model-specific build/query
  transform. RAGOPT only records the candidate, runs paired cells, and applies
  product-defined gates.
- **No automatic promotion.** A passing report produces a promotion plan; the
  product still builds, verifies, activates, and rolls back bundles.

## Alternatives Considered

- **Add only `search_query:` at runtime:** rejected because it creates an
  asymmetric deployment without rebuilding the indexed vectors as prescribed.
- **Assume Ollama handles task intent automatically:** rejected until verified;
  the current API call sends only model and raw text, while the upstream model
  card puts prefix responsibility on the caller.
- **Change embeddings, RRF weights, and document diversity together:** rejected
  because a multi-surface mutation cannot identify which mechanism moved the
  result.
- **Run answer/judge evaluation immediately:** deferred because retrieval
  metrics can reject the candidate without provider spending or generation
  noise.
- **Generalize task-prefix APIs into ragkit immediately:** deferred. First prove
  the policy in GEC; move a generic transform contract only after another
  consumer demonstrates the same requirement.

## Implementation Plan

### Track P1 — Verify serving semantics

- Record the exact Ollama model digest/version and inspect its Modelfile.
- Embed one fixed text raw and prefixed; verify the resulting vectors differ.
- Record the preflight as a private native artifact with no business data.

### Track P2 — Freeze experiment identity

- Add a typed embedding-transform identity to the GEC bundle manifest.
- Freeze the corpus, chunks, representations, model, dimensions, retrieval
  configuration, suite, evaluator, and gates.
- Add comparison document-coverage expectations before building the candidate.

### Track P3 — Implement product-owned transforms

- Apply `search_document: ` immediately before build-time embedding.
- Apply `search_query: ` immediately before query embedding.
- Keep lexical queries byte-identical.
- Namespace embedding cache entries by transform version and exact prefix.
- Reject missing or unsupported query transforms while opening the bundle.

### Track P4 — Retrieval proof

- Build fresh incumbent and challenger immutable bundles.
- Verify both bundles and their semantic identities.
- Run the frozen retrieval suite and compare channel, group, and coverage
  metrics plus latency.
- Write a paired rejection or promotion report.

### Track P5 — Conditional chatbot proof

- Proceed only if P4 passes.
- Run the feedback suite twice from fresh RAGOPT roots under declared budgets.
- Stop before validation unless feedback gates pass.

## Open Questions

- Which exact upstream `nomic-embed-text` revision does the installed Ollama
  model digest contain?
- Does the current Ollama Modelfile add any prefix or prompt template?
- What frozen regression threshold should protect exact-product and jargon
  groups while targeting comparison coverage?
- Should document-diversity become a separate later candidate if prefixing
  improves similarity but not comparison breadth?

## References

- [Nomic `nomic-embed-text-v1.5` model card](https://huggingface.co/nomic-ai/nomic-embed-text-v1.5) — authoritative task-prefix instructions and examples.
- [Nomic Embed technical report](https://arxiv.org/abs/2402.01613) — model design, training, and retrieval evaluation.
- [Ollama `nomic-embed-text` library entry](https://ollama.com/library/nomic-embed-text) — packaged model and embedding API.
- `data/knowledge-manifest-hybrid.yaml:31` — current GEC provider, endpoint,
  model, dimensions, and batching configuration.
- `internal/knowledgebuild/build.go:147` — raw and breadcrumb representation
  composition; `:170` begins build-time embedding.
- `internal/knowledgebuild/embed.go:36` — Geppetto-to-ragkit embedder adapter;
  `:74` reconstructs the query embedder.
- `internal/knowledge/service.go:47` — query embedder reconstruction from bundle
  identity; `:266` begins lexical/vector retrieval and fusion.
- `reference/09-first-gec-feedback-proof-and-source-role-candidate-rejection.md`
  — first paired feedback result that motivated this follow-up investigation.
