//go:build !debug

package debug

import (
	"context"
	"errors"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/hmoragrega/longblock/pkg/pauser"
	"strings"
	"testing"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

var errDummy = errors.New("dummy error")

func Test_NewAppModuleFromOpts(t *testing.T) {
	tests := []struct {
		name                     string
		appOpts                  servertypes.AppOptions
		expectedPauseAllowed     bool
		expectedPauseOnEachBlock bool
	}{
		{
			name: "PauseAllowed and PauseOnEachBlock are true",
			appOpts: &mockAppOptions{
				options: map[string]interface{}{
					FlagPauseAllowed:     true,
					FlagPauseOnEachBlock: true,
				},
			},
			expectedPauseAllowed:     true,
			expectedPauseOnEachBlock: true,
		},
		{
			name: "PauseAllowed is false, PauseOnEachBlock is true",
			appOpts: &mockAppOptions{
				options: map[string]interface{}{
					FlagPauseAllowed:     false,
					FlagPauseOnEachBlock: true,
				},
			},
			expectedPauseAllowed:     false,
			expectedPauseOnEachBlock: true,
		},
		{
			name: "PauseAllowed is true, PauseOnEachBlock is false",
			appOpts: &mockAppOptions{
				options: map[string]interface{}{
					FlagPauseAllowed:     true,
					FlagPauseOnEachBlock: false,
				},
			},
			expectedPauseAllowed:     true,
			expectedPauseOnEachBlock: false,
		},
		{
			name: "PauseAllowed and PauseOnEachBlock are false",
			appOpts: &mockAppOptions{
				options: map[string]interface{}{
					FlagPauseAllowed:     false,
					FlagPauseOnEachBlock: false,
				},
			},
			expectedPauseAllowed:     false,
			expectedPauseOnEachBlock: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appModule := NewAppModuleFromOpts(tt.appOpts)

			if _, ok := interface{}(appModule).(module.AppModuleBasic); !ok {
				t.Errorf("AppModule should implement module.AppModuleBasic")
			}
			if _, ok := interface{}(appModule).(module.HasServices); !ok {
				t.Errorf("AppModule should implement module.HasServices")
			}

			p, ok := appModule.pauseSvc.(*pauser.Pauser)
			if !ok {
				t.Fatalf("pauseSvc should be of type *Pauser")
			}
			if p.PauseAllowed() != tt.expectedPauseAllowed {
				t.Errorf("expected allowed to be %v, got %v", tt.expectedPauseAllowed, p.PauseAllowed())
			}
			if p.PauseOnEachBlock() != tt.expectedPauseOnEachBlock {
				t.Errorf("expected pauseOnEachBlock to be %v, got %v", tt.expectedPauseOnEachBlock, p.PauseOnEachBlock())
			}
		})
	}
}

func Test_RegisterLegacyAminoCodec(t *testing.T) {
	cdc := codec.NewLegacyAmino()

	AppModule{}.RegisterLegacyAminoCodec(cdc)

	var builder strings.Builder
	err := cdc.PrintTypes(&builder)
	if err != nil {
		t.Fatalf("failed to print amino registered types: %v", err)
	}
	expectedTypes := []string{
		"debug/PauseResponse",
		"debug/ResumeResponse",
	}
	for _, expectedType := range expectedTypes {
		if !strings.Contains(builder.String(), expectedType) {
			t.Errorf("expected type %s not found in registered types", expectedType)
		}
	}
}

func TestRegisterInterfaces(t *testing.T) {
	interfaceRegistry := types.NewInterfaceRegistry()

	AppModule{}.RegisterInterfaces(interfaceRegistry)

	expectedInterfaces := []string{
		"/hmoragrega.longblock.debug.v1.PauseResponse",
		"/hmoragrega.longblock.debug.v1.ResumeResponse",
	}

	for _, expectedInterface := range expectedInterfaces {
		if _, err := interfaceRegistry.Resolve(expectedInterface); err != nil {
			t.Errorf("expected interface %s to be registered, but it was not", expectedInterface)
		}
	}
}

func TestAppModule_Name(t *testing.T) {
	appModule := AppModule{}
	expectedName := ModuleName

	if appModule.Name() != expectedName {
		t.Errorf("expected %v, got %v", expectedName, appModule.Name())
	}
}

func TestAppModule_BeginBlock(t *testing.T) {
	tests := []struct {
		name          string
		pauseSvc      PauseService
		expectedError error
	}{
		{
			name: "HoldIfPaused returns no error",
			pauseSvc: &mockPauseService{
				holdIfPausedErr: nil,
			},
			expectedError: nil,
		},
		{
			name: "HoldIfPaused returns error",
			pauseSvc: &mockPauseService{
				holdIfPausedErr: errDummy,
			},
			expectedError: errDummy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appModule := NewAppModule(tt.pauseSvc)
			err := appModule.BeginBlock(context.Background())
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

type mockPauseService struct {
	holdIfPausedErr error
	pauseErr        error
	resumeErr       error
}

func (m *mockPauseService) HoldIfPaused(ctx context.Context) error {
	return m.holdIfPausedErr
}

func (m *mockPauseService) Pause() error {
	return m.pauseErr
}

func (m *mockPauseService) Resume() error {
	return m.resumeErr
}

type mockAppOptions struct {
	options map[string]interface{}
}

func (m *mockAppOptions) Get(key string) interface{} {
	return m.options[key]
}
