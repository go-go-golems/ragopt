package review

import (
	"os"
	"path/filepath"
	"strings"
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

func TestBuildArtifactsValidatesProtocolAndOwnsPayload(t *testing.T) {
	valid := Protocol{SchemaVersion: "test/v1", Dimensions: []Dimension{{Name: "quality", Min: 0, Max: 3}}}
	for name, protocol := range map[string]Protocol{
		"empty name":      {SchemaVersion: "test/v1", Dimensions: []Dimension{{Name: "", Min: 0, Max: 3}}},
		"duplicate name":  {SchemaVersion: "test/v1", Dimensions: []Dimension{{Name: "quality"}, {Name: "quality"}}},
		"inverted bounds": {SchemaVersion: "test/v1", Dimensions: []Dimension{{Name: "quality", Min: 4, Max: 3}}},
	} {
		t.Run(name, func(t *testing.T) {
			_, _, err := BuildArtifacts(protocol, []Candidate{{SubjectID: "q1", Variant: "a", Payload: []byte(`{}`)}})
			require.Error(t, err)
		})
	}
	for name, candidate := range map[string]Candidate{
		"subject whitespace": {SubjectID: " q1", Variant: "a", Payload: []byte(`{}`)},
		"variant whitespace": {SubjectID: "q1", Variant: "a ", Payload: []byte(`{}`)},
	} {
		t.Run(name, func(t *testing.T) {
			_, _, err := BuildArtifacts(valid, []Candidate{candidate})
			require.ErrorContains(t, err, "surrounding whitespace")
		})
	}
	payload := []byte(`{"answer":"original"}`)
	queue, _, err := BuildArtifacts(valid, []Candidate{{SubjectID: "q1", Variant: "a", Payload: payload}})
	require.NoError(t, err)
	payload[11] = 'X'
	require.JSONEq(t, `{"answer":"original"}`, string(queue[0].Payload))
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

func TestLoadAnnotationsRejectsTrailingJSONValue(t *testing.T) {
	keys := []KeyEntry{{ReviewID: "r1", SubjectID: "q1", Variant: "v1"}}
	dimensions := []Dimension{{Name: "quality", Min: 0, Max: 3}}
	path := filepath.Join(t.TempDir(), "annotations.jsonl")
	line := `{"review_id":"r1","reviewer":"alice","scores":{"quality":3}} {"ignored":true}` + "\n"
	require.NoError(t, os.WriteFile(path, []byte(line), 0o600))
	_, err := LoadAnnotations(path, keys, dimensions)
	require.ErrorContains(t, err, "trailing JSON value")
}

func TestLoadAnnotationsRejectsReviewerAliasesAndUndeclaredScores(t *testing.T) {
	keys := []KeyEntry{{ReviewID: "r1", SubjectID: "q1", Variant: "v1"}}
	dimensions := []Dimension{{Name: "quality", Min: 0, Max: 3}}
	path := filepath.Join(t.TempDir(), "annotations.jsonl")
	for name, line := range map[string]string{
		"reviewer whitespace": `{"review_id":"r1","reviewer":" alice ","scores":{"quality":3}}`,
		"undeclared score":    `{"review_id":"r1","reviewer":"alice","scores":{"quality":3,"qualty":2}}`,
	} {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.WriteFile(path, []byte(line+"\n"), 0o600))
			_, err := LoadAnnotations(path, keys, dimensions)
			require.Error(t, err)
		})
	}
}

func TestLoadAnnotationsAcceptsLargeNotesAndRejectsInvalidDimensions(t *testing.T) {
	keys := []KeyEntry{{ReviewID: "r1", SubjectID: "q1", Variant: "v1"}}
	dimensions := []Dimension{{Name: "quality", Min: 0, Max: 3}}
	path := filepath.Join(t.TempDir(), "annotations.jsonl")
	notes := strings.Repeat("n", 128*1024)
	line := `{"review_id":"r1","reviewer":"alice","scores":{"quality":3},"notes":"` + notes + `"}` + "\n"
	require.NoError(t, os.WriteFile(path, []byte(line), 0o600))
	annotations, err := LoadAnnotations(path, keys, dimensions)
	require.NoError(t, err)
	require.Equal(t, notes, annotations[0].Notes)
	_, err = LoadAnnotations(path, keys, []Dimension{{Name: "quality"}, {Name: "quality"}})
	require.ErrorContains(t, err, "duplicate")
}

func TestReviewBoundariesRejectUnsafeRangesAndDuplicateKeys(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	if err := ValidateDimensions([]Dimension{{Name: "quality", Min: -1, Max: maxInt}}); err == nil || !strings.Contains(err.Error(), "too wide") {
		t.Fatalf("expected unsafe disagreement range rejection, got %v", err)
	}
	keys := []KeyEntry{
		{ReviewID: "r1", SubjectID: "q1", Variant: "control"},
		{ReviewID: "r1", SubjectID: "q1", Variant: "candidate"},
	}
	_, err := LoadAnnotations("", keys, []Dimension{{Name: "quality", Min: 0, Max: 3}})
	if err == nil || !strings.Contains(err.Error(), "duplicate unblinding key") {
		t.Fatalf("expected duplicate key rejection, got %v", err)
	}
}

func TestValidateKeysRejectsMalformedIdentities(t *testing.T) {
	for name, key := range map[string]KeyEntry{
		"empty review ID":         {SubjectID: "q1", Variant: "control"},
		"spaced review ID":        {ReviewID: " r1", SubjectID: "q1", Variant: "control"},
		"empty subject ID":        {ReviewID: "r1", Variant: "control"},
		"spaced subject ID":       {ReviewID: "r1", SubjectID: "q1 ", Variant: "control"},
		"empty variant":           {ReviewID: "r1", SubjectID: "q1"},
		"whitespace-only variant": {ReviewID: "r1", SubjectID: "q1", Variant: " "},
	} {
		t.Run(name, func(t *testing.T) {
			if err := ValidateKeys([]KeyEntry{key}); err == nil {
				t.Fatalf("ValidateKeys accepted %+v", key)
			}
		})
	}
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

func TestAggregateReviewerDifferenceDoesNotOverflow(t *testing.T) {
	maximum := int(^uint(0) >> 1)
	dimensions := []Dimension{{Name: "quality", Min: 0, Max: maximum}}
	keys := []KeyEntry{
		{ReviewID: "a", SubjectID: "q1", Variant: "control"},
		{ReviewID: "b", SubjectID: "q2", Variant: "control"},
	}
	annotations := []Annotation{
		{ReviewID: "a", Reviewer: "alice", Scores: map[string]int{"quality": maximum}},
		{ReviewID: "a", Reviewer: "bob", Scores: map[string]int{"quality": 0}},
		{ReviewID: "b", Reviewer: "alice", Scores: map[string]int{"quality": maximum}},
		{ReviewID: "b", Reviewer: "bob", Scores: map[string]int{"quality": 0}},
	}
	report := Aggregate(keys, annotations, dimensions)
	require.Len(t, report.ReviewerOverlaps, 1)
	require.Equal(t, float64(maximum), report.ReviewerOverlaps[0].Dimensions["quality"].MeanAbsoluteDifference)
}

func TestAggregatePreservesDuplicateSubjectVariantItems(t *testing.T) {
	dimensions := []Dimension{{Name: "quality", Min: 0, Max: 4}}
	keys := []KeyEntry{
		{ReviewID: "a1", SubjectID: "q1", Variant: "a"},
		{ReviewID: "a2", SubjectID: "q1", Variant: "a"},
		{ReviewID: "b1", SubjectID: "q1", Variant: "b"},
	}
	annotations := []Annotation{
		{ReviewID: "a1", Reviewer: "alice", Scores: map[string]int{"quality": 2}},
		{ReviewID: "a2", Reviewer: "alice", Scores: map[string]int{"quality": 4}},
		{ReviewID: "b1", Reviewer: "alice", Scores: map[string]int{"quality": 1}},
	}
	report := Aggregate(keys, annotations, dimensions)
	require.Len(t, report.Comparisons, 1)
	require.Equal(t, 1, report.Comparisons[0].ReviewedPairs)
	require.Equal(t, 2.0, report.Comparisons[0].MeanDelta, "mean(a items)=3 minus b=1")
}

func TestAggregatePreservesLargeIntegerPairDifference(t *testing.T) {
	maximum := int(^uint(0) >> 1)
	dimensions := []Dimension{{Name: "quality", Min: 0, Max: maximum}}
	keys := []KeyEntry{
		{ReviewID: "a", SubjectID: "q1", Variant: "a"},
		{ReviewID: "b", SubjectID: "q1", Variant: "b"},
	}
	annotations := []Annotation{
		{ReviewID: "a", Reviewer: "alice", Scores: map[string]int{"quality": maximum}},
		{ReviewID: "b", Reviewer: "alice", Scores: map[string]int{"quality": maximum - 1}},
	}
	report := Aggregate(keys, annotations, dimensions)
	require.Len(t, report.Comparisons, 1)
	require.Equal(t, 1, report.Comparisons[0].Wins)
	require.Equal(t, 1.0, report.Comparisons[0].MeanDelta)
}
