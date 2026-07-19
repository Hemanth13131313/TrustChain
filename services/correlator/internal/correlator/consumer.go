package correlator

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

type Consumer struct {
	reader   *kafka.Reader
	analyzer *Analyzer
}

func NewConsumer(brokers []string, topic, groupID string, analyzer *Analyzer) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	return &Consumer{reader: r, analyzer: analyzer}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("Starting observation consumer...")
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("error reading message: %v", err)
			continue
		}

		var obs WorkloadObservation
		if err := json.Unmarshal(m.Value, &obs); err != nil {
			log.Printf("invalid observation payload: %v", err)
			continue
		}

		// Analyze the observation
		_ = c.analyzer.CheckDrift(ctx, obs)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
