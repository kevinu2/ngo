package KafkaMq

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
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

// ConsumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	service ConsumerI
	config  *Config
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var mqMsg Msg
		mqMsg.Topic = msg.Topic
		mqMsg.Msg = string(msg.Value)
		h.service.Consume(mqMsg)
		if h.config.IsDebug {
			fmt.Printf("%s/%d/%d/%s \n", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (m *MsgQueue) Consumer() {
	dc := DefaultConsumerConfig()
	config := sarama.NewConfig()
	config.Version, _ = sarama.ParseKafkaVersion(dc.Version)
	config.ClientID = dc.ClientID
	config.Metadata.Timeout = time.Duration(dc.MetadataMaxAgeMS) * time.Millisecond
	config.Consumer.Group.Session.Timeout = time.Duration(dc.SessionTimeoutMS) * time.Millisecond
	config.Consumer.MaxWaitTime = time.Duration(dc.FetchMaxWaitMS) * time.Millisecond
	config.Consumer.Fetch.Max = dc.FetchMaxBytes
	config.Consumer.Fetch.Min = dc.FetchMinBytes
	config.Consumer.Offsets.AutoCommit.Enable = dc.AutoCommit
	config.Consumer.Offsets.AutoCommit.Interval = time.Duration(dc.AutoCommitIntervalMS) * time.Millisecond
	config.Net.DialTimeout = time.Duration(dc.NetConfig.TimeoutMS) * time.Millisecond
	config.Net.KeepAlive = time.Duration(dc.NetConfig.KeepAliveMS) * time.Millisecond

	if dc.SaslConfig != nil {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = dc.SaslConfig.SaslUser
		config.Net.SASL.Password = dc.SaslConfig.SaslPassword
		config.Net.SASL.Mechanism = sarama.SASLMechanism(dc.SaslConfig.SaslMechanism)
	}

	if dc.FromBeginning {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	group, err := sarama.NewConsumerGroup(m.Config.Host, m.Config.Group, config)
	if err != nil {
		fmt.Printf("Failed to start consumer group: %s\n", err)
		return
	}
	defer group.Close()

	handler := &consumerGroupHandler{
		service: m.Service,
		config:  m.Config,
	}

	ctx := context.Background()
	for {
		err := group.Consume(ctx, m.Config.Topics, handler)
		if err != nil {
			fmt.Printf("Error from consumer: %s\n", err)
			break
		}
		if ctx.Err() != nil {
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
		fmt.Printf("sarama.NewAsyncProducer fails, err %s \n", err.Error())
		panic(err)
	}
	m.AsyncProducer = ap
	return ap
}

func (m *MsgQueue) Producer(message []byte, topic string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	m.AsyncProducer.Input() <- msg
	select {
	case err := <-m.AsyncProducer.Errors():
		if m.Config.IsDebug {
			fmt.Printf("Send fails (%s), err %s \n", message, err)
		}
		return err
	case <-m.AsyncProducer.Successes():
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
			fmt.Printf("Send succeed(%d/%d/%s/%s) \n", part, offset, topic, message)
		}
	}
	return nil
}
