package candidate

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func resolveBundleRoot(root string) (string, error) {
	if strings.TrimSpace(root) == "" {
		return "", errors.New("bundle root is required")
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", errors.Wrap(err, "resolve bundle root")
	}
	resolved, err := filepath.EvalSymlinks(absRoot)
	if err != nil {
		return "", errors.Wrap(err, "resolve bundle root symlinks")
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", errors.Wrap(err, "stat bundle root")
	}
	if !info.IsDir() {
		return "", errors.New("bundle root is not a directory")
	}
	return resolved, nil
}

func resolveBundleFile(root, relative string) (string, error) {
	if relative == "" || filepath.IsAbs(relative) {
		return "", errors.New("bundle path must be non-empty and relative")
	}
	clean := filepath.Clean(relative)
	if clean != relative || clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", errors.Errorf("bundle path %q is not canonical and confined", relative)
	}
	path := filepath.Join(root, clean)
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", errors.Wrapf(err, "resolve bundle path %q", relative)
	}
	relativeToRoot, err := filepath.Rel(root, resolved)
	if err != nil {
		return "", errors.Wrap(err, "compare resolved bundle path")
	}
	if relativeToRoot == ".." || strings.HasPrefix(relativeToRoot, ".."+string(filepath.Separator)) {
		return "", errors.Errorf("bundle path %q resolves outside bundle root", relative)
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", errors.Wrapf(err, "stat bundle path %q", relative)
	}
	if !info.Mode().IsRegular() {
		return "", errors.Errorf("bundle path %q is not a regular file", relative)
	}
	return resolved, nil
}
