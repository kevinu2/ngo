package Utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type stringsUtils struct {
}

func Strings() stringsUtils {
	return stringsUtils{}
}

func (stringsUtils) Reverse(vS string) string {
	bytes := []rune(vS)
	for from, to := 0, len(bytes)-1; from < to; from, to = from+1, to-1 {
		bytes[from], bytes[to] = bytes[to], bytes[from]
	}
	return string(bytes)
}

func (stringsUtils) IsEmpty(vS string) bool {
	return vS == ""
}

func (stringsUtils) ZeroAdd(vStr string, zeroLen int64) string {
	var (
		n       int64
		zeroStr = ""
		zero    = "0"
	)

	if zeroLen < 0 {
		return zeroStr
	}
	for n = 0; n < zeroLen; n++ {
		zeroStr = zeroStr + zero
	}
	return zeroStr + vStr
}

func (stringsUtils) RandCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var buffer strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&buffer, "%d", numeric[rand.Intn(r)])
	}
	return buffer.String()
}
