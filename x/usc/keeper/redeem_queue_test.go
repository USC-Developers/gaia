package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestUSCKeeperRedeemingQueue(t *testing.T) {
	const blockDur = 5 * time.Second

	te := NewTestEnv(t)
	keeper, bankKeeper := te.app.USCKeeper, te.app.BankKeeper

	// Fixtures
	accAddr1, _ := te.AddAccount(t, "")
	redeemCoins1 := sdk.NewCoins(
		sdk.NewCoin("uusdt", sdk.NewInt(1000)),
		sdk.NewCoin("musdc", sdk.NewInt(100)),
	)

	accAddr2, _ := te.AddAccount(t, "")
	redeemCoins2 := sdk.NewCoins(
		sdk.NewCoin("abusd", sdk.NewInt(100000)),
	)

	accAddr3, _ := te.AddAccount(t, "")
	redeemCoins3 := sdk.NewCoins(
		sdk.NewCoin("abusd", sdk.NewInt(1)),
		sdk.NewCoin("musdc", sdk.NewInt(10)),
	)

	curRedeemingPoolCoins := redeemCoins1.Add(redeemCoins2...).Add(redeemCoins3...)
	te.AddRedeemingPoolBalance(t, curRedeemingPoolCoins.String())

	ctx, redeemDur := te.ctx, keeper.RedeemDur(te.ctx)
	var redeemTimestamp1, redeemTimestamp2 time.Time

	// Block 1: set 1st timeSlice with 2 entries
	{
		ctx = ctx.WithBlockTime(ctx.BlockTime().Add(blockDur))
		redeemTimestamp1 = ctx.BlockTime().Add(redeemDur)

		assert.Equal(t,
			redeemTimestamp1,
			keeper.BeginRedeeming(ctx, accAddr1, redeemCoins1),
		)
		assert.Equal(t,
			redeemTimestamp1,
			keeper.BeginRedeeming(ctx, accAddr2, redeemCoins2),
		)
	}

	// Block 2: set 2nd timeSlice with 1 entry
	{
		ctx = ctx.WithBlockTime(ctx.BlockTime().Add(blockDur))
		redeemTimestamp2 = ctx.BlockTime().Add(redeemDur)

		assert.Equal(t,
			redeemTimestamp2,
			keeper.BeginRedeeming(ctx, accAddr3, redeemCoins3),
		)
	}

	// Block 3: ensure that none of the queue entries were triggered
	{
		ctx = ctx.WithBlockTime(ctx.BlockTime().Add(blockDur))

		keeper.EndRedeeming(ctx)
		assert.True(t, bankKeeper.GetAllBalances(ctx, accAddr1).IsZero())
		assert.True(t, bankKeeper.GetAllBalances(ctx, accAddr2).IsZero())
		assert.True(t, bankKeeper.GetAllBalances(ctx, accAddr3).IsZero())

		assert.EqualValues(t, curRedeemingPoolCoins, keeper.RedeemingPool(ctx))
	}

	// Block X (redeemTimestamp1): check that 2 queue entries were processed
	{
		ctx = ctx.WithBlockTime(redeemTimestamp1)

		keeper.EndRedeeming(ctx)
		assert.EqualValues(t, redeemCoins1, bankKeeper.GetAllBalances(ctx, accAddr1))
		assert.EqualValues(t, redeemCoins2, bankKeeper.GetAllBalances(ctx, accAddr2))
		assert.True(t, bankKeeper.GetAllBalances(ctx, accAddr3).IsZero())

		curRedeemingPoolCoins = curRedeemingPoolCoins.Sub(redeemCoins1).Sub(redeemCoins2)
		assert.EqualValues(t, curRedeemingPoolCoins, keeper.RedeemingPool(ctx))
	}

	// Block Y (redeemTimestamp2 + 1 second): check that 1 queue entry was processed
	{
		ctx = ctx.WithBlockTime(redeemTimestamp2.Add(1 * time.Second))

		keeper.EndRedeeming(ctx)
		assert.EqualValues(t, redeemCoins1, bankKeeper.GetAllBalances(ctx, accAddr1))
		assert.EqualValues(t, redeemCoins2, bankKeeper.GetAllBalances(ctx, accAddr2))
		assert.EqualValues(t, redeemCoins3, bankKeeper.GetAllBalances(ctx, accAddr3))

		curRedeemingPoolCoins = sdk.NewCoins()
		assert.EqualValues(t, curRedeemingPoolCoins, keeper.RedeemingPool(ctx))
	}

	// Block Z: check that the queue doesn't panic when emptied (just in case)
	{
		ctx = ctx.WithBlockTime(ctx.BlockTime().Add(blockDur))

		keeper.EndRedeeming(ctx)
	}
}
