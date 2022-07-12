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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmProto "github.com/tendermint/tendermint/proto/tendermint/types"
)

var (
	MockTimestamp = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

	BUSDMeta = types.TokenMeta{Denom: "nbusd", Decimals: 9}
	USDTMeta = types.TokenMeta{Denom: "uusdt", Decimals: 6}
	USDCMeta = types.TokenMeta{Denom: "musdc", Decimals: 3}

	USCDerivativeDenom1 = "ibc/000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F"
)

type TestEnv struct {
	app         *gaiaapp.GaiaApp
	ctx         sdk.Context
	queryClient types.QueryClient
	msgServer   types.MsgServer
}

func NewTestEnv(t *testing.T) *TestEnv {
	// Base app
	app := helpers.Setup(t, false, 1)
	ctx := app.BaseApp.NewContext(false, tmProto.Header{Time: MockTimestamp})

	// Add collateral metas
	uscParams := app.USCKeeper.GetParams(ctx)
	uscParams.CollateralMetas = []types.TokenMeta{BUSDMeta, USDTMeta, USDCMeta}
	uscParams.UscIbcDenoms = []string{USCDerivativeDenom1}
	app.USCKeeper.SetParams(ctx, uscParams)

	// Build services
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQueryServerImpl(app.USCKeeper))
	queryClient := types.NewQueryClient(queryHelper)

	msgServer := keeper.NewMsgServerImpl(app.USCKeeper)

	return &TestEnv{
		app:         app,
		ctx:         ctx,
		queryClient: queryClient,
		msgServer:   msgServer,
	}
}

func (te *TestEnv) AddAccount(t *testing.T, coinsRaw string) (sdk.AccAddress, sdk.Coins) {
	coins, err := sdk.ParseCoinsNormalized(coinsRaw)
	require.NoError(t, err)

	genAddrs := helpers.AddTestAddrs(te.app, te.ctx, 1, coins)

	return genAddrs[0], coins
}

func (te *TestEnv) AddActivePoolBalance(t *testing.T, coinsRaw string) sdk.Coins {
	coins, err := sdk.ParseCoinsNormalized(coinsRaw)
	require.NoError(t, err)

	require.NoError(t, te.app.BankKeeper.MintCoins(te.ctx, types.ActivePoolName, coins))

	return coins
}

func (te *TestEnv) AddRedeemingPoolBalance(t *testing.T, coinsRaw string) sdk.Coins {
	coins := te.AddActivePoolBalance(t, coinsRaw)
	require.NoError(t, te.app.BankKeeper.SendCoinsFromModuleToModule(te.ctx, types.ActivePoolName, types.RedeemingPoolName, coins))

	return coins
}

func (te *TestEnv) CheckInvariants(t *testing.T) {
	msg, broken := keeper.AllInvariants(te.app.USCKeeper)(te.ctx)
	assert.Falsef(t, broken, msg)
}
