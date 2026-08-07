---
Title: GEC CoinVault RAG End-to-End Failure Investigation and Pragmatic Recovery Guide
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
    - Path: /tmp/gec-ragopt-phase5/internal/knowledgebuild/connectors.go
      Note: MySQL and curated SQL-document extraction into normalized knowledge documents
    - Path: /tmp/gec-ragopt-phase5/internal/knowledgebuild/build.go
      Note: Corpus, chunk, representation, embedding, and immutable bundle construction
    - Path: /tmp/gec-ragopt-phase5/internal/knowledge/service.go
      Note: Authorized lexical, vector, fusion, rerank, and result-materialization runtime
    - Path: /tmp/gec-ragopt-phase5/internal/knowledge/tool.go
      Note: Production knowledge_search tool schema and execution
    - Path: /tmp/gec-ragopt-phase5/internal/knowledge/evidence.go
      Note: Per-session evidence admission and stable evidence identifiers
    - Path: /tmp/gec-ragopt-phase5/internal/webchat/runtime_prompts.go
      Note: Answer, citation, source, and self-grade instructions
    - Path: /tmp/gec-ragopt-phase5/internal/webchat/coinvault_projection_feature.go
      Note: Answer projection, evidence resolution, and projection-error events
    - Path: /tmp/gec-ragopt-phase5/internal/knowledge/judge.go
      Note: Statement decomposition, entailment verdicts, relevance, and faithfulness
    - Path: /tmp/gec-ragopt-phase5/cmd/coinvault/cmds/knowledge_ragopt.go
      Note: Product-owned paired adapter, budgets, native artifacts, and RAGOPT outcomes
    - Path: /tmp/gec-ragopt-phase5/cmd/coinvault/cmds/knowledge_ragopt_trace.go
      Note: Native runtime observation and source-role metric projection
ExternalSources:
    - https://huggingface.co/nomic-ai/nomic-embed-text-v1.5
    - https://arxiv.org/abs/2402.01613
Summary: Intern-facing investigation of the GEC corpus, index, hybrid retrieval, tool loop, evidence contract, judge, RAGOPT proof, observed failure mechanisms, and ordered recovery design.
LastUpdated: 2026-08-07T09:17:22.133381905-04:00
WhatFor: ""
WhenToUse: ""
---

# GEC CoinVault RAG End-to-End Failure Investigation and Pragmatic Recovery Guide

## Executive Summary

This report investigates the first paired RAGOPT feedback run against the GEC
CoinVault knowledge and SQL runtime. It is written for an intern who needs to
understand the product, corpus, immutable index bundle, hybrid retrieval,
Geppetto tool loop, evidence contract, answer projection, judge, and RAGOPT
promotion boundary before changing any of them.

The short conclusion is:

> The evidence does not show that the whole GEC RAG system is broken. It shows
> that the current system performs well on schema discovery and reasonably on
> an exact-product-plus-explanation question, while a multi-document historical
> comparison fails because the admitted evidence lacks document breadth and the
> answer model fills the gaps. The current evaluator then obscures the diagnosis
> by conflating tool routing, knowledge roles, document coverage, and contract
> validity.

The first candidate changed only the prose description of `knowledge_search`.
It asked the model to choose `schema_doc`, `product`, or `guide` roles more
explicitly. The candidate was correctly rejected:

- all six paired feedback cells completed;
- the challenger lost relevance, faithfulness, and source-role score;
- its worst faithfulness was `0.45`, below the hard `0.80` floor;
- validation was correctly left unrun;
- the raw native evidence remained outside the documentation repository.

The most important confirmed problems are:

1. **Comparison evidence lacks document breadth.** Four of the challenger's
   five chunks came from one general guide. The specific Morgan and Peace guides
   existed in the bundle but did not enter admitted evidence.
2. **Chunk-level collapse is not document-level diversity.** Raw and breadcrumb
   representation hits collapse to chunks, but multiple chunks from one
   document can still consume the result budget.
3. **The generator is insufficiently constrained by admitted evidence.** It
   supplied plausible facts from model memory and labeled itself measured even
   when the external judge found only 9 of 20 factual statements supported.
4. **The route metric is structurally wrong for a multi-tool agent.**
   `source_role_match` looks only at `knowledge_search` roles and ignores facts
   already obtained from `sql_doc` or `sql_query`.
5. **The answer contract is observed but not enforced.** Projection errors,
   invalid evidence citations, malformed UI blocks, and tool failures do not
   necessarily make a RAGOPT cell contract-invalid.
6. **Comparison retrieval evaluation measures any-document hit, not required
   breadth.** A broad category guide counts as success even when the two
   specific comparison documents are absent.
7. **The Nomic embedding input policy deserves a separate experiment.** The
   code sends raw document and query text to `nomic-embed-text`; Nomic's model
   card requires `search_document:` and `search_query:` prefixes. This is a
   strong hypothesis, not yet a proven root cause.

The pragmatic recovery order is measurement first, then one targeted retrieval
mutation, then generation grounding. Do not begin by changing the embedding
model, adding an LLM reranker, rewriting the tool loop, or building an autonomous
optimizer.

## 1. Scope and System Identity

### 1.1 Which chatbot is this?

This report is about the **GEC CoinVault Admin/Analyst chatbot**. It supports
backend commerce and logistics questions over GEC product/catalog knowledge and
live relational data. It is not the TTC customer-facing Garden Assistant with
small plant-comparison cards.

The two products may eventually share ragkit retrieval mechanisms and RAGOPT
experiment infrastructure, but they have different source systems, evidence
needs, authorization policies, response contracts, and evaluation suites.

### 1.2 Repository and revision boundary

The executable product proof was frozen in the isolated worktree:

```text
/tmp/gec-ragopt-phase5
branch: codex/ragopt-phase5-gec
revision: e2d1997151abf17a779559acbf8d0eea976d5504
```

The generic experiment harness is:

```text
/home/manuel/code/wesen/go-go-golems/ragopt
ticket: RAGOPT-001
```

The native first-run artifacts are private temporary evidence under:

```text
/tmp/gec-ragopt-feedback-proof-1/
  20260807T012544.623563776Z-gec-source-role-routing-feedback-9aaf797db8c0/
```

They are not committed because prompts, SQL-derived context, retrieval chunks,
answers, and judge inputs can contain business-sensitive material.

### 1.3 What is locked, and what changed?

The proof compares two complete snapshots over the same frozen suite. Locked
inputs include:

- GEC source revision;
- RAGOPT revision;
- corpus and immutable bundle identity;
- lexical and vector indexes;
- embedding identity;
- answer and judge profiles;
- tool catalog and authorization behavior;
- response prompt and evaluator;
- feedback suite and gate policy;
- operation and token ceilings.

The only candidate mutation was a complete replacement of the
`knowledge_search` tool description. The incumbent described available roles
generically. The challenger instructed the model to use:

- `schema_doc` for schema or column questions;
- `product` for exact item facts;
- `guide` for explanatory, historical, or comparison questions.

It did not change search code, index content, top-k, fusion, reranking,
embeddings, answer prompts, judge prompts, or the shared tool loop.

## 2. Mental Model and Vocabulary

An intern should keep these objects separate:

| Term | Meaning | Owned by |
|---|---|---|
| Source row | A MySQL product/category record | GEC database |
| Curated SQL document | Repository-authored schema/analyst documentation | GEC product |
| Knowledge document | Stable normalized text plus metadata and source role | GEC build pipeline |
| Chunk | Heading-aware section fragment from a document | ragkit/GEC build |
| Representation | Text embedded or indexed for one chunk; raw or breadcrumb | ragkit/GEC build |
| Bundle | Immutable content-addressed corpus, chunks, indexes, and manifest | GEC build/runtime |
| Retrieval hit | Ranked representation or chunk identity from a channel | ragkit runtime |
| Evidence item | Authorized material admitted to one chat session as `[E#]` | GEC runtime |
| Tool route | Sequence of `sql_doc`, `sql_query`, and `knowledge_search` calls | Geppetto + model |
| Native trace | Product-specific record of the real turn | GEC adapter |
| Outcome | Small common projection used for paired comparison | RAGOPT |
| Candidate | One immutable mutable-asset replacement plus hypothesis | RAGOPT contract |
| Gate | Product-defined deterministic promotion condition | GEC policy + RAGOPT |

These distinctions matter because a relevant document can exist in the corpus
yet fail to become a retrieved chunk; a retrieved chunk can fail authorization;
an authorized chunk can fail to be admitted; admitted evidence can fail to
support a generated claim; and an evaluator can fail to notice a malformed
answer even when the runtime recorded the error.

## 3. End-to-End Architecture

```text
                         OFFLINE BUILD

  MySQL products/categories       Curated SQL documentation
             |                               |
             +----------- connectors --------+
                             |
                     normalized documents
                  stable ID + role + digest
                             |
              optional repeated-furniture removal
                             |
                   Markdown heading chunker
                  1600 runes / 120 overlap
                             |
             +---------------+----------------+
             |                                |
      raw representation              breadcrumb representation
             |                                |
             +---------- lexical text --------+
             +---------- vector text ---------+
                             |
                   BM25 + exact SQLite vector
                             |
             immutable content-addressed bundle


                        ONLINE CHAT TURN

  Admin question
       |
  Geppetto tool loop, auto tool choice, maximum 20 iterations
       |
       +---- sql_doc --------> curated schema knowledge
       +---- sql_query ------> live SELECT-only bounded SQL
       +---- knowledge_search
                    |
          lexical BM25 ranking
          vector query embedding + ranking
                    |
          weighted reciprocal-rank fusion
          optional reranker (disabled here)
                    |
          authorization + chunk materialization
                    |
          evidence ledger assigns [E1], [E2], ...
       |
  model writes final structured answer
       |
  projection resolves citations and UI blocks
       |
  websocket/session presentation


                    OFFLINE PAIRED EVALUATION

  frozen feedback case
       |
  incumbent turn + challenger turn
       |
  native traces + product judge
       |
  common RAGOPT outcomes
       |
  strict pair join -> aggregate -> gates -> report
       |
  reject OR non-applying promotion plan
```

## 4. Corpus Construction

### 4.1 Product connector

`internal/knowledgebuild/connectors.go:51` performs full-source extraction for
active products. It joins or projects product details/facets, converts HTML to
plain normalized text, assigns stable identifiers such as `gec:product:123`,
and records content digest and `updated_at` metadata.

Product documents are assigned the `product` source role. They are appropriate
for durable item identity, descriptions, specifications, and category/facet
context.

Products with descriptions shorter than the configured minimum are excluded.
The hybrid manifest currently uses 200 runes.

### 4.2 Category/guide connector

`internal/knowledgebuild/connectors.go:190` extracts active categories with
substantive descriptions. These descriptions contain long-form collecting
guides, category history, terminology, designers, dates, and other explanatory
material.

Category documents are assigned the `guide` source role. The hybrid manifest
uses a 120-rune minimum.

The Morgan and Peace evidence needed by the failed comparison exists here:

- `gec:category:69` contains the Morgan guide, including George T. Morgan and
  the Liberty/eagle design;
- `gec:category:70` contains the Peace guide, including Anthony de Francisci,
  dates, obverse/reverse design, and historical context;
- `gec:category:8` is a broader silver-dollar guide that contains portions of
  both topics but is not a complete substitute for the two focused guides.

### 4.3 Curated SQL documents

`internal/knowledgebuild/connectors.go:258` turns repository-authored analyst
documentation into `schema_doc` knowledge documents. Separately, the runtime
also exposes the `sql_doc` tool. Do not confuse the knowledge role with the tool
name: one is metadata used to filter retrieval; the other is an executable tool
route.

### 4.4 Facts intentionally excluded from the vector corpus

Price, cost, and inventory quantity are deliberately not indexed. They remain
live SQL facts because they change frequently and need current access controls.

That creates a deliberate evidence split:

```text
durable descriptive fact -> immutable knowledge bundle
volatile operational fact -> live SQL query
```

A good Admin Chat answer often requires both. An exact product's current stock
and price may come from SQL, while the explanation of grade terminology comes
from a guide.

### 4.5 Current extraction limitations

- Extraction is a deterministic full scan; `updated_at` is recorded but not a
  watermark.
- A deleted source disappears only after a complete new bundle is built.
- There is no `cms_entries` connector yet.
- Production scheduling, publication, atomic activation, and hot reload are
  separate missing operational work; see
  `design-doc/02-production-index-build-scheduling-resumability-and-ragopt-integration.md`.

These are real production gaps, but the first feedback failure used a frozen,
verified bundle containing the needed documents. Freshness was not the cause of
that failure.

## 5. Chunking and Representations

### 5.1 Chunker configuration

`internal/knowledgebuild/build.go:138` invokes the heading-aware Markdown
chunker. The frozen manifest uses:

```yaml
chunking:
  max_section_runes: 1600
  min_section_runes: 200
  overlap_runes: 120
```

Heading awareness preserves document structure better than fixed byte windows,
but long sections still split. In the comparison trace, some chunks begin in
the middle of sentences, for example with fragments equivalent to:

```text
"was an assistant engraver ..."
"the obverse of the Peace dollar ..."
```

The missing subject can make a true fragment weak as standalone evidence. The
breadcrumb representation adds context for retrieval, but the materialized
answer evidence still needs enough local text to support a claim.

### 5.2 Two representations per chunk

`internal/knowledgebuild/build.go:147` creates raw representations.
`internal/knowledgebuild/build.go:151` creates breadcrumb representations.
They are composed at line 155.

Conceptually:

```text
raw:
  <chunk text>

breadcrumb:
  <document title> > <section path>
  <chunk text>
```

Both representations point back to the same chunk. This improves recall because
one channel can match exact prose while the other can match document/section
context.

The frozen bundle contains:

```text
documents:       16,032
chunks:          44,175
representations: 88,350
```

### 5.3 Embeddings

The hybrid manifest declares:

```yaml
embedding:
  enabled: true
  provider: ollama
  base_url: http://127.0.0.1:11435
  model: nomic-embed-text
  dimensions: 768
  batch_size: 16
```

`internal/knowledgebuild/build.go:170` constructs the embedder and wraps it in a
durable file cache. Four workers embed raw and breadcrumb representations in
batches. Unchanged representation text is a cache hit on later builds.

At runtime, `internal/knowledge/service.go:47` reads the vector identity from
the bundle manifest and reconstructs a query embedder through
`internal/knowledgebuild/embed.go:74`. Runtime endpoint selection uses
`COINVAULT_EMBEDDING_BASE_URL`, falling back to local Ollama on port 11434.

Provider, model, and dimensions are checked. The current identity does not
record an input task transform.

### 5.4 Prefix hypothesis

Nomic's upstream model card says RAG inputs should be transformed as:

```text
search_document: <indexed passage>
search_query: <user question>
```

The GEC code currently passes raw strings on both paths. That is preserved as a
follow-up optimization track in
`design-doc/03-nomic-retrieval-prefix-optimization-investigation.md`.

It must remain separate from this diagnosis:

- it is credible because it follows the model's authoritative instructions;
- it is not proven because the existing benchmark was strong overall and the
  failed case also shows an independent document-diversity problem;
- both document and query paths must change together;
- BM25 must remain unprefixed;
- the vector bundle must be rebuilt;
- transform identity must enter bundle identity and cache keys.

## 6. Immutable Bundle and Runtime Opening

The build passes documents, chunks, representations, chunker identity, vectors,
and embedding identity to ragkit's index-bundle builder. The resulting directory
is content-addressed and immutable.

At runtime `internal/knowledge/service.go:41`:

1. trims and validates the bundle path;
2. loads the manifest before opening indexes;
3. reconstructs the query embedder when a vector channel exists;
4. opens and verifies the bundle;
5. loads verified corpus documents;
6. creates document and chunk lookup maps;
7. records bundle ID and corpus root.

This is a good boundary. Query serving cannot silently point at arbitrary corpus
text after the index was built. The remaining production concern is activation:
CoinVault opens one verified bundle at startup and does not currently hot-reload
a newly published bundle.

## 7. Search Runtime

### 7.1 Tool input

`internal/knowledge/tool.go:17` defines the knowledge-search input. Important
fields include query text, limit, channel selection, and optional source-role
filters. `NewToolEntry` registers the tool with Geppetto. `runSearch` executes
the service call and admits returned chunks to the evidence ledger.

### 7.2 Search sequence

`internal/knowledge/service.go:199` is the public search boundary. In simplified
form:

```go
func Search(ctx, request) ([]Result, error) {
    validate(request)
    scopes, roles := authorize(request, caller)
    depth := request.Limit * searchDepth
    candidates := retrieve(ctx, request.Query, depth, request.Channels, scopes, roles)
    candidates = filterAndMaterialize(candidates, request.Limit)
    return resultsWithDocumentAndChunkMetadata(candidates)
}
```

The service over-fetches because authorization and role filtering can discard
ranked results.

### 7.3 Lexical channel

`internal/knowledge/service.go:317` sends the query to the bundle's lexical
index. Optional synonym expansion applies only here. The frozen proof had
synonyms disabled.

Representation hits are collapsed to one hit per chunk. The lexical score is
BM25.

### 7.4 Vector channel

`internal/knowledge/service.go:333` sends the raw query to the bundle's vector
index. The vector index embeds the query with the reconstructed query embedder,
performs exact SQLite vector search, and returns representation hits. Those are
also collapsed to one hit per chunk.

### 7.5 Fusion

`internal/knowledge/service.go:266` retrieves both channel rankings when the
bundle has vectors and the request is not lexical-only. It uses deterministic
weighted reciprocal-rank fusion:

```text
score(document) = sum over channels(weight / (rank_constant + rank))
```

The frozen defaults are:

```text
rank constant: 60
vector weight: 1.0
reranker:      disabled
synonyms:      disabled
```

Previous GEC experiments found the RRF parameter surface fairly flat and did
not promote tested rerankers or synonym variants. The first feedback proof does
not overturn those results.

### 7.6 Collapse is not diversity

Raw and breadcrumb hits collapse to chunks. That prevents two representations
of one chunk from appearing twice, but it does not prevent several chunks from
the same document from occupying top-k.

```text
representation diversity: handled
chunk diversity:          handled
document diversity:       not handled
```

This distinction directly explains the challenger comparison result: five
chunks looked like five evidence items, but four came from one broad guide.

## 8. Authorization and Evidence Admission

Search applies allowed scope and allowed source-role filters before evidence is
returned to the model. This is a security boundary, not merely a relevance
feature.

`internal/knowledge/evidence.go:17` defines the session evidence ledger.
`Admit` assigns stable identifiers such as `[E1]` and deduplicates by chunk
identity.

Conceptually:

```go
for result in authorizedSearchResults {
    if ledger.hasChunk(result.ChunkID) {
        reuse existing evidence ID
    } else {
        ledger.add(nextEvidenceID, result)
    }
}
```

Again, the ledger is chunk-aware, not document-budget-aware. That is reasonable
for citation identity but insufficient as a top-k coverage policy.

The runtime also has an evidence cache in
`internal/webchat/evidence_cache.go:21`, allowing projection to resolve cited
evidence IDs after tool execution.

## 9. The Three-Tool Agent

The `analyst-rag` profile exposes three complementary tools:

| Tool | Use | Typical authority |
|---|---|---|
| `sql_doc` | Discover tables, columns, relationships, and analyst guidance | Curated schema docs |
| `sql_query` | Execute bounded read-only SQL for current operational facts | Live database snapshot |
| `knowledge_search` | Retrieve durable product/category/guide evidence | Immutable bundle |

The shared Geppetto tool loop uses automatic tool choice with a maximum of 20
iterations. It is the same reusable tool-loop framework used elsewhere, but the
tool catalog, descriptions, authorization, state, and answer protocol are
product-owned. Therefore two products can share the loop implementation while
exhibiting different route behavior.

A good route depends on evidence already collected:

```text
Question: exact product price plus meaning of grade MS69

sql_doc   -> find relevant product/price schema
sql_query -> establish exact live product, price, inventory
knowledge_search(role=guide) -> explain MS69
answer    -> combine current fact and durable explanation
```

It would be wasteful to demand `knowledge_search(role=product)` merely because
the question mentions an exact product if SQL has already established the exact
product facts. This is why a static per-case expected knowledge role is too
coarse.

## 10. Answer and Citation Contract

`internal/webchat/runtime_prompts.go:129` tells the model how to use sources and
citations. Knowledge evidence must be cited using `[E#]`; SQL claims use the
tool result rather than inventing evidence IDs.

`internal/webchat/runtime_prompts.go:146` defines final answer metadata and the
self-grade. The model is asked to label an answer measured only when material
claims come directly from tools.

The first proof demonstrates that self-grade is telemetry, not a gate. The
comparison challenger labeled itself measured while the external judge found
11 of 20 factual statements unsupported.

`internal/webchat/coinvault_projection_feature.go:100` processes the final
answer into product events. Around lines 331–345 it resolves cited evidence IDs.
When a citation does not correspond to admitted `knowledge_search` evidence, it
emits `CoinVaultProjectionErrorProjected`.

That observability worked: the schema incumbent cited `[E1]` even though it had
no knowledge-search evidence, and the projection emitted:

```text
cited evidence ids E1 did not match any knowledge_search result in this run
```

The product evaluator failed to consume the event as a contract failure. The
challenger schema answer also used a malformed pills closing tag while remaining
contract-valid.

## 11. Native Judge

`internal/knowledge/judge.go` implements a decomposed judge in two provider
steps.

### 11.1 Statement extraction

The first prompt receives the user question and answer. It extracts factual
statements that should be checked. It does not see evidence at this stage.

### 11.2 Evidence verdicts and relevance

The second prompt receives:

- extracted statements;
- admitted knowledge chunks;
- bounded/truncated SQL results;
- the original question and answer context.

For each statement it decides whether the supplied evidence supports it and
records a reason. It also scores answer relevance.

Faithfulness is then computed deterministically:

```text
faithfulness = supported_statement_count / factual_statement_count
```

This makes the comparison failure interpretable:

```text
incumbent:  11 supported / 17 statements = 0.647...
challenger:  9 supported / 20 statements = 0.45
```

### 11.3 Judge artifact weakness

The native artifact projected into RAGOPT retains counts and scores but not the
statement text, individual verdict, or reason. Raw cache records may contain
more, but the primary run artifact should preserve a private statement ledger:

```go
type StatementVerdict struct {
    StatementID string   `json:"statement_id"`
    Text        string   `json:"text"`
    Supported   bool     `json:"supported"`
    EvidenceIDs []string `json:"evidence_ids,omitempty"`
    Reason      string   `json:"reason"`
}
```

This stays product-owned and private. RAGOPT needs only the aggregate metric and
the digest/path of the native artifact.

## 12. RAGOPT Boundary

RAGOPT does not retrieve documents or run the chat loop. It supplies reusable
experiment mechanics:

- immutable run directories;
- copied and hashed declared inputs;
- semantic identities;
- one-mutation candidate validation;
- paired incumbent/challenger scheduling;
- synced result-cell append and resume;
- common outcome projection;
- strict pairing and aggregate comparison;
- product-defined gates;
- deterministic reports and non-applying promotion plans.

The GEC adapter in `cmd/coinvault/cmds/knowledge_ragopt.go:135` composes the
real product runtime. Around line 328 it executes each cell. Around line 381 it
runs the product judge. Around line 413 it projects metrics. Around line 433 it
validates the native trace.

The common outcome deliberately does not copy every native event:

```go
type Outcome struct {
    Complete       bool
    ContractValid  bool
    Failed         bool
    Abstained      bool
    Metrics        map[string]float64
    Calls          int
    Tokens         int64
    Duration       time.Duration
    NativeArtifact ArtifactRef
    SemanticID     string
}
```

This separation is correct. The defect is not that RAGOPT lacks GEC semantics;
it is that the GEC adapter currently projects insufficient semantics into
`ContractValid` and its product metrics.

## 13. First Feedback Proof

### 13.1 Frozen run

```text
run ID:           20260807T012544.623563776Z-gec-source-role-routing-feedback-9aaf797db8c0
paired cases:     3
cells completed:  6 / 6
answer calls:     19 / 24
embedding calls:  4 / 18
judge calls:      12 / 12
answer tokens:    192,407 / 500,000
tool results:     13
validation cells: 0
decision:         fail
```

Embedding calls in this accounting are knowledge-search query embeddings, not
the already-built document index.

### 13.2 Aggregate result

| Metric | Incumbent | Challenger | Delta |
|---|---:|---:|---:|
| Relevance | 0.953333 | 0.866667 | -0.086667 |
| Faithfulness | 0.862745 | 0.800000 | -0.062745 |
| Source-role match | 0.666667 | 0.333333 | -0.333333 |

The aggregate challenger faithfulness of 0.80 hides the hard failure: the
minimum candidate cell was 0.45, below the required 0.80 floor. RAGOPT correctly
evaluated the hard condition before considering aggregate stories or cost.

## 14. Case Reconstruction

### 14.1 Orders-table/schema case

Both arms used:

```text
question
  -> sql_doc
  -> sql_query
  -> answer
```

Neither arm needed `knowledge_search`. Both achieved relevance 1.0 and
faithfulness 1.0.

What worked:

- schema discovery and live query planning were effective;
- the tool loop reached a correct, grounded answer;
- no vector retrieval was required.

What the metrics got wrong:

- `source_role_match` scored both arms zero because no knowledge role was used;
- this is not a route failure—the correct route was SQL-only.

What the contract got wrong:

- the incumbent cited `[E1]` without knowledge evidence;
- projection emitted an invalid-citation error;
- the GEC trace validator still marked the cell contract-valid;
- the challenger contained a malformed pills closing tag and also remained
  contract-valid.

Conclusion: the product answered the question well, while route and contract
metrics mischaracterized parts of the turn.

### 14.2 Exact-product/facet case

The question required a current exact product fact plus an explanation of a
grading facet.

Incumbent route:

```text
sql_doc
  -> sql_query attempt 1 fails
  -> sql_query attempt 2 succeeds
  -> knowledge_search(roles=[guide, product])
  -> answer
```

Challenger route:

```text
sql_doc
  -> sql_query succeeds
  -> knowledge_search(roles=[guide])
  -> answer
```

Results:

| Arm | Relevance | Faithfulness | Role match | Answer calls | Input tokens |
|---|---:|---:|---:|---:|---:|
| Incumbent | 1.00 | 0.941 | 1 | 5 | 57,374 |
| Challenger | 0.95 | 0.950 | 0 | 4 | 44,173 |

The challenger SQL result established the exact product, current price/stock,
and product ID. Guide evidence explained the grade. It used fewer calls and
tokens and marginally improved faithfulness.

The zero role score is misleading. The case expected `product`, but product
identity had already been grounded by SQL. A route-aware evaluator should ask:

- Did the route use an authoritative source for the exact product fact?
- Did it use guide evidence for the explanatory facet?
- Did the final claims cite or derive from those sources?

It should not demand redundant product retrieval independently of route state.

### 14.3 Morgan-versus-Peace comparison

Incumbent search:

```text
query:
  difference between Morgan and Peace silver dollars design dates history specifications
roles: [guide, product]
limit: 5
```

Challenger search:

```text
query:
  difference between Morgan and Peace silver dollars history design dates composition collecting comparison
roles: [guide]
limit: 5
```

The incumbent returned a broad guide plus four product documents. The
challenger returned four chunks from the broad silver-dollar guide and one from
another guide. The challenger did not admit the specific Morgan and Peace guide
documents despite both existing in the corpus.

```text
candidate result budget = 5 chunks

challenger:
  category 8, chunk A
  category 8, chunk B
  category 8, chunk C
  category 8, chunk D
  category 81, chunk E

needed for complete comparison:
  category 69, Morgan guide
  category 70, Peace guide
```

The current retrieval evaluation lists category 69, category 8, and category
70 as expected documents and counts any one as a hit. Category 8 therefore
makes the case look retrieved even though evidence breadth is incomplete.

The model then wrote details about designers, obverse/reverse imagery, dates,
composition, and collecting history. Some were true in the corpus but absent
from admitted evidence. The judge can only grade supplied evidence, so those
claims were unsupported.

| Arm | Relevance | Faithfulness | Supported statements | Knowledge document breadth |
|---|---:|---:|---:|---:|
| Incumbent | 0.86 | 0.647 | 11 / 17 | 5 docs / 5 chunks |
| Challenger | 0.65 | 0.450 | 9 / 20 | 2 docs / 5 chunks |

The causal chain is:

```text
role-description mutation
  -> guide-only search
  -> top-k concentrated in broad guide chunks
  -> specific Morgan and Peace guides absent
  -> answer model fills missing comparison facts
  -> unsupported statement count rises
  -> faithfulness floor fails
  -> RAGOPT rejects candidate
```

This is the strongest proven failure in the run.

## 15. Findings by Confidence

### 15.1 Confirmed by native evidence

| Finding | Evidence | Consequence |
|---|---|---|
| Specific comparison documents exist | Frozen corpus contains categories 69 and 70 | Not a source-extraction absence |
| Challenger evidence is document-concentrated | 4 of 5 chunks from category 8 | Top-k does not guarantee breadth |
| Generator asserts unsupported facts | Judge supports 9 of 20 statements | Grounding failure is real |
| Source-role metric ignores SQL route | Schema and exact-product traces | Metric cannot represent multi-tool correctness |
| Contract-valid ignores projection errors | Invalid `[E1]` citation remains valid | Gate can accept malformed protocol behavior |
| Any-document hit masks comparison incompleteness | Expected-doc logic accepts category 8 alone | Retrieval score is optimistic |
| Candidate description is a weak control | Different routes and worse quality from prose only | Candidate should remain rejected |

### 15.2 Strong hypotheses requiring experiments

| Hypothesis | Why credible | Required test |
|---|---|---|
| Nomic task prefixes improve vector retrieval | Upstream model card requires them; code omits them | Fresh paired vector bundles and frozen retrieval suite |
| Document diversity improves comparisons | Correct docs lose top-k to repeated broad-guide chunks | One per-document cap or diversity candidate |
| Better chunk self-containment improves entailment | Some chunks begin mid-sentence | Frozen rechunk candidate or evidence-context prefix |
| Stricter answer prompt reduces unsupported claims | Model fills gaps despite self-grade | Paired grounding prompt after retrieval is measured |

### 15.3 Not supported as current root causes

- database extraction failure;
- missing Morgan/Peace source content;
- stale prices in the vector corpus;
- global BM25 failure;
- global vector-model failure;
- RRF constant or vector weight;
- absence of a reranker;
- absence of synonym expansion;
- RAGOPT paired-run mechanics;
- sessionstream itself.

Those may deserve future work, but changing them now would not follow the
evidence from this proof.

## 16. Evaluator Redesign

The current evaluator needs separate metrics for separate questions.

### 16.1 Tool-route correctness

Did the agent use an allowed, authoritative route for each evidence need?

```go
type EvidenceNeed struct {
    ID             string
    FactClass      string
    AllowedTools   []string
    Required       bool
}

type RouteExpectation struct {
    Needs []EvidenceNeed
}
```

Example:

```yaml
needs:
  - id: current-product-fact
    allowed_tools: [sql_query]
    required: true
  - id: grade-explanation
    allowed_tools: [knowledge_search]
    allowed_roles: [guide]
    required: true
```

Metric:

```text
tool_route_recall = required needs satisfied by allowed tool / required needs
```

### 16.2 Knowledge-role correctness

For calls that actually use knowledge search, score allowed and required roles
without pretending that every case must call knowledge search.

```text
knowledge_role_precision = allowed requested roles / requested roles
knowledge_role_recall    = required roles observed / required roles
```

The metric should be nullable or not-applicable for SQL-only cases, not zero.

### 16.3 Document coverage

Comparison cases require a set or groups of documents:

```yaml
required_document_groups:
  - any_of: [gec:category:69]
  - any_of: [gec:category:70]
```

Metrics:

```text
document_coverage@k
complete_document_coverage@k
unique_documents@k
maximum_chunks_from_one_document@k
```

### 16.4 Citation validity and coverage

```text
citation_validity = resolved cited evidence IDs / all cited evidence IDs
citation_coverage = material supported claims with citations / material claims
```

SQL-grounded claims need a product-specific lineage marker or trace association;
they should not invent `[E#]` citations.

### 16.5 Answer faithfulness

Keep the decomposed judge, but preserve its statement ledger. Report:

- statement count;
- supported count;
- unsupported count;
- unsupported-claim rate;
- per-statement evidence IDs and reason;
- SQL truncation indicator;
- judge failure separately from a low score.

### 16.6 Contract validity

Replace the current minimal check with an explicit product contract:

```go
type ContractReport struct {
    TerminalFinished       bool
    SessionIdentityValid   bool
    ProviderObserved       bool
    RequiredBlocksValid    bool
    CitationsResolved      bool
    ProjectionErrors       []string
    ToolErrors             []string
    UnauthorizedEvents     []string
}

func (r ContractReport) Valid() bool {
    return r.TerminalFinished &&
        r.SessionIdentityValid &&
        r.ProviderObserved &&
        r.RequiredBlocksValid &&
        r.CitationsResolved &&
        len(r.ProjectionErrors) == 0 &&
        len(r.UnauthorizedEvents) == 0
}
```

Some tool errors may be recoverable if the model corrects them. The contract
must distinguish recovered tool errors from unresolved final failures rather
than treating every intermediate error identically.

## 17. Retrieval Trace Design

The native trace should make ranking loss visible without copying data into
RAGOPT.

```go
type RetrievalTrace struct {
    Query              string
    RequestedRoles     []string
    RequestedScopes    []string
    RequestedLimit     int
    Lexical            []RankedHit
    Vector             []RankedHit
    Fused              []RankedHit
    AfterAuthorization []RankedHit
    AfterRerank        []RankedHit
    Returned           []RankedHit
    FallbackUsed       bool
}

type RankedHit struct {
    Rank             int
    RepresentationID string
    ChunkID          string
    DocumentID       string
    SourceRole       string
    Score            float64
}
```

For each stage, preserve IDs, ranks, roles, scopes, and scores. Text can remain
in a separate private payload. This allows an investigator to ask:

- Did Morgan rank lexically but disappear in fusion?
- Did Peace rank in vector search but fail role authorization?
- Did two representations collapse correctly?
- Did one document dominate after collapse?
- Did reranking or fallback change the final set?

## 18. Pragmatic Recovery Plan

### Phase A — Make existing evidence trustworthy

Goal: diagnose the next candidate without manually reverse-engineering caches.

Tasks:

- preserve statement-level judge verdicts in the private native artifact;
- add retrieval-stage IDs/ranks/document roles to the native trace;
- make projection errors and invalid citations visible to contract validation;
- validate required final-answer blocks by query type;
- distinguish recovered tool errors from unresolved failures;
- replace `source_role_match` with route, role, coverage, and citation metrics;
- rerun only provider-free fixtures and existing artifacts while implementing
  measurement.

Exit criterion:

> A reviewer can explain why each cell passed or failed using one native run
> directory, without reading provider caches or reconstructing hidden rankings.

### Phase B — Grow a small stratified feedback suite

Do not jump from three cases to hundreds. Start with roughly 15–25 curated
feedback cases:

- SQL-only schema discovery;
- exact live product facts;
- retrieval-only exact guide questions;
- product-description retrieval;
- mixed SQL-plus-guide questions;
- multi-document comparisons;
- domain jargon/paraphrases;
- ambiguous entity names;
- appropriate abstention;
- authorization and protected-data boundaries.

Each case should declare:

```yaml
id: compare-morgan-peace
group: multi-document-comparison
question: ...
evidence_needs:
  - id: morgan-history
    allowed_tools: [knowledge_search]
    allowed_roles: [guide]
    required_documents: [gec:category:69]
  - id: peace-history
    allowed_tools: [knowledge_search]
    allowed_roles: [guide]
    required_documents: [gec:category:70]
protected_assertions:
  - no live customer data
```

Freeze a separate held-out validation set before evaluating candidates.

Exit criterion:

> Every important runtime mode has multiple cases, and comparison breadth can
> fail independently of generic hit@k.

### Phase C — Run one targeted retrieval candidate

There are two credible first candidates. Run them separately.

Candidate C1: Nomic task prefixes

- verify the installed Ollama model and Modelfile;
- record prefix transform in bundle identity and embedding cache namespace;
- prefix all document representations with `search_document: `;
- prefix vector queries with `search_query: `;
- leave lexical input unchanged;
- build fresh incumbent and challenger vector bundles;
- gate on frozen retrieval and document-coverage metrics.

Candidate C2: document-diverse evidence

- preserve the same rankings;
- admit at most a small configurable number of chunks per document, or use a
  simple document-diversification pass;
- do not add an LLM reranker;
- target multi-document comparison coverage;
- protect exact-product and single-guide groups from regression.

Suggested simplest policy:

```go
func DiverseTopK(hits []Hit, k, maxPerDocument int) []Hit {
    counts := map[string]int{}
    out := make([]Hit, 0, k)
    for _, hit := range hits {
        if counts[hit.DocumentID] >= maxPerDocument {
            continue
        }
        counts[hit.DocumentID]++
        out = append(out, hit)
        if len(out) == k {
            break
        }
    }
    return out
}
```

Start with `maxPerDocument=2` only as a candidate, not a default. The frozen
suite decides.

Exit criterion:

> One immutable candidate improves its declared retrieval target, passes all
> group regression limits, and reproduces from a fresh bundle/run root.

### Phase D — Tighten answer grounding

Only after retrieval evidence is adequate:

- tell the model not to complete missing facts from memory;
- require it to state when supplied evidence does not establish a requested
  distinction;
- require citations adjacent to material knowledge claims;
- keep SQL lineage separate from knowledge evidence IDs;
- evaluate unsupported-claim rate and citation coverage;
- avoid adding a second model-based verifier inside the production turn until
  offline evidence shows it is necessary.

Example instruction:

```text
Use only tool results for material factual claims. If the supplied evidence
does not establish a requested date, designer, specification, or comparison,
say that it is not established by the current sources. Do not complete the
answer from general model knowledge.
```

Exit criterion:

> When retrieval omits a fact, the answer omits or qualifies it rather than
> asserting it unsupported.

### Phase E — Paired feedback and conditional validation

- run incumbent and challenger over the frozen feedback suite;
- enforce operation and token ceilings;
- inspect every failed or regressed case;
- repeat from a fresh root to establish semantic reproducibility;
- run validation only if feedback gates pass;
- produce a non-applying promotion plan;
- activate through the product's bundle publication path, not RAGOPT.

## 19. What Not to Build Yet

Pragmatism is a design constraint. Do not add:

- a generic plugin or subprocess protocol to RAGOPT;
- an autonomous prompt/candidate generator;
- an LLM router in front of the existing LLM tool planner;
- an LLM reranker before cheap deterministic diversity is tested;
- a workflow engine inside RAGOPT;
- a cross-product universal quality threshold;
- automatic production activation;
- incremental CDC extraction before nightly deterministic rebuild cost is
  measured;
- a compatibility adapter for old bundle identities.

If a bundle schema changes before v0.1, make a clean versioned break. Do not add
an implicit fallback that guesses prefix or evaluator semantics.

## 20. Ownership Matrix

| Concern | GEC/CoinVault | ragkit | RAGOPT | Deployment |
|---|---|---|---|---|
| Source extraction | Owns | — | Records identity | Schedules job |
| Chunk/representation configuration | Chooses | Implements primitives | Locks candidate | — |
| Embedding task transform | Owns initially | May generalize after proof | Locks/compares | Provides endpoint |
| Lexical/vector retrieval | Configures | Implements | Observes metrics | Hosts bundle |
| Tool catalog and authorization | Owns | — | Locks identity | Secrets/network |
| Tool loop | Composes shared Geppetto | — | Observes calls | Hosts runtime |
| Answer protocol/projection | Owns | — | Receives contract flag | Serves UI |
| Judge and suite semantics | Owns | May provide utilities | Executes paired cells | Provider credentials |
| Run custody and gates | Supplies native artifacts/policy | — | Owns | Durable artifact store |
| Activation and rollback | Defines health | — | Produces plan only | Owns mutation |

The correct reuse boundary is mechanism, not meaning. ragkit can own generic
retrieval algorithms. RAGOPT can own generic experiment durability. GEC owns
what counts as an authoritative route and acceptable answer.

## 21. Testing Strategy

### 21.1 Unit tests

- representation collapse preserves best rank per chunk;
- document-diversity candidate preserves rank order among admitted hits;
- prefix transform applies exactly once;
- lexical query remains byte-identical;
- embedding cache keys change with transform identity;
- bundle open rejects missing/unsupported query transform;
- route metrics handle not-applicable knowledge search;
- comparison coverage requires every declared evidence group;
- contract validator rejects unresolved citations and projection errors;
- recovered SQL retry can remain valid when final evidence is correct.

### 21.2 Fixture integration tests

Build a tiny corpus containing:

- one broad comparison document;
- one focused Morgan document;
- one focused Peace document;
- multiple chunks from the broad document;
- one live-SQL fixture result.

Prove:

- chunk collapse does not imply document diversity;
- the coverage metric fails when only the broad guide is returned;
- the diversity candidate can admit both focused documents;
- an SQL-only case has route success and knowledge-role not-applicable;
- invalid `[E1]` without admitted evidence makes the contract invalid;
- paired interruption/resume produces the same canonical outcomes.

### 21.3 Frozen retrieval evaluation

Report at minimum:

- hit@1, hit@5, and MRR;
- document coverage@5 and complete coverage@5;
- unique documents@5;
- lexical, vector, and fused results;
- group deltas;
- query embedding latency and failures;
- exact semantic bundle IDs.

### 21.4 Provider-backed feedback

Provider evaluation is the final feedback stage, not the first debugging tool.
For each cell preserve:

- tool timeline;
- retrieval trace;
- evidence ledger;
- projection events;
- answer;
- statement verdicts;
- operation/token counts;
- common outcome and native artifact digest.

## 22. Operational Context

The production refresh design remains:

```text
nightly EventBridge schedule
  -> product refresh job
  -> consistent read-only snapshot/full scan
  -> deterministic documents/chunks
  -> cached embedding reuse
  -> new immutable bundle
  -> verify
  -> retrieval + product gates
  -> publish candidate bundle
  -> conditional active pointer update
  -> restart/reload CoinVault
  -> health check
  -> retain previous bundle for rollback
```

At approximately 16,000 source documents, a deterministic full rebuild with
incremental embedding-cache reuse is the pragmatic first production approach.
Do not introduce CDC until measurements show the full scan/build is the actual
bottleneck.

RAGOPT should supply build/run custody contracts and gates, not EventBridge,
AWS Batch, River, or deployment mutation. See the separate production refresh
design for resumability and registry interfaces.

## 23. Intern Review Path

Read in this order:

1. This report through the three case reconstructions.
2. `reference/09-first-gec-feedback-proof-and-source-role-candidate-rejection.md`
   for the sanitized proof record.
3. `internal/knowledgebuild/connectors.go:51` and `:190` for source extraction.
4. `internal/knowledgebuild/build.go:57`, `:138`, `:147`, and `:170` for bundle
   construction.
5. `internal/knowledge/service.go:199`, `:266`, `:317`, and `:333` for search.
6. `internal/knowledge/tool.go:74` and `internal/knowledge/evidence.go:50` for
   tool registration and evidence admission.
7. `internal/webchat/runtime_prompts.go:129` and
   `coinvault_projection_feature.go:331` for answer/citation behavior.
8. `internal/knowledge/judge.go:382` and `:435` for judge execution.
9. `cmd/coinvault/cmds/knowledge_ragopt.go:328`, `:413`, and `:433` for the
   product adapter and weak validation seam.
10. ragopt `pkg/eval`, `pkg/compare`, `pkg/gate`, and `pkg/report` for generic
    paired experiment mechanics.

Then reproduce the gate decision without provider calls:

```bash
cd /home/manuel/code/wesen/go-go-golems/ragopt
go run ./cmd/ragopt compare \
  --run /tmp/gec-ragopt-feedback-proof-1/20260807T012544.623563776Z-gec-source-role-routing-feedback-9aaf797db8c0 \
  --format json
```

Do not copy the raw result records into a ticket or public issue.

## 24. File and API Reference

### GEC source and build

- `internal/knowledgebuild/connectors.go:51` — product extraction.
- `internal/knowledgebuild/connectors.go:190` — category/guide extraction.
- `internal/knowledgebuild/connectors.go:258` — curated SQL documents.
- `internal/knowledgebuild/build.go:57` — build entry point.
- `internal/knowledgebuild/build.go:138` — heading-aware chunking.
- `internal/knowledgebuild/build.go:147` — raw representations.
- `internal/knowledgebuild/build.go:151` — breadcrumb representations.
- `internal/knowledgebuild/build.go:170` — cached embeddings.
- `internal/knowledgebuild/embed.go:36` — `rag.Embedder` adapter.
- `data/knowledge-manifest-hybrid.yaml:15` — frozen hybrid configuration.

### GEC retrieval and tools

- `internal/knowledge/service.go:22` — `Service` state.
- `internal/knowledge/service.go:41` — verified bundle open.
- `internal/knowledge/service.go:199` — public `Search`.
- `internal/knowledge/service.go:266` — retrieve and channel fusion.
- `internal/knowledge/service.go:317` — lexical ranking.
- `internal/knowledge/service.go:333` — vector ranking.
- `internal/knowledge/service.go:422` — weighted RRF helper.
- `internal/knowledge/service.go:440` — optional rerank path.
- `internal/knowledge/tool.go:17` — tool input.
- `internal/knowledge/tool.go:49` — tool configuration.
- `internal/knowledge/tool.go:74` — Geppetto registration.
- `internal/knowledge/tool.go:99` — tool search execution.
- `internal/knowledge/evidence.go:50` — evidence admission.

### GEC chat, projection, and evaluation

- `internal/webchat/runtime_prompts.go:129` — citation/source contract.
- `internal/webchat/runtime_prompts.go:146` — final answer metadata/self-grade.
- `internal/webchat/evidence_cache.go:21` — evidence lookup.
- `internal/webchat/coinvault_projection_feature.go:100` — projection feature.
- `internal/webchat/coinvault_projection_feature.go:331` — citation resolution.
- `internal/knowledge/judge.go:37` — statement prompt.
- `internal/knowledge/judge.go:51` — verdict prompt.
- `internal/knowledge/judge.go:382` — judge verdict execution.
- `internal/knowledge/judge.go:435` — answer judge.
- `cmd/coinvault/cmds/knowledge_ragopt.go:135` — adapter composition.
- `cmd/coinvault/cmds/knowledge_ragopt.go:328` — cell executor.
- `cmd/coinvault/cmds/knowledge_ragopt.go:381` — judge invocation.
- `cmd/coinvault/cmds/knowledge_ragopt.go:413` — metric projection.
- `cmd/coinvault/cmds/knowledge_ragopt.go:433` — native trace validation.
- `cmd/coinvault/cmds/knowledge_ragopt_trace.go:80` — trace observer.
- `cmd/coinvault/cmds/knowledge_ragopt_trace.go:205` — current role metric.

### Frozen candidate contracts

- `configs/ragopt/source-role-routing-v1/feedback-suite.json`
- `configs/ragopt/source-role-routing-v1/validation-suite.json`
- `configs/ragopt/source-role-routing-v1/gate-policy.yaml`
- `configs/ragopt/source-role-routing-v1/runtime-contract.yaml`
- `configs/ragopt/source-role-routing-v1/source-lock.yaml`
- `configs/ragopt/source-role-routing-v1/answer-contract.json`
- `configs/ragopt/source-role-routing-v1/adapter-contract.yaml`

The `answer-contract.json` file is copied and locked but is not currently
executed as a validator by the adapter. That gap must be resolved in Phase A.

## 25. Decision Record

### Decision 1: Keep the source-role candidate rejected

Reason: it breached the faithfulness floor, regressed aggregate quality, and
did not produce a reliable route improvement.

### Decision 2: Do not run validation for the rejected candidate

Reason: validation is a protected spending and information gate. Feedback did
not pass.

### Decision 3: Repair product measurement before tuning retrieval broadly

Reason: current route and contract metrics can label correct SQL-only behavior
wrong and malformed protocol behavior valid.

### Decision 4: Preserve Nomic prefixing as a separate candidate

Reason: it follows authoritative model instructions, changes the entire vector
space, and can be tested cheaply at retrieval level.

### Decision 5: Test deterministic document diversity before an LLM reranker

Reason: the observed failure is repeated chunks from one document. A small
diversity policy targets that mechanism directly and is cheaper, deterministic,
and easier to attribute.

### Decision 6: Keep product semantics out of RAGOPT

Reason: RAGOPT should make experiments real and reproducible; it should not
decide what GEC SQL, source roles, citations, or answer blocks mean.

## 26. Final Checklist Before the Next Provider Run

- [ ] Native retrieval trace records every ranking stage and document identity.
- [ ] Native judge artifact records statement-level verdicts and reasons.
- [ ] Contract validation consumes projection and citation errors.
- [ ] Route metrics distinguish SQL, knowledge role, and evidence coverage.
- [ ] Multi-document cases require complete evidence groups.
- [ ] Feedback and validation suites are frozen separately.
- [ ] Candidate changes exactly one semantic asset.
- [ ] Bundle identity includes every embedding transform.
- [ ] Operation and token budgets are declared.
- [ ] Provider data-flow authorization covers the actual tool surface.
- [ ] Raw artifacts have a private durable retention location.
- [ ] Validation cannot execute unless feedback passes.
- [ ] Promotion remains a plan, not an automatic deployment.

## 27. Closing Perspective

The failed source-role candidate is useful evidence. It proves the experiment
system can reject a plausible-sounding change and preserve enough native
artifacts to discover why. The failure is not an invitation to add more agentic
machinery. It is a request for sharper contracts:

```text
measure route truthfully
measure evidence breadth
enforce citation and projection contracts
retrieve the required documents
then constrain the answer to those documents
```

Once those foundations are in place, Nomic task prefixes and document diversity
are small, falsifiable candidates. That is the self-optimization approach we
want: one hypothesis, one mutation, frozen evidence, explicit budgets, hard
regression gates, retained failures, and no production change without a
reviewable promotion plan.

## References

- [Nomic `nomic-embed-text-v1.5` model card](https://huggingface.co/nomic-ai/nomic-embed-text-v1.5)
- [Nomic Embed technical report](https://arxiv.org/abs/2402.01613)
- `reference/09-first-gec-feedback-proof-and-source-role-candidate-rejection.md`
- `design-doc/02-production-index-build-scheduling-resumability-and-ragopt-integration.md`
- `design-doc/03-nomic-retrieval-prefix-optimization-investigation.md`
- GEC handoff:
  `ttmp/2026/08/05/GEC-RAG-OPT-001--retrieval-optimization-reranker-eval-growth-and-benchmarked-retrieval-experiments/reference/05-handoff-to-optimizer-open-tracks-review-guide-and-known-weaknesses.md`
