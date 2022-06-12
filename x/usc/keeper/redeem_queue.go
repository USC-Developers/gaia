package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/v7/x/usc/types"
)

// BeginRedeeming creates a new redeem entry and enqueues it.
func (k Keeper) BeginRedeeming(ctx sdk.Context, accAddr sdk.AccAddress, amount sdk.Coins) time.Time {
	redeemEntry := types.RedeemEntry{
		Address:          accAddr.String(),
		CollateralAmount: amount,
	}
	completionTime := ctx.BlockTime().Add(k.RedeemDur(ctx))

	k.InsertToRedeemQueue(ctx, redeemEntry, completionTime)

	ctx.EventManager().EmitEvent(
		types.NewRedeemQueuedEvent(accAddr, amount, completionTime),
	)

	return completionTime
}

// EndRedeeming dequeues all mature redeem entries and sends collaterals to a requester from the module's Redeeming pool.
func (k Keeper) EndRedeeming(ctx sdk.Context) {
	for _, entry := range k.DequeueAllMatureFromRedeemQueue(ctx, ctx.BlockTime()) {
		accAddr, err := sdk.AccAddressFromBech32(entry.Address)
		if err != nil {
			panic(fmt.Errorf("converting redeeming entry account (%s) from Bech32: %w", accAddr, err))
		}

		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.RedeemingPoolName, accAddr, entry.CollateralAmount); err != nil {
			panic(fmt.Errorf("sending collateral coins (%s) from module to account (%s): %w", entry.CollateralAmount, accAddr, err))
		}

		ctx.EventManager().EmitEvent(
			types.NewRedeemDoneEvent(accAddr, entry.CollateralAmount, ctx.BlockTime()),
		)
	}
}

// SetRedeemQueueTimeSlice sets redeeming queue timeSlice at a given timestamp key.
func (k Keeper) SetRedeemQueueTimeSlice(ctx sdk.Context, timestamp time.Time, entries []types.RedeemEntry) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRedeemingQueueKey(timestamp)

	bz := k.cdc.MustMarshal(&types.RedeemEntries{
		Entries: entries,
	})
	store.Set(key, bz)
}

// GetRedeemQueueTimeSlice returns redeeming queue timeSlice at a given timestamp key.
func (k Keeper) GetRedeemQueueTimeSlice(ctx sdk.Context, timestamp time.Time) []types.RedeemEntry {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRedeemingQueueKey(timestamp)

	bz := store.Get(key)
	if bz == nil {
		return []types.RedeemEntry{}
	}

	var data types.RedeemEntries
	k.cdc.MustUnmarshal(bz, &data)

	return data.Entries
}

// InsertToRedeemQueue adds redeem entry to the redeeming queue timeSlice.
func (k Keeper) InsertToRedeemQueue(ctx sdk.Context, entry types.RedeemEntry, completionTime time.Time) {
	timeSlice := k.GetRedeemQueueTimeSlice(ctx, completionTime)
	timeSlice = append(timeSlice, entry)

	k.SetRedeemQueueTimeSlice(ctx, completionTime, timeSlice)
}

// DequeueAllMatureFromRedeemQueue returns all redeeming queue entries whose timestamp key is LTE a given timestamp
// removing them from the queue.
func (k Keeper) DequeueAllMatureFromRedeemQueue(ctx sdk.Context, endTime time.Time) (matureEntries []types.RedeemEntry) {
	store := ctx.KVStore(k.storeKey)

	iterator := store.Iterator(types.RedeemingQueueKey, sdk.InclusiveEndBytes(types.GetRedeemingQueueKey(endTime)))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var data types.RedeemEntries
		k.cdc.MustUnmarshal(iterator.Value(), &data)

		matureEntries = append(matureEntries, data.Entries...)

		store.Delete(iterator.Key())
	}

	return
}

// IterateRedeemQueue iterates over all redeeming queue entries.
func (k Keeper) IterateRedeemQueue(ctx sdk.Context, fn func(timestamp time.Time, entry types.RedeemEntry) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStoreReversePrefixIterator(store, types.RedeemingQueueKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		timestamp := types.ParseRedeemingQueueKey(iterator.Key())

		var data types.RedeemEntries
		k.cdc.MustUnmarshal(iterator.Value(), &data)

		for _, entry := range data.Entries {
			if fn(timestamp, entry) {
				break
			}
		}
	}
}
