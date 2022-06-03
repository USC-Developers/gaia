package types_test

import (
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/v7/x/usc/types"
	"testing"
)

var (
	_, _, validAddress = testdata.KeyTestPubAddr()
	inValidAddress     = "InValidAddress"

	validateBasicCases = []validateBasicCase{
		{
			description:      "case for valid test",
			address:          validAddress.String(),
			collateralAmount: sdk.NewCoins(sdk.NewCoin("Normal", sdk.NewInt(100))),
			expectedError:    false,
		},
		{
			description:      "case for sending zero coins",
			address:          validAddress.String(),
			collateralAmount: sdk.NewCoins(sdk.NewCoin("Zero", sdk.NewInt(0))),
			expectedError:    true,
		},
		{
			description:      "case for sending no coins",
			address:          validAddress.String(),
			collateralAmount: sdk.NewCoins(),
			expectedError:    true,
		},
		{
			description:      "case for invalid address",
			address:          inValidAddress,
			collateralAmount: sdk.NewCoins(sdk.NewCoin("Normal", sdk.NewInt(100))),
			expectedError:    true,
		},
		{
			description:   "case for negative amount",
			address:       inValidAddress,
			denom:         "Normal",
			amount:        -100,
			expectedPanic: true,
		},
		{
			description:   "case for invalid Denom",
			address:       inValidAddress,
			denom:         "#Normal",
			amount:        100,
			expectedPanic: true,
		},
	}
)

type validateBasicCase struct {
	description      string
	address          string
	collateralAmount []sdk.Coin
	expectedError    bool
	expectedPanic    bool
	denom            string
	amount           int64
}

func (s *TestSuite) TestMsgMintValidateBasic() {
	// Amount 100
	assert, require := s.Assert(), s.Require()

	for _, scenario := range validateBasicCases {
		s.T().Run(scenario.description, func(t *testing.T) {
			if scenario.expectedPanic {
				assert.Panics(func() {
					types.MsgMintUSC{Address: scenario.address, CollateralAmount: sdk.NewCoins(sdk.NewCoin(scenario.denom, sdk.NewInt(scenario.amount)))}.ValidateBasic()
				}, "The code did not panic")
			} else if scenario.expectedError {
				require.Error(types.MsgMintUSC{
					Address:          scenario.address,
					CollateralAmount: scenario.collateralAmount,
				}.ValidateBasic())
			} else {
				require.NoError(types.MsgMintUSC{
					Address:          scenario.address,
					CollateralAmount: scenario.collateralAmount,
				}.ValidateBasic())
			}
		})
	}
}
