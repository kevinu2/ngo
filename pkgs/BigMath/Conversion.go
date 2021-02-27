package BigMath

import (
	"fmt"
	"github.com/shopspring/decimal"
	"ngo/constant"
)

func DecToBin(d1Str string) (decStr string) {
	bin := decimal.NewFromInt(2)
	d1, _ := decimal.NewFromString(d1Str)
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

func BinToDec(binStr string) (decStr string) {
	//bin, _ := decimal.NewFromString(binStr)
	var (
		bin     = "2"
		counter = fmt.Sprint(len(binStr) - 1)
	)

	for _, v := range []byte(binStr) {
		if string(v) == "1" {
			decStr = Add(decStr, Pow(bin, counter))
		}
		counter = Sub(counter, "1")
	}

	return decStr
}
