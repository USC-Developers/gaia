package types

import (
	"fmt"
	"strings"
	"time"

	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	ibcTransferTypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	yaml "gopkg.in/yaml.v2"
)

// Params default values.
const (
	DefaultRedeemPeriod     = 2 * 7 * 24 * time.Hour // 2 weeks
	DefaultMaxRedeemEntries = 7

	DefaultUSCDenom    = "uusc"
	DefaultUSCDecimals = 6
)

// Params storage keys.
var (
	ParamsKeyRedeemDur        = []byte("RedeemDur")
	ParamsKeyMaxRedeemEntries = []byte("MaxRedeemEntries")
	ParamsKeyCollateralMetas  = []byte("CollateralMetas")
	ParamsKeyUscIbcDenoms     = []byte("UscIbcDenoms")
)

// USCMeta is a hardcoded token meta for the USC native token.
var USCMeta = TokenMeta{
	Denom:    DefaultUSCDenom,
	Decimals: DefaultUSCDecimals,
}

var _ paramsTypes.ParamSet = &Params{}

// NewParams creates a new Params object.
func NewParams(redeemDur time.Duration, maxRedeemEntries uint32, collateralMetas []TokenMeta, uscIbcDenoms []string) Params {
	return Params{
		RedeemDur:        redeemDur,
		MaxRedeemEntries: maxRedeemEntries,
		CollateralMetas:  collateralMetas,
		UscIbcDenoms:     uscIbcDenoms,
	}
}

// DefaultParams returns Params with defaults.
func DefaultParams() Params {
	return Params{
		RedeemDur:        DefaultRedeemPeriod,
		MaxRedeemEntries: DefaultMaxRedeemEntries,
		CollateralMetas:  []TokenMeta{},
		UscIbcDenoms:     []string{},
	}
}

// ParamKeyTable returns module params table.
func ParamKeyTable() paramsTypes.KeyTable {
	return paramsTypes.NewKeyTable().RegisterParamSet(&Params{})
}

// String implements the fmt.Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)

	return string(out)
}

// ParamSetPairs implements the paramsTypes.ParamSet interface.
func (p *Params) ParamSetPairs() paramsTypes.ParamSetPairs {
	return paramsTypes.ParamSetPairs{
		paramsTypes.NewParamSetPair(ParamsKeyRedeemDur, &p.RedeemDur, validateRedeemDurParam),
		paramsTypes.NewParamSetPair(ParamsKeyMaxRedeemEntries, &p.MaxRedeemEntries, validateMaxRedeemEntriesParam),
		paramsTypes.NewParamSetPair(ParamsKeyCollateralMetas, &p.CollateralMetas, validateCollateralMetasParam),
		paramsTypes.NewParamSetPair(ParamsKeyUscIbcDenoms, &p.UscIbcDenoms, validateUscIbcDenomsParam),
	}
}

// Validate perform a Params object validation.
func (p Params) Validate() error {
	// Basic
	if err := validateRedeemDurParam(p.RedeemDur); err != nil {
		return err
	}

	if err := validateMaxRedeemEntriesParam(p.MaxRedeemEntries); err != nil {
		return err
	}

	if err := validateCollateralMetasParam(p.CollateralMetas); err != nil {
		return err
	}

	if err := validateUscIbcDenomsParam(p.UscIbcDenoms); err != nil {
		return err
	}

	return nil
}

// validateRedeemDurParam validates the RedeemDur param.
func validateRedeemDurParam(i interface{}) (retErr error) {
	defer func() {
		if retErr != nil {
			retErr = fmt.Errorf("redeem_dur param: %w", retErr)
		}
	}()

	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type (%T, time.Duration is expected)", i)
	}

	if v < 0 {
		return fmt.Errorf("must be GTE 0 (%d)", v)
	}

	return
}

// validateMaxRedeemEntriesParam validates the MaxRedeemEntries param.
func validateMaxRedeemEntriesParam(i interface{}) (retErr error) {
	defer func() {
		if retErr != nil {
			retErr = fmt.Errorf("max_redeem_entries param: %w", retErr)
		}
	}()

	if _, ok := i.(uint32); !ok {
		return fmt.Errorf("invalid parameter type (%T, uint32 is expected)", i)
	}

	return
}

// validateCollateralMetasParam validates the CollateralMetas param.
func validateCollateralMetasParam(i interface{}) (retErr error) {
	defer func() {
		if retErr != nil {
			retErr = fmt.Errorf("collateral_metas param: %w", retErr)
		}
	}()

	v, ok := i.([]TokenMeta)
	if !ok {
		return fmt.Errorf("invalid parameter type (%T, []TokenMeta is expected)", i)
	}

	denomSet := make(map[string]struct{})
	for _, meta := range v {
		if err := meta.Validate(); err != nil {
			return err
		}

		if meta.Denom == USCMeta.Denom {
			return fmt.Errorf("usc denom (%s) found", USCMeta.Denom)
		}

		if _, ok := denomSet[meta.Denom]; ok {
			return fmt.Errorf("tokenMeta (%s): duplicated", meta.Denom)
		}
		denomSet[meta.Denom] = struct{}{}
	}

	return
}

// validateUscIbcDenomsParam validates the UscIbcDenoms param.
func validateUscIbcDenomsParam(i interface{}) (retErr error) {
	defer func() {
		if retErr != nil {
			retErr = fmt.Errorf("usc_ibc_denoms param: %w", retErr)
		}
	}()

	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type (%T, []string is expected)", i)
	}

	denomSet := make(map[string]struct{})
	for _, denom := range v {
		if !strings.HasPrefix(denom, ibcTransferTypes.DenomPrefix) {
			return fmt.Errorf("denom (%s): not an IBC token (%s/ prefix )", denom, ibcTransferTypes.DenomPrefix)
		}

		if err := ibcTransferTypes.ValidateIBCDenom(denom); err != nil {
			return fmt.Errorf("denom (%s): validation: %w", denom, err)
		}

		if _, ok := denomSet[denom]; ok {
			return fmt.Errorf("denom (%s): duplicated", denom)
		}
		denomSet[denom] = struct{}{}
	}

	return
}
