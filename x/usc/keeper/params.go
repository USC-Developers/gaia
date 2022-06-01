package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/v7/x/usc/types"
)

// RedeemDur returns USC -> collateral coins redeem timeout duration.
func (k Keeper) RedeemDur(ctx sdk.Context) (res time.Duration) {
	k.paramStore.Get(ctx, types.ParamsKeyRedeemDur, &res)
	return
}

// CollateralDenoms returns supported collateral coin denoms.
func (k Keeper) CollateralDenoms(ctx sdk.Context) (res []string) {
	k.paramStore.Get(ctx, types.ParamsKeyCollateralDenoms, &res)
	return
}

// CollateralDenomSet returns supported collateral coin denoms set.
func (k Keeper) CollateralDenomSet(ctx sdk.Context) map[string]struct{} {
	denoms := k.CollateralDenoms(ctx)

	set := make(map[string]struct{}, len(denoms))
	for _, denom := range denoms {
		set[denom] = struct{}{}
	}

	return set
}

// USCDenom returns the USC coin denom.
func (k Keeper) USCDenom(ctx sdk.Context) (res string) {
	k.paramStore.Get(ctx, types.ParamsKeyUSCDenom, &res)
	return
}

// GetParams returns all module parameters.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.RedeemDur(ctx),
		k.CollateralDenoms(ctx),
		k.USCDenom(ctx),
	)
}

// SetParams sets all module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramStore.SetParamSet(ctx, &params)
}
