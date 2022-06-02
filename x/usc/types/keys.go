package types

import (
	"bytes"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the module name.
	ModuleName = "usc"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the msg router key for the module.
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the module.
	QuerierRoute = ModuleName

	// ActivePoolName defines module name for storing collateral coins.
	ActivePoolName = "usc_active_pool"

	// RedeemingPoolName defines module name for storing collateral coins which are queued to be redeemed.
	RedeemingPoolName = "usc_redeeming_pool"
)

var (
	// RedeemingQueueKey prefix for timestamps in redeeming queue.
	RedeemingQueueKey = []byte{0x10}
)

// GetRedeemingQueueKey creates storage key for all redeem request for a specific timeSlice.
func GetRedeemingQueueKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)

	return append(RedeemingQueueKey, bz...)
}

// ParseRedeemingQueueKey parses timeSlice storage key.
func ParseRedeemingQueueKey(key []byte) time.Time {
	if len(key) == 0 {
		panic(fmt.Errorf("parsing timeSlice key (%v): empty key"))
	}

	prefix, bz := key[:1], key[1:]
	if !bytes.Equal(prefix, RedeemingQueueKey) {
		panic(fmt.Errorf("parsing timeSlice key (%v): unexpected prefix", key))
	}

	timestamp, err := sdk.ParseTimeBytes(bz)
	if err != nil {
		panic(fmt.Errorf("parsing timeSlice key (%v): %w", key, err))
	}

	return timestamp
}
