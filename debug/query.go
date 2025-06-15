package debug

import (
	"context"
	"fmt"

	"github.com/hmoragrega/longblock/debug/types"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ types.QueryServer = (*QueryServer)(nil)

type QueryServer struct {
	pausable Pausable
}

// NewQueryServerImpl returns an implementation of the debug QueryServer interface
func NewQueryServerImpl(pausable Pausable) *QueryServer {
	return &QueryServer{
		pausable: pausable,
	}
}

func (s QueryServer) Resume(_ context.Context, _ *emptypb.Empty) (res *types.ResumeResponse, err error) {
	res = new(types.ResumeResponse)
	err = s.pausable.Resume()
	if err != nil {
		res.Msg = fmt.Sprintf("Resume failed: %s", err.Error())
		return res, err
	}

	res.Success = true
	res.Msg = "Node resumed successfully"

	return res, nil
}

func (s QueryServer) Pause(_ context.Context, _ *emptypb.Empty) (res *types.PauseResponse, err error) {
	res = new(types.PauseResponse)
	err = s.pausable.Pause()
	if err != nil {
		res.Msg = fmt.Sprintf("Pause failed: %s", err.Error())
		return res, err
	}

	res.Success = true
	res.Msg = "Node paused successfully"

	return res, nil
}
