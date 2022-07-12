package keeper

import (
	"sort"
	"strings"

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

// USCMetaDerivative returns USC token meta for the native coin or from the list of supported IBC derivatives.
func (k Keeper) USCMetaDerivative(ctx sdk.Context, denom string) (types.TokenMeta, error) {
	uscMeta := k.USCMeta()

	if denom == uscMeta.Denom {
		return uscMeta, nil
	}

	uscIbcDenoms := k.USCIbcDenoms(ctx)
	for _, ibcDenom := range uscIbcDenoms {
		if denom != ibcDenom {
			continue
		}

		uscMeta.Denom = ibcDenom
		return uscMeta, nil
	}

	return types.TokenMeta{}, sdkErrors.Wrapf(
		types.ErrInvalidUSC,
		"got (%s), expected (%s / [%s])",
		denom,
		uscMeta.Denom, strings.Join(uscIbcDenoms, ","),
	)
}

// ConvertCollateralsToUSC converts collateral coins to USC coin in 1:1 relation.
func (k Keeper) ConvertCollateralsToUSC(ctx sdk.Context, colCoins sdk.Coins) (uscCoin sdk.Coin, colUsedCoins sdk.Coins, retErr error) {
	uscMeta, colMetas := k.USCMeta(), k.CollateralMetasSet(ctx)

	uscCoin = uscMeta.NewZeroCoin()
	for _, colCoin := range colCoins {
		// Check if denom is supported
		colMeta, ok := colMetas[colCoin.Denom]
		if !ok {
			supportedColDenoms := make([]string, 0, len(colMetas))
			for colDenom := range colMetas {
				supportedColDenoms = append(supportedColDenoms, colDenom)
			}

			retErr = sdkErrors.Wrapf(types.ErrUnsupportedCollateral,
				"got (%s), expected ([%s])",
				colCoin.Denom,
				strings.Join(supportedColDenoms, ","),
			)
			return
		}

		// Convert collateral -> USC, note actually used collateral amount
		colConvertedCoin, colUsedCoin, err := colMeta.ConvertCoin2(colCoin, uscMeta)
		if err != nil {
			retErr = sdkErrors.Wrapf(types.ErrInternal, "converting collateral token (%s) to USC: %v", colCoin, err)
			return
		}
		uscCoin = uscCoin.Add(colConvertedCoin)
		colUsedCoins = colUsedCoins.Add(colUsedCoin)
	}

	return
}

// ConvertUSCToCollaterals converts USC coin to collateral coins in 1:1 relation iterating module's Active pool from the highest supply to the lowest.
// Returns converted USC (equals to input if there are no leftovers)  and collaterals coins.
func (k Keeper) ConvertUSCToCollaterals(ctx sdk.Context, uscCoin sdk.Coin) (uscUsedCoin sdk.Coin, colCoins sdk.Coins, retErr error) {
	colMetas := k.CollateralMetasSet(ctx)

	// Check source denom and get USC meta
	uscMeta, err := k.USCMetaDerivative(ctx, uscCoin.Denom)
	if err != nil {
		retErr = err
		return
	}

	// Sort active pool coins from the highest supply to the lowest normalizing amounts
	poolCoins := k.ActivePool(ctx)

	baseMeta := k.BaseMeta(ctx)
	poolCoinsNormalizedSet := make(map[string]sdk.Int)
	for _, poolCoin := range poolCoins {
		poolMeta, ok := colMetas[poolCoin.Denom]
		if !ok {
			k.Logger(ctx).Info("Collateral meta not found for ActivePool coin (skip)", "denom", poolCoin.Denom)
			continue
		}

		normalizedCoin, err := poolMeta.NormalizeCoin(poolCoin, baseMeta)
		if err != nil {
			retErr = sdkErrors.Wrapf(types.ErrInternal, "normalizing ActivePool coin (%s): %v", poolCoin, err)
			return
		}
		poolCoinsNormalizedSet[poolCoin.Denom] = normalizedCoin.Amount
	}
	sort.Slice(poolCoins, func(i, j int) bool {
		iDenom, jDenom := poolCoins[i].Denom, poolCoins[j].Denom
		iAmt, jAmt := poolCoinsNormalizedSet[poolCoins[i].Denom], poolCoinsNormalizedSet[poolCoins[j].Denom]

		if iAmt.GT(jAmt) {
			return true
		}
		if iAmt.Equal(jAmt) && iDenom > jDenom {
			return true
		}

		return false
	})

	// Fill up the desired USC amount with the current Active pool collateral supply
	uscLeftToFillCoin := uscCoin
	for _, poolCoin := range poolCoins {
		poolMeta := colMetas[poolCoin.Denom] // no need to check the error, since it is checked above

		// Convert collateral -> USC to make it comparable
		poolConvertedCoin, err := poolMeta.ConvertCoin(poolCoin, uscMeta)
		if err != nil {
			retErr = sdkErrors.Wrapf(types.ErrInternal, "converting pool token (%s) to USC: %v", poolCoin, err)
			return
		}

		// Define USC left to fill reduce amount (how much could be covered by this collateral)
		uscReduceCoin := uscLeftToFillCoin
		if poolConvertedCoin.IsLT(uscReduceCoin) {
			uscReduceCoin = poolConvertedCoin
		}

		// Convert USC reduce amount to collateral
		colCoin, uscReduceUsedCoin, err := uscMeta.ConvertCoin2(uscReduceCoin, poolMeta)
		if err != nil {
			retErr = sdkErrors.Wrapf(types.ErrInternal, "converting USC reduce token (%s) to collateral denom (%s): %v", uscReduceCoin, poolMeta.Denom, err)
			return
		}

		// Skip the current collateral if its supply can't cover USC reduce amount, try the next one
		if colCoin.Amount.IsZero() {
			continue
		}

		// Apply current results
		uscLeftToFillCoin = uscLeftToFillCoin.Sub(uscReduceUsedCoin)
		colCoins = colCoins.Add(colCoin)

		// Check if redeem amount is filled up
		if uscLeftToFillCoin.IsZero() {
			break
		}
	}
	uscUsedCoin = uscCoin.Sub(uscLeftToFillCoin)

	return
}
