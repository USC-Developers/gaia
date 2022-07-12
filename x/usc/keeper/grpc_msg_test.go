package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gaia/v7/x/usc/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUSCKeeperMsgMintUSC(t *testing.T) {
	type testCase struct {
		name           string
		colCoinsToSwap string
		//
		errExpected          error
		uscCoinExpected      string
		colCoinsUsedExpected string
	}

	testCases := []testCase{
		{
			name:                 "OK: mint 3.0 USC tokens (no collateral leftovers)",
			colCoinsToSwap:       "1000000000nbusd,1000000uusdt,1000musdc", // 1.0 BUSD, 1.0 USDT, 1.0 USDC
			errExpected:          nil,
			uscCoinExpected:      "3000000uusc",                            // 3.0 USC
			colCoinsUsedExpected: "1000000000nbusd,1000000uusdt,1000musdc", // all
		},
		{
			name:                 "OK: mint 3.0 USC tokens (with collateral leftovers)",
			colCoinsToSwap:       "1000000100nbusd,1000000uusdt,1000musdc", // 1.000000100 BUSD, 1.0 USDT, 1.0 USDC
			errExpected:          nil,
			uscCoinExpected:      "3000000uusc",                            // 3.0 USC
			colCoinsUsedExpected: "1000000000nbusd,1000000uusdt,1000musdc", // 1.0 BUSD, 1.0 USDT, 1.0 USDC
		},
		{
			name:           "Fail: unsupported collateral",
			colCoinsToSwap: "1000000000mbusd", // 1.0 BUSD (unsupported)
			errExpected:    types.ErrUnsupportedCollateral,
		},
		{
			name:           "Fail: collateral is too small",
			colCoinsToSwap: "100nbusd", // 1.000000100 BUSD
			errExpected:    sdkErrors.ErrInsufficientFunds,
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
			uscCoinExpected, err := sdk.ParseCoinNormalized(tc.uscCoinExpected)
			require.NoError(t, err)

			assert.Equal(t,
				uscCoinExpected.String(),
				res.MintedAmount.String(),
			)

			// Verify collaterals used
			colUsedExpected, err := sdk.ParseCoinsNormalized(tc.colCoinsUsedExpected)
			require.NoError(t, err)

			assert.Equal(t,
				colUsedExpected.String(),
				sdk.NewCoins(res.CollateralsAmount...).String(),
			)

			// Verify account balance
			assert.Equal(t, uscCoinExpected.String(),
				te.app.BankKeeper.GetBalance(te.ctx, accAddr, uscCoinExpected.Denom).String(),
			)

			// Verify Active pool balance
			assert.Equal(t,
				colUsedExpected.String(),
				te.app.USCKeeper.ActivePool(te.ctx).String(),
			)
		})
	}
}

func TestUSCKeeperMsgRedeemCollateral(t *testing.T) {
	type testCase struct {
		name            string
		uscCoinToRedeem string
		activePoolCoins string
		//
		errExpected              error
		uscCoinLeftExpected      string
		colCoinsRedeemedExpected string
	}

	testCases := []testCase{
		{
			name:                     "OK: native: partially filled",
			uscCoinToRedeem:          "10020uusc",      // 0.010020 USC
			activePoolCoins:          "5musdc,10uusdt", // 0.005 USDC, 0.000010 USDT
			uscCoinLeftExpected:      "5010uusc",       // 0.005010 USC
			colCoinsRedeemedExpected: "5musdc,10uusdt",
		},
		{
			name:                     "OK: native: fully filled",
			uscCoinToRedeem:          "130000000uusc",                             // 130.0 USC
			activePoolCoins:          "75000000000nbusd,50000000uusdt,25000musdc", // 75.0 BUSD, 50.0 USDT, 25.0 USDC
			uscCoinLeftExpected:      "0uusc",                                     // none
			colCoinsRedeemedExpected: "75000000000nbusd,50000000uusdt,5000musdc",  // 75.0 BUSD, 50.0 USDT, 5.0 USDC
		},
		{
			name:                     "OK: derivative: fully filled",
			uscCoinToRedeem:          "100000000" + USCDerivativeDenom1, // 100.0 IBC USC
			activePoolCoins:          "99500000uusdt,500musdc",          // 99.6 USDT, 0.5 USDC
			uscCoinLeftExpected:      "0" + USCDerivativeDenom1,         // none
			colCoinsRedeemedExpected: "99500000uusdt,500musdc",          // 99.6 USDT, 0.5 USDC
		},
		{
			name:            "Fail: unsupported USC derivative",
			uscCoinToRedeem: "10000ibc/7F1D3FCF4AE79E1554D670D1AD949A9BA4E4A3C76C63093E17E446A46061A7A2", // 0.010000 IBC USC (unsupported)
			activePoolCoins: "999nbusd",                                                                  // 0.000000999 BUSD
			errExpected:     types.ErrInvalidUSC,
		},
		{
			name:            "Fail: USC amount is too small",
			uscCoinToRedeem: "1uusc",    // 0.000001 USC
			activePoolCoins: "999nbusd", // 0.000000999 BUSD
			errExpected:     sdkErrors.ErrInsufficientFunds,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Fixtures
			te := NewTestEnv(t)

			accAddr, accCoins := te.AddAccount(t, tc.uscCoinToRedeem)
			uscRedeemCoin := accCoins[0]

			activePoolCoins := te.AddActivePoolBalance(t, tc.activePoolCoins)

			// Send msg
			msg := types.NewMsgRedeemCollateral(accAddr, uscRedeemCoin)
			require.NoError(t, msg.ValidateBasic())

			res, err := te.msgServer.RedeemCollateral(sdk.WrapSDKContext(te.ctx), msg)
			if tc.errExpected != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.errExpected)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, res)

			// Build expect value
			uscLeftExpectedCoin, err := sdk.ParseCoinNormalized(tc.uscCoinLeftExpected)
			require.NoError(t, err)

			uscBurnedExpected := uscRedeemCoin.Sub(uscLeftExpectedCoin)
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
				uscLeftExpectedCoin.String(),
				te.app.BankKeeper.GetBalance(te.ctx, accAddr, uscRedeemCoin.Denom).String(),
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
