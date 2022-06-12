package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestMsgMintUSCValidateBasic(t *testing.T) {
	type testCase struct {
		name             string
		address          string
		collateralAmount []sdk.Coin
		//
		errExpected bool
	}

	validDenom, invalidDenom := "usdc", "#Invalid"

	_, _, validAddr := testdata.KeyTestPubAddr()
	invalidAddr := "InvalidAddress"

	testCases := []testCase{
		{
			name:             "OK",
			address:          validAddr.String(),
			collateralAmount: sdk.NewCoins(sdk.NewCoin(validDenom, sdk.OneInt())),
		},
		{
			name:             "Fail: coins with zero amt",
			address:          validAddr.String(),
			collateralAmount: sdk.NewCoins(sdk.NewCoin(validDenom, sdk.ZeroInt())),
			errExpected:      true,
		},
		{
			name:             "Fail: empty coins 1",
			address:          validAddr.String(),
			collateralAmount: sdk.NewCoins(),
			errExpected:      true,
		},
		{
			name:             "Fail: empty coins 2",
			address:          validAddr.String(),
			collateralAmount: nil,
			errExpected:      true,
		},
		{
			name:             "Fail: invalid address",
			address:          invalidAddr,
			collateralAmount: sdk.NewCoins(sdk.NewCoin(validDenom, sdk.OneInt())),
			errExpected:      true,
		},
		{
			name:             "Fail: negative amount",
			address:          invalidAddr,
			collateralAmount: []sdk.Coin{{Denom: validDenom, Amount: sdk.NewInt(-1)}},
			errExpected:      true,
		},
		{
			name:             "Fail: invalid Denom",
			address:          invalidAddr,
			collateralAmount: []sdk.Coin{{Denom: invalidDenom, Amount: sdk.OneInt()}},
			errExpected:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := MsgMintUSC{
				Address:          tc.address,
				CollateralAmount: tc.collateralAmount,
			}

			if tc.errExpected {
				assert.Error(t, msg.ValidateBasic())
				return
			}
			assert.NoError(t, msg.ValidateBasic())
		})
	}
}

func TestMsgRedeemCollateralValidateBasic(t *testing.T) {
	type testCase struct {
		name      string
		address   string
		usdAmount sdk.Coin
		//
		errExpected bool
	}

	validDenom, invalidDenom := "usdc", "#Invalid"

	_, _, validAddr := testdata.KeyTestPubAddr()
	invalidAddr := "InvalidAddress"

	testCases := []testCase{
		{
			name:      "OK",
			address:   validAddr.String(),
			usdAmount: sdk.NewCoin(validDenom, sdk.OneInt()),
		},
		{
			name:        "Fail: coin with zero amt",
			address:     validAddr.String(),
			usdAmount:   sdk.NewCoin(validDenom, sdk.ZeroInt()),
			errExpected: true,
		},
		{
			name:        "Fail: invalid address",
			address:     invalidAddr,
			usdAmount:   sdk.NewCoin(validDenom, sdk.OneInt()),
			errExpected: true,
		},
		{
			name:        "Fail: negative amount",
			address:     invalidAddr,
			usdAmount:   sdk.Coin{Denom: validDenom, Amount: sdk.NewInt(-1)},
			errExpected: true,
		},
		{
			name:        "Fail: invalid Denom",
			address:     invalidAddr,
			usdAmount:   sdk.Coin{Denom: invalidDenom, Amount: sdk.OneInt()},
			errExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := MsgRedeemCollateral{
				Address:   tc.address,
				UscAmount: tc.usdAmount,
			}

			if tc.errExpected {
				assert.Error(t, msg.ValidateBasic())
				return
			}
			assert.NoError(t, msg.ValidateBasic())
		})
	}
}
