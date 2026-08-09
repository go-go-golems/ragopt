// Package report renders deterministic, human-reviewable promotion evidence
// and a machine-readable plan that cannot apply a product change.
package report

import (
	"github.com/go-go-golems/ragopt/pkg/candidate"
	"github.com/go-go-golems/ragopt/pkg/gate"
)

const PromotionPlanAPIVersion = "ragopt-promotion-plan/v1"

// PromotionPlan describes the evaluated change and decision. The fixed
// review_required state makes human application an explicit external step.
type PromotionPlan struct {
	APIVersion         string                        `json:"api_version" yaml:"api_version"`
	State              string                        `json:"state" yaml:"state"`
	HumanApplyRequired bool                          `json:"human_apply_required" yaml:"human_apply_required"`
	RunID              string                        `json:"run_id" yaml:"run_id"`
	CandidateID        string                        `json:"candidate_id" yaml:"candidate_id"`
	CandidateDigest    string                        `json:"candidate_digest" yaml:"candidate_digest"`
	ParentSnapshot     string                        `json:"parent_snapshot" yaml:"parent_snapshot"`
	ChildSnapshot      string                        `json:"child_snapshot" yaml:"child_snapshot"`
	Mutation           candidate.MutationDeclaration `json:"mutation" yaml:"mutation"`
	ChangedAsset       string                        `json:"changed_asset" yaml:"changed_asset"`
	ParentAssetDigest  string                        `json:"parent_asset_digest" yaml:"parent_asset_digest"`
	ChildAssetDigest   string                        `json:"child_asset_digest" yaml:"child_asset_digest"`
	PolicyName         string                        `json:"policy_name" yaml:"policy_name"`
	PolicyDigest       string                        `json:"policy_digest" yaml:"policy_digest"`
	Decision           gate.DecisionStatus           `json:"decision" yaml:"decision"`
	Reasons            []string                      `json:"reasons,omitempty" yaml:"reasons,omitempty"`
}

// Document is the two-format promotion-review deliverable.
type Document struct {
	Markdown string
	Plan     PromotionPlan
}
