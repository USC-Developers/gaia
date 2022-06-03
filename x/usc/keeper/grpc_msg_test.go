package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/v7/x/usc/types"
)

func (s *TestSuite) TestMsgMintUSC() {
	assert, require := s.Assert(), s.Require()
	bankKeeper, uscKeeper := s.app.BankKeeper, s.app.USCKeeper
	accAddr := s.accAddrs[0]

	coinsToSwap := sdk.NewCoins(
		sdk.NewCoin("busd", sdk.NewInt(10)),
	)

	msg := types.NewMsgMintUSC(accAddr, coinsToSwap)
	require.NoError(msg.ValidateBasic())

	// Send Msg
	res, err := s.msgServer.MintUSC(sdk.WrapSDKContext(s.ctx), msg)
	require.NoError(err)

	// Check result
	require.NotNil(res)
	assert.Equal("10usc", res.MintedAmount.String())

	// Check account balance
	assert.EqualValues(10, bankKeeper.GetBalance(s.ctx, accAddr, "usc").Amount.Int64())
	assert.EqualValues(90, bankKeeper.GetBalance(s.ctx, accAddr, "busd").Amount.Int64())

	// Check pools
	assert.Equal("10busd", uscKeeper.ActivePool(s.ctx).String())
	assert.True(uscKeeper.RedeemingPool(s.ctx).IsZero())
}
