package runstore

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// joinWithin validates a relative artifact path and rejects existing symlink
// components. It protects normal single-process use; like other lexical path
// checks, it is not a defense against an attacker racing filesystem changes.
func joinWithin(root, relative string) (string, error) {
	if strings.TrimSpace(root) == "" {
		return "", errors.New("run root is required")
	}
	if relative == "" || filepath.IsAbs(relative) {
		return "", errors.New("artifact path must be non-empty and relative")
	}
	clean := filepath.Clean(relative)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", errors.New("artifact path escapes run root")
	}

	root = filepath.Clean(root)
	path := filepath.Join(root, clean)
	relativeToRoot, err := filepath.Rel(root, path)
	if err != nil {
		return "", errors.Wrap(err, "compare artifact path with run root")
	}
	if relativeToRoot == ".." || strings.HasPrefix(relativeToRoot, ".."+string(filepath.Separator)) {
		return "", errors.New("artifact path escapes run root")
	}

	current := root
	for _, component := range strings.Split(relativeToRoot, string(filepath.Separator)) {
		current = filepath.Join(current, component)
		info, statErr := os.Lstat(current)
		if statErr != nil {
			if os.IsNotExist(statErr) {
				break
			}
			return "", errors.Wrapf(statErr, "inspect artifact path component %q", current)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", errors.Errorf("artifact path contains symlink component %q", current)
		}
	}
	return path, nil
}

func safeName(label, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.Errorf("%s is required", label)
	}
	var builder strings.Builder
	lastDash := false
	for _, char := range strings.ToLower(value) {
		switch {
		case char >= 'a' && char <= 'z', char >= '0' && char <= '9':
			builder.WriteRune(char)
			lastDash = false
		case char == '-', char == '_', char == ' ':
			if builder.Len() > 0 && !lastDash {
				builder.WriteByte('-')
				lastDash = true
			}
		}
	}
	name := strings.Trim(builder.String(), "-")
	if name == "" {
		return "", errors.Errorf("%s has no safe characters", label)
	}
	return name, nil
}
