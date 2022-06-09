package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"
)

// Validate performs a basic TokenMeta validation.
func (m TokenMeta) Validate() error {
	if err := sdk.ValidateDenom(m.Denom); err != nil {
		return fmt.Errorf("tokenMeta (%s): validation: %w", m.Denom, err)
	}

	return nil
}

// String implements the fmt.Stringer interface.
func (m TokenMeta) String() string {
	out, _ := yaml.Marshal(m)

	return string(out)
}

// DecUnit returns minimal token sdk.Dec value.
func (m TokenMeta) DecUnit() sdk.Dec {
	return sdk.NewDecWithPrec(1, int64(m.Decimals))
}

// NewZeroCoin returns a new empty sdk.Coin for the meta.
func (m TokenMeta) NewZeroCoin() sdk.Coin {
	return sdk.NewCoin(m.Denom, sdk.ZeroInt())
}

// ConvertCoin converts sdk.Coin to a given TokenMeta.
// Function is a variation of sdk.ConvertCoin.
func (m TokenMeta) ConvertCoin(coin sdk.Coin, dstMeta TokenMeta) (sdk.Coin, error) {
	if coin.Denom != m.Denom {
		return sdk.Coin{}, fmt.Errorf("coin.Denom (%s) is NE srcMeta.Denom (%s)", coin.Denom, m.Denom)
	}
	if err := m.Validate(); err != nil {
		return sdk.Coin{}, fmt.Errorf("invalid srcMeta: %w", err)
	}
	if err := dstMeta.Validate(); err != nil {
		return sdk.Coin{}, fmt.Errorf("invalid dstMeta: %w", err)
	}

	srcUnit, dstUnit := m.DecUnit(), dstMeta.DecUnit()
	if srcUnit.Equal(dstUnit) {
		return sdk.NewCoin(dstMeta.Denom, coin.Amount), nil
	}

	return sdk.NewCoin(dstMeta.Denom, coin.Amount.ToDec().Mul(srcUnit).Quo(dstUnit).TruncateInt()), nil
}

// NormalizeCoin converts sdk.Coin to a smaller decimals unit.
// Function is a variation of sdk.NormalizeCoin.
func (m TokenMeta) NormalizeCoin(coin sdk.Coin, dstMeta TokenMeta) (sdk.Coin, error) {
	if dstMeta.Decimals < m.Decimals {
		return sdk.Coin{}, fmt.Errorf("dstMeta.Decimals (%d) is LT srcMeta.Decimals (%s)", dstMeta.Decimals, m.Decimals)
	}

	coinNormalized, err := m.ConvertCoin(coin, dstMeta)
	if err != nil {
		return sdk.Coin{}, err
	}

	return sdk.NewCoin(m.Denom, coinNormalized.Amount), nil
}
