package pauser

import (
	"errors"
)

var (
	ErrDebugIsNotEnabled = errors.New("debug is not enabled")
	ErrPauseNotAllowed   = errors.New("pause is not allowed")
	ErrAlreadyPaused     = errors.New("app is already paused")
	ErrPauseIsAutomatic  = errors.New("app is set to pause automatically on each block")
	ErrAppNotPaused      = errors.New("app is not paused")
)
