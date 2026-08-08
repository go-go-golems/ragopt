package review

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArtifactsAreDeterministicAndBlinded(t *testing.T) {
	protocol := Protocol{SchemaVersion: "test/v1", Dimensions: []Dimension{{Name: "quality", Min: 0, Max: 3}}}
	candidate := Candidate{SubjectID: "q1", Variant: "vector", Payload: []byte(`{"answer":"hello"}`), Identity: map[string]string{"evidence": "c1"}}
	queue, keys, err := BuildArtifacts(protocol, []Candidate{candidate})
	require.NoError(t, err)
	again, againKeys, err := BuildArtifacts(protocol, []Candidate{candidate})
	require.NoError(t, err)
	require.Equal(t, queue, again)
	require.Equal(t, keys, againKeys)
	require.NotContains(t, string(queue[0].Payload), "vector")
	require.Equal(t, "vector", keys[0].Variant)
}

func TestLoadAnnotationsValidatesKnownIDsRangesAndDuplicates(t *testing.T) {
	keys := []KeyEntry{{ReviewID: "r1", SubjectID: "q1", Variant: "v1"}}
	dimensions := []Dimension{{Name: "quality", Min: 0, Max: 3}}
	path := filepath.Join(t.TempDir(), "annotations.jsonl")
	require.NoError(t, os.WriteFile(path, []byte("{\"review_id\":\"r1\",\"reviewer\":\"alice\",\"scores\":{\"quality\":3}}\n"), 0o600))
	annotations, err := LoadAnnotations(path, keys, dimensions)
	require.NoError(t, err)
	require.Len(t, annotations, 1)
	require.NoError(t, os.WriteFile(path, []byte("{\"review_id\":\"r1\",\"reviewer\":\"alice\",\"scores\":{\"quality\":4}}\n"), 0o600))
	_, err = LoadAnnotations(path, keys, dimensions)
	require.Error(t, err)
}
