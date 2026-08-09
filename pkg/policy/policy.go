// Package policy strictly loads product-authored gate policies independently
// of evaluation and decision packages.
package policy

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const PolicyAPIVersion = "ragopt-gate-policy/v1"

type HardGates struct {
	RequireAllCells      bool               `json:"require_all_cells" yaml:"require_all_cells"`
	RequireCompleted     bool               `json:"require_completed" yaml:"require_completed"`
	RequireContractValid bool               `json:"require_contract_valid" yaml:"require_contract_valid"`
	MaxFailureRate       *float64           `json:"max_failure_rate" yaml:"max_failure_rate"`
	MetricFloors         map[string]float64 `json:"metric_floors,omitempty" yaml:"metric_floors,omitempty"`
}

type Target struct {
	Metric                    string   `json:"metric" yaml:"metric"`
	Groups                    []string `json:"groups,omitempty" yaml:"groups,omitempty"`
	MinimumMeanDelta          float64  `json:"minimum_mean_delta" yaml:"minimum_mean_delta"`
	RequirePositiveEachRepeat bool     `json:"require_positive_each_repeat" yaml:"require_positive_each_repeat"`
}

type Regressions struct {
	MaximumCaseDelta map[string]float64            `json:"maximum_case_delta,omitempty" yaml:"maximum_case_delta,omitempty"`
	MaximumMeanDelta map[string]map[string]float64 `json:"maximum_mean_delta,omitempty" yaml:"maximum_mean_delta,omitempty"`
}

type Policy struct {
	APIVersion  string      `json:"api_version" yaml:"api_version"`
	Name        string      `json:"name" yaml:"name"`
	HardGates   HardGates   `json:"hard_gates" yaml:"hard_gates"`
	Target      Target      `json:"target" yaml:"target"`
	Regressions Regressions `json:"regressions" yaml:"regressions"`
	TieBreakers []string    `json:"tie_breakers,omitempty" yaml:"tie_breakers,omitempty"`
}

type Document struct {
	Policy     Policy `json:"policy"`
	Digest     string `json:"digest"`
	ByteDigest string `json:"byte_digest"`
	Path       string `json:"path"`
}

var identifierPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)

// Load strictly loads one policy without supplying metric or threshold
// defaults.
func Load(ctx context.Context, path string) (*Document, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrap(err, "resolve gate policy path")
	}
	data, err := os.ReadFile(absolute)
	if err != nil {
		return nil, errors.Wrap(err, "read gate policy")
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	var policy Policy
	if err := decoder.Decode(&policy); err != nil {
		return nil, errors.Wrap(err, "decode strict gate policy YAML")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, errors.New("gate policy contains multiple YAML documents")
		}
		return nil, errors.Wrap(err, "check gate policy trailing document")
	}
	if err := validateRequiredPolicyFields(data); err != nil {
		return nil, err
	}
	if err := Validate(policy); err != nil {
		return nil, err
	}
	semanticDigest, err := Digest(policy)
	if err != nil {
		return nil, err
	}
	return &Document{
		Policy: policy, Digest: semanticDigest, ByteDigest: digestBytes(data), Path: absolute,
	}, nil
}

func validateRequiredPolicyFields(data []byte) error {
	var required struct {
		Target struct {
			MinimumMeanDelta *float64 `yaml:"minimum_mean_delta"`
		} `yaml:"target"`
	}
	if err := yaml.Unmarshal(data, &required); err != nil {
		return errors.Wrap(err, "decode required gate policy fields")
	}
	if required.Target.MinimumMeanDelta == nil {
		return errors.New("target minimum_mean_delta is required")
	}
	return nil
}

// Digest returns the canonical semantic identity of a validated policy.
func Digest(policy Policy) (string, error) {
	if err := Validate(policy); err != nil {
		return "", err
	}
	semantic, err := json.Marshal(policy)
	if err != nil {
		return "", errors.Wrap(err, "marshal semantic gate policy")
	}
	return digestBytes(semantic), nil
}

// Validate checks the complete policy semantic contract.
func Validate(policy Policy) error {
	if policy.APIVersion != PolicyAPIVersion {
		return errors.Errorf("unsupported gate policy API version %q", policy.APIVersion)
	}
	if !identifierPattern.MatchString(policy.Name) {
		return errors.Errorf("invalid gate policy name %q", policy.Name)
	}
	if policy.HardGates.MaxFailureRate == nil {
		return errors.New("gate policy max_failure_rate is required")
	}
	if value := *policy.HardGates.MaxFailureRate; !finite(value) || value < 0 || value > 1 {
		return errors.New("gate policy max_failure_rate must be finite and between zero and one")
	}
	if err := validateMetricMap("metric floor", policy.HardGates.MetricFloors); err != nil {
		return err
	}
	if !identifierPattern.MatchString(policy.Target.Metric) {
		return errors.Errorf("invalid target metric %q", policy.Target.Metric)
	}
	if !finite(policy.Target.MinimumMeanDelta) {
		return errors.New("target minimum_mean_delta must be finite")
	}
	if err := validateIdentifiers("target group", policy.Target.Groups); err != nil {
		return err
	}
	if err := validateMetricMap("maximum case delta", policy.Regressions.MaximumCaseDelta); err != nil {
		return err
	}
	for group, metrics := range policy.Regressions.MaximumMeanDelta {
		if group != "all" && !identifierPattern.MatchString(group) {
			return errors.Errorf("invalid regression group %q", group)
		}
		if err := validateMetricMap("maximum mean delta for "+group, metrics); err != nil {
			return err
		}
	}
	allowedTieBreakers := map[string]struct{}{
		"provider_calls": {}, "tool_calls": {}, "total_tokens": {}, "duration": {},
	}
	seen := map[string]struct{}{}
	for _, tieBreaker := range policy.TieBreakers {
		if _, ok := allowedTieBreakers[tieBreaker]; !ok {
			return errors.Errorf("unsupported tie breaker %q", tieBreaker)
		}
		if _, exists := seen[tieBreaker]; exists {
			return errors.Errorf("duplicate tie breaker %q", tieBreaker)
		}
		seen[tieBreaker] = struct{}{}
	}
	return nil
}

func validateMetricMap(label string, values map[string]float64) error {
	for metric, value := range values {
		if !identifierPattern.MatchString(metric) {
			return errors.Errorf("%s has invalid metric %q", label, metric)
		}
		if !finite(value) {
			return errors.Errorf("%s %q is not finite", label, metric)
		}
	}
	return nil
}

func validateIdentifiers(label string, values []string) error {
	seen := map[string]struct{}{}
	for _, value := range values {
		if !identifierPattern.MatchString(value) {
			return errors.Errorf("invalid %s %q", label, value)
		}
		if _, exists := seen[value]; exists {
			return errors.Errorf("duplicate %s %q", label, value)
		}
		seen[value] = struct{}{}
	}
	return nil
}

func finite(value float64) bool { return !math.IsNaN(value) && !math.IsInf(value, 0) }

func digestBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}
