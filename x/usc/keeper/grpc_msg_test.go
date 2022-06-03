package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/v7/x/usc/types"
	"testing"
	"time"
)

var (
	validCoinsToSwap = []sdk.Coin{
		sdk.NewCoin("busd", sdk.NewInt(10)),
		sdk.NewCoin("usdc", sdk.NewInt(10)),
		sdk.NewCoin("usdt", sdk.NewInt(10)),
	}

	zeroCoinsToSwap = []sdk.Coin{
		sdk.NewCoin("busd", sdk.NewInt(0)),
		sdk.NewCoin("usdc", sdk.NewInt(0)),
		sdk.NewCoin("usdt", sdk.NewInt(0)),
	}

	outOfBalanceCoinsToSwap = []sdk.Coin{
		sdk.NewCoin("busd", sdk.NewInt(1000)),
		sdk.NewCoin("usdc", sdk.NewInt(1000)),
		sdk.NewCoin("usdt", sdk.NewInt(1000)),
	}

	inValidCoinsToSwap = []sdk.Coin{
		sdk.NewCoin("Invalid1", sdk.NewInt(10)),
		sdk.NewCoin("Invalid2", sdk.NewInt(10)),
		sdk.NewCoin("Invalid3", sdk.NewInt(10)),
	}

	uscCoin             = sdk.NewCoin("usc", sdk.NewInt(20000000000010))
	outOfBalanceUscCoin = sdk.NewCoin("usc", sdk.NewInt(20000000000011))
	randomCoin          = sdk.NewCoin("Random", sdk.NewInt(10))

	beforeMintBalances = []balanceCase{
		{description: "case for zero value of usc",
			input: 0,
			denom: "usc",
		}, {
			description: "case for 100 initial value of busd",
			input:       100,
			denom:       "busd",
		},
		{description: "case for 100 initial value of usdc",
			input: 100,
			denom: "usdc",
		},
		{description: "case for 100 initial value of usdt",
			input: 100,
			denom: "usdt",
		},
	}

	afterMintBalances = []balanceCase{
		{description: "case for equal value of usc after mint",
			input: 20000000000010,
			denom: "usc",
		}, {
			description: "case for value of busd after mint",
			input:       90,
			denom:       "busd",
		},
		{description: "case for value of usdc after mint",
			input: 90,
			denom: "usdc",
		},
		{description: "case for value of usdt after mint",
			input: 90,
			denom: "usdt",
		},
	}

	afterRedeemRequestBalances = []balanceCase{
		{description: "case for 0 usc after redeem request",
			input: 0,
			denom: "usc",
		},
		{
			description: "case for value of busd after redeem request",
			input:       90,
			denom:       "busd",
		},
		{description: "case for value of usdc after redeem request",
			input: 90,
			denom: "usdc",
		},
		{description: "case for value of usdt after redeem request",
			input: 90,
			denom: "usdt",
		},
	}

	mintMessageCoinsCasesData = []mintMessageCoinsCase{
		{
			description:   "case for valid coins for minting",
			coins:         validCoinsToSwap,
			expectedError: false,
		},

		{
			description:   "case for invalid coins for minting",
			coins:         inValidCoinsToSwap,
			expectedError: true,
		},

		{
			description:   "case for out of balance coins for minting",
			coins:         outOfBalanceCoinsToSwap,
			expectedError: true,
		},

		{
			description:   "case for zero coins for minting",
			coins:         zeroCoinsToSwap,
			expectedError: true,
		},
	}

	redeemMessageCoinsCasesData = []redeemMessageCoinCase{

		{
			description:   "case for valid redeem message",
			coin:          uscCoin,
			expectedError: false,
		},

		{
			description:   "case for out of balance redeem",
			coin:          outOfBalanceUscCoin,
			expectedError: true,
		}, {
			description:    "case for account with zero usc redeem",
			coin:           outOfBalanceUscCoin,
			expectedError:  true,
			uscZeroAddress: true,
		},

		{
			description:   "case for random coin redeem",
			coin:          randomCoin,
			expectedError: true,
		},
	}
)

type balanceCase struct {
	description string
	input       int64
	denom       string
}

type mintMessageCoinsCase struct {
	description   string
	coins         []sdk.Coin
	expectedError bool
}

type redeemMessageCoinCase struct {
	description    string
	coin           sdk.Coin
	expectedError  bool
	uscZeroAddress bool
}

func (s *TestSuite) TestMessages() {
	/* ----------- Setup --------------- */

	assert, require := s.Assert(), s.Require()
	bankKeeper, uscKeeper := s.app.BankKeeper, s.app.USCKeeper
	accAddr := s.accAddrs[0]
	accAddr1 := s.accAddrs[1]

	s.verifyPool()
	// Check account balance before minting
	for _, scenario := range beforeMintBalances {
		s.T().Run(scenario.description, func(t *testing.T) {
			assert.EqualValues(scenario.input, bankKeeper.GetBalance(s.ctx, accAddr, scenario.denom).Amount.Int64())
		})
	}

	/* ----------- Mint Message tests --------------- */
	for _, scenario := range mintMessageCoinsCasesData {
		s.T().Run(scenario.description, func(t *testing.T) {
			msgMint := types.NewMsgMintUSC(accAddr, scenario.coins)
			// Send Msg
			resMint, errMint := s.msgServer.MintUSC(sdk.WrapSDKContext(s.ctx), msgMint)
			if scenario.expectedError {
				require.Nil(resMint)
				require.Error(errMint)
			} else {
				require.NoError(errMint)
				require.NotNil(resMint)
				assert.Equal("20000000000010usc", resMint.MintedAmount.String())
				// Check pools
				assert.Equal("10busd,10usdc,10usdt", uscKeeper.ActivePool(s.ctx).String())
				assert.True(uscKeeper.RedeemingPool(s.ctx).IsZero())
			}
		})
	}

	s.verifyPool()
	// Check account balance after minting
	for _, scenario := range afterMintBalances {
		s.T().Run(scenario.description, func(t *testing.T) {
			assert.EqualValues(scenario.input, bankKeeper.GetBalance(s.ctx, accAddr, scenario.denom).Amount.Int64())
		})

	}

	/* ----------- Redeem Tests --------------- */
	for _, scenario := range redeemMessageCoinsCasesData {
		s.T().Run(scenario.description, func(t *testing.T) {
			accAddr_ := accAddr
			if scenario.uscZeroAddress {
				accAddr_ = accAddr1
			}
			msgRedeem := types.NewMsgRedeemCollateral(accAddr_, scenario.coin)
			// Send Msg
			resRedeem, errRedeem := s.msgServer.RedeemCollateral(sdk.WrapSDKContext(s.ctx), msgRedeem)
			if scenario.expectedError {
				require.Nil(resRedeem)
				require.Error(errRedeem)

			} else {
				require.NotNil(resRedeem)
				require.NoError(errRedeem)

				assert.Equal(validCoinsToSwap, resRedeem.RedeemedAmount)
				assert.True(uscKeeper.ActivePool(s.ctx).IsZero())
				assert.Equal("10busd,10usdc,10usdt", uscKeeper.RedeemingPool(s.ctx).String())

				s.verifyPool()

				for _, scenario := range afterRedeemRequestBalances {
					s.T().Run(scenario.description, func(t *testing.T) {
						assert.EqualValues(scenario.input, bankKeeper.GetBalance(s.ctx, accAddr, scenario.denom).Amount.Int64())
					})

				}

				// clear redeeming queue
				s.app.USCKeeper.EndRedeeming(s.ctx.WithBlockTime(time.Now()))

				s.verifyPool()
				// should be equal to balance before minting
				for _, scenario := range beforeMintBalances {
					s.T().Run(scenario.description, func(t *testing.T) {
						assert.EqualValues(scenario.input, bankKeeper.GetBalance(s.ctx, accAddr, scenario.denom).Amount.Int64())
					})

				}
			}
		})
	}
}
