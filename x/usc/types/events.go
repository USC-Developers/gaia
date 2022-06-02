package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeRedeemQueued = "collateral_redeem_queued"
	EventTypeRedeemDone   = "collateral_redeem_done"

	AttributeKeyCompletionTime = "completion_time"
)

// NewRedeemQueuedEvent creates a new redeem enqueue event.
func NewRedeemQueuedEvent(accAddr sdk.AccAddress, amount sdk.Coins, completionTime time.Time) sdk.Event {
	return sdk.NewEvent(
		EventTypeRedeemQueued,
		sdk.NewAttribute(sdk.AttributeKeySender, accAddr.String()),
		sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
		sdk.NewAttribute(AttributeKeyCompletionTime, completionTime.String()),
	)
}

// NewRedeemDoneEvent creates a new redeem dequeue event.
func NewRedeemDoneEvent(accAddr sdk.AccAddress, amount sdk.Coins, completionTime time.Time) sdk.Event {
	return sdk.NewEvent(
		EventTypeRedeemDone,
		sdk.NewAttribute(sdk.AttributeKeySender, accAddr.String()),
		sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
		sdk.NewAttribute(AttributeKeyCompletionTime, completionTime.String()),
	)
}
