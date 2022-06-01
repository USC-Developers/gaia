package keeper

import (
	"fmt"
	"math"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

// Pool returns current module pool collateral balance.
func (k Keeper) Pool(ctx sdk.Context) sdk.Coins {
	poolAcc := k.authKeeper.GetModuleAccount(ctx, types.ModuleName)

	return k.bankKeeper.GetAllBalances(ctx, poolAcc.GetAddress())
}
