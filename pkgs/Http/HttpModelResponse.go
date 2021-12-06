package Http

type Response struct {
	Err    Error       `json:"error"`
	Result interface{} `json:"result"`
}

type RightResponse struct {
	Result interface{} `json:"result"`
}

type Error struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
}

type PageResponse struct {
	Items []interface{} `json:"items"`
	Count int64         `json:"count"`
}
