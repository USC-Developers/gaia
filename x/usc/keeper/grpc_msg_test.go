package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/v7/x/usc/types"
	"testing"
	"time"
)

var (
	validCoinsToSwap = []sdk.Coin{
		sdk.NewCoin("abusd", sdk.NewInt(10)),
		sdk.NewCoin("musdc", sdk.NewInt(10)),
		sdk.NewCoin("uusdt", sdk.NewInt(10)),
	}

	zeroCoinsToSwap = []sdk.Coin{
		sdk.NewCoin("abusd", sdk.NewInt(0)),
		sdk.NewCoin("musdc", sdk.NewInt(0)),
		sdk.NewCoin("uusdt", sdk.NewInt(0)),
	}

	outOfBalanceCoinsToSwap = []sdk.Coin{
		sdk.NewCoin("abusd", sdk.NewInt(1000)),
		sdk.NewCoin("musdc", sdk.NewInt(1000)),
		sdk.NewCoin("uusdt", sdk.NewInt(1000)),
	}

	inValidCoinsToSwap = []sdk.Coin{
		sdk.NewCoin("Invalid1", sdk.NewInt(10)),
		sdk.NewCoin("Invalid2", sdk.NewInt(10)),
		sdk.NewCoin("Invalid3", sdk.NewInt(10)),
	}

	validUscCoinAmount  = sdk.NewCoin("ausc", sdk.NewInt(10010000000000010))
	outOfBalanceUscCoin = sdk.NewCoin("ausc", sdk.NewInt(10010000000000011))
	randomCoin          = sdk.NewCoin("Random", sdk.NewInt(10))

	beforeMintBalances = []balanceCase{
		{description: "case for zero value of ausc",
			input: 0,
			denom: "ausc",
		}, {
			description: "case for 100 initial value of abusd",
			input:       100,
			denom:       "abusd",
		},
		{description: "case for 100 initial value of musdc",
			input: 100,
			denom: "musdc",
		},
		{description: "case for 100 initial value of uusdt",
			input: 100,
			denom: "uusdt",
		},
	}

	afterMintBalances = []balanceCase{
		{description: "case for equal value of usc after mint",
			input: 10010000000000010,
			denom: "ausc",
		}, {
			description: "case for value of busd after mint",
			input:       90,
			denom:       "abusd",
		},
		{description: "case for value of usdc after mint",
			input: 90,
			denom: "musdc",
		},
		{description: "case for value of usdt after mint",
			input: 90,
			denom: "uusdt",
		},
	}

	afterRedeemRequestBalances = []balanceCase{
		{description: "case for 0 usc after redeem request",
			input: 0,
			denom: "ausc",
		},
		{
			description: "case for value of busd after redeem request",
			input:       90,
			denom:       "abusd",
		},
		{description: "case for value of usdc after redeem request",
			input: 90,
			denom: "musdc",
		},
		{description: "case for value of usdt after redeem request",
			input: 90,
			denom: "uusdt",
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
			coin:          validUscCoinAmount,
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
	s.verifyUSCSupplyInvariant()
	s.verifyRedeemingQueueInvariant()
	// Check account balance before minting
	for _, scenario := range beforeMintBalances {
		s.T().Run(scenario.description, func(t *testing.T) {
			assert.EqualValues(scenario.input, bankKeeper.GetBalance(s.ctx, accAddr, scenario.denom).Amount.Uint64())
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

				assert.Equal(validUscCoinAmount.String(), resMint.MintedAmount.String())
				// Check pools
				var activePool []sdk.Coin = uscKeeper.ActivePool(s.ctx)
				assert.Equal(scenario.coins, activePool)
				assert.True(uscKeeper.RedeemingPool(s.ctx).IsZero())
			}
		})
	}

	s.verifyPool()
	s.verifyUSCSupplyInvariant()
	s.verifyRedeemingQueueInvariant()
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
				var redeemingPool []sdk.Coin = uscKeeper.RedeemingPool(s.ctx)
				assert.Equal(validCoinsToSwap, redeemingPool)

				s.verifyPool()
				s.verifyUSCSupplyInvariant()
				s.verifyRedeemingQueueInvariant()

				for _, scenario := range afterRedeemRequestBalances {
					s.T().Run(scenario.description, func(t *testing.T) {
						assert.EqualValues(scenario.input, bankKeeper.GetBalance(s.ctx, accAddr, scenario.denom).Amount.Int64())
					})

				}

				// clear redeeming queue
				uscKeeper.EndRedeeming(s.ctx.WithBlockTime(time.Now()))

				s.verifyPool()
				s.verifyUSCSupplyInvariant()
				s.verifyRedeemingQueueInvariant()
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
