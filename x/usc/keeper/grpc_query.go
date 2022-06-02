package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/v7/x/usc/types"
)

var _ types.QueryServer = (*queryServer)(nil)

// queryServer implements the gRPC querier service.
type queryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the types.QueryServer interface.
func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

// Pool implements the types.QueryServer interface.
func (k queryServer) Pool(goCtx context.Context, req *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryPoolResponse{
		ActivePool:    k.Keeper.ActivePool(ctx),
		RedeemingPool: k.Keeper.RedeemingPool(ctx),
	}, nil
}

// Params implements the types.QueryServer interface.
func (k queryServer) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryParamsResponse{
		Params: k.GetParams(ctx),
	}, nil
}
