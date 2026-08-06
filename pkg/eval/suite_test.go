package eval

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSuiteNormalizesGroupsAndOpaqueInput(t *testing.T) {
	path := filepath.Join(t.TempDir(), "suite.json")
	writeFile(t, path, []byte(`{
  "api_version": "ragopt-suite/v1",
  "name": "fixture-suite",
  "cases": [
    {"id":"a","groups":["z","a"],"input":{"z":2,"a":1}},
    {"id":"b","input":[true,3]}
  ]
}`))
	document, err := LoadSuite(t.Context(), path)
	mustNoError(t, err)
	if got := strings.Join(document.Suite.Cases[0].Groups, ","); got != "a,z" {
		t.Fatalf("groups were not normalized: %q", got)
	}
	if got := string(document.Suite.Cases[0].Input); got != `{"a":1,"z":2}` {
		t.Fatalf("opaque input was not canonicalized: %s", got)
	}
	if !strings.HasPrefix(document.Digest, "sha256:") || !filepath.IsAbs(document.SourcePath) {
		t.Fatalf("incomplete suite document: %#v", document)
	}
}

func TestLoadSuiteRejectsSchemaAndIdentityErrors(t *testing.T) {
	tests := map[string]string{
		"unknown field":   `{"api_version":"ragopt-suite/v1","name":"suite","unknown":true,"cases":[{"id":"a","input":{}}]}`,
		"duplicate case":  `{"api_version":"ragopt-suite/v1","name":"suite","cases":[{"id":"a","input":{}},{"id":"a","input":{}}]}`,
		"duplicate group": `{"api_version":"ragopt-suite/v1","name":"suite","cases":[{"id":"a","groups":["g","g"],"input":{}}]}`,
		"missing input":   `{"api_version":"ragopt-suite/v1","name":"suite","cases":[{"id":"a"}]}`,
		"multiple values": `{"api_version":"ragopt-suite/v1","name":"suite","cases":[{"id":"a","input":{}}]} {}`,
	}
	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "suite.json")
			writeFile(t, path, []byte(body))
			_, err := LoadSuite(t.Context(), path)
			if err == nil {
				t.Fatal("expected suite rejection")
			}
		})
	}
}

func TestSuiteDigestPreservesCaseOrder(t *testing.T) {
	root := t.TempDir()
	firstPath := filepath.Join(root, "first.json")
	secondPath := filepath.Join(root, "second.json")
	writeFile(t, firstPath, []byte(`{"api_version":"ragopt-suite/v1","name":"suite","cases":[{"id":"a","input":{}},{"id":"b","input":{}}]}`))
	writeFile(t, secondPath, []byte(`{"api_version":"ragopt-suite/v1","name":"suite","cases":[{"id":"b","input":{}},{"id":"a","input":{}}]}`))
	first, err := LoadSuite(t.Context(), firstPath)
	mustNoError(t, err)
	second, err := LoadSuite(t.Context(), secondPath)
	mustNoError(t, err)
	if first.Digest == second.Digest {
		t.Fatal("suite digest ignored declared case order")
	}
}

func writeFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}
