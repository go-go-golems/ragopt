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
	if err := ValidateProtocol(protocol); err != nil {
		return nil, nil, err
	}
	queue := make([]QueueEntry, 0, len(candidates))
	keys := make([]KeyEntry, 0, len(candidates))
	seen := map[string]struct{}{}
	for _, candidate := range candidates {
		if strings.TrimSpace(candidate.SubjectID) == "" || strings.TrimSpace(candidate.Variant) == "" {
			return nil, nil, errors.New("review subject ID and variant are required")
		}
		if strings.TrimSpace(candidate.SubjectID) != candidate.SubjectID || strings.TrimSpace(candidate.Variant) != candidate.Variant {
			return nil, nil, errors.New("review subject ID and variant must not have surrounding whitespace")
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
		queue = append(queue, QueueEntry{SchemaVersion: protocol.SchemaVersion, ReviewID: id, SubjectID: candidate.SubjectID, Payload: append(json.RawMessage(nil), candidate.Payload...)})
		keys = append(keys, KeyEntry{ReviewID: id, SubjectID: candidate.SubjectID, Variant: candidate.Variant})
	}
	sort.Slice(queue, func(i, j int) bool { return queue[i].ReviewID < queue[j].ReviewID })
	sort.Slice(keys, func(i, j int) bool { return keys[i].ReviewID < keys[j].ReviewID })
	return queue, keys, nil
}

// ValidateProtocol rejects ambiguous review identities and unusable scoring
// dimensions before either queues or annotations cross a trust boundary.
func ValidateProtocol(protocol Protocol) error {
	if strings.TrimSpace(protocol.SchemaVersion) == "" {
		return errors.New("review schema version is required")
	}
	if strings.TrimSpace(protocol.SchemaVersion) != protocol.SchemaVersion {
		return errors.New("review schema version must not have surrounding whitespace")
	}
	return ValidateDimensions(protocol.Dimensions)
}

// ValidateDimensions validates the shared scoring schema used by queue
// construction and annotation loading.
func ValidateDimensions(dimensions []Dimension) error {
	if len(dimensions) == 0 {
		return errors.New("at least one review dimension is required")
	}
	seen := make(map[string]struct{}, len(dimensions))
	for _, dimension := range dimensions {
		if strings.TrimSpace(dimension.Name) == "" {
			return errors.New("review dimension name is required")
		}
		if strings.TrimSpace(dimension.Name) != dimension.Name {
			return errors.Errorf("review dimension %q must not have surrounding whitespace", dimension.Name)
		}
		if _, exists := seen[dimension.Name]; exists {
			return errors.Errorf("duplicate review dimension %q", dimension.Name)
		}
		seen[dimension.Name] = struct{}{}
		if dimension.Min > dimension.Max {
			return errors.Errorf("review dimension %q has minimum %d greater than maximum %d", dimension.Name, dimension.Min, dimension.Max)
		}
	}
	return nil
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
	if err := ValidateDimensions(dimensions); err != nil {
		return err
	}
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
	if err := ValidateDimensions(dimensions); err != nil {
		return nil, err
	}
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
	reader := bufio.NewReader(file)
	line := 0
	for {
		text, readErr := reader.ReadString('\n')
		if len(text) == 0 && readErr == io.EOF {
			break
		}
		line++
		text = strings.TrimSuffix(text, "\n")
		text = strings.TrimSuffix(text, "\r")
		decoder := json.NewDecoder(strings.NewReader(text))
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
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return nil, errors.Wrap(readErr, "read annotations")
		}
	}
	sort.Slice(annotations, func(i, j int) bool {
		if annotations[i].ReviewID != annotations[j].ReviewID {
			return annotations[i].ReviewID < annotations[j].ReviewID
		}
		return annotations[i].Reviewer < annotations[j].Reviewer
	})
	return annotations, nil
}
