// Package review implements a structurally blinded human-review protocol.
// Product code supplies opaque payload and identity values; this package never
// interprets either one and keeps the variant only in the separate key file.
package review

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type Dimension struct {
	Name     string
	Min, Max int
}
type Protocol struct {
	SchemaVersion string
	Dimensions    []Dimension
}
type Candidate struct {
	SubjectID, Variant string
	Payload            json.RawMessage
	Identity           any
}
type QueueEntry struct {
	SchemaVersion string          `json:"schema_version"`
	ReviewID      string          `json:"review_id"`
	SubjectID     string          `json:"subject_id"`
	Payload       json.RawMessage `json:"payload"`
}
type KeyEntry struct {
	ReviewID  string `json:"review_id"`
	SubjectID string `json:"subject_id"`
	Variant   string `json:"variant"`
}
type Annotation struct {
	ReviewID string         `json:"review_id"`
	Reviewer string         `json:"reviewer"`
	Scores   map[string]int `json:"scores"`
	Notes    string         `json:"notes,omitempty"`
}

// BuildArtifacts creates a blind queue and a separately stored unblinding key.
// QueueEntry cannot express Candidate.Variant, which prevents accidental leaks
// through a typed reviewer transport.
func BuildArtifacts(protocol Protocol, candidates []Candidate) ([]QueueEntry, []KeyEntry, error) {
	if strings.TrimSpace(protocol.SchemaVersion) == "" {
		return nil, nil, errors.New("review schema version is required")
	}
	queue := make([]QueueEntry, 0, len(candidates))
	keys := make([]KeyEntry, 0, len(candidates))
	seen := map[string]struct{}{}
	for _, candidate := range candidates {
		if strings.TrimSpace(candidate.SubjectID) == "" || strings.TrimSpace(candidate.Variant) == "" {
			return nil, nil, errors.New("review subject ID and variant are required")
		}
		if !json.Valid(candidate.Payload) {
			return nil, nil, errors.Errorf("review payload for subject %q is not valid JSON", candidate.SubjectID)
		}
		id, err := deterministicID(protocol.SchemaVersion, candidate)
		if err != nil {
			return nil, nil, err
		}
		if _, ok := seen[id]; ok {
			return nil, nil, errors.Errorf("duplicate deterministic review ID %q", id)
		}
		seen[id] = struct{}{}
		queue = append(queue, QueueEntry{SchemaVersion: protocol.SchemaVersion, ReviewID: id, SubjectID: candidate.SubjectID, Payload: candidate.Payload})
		keys = append(keys, KeyEntry{ReviewID: id, SubjectID: candidate.SubjectID, Variant: candidate.Variant})
	}
	sort.Slice(queue, func(i, j int) bool { return queue[i].ReviewID < queue[j].ReviewID })
	sort.Slice(keys, func(i, j int) bool { return keys[i].ReviewID < keys[j].ReviewID })
	return queue, keys, nil
}

func deterministicID(version string, candidate Candidate) (string, error) {
	payload, err := json.Marshal(struct {
		Version, SubjectID, Variant string
		Identity                    any
	}{version, candidate.SubjectID, candidate.Variant, candidate.Identity})
	if err != nil {
		return "", errors.Wrap(err, "marshal review identity")
	}
	digest := sha256.Sum256(payload)
	return "review-" + hex.EncodeToString(digest[:12]), nil
}

func KnownIDs(keys []KeyEntry) map[string]struct{} {
	known := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		known[key.ReviewID] = struct{}{}
	}
	return known
}
func ValidateAnnotation(annotation Annotation, known map[string]struct{}, dimensions []Dimension) error {
	if _, ok := known[annotation.ReviewID]; !ok {
		return errors.Errorf("unknown review ID %q", annotation.ReviewID)
	}
	trimmedReviewer := strings.TrimSpace(annotation.Reviewer)
	if trimmedReviewer == "" {
		return errors.New("reviewer is required")
	}
	if trimmedReviewer != annotation.Reviewer {
		return errors.New("reviewer must not have surrounding whitespace")
	}
	declared := make(map[string]struct{}, len(dimensions))
	for _, dimension := range dimensions {
		declared[dimension.Name] = struct{}{}
		score, ok := annotation.Scores[dimension.Name]
		if !ok {
			return errors.Errorf("%s score is required", dimension.Name)
		}
		if score < dimension.Min || score > dimension.Max {
			return errors.Errorf("%s score %d is outside [%d,%d]", dimension.Name, score, dimension.Min, dimension.Max)
		}
	}
	for name := range annotation.Scores {
		if _, ok := declared[name]; !ok {
			return errors.Errorf("score dimension %q is not declared", name)
		}
	}
	return nil
}

func LoadAnnotations(path string, keys []KeyEntry, dimensions []Dimension) ([]Annotation, error) {
	if strings.TrimSpace(path) == "" {
		return nil, nil
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open annotations")
	}
	defer func() { _ = file.Close() }()
	known, seen := KnownIDs(keys), map[string]struct{}{}
	var annotations []Annotation
	scanner := bufio.NewScanner(file)
	line := 0
	for scanner.Scan() {
		line++
		decoder := json.NewDecoder(strings.NewReader(scanner.Text()))
		decoder.DisallowUnknownFields()
		var annotation Annotation
		if err := decoder.Decode(&annotation); err != nil {
			return nil, errors.Wrapf(err, "decode annotation line %d", line)
		}
		var trailing any
		if err := decoder.Decode(&trailing); err != io.EOF {
			if err == nil {
				return nil, errors.Errorf("decode annotation line %d: trailing JSON value", line)
			}
			return nil, errors.Wrapf(err, "check annotation line %d trailing data", line)
		}
		if err := ValidateAnnotation(annotation, known, dimensions); err != nil {
			return nil, errors.Wrapf(err, "validate annotation line %d", line)
		}
		identity := annotation.ReviewID + "\x00" + annotation.Reviewer
		if _, ok := seen[identity]; ok {
			return nil, errors.Errorf("duplicate annotation for review %q by reviewer %q", annotation.ReviewID, annotation.Reviewer)
		}
		seen[identity] = struct{}{}
		annotations = append(annotations, annotation)
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "scan annotations")
	}
	sort.Slice(annotations, func(i, j int) bool {
		if annotations[i].ReviewID != annotations[j].ReviewID {
			return annotations[i].ReviewID < annotations[j].ReviewID
		}
		return annotations[i].Reviewer < annotations[j].Reviewer
	})
	return annotations, nil
}
