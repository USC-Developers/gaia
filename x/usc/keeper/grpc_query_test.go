package keeper_test

import (
	gocontext "context"
	"github.com/cosmos/gaia/v7/x/usc/types"
)

func (s *TestSuite) TestQueryPool() {
	s.verifyPool()
}

func (s *TestSuite) TestQueryParams() {
	app, ctx, queryClient := s.app, s.ctx, s.queryClient
	_, require := s.Assert(), s.Require()

	// Query Params
	resp, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	require.NoError(err)
	require.Equal(app.USCKeeper.GetParams(ctx), resp.Params)
}
