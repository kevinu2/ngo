package BigMath

import (
	"github.com/shopspring/decimal"
	"ngo/constant"
)

func DecToBin(d1Str string) string {
	bin := decimal.NewFromInt(2)
	d1, _ := decimal.NewFromString(d1Str)
	decStr := ""
	if d1 == decimal.Zero {
		return constant.DefaultZero
	}

	// num /= 2 每次循环的时候 都将num除以2  再把结果赋值给 num
	for ; d1.GreaterThan(decimal.Zero); d1 = d1.Div(bin).Truncate(0) {
		//fmt.Println(num)
		lsb := d1.Mod(bin)
		decStr = lsb.String() + decStr
	}
	return decStr
}
