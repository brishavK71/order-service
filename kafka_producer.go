package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
	topic  string
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false,
	}
	return &KafkaProducer{writer: w, topic: topic}
}

func (kp *KafkaProducer) Close() error {
	return kp.writer.Close()
}

func (kp *KafkaProducer) PublishOrderCreated(ctx context.Context, evt OrderEvent) error {
	b, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(evt.OrderID),
		Value: b,
		Time:  time.Now(),
	}

	// WriteMessages retries internally on temporary failures
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return kp.writer.WriteMessages(ctx, msg)
}

// helper to build producer config from env (done in main)
func BrokerListFromEnv() []string {
	// default and environment override handled in main
	return []string{"kafka:9092"}
}
