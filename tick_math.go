package uniswap_v3_simulator

import (
	"errors"
	"github.com/shopspring/decimal"
	"math"
	"math/big"
)

var POWERS_OF_2 []struct {
	i   int
	pow decimal.Decimal
}

func init() {
	for _, i := range []int{128, 64, 32, 16, 8, 4, 2, 1} {
		POWERS_OF_2 = append(POWERS_OF_2, struct {
			i   int
			pow decimal.Decimal
		}{
			i:   i,
			pow: decimal.NewFromInt(2).Pow(decimal.NewFromInt(int64(i))),
		})
	}

}
func TickSpacingToMaxLiquidityPerTick(tickSpacing int) decimal.Decimal {
	ts := decimal.NewFromInt(int64(tickSpacing))
	minTick := decimal.NewFromInt(int64(MIN_TICK)).Div(ts).RoundDown(0).Mul(ts)
	maxTick := decimal.NewFromInt(int64(MAX_TICK)).Div(ts).RoundDown(0).Mul(ts)
	numTicks := maxTick.Sub(minTick).Div(ts).RoundDown(0).Add(decimal.NewFromInt(1))
	return MaxUint128.Div(numTicks).RoundDown(0)
}

func GetTickAtSqrtRatio(sqrtRatioX96 decimal.Decimal) (int, error) {
	if sqrtRatioX96.LessThan(MIN_SQRT_RATIO) || sqrtRatioX96.GreaterThanOrEqual(MAX_SQRT_RATIO) {
		return 0, errors.New("SQRT_RATIO")
	}
	sqrtRatioX128 := sqrtRatioX96.Mul(decimal.NewFromInt(2).Pow(decimal.NewFromInt(32)))
	msb, err := MostSignificantBit(sqrtRatioX128)
	if err != nil {
		return 0, err
	}
	var r decimal.Decimal
	b := sqrtRatioX128.BigInt()
	if msb >= 128 {
		b = b.Rsh(b, uint(msb-127))
	} else {
		b = b.Lsh(b, uint(127-msb))
	}
	r = decimal.NewFromBigInt(b, 0)
	log_2 := big.NewInt(int64(msb - 128))
	log_2 = log_2.Lsh(log_2, 64)
	log2 := decimal.NewFromBigInt(log_2, 0)
	for i := 0; i < 14; i++ {
		tmp := r.Mul(r).BigInt()
		tmp.Rsh(tmp, 127)
		r = decimal.NewFromBigInt(tmp, 0)
		f := r.BigInt()
		f = f.Rsh(f, 128)
		log2_bigInt := log2.BigInt()
		log2_bigInt = log2_bigInt.Or(log2_bigInt, f.Lsh(f, uint(63-i)))
		r_bi := r.BigInt()
		r_bi = r_bi.Rsh(r_bi, uint(f.Int64()))
		r = decimal.NewFromBigInt(r_bi, 0)
	}
	c1, _ := decimal.NewFromString("255738958999603826347141")
	c2, _ := decimal.NewFromString("3402992956809132418596140100660247210")
	log_sqrt10001 := log2.Mul(c1)
	tickLow_bi := log_sqrt10001.Sub(c2).BigInt()
	tickLow_bi = tickLow_bi.Rsh(tickLow_bi, 128)
	tickLow := decimal.NewFromBigInt(tickLow_bi, 0)

	c3, _ := decimal.NewFromString("291339464771989622907027621153398088495")
	tickHigh_bi := log_sqrt10001.Add(c3).BigInt()
	tickHigh_bi = tickLow_bi.Rsh(tickHigh_bi, 128)
	tickHigh := decimal.NewFromBigInt(tickHigh_bi, 0)
	if tickLow.Equal(tickHigh) {
		return int(tickLow.IntPart()), nil
	} else {
		sqrt, err := GetSqrtRatioAtTick(int(tickHigh.IntPart()))
		if err != nil {
			return 0, err
		}
		if sqrt.LessThanOrEqual(sqrtRatioX96) {
			return int(tickHigh.IntPart()), nil
		} else {
			return int(tickLow.IntPart()), nil
		}
	}
}
func mulShift(val decimal.Decimal, mulBy string) decimal.Decimal {
	byBi, _ := big.NewInt(0).SetString(mulBy, 16)
	by := decimal.NewFromBigInt(byBi, 0)
	tmp := val.Mul(by).BigInt()
	tmp = tmp.Rsh(tmp, 128)
	return decimal.NewFromBigInt(tmp, 0)
}
func GetSqrtRatioAtTick(tick int) (decimal.Decimal, error) {
	if tick < MIN_TICK || tick > MAX_TICK {
		return decimal.Zero, errors.New("TICK")
	}
	var absTick int = int(math.Abs(float64(tick)))
	var ratio decimal.Decimal
	var ratioBi *big.Int
	if absTick&0x1 != 0 {
		ratioBi, _ = big.NewInt(0).SetString("fffcb933bd6fad37aa2d162d1a594001", 16)
	} else {
		ratioBi, _ = big.NewInt(0).SetString("100000000000000000000000000000000", 16)
	}
	ratio = decimal.NewFromBigInt(ratioBi, 0)

	if (absTick & 0x2) != 0 {
		ratio = mulShift(ratio, "fff97272373d413259a46990580e213a")
	}
	if (absTick & 0x4) != 0 {

		ratio = mulShift(ratio, "fff2e50f5f656932ef12357cf3c7fdcc")
	}
	if (absTick & 0x8) != 0 {

		ratio = mulShift(ratio, "ffe5caca7e10e4e61c3624eaa0941cd0")
	}
	if (absTick & 0x10) != 0 {

		ratio = mulShift(ratio, "ffcb9843d60f6159c9db58835c926644")
	}
	if (absTick & 0x20) != 0 {
		ratio = mulShift(ratio, "ff973b41fa98c081472e6896dfb254c0")

	}
	if (absTick & 0x40) != 0 {

		ratio = mulShift(ratio, "ff2ea16466c96a3843ec78b326b52861")
	}
	if (absTick & 0x80) != 0 {

		ratio = mulShift(ratio, "fe5dee046a99a2a811c461f1969c3053")
	}
	if (absTick & 0x100) != 0 {

		ratio = mulShift(ratio, "fcbe86c7900a88aedcffc83b479aa3a4")
	}
	if (absTick & 0x200) != 0 {

		ratio = mulShift(ratio, "f987a7253ac413176f2b074cf7815e54")
	}
	if (absTick & 0x400) != 0 {

		ratio = mulShift(ratio, "f3392b0822b70005940c7a398e4b70f3")
	}
	if (absTick & 0x800) != 0 {

		ratio = mulShift(ratio, "e7159475a2c29b7443b29c7fa6e889d9")
	}
	if (absTick & 0x1000) != 0 {

		ratio = mulShift(ratio, "d097f3bdfd2022b8845ad8f792aa5825")
	}
	if (absTick & 0x2000) != 0 {

		ratio = mulShift(ratio, "a9f746462d870fdf8a65dc1f90e061e5")
	}
	if (absTick & 0x4000) != 0 {

		ratio = mulShift(ratio, "70d869a156d2a1b890bb3df62baf32f7")
	}
	if (absTick & 0x8000) != 0 {

		ratio = mulShift(ratio, "31be135f97d08fd981231505542fcfa6")
	}
	if (absTick & 0x10000) != 0 {

		ratio = mulShift(ratio, "9aa508b5b7a84e1c677de54f3e99bc9")
	}
	if (absTick & 0x20000) != 0 {
		ratio = mulShift(ratio, "5d6af8dedb81196699c329225ee604")
	}
	if (absTick & 0x40000) != 0 {

		ratio = mulShift(ratio, "2216e584f5fa1ea926041bedfe98")
	}
	if (absTick & 0x80000) != 0 {

		ratio = mulShift(ratio, "48a170391f7dc42444e8fa2")
	}
	if tick > 0 {
		ratio = MaxUint256.Div(ratio).RoundDown(0)
	}
	_, remainder := ratio.QuoRem(Q32, 0)
	remainder = remainder.RoundDown(0)
	if remainder.GreaterThan(decimal.Zero) {
		return ratio.Div(Q32).Add(decimal.NewFromInt(1)).RoundDown(0), nil
	} else {
		return ratio.Div(Q32).RoundDown(0), nil
	}
}

func MostSignificantBit(x decimal.Decimal) (int, error) {
	if !x.GreaterThan(decimal.Zero) {
		return 0, errors.New("ZERO")
	}
	if !x.LessThanOrEqual(MaxUint256) {
		return 0, errors.New("MAX")
	}
	var msb int = 0
	for _, s := range POWERS_OF_2 {
		if x.GreaterThanOrEqual(s.pow) {
			x = x.Div(decimal.NewFromInt(2).Pow(decimal.NewFromInt(int64(s.i)))).RoundDown(0)
			msb += s.i
		}
	}
	return msb, nil
}
