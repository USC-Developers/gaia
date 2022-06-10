package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/v7/x/usc/types"
	"testing"
	"time"
)

var (
	ConversionCasesData = []conversionCase{
		{
			description:   "case for 10 busd 20 usdt 30 usdc",
			busdAmount:    10,
			usdtAmount:    20,
			usdcAmount:    30,
			mint:          true,
			expectedError: false,
		},
		{
			description:   "case for 30 busd 40 usdt 50 usdc",
			busdAmount:    30,
			usdtAmount:    40,
			usdcAmount:    50,
			mint:          true,
			expectedError: false,
		},
		{
			description:   "case for 100 each coins",
			busdAmount:    100,
			usdtAmount:    100,
			usdcAmount:    100,
			mint:          true,
			expectedError: false,
		},
		{
			description:   "case for conversion without minting",
			busdAmount:    100,
			usdtAmount:    50,
			usdcAmount:    10,
			mint:          false,
			expectedError: true,
		},
	}
)

type conversionCase struct {
	description   string
	busdAmount    int64
	usdtAmount    int64
	usdcAmount    int64
	expectedError bool
	mint          bool
}

func (s *TestSuite) TestConversion() {
	assert, require := s.Assert(), s.Require()
	uscKeeper, ctx := s.app.USCKeeper, s.ctx
	accAddr := s.accAddrs[0]
	for _, scenario := range ConversionCasesData {
		s.T().Run(scenario.description, func(t *testing.T) {
			GenBUSDCoin_ := sdk.NewCoin("abusd", sdk.NewInt(scenario.busdAmount))
			GenUSDTCoin_ := sdk.NewCoin("uusdt", sdk.NewInt(scenario.usdtAmount))
			GenUSDCCoin_ := sdk.NewCoin("musdc", sdk.NewInt(scenario.usdcAmount))
			collateralsExpected := sdk.NewCoins(GenBUSDCoin_, GenUSDTCoin_, GenUSDCCoin_)

			usc, err := uscKeeper.ConvertCollateralsToUSC(ctx, collateralsExpected)
			require.Nil(err)
			require.NotNil(usc)

			if scenario.mint {
				// Mint expected collateral
				msgMint := types.NewMsgMintUSC(accAddr, collateralsExpected)
				// Send Msg
				mintUSC, err := s.msgServer.MintUSC(sdk.WrapSDKContext(s.ctx), msgMint)
				require.Nil(err)
				require.NotNil(mintUSC)
			}

			// convert usc to collaterals
			collateralsResponse, err := uscKeeper.ConvertUSCToCollaterals(ctx, usc)
			if scenario.expectedError {
				require.NotNil(err)
				require.Nil(collateralsResponse)
				assert.Error(err)
				// collateral response should not be same as collateral expected
				assert.NotEqual(collateralsExpected, collateralsResponse)
			} else {
				require.Nil(err)
				require.NotNil(collateralsResponse)
				// collateral response should be same as collateral expected
				assert.Equal(collateralsExpected, collateralsResponse)

				// Redeem after tests
				msgRedeem := types.NewMsgRedeemCollateral(accAddr, sdk.NewCoin("ausc", sdk.NewInt(usc.Amount.Int64())))
				// Send Msg
				resRedeem, errRedeem := s.msgServer.RedeemCollateral(sdk.WrapSDKContext(s.ctx), msgRedeem)
				require.Nil(errRedeem)
				require.NotNil(resRedeem)

				// Redeeming pool should be equal to collateralsExpected
				assert.Equal(uscKeeper.RedeemingPool(ctx), collateralsExpected)
				uscKeeper.EndRedeeming(s.ctx.WithBlockTime(time.Now()))
			}
		})
	}
}
