package pauser

import (
	"context"
	"errors"
	"testing"
)

func TestPauser(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *Pauser
		method        func(p *Pauser) error
		expectedError error
	}{
		{
			name: "HoldIfPaused - not allowed",
			setup: func() *Pauser {
				return NewPauser(false, false, 0)
			},
			method: func(p *Pauser) error {
				return p.HoldIfPaused(context.Background())
			},
			expectedError: nil,
		},
		{
			name: "HoldIfPaused - skip first",
			setup: func() *Pauser {
				return NewPauser(true, false, 1)
			},
			method: func(p *Pauser) error {
				return p.HoldIfPaused(context.Background())
			},
			expectedError: nil,
		},
		{
			name: "HoldIfPaused - pause on each block",
			setup: func() *Pauser {
				return NewPauser(true, true, 0)
			},
			method: func(p *Pauser) error {
				go func() {
					p.resume <- struct{}{} // simulate resuming the app
				}()
				return p.HoldIfPaused(context.Background())
			},
			expectedError: nil,
		},
		{
			name: "HoldIfPaused - not paused",
			setup: func() *Pauser {
				return NewPauser(true, false, 0)
			},
			method: func(p *Pauser) error {
				return p.HoldIfPaused(context.Background())
			},
			expectedError: nil,
		},
		{
			name: "HoldIfPaused - user paused",
			setup: func() *Pauser {
				return NewPauser(true, false, 0)
			},
			method: func(p *Pauser) error {
				go func() {
					p.resume <- struct{}{} // simulate resuming the app
				}()
				p.pause.Store(true)
				return p.HoldIfPaused(context.Background())
			},
			expectedError: nil,
		},
		{
			name: "HoldIfPaused - context canceled",
			setup: func() *Pauser {
				return NewPauser(true, true, 0)
			},
			method: func(p *Pauser) error {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return p.HoldIfPaused(ctx)
			},
			expectedError: context.Canceled,
		},
		{
			name: "Pause - not allowed",
			setup: func() *Pauser {
				return NewPauser(false, false, 0)
			},
			method: func(p *Pauser) error {
				return p.Pause()
			},
			expectedError: ErrPauseNotAllowed,
		},
		{
			name: "Pause - already paused",
			setup: func() *Pauser {
				p := NewPauser(true, false, 0)
				p.pause.Store(true)
				return p
			},
			method: func(p *Pauser) error {
				return p.Pause()
			},
			expectedError: ErrAlreadyPaused,
		},
		{
			name: "Pause - pause is automatic",
			setup: func() *Pauser {
				return NewPauser(true, true, 0)
			},
			method: func(p *Pauser) error {
				return p.Pause()
			},
			expectedError: ErrPauseIsAutomatic,
		},
		{
			name: "Resume - not allowed",
			setup: func() *Pauser {
				return NewPauser(false, false, 0)
			},
			method: func(p *Pauser) error {
				return p.Resume()
			},
			expectedError: ErrPauseNotAllowed,
		},
		{
			name: "Resume - app not paused",
			setup: func() *Pauser {
				return NewPauser(true, false, 0)
			},
			method: func(p *Pauser) error {
				return p.Resume()
			},
			expectedError: ErrAppNotPaused,
		},
		{
			name: "Resume - success",
			setup: func() *Pauser {
				p := NewPauser(true, false, 0)
				p.resume = make(chan struct{}, 1)
				_ = p.Pause()
				return p
			},
			method: func(p *Pauser) error {
				return p.Resume()
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.setup()
			err := tt.method(p)
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}
