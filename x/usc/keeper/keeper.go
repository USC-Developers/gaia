package keeper

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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
	//
	initialized bool
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

// Init performs module's genesis validation and registers USC and collateral denoms.
// Since AppModuleBasic.ValidateGenesis doesn't have an app state (x/bank genesis in our case), we have to
// perform a full genesis validation here (in runtime).
func (k Keeper) Init(ctx sdk.Context) error {
	if k.initialized {
		return nil
	}

	// Build USC and collateral coin denoms set
	uscDenom := k.USCDenom(ctx)
	targetDenomSet := k.CollateralDenomSet(ctx)
	targetDenomSet[uscDenom] = struct{}{}

	// Iterate over all registered x/bank metadata and ensure target denoms are registered
	uscExponent, minExponent := uint32(0), uint32(math.MaxUint32)
	k.bankKeeper.IterateAllDenomMetaData(ctx, func(meta banktypes.Metadata) bool {
		for _, unit := range meta.DenomUnits {
			if _, ok := targetDenomSet[unit.Denom]; !ok {
				continue
			}

			sdk.RegisterDenom(unit.Denom, sdk.NewDecWithPrec(1, int64(unit.Exponent)))

			delete(targetDenomSet, unit.Denom)
			if unit.Denom == uscDenom {
				uscExponent = unit.Exponent
			}
			if unit.Exponent < minExponent {
				minExponent = unit.Exponent
			}
		}
		return false
	})
	if len(targetDenomSet) > 0 {
		missingDenoms := make([]string, 0, len(targetDenomSet))
		for denom := range targetDenomSet {
			missingDenoms = append(missingDenoms, denom)
		}
		return fmt.Errorf("x/bank metadata not found for denoms: [%s]", strings.Join(missingDenoms, ", "))
	}

	// Check that USC precision is high enough
	if uscExponent < minExponent {
		return fmt.Errorf("usc precision (%d) must be GTE that min collateral precision (%d)", uscExponent, minExponent)
	}

	k.initialized = true

	return nil
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
	uscCoin := sdk.NewCoin(k.USCDenom(ctx), sdk.ZeroInt())
	collateralDenomSet := k.CollateralDenomSet(ctx)
	for _, colCoin := range colCoins {
		// Check if denom is supported
		if _, ok := collateralDenomSet[colCoin.Denom]; !ok {
			return sdk.Coin{}, sdkErrors.Wrapf(types.ErrUnsupportedCollateral, "denom (%s)", colCoin.Denom)
		}

		// Convert collateral -> USC, add USCs
		colConvertedCoin, err := sdk.ConvertCoin(colCoin, uscCoin.Denom)
		if err != nil {
			return sdk.Coin{}, sdkErrors.Wrapf(types.ErrInternal, "converting collateral denom (%s) to USC: %v", colCoin.Denom, err)
		}
		uscCoin = uscCoin.Add(colConvertedCoin)
	}

	return uscCoin, nil
}

// ConvertUSCToCollaterals converts USC coin to collateral coins in 1:1 relation iterating module's Active pool from the highest supply to the lowest.
func (k Keeper) ConvertUSCToCollaterals(ctx sdk.Context, uscCoin sdk.Coin) (sdk.Coins, error) {
	// Check source denom
	uscDenom := k.USCDenom(ctx)
	if uscCoin.Denom != uscDenom {
		return nil, sdkErrors.Wrapf(types.ErrInvalidUSC, "got (%s), expected (%s)", uscCoin.Denom, uscDenom)
	}

	// Sort active pool coins from the highest supply to the lowest normalizing amounts
	poolCoins := k.ActivePool(ctx)

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

	// Convert and add collateral coins
	uscLeft, colCoins := uscCoin, sdk.NewCoins()
	for _, poolCoin := range poolCoins {
		// Convert collateral -> USC
		poolConvertedCoin, err := sdk.ConvertCoin(poolCoin, uscDenom)
		if err != nil {
			return nil, sdkErrors.Wrapf(types.ErrInternal, "converting pool denom (%s) to USC: %v", poolCoin.Denom, err)
		}

		// Reduce USC convert amount
		uscSubCoin := uscLeft
		if poolConvertedCoin.IsLT(uscSubCoin) {
			uscSubCoin = poolConvertedCoin
		}
		uscLeft = uscLeft.Sub(uscSubCoin)

		// Convert sub amount to collateral
		colCoin, err := sdk.ConvertCoin(uscSubCoin, poolCoin.Denom)
		if err != nil {
			return nil, sdkErrors.Wrapf(types.ErrInternal, "converting USC to collateral denom (%s): %v", poolCoin.Denom, err)
		}
		colCoins = colCoins.Add(colCoin)

		// Check if redeem amount is filled up
		if uscLeft.IsZero() {
			break
		}
	}
	if !uscLeft.IsZero() {
		return nil, sdkErrors.Wrapf(types.ErrInternal, "usc amount (%s) can not be filled up with the pool supply (%s)", uscCoin, poolCoins)
	}

	return colCoins, nil
}
