package pauser

import (
	"context"
	"sync/atomic"
)

// Pauser is a utility that allows to pause and resume processing
type Pauser struct {
	// pauseAllowed indicates if the app is allowed to be paused.
	pauseAllowed bool

	// pauseOnEachBlock indicates if the app should pause on each block.
	pauseOnEachBlock bool

	// resume is used to indicate if the BaseApp is currently paused.
	pause atomic.Bool

	// resume is used to signal the BaseApp to continue processing after a pause.
	resume chan struct{}

	// skip is set, the first N blocks won't be paused
	skip int
}

// NewPauser creates a new Pauser instance.
func NewPauser(pauseAllowed, pauseOnEachBlock bool, skip int) *Pauser {
	return &Pauser{
		pauseAllowed:     pauseAllowed,
		pauseOnEachBlock: pauseOnEachBlock,
		skip:             skip,
		resume:           make(chan struct{}),
	}
}

func (p *Pauser) PauseAllowed() bool {
	return p.pauseAllowed
}

func (p *Pauser) PauseOnEachBlock() bool {
	return p.pauseOnEachBlock
}

func (p *Pauser) HoldIfPaused(ctx context.Context) error {
	if !p.allowed() {
		return nil
	}

	if p.skip > 0 {
		p.skip--
		return nil
	}

	if !p.pauseOnEachBlock && !p.pause.Load() {
		return nil
	}

	// reset the pause signal to allow for a new pause
	defer p.pause.Store(false)

	select {
	case <-p.resume:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Pauser) Pause() error {
	if !p.allowed() {
		return ErrPauseNotAllowed
	}

	if p.pauseOnEachBlock {
		return ErrPauseIsAutomatic
	}

	if !p.pause.CompareAndSwap(false, true) {
		return ErrAlreadyPaused
	}

	return nil
}

func (p *Pauser) Resume() error {
	if !p.allowed() {
		return ErrPauseNotAllowed
	}

	select {
	case p.resume <- struct{}{}:
		return nil
	default:
		return ErrAppNotPaused
	}
}

func (p *Pauser) allowed() bool {
	return p.pauseAllowed
}
