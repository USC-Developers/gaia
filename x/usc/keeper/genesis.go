package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/v7/x/usc/types"
)

// InitGenesis initializes the module genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the current module genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := types.NewParams(
		k.RedeemDur(ctx),
		k.CollateralDenoms(ctx),
		k.USCDenom(ctx),
	)

	return &types.GenesisState{
		Params: params,
	}
}
