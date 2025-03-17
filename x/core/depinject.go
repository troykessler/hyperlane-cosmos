package core

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"

	"github.com/troykessler/hyperlane-cosmos/x/core/keeper"
	"github.com/troykessler/hyperlane-cosmos/x/core/types"

	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	modulev1 "github.com/troykessler/hyperlane-cosmos/api/core/module/v1"
)

var _ appmodule.AppModule = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

func init() {
	appmodule.Register(
		&modulev1.Module{},
		appmodule.Provide(ProvideModule),
		appmodule.Invoke(),
	)
}

type ModuleInputs struct {
	depinject.In

	Cdc          codec.Codec
	StoreService store.KVStoreService
	AddressCodec address.Codec

	Config *modulev1.Module

	BankKeeper types.BankKeeper
}

type ModuleOutputs struct {
	depinject.Out

	Module appmodule.AppModule
	Keeper *keeper.Keeper
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to governance as authority if not provided
	authority := authtypes.NewModuleAddress("gov")
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	k := keeper.NewKeeper(in.Cdc, in.AddressCodec, in.StoreService, authority.String(), in.BankKeeper)
	m := NewAppModule(in.Cdc, &k)

	return ModuleOutputs{Module: m, Keeper: &k}
}

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{}
}
