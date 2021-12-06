package Mq

type MQueueMsg struct {
	Topic string
	Msg   string
}

type KafkaConfig struct {
	Host    []string `json:"host"`
	GroupId string   `json:"group_id"`
}

type ConsumerI interface {
	Consume(msg MQueueMsg)
}
