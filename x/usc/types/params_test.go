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

	invalidDenom := "#Invalid"
	validParams := Params{
		RedeemDur: 1 * time.Second,
		CollateralMetas: []TokenMeta{
			{Denom: "usdt", Decimals: 6},
			{Denom: "musdc", Decimals: 3},
		},
		UscMeta: TokenMeta{Denom: "ausc", Decimals: 8},
	}

	testCases := []testCase{
		{
			name: "OK",
			buildInput: func() Params {
				return validParams
			},
		},
		{
			name: "Fail: UscMeta: invalid denom",
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
				p.UscMeta.Denom = invalidDenom
				return p
			},
			errExpected: true,
		},
		{
			name: "Fail: CollateralMetas: usc denom included",
			buildInput: func() Params {
				p := validParams
				p.CollateralMetas = append(p.CollateralMetas, p.UscMeta)
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
			name: "Fail: UscMeta: not the highest decimals value",
			buildInput: func() Params {
				p := validParams
				p.UscMeta.Decimals = 2
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
