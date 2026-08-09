package runstore

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunLifecycleAndReader(t *testing.T) {
	root := t.TempDir()
	inputPath := filepath.Join(t.TempDir(), "suite.json")
	mustWriteFile(t, inputPath, []byte(`{"cases":["a"]}`))

	dimensions := map[string]string{"suite": "suite-v1", "model": "model-a"}
	run, err := Create(t.Context(), Options{
		Root:        root,
		Name:        "Paired Evaluation",
		Description: "lifecycle fixture",
		Dimensions:  dimensions,
	}, map[string]any{"repeats": 2, "enabled": true})
	mustNoError(t, err)
	dimensions["model"] = "mutated-by-caller"
	if got := run.Manifest().Dimensions["model"]; got != "model-a" {
		t.Fatalf("manifest dimension alias: got %q", got)
	}

	input, err := run.CopyInput(t.Context(), "evaluation suite", inputPath)
	mustNoError(t, err)
	if input.SHA256 == "" || input.SizeBytes == 0 || !filepath.IsAbs(input.OriginalPath) {
		t.Fatalf("incomplete input reference: %#v", input)
	}
	mustNoError(t, run.AppendJSONL(t.Context(), "results/cells.jsonl", map[string]any{"case_id": "a", "repeat": 0}))
	mustNoError(t, run.AppendJSONL(t.Context(), "results/cells.jsonl", map[string]any{"case_id": "a", "repeat": 1}))

	active, err := Open(run.Dir())
	mustNoError(t, err)
	if active.Status().State != StateActive {
		t.Fatalf("state: got %q", active.Status().State)
	}
	if len(active.Inputs()) != 1 {
		t.Fatalf("inputs: got %d", len(active.Inputs()))
	}
	if countJSONLLines(t, filepath.Join(run.Dir(), "results", "cells.jsonl")) != 2 {
		t.Fatal("expected two durable JSONL records")
	}

	mustNoError(t, run.Complete(t.Context(), Summary{
		Message: "fixture complete",
		Metrics: map[string]any{"cells": 2},
	}))
	completed, err := Open(run.Dir())
	mustNoError(t, err)
	if completed.Status().State != StateComplete || completed.Status().FinishedAt == nil {
		t.Fatalf("invalid completed status: %#v", completed.Status())
	}
	if err := run.WriteBytes(t.Context(), "results/late.txt", []byte("late")); err == nil {
		t.Fatal("expected terminal write rejection")
	}
	if err := run.Complete(t.Context(), Summary{}); err == nil {
		t.Fatal("expected duplicate completion rejection")
	}
}

func TestCopyInputRejectsDuplicateRoleAndPath(t *testing.T) {
	run := mustCreateRun(t, map[string]any{"x": 1})
	first := filepath.Join(t.TempDir(), "one.txt")
	second := filepath.Join(t.TempDir(), "two.txt")
	mustWriteFile(t, first, []byte("one"))
	mustWriteFile(t, second, []byte("two"))
	_, err := run.CopyInput(t.Context(), "suite", first)
	mustNoError(t, err)
	_, err = run.CopyInput(t.Context(), "suite", second)
	if err == nil || !strings.Contains(err.Error(), "already registered") {
		t.Fatalf("expected duplicate role error, got %v", err)
	}
	_, err = run.CopyInput(t.Context(), "SUITE", second)
	if err == nil || !strings.Contains(err.Error(), "copied path") {
		t.Fatalf("expected duplicate copied path error, got %v", err)
	}
}

func TestCopyInputRejectsReservedManifestPath(t *testing.T) {
	run := mustCreateRun(t, map[string]any{"x": 1})
	source := filepath.Join(t.TempDir(), "source.json")
	mustWriteFile(t, source, []byte(`{}`))
	_, err := run.CopyInput(t.Context(), "manifest", source)
	if err == nil || !strings.Contains(err.Error(), "reserved input manifest") {
		t.Fatalf("expected reserved manifest collision, got %v", err)
	}
}

func TestCopyInputBoundsCompleteDestinationComponent(t *testing.T) {
	run := mustCreateRun(t, map[string]any{"x": 1})
	extension := "." + strings.Repeat("e", 240)
	source := filepath.Join(t.TempDir(), "s"+extension)
	mustWriteFile(t, source, []byte("bounded"))
	ref, err := run.CopyInput(t.Context(), strings.Repeat("role", 60), source)
	mustNoError(t, err)
	if got := len(filepath.Base(ref.CopiedPath)); got > maximumFileComponentBytes {
		t.Fatalf("copied input component has %d bytes", got)
	}
	data, err := os.ReadFile(filepath.Join(run.Dir(), ref.CopiedPath))
	mustNoError(t, err)
	if string(data) != "bounded" {
		t.Fatalf("copied input = %q", data)
	}
}

func TestCreateBoundsRunDirectoryComponentAndRetainsName(t *testing.T) {
	name := strings.Repeat("descriptive-name-", 40)
	run, err := Create(t.Context(), Options{Root: t.TempDir(), Name: name}, map[string]any{"x": 1})
	mustNoError(t, err)
	if got := len(filepath.Base(run.Dir())); got > maximumFileComponentBytes {
		t.Fatalf("run directory component has %d bytes", got)
	}
	if run.Manifest().Name != name {
		t.Fatal("manifest name was changed")
	}
}

func TestOpenDoesNotRecoverPendingInputTransactions(t *testing.T) {
	t.Run("committed manifest remains pending", func(t *testing.T) {
		run := mustCreateRun(t, map[string]any{})
		data := []byte("suite")
		ref := InputRef{Role: "suite", CopiedPath: "inputs/suite.json", SHA256: digestBytes(data), SizeBytes: int64(len(data))}
		pending := filepath.Join(run.Dir(), "inputs", ".pending-suite.json")
		mustWriteFile(t, pending, data)
		mustWriteJSON(t, filepath.Join(run.Dir(), "inputs", "manifest.json"), []InputRef{ref})
		if _, err := Open(run.Dir()); err == nil {
			t.Fatal("Open accepted a manifest whose publication is still pending")
		}
		if _, err := os.Stat(pending); err != nil {
			t.Fatalf("reader changed writer-owned pending input: %v", err)
		}
	})

	t.Run("uncommitted bytes remain pending", func(t *testing.T) {
		run := mustCreateRun(t, map[string]any{})
		pending := filepath.Join(run.Dir(), "inputs", ".pending-suite.json")
		mustWriteFile(t, pending, []byte("suite"))
		reader, err := Open(run.Dir())
		mustNoError(t, err)
		if len(reader.Inputs()) != 0 {
			t.Fatalf("inputs = %d, want 0", len(reader.Inputs()))
		}
		if _, err := os.Stat(pending); err != nil {
			t.Fatalf("reader changed writer-owned pending input: %v", err)
		}
	})
}

func TestResumeRecoversPendingInputTransactions(t *testing.T) {
	config := map[string]any{"x": 1}
	t.Run("committed manifest completes publication", func(t *testing.T) {
		run := mustCreateRun(t, config)
		data := []byte("suite")
		ref := InputRef{Role: "suite", CopiedPath: "inputs/suite.json", SHA256: digestBytes(data), SizeBytes: int64(len(data))}
		mustWriteFile(t, filepath.Join(run.Dir(), "inputs", ".pending-suite.json"), data)
		mustWriteJSON(t, filepath.Join(run.Dir(), "inputs", "manifest.json"), []InputRef{ref})
		resumed, err := Resume(t.Context(), run.Dir(), config)
		mustNoError(t, err)
		if len(resumed.Inputs()) != 1 {
			t.Fatalf("recovered inputs = %d, want 1", len(resumed.Inputs()))
		}
		if _, err := os.Stat(filepath.Join(run.Dir(), ref.CopiedPath)); err != nil {
			t.Fatalf("committed pending input was not published: %v", err)
		}
	})

	t.Run("uncommitted bytes are discarded", func(t *testing.T) {
		run := mustCreateRun(t, config)
		pending := filepath.Join(run.Dir(), "inputs", ".pending-suite.json")
		mustWriteFile(t, pending, []byte("suite"))
		resumed, err := Resume(t.Context(), run.Dir(), config)
		mustNoError(t, err)
		if len(resumed.Inputs()) != 0 {
			t.Fatalf("recovered inputs = %d, want 0", len(resumed.Inputs()))
		}
		if _, err := os.Stat(pending); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("uncommitted pending input remains: %v", err)
		}
	})
}

func TestCreateRetainsAbsoluteRunDirectory(t *testing.T) {
	base := t.TempDir()
	t.Chdir(base)
	mustNoError(t, os.Mkdir("runs", 0o700))
	run, err := Create(t.Context(), Options{Root: "runs", Name: "absolute"}, map[string]any{"x": 1})
	mustNoError(t, err)
	if !filepath.IsAbs(run.Dir()) {
		t.Fatalf("run directory = %q, want absolute", run.Dir())
	}
	t.Chdir(t.TempDir())
	mustNoError(t, run.WriteBytes(t.Context(), "results/after-chdir.txt", []byte("ok")))
	if _, err := os.Stat(filepath.Join(run.Dir(), "results", "after-chdir.txt")); err != nil {
		t.Fatalf("write after chdir: %v", err)
	}
}

func TestConfigDigestIsStableAcrossMapOrder(t *testing.T) {
	first := mustCreateRun(t, map[string]any{"z": 2, "a": map[string]any{"b": true, "a": 1}})
	second := mustCreateRun(t, map[string]any{"a": map[string]any{"a": 1, "b": true}, "z": 2})
	if first.Manifest().ConfigDigest != second.Manifest().ConfigDigest {
		t.Fatalf("config digests differ: %s != %s", first.Manifest().ConfigDigest, second.Manifest().ConfigDigest)
	}
}

func TestResumeRequiresActiveRunAndExactConfig(t *testing.T) {
	config := map[string]any{"suite": "suite-v1", "repeats": 2}
	run := mustCreateRun(t, config)
	mustNoError(t, run.AppendJSONL(t.Context(), "results/cells.jsonl", map[string]any{"cell": 1}))

	resumed, err := Resume(t.Context(), run.Dir(), map[string]any{"repeats": 2, "suite": "suite-v1"})
	mustNoError(t, err)
	mustNoError(t, resumed.AppendJSONL(t.Context(), "results/cells.jsonl", map[string]any{"cell": 2}))
	if countJSONLLines(t, filepath.Join(run.Dir(), "results", "cells.jsonl")) != 2 {
		t.Fatal("resumed writer did not preserve and append result cells")
	}

	_, err = Resume(t.Context(), run.Dir(), map[string]any{"suite": "different", "repeats": 2})
	if err == nil || !strings.Contains(err.Error(), "resume config digest mismatch") {
		t.Fatalf("expected resume identity error, got %v", err)
	}
	mustNoError(t, resumed.Complete(t.Context(), Summary{}))
	_, err = Resume(t.Context(), run.Dir(), config)
	if err == nil || !strings.Contains(err.Error(), "cannot resume run in state \"complete\"") {
		t.Fatalf("expected terminal resume error, got %v", err)
	}
}

func TestPathsRejectEscapesAndSymlinks(t *testing.T) {
	run := mustCreateRun(t, map[string]any{})
	for _, invalid := range []string{"", ".", "..", "../outside", filepath.Join(string(filepath.Separator), "absolute")} {
		if _, err := run.Path(invalid); err == nil {
			t.Fatalf("expected path %q to fail", invalid)
		}
	}
	outside := t.TempDir()
	link := filepath.Join(run.Dir(), "results", "link")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if err := run.WriteBytes(t.Context(), "results/link/escape.txt", []byte("no")); err == nil {
		t.Fatal("expected symlink path rejection")
	}
	if _, err := os.Stat(filepath.Join(outside, "escape.txt")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("outside artifact unexpectedly exists: %v", err)
	}
}

func TestFailPreservesPriorArtifacts(t *testing.T) {
	run := mustCreateRun(t, map[string]any{"mode": "failure"})
	mustNoError(t, run.AppendJSONL(t.Context(), "results/cells.jsonl", map[string]any{"case_id": "kept"}))
	mustNoError(t, run.Fail(t.Context(), errors.New("arm failed")))
	reader, err := Open(run.Dir())
	mustNoError(t, err)
	if reader.Status().State != StateFailed || reader.Status().Error != "arm failed" {
		t.Fatalf("failed status: %#v", reader.Status())
	}
	if countJSONLLines(t, filepath.Join(run.Dir(), "results", "cells.jsonl")) != 1 {
		t.Fatal("prior result was not preserved")
	}
}

func TestOpenRejectsConfigAndInputDrift(t *testing.T) {
	t.Run("config", func(t *testing.T) {
		run := mustCreateRun(t, map[string]any{"locked": true})
		mustWriteFile(t, filepath.Join(run.Dir(), "config.json"), []byte(`{"locked":false}`))
		if _, err := Open(run.Dir()); err == nil || !strings.Contains(err.Error(), "config digest mismatch") {
			t.Fatalf("expected config digest error, got %v", err)
		}
	})
	t.Run("input", func(t *testing.T) {
		run := mustCreateRun(t, map[string]any{})
		source := filepath.Join(t.TempDir(), "suite.txt")
		mustWriteFile(t, source, []byte("original"))
		input, err := run.CopyInput(t.Context(), "suite", source)
		mustNoError(t, err)
		mustWriteFile(t, filepath.Join(run.Dir(), input.CopiedPath), []byte("changed"))
		if _, err := Open(run.Dir()); err == nil || !strings.Contains(err.Error(), "input \"suite\"") {
			t.Fatalf("expected copied input integrity error, got %v", err)
		}
	})
}

func TestOpenRejectsInvalidDimensionsAndUnknownFields(t *testing.T) {
	t.Run("dimension", func(t *testing.T) {
		run := mustCreateRun(t, map[string]any{})
		path := filepath.Join(run.Dir(), "manifest.json")
		var manifest map[string]any
		mustReadJSON(t, path, &manifest)
		manifest["dimensions"] = map[string]string{"suite": " "}
		mustWriteJSON(t, path, manifest)
		if _, err := Open(run.Dir()); err == nil || !strings.Contains(err.Error(), "empty value") {
			t.Fatalf("expected dimension error, got %v", err)
		}
	})
	t.Run("unknown status field", func(t *testing.T) {
		run := mustCreateRun(t, map[string]any{})
		path := filepath.Join(run.Dir(), "status.json")
		var status map[string]any
		mustReadJSON(t, path, &status)
		status["mystery"] = true
		mustWriteJSON(t, path, status)
		if _, err := Open(run.Dir()); err == nil || !strings.Contains(err.Error(), "unknown field") {
			t.Fatalf("expected strict schema error, got %v", err)
		}
	})
}

func TestCreateValidatesAndClonesDimensions(t *testing.T) {
	for name, dimensions := range map[string]map[string]string{
		"blank key":   {"": "value"},
		"spaced key":  {" suite ": "value"},
		"blank value": {"suite": " "},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := Create(t.Context(), Options{Root: t.TempDir(), Name: "fixture", Dimensions: dimensions}, map[string]any{})
			if err == nil {
				t.Fatal("expected invalid dimension error")
			}
		})
	}
}

func mustCreateRun(t *testing.T, config any) *Run {
	t.Helper()
	run, err := Create(context.Background(), Options{Root: t.TempDir(), Name: "fixture"}, config)
	mustNoError(t, err)
	return run
}

func mustNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func mustWriteFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func mustReadJSON(t *testing.T, path string, destination any) {
	t.Helper()
	data, err := os.ReadFile(path)
	mustNoError(t, err)
	mustNoError(t, json.Unmarshal(data, destination))
}

func mustWriteJSON(t *testing.T, path string, value any) {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	mustNoError(t, err)
	mustWriteFile(t, path, append(data, '\n'))
}

func countJSONLLines(t *testing.T, path string) int {
	t.Helper()
	file, err := os.Open(path)
	mustNoError(t, err)
	defer func() { mustNoError(t, file.Close()) }()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
		var value map[string]any
		mustNoError(t, json.Unmarshal(scanner.Bytes(), &value))
	}
	mustNoError(t, scanner.Err())
	return count
}
