package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/v7/x/usc/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUSCKeeperMsgMintUSC(t *testing.T) {
	type testCase struct {
		name           string
		colCoinsToSwap string
		//
		errExpected    error
		uscAmtExpected string
	}

	testCases := []testCase{
		{
			name:           "Mint 3.0 USC tokens",
			colCoinsToSwap: "1000000000000000000abusd,1000000uusdt,1000musdc", // 1.0 BUSD, 1.0 USDT, 1.0 USDC
			errExpected:    nil,
			uscAmtExpected: "3000000000000000000", // 3.0 USC
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Fixtures
			te := NewTestEnv(t)

			accAddr, accCoins := te.AddAccount(t, tc.colCoinsToSwap)
			swapColCoins := accCoins

			// Send msg
			msg := types.NewMsgMintUSC(accAddr, swapColCoins)
			require.NoError(t, msg.ValidateBasic())

			res, err := te.msgServer.MintUSC(sdk.WrapSDKContext(te.ctx), msg)
			if tc.errExpected != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.errExpected)

				return
			}
			require.NoError(t, err)
			require.NotNil(t, res)

			// Verify minted USC
			uscAmtExpected, ok := sdk.NewIntFromString(tc.uscAmtExpected)
			require.True(t, ok)
			uscMintedCoinExpected := sdk.NewCoin(types.DefaultUSCDenom, uscAmtExpected)

			assert.Equal(t,
				uscMintedCoinExpected.String(),
				res.MintedAmount.String(),
			)

			// Verify account balance
			assert.Equal(t, uscMintedCoinExpected.String(),
				te.app.BankKeeper.GetBalance(te.ctx, accAddr, types.DefaultUSCDenom).String(),
			)

			// Verify Active pool balance
			assert.Equal(t,
				swapColCoins.String(),
				te.app.USCKeeper.ActivePool(te.ctx).String(),
			)
		})
	}
}

func TestUSCKeeperMsgRedeemCollateral(t *testing.T) {
	type testCase struct {
		name            string
		uscAmtToRedeem  string
		activePoolCoins string
		//
		errExpected              error
		uscAmtLeftExpected       string
		colCoinsRedeemedExpected string
	}

	testCases := []testCase{
		{
			name:                     "Declined redeem: USC amount is too small",
			uscAmtToRedeem:           "100000",         // 0.0000000000001 USC
			activePoolCoins:          "100000000uusdt", // 100.000000 USDT
			errExpected:              types.ErrRedeemDeclined,
			uscAmtLeftExpected:       "100000", // same amount
			colCoinsRedeemedExpected: "",
		},
		{
			name:                     "Partially filled",
			uscAmtToRedeem:           "10020000000000000", // 0.010020 USC
			activePoolCoins:          "5musdc,10uusdt",    // 0.005 USDC, 0.000010 USDT
			errExpected:              nil,
			uscAmtLeftExpected:       "5010000000000000", // 0.005010 USC
			colCoinsRedeemedExpected: "5musdc,10uusdt",
		},
		{
			name:                     "Fully filled",
			uscAmtToRedeem:           "130000000000000000000",                              // 130.0 USC
			activePoolCoins:          "75000000000000000000abusd,50000000uusdt,25000musdc", // 75.0 BUSD, 50.0 USDT, 25.0 USDC,
			errExpected:              nil,
			uscAmtLeftExpected:       "0", // none
			colCoinsRedeemedExpected: "75000000000000000000abusd,50000000uusdt,5000musdc",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Fixtures
			te := NewTestEnv(t)

			accAddr, accCoins := te.AddAccount(t, tc.uscAmtToRedeem+types.DefaultUSCDenom)
			uscRedeemCoin := accCoins[0]

			activePoolCoins := te.AddActivePoolBalance(t, tc.activePoolCoins)

			// Send msg
			msg := types.NewMsgRedeemCollateral(accAddr, uscRedeemCoin)
			require.NoError(t, msg.ValidateBasic())

			res, err := te.msgServer.RedeemCollateral(sdk.WrapSDKContext(te.ctx), msg)
			if tc.errExpected != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.errExpected)

				// Ensure balances have not chagned
				assert.Equal(t,
					accCoins.String(),
					te.app.BankKeeper.GetAllBalances(te.ctx, accAddr).String(),
				)
				assert.Equal(t,
					activePoolCoins.String(),
					te.app.USCKeeper.ActivePool(te.ctx).String(),
				)
				assert.True(t,
					te.app.USCKeeper.RedeemingPool(te.ctx).IsZero(),
				)

				return
			}
			require.NoError(t, err)
			require.NotNil(t, res)

			// Build expect value
			uscLeftAmtExpected, ok := sdk.NewIntFromString(tc.uscAmtLeftExpected)
			require.True(t, ok)
			uscLeftCoinExpected := sdk.NewCoin(types.DefaultUSCDenom, uscLeftAmtExpected)

			uscBurnedExpected := uscRedeemCoin.Sub(uscLeftCoinExpected)
			completionTimeExpected := MockTimestamp.Add(te.app.USCKeeper.RedeemDur(te.ctx))

			colRedeemedCoinsExpected, err := sdk.ParseCoinsNormalized(tc.colCoinsRedeemedExpected)
			require.NoError(t, err)

			// Verify the result
			assert.Equal(t,
				uscBurnedExpected.String(),
				res.BurnedAmount.String(),
			)
			assert.Equal(t,
				colRedeemedCoinsExpected.String(),
				sdk.NewCoins(res.RedeemedAmount...).String(),
			)
			assert.EqualValues(t,
				completionTimeExpected,
				res.CompletionTime,
			)

			// Verify account balance
			assert.Equal(t,
				uscLeftCoinExpected.String(),
				te.app.BankKeeper.GetBalance(te.ctx, accAddr, types.DefaultUSCDenom).String(),
			)

			// Verify Active pool balance
			assert.Equal(t,
				activePoolCoins.Sub(colRedeemedCoinsExpected).String(),
				te.app.USCKeeper.ActivePool(te.ctx).String(),
			)

			// Verify Redeeming pool balance
			assert.Equal(t,
				colRedeemedCoinsExpected.String(),
				te.app.USCKeeper.RedeemingPool(te.ctx).String(),
			)
		})
	}
}
