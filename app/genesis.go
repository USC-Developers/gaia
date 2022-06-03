package gaia

import (
	"encoding/json"

	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() GenesisState {
	encCfg := MakeEncodingConfig()

	// Default state
	state := ModuleBasics.DefaultGenesis(encCfg.Marshaler)

	// Add x/bank default metadata (this is needed since x/usc can't start without those registered)
	bankState := bankTypes.DefaultGenesisState()
	bankState.DenomMetadata = []bankTypes.Metadata{
		{
			Base:    "ausc",
			Name:    "USC token",
			Symbol:  "USC",
			Display: "usc",
			DenomUnits: []*bankTypes.DenomUnit{
				{Denom: "ausc", Exponent: 0, Aliases: []string{"attousc"}},
				{Denom: "usc", Exponent: 18},
			},
		},
		{
			Base:    "uusdt",
			Name:    "USDT token",
			Symbol:  "USDT",
			Display: "usdt",
			DenomUnits: []*bankTypes.DenomUnit{
				{Denom: "musdt", Exponent: 0, Aliases: []string{"microusdt"}},
				{Denom: "usdt", Exponent: 6},
			},
		},
		{
			Base:    "uusdc",
			Name:    "USDC token",
			Symbol:  "USDC",
			Display: "usdC",
			DenomUnits: []*bankTypes.DenomUnit{
				{Denom: "musdc", Exponent: 0, Aliases: []string{"microusdc"}},
				{Denom: "usdc", Exponent: 6},
			},
		},
		{
			Base:    "abusd",
			Name:    "BUSD token",
			Symbol:  "BUSD",
			Display: "busd",
			DenomUnits: []*bankTypes.DenomUnit{
				{Denom: "abusd", Exponent: 0, Aliases: []string{"attobusd"}},
				{Denom: "busd", Exponent: 18},
			},
		},
	}
	state[bankTypes.ModuleName] = encCfg.Marshaler.MustMarshalJSON(bankState)

	return state
}
