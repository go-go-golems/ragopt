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
	extension := filepath.Ext(absoluteSource)
	relative := filepath.Join("inputs", copiedInputName(name, extension))
	if filepath.Clean(relative) == filepath.Join("inputs", "manifest.json") {
		return InputRef{}, errors.Errorf("input role %q with extension %q collides with reserved input manifest path", role, extension)
	}
	pendingRelative := filepath.Join("inputs", ".pending-"+filepath.Base(relative))

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
	ref := InputRef{
		Role:         role,
		OriginalPath: absoluteSource,
		CopiedPath:   relative,
		SHA256:       digestBytes(data),
		SizeBytes:    int64(len(data)),
	}
	if err := run.writeBytes(ctx, pendingRelative, data); err != nil {
		return InputRef{}, err
	}
	run.inputs = append(run.inputs, ref)
	if err := run.writeJSON(ctx, "inputs/manifest.json", run.inputs); err != nil {
		return InputRef{}, err
	}
	if err := ctx.Err(); err != nil {
		return InputRef{}, err
	}
	pendingPath, err := joinWithin(run.dir, pendingRelative)
	if err != nil {
		return InputRef{}, err
	}
	finalPath, err := joinWithin(run.dir, relative)
	if err != nil {
		return InputRef{}, err
	}
	if err := os.Rename(pendingPath, finalPath); err != nil {
		run.inputs = run.inputs[:len(run.inputs)-1]
		if rollbackErr := run.writeJSON(ctx, "inputs/manifest.json", run.inputs); rollbackErr != nil {
			return InputRef{}, errors.Wrapf(rollbackErr, "roll back copied input manifest after publication failure: %v", err)
		}
		return InputRef{}, errors.Wrap(err, "publish copied input")
	}
	if err := syncDirectory(filepath.Dir(finalPath)); err != nil {
		return InputRef{}, errors.Wrap(err, "sync copied input directory")
	}
	return ref, nil
}
