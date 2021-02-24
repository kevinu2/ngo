package utils

import (
	"github.com/shopspring/decimal"
	"math"
	"strconv"
)

// 加法
func Add(d1Str, d2Str string) string {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Add(d2).Truncate(8).String()
}

// 减法
func Sub(d1Str, d2Str string) string {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Sub(d2).Truncate(8).String()
}

//

func SubF(d1Str, d2Str string) float64 {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	d3, _ := d1.Sub(d2).Truncate(8).Float64()
	return d3
}

// 乘法
func Mul(d1Str, d2Str string) string {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Mul(d2).Truncate(8).String()
}
func AbsMul(d1Str, d2Str string) string {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Mul(d2).Truncate(8).Abs().String()
}

// 除法
func Div(d1Str, d2Str string) string {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	if d2.IsZero() {
		return d1Str
	}
	return d1.Div(d2).Truncate(8).String()
}

func DivInt(d1Str, d2Str string) int {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	if d2.IsZero() {
		return 0
	}
	rsStr := d1.Div(d2).Truncate(0).String()
	n, _ := strconv.ParseInt(rsStr, 0, 64)
	return int(n)
}

func Div4(d1Str, d2Str string) string {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	if d2.IsZero() {
		return d1Str
	}
	return d1.Div(d2).Truncate(4).String()
}

func DivF(d1Str, d2Str string) float64 {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	if d2.IsZero() {
		return 0.0
	}
	d3, _ := d1.Div(d2).Truncate(8).Float64()
	return d3
}

// 等于
func Eq(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Equal(d2)
}

// 绝对值等于
func AbsEq(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Abs().Equal(d2.Abs())
}

// 大于
func Gt(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.GreaterThan(d2)
}

// 绝对值大于
func AbsGt(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Abs().GreaterThan(d2.Abs())
}

// 大于等于
func Gte(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.GreaterThanOrEqual(d2)
}

// 绝对值大于等于
func AbsGte(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Abs().GreaterThanOrEqual(d2.Abs())
}

// 小于
func Lt(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.LessThan(d2)
}

// 绝对值小于
func AbsLt(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Abs().LessThan(d2.Abs())
}

// 小于等于
func Lte(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.LessThanOrEqual(d2)
}

// 绝对值小于等于
func AbsLte(d1Str, d2Str string) bool {
	d1, _ := decimal.NewFromString(d1Str)
	d2, _ := decimal.NewFromString(d2Str)
	return d1.Abs().LessThanOrEqual(d2.Abs())
}

// 反数
func Inverse(d1Str string) string {
	minus := decimal.NewFromInt(-1)
	d1, _ := decimal.NewFromString(d1Str)
	return d1.Mul(minus).Truncate(8).String()
}

// 绝对值
func Abs(d1Str string) string {
	d1, _ := decimal.NewFromString(d1Str)
	return d1.Abs().Truncate(8).String()
}

// n次方
func Pow(d1Str string, n string) string {
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
func Sqrt(d1Str string) string {
	d1, _ := decimal.NewFromString(d1Str)
	d1Float, _ := d1.Float64()
	d1Sqrt := math.Sqrt(d1Float)
	d1Result := decimal.NewFromFloat(d1Sqrt)
	//log.Logger().Infof("String: %v, Dec: %v, Float: %v, Sqrt: %v, result: %v", d1, d1Float, d1Sqrt, d1Result)
	return d1Result.String()
}
