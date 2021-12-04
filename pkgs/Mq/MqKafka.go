package Mq

import (
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/spf13/viper"
	"time"
)

var m *MsgQueue

type MsgQueue struct {
	topic   string
	service ConsumerI
	config  *KafkaConfig
}

func init() {
	m = new(MsgQueue)
}

func AddConsumer(topic string, service ConsumerI, mqGroup uint8) {
	m.AddConsumer(topic, service, mqGroup)
}
func (m *MsgQueue) AddConsumer(topic string, service ConsumerI, mqGroup uint8) {
	m.topic = topic
	m.service = service
	m.config = InitKafka(mqGroup)
}

func InitKafka(mqGroup uint8) *KafkaConfig {
	switch mqGroup {
	case MqGroup.GetCode():
		return &KafkaConfig{
			Host:    viper.GetStringSlice("kafka.host"),
			GroupId: viper.GetString("kafka.group_id"),
		}
	default:
		return &KafkaConfig{
			Host:    viper.GetStringSlice("kafka_default.host"),
			GroupId: viper.GetString("kafka_default.group_id"),
		}
	}
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

	groupID := fmt.Sprintf("%s%s", MqGroupPrefix, m.config.GroupId)
	// Init consumer, consume errors & messages
	consumer, err := cluster.NewConsumer(m.config.Host, groupID, []string{m.topic}, config)
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
			mqMsg.Topic = m.topic
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
