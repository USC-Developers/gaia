package keeper

import (
	"context"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gaia/v7/x/usc/types"
)

var _ types.MsgServer = (*msgServer)(nil)

// msgServer implements the gRPC SDK messages service.
type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the types.MsgServer interface.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// MintUSC implements the types.MsgServer interface.
func (k msgServer) MintUSC(goCtx context.Context, req *types.MsgMintUSC) (*types.MsgMintUSCResponse, error) {
	if req == nil {
		return nil, sdkErrors.Wrapf(types.ErrInternal, "req is nil")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check input
	accAddr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	collateralDenomSet := k.CollateralDenomSet(ctx)
	for _, colCoin := range req.CollateralAmount {
		if _, ok := collateralDenomSet[colCoin.Denom]; !ok {
			return nil, sdkErrors.Wrapf(types.ErrUnsupportedCollateral, "denom (%s)", colCoin.Denom)
		}
	}

	// Convert collateral coins 1:1 to USC coin
	uscCoin := sdk.NewCoin(k.USCDenom(ctx), sdk.ZeroInt())
	for _, colCoin := range req.CollateralAmount {
		colConvertedCoin, err := sdk.ConvertCoin(colCoin, uscCoin.Denom)
		if err != nil {
			return nil, sdkErrors.Wrapf(types.ErrInternal, "converting collateral denom (%s) to USC: %v", colCoin.Denom, err)
		}

		uscCoin = uscCoin.Add(colConvertedCoin)
	}

	// Transfer account's collateral coins to the module's pool
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, accAddr, types.ModuleName, req.CollateralAmount); err != nil {
		return nil, err
	}

	// Mint USC coin and transfer to client's account
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(uscCoin)); err != nil {
		return nil, sdkErrors.Wrapf(types.ErrInternal, "minting USC coin (%s): %v", uscCoin, err)
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, accAddr, sdk.NewCoins(uscCoin)); err != nil {
		return nil, sdkErrors.Wrapf(types.ErrInternal, "sending USC coin (%s) from module to account: %v", uscCoin, err)
	}

	return &types.MsgMintUSCResponse{
		MintedAmount: uscCoin,
	}, nil
}

// RedeemCollateral implements the types.MsgServer interface.
func (k msgServer) RedeemCollateral(goCtx context.Context, req *types.MsgRedeemCollateral) (*types.MsgRedeemCollateralResponse, error) {
	if req == nil {
		return nil, sdkErrors.Wrapf(types.ErrInternal, "req is nil")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check input
	accAddr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	uscDenom := k.USCDenom(ctx)
	if req.UscAmount.Denom != uscDenom {
		return nil, sdkErrors.Wrapf(types.ErrInvalidUSC, "got (%s), expected (%s)", req.UscAmount.Denom, uscDenom)
	}

	// Convert USC coin 1:1 to collateral coins iterating module's pool from the highest supply to the lowest
	poolCoins := k.Pool(ctx)

	poolCoinsNormalized := make(sdk.Coins, len(poolCoins))
	for _, poolCoin := range poolCoins {
		poolCoinsNormalized = append(poolCoinsNormalized, sdk.NormalizeCoin(poolCoin))
	}
	sort.Slice(poolCoins, func(i, j int) bool {
		if poolCoinsNormalized[i].Amount.GT(poolCoinsNormalized[j].Amount) {
			return true
		}
		if poolCoinsNormalized[i].Amount.Equal(poolCoinsNormalized[j].Amount) && poolCoinsNormalized[i].Denom > poolCoinsNormalized[j].Denom {
			return true
		}
		return false
	})

	uscRedeemCoin := req.UscAmount
	colCoins := sdk.NewCoins()
	for _, poolCoin := range poolCoins {
		poolConvertedCoin, err := sdk.ConvertCoin(poolCoin, uscDenom)
		if err != nil {
			return nil, sdkErrors.Wrapf(types.ErrInternal, "converting pool denom (%s) to USC: %v", poolCoin.Denom, err)
		}

		// Reduce USC redeem amount
		uscSubCoin := uscRedeemCoin
		if poolConvertedCoin.IsLT(uscSubCoin) {
			uscSubCoin = poolConvertedCoin
		}
		uscRedeemCoin = uscRedeemCoin.Sub(uscSubCoin)

		// Convert sub amount to collateral
		colCoin, err := sdk.ConvertCoin(uscSubCoin, poolCoin.Denom)
		if err != nil {
			return nil, sdkErrors.Wrapf(types.ErrInternal, "converting USC back to collateral denom (%s): %v", poolCoin.Denom, err)
		}
		colCoins = colCoins.Add(colCoin)

		// Check if redeem amount is filled up
		if uscRedeemCoin.IsZero() {
			break
		}
	}
	if !uscRedeemCoin.IsZero() {
		return nil, sdkErrors.Wrapf(types.ErrInternal, "usc amount (%s) can not be filled up with pool (%s)", req.UscAmount, poolCoins)
	}

	// Transfer account's USC coin to the module's pool
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, accAddr, types.ModuleName, sdk.NewCoins(req.UscAmount)); err != nil {
		return nil, err
	}

	// Burn USC coin and transfer collaterals to client's account
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(req.UscAmount)); err != nil {
		return nil, sdkErrors.Wrapf(types.ErrInternal, "burning USC coin (%s): %v", req.UscAmount, err)
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, accAddr, colCoins); err != nil {
		return nil, sdkErrors.Wrapf(types.ErrInternal, "sending collateral coins (%s) from module to account: %v", colCoins, err)
	}

	return &types.MsgRedeemCollateralResponse{
		RedeemedAmount: colCoins,
	}, nil
}
