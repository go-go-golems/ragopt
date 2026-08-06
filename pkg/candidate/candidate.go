package candidate

import (
	"bytes"
	"context"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// LoadCandidate strictly loads one bundle and returns it only after all schema,
// path, digest, and exactly-one-mutation invariants pass.
func LoadCandidate(ctx context.Context, bundleRoot, manifestPath string) (*Candidate, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	root, err := resolveBundleRoot(bundleRoot)
	if err != nil {
		return nil, err
	}
	path, err := resolveBundleFile(root, manifestPath)
	if err != nil {
		return nil, errors.Wrap(err, "resolve candidate manifest")
	}
	var manifest CandidateManifest
	if err := readStrictYAML(path, &manifest); err != nil {
		return nil, errors.Wrap(err, "read candidate manifest")
	}
	if err := validateManifest(manifest); err != nil {
		return nil, err
	}
	parent, err := loadSnapshot(ctx, root, manifest.ParentSnapshot)
	if err != nil {
		return nil, errors.Wrap(err, "load parent snapshot")
	}
	child, err := loadSnapshot(ctx, root, manifest.CandidateSnapshot)
	if err != nil {
		return nil, errors.Wrap(err, "load candidate snapshot")
	}
	mutation, err := validateMutation(manifest, *parent, *child)
	if err != nil {
		return nil, err
	}
	digest, err := digestCandidate(manifest, parent.SnapshotID, child.SnapshotID)
	if err != nil {
		return nil, err
	}
	return &Candidate{
		Manifest: manifest,
		Parent:   *parent,
		Child:    *child,
		Mutation: mutation,
		Digest:   digest,
		Root:     root,
	}, nil
}

func validateManifest(manifest CandidateManifest) error {
	if manifest.APIVersion != CandidateAPIVersion {
		return errors.Errorf("unsupported candidate API version %q", manifest.APIVersion)
	}
	if !logicalNamePattern.MatchString(manifest.CandidateID) {
		return errors.Errorf("invalid candidate ID %q", manifest.CandidateID)
	}
	if manifest.ParentSnapshot == manifest.CandidateSnapshot {
		return errors.New("parent and candidate snapshot paths must differ")
	}
	if strings.TrimSpace(manifest.Proposer.Kind) == "" || len(manifest.Proposer.Kind) > 64 {
		return errors.New("candidate proposer kind is required")
	}
	if strings.TrimSpace(manifest.Proposer.Identity) == "" || len(manifest.Proposer.Identity) > 255 {
		return errors.New("candidate proposer identity is required")
	}
	if !logicalNamePattern.MatchString(manifest.Mutation.Asset) {
		return errors.Errorf("invalid declared mutation asset %q", manifest.Mutation.Asset)
	}
	if strings.TrimSpace(manifest.Mutation.Hypothesis) == "" {
		return errors.New("candidate mutation hypothesis is required")
	}
	if strings.TrimSpace(manifest.Mutation.ExpectedImprovement.Metric) == "" {
		return errors.New("candidate expected improvement metric is required")
	}
	if len(manifest.Mutation.RegressionRisks) == 0 {
		return errors.New("candidate requires at least one regression risk")
	}
	if err := validateNonemptyUnique("expected improvement group", manifest.Mutation.ExpectedImprovement.Groups); err != nil {
		return err
	}
	if err := validateNonemptyUnique("regression risk", manifest.Mutation.RegressionRisks); err != nil {
		return err
	}
	if manifest.Evidence.DiagnosticManifestDigest != "" && !digestPattern.MatchString(manifest.Evidence.DiagnosticManifestDigest) {
		return errors.Errorf("invalid diagnostic manifest digest %q", manifest.Evidence.DiagnosticManifestDigest)
	}
	return validateNonemptyUnique("selected case ID", manifest.Evidence.SelectedCaseIDs)
}

func validateNonemptyUnique(label string, values []string) error {
	seen := make(map[string]struct{}, len(values))
	for index, value := range values {
		if strings.TrimSpace(value) != value || value == "" {
			return errors.Errorf("%s %d is empty or has surrounding whitespace", label, index)
		}
		if _, exists := seen[value]; exists {
			return errors.Errorf("duplicate %s %q", label, value)
		}
		seen[value] = struct{}{}
	}
	return nil
}

func validateMutation(manifest CandidateManifest, parent, child Snapshot) (Mutation, error) {
	if parent.System != child.System {
		return Mutation{}, errors.Errorf("snapshot system changed: parent=%q candidate=%q", parent.System, child.System)
	}
	if !reflect.DeepEqual(parent.Dimensions, child.Dimensions) {
		return Mutation{}, errors.New("candidate changes locked snapshot dimensions")
	}
	if err := requireIdenticalAssets("locked", parent.LockedAssets, child.LockedAssets, parent.assetBytes, child.assetBytes); err != nil {
		return Mutation{}, err
	}

	parentMutable := indexAssets(parent.MutableAssets)
	childMutable := indexAssets(child.MutableAssets)
	if len(parentMutable) != len(childMutable) {
		return Mutation{}, errors.New("candidate changes the mutable asset set")
	}
	changed := make([]Mutation, 0, 2)
	for name, parentRef := range parentMutable {
		childRef, exists := childMutable[name]
		if !exists {
			return Mutation{}, errors.Errorf("candidate removes mutable asset %q", name)
		}
		parentBytes := parent.assetBytes[name]
		childBytes := child.assetBytes[name]
		if bytes.Equal(parentBytes, childBytes) {
			if parentRef != childRef {
				return Mutation{}, errors.Errorf("unchanged mutable asset %q changes metadata or path", name)
			}
			continue
		}
		if parentRef.MediaType != childRef.MediaType {
			return Mutation{}, errors.Errorf("mutable asset %q changes media type", name)
		}
		changed = append(changed, Mutation{
			AssetName:    name,
			Parent:       parentRef,
			Candidate:    childRef,
			ParentDigest: parentRef.SHA256,
			ChildDigest:  childRef.SHA256,
		})
	}
	if len(changed) != 1 {
		return Mutation{}, errors.Errorf("candidate must change exactly one mutable asset, changed=%d", len(changed))
	}
	if manifest.Mutation.Asset != changed[0].AssetName {
		return Mutation{}, errors.Errorf("declared mutation asset %q does not match changed asset %q", manifest.Mutation.Asset, changed[0].AssetName)
	}
	if parent.SnapshotID == child.SnapshotID {
		return Mutation{}, errors.New("changed candidate has the same snapshot ID as its parent")
	}
	return changed[0], nil
}

func requireIdenticalAssets(class string, parent, child []AssetRef, parentBytes, childBytes map[string][]byte) error {
	parentByName := indexAssets(parent)
	childByName := indexAssets(child)
	if len(parentByName) != len(childByName) {
		return errors.Errorf("candidate changes the %s asset set", class)
	}
	for name, parentRef := range parentByName {
		childRef, exists := childByName[name]
		if !exists {
			return errors.Errorf("candidate removes %s asset %q", class, name)
		}
		if parentRef != childRef || !bytes.Equal(parentBytes[name], childBytes[name]) {
			return errors.Errorf("candidate changes %s asset %q", class, name)
		}
	}
	return nil
}

func indexAssets(assets []AssetRef) map[string]AssetRef {
	result := make(map[string]AssetRef, len(assets))
	for _, asset := range assets {
		result[asset.Name] = asset
	}
	return result
}
