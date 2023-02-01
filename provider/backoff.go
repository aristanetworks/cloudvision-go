// Copyright (c) 2022 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package provider

import "time"

// BackoffTimer implements a backoff mechanism to be used in retries.
type BackoffTimer struct {
	BackoffBase           time.Duration
	BackoffMax            time.Duration
	BackoffMult           float64
	BackoffResetThreshold time.Duration

	backoffCurrent time.Duration
	timer          *time.Timer
	lastBackoff    time.Time
}

// BackoffOption defines options to be used with BackoffTimer
type BackoffOption func(b *BackoffTimer)

// WithBackoffBase configures the base backoff
func WithBackoffBase(base time.Duration) BackoffOption {
	return func(b *BackoffTimer) {
		b.BackoffBase = base
	}
}

// WithBackoffMax configures the maximum backoff
func WithBackoffMax(base time.Duration) BackoffOption {
	return func(b *BackoffTimer) {
		b.BackoffMax = base
	}
}

// NewBackoffTimer creates a new BackoffTimer with defaults settings
func NewBackoffTimer(opts ...BackoffOption) *BackoffTimer {
	b := &BackoffTimer{
		BackoffBase:           5 * time.Second,
		BackoffMax:            3 * time.Minute,
		BackoffMult:           1.8,
		BackoffResetThreshold: time.Hour,
		timer:                 time.NewTimer(time.Nanosecond),
		lastBackoff:           time.Now(),
	}
	for _, o := range opts {
		o(b)
	}
	b.backoffCurrent = b.BackoffBase
	return b
}

// Wait returns a channel that must be waited on, which will signify the backoff
// period has passed.
func (b *BackoffTimer) Wait() <-chan time.Time {
	return b.timer.C
}

// Reset resets the timer with BackoffBase
func (b *BackoffTimer) Reset() {
	if !b.timer.Stop() {
		<-b.timer.C
	}
	b.timer.Reset(b.BackoffBase)
	b.backoffCurrent = b.BackoffBase
}

// Backoff will backoff, so retries are spaced out.
// Returns the current backoff delay, that is set to trigger the next timer.
func (b *BackoffTimer) Backoff() time.Duration {
	if time.Since(b.lastBackoff) > b.BackoffResetThreshold {
		// reset backoff if running fine for a while
		b.backoffCurrent = b.BackoffBase
	}
	waitTime := b.backoffCurrent
	b.timer.Reset(waitTime)
	b.backoffCurrent = time.Duration(float64(b.backoffCurrent) * b.BackoffMult).Round(time.Second)
	if b.backoffCurrent > b.BackoffMax {
		b.backoffCurrent = b.BackoffMax
	}

	b.lastBackoff = time.Now()
	return waitTime
}
