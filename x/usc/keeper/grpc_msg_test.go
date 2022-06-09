package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/v7/x/usc/types"
)

func (s *TestSuite) TestMsgMintUSC() {
	assert, require := s.Assert(), s.Require()
	bankKeeper, uscKeeper := s.app.BankKeeper, s.app.USCKeeper
	accAddr := s.accAddrs[0]

	busdSwapAmt, _ := sdk.NewIntFromString("1000000000000000000") // 1.0
	usdtSwapAmt, _ := sdk.NewIntFromString("1000000")             // 1.0
	usdcSwapAmt, _ := sdk.NewIntFromString("1000")                // 1.0

	accSwapCoins := sdk.NewCoins(
		sdk.NewCoin("abusd", busdSwapAmt),
		sdk.NewCoin("uusdt", usdtSwapAmt),
		sdk.NewCoin("musdc", usdcSwapAmt),
	)

	uscAmt, _ := sdk.NewIntFromString("3000000000000000000") // 1.0
	uscExpected := sdk.NewCoin("ausc", uscAmt)

	accBalanceExpected := GenCoins.Sub(accSwapCoins).Add(uscExpected)

	msg := types.NewMsgMintUSC(accAddr, accSwapCoins)
	require.NoError(msg.ValidateBasic())

	// Send Msg
	res, err := s.msgServer.MintUSC(sdk.WrapSDKContext(s.ctx), msg)
	require.NoError(err)

	// Check result
	require.NotNil(res)
	assert.Equal(uscExpected.String(), res.MintedAmount.String())

	// Check account balance
	assert.Equal(accBalanceExpected.String(), bankKeeper.GetAllBalances(s.ctx, accAddr).String())

	// Check pools
	assert.Equal(accSwapCoins.String(), uscKeeper.ActivePool(s.ctx).String())
	assert.True(uscKeeper.RedeemingPool(s.ctx).IsZero())
}
