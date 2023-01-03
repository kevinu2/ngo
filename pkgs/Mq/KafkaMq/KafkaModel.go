package KafkaMq

type Msg struct {
	Topic string
	Msg   string
}

type Config struct {
	Topic      []string `json:"topic"`
	Host       []string `json:"host"`
	Group      string   `json:"group"`
	AutoCommit bool     `json:"auto_commit"`
}

type ConsumerI interface {
	Consume(msg Msg)
}
