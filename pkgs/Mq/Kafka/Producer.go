package Kafka

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/kevinu2/ngo/v2/pkgs/Log"
	"github.com/kevinu2/ngo/v2/pkgs/Utils"

	kafkaGo "github.com/segmentio/kafka-go"
	"github.com/ugorji/go/codec"
)

type Producer struct {
	config       map[interface{}]interface{}
	encoder      codec.Encoder
	producer     *kafkaGo.Writer
	Brokers      []string
	Topic        string
	mu           sync.Mutex
	recovering   int32
	closed       int32
	maxRetries   int
	retryBackoff time.Duration
}

func (o *Producer) newProducer() (*kafkaGo.Writer, error) {
	w := &kafkaGo.Writer{
		Addr:                   kafkaGo.TCP(o.Brokers...),
		Topic:                  o.Topic,
		Balancer:               &kafkaGo.LeastBytes{},
		MaxAttempts:            1,
		WriteBackoffMax:        o.retryBackoff,
		RequiredAcks:           kafkaGo.RequireAll,
		Async:                  true,
		AllowAutoTopicCreation: true,
	}
	return w, nil
}

func NewProducer(config map[interface{}]interface{}) (*Producer, error) {
	n := &Producer{
		config:       config,
		Brokers:      make([]string, 0),
		maxRetries:   10,
		retryBackoff: 1000 * time.Millisecond,
	}
	// Brokers
	if v, ok := config["Brokers"]; ok {
		if vv, ok := v.([]interface{}); ok {
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
		return nil, errors.New("brokers must be set in kafka output")
	}
	// check environment variable for Brokers
	if v, ok := os.LookupEnv("KAFKA_BROKERS"); ok {
		if Utils.IsValidBrokers(v, ",") {
			n.Brokers = strings.Split(v, ",")
		} else {
			log.Logger().Error("KAFKA_BROKERS environment variable is empty or invalid, should be in the format 'ip1:port,ip2:port'")
		}
	}
	// Topic
	if v, ok := config["Topic"].(string); ok {
		n.Topic = v
	} else {
		return nil, errors.New("topic must be set in kafka output")
	}
	// retries
	if v, ok := config["MaxRetries"].(int); ok {
		n.maxRetries = v
	}
	// backoff
	if v, ok := config["RetryBackoffMs"].(int); ok {
		n.retryBackoff = time.Duration(v) * time.Millisecond
	}
	p, err := n.newProducer()
	if err != nil {
		log.Logger().Error("Failed to create kafka producer %s, %s", n.Topic, err.Error())
		return nil, err
	}
	n.producer = p
	return n, nil
}

// helper to detect leader related errors
func (o *Producer) isNotLeaderErr(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "not leader for partition") ||
		strings.Contains(s, "leader not available") ||
		strings.Contains(s, "unknown Topic") ||
		strings.Contains(s, "no leader")
}

func (o *Producer) recreateWriter() {
	if atomic.LoadInt32(&o.closed) == 1 {
		return
	}
	if !atomic.CompareAndSwapInt32(&o.recovering, 0, 1) {
		return
	}
	defer atomic.StoreInt32(&o.recovering, 0)

	o.mu.Lock()
	defer o.mu.Unlock()
	if o.producer != nil {
		_ = o.producer.Close()
	}
	p, err := o.newProducer()
	if err != nil {
		log.Logger().Fatalf("Failed to recreate kafka producer for Topic=%s, %s", o.Topic, err.Error())
	}
	o.producer = p
	log.Logger().Infof("kafka writer recreated for Topic=%s", o.Topic)
}

func (o *Producer) Write(topic string, message []byte) error {
	var err error
	attempt := 0
	for {
		if atomic.LoadInt32(&o.closed) == 1 {
			return errors.New("kafka producer is closed")
		}
		o.mu.Lock()
		writer := o.producer
		o.mu.Unlock()
		if writer == nil {
			return errors.New("kafka writer is nil")
		}

		err = writer.WriteMessages(context.Background(), kafkaGo.Message{
			Topic: topic,
			Value: message,
		})
		if err == nil {
			return nil
		}

		if attempt == 0 && o.isNotLeaderErr(err) {
			log.Logger().Warnf("kafka not-leader error (will retry): %v", err)
			go o.recreateWriter()
		} else {
			log.Logger().Warnf("kafka write attempt %d failed: %v", attempt+1, err)
		}

		attempt++
		if attempt > o.maxRetries {
			log.Logger().Errorf("kafka write failed after %d retries: %v", attempt-1, err)
			return err
		}

		time.Sleep(o.retryBackoff * time.Duration(attempt))
	}
}

func (o *Producer) Shutdown() {
	atomic.StoreInt32(&o.closed, 1)
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.producer != nil {
		_ = o.producer.Close()
	}
}
