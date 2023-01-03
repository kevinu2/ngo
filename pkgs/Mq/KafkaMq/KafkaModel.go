package KafkaMq

type Msg struct {
	Topic string
	Msg   string
}

type Config struct {
	Topics  []string `json:"topics"`
	Topic   string   `json:"topic"`
	Host    []string `json:"host"`
	Group   string   `json:"group"`
	IsDebug bool     `json:"is_debug"`
}

type ConsumerI interface {
	Consume(msg Msg)
}
