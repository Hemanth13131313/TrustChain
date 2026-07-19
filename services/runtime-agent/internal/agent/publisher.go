package agent

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type WorkloadObservation struct {
	NodeName    string    `json:"node_name"`
	PodName     string    `json:"pod_name"`
	Namespace   string    `json:"namespace"`
	ImageDigest string    `json:"image_digest"`
	ObservedAt  time.Time `json:"observed_at"`
}

type Publisher struct {
	writer *kafka.Writer
}

func NewPublisher(brokers []string, topic string) *Publisher {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &Publisher{writer: w}
}

func (p *Publisher) Publish(ctx context.Context, obs WorkloadObservation) error {
	obs.ObservedAt = time.Now()
	val, err := json.Marshal(obs)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(obs.PodName),
			Value: val,
		},
	)
	
	if err != nil {
		log.Printf("failed to write observation to kafka: %v", err)
		return err
	}
	
	log.Printf("Published observation for %s (digest: %s)", obs.PodName, obs.ImageDigest)
	return nil
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}
