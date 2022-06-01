package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	yaml "gopkg.in/yaml.v2"
)

const (
	DefaultParamSpace = ModuleName
)

// Params default values.
const (
	DefaultRedeemPeriod = 2 * 7 * 24 * time.Hour // 2 weeks
	DefaultUSCDenom     = "usc"
)

var (
	DefaultCollateralDenoms = []string{"usdt", "usdc", "busd"}
)

// Params storage keys.
var (
	ParamsKeyRedeemDur        = []byte("RedeemDur")
	ParamsKeyCollateralDenoms = []byte("CollateralDenoms")
	ParamsKeyUSCDenom         = []byte("USCDenom")
)

var _ paramsTypes.ParamSet = &Params{}

// NewParams creates a new Params object.
func NewParams(redeemDur time.Duration, collateralDenoms []string, uscDenom string) Params {
	return Params{
		RedeemDur:        redeemDur,
		CollateralDenoms: collateralDenoms,
		UscDenom:         uscDenom,
	}
}

// DefaultParams returns Params with defaults.
func DefaultParams() Params {
	return Params{
		RedeemDur:        DefaultRedeemPeriod,
		CollateralDenoms: DefaultCollateralDenoms,
		UscDenom:         DefaultUSCDenom,
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
		paramsTypes.NewParamSetPair(ParamsKeyCollateralDenoms, &p.CollateralDenoms, validateCollateralDenomsParam),
		paramsTypes.NewParamSetPair(ParamsKeyUSCDenom, &p.UscDenom, validateUSCDenomParam),
	}
}

// ValidateBasic perform a basic Params values validation.
func (p Params) ValidateBasic() error {
	if err := validateRedeemDurParam(p.RedeemDur); err != nil {
		return err
	}

	if err := validateCollateralDenomsParam(p.CollateralDenoms); err != nil {
		return err
	}

	if err := validateUSCDenomParam(p.UscDenom); err != nil {
		return err
	}

	for _, colDenom := range p.CollateralDenoms {
		if colDenom == p.UscDenom {
			return fmt.Errorf("usc_denom is used within collateral_denoms")
		}
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

// validateCollateralDenomsParam validates the CollateralDenoms param.
func validateCollateralDenomsParam(i interface{}) (retErr error) {
	defer func() {
		if retErr != nil {
			retErr = fmt.Errorf("collateral_denoms param: %w", retErr)
		}
	}()

	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type (%T, []string is expected)", i)
	}

	if len(v) == 0 {
		return fmt.Errorf("empty")
	}

	denomSet := make(map[string]struct{})
	for _, denom := range v {
		if err := sdk.ValidateDenom(denom); err != nil {
			return fmt.Errorf("%s: validation: %w", denom, err)
		}

		if _, ok := denomSet[denom]; ok {
			return fmt.Errorf("%s: duplicated", denom)
		}
		denomSet[denom] = struct{}{}
	}

	return
}

// validateUSCDenomParam validates the USCDenom param.
func validateUSCDenomParam(i interface{}) (retErr error) {
	defer func() {
		if retErr != nil {
			retErr = fmt.Errorf("usc_denom param: %w", retErr)
		}
	}()

	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type (%T, []string is expected)", i)
	}

	if err := sdk.ValidateDenom(v); err != nil {
		return fmt.Errorf("%s: validation: %w", v, err)
	}

	return
}
