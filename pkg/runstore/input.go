package runstore

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// CopyInput copies an immutable input into inputs/ and records its identity.
func (run *Run) CopyInput(ctx context.Context, role, source string) (InputRef, error) {
	if run == nil {
		return InputRef{}, errors.New("run is nil")
	}
	if err := ctx.Err(); err != nil {
		return InputRef{}, err
	}
	name, err := safeName("input role", role)
	if err != nil {
		return InputRef{}, err
	}
	absoluteSource, err := filepath.Abs(source)
	if err != nil {
		return InputRef{}, errors.Wrap(err, "resolve input path")
	}
	data, err := os.ReadFile(absoluteSource)
	if err != nil {
		return InputRef{}, errors.Wrap(err, "read input")
	}
	relative := filepath.Join("inputs", name+filepath.Ext(absoluteSource))

	run.mu.Lock()
	defer run.mu.Unlock()
	if run.terminal {
		return InputRef{}, errors.New("run is terminal")
	}
	for _, existing := range run.inputs {
		if existing.Role == role {
			return InputRef{}, errors.Errorf("input role %q is already registered", role)
		}
		if existing.CopiedPath == relative {
			return InputRef{}, errors.Errorf("input copied path %q is already registered", relative)
		}
	}
	if err := run.writeBytes(ctx, relative, data); err != nil {
		return InputRef{}, err
	}
	ref := InputRef{
		Role:         role,
		OriginalPath: absoluteSource,
		CopiedPath:   relative,
		SHA256:       digestBytes(data),
		SizeBytes:    int64(len(data)),
	}
	run.inputs = append(run.inputs, ref)
	if err := run.writeJSON(ctx, "inputs/manifest.json", run.inputs); err != nil {
		return InputRef{}, err
	}
	return ref, nil
}
