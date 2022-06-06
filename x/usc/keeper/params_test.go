package keeper_test

import (
	"github.com/cosmos/gaia/v7/x/usc/types"
)

func (s *TestSuite) TestParams() {
	require := s.Require()
	ctx, uscKeeper := s.ctx, s.app.USCKeeper

	expParams := types.DefaultParams()

	// check that the empty keeper loads the default
	resParams := uscKeeper.GetParams(ctx)
	require.True(expParams.Equal(resParams))

	// modify a params, save, and retrieve
	expParams.CollateralDenoms = []string{"test", "test1"}
	uscKeeper.SetParams(ctx, expParams)
	resParams = uscKeeper.GetParams(ctx)
	require.True(expParams.Equal(resParams))
}
