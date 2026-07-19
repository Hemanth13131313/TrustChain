package event

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/trustchain/verification/internal/domain"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string, topic string) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaPublisher) PublishVerificationResult(ctx context.Context, result domain.VerificationResult) error {
	payload, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshaling verification result: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(result.Digest), // use digest as partition key
		Value: payload,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("writing kafka message: %w", err)
	}
	return nil
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}
