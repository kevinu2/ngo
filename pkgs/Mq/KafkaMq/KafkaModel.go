package KafkaMq

type Msg struct {
	Topic string
	Msg   string
}

type Config struct {
	Topic []string `json:"topic"`
	Host  []string `json:"host"`
	Group string   `json:"group"`
}

type ConsumerI interface {
	Consume(msg Msg)
}
