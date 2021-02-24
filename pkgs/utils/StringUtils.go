package utils

type String struct {
}

//将字符串倒序
func (String) Reverse(vS string) string {
	bytes := []rune(vS)
	for from, to := 0, len(bytes)-1; from < to; from, to = from+1, to-1 {
		bytes[from], bytes[to] = bytes[to], bytes[from]
	}
	return string(bytes)
}

func (String) IsEmpty(vS string) bool {
	return vS == ""
}
