package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUSCParamsTest(t *testing.T) {
	type testCase struct {
		name       string
		buildInput func() Params
		//
		errExpected bool
	}

	validParams := Params{
		RedeemDur: 1 * time.Second,
		CollateralMetas: []TokenMeta{
			{Denom: "usdt", Decimals: 6},
			{Denom: "musdc", Decimals: 3},
		},
		UscIbcDenoms: []string{"ibc/7F1D3FCF4AE79E1554D670D1AD949A9BA4E4A3C76C63093E17E446A46061A7A2"},
	}

	testCases := []testCase{
		{
			name: "OK",
			buildInput: func() Params {
				return validParams
			},
		},
		{
			name: "Fail: RedeemDur: negative",
			buildInput: func() Params {
				p := validParams
				p.RedeemDur = -1 * time.Second
				return p
			},
			errExpected: true,
		},
		{
			name: "Fail: Meta: invalid denom",
			buildInput: func() Params {
				p := validParams
				p.CollateralMetas[0].Denom = InvalidDenom
				return p
			},
			errExpected: true,
		},
		{
			name: "Fail: Meta: zero decimals",
			buildInput: func() Params {
				p := validParams
				p.CollateralMetas[0].Decimals = 0
				return p
			},
			errExpected: true,
		},
		{
			name: "Fail: CollateralMetas: usc denom included",
			buildInput: func() Params {
				p := validParams
				p.CollateralMetas = append(p.CollateralMetas, USCMeta)
				return p
			},
			errExpected: true,
		},
		{
			name: "Fail: CollateralMetas: duplication",
			buildInput: func() Params {
				p := validParams
				p.CollateralMetas = append(p.CollateralMetas, p.CollateralMetas[0])
				return p
			},
			errExpected: true,
		},
		{
			name: "Fail: UscIbcDenoms: duplication",
			buildInput: func() Params {
				p := validParams
				p.UscIbcDenoms = append(p.UscIbcDenoms, "ibc/7F1D3FCF4AE79E1554D670D1AD949A9BA4E4A3C76C63093E17E446A46061A7A2")
				return p
			},
			errExpected: true,
		},
		{
			name: "Fail: UscIbcDenoms: non-IBC denom",
			buildInput: func() Params {
				p := validParams
				p.UscIbcDenoms = append(p.UscIbcDenoms, "uatom")
				return p
			},
			errExpected: true,
		},
		{
			name: "Fail: UscIbcDenoms: invalid hash",
			buildInput: func() Params {
				p := validParams
				p.UscIbcDenoms = append(p.UscIbcDenoms, "ibc/7F1D3FCF4AE79E1554D670D1AD949A9BA4E4A3C76C63093E17E446A46061A7A200")
				return p
			},
			errExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.errExpected {
				assert.Error(t, tc.buildInput().Validate())
				return
			}
			assert.NoError(t, tc.buildInput().Validate())
		})
	}
}
