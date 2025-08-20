package BigMath

import (
	"fmt"
	"strings"

	"github.com/kevinu2/ngo/v2/pkgs/Default"
	"github.com/shopspring/decimal"
)

func DecToBin(d1Str string) (decStr string) {
	bin := decimal.NewFromInt(2)
	d1, _ := decimal.NewFromString(d1Str)
	if d1 == decimal.Zero {
		return Default.StringZero
	}
	for ; d1.GreaterThan(decimal.Zero); d1 = d1.Div(bin).Truncate(0) {
		lsb := d1.Mod(bin)
		decStr = lsb.String() + decStr
	}
	return decStr
}

func BinToDec(binStr string) (decStr string) {
	var (
		bin      = "2"
		binSlice = strings.Split(binStr, "")
		counter  = fmt.Sprint(len(binSlice) - 1)
	)

	for _, v := range binSlice {
		if v == "1" {
			decStr = Add(decStr, Pow(bin, counter))
		}
		counter = Sub(counter, "1")
	}
	return decStr
}
