package runstore

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

// Complete writes the terminal summary and marks the run complete.
func (run *Run) Complete(ctx context.Context, summary Summary) error {
	return run.finish(ctx, StateComplete, "", summary)
}

// Fail marks the run failed while preserving every existing artifact.
func (run *Run) Fail(ctx context.Context, cause error) error {
	message := ""
	if cause != nil {
		message = cause.Error()
	}
	return run.finish(ctx, StateFailed, message, Summary{})
}

func (run *Run) finish(ctx context.Context, state, message string, summary Summary) error {
	if run == nil {
		return errors.New("run is nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	run.mu.Lock()
	defer run.mu.Unlock()
	if run.terminal {
		return errors.New("run is already terminal")
	}
	if state == StateComplete {
		if err := run.writeJSON(ctx, "results/summary.json", summary); err != nil {
			return err
		}
	}
	finishedAt := time.Now().UTC()
	run.status.State = state
	run.status.FinishedAt = &finishedAt
	run.status.Error = message
	if err := run.writeJSON(ctx, "status.json", run.status); err != nil {
		return err
	}
	run.terminal = true
	return nil
}
