package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gaiaapp "github.com/cosmos/gaia/v7/app"
	"github.com/cosmos/gaia/v7/app/helpers"
	"github.com/cosmos/gaia/v7/x/usc/keeper"
	"github.com/cosmos/gaia/v7/x/usc/types"
	"github.com/stretchr/testify/suite"
	tmProto "github.com/tendermint/tendermint/proto/tendermint/types"
)

var (
	MockTimestamp = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

	BUSDMeta = types.TokenMeta{Denom: "abusd", Decimals: 18}
	USDTMeta = types.TokenMeta{Denom: "uusdt", Decimals: 6}
	USDCMeta = types.TokenMeta{Denom: "musdc", Decimals: 3}

	GenBUSDAmt, _ = sdk.NewIntFromString("10000000000000000000") // 10.0
	GenUSDTAmt, _ = sdk.NewIntFromString("10000000")             // 10.0
	GenUSDCAmt, _ = sdk.NewIntFromString("10000")                // 10.0

	GenCoins = sdk.NewCoins(
		sdk.NewCoin("abusd", GenBUSDAmt),
		sdk.NewCoin("uusdt", GenUSDTAmt),
		sdk.NewCoin("musdc", GenUSDCAmt),
	)
)

type TestSuite struct {
	suite.Suite

	app         *gaiaapp.GaiaApp
	ctx         sdk.Context
	queryClient types.QueryClient
	msgServer   types.MsgServer

	accAddrs []sdk.AccAddress
}

func (s *TestSuite) SetupTest() {
	app := helpers.Setup(s.T(), false, 1)
	ctx := app.BaseApp.NewContext(false, tmProto.Header{Time: MockTimestamp})

	uscParams := app.USCKeeper.GetParams(ctx)
	uscParams.CollateralMetas = []types.TokenMeta{BUSDMeta, USDTMeta, USDCMeta}
	app.USCKeeper.SetParams(ctx, uscParams)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQueryServerImpl(app.USCKeeper))
	queryClient := types.NewQueryClient(queryHelper)

	msgServer := keeper.NewMsgServerImpl(app.USCKeeper)

	genCoins := GenCoins
	genAddrs := helpers.AddTestAddrs(app, ctx, 2, genCoins)

	s.app, s.ctx, s.queryClient, s.msgServer, s.accAddrs = app, ctx, queryClient, msgServer, genAddrs
}

func TestUSCKeeper(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
