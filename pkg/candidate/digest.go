package candidate

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"

	"github.com/pkg/errors"
)

func digestBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// DigestSnapshot computes the semantic identity encoded by snapshot_id.
// Asset order is normalized by logical name; paths remain part of identity.
func DigestSnapshot(snapshot Snapshot) (string, error) {
	locked := append([]AssetRef(nil), snapshot.LockedAssets...)
	mutable := append([]AssetRef(nil), snapshot.MutableAssets...)
	sort.Slice(locked, func(i, j int) bool { return locked[i].Name < locked[j].Name })
	sort.Slice(mutable, func(i, j int) bool { return mutable[i].Name < mutable[j].Name })
	payload := struct {
		APIVersion    string            `json:"api_version"`
		System        string            `json:"system"`
		LockedAssets  []AssetRef        `json:"locked_assets"`
		MutableAssets []AssetRef        `json:"mutable_assets"`
		Dimensions    map[string]string `json:"dimensions"`
	}{
		APIVersion:    snapshot.APIVersion,
		System:        snapshot.System,
		LockedAssets:  locked,
		MutableAssets: mutable,
		Dimensions:    snapshot.Dimensions,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", errors.Wrap(err, "marshal snapshot identity")
	}
	return digestBytes(data), nil
}

func digestCandidate(manifest CandidateManifest, parentID, childID string) (string, error) {
	payload := struct {
		APIVersion  string              `json:"api_version"`
		CandidateID string              `json:"candidate_id"`
		ParentID    string              `json:"parent_snapshot_id"`
		ChildID     string              `json:"candidate_snapshot_id"`
		Proposer    Proposer            `json:"proposer"`
		Mutation    MutationDeclaration `json:"mutation"`
		Evidence    Evidence            `json:"evidence,omitempty"`
	}{
		APIVersion:  manifest.APIVersion,
		CandidateID: manifest.CandidateID,
		ParentID:    parentID,
		ChildID:     childID,
		Proposer:    manifest.Proposer,
		Mutation:    manifest.Mutation,
		Evidence:    manifest.Evidence,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", errors.Wrap(err, "marshal candidate identity")
	}
	return digestBytes(data), nil
}
