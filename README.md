# ragopt

`ragopt` is the repository for a Go library and Glazed command-line application
for reproducible, evidence-driven optimization experiments over retrieval
systems.

The implementation is intentionally being designed from mechanisms already
proved in the GEC RAG and RAG-TTC work. See the active docmgr ticket under
`ttmp/` for scope, architecture, and implementation phases.

## Development

```bash
GOWORK=off go generate ./...
GOWORK=off go test ./...
GOWORK=off go run ./cmd/ragopt
```
