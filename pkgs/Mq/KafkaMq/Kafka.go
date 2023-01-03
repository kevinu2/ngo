package KafkaMq

import (
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"time"
)

var m *MsgQueue

type MsgQueue struct {
	AsyncProducer sarama.AsyncProducer
	Service       ConsumerI
	Config        *Config
}

func init() {
	m = New()
}

func New() *MsgQueue { return m.New() }
func (m *MsgQueue) New() *MsgQueue {
	v := new(MsgQueue)
	return v
}

func AddConfig(topic map[string]string, topics, host []string, group string, debug bool) {
	m.AddConfig(topic, topics, host, group, debug)
}
func (m *MsgQueue) AddConfig(topic map[string]string, topics, host []string, group string, debug bool) {
	c := &Config{
		Topic:   topic,
		Topics:  topics,
		Host:    host,
		Group:   group,
		IsDebug: debug,
	}
	m.Config = c
}

func AddConsumer(service ConsumerI) { m.AddConsumer(service) }
func (m *MsgQueue) AddConsumer(service ConsumerI) {
	m.Service = service
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
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Init consumer, consume errors & messages
	consumer, err := cluster.NewConsumer(m.Config.Host, m.Config.Group, m.Config.Topics, config)
	if err != nil {
		fmt.Printf("Failed to start consumer: %s", err)
		return
	}
	defer consumer.Close()

	// Consume all channels, wait for signal to exit
	for {
		select {
		case msg, more := <-consumer.Messages():
			var mqMsg Msg
			mqMsg.Topic = msg.Topic
			mqMsg.Msg = string(msg.Value)
			m.Service.Consume(mqMsg)

			if more {
				if m.Config.IsDebug {
					fmt.Printf("%s/%d/%d\t%s\n", msg.Topic, string(msg.Value), msg.Partition, msg.Offset)
				}
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

func (m *MsgQueue) AddProducer() sarama.AsyncProducer {
	if m.AsyncProducer != nil {
		return m.AsyncProducer
	}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	ap, err := sarama.NewAsyncProducer(m.Config.Host, config)
	if err != nil {
		fmt.Printf("sarama.NewSyncProducer fails, err %s \n", err.Error())
		panic(err)
	}
	return ap

}

func (m *MsgQueue) Producer(message []byte, topic string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	m.AsyncProducer.Input() <- msg
	if <-m.AsyncProducer.Errors() != nil {
		if m.Config.IsDebug {
		}
		fmt.Printf("Send fails (%s), err %s \n", message, <-m.AsyncProducer.Errors())
		return <-m.AsyncProducer.Errors()
	} else {
		if m.Config.IsDebug {
			fmt.Printf("Send succeed(%s) \n", message)
		}
	}
	return nil
}

func (m *MsgQueue) CloseProducer() {
	m.AsyncProducer.AsyncClose()
}

func (m *MsgQueue) SyncProducer(message []byte, topic string) error {
	//指定配置
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewSyncProducer(m.Config.Host, config)
	if err != nil {
		if m.Config.IsDebug {
			fmt.Printf("sarama.NewSyncProducer (%v), err %s \n", m.Config.Host, err.Error())
		}
		return err
	}
	defer p.Close()
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	part, offset, err := p.SendMessage(msg)
	if err != nil {
		if m.Config.IsDebug {
			fmt.Printf("Send fails (%s/%s), err %s \n", topic, message, err.Error())
		}
		return err
	} else {
		if m.Config.IsDebug {
			fmt.Printf("Send succeed(%s/%s/%s/%s) \n", topic, message, part, offset)
		}
	}
	return nil
}
