package Utils

type Strings struck {
}

func (Strings) Reverse(vS string) string {
	bytes := []rune(vS)
	for from, to := 0, len(bytes)-1; from < to; from, to = from+1, to-1 {
		bytes[from], bytes[to] = bytes[to], bytes[from]
	}
	return string(bytes)
}

func (Strings) IsEmpty(vS string) bool {
	return vS == ""
}

func (Strings) ZeroAdd(vStr string, zeroLen int64) string {
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
