package Http

const (
	SessionId   = "session_id"
	HttpCode200 = 200
)

var (
	CodeUnLogin = 2013
)

const (
	NoErr         uint16 = 0
	SystemCodeErr uint16 = 5000
)

var (
	NoError  = Error{Code: 0, Msg: "success", Result: ""}
	NoResult = ""
)
