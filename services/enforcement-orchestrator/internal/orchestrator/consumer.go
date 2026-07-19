package orchestrator

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type DriftEvent struct {
	Namespace   string `json:"namespace"`
	PodName     string `json:"pod_name"`
	ImageDigest string `json:"image_digest"`
	Reason      string `json:"reason"`
}

type Consumer struct {
	reader  *kafka.Reader
	auditor *Auditor
}

func NewConsumer(brokers []string, topic, groupID string, auditor *Auditor) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	return &Consumer{reader: r, auditor: auditor}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("Starting enforcement orchestrator consumer...")
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("error reading message: %v", err)
			continue
		}

		var event DriftEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("invalid drift event payload: %v", err)
			continue
		}

		log.Printf("Received drift event for pod %s/%s. Initiating quarantine...", event.Namespace, event.PodName)
		
		// 1. (Simulated) Quarantine Action
		log.Printf("[ACTION] Applied restrictive NetworkPolicy to isolate pod %s", event.PodName)

		// 2. (Simulated) Ticketing
		log.Printf("[TICKET] Created ServiceNow Incident INC-10293 for unauthorized image digest %s", event.ImageDigest)

		// 3. Cryptographic Audit Log
		payload := map[string]interface{}{
			"action":       "QUARANTINE",
			"target_pod":   event.PodName,
			"target_ns":    event.Namespace,
			"drift_reason": event.Reason,
			"ticket_ref":   "INC-10293",
		}
		
		if err := c.auditor.LogIncident(ctx, "DRIFT_QUARANTINED", payload); err != nil {
			log.Printf("CRITICAL: Failed to write audit log for pod %s: %v", event.PodName, err)
		} else {
			log.Printf("Audit log secured.")
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
