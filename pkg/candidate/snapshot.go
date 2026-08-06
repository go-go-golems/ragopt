package candidate

import (
	"context"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

var (
	logicalNamePattern  = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,127}$`)
	dimensionKeyPattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_.-]{0,127}$`)
	digestPattern       = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
)

// LoadSnapshot strictly loads and validates one snapshot manifest and all of
// its declared asset bytes relative to bundleRoot.
func LoadSnapshot(ctx context.Context, bundleRoot, manifestPath string) (*Snapshot, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	root, err := resolveBundleRoot(bundleRoot)
	if err != nil {
		return nil, err
	}
	return loadSnapshot(ctx, root, manifestPath)
}

func loadSnapshot(ctx context.Context, resolvedRoot, manifestPath string) (*Snapshot, error) {
	path, err := resolveBundleFile(resolvedRoot, manifestPath)
	if err != nil {
		return nil, errors.Wrap(err, "resolve snapshot manifest")
	}
	var snapshot Snapshot
	if err := readStrictYAML(path, &snapshot); err != nil {
		return nil, errors.Wrap(err, "read snapshot manifest")
	}
	if snapshot.APIVersion != SnapshotAPIVersion {
		return nil, errors.Errorf("unsupported snapshot API version %q", snapshot.APIVersion)
	}
	if !logicalNamePattern.MatchString(snapshot.System) {
		return nil, errors.Errorf("invalid snapshot system %q", snapshot.System)
	}
	if !digestPattern.MatchString(snapshot.SnapshotID) {
		return nil, errors.Errorf("invalid snapshot ID %q", snapshot.SnapshotID)
	}
	if err := validateDimensions(snapshot.Dimensions); err != nil {
		return nil, err
	}

	assets := make(map[string][]byte, len(snapshot.LockedAssets)+len(snapshot.MutableAssets))
	if err := loadAssets(ctx, resolvedRoot, "locked", snapshot.LockedAssets, assets); err != nil {
		return nil, err
	}
	if err := loadAssets(ctx, resolvedRoot, "mutable", snapshot.MutableAssets, assets); err != nil {
		return nil, err
	}
	if len(snapshot.MutableAssets) == 0 {
		return nil, errors.New("snapshot requires at least one mutable asset")
	}
	sort.Slice(snapshot.LockedAssets, func(i, j int) bool { return snapshot.LockedAssets[i].Name < snapshot.LockedAssets[j].Name })
	sort.Slice(snapshot.MutableAssets, func(i, j int) bool { return snapshot.MutableAssets[i].Name < snapshot.MutableAssets[j].Name })
	expected, err := DigestSnapshot(snapshot)
	if err != nil {
		return nil, err
	}
	if snapshot.SnapshotID != expected {
		return nil, errors.Errorf("snapshot ID mismatch: manifest=%s actual=%s", snapshot.SnapshotID, expected)
	}
	snapshot.Dimensions = cloneStrings(snapshot.Dimensions)
	snapshot.assetBytes = assets
	return &snapshot, nil
}

func loadAssets(ctx context.Context, root, class string, refs []AssetRef, assets map[string][]byte) error {
	for index, ref := range refs {
		if err := ctx.Err(); err != nil {
			return err
		}
		if !logicalNamePattern.MatchString(ref.Name) {
			return errors.Errorf("%s asset %d has invalid name %q", class, index, ref.Name)
		}
		if _, exists := assets[ref.Name]; exists {
			return errors.Errorf("duplicate asset name %q", ref.Name)
		}
		if strings.TrimSpace(ref.MediaType) == "" || len(ref.MediaType) > 255 {
			return errors.Errorf("asset %q has invalid media type", ref.Name)
		}
		if !digestPattern.MatchString(ref.SHA256) {
			return errors.Errorf("asset %q has invalid SHA-256 %q", ref.Name, ref.SHA256)
		}
		if ref.SizeBytes < 0 {
			return errors.Errorf("asset %q has negative size", ref.Name)
		}
		path, err := resolveBundleFile(root, ref.Path)
		if err != nil {
			return errors.Wrapf(err, "resolve asset %q", ref.Name)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "read asset %q", ref.Name)
		}
		if int64(len(data)) != ref.SizeBytes {
			return errors.Errorf("asset %q size mismatch: manifest=%d actual=%d", ref.Name, ref.SizeBytes, len(data))
		}
		if actual := digestBytes(data); actual != ref.SHA256 {
			return errors.Errorf("asset %q digest mismatch: manifest=%s actual=%s", ref.Name, ref.SHA256, actual)
		}
		assets[ref.Name] = append([]byte(nil), data...)
	}
	return nil
}

func validateDimensions(dimensions map[string]string) error {
	for key, value := range dimensions {
		if !dimensionKeyPattern.MatchString(key) {
			return errors.Errorf("invalid snapshot dimension key %q", key)
		}
		if strings.TrimSpace(value) != value || value == "" || len(value) > 1024 {
			return errors.Errorf("snapshot dimension %q has an invalid value", key)
		}
	}
	return nil
}

func cloneStrings(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}
	result := make(map[string]string, len(values))
	for key, value := range values {
		result[key] = value
	}
	return result
}
