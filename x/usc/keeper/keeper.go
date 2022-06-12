package keeper

import (
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/gaia/v7/x/usc/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the usc store.
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	authKeeper types.AccountKeeper
	bankKeeper types.BankKeeper
	paramStore paramsTypes.Subspace
}

// NewKeeper creates a new usc Keeper instance.
func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, ak types.AccountKeeper, bk types.BankKeeper, ps paramsTypes.Subspace) Keeper {
	// Set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		authKeeper: ak,
		bankKeeper: bk,
		paramStore: ps,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// ActivePool returns current module's Active pool collateral balance.
func (k Keeper) ActivePool(ctx sdk.Context) sdk.Coins {
	poolAcc := k.authKeeper.GetModuleAccount(ctx, types.ActivePoolName)

	return k.bankKeeper.GetAllBalances(ctx, poolAcc.GetAddress())
}

// RedeemingPool returns current module's Redeeming pool collateral balance.
func (k Keeper) RedeemingPool(ctx sdk.Context) sdk.Coins {
	poolAcc := k.authKeeper.GetModuleAccount(ctx, types.RedeemingPoolName)

	return k.bankKeeper.GetAllBalances(ctx, poolAcc.GetAddress())
}

// ConvertCollateralsToUSC converts collateral coins to USC coin in 1:1 relation.
func (k Keeper) ConvertCollateralsToUSC(ctx sdk.Context, colCoins sdk.Coins) (sdk.Coin, error) {
	uscMeta, colMetas := k.USCMeta(ctx), k.CollateralMetasSet(ctx)

	uscCoin := uscMeta.NewZeroCoin()
	for _, colCoin := range colCoins {
		// Check if denom is supported
		colMeta, ok := colMetas[colCoin.Denom]
		if !ok {
			return sdk.Coin{}, sdkErrors.Wrapf(types.ErrUnsupportedCollateral, "denom (%s)", colCoin.Denom)
		}

		// Convert collateral -> USC, add USCs
		colConvertedCoin, err := colMeta.ConvertCoin(colCoin, uscMeta)
		if err != nil {
			return sdk.Coin{}, sdkErrors.Wrapf(types.ErrInternal, "converting collateral token (%s) to USC: %v", colCoin, err)
		}
		uscCoin = uscCoin.Add(colConvertedCoin)
	}

	return uscCoin, nil
}

// ConvertUSCToCollaterals converts USC coin to collateral coins in 1:1 relation iterating module's Active pool from the highest supply to the lowest.
// Returns converted USC (equals to input if there are no leftovers)  and collaterals coins.
func (k Keeper) ConvertUSCToCollaterals(ctx sdk.Context, uscCoin sdk.Coin) (sdk.Coin, sdk.Coins, error) {
	uscMeta, colMetas := k.USCMeta(ctx), k.CollateralMetasSet(ctx)

	// Check source denom
	if uscCoin.Denom != uscMeta.Denom {
		return sdk.Coin{}, nil, sdkErrors.Wrapf(types.ErrInvalidUSC, "got (%s), expected (%s)", uscCoin.Denom, uscMeta.Denom)
	}

	// Sort active pool coins from the highest supply to the lowest normalizing amounts
	poolCoins := k.ActivePool(ctx)

	poolCoinsNormalized := make(sdk.Coins, 0, len(poolCoins))
	for _, poolCoin := range poolCoins {
		poolMeta, ok := colMetas[poolCoin.Denom]
		if !ok {
			k.Logger(ctx).Info("Collateral meta not found for ActivePool coin (skip)", "denom", poolCoin.Denom)
			continue
		}

		poolCoinNormalized, err := poolMeta.NormalizeCoin(poolCoin, uscMeta)
		if err != nil {
			return sdk.Coin{}, nil, sdkErrors.Wrapf(types.ErrInternal, "normalizing ActivePool coin (%s): %v", poolCoin, err)
		}
		poolCoinsNormalized = append(poolCoinsNormalized, poolCoinNormalized)
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

	// Fill up the desired USC amount with the current Active pool collateral supply
	uscLeftToFillCoin, colCoins := uscCoin, sdk.NewCoins()
	for _, poolCoin := range poolCoins {
		poolMeta, _ := colMetas[poolCoin.Denom] // no need to check the error, since it is checked above

		// Convert collateral -> USC to make it comparable (no amt loss here, since USC decimals are always GTE collateral's)
		poolConvertedCoin, err := poolMeta.ConvertCoin(poolCoin, uscMeta)
		if err != nil {
			return sdk.Coin{}, nil, sdkErrors.Wrapf(types.ErrInternal, "converting pool token (%s) to USC: %v", poolCoin, err)
		}

		// Define USC left to fill reduce amount (how much could be covered by this collateral)
		uscReduceCoin := uscLeftToFillCoin
		if poolConvertedCoin.IsLT(uscReduceCoin) {
			uscReduceCoin = poolConvertedCoin
		}

		// Convert USC reduce amount to collateral (amt loss could happen here)
		colCoin, err := uscMeta.ConvertCoin(uscReduceCoin, poolMeta)
		if err != nil {
			return sdk.Coin{}, nil, sdkErrors.Wrapf(types.ErrInternal, "converting USC reduce token (%s) to collateral denom (%s): %v", uscReduceCoin, poolMeta.Denom, err)
		}

		// Skip the current collateral if its supply can't cover USC reduce amount, try the next one
		if colCoin.Amount.IsZero() {
			continue
		}

		// Apply current results
		uscLeftToFillCoin = uscLeftToFillCoin.Sub(uscReduceCoin)
		colCoins = colCoins.Add(colCoin)

		// Check if redeem amount is filled up
		if uscLeftToFillCoin.IsZero() {
			break
		}
	}

	return uscCoin.Sub(uscLeftToFillCoin), colCoins, nil
}
