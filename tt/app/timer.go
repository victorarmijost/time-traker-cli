package app

import (
	"context"
	"time"
)

// StartTimer starts a standalone countdown, independent of the working/idle tracking state.
// If a work record is open while it runs, that time is tracked as usual; if not, nothing is recorded.
func (kern *App) StartTimer(ctx context.Context, duration time.Duration) error {
	kern.timer.Start(duration)
	return nil
}

func (kern *App) StopTimer(ctx context.Context) error {
	kern.timer.Stop()
	return nil
}

func (kern *App) GetTimerRemaining() time.Duration {
	return kern.timer.Remaining()
}
