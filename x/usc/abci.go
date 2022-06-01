package usc

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/v7/x/usc/keeper"
	"github.com/cosmos/gaia/v7/x/usc/types"
)

// BeginBlocker performs a post keeper initialization.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	if err := k.Init(ctx); err != nil {
		panic(fmt.Errorf("module x/%s: post keeper initialization: %w", types.ModuleName, err))
	}
}
