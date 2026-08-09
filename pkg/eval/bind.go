package eval

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/go-go-golems/ragopt/pkg/candidate"
	"github.com/go-go-golems/ragopt/pkg/runstore"
)

const (
	roleSuite             = "evaluation-suite"
	rolePolicy            = "gate-policy"
	roleCandidateManifest = "candidate-manifest"
	roleParentSnapshot    = "parent-snapshot"
	roleChildSnapshot     = "candidate-snapshot"
)

func bindInputs(
	ctx context.Context,
	run *runstore.Run,
	prepared *preparedRequest,
) (CandidateView, CandidateView, error) {
	suiteRef, err := ensureBoundInput(ctx, run, roleSuite, prepared.suite.SourcePath, prepared.suite.ByteDigest, 0)
	if err != nil {
		return CandidateView{}, CandidateView{}, errors.Wrap(err, "copy evaluation suite")
	}
	if suiteRef.SHA256 != prepared.suite.ByteDigest {
		return CandidateView{}, CandidateView{}, errors.Errorf("evaluation suite changed during binding: expected=%s copied=%s", prepared.suite.ByteDigest, suiteRef.SHA256)
	}
	policyRef, err := ensureBoundInput(ctx, run, rolePolicy, prepared.policyPath, prepared.config.PolicyDigest, 0)
	if err != nil {
		return CandidateView{}, CandidateView{}, errors.Wrap(err, "copy gate policy")
	}
	if policyRef.SHA256 != prepared.config.PolicyDigest {
		return CandidateView{}, CandidateView{}, errors.Errorf("gate policy changed during binding: expected=%s copied=%s", prepared.config.PolicyDigest, policyRef.SHA256)
	}
	candidateValue := prepared.candidate
	candidateRef, err := ensureBoundInput(ctx, run, roleCandidateManifest, filepath.Join(candidateValue.Root, candidateValue.ManifestPath), candidateValue.ManifestByteDigest, 0)
	if err != nil {
		return CandidateView{}, CandidateView{}, errors.Wrap(err, "copy candidate manifest")
	}
	if candidateRef.SHA256 != candidateValue.ManifestByteDigest {
		return CandidateView{}, CandidateView{}, errors.Errorf("candidate manifest changed during binding: expected=%s copied=%s", candidateValue.ManifestByteDigest, candidateRef.SHA256)
	}
	parentRef, err := ensureBoundInput(ctx, run, roleParentSnapshot, filepath.Join(candidateValue.Root, candidateValue.Manifest.ParentSnapshot), candidateValue.Parent.ByteDigest, 0)
	if err != nil {
		return CandidateView{}, CandidateView{}, errors.Wrap(err, "copy parent snapshot")
	}
	if parentRef.SHA256 != candidateValue.Parent.ByteDigest {
		return CandidateView{}, CandidateView{}, errors.Errorf("parent snapshot changed during binding: expected=%s copied=%s", candidateValue.Parent.ByteDigest, parentRef.SHA256)
	}
	childRef, err := ensureBoundInput(ctx, run, roleChildSnapshot, filepath.Join(candidateValue.Root, candidateValue.Manifest.CandidateSnapshot), candidateValue.Child.ByteDigest, 0)
	if err != nil {
		return CandidateView{}, CandidateView{}, errors.Wrap(err, "copy candidate snapshot")
	}
	if childRef.SHA256 != candidateValue.Child.ByteDigest {
		return CandidateView{}, CandidateView{}, errors.Errorf("candidate snapshot changed during binding: expected=%s copied=%s", candidateValue.Child.ByteDigest, childRef.SHA256)
	}
	if err := copySnapshotAssets(ctx, run, candidateValue.Root, "parent", candidateValue.Parent); err != nil {
		return CandidateView{}, CandidateView{}, err
	}
	if err := copySnapshotAssets(ctx, run, candidateValue.Root, "candidate", candidateValue.Child); err != nil {
		return CandidateView{}, CandidateView{}, err
	}
	if err := verifyBoundInputs(run, prepared.config.InputDigests); err != nil {
		return CandidateView{}, CandidateView{}, err
	}
	return viewsFromInputs(run, candidateValue)
}

func ensureBoundInput(ctx context.Context, run *runstore.Run, role, source, expectedDigest string, expectedSize int64) (runstore.InputRef, error) {
	for _, input := range run.Inputs() {
		if input.Role != role {
			continue
		}
		if input.SHA256 != expectedDigest || (expectedSize > 0 && input.SizeBytes != expectedSize) {
			return runstore.InputRef{}, errors.Errorf("existing bound input %q differs from requested source identity", role)
		}
		path, err := run.Path(input.CopiedPath)
		if err != nil {
			return runstore.InputRef{}, err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return runstore.InputRef{}, errors.Wrapf(err, "read existing bound input %q", role)
		}
		if digestBytes(data) != expectedDigest || int64(len(data)) != input.SizeBytes {
			return runstore.InputRef{}, errors.Errorf("existing bound input %q content changed", role)
		}
		return input, nil
	}
	return run.CopyInput(ctx, role, source)
}

func expectedInputDigests(suite *SuiteDocument, policyDigest string, candidateValue *candidate.Candidate) (map[string]string, error) {
	if suite == nil || candidateValue == nil {
		return nil, errors.New("suite and candidate are required")
	}
	result := map[string]string{
		roleSuite:             suite.ByteDigest,
		rolePolicy:            policyDigest,
		roleCandidateManifest: candidateValue.ManifestByteDigest,
		roleParentSnapshot:    candidateValue.Parent.ByteDigest,
		roleChildSnapshot:     candidateValue.Child.ByteDigest,
	}
	for _, side := range []struct {
		name     string
		snapshot candidate.Snapshot
	}{
		{name: "parent", snapshot: candidateValue.Parent},
		{name: "candidate", snapshot: candidateValue.Child},
	} {
		for _, item := range snapshotAssets(side.snapshot) {
			role := assetRole(side.name, item.ref, item.locked)
			if _, exists := result[role]; exists {
				return nil, errors.Errorf("duplicate expected input role %q", role)
			}
			result[role] = item.ref.SHA256
		}
	}
	for role, digest := range result {
		if strings.TrimSpace(digest) == "" {
			return nil, errors.Errorf("expected input role %q has no captured byte digest", role)
		}
	}
	return result, nil
}

func verifyBoundInputs(run *runstore.Run, expected map[string]string) error {
	actual := make(map[string]runstore.InputRef)
	for _, input := range run.Inputs() {
		if _, exists := actual[input.Role]; exists {
			return errors.Errorf("duplicate bound input role %q", input.Role)
		}
		actual[input.Role] = input
	}
	if len(actual) != len(expected) {
		return errors.Errorf("bound input count mismatch: run=%d expected=%d", len(actual), len(expected))
	}
	for role, digest := range expected {
		input, exists := actual[role]
		if !exists {
			return errors.Errorf("run is missing bound input role %q", role)
		}
		if input.SHA256 != digest {
			return errors.Errorf("bound input %q digest mismatch: run=%s expected=%s", role, input.SHA256, digest)
		}
		path, err := run.Path(input.CopiedPath)
		if err != nil {
			return errors.Wrapf(err, "resolve bound input %q", role)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "read bound input %q", role)
		}
		if actual := digestBytes(data); actual != digest || int64(len(data)) != input.SizeBytes {
			return errors.Errorf("bound input %q content changed: manifest=%s actual=%s", role, digest, actual)
		}
	}
	return nil
}

func copySnapshotAssets(ctx context.Context, run *runstore.Run, root, side string, snapshot candidate.Snapshot) error {
	for _, item := range snapshotAssets(snapshot) {
		role := assetRole(side, item.ref, item.locked)
		copied, err := ensureBoundInput(ctx, run, role, filepath.Join(root, item.ref.Path), item.ref.SHA256, item.ref.SizeBytes)
		if err != nil {
			return errors.Wrapf(err, "copy %s asset %q", side, item.ref.Name)
		}
		if copied.SHA256 != item.ref.SHA256 || copied.SizeBytes != item.ref.SizeBytes {
			return errors.Errorf("%s asset %q changed during binding", side, item.ref.Name)
		}
	}
	return nil
}

func viewsFromInputs(run *runstore.Run, candidateValue *candidate.Candidate) (CandidateView, CandidateView, error) {
	inputs := make(map[string]runstore.InputRef)
	for _, input := range run.Inputs() {
		if _, exists := inputs[input.Role]; exists {
			return CandidateView{}, CandidateView{}, errors.Errorf("duplicate copied input role %q", input.Role)
		}
		inputs[input.Role] = input
	}
	parent, err := viewFromSnapshot(run, inputs, "incumbent", "parent", candidateValue, candidateValue.Parent)
	if err != nil {
		return CandidateView{}, CandidateView{}, err
	}
	child, err := viewFromSnapshot(run, inputs, "candidate", "candidate", candidateValue, candidateValue.Child)
	if err != nil {
		return CandidateView{}, CandidateView{}, err
	}
	return parent, child, nil
}

func viewFromSnapshot(
	run *runstore.Run,
	inputs map[string]runstore.InputRef,
	role string,
	side string,
	candidateValue *candidate.Candidate,
	snapshot candidate.Snapshot,
) (CandidateView, error) {
	assets := make(map[string]ResolvedAsset, len(snapshot.LockedAssets)+len(snapshot.MutableAssets))
	for _, item := range snapshotAssets(snapshot) {
		inputRole := assetRole(side, item.ref, item.locked)
		input, exists := inputs[inputRole]
		if !exists {
			return CandidateView{}, errors.Errorf("run is missing copied input role %q", inputRole)
		}
		if input.SHA256 != item.ref.SHA256 || input.SizeBytes != item.ref.SizeBytes {
			return CandidateView{}, errors.Errorf("copied input %q does not match snapshot asset %q", inputRole, item.ref.Name)
		}
		path, err := run.Path(input.CopiedPath)
		if err != nil {
			return CandidateView{}, errors.Wrapf(err, "resolve copied asset %q", item.ref.Name)
		}
		assets[item.ref.Name] = ResolvedAsset{Ref: item.ref, Path: path, Locked: item.locked}
	}
	return CandidateView{
		Role:            role,
		CandidateID:     candidateValue.Manifest.CandidateID,
		CandidateDigest: candidateValue.Digest,
		SnapshotDigest:  snapshot.SnapshotID,
		System:          snapshot.System,
		Assets:          assets,
		Dimensions:      cloneStrings(snapshot.Dimensions),
	}, nil
}

type snapshotAsset struct {
	ref    candidate.AssetRef
	locked bool
}

func snapshotAssets(snapshot candidate.Snapshot) []snapshotAsset {
	result := make([]snapshotAsset, 0, len(snapshot.LockedAssets)+len(snapshot.MutableAssets))
	for _, ref := range snapshot.LockedAssets {
		result = append(result, snapshotAsset{ref: ref, locked: true})
	}
	for _, ref := range snapshot.MutableAssets {
		result = append(result, snapshotAsset{ref: ref, locked: false})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].locked != result[j].locked {
			return result[i].locked
		}
		return result[i].ref.Name < result[j].ref.Name
	})
	return result
}

func assetRole(side string, ref candidate.AssetRef, locked bool) string {
	class := "mutable"
	if locked {
		class = "locked"
	}
	name := identityComponent("asset", ref.Name)
	digestSuffix := strings.TrimPrefix(ref.SHA256, "sha256:")
	if len(digestSuffix) > 12 {
		digestSuffix = digestSuffix[:12]
	}
	return fmt.Sprintf("%s-%s-%s-%s", side, class, name, digestSuffix)
}

func cloneStrings(values map[string]string) map[string]string {
	result := make(map[string]string, len(values))
	for key, value := range values {
		result[key] = value
	}
	return result
}
