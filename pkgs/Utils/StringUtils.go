package Utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
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
		_, _ = fmt.Fprintf(&buffer, "%d", numeric[rand.Intn(r)])
	}
	return buffer.String()
}

func (stringsUtils) IsDigit(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}

func (stringsUtils) IsLetter(str string) bool {
	pattern := "^[A-Za-z]+$"
	result, _ := regexp.MatchString(pattern, str)
	return result
}

func (stringsUtils) IsEmail(email string) bool {
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func (stringsUtils) Trim(str string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(str, "")
}
