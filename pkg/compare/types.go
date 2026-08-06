// Package compare builds strict incumbent/candidate pairs and transparent
// aggregates from one validated evaluation artifact run.
package compare

import "github.com/go-go-golems/ragopt/pkg/eval"

const ReportAPIVersion = "ragopt-comparison/v1"

type PairKey struct {
	CaseID      string `json:"case_id"`
	RepeatIndex int    `json:"repeat_index"`
}

type MetricDelta struct {
	Metric    string  `json:"metric"`
	Incumbent float64 `json:"incumbent"`
	Candidate float64 `json:"candidate"`
	Delta     float64 `json:"delta"`
}

type MetricPresence struct {
	Metric           string `json:"metric"`
	IncumbentPresent bool   `json:"incumbent_present"`
	CandidatePresent bool   `json:"candidate_present"`
}

type CostDelta struct {
	ProviderCalls int   `json:"provider_calls"`
	ToolCalls     int   `json:"tool_calls"`
	TotalTokens   int   `json:"total_tokens"`
	DurationNanos int64 `json:"duration_nanos"`
}

type Pair struct {
	Key            PairKey          `json:"key"`
	Groups         []string         `json:"groups,omitempty"`
	Incumbent      eval.Cell        `json:"incumbent"`
	Candidate      eval.Cell        `json:"candidate"`
	Deltas         []MetricDelta    `json:"deltas,omitempty"`
	MetricPresence []MetricPresence `json:"metric_presence,omitempty"`
	Costs          CostDelta        `json:"costs"`
}

type MissingPair struct {
	Key              PairKey  `json:"key"`
	Groups           []string `json:"groups,omitempty"`
	MissingIncumbent bool     `json:"missing_incumbent"`
	MissingCandidate bool     `json:"missing_candidate"`
}

type MetricAggregate struct {
	Group           string  `json:"group"`
	Metric          string  `json:"metric"`
	ExpectedPairs   int     `json:"expected_pairs"`
	CompletePairs   int     `json:"complete_pairs"`
	PairsWithMetric int     `json:"pairs_with_metric"`
	MeanIncumbent   float64 `json:"mean_incumbent"`
	MeanCandidate   float64 `json:"mean_candidate"`
	MeanDelta       float64 `json:"mean_delta"`
	Wins            int     `json:"wins"`
	Ties            int     `json:"ties"`
	Losses          int     `json:"losses"`
}

type GroupAggregate struct {
	Group                  string  `json:"group"`
	ExpectedPairs          int     `json:"expected_pairs"`
	CompletePairs          int     `json:"complete_pairs"`
	IncumbentCompleted     int     `json:"incumbent_completed"`
	CandidateCompleted     int     `json:"candidate_completed"`
	IncumbentContractValid int     `json:"incumbent_contract_valid"`
	CandidateContractValid int     `json:"candidate_contract_valid"`
	IncumbentFailures      int     `json:"incumbent_failures"`
	CandidateFailures      int     `json:"candidate_failures"`
	IncumbentAbstentions   int     `json:"incumbent_abstentions"`
	CandidateAbstentions   int     `json:"candidate_abstentions"`
	CandidateFailureRate   float64 `json:"candidate_failure_rate"`
	MeanProviderCallsDelta float64 `json:"mean_provider_calls_delta"`
	MeanToolCallsDelta     float64 `json:"mean_tool_calls_delta"`
	MeanTotalTokensDelta   float64 `json:"mean_total_tokens_delta"`
	MeanDurationNanosDelta float64 `json:"mean_duration_nanos_delta"`
}

type Report struct {
	APIVersion      string            `json:"api_version"`
	RunID           string            `json:"run_id"`
	SuiteDigest     string            `json:"suite_digest"`
	PolicyDigest    string            `json:"policy_digest"`
	CandidateID     string            `json:"candidate_id"`
	CandidateDigest string            `json:"candidate_digest"`
	ParentSnapshot  string            `json:"parent_snapshot"`
	ChildSnapshot   string            `json:"child_snapshot"`
	IncumbentArm    string            `json:"incumbent_arm"`
	ChallengerArm   string            `json:"challenger_arm"`
	ExpectedPairs   int               `json:"expected_pairs"`
	CompletePairs   int               `json:"complete_pairs"`
	Pairs           []Pair            `json:"pairs"`
	MissingPairs    []MissingPair     `json:"missing_pairs,omitempty"`
	Metrics         []MetricAggregate `json:"metrics"`
	Groups          []GroupAggregate  `json:"groups"`
}
