package Kafka

import (
	"context"
	"crypto/tls"
	"errors"
	"os"
	"strings"
	"time"

	log "github.com/kevinu2/ngo/v2/pkgs/Log"
	"github.com/kevinu2/ngo/v2/pkgs/Utils"

	kafkaGo "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

// Consumer 使用的Kafka-go的input插件
type Consumer struct {
	config         map[any]any
	decorateEvents bool
	Messages       chan *Msg
	reader         *kafkaGo.Reader
	readerConfig   *kafkaGo.ReaderConfig
}

func (c *Consumer) ReaderConfig(config map[any]any) (*kafkaGo.ReaderConfig, error) {

	n := &kafkaGo.ReaderConfig{
		Brokers:        make([]string, 0),
		GroupID:        "",
		GroupTopics:    make([]string, 0),
		GroupBalancers: []kafkaGo.GroupBalancer{kafkaGo.RangeGroupBalancer{}},

		MinBytes: 1,
		MaxBytes: 10 * 1024 * 1024,
		MaxWait:  1000 * time.Millisecond,

		Dialer: &kafkaGo.Dialer{
			//Timeout:   10 * time.Second,
			//KeepAlive: 30 * time.Second,
			DualStack: true,
		},

		ReadBackoffMin: 500 * time.Millisecond,
		ReadBackoffMax: 5 * time.Second,

		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    30 * time.Second,
		RebalanceTimeout:  60 * time.Second,

		StartOffset:    kafkaGo.LastOffset,
		CommitInterval: 1 * time.Second,

		IsolationLevel: kafkaGo.ReadUncommitted,

		WatchPartitionChanges: true,
		ReadLagInterval:       10 * time.Second,

		Logger: kafkaGo.LoggerFunc(func(msg string, args ...any) {
			log.Logger().Debugf(msg, args...)
		}),
		ErrorLogger: kafkaGo.LoggerFunc(func(msg string, args ...any) {
			log.Logger().Errorf(msg, args...)
		}),
	}
	// check if GroupID is set in config
	if v, ok := config["GroupID"]; ok {
		n.GroupID = v.(string)
	} else {
		return nil, errors.New("GroupID must be set in kafka input")
	}
	// check if Brokers are set in config
	if v, ok := config["Brokers"]; ok {
		if vv, ok := v.([]any); ok {
			for _, broker := range vv {
				if brokerStr, ok := broker.(string); ok {
					n.Brokers = append(n.Brokers, brokerStr)
				} else {
					return nil, errors.New("brokers must be a list of strings in kafka output")
				}
			}
		} else if vv, ok := v.([]string); ok {
			n.Brokers = vv
		} else {
			return nil, errors.New("brokers must be a list of strings in kafka output")
		}
	} else {
		return nil, errors.New("brokers must be a list of strings in kafka output")
	}
	// check environment variable for Brokers
	if v, ok := os.LookupEnv("KAFKA_BROKERS"); ok {
		if Utils.IsValidBrokers(v, ",") {
			n.Brokers = strings.Split(v, ",")
		} else {
			log.Logger().Error("KAFKA_BROKERS environment variable is empty or invalid, should be in the format 'ip1:port,ip2:port'")
		}
	}
	if v, ok := config["Topics"]; ok {
		if vv, ok := v.([]any); ok {
			for _, topic := range vv {
				if topicStr, ok := topic.(string); ok {
					n.GroupTopics = append(n.GroupTopics, topicStr)
				} else {
					return nil, errors.New("topics must be a list of strings in kafka output")
				}
			}
		} else if vv, ok := v.([]string); ok {
			n.GroupTopics = vv
		} else {
			return nil, errors.New("topics must be a list of strings in kafka output")
		}
	} else {
		return nil, errors.New("topics must be a list of strings in kafka output")
	}
	if v, ok := config["MinBytes"]; ok {
		n.MinBytes = v.(int)
	}
	if v, ok := config["MaxBytes"]; ok {
		n.MaxBytes = v.(int)
	}
	if v, ok := config["HeartbeatInterval"]; ok {
		n.HeartbeatInterval = time.Duration(v.(int)) * time.Second
	}
	if v, ok := config["CommitInterval"]; ok {
		n.CommitInterval = time.Duration(v.(int)) * time.Second
	}
	if v, ok := config["MaxWait"]; ok {
		n.MaxWait = time.Duration(v.(int)) * time.Second
	}
	if v, ok := config["SessionTimeout"]; ok {
		n.SessionTimeout = time.Duration(v.(int)) * time.Second
	}
	if v, ok := config["RebalanceTimeout"]; ok {
		n.RebalanceTimeout = time.Duration(v.(int)) * time.Second
	}
	if v, ok := config["Timeout"].(int); ok {
		n.Dialer.Timeout = time.Duration(v) * time.Second
	}
	if v, ok := config["KeepAlive"].(int); ok {
		n.Dialer.KeepAlive = time.Duration(v) * time.Second
	}

	if v, ok := config["SASL"]; ok {
		vh := v.(map[string]any)
		v1, ok1 := vh["Type"]
		v2, ok2 := vh["Username"]
		v3, ok3 := vh["Password"]
		if ok1 && ok2 && ok3 {
			var (
				mechanism sasl.Mechanism
				err       error
			)
			switch v1.(string) {
			case "Plain":
				mechanism = plain.Mechanism{
					Username: v2.(string),
					Password: v3.(string),
				}
			case "SCRAM":
				mechanism, err = scram.Mechanism(
					scram.SHA512,
					v2.(string),
					v3.(string),
				)
				if err != nil {
					return nil, err
				}
			default:
				return nil, errors.New("unsupported SASL type: " + v1.(string))
			}
			n.Dialer.SASLMechanism = mechanism
		} else {
			return nil, errors.New("SASL configuration is incomplete, missing Type, Username or Password")
		}
	}
	//TODO TLS support
	if v, ok := config["TLS"]; ok {
		vh := v.(map[string]any)
		var t *tls.Config
		if _, ok := vh["PrivateKey"]; ok {

		}
		n.Dialer.TLS = t
	}
	if err := n.Validate(); err != nil {
		log.Logger().Errorf("ReadConfig Validate error: %s - %v", err.Error(), config)
		return nil, err
	}
	return n, nil
}

func NewConsumer(config map[any]any) (*Consumer, error) {
	c := &Consumer{
		Messages:       make(chan *Msg, 1024),
		decorateEvents: false,
		reader:         nil,
	}
	var err error

	c.readerConfig, err = c.ReaderConfig(config)
	if err != nil {
		log.Logger().Fatalf("consumer configuration error: %s", err.Error())
		return nil, err
	}

	c.reader = kafkaGo.NewReader(*c.readerConfig)

	go func() {
		failCount := 0
		for {
			m, err := c.reader.ReadMessage(context.Background())
			if err != nil {
				if strings.Contains(err.Error(), "Not Coordinator") {
					failCount++
					if failCount > 10 {
						log.Logger().Error("Kafka Not Coordinator Error 10 times, Exit")
						os.Exit(1)
					}
					log.Logger().Warn("Kafka Not Coordinator Error, Reconnect: %s", err.Error())
					_ = c.reader.Close()
					time.Sleep(2 * time.Second)
					if c.readerConfig, err = c.ReaderConfig(config); err == nil {
						c.reader = kafkaGo.NewReader(*c.readerConfig)
					} else {
						log.Logger().Error("consumer_settings wrong")
						os.Exit(1)
					}
					continue
				}
				log.Logger().Errorf("ReadMessage Error: %s", err.Error())
				break
			}
			c.Messages <- (*Msg)(&m)
		}
	}()
	return c, nil
}

func (c *Consumer) ReadOneEvent() []byte {
	message, ok := <-c.Messages
	if ok {
		return message.Value
	}
	return nil
}

func (c *Consumer) Shutdown() {
	if err := c.reader.Close(); err != nil {
		log.Logger().Errorf("failed to close reader: %s", err.Error())
	}
}
