package domain

import (
	"sync"
	"time"
)

// Timer is a standalone countdown, independent of the working/idle tracking state.
type Timer struct {
	mu  sync.RWMutex
	end *time.Time
}

func NewTimer() *Timer {
	return &Timer{}
}

func (t *Timer) Start(duration time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	end := time.Now().Add(duration)
	t.end = &end
}

func (t *Timer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.end = nil
}

func (t *Timer) Remaining() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.end == nil {
		return 0
	}

	left := time.Until(*t.end)
	if left < 0 {
		return 0
	}

	return left
}
