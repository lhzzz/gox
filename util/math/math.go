package math

import (
	"math"
	// 该库内部做了优化，提前签名，使用的时候重新计算可用token数量
)

// FloatDecimal 浮点数保留小数点后prec位，去尾法
func FloatDecimal(f float64, prec int) float64 {
	precValue := math.Pow(10, float64(prec))
	return math.Trunc(f*precValue) / precValue
}

// FloatDecimalRound 浮点数保留小数点后prec位，四舍五入
func FloatDecimalRound(f float64, prec int) float64 {
	return FloatDecimal(f+5/math.Pow(10, float64(prec+1)), prec)
}

// FloatEqual 判断浮点数是否相等，prec控制精度
func FloatEqual(x, y, prec float64) bool {
	return math.Dim(x, y) < prec
}

// Float32Equal 判断float32是否相等，prec控制精度
func Float32Equal(x, y, prec float32) bool {
	return FloatEqual(float64(x), float64(y), float64(prec))
}
