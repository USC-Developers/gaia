package usc

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/v7/x/usc/keeper"
	"github.com/cosmos/gaia/v7/x/usc/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker performs a post keeper initialization.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	if err := k.Init(ctx); err != nil {
		panic(fmt.Errorf("x/%s module BeginBlocker: %w", types.ModuleName, err))
	}
}

// EndBlocker manages redeeming queue.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	k.EndRedeeming(ctx)

	return []abci.ValidatorUpdate{}
}
