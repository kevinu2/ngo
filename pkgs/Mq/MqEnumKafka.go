package Mq

type MqSrvGroup uint8

const (
	MqGroup MqSrvGroup = iota + 1
)

func (g MqSrvGroup) GetCode() uint8 {
	switch g {
	case MqGroup:
		return 1
	default:
		return 0
	}
}
