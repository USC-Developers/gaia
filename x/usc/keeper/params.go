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

// CollateralMetas returns supported collateral token metas.
func (k Keeper) CollateralMetas(ctx sdk.Context) (res []types.TokenMeta) {
	k.paramStore.Get(ctx, types.ParamsKeyCollateralMetas, &res)
	return
}

// CollateralMetasSet returns supported collateral token metas set (key: denom).
func (k Keeper) CollateralMetasSet(ctx sdk.Context) map[string]types.TokenMeta {
	metas := k.CollateralMetas(ctx)

	set := make(map[string]types.TokenMeta, len(metas))
	for _, meta := range metas {
		set[meta.Denom] = meta
	}

	return set
}

// USCMeta returns the USC token meta.
func (k Keeper) USCMeta(ctx sdk.Context) (res types.TokenMeta) {
	k.paramStore.Get(ctx, types.ParamsKeyUSCMeta, &res)
	return
}

// GetParams returns all module parameters.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.RedeemDur(ctx),
		k.CollateralMetas(ctx),
		k.USCMeta(ctx),
	)
}

// SetParams sets all module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramStore.SetParamSet(ctx, &params)
}
