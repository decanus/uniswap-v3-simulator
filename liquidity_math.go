package uniswap_v3_simulator

import (
	"errors"
	"github.com/shopspring/decimal"
)

var OVERFLOW = errors.New("OVERFLOW")
var UNDERFLOW = errors.New("UNDERFLOW")

func LiquidityAddDelta(x decimal.Decimal, y decimal.Decimal) (decimal.Decimal, error) {
	if x.GreaterThan(MaxUint128) {
		return decimal.Decimal{}, OVERFLOW
	}
	if y.GreaterThan(MaxUint128) {
		return decimal.Decimal{}, OVERFLOW
	}
	if y.IsNegative() {
		negy := y.Neg()
		if negy.GreaterThan(x) {
			return decimal.Decimal{}, UNDERFLOW
		}
		return x.Sub(negy), nil
	} else {
		if x.Add(y).GreaterThan(MaxUint128) {
			return decimal.Decimal{}, OVERFLOW
		}
		return x.Add(y), nil
	}
}
