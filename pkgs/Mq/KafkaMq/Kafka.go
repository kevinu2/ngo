package KafkaMq

import (
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"time"
)

var m *MsgQueue

type MsgQueue struct {
	service ConsumerI
	config  *KafkaConfig
}

func init() {
	m = New()
}

func New() *MsgQueue { return m.New() }
func (m *MsgQueue) New() *MsgQueue {
	v := new(MsgQueue)
	return v
}

func AddConsumer(topic, host []string, group string, service ConsumerI) {
	m.AddConsumer(topic, host, group, service)
}
func (m *MsgQueue) AddConsumer(topic, host []string, group string, service ConsumerI) {
	m.service = service
	c := KafkaConfig{
		Topic: topic,
		Host:  host,
		Group: group,
	}
	m.config = &c
}

func (m *MsgQueue) ConsumeLoop() {
	for {
		m.Consumer()
		time.Sleep(3 * time.Second)
		fmt.Println("reconnect kafka publish...")
	}
}

func (m *MsgQueue) Consumer() {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Init consumer, consume errors & messages
	consumer, err := cluster.NewConsumer(m.config.Host, m.config.Group, m.config.Topic, config)
	if err != nil {
		fmt.Printf("Failed to start consumer: %s", err)
		return
	}
	defer consumer.Close()

	// Consume all channels, wait for signal to exit
	for {
		select {
		case msg, more := <-consumer.Messages():
			var mqMsg MQueueMsg
			mqMsg.Topic = msg.Topic
			mqMsg.Msg = string(msg.Value)
			m.service.Consume(mqMsg)

			if more {
				//fmt.Printf("%s/%d/%d\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
				consumer.MarkOffset(msg, "")
			}
		case ntf, more := <-consumer.Notifications():
			if more {
				fmt.Printf("Rebalanced: %+v\n", ntf)
			}
		case err, more := <-consumer.Errors():
			if more {
				fmt.Printf("Error: %s\n", err.Error())
			}
			break
		}
	}
}
