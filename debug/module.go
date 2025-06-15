package debug

import (
	"context"

	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"

	"github.com/hmoragrega/longblock/debug/types"
	"github.com/hmoragrega/longblock/pkg/pauser"
)

const ModuleName = "debug"

const (
	FlagPauseOnEachBlock = "debug.pause-on-each-block"
	FlagPauseAllowed     = "debug.pause-allowed"
	FlagPauseSkip        = "debug.pause-skip"
)

type Cli interface {
	Flags() *pflag.FlagSet
}

// AddModuleInitFlags adds the flags for the debug module to the start command.
func AddModuleInitFlags(cli Cli) {
	cli.Flags().Bool(FlagPauseAllowed, false, "If true, the node will allow pausing and resuming the app (default: false)")
	cli.Flags().Bool(FlagPauseOnEachBlock, false, "If true, the node will pause on each block and wait for user input before continuing (default: false)")
	cli.Flags().Int(FlagPauseSkip, 0, "If set, the first N blocks won't be paused (default: 0)")
}

var (
	_ module.AppModuleBasic = AppModule{}
	_ module.HasServices    = AppModule{}
	_ appmodule.AppModule   = AppModule{}
)

type Pausable interface {
	// Pause pauses the BaseApp, preventing it from processing any further blocks.
	Pause() error
	// Resume resumes the BaseApp, allowing it to process blocks again.
	Resume() error
}

type PauseService interface {
	Pausable
	HoldIfPaused(ctx context.Context) error
}

type AppModule struct {
	pauseSvc PauseService
}

// RegisterLegacyAminoCodec registers the debug module's types.
func (am AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	if !Enabled {
		return
	}
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers interfaces and implementations of the debug module.
func (am AppModule) RegisterInterfaces(interfaceRegistry codectypes.InterfaceRegistry) {
	if !Enabled {
		return
	}
	types.RegisterInterfaces(interfaceRegistry)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the debug module.
func (am AppModule) RegisterGRPCGatewayRoutes(c client.Context, serveMux *runtime.ServeMux) {
	if !Enabled {
		return
	}
	if err := types.RegisterQueryHandlerClient(context.Background(), serveMux, types.NewQueryClient(c)); err != nil {
		panic(err)
	}
}

// NewAppModule creates a new AppModule Object
func NewAppModule(pauseSvc PauseService) AppModule {
	if !Enabled {
		pauseSvc = pauser.NewNoOpPauser()
	}
	return AppModule{
		pauseSvc: pauseSvc,
	}
}

// NewAppModuleFromOpts creates a new AppModule Object from server options
func NewAppModuleFromOpts(appOpts servertypes.AppOptions) AppModule {
	return NewAppModule(pauser.NewPauser(
		cast.ToBool(appOpts.Get(FlagPauseAllowed)),
		cast.ToBool(appOpts.Get(FlagPauseOnEachBlock)),
		cast.ToInt(appOpts.Get(FlagPauseSkip)),
	))
}

func (am AppModule) IsOnePerModuleType() {}

func (am AppModule) IsAppModule() {}

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (am AppModule) Name() string {
	return ModuleName
}

func (am AppModule) BeginBlock(ctx context.Context) error {
	if !Enabled {
		return nil
	}
	return am.pauseSvc.HoldIfPaused(ctx)
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	if !Enabled {
		return
	}

	types.RegisterQueryServer(cfg.QueryServer(), NewQueryServerImpl(am.pauseSvc))
}
