package pollers

import (
	"context"
)

type PlayerPoller interface {
	Name() string
	Start(ctx context.Context)
	// PollPlaying polls player every X seconds and computes the current played track
	PollPlaying(ctx context.Context) error
}
