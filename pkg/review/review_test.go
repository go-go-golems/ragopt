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

func TestAggregatePairsVariantsAndReviewerOverlap(t *testing.T) {
	dimensions := []Dimension{{Name: "quality", Min: 0, Max: 3}}
	keys := []KeyEntry{{ReviewID: "a", SubjectID: "q1", Variant: "control"}, {ReviewID: "b", SubjectID: "q1", Variant: "candidate"}}
	annotations := []Annotation{{ReviewID: "a", Reviewer: "alice", Scores: map[string]int{"quality": 3}}, {ReviewID: "a", Reviewer: "bob", Scores: map[string]int{"quality": 2}}, {ReviewID: "b", Reviewer: "alice", Scores: map[string]int{"quality": 1}}}
	report := Aggregate(keys, annotations, dimensions)
	require.Len(t, report.Variants, 2)
	require.Len(t, report.Comparisons, 1)
	require.Equal(t, 1, report.Comparisons[0].ReviewedPairs)
	require.Equal(t, 1, report.OverlappingItems)
	require.Len(t, report.ReviewerOverlaps, 1)
	require.Equal(t, 1, report.ReviewerOverlaps[0].Dimensions["quality"].Count)
}
