package utils

import (
	"github.com/shopspring/decimal"
	"math"
	"strconv"
	"wanxiang/constant"
)

type BigMath struct {
}

// 加法
func (b BigMath) Add(d1Str, d2Str string) string {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Add(d2).Truncate(8).String()
}

// 减法
func (b BigMath) Sub(d1Str, d2Str string) string {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Sub(d2).Truncate(8).String()
}

//

func (b BigMath) SubF(d1Str, d2Str string) float64 {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	d3, _ := d1.Sub(d2).Truncate(8).Float64()
	return d3
}

// 乘法
func (b BigMath) Mul(d1Str, d2Str string) string {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Mul(d2).Truncate(8).String()
}
func (b BigMath) AbsMul(d1Str, d2Str string) string {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Mul(d2).Truncate(8).Abs().String()
}

// 除法
func (b BigMath) Div(d1Str, d2Str string) string {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	if d2.IsZero() {
		return d1Str
	}
	return d1.Div(d2).Truncate(8).String()
}

func (b BigMath) DivInt(d1Str, d2Str string) int {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	if d2.IsZero() {
		return 0
	}
	rsStr := d1.Div(d2).Truncate(0).String()
	n, _ := strconv.ParseInt(rsStr, 0, 64)
	return int(n)
}

func (b BigMath) Div4(d1Str, d2Str string) string {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	if d2.IsZero() {
		return d1Str
	}
	return d1.Div(d2).Truncate(4).String()
}

func (b BigMath) DivF(d1Str, d2Str string) float64 {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	if d2.IsZero() {
		return 0.0
	}
	d3, _ := d1.Div(d2).Truncate(8).Float64()
	return d3
}

// 等于
func (b BigMath) Eq(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Equal(d2)
}

// 绝对值等于
func (b BigMath) AbsEq(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Abs().Equal(d2.Abs())
}

// 大于
func (b BigMath) Gt(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.GreaterThan(d2)
}

// 绝对值大于
func (b BigMath) AbsGt(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Abs().GreaterThan(d2.Abs())
}

// 大于等于
func (b BigMath) Gte(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.GreaterThanOrEqual(d2)
}

// 绝对值大于等于
func (b BigMath) AbsGte(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Abs().GreaterThanOrEqual(d2.Abs())
}

// 小于
func (b BigMath) Lt(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.LessThan(d2)
}

// 绝对值小于
func (b BigMath) AbsLt(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Abs().LessThan(d2.Abs())
}

// 小于等于
func (b BigMath) Lte(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.LessThanOrEqual(d2)
}

// 绝对值小于等于
func (b BigMath) AbsLte(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Abs().LessThanOrEqual(d2.Abs())
}

// 反数
func (b BigMath) Inverse(d1Str string) string {
	minus := decimal.NewFromInt(-1)
	d1, _ := decimal.NewFromString(d1Str)
	return d1.Mul(minus).Truncate(8).String()
}

// 绝对值
func (b BigMath) Abs(d1Str string) string {
	d1, _ := decimal.NewFromString(d1Str)
	return d1.Abs().Truncate(8).String()
}

// n次方
func (b BigMath) Pow(d1Str string, n string) string {
	n1, _ := strconv.Atoi(n)
	d1, _ := decimal.NewFromString(d1Str)
	d2 := d1 //decimal.NewFromString("1")
	if n1 < 0 {
		return ""
	}
	if n1 == 1 {
		return d1.Truncate(8).String()
	}
	for i := 2; i <= n1; i++ {
		d1 = d1.Mul(d2)
		//fmt.Println(d1Str, i, d1.String())
	}
	return d1.Truncate(8).String()
}

// 平方根
func (b BigMath) Sqrt(d1Str string) string {
	d1, _ := decimal.NewFromString(d1Str)
	d1Float, _ := d1.Float64()
	d1Sqrt := math.Sqrt(d1Float)
	d1Result := decimal.NewFromFloat(d1Sqrt)
	//log.Logger().Infof("String: %v, Dec: %v, Float: %v, Sqrt: %v, result: %v", d1, d1Float, d1Sqrt, d1Result)
	return d1Result.String()
}

func (b BigMath) DecToBin(d1Str string) string {
	d1, _ := decimal.NewFromString(d1Str)
	decStr := ""
	if d1 == decimal.Zero {
		return "0"
	}

	// num /= 2 每次循环的时候 都将num除以2  再把结果赋值给 num
	for ; d1.GreaterThan(decimal.Zero); d1 = d1.Div(constant.Binary).Truncate(0) {
		//fmt.Println(num)
		lsb := d1.Mod(constant.Binary)
		decStr = lsb.String() + decStr
	}
	return decStr
}
