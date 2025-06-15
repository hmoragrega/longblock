package pauser

import (
	"context"
)

type NoOpPauser struct{}

func NewNoOpPauser() *NoOpPauser {
	return &NoOpPauser{}
}

func (p *NoOpPauser) HoldIfPaused(_ context.Context) error {
	return nil
}

func (p *NoOpPauser) Pause() error {
	return ErrDebugIsNotEnabled
}

func (p *NoOpPauser) Resume() error {
	return ErrDebugIsNotEnabled
}
