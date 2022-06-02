package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/v7/x/usc/types"
)

// InitGenesis performs module's genesis post validation, registers USC and collateral denoms and sets module params.
// Since during genesis state ValidateBasic we don't have an access to the app state (x/bank genesis in our case),
// we have to perform a full genesis validation here (in runtime).
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) error {
	// Set params
	k.SetParams(ctx, genState.Params)

	// Build USC and collateral coin denoms set
	if err := k.Init(ctx); err != nil {
		return err
	}

	return nil
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
