package eval

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

var identifierPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)

// LoadSuite strictly loads, normalizes, and identifies one suite JSON file.
func LoadSuite(ctx context.Context, path string) (*SuiteDocument, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrap(err, "resolve suite path")
	}
	data, err := os.ReadFile(absolute)
	if err != nil {
		return nil, errors.Wrap(err, "read suite")
	}
	var suite Suite
	if err := decodeStrictJSON(data, &suite); err != nil {
		return nil, errors.Wrap(err, "decode suite")
	}
	if err := normalizeSuite(&suite); err != nil {
		return nil, err
	}
	digest, err := digestJSON(suite)
	if err != nil {
		return nil, errors.Wrap(err, "digest suite")
	}
	return &SuiteDocument{Suite: suite, Digest: digest, ByteDigest: digestBytes(data), SourcePath: absolute}, nil
}

func normalizeSuite(suite *Suite) error {
	if suite == nil {
		return errors.New("suite is nil")
	}
	if suite.APIVersion != SuiteAPIVersion {
		return errors.Errorf("unsupported suite API version %q", suite.APIVersion)
	}
	if !identifierPattern.MatchString(suite.Name) {
		return errors.Errorf("invalid suite name %q", suite.Name)
	}
	if len(suite.Cases) == 0 {
		return errors.New("suite requires at least one case")
	}
	caseIDs := make(map[string]struct{}, len(suite.Cases))
	for index := range suite.Cases {
		caseValue := &suite.Cases[index]
		if !identifierPattern.MatchString(caseValue.ID) {
			return errors.Errorf("case %d has invalid ID %q", index, caseValue.ID)
		}
		if _, exists := caseIDs[caseValue.ID]; exists {
			return errors.Errorf("duplicate case ID %q", caseValue.ID)
		}
		caseIDs[caseValue.ID] = struct{}{}
		if len(caseValue.Input) == 0 || !json.Valid(caseValue.Input) {
			return errors.Errorf("case %q input is not valid JSON", caseValue.ID)
		}
		canonicalInput, err := canonicalRawJSON(caseValue.Input)
		if err != nil {
			return errors.Wrapf(err, "canonicalize case %q input", caseValue.ID)
		}
		caseValue.Input = canonicalInput
		groups := make(map[string]struct{}, len(caseValue.Groups))
		for _, group := range caseValue.Groups {
			if !identifierPattern.MatchString(group) {
				return errors.Errorf("case %q has invalid group %q", caseValue.ID, group)
			}
			if group == "all" {
				return errors.Errorf("case %q uses reserved group %q", caseValue.ID, group)
			}
			if _, exists := groups[group]; exists {
				return errors.Errorf("case %q has duplicate group %q", caseValue.ID, group)
			}
			groups[group] = struct{}{}
		}
		sort.Strings(caseValue.Groups)
	}
	return nil
}

func canonicalRawJSON(raw []byte) ([]byte, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, errors.Wrap(err, "decode JSON value")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, errors.New("JSON contains multiple values")
		}
		return nil, errors.Wrap(err, "check JSON trailing data")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil, errors.Wrap(err, "marshal canonical JSON value")
	}
	return data, nil
}

func digestJSON(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", errors.Wrap(err, "marshal JSON identity")
	}
	return digestBytes(data), nil
}

func digestBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func decodeStrictJSON(data []byte, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return errors.Wrap(err, "decode JSON")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("JSON contains multiple values")
		}
		return errors.Wrap(err, "check JSON trailing data")
	}
	return nil
}

func cloneCase(caseValue Case) Case {
	return Case{
		ID:     caseValue.ID,
		Groups: append([]string(nil), caseValue.Groups...),
		Input:  append(json.RawMessage(nil), caseValue.Input...),
	}
}

func validateSuiteDocument(document *SuiteDocument) error {
	if document == nil {
		return errors.New("suite document is required")
	}
	copySuite := document.Suite
	copySuite.Cases = append([]Case(nil), document.Suite.Cases...)
	for index := range copySuite.Cases {
		copySuite.Cases[index] = cloneCase(copySuite.Cases[index])
	}
	if err := normalizeSuite(&copySuite); err != nil {
		return err
	}
	digest, err := digestJSON(copySuite)
	if err != nil {
		return err
	}
	if digest != document.Digest {
		return errors.Errorf("suite digest mismatch: document=%s actual=%s", document.Digest, digest)
	}
	if strings.TrimSpace(document.SourcePath) == "" {
		return errors.New("suite source path is required")
	}
	return nil
}
