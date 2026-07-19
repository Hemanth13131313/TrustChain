package main

import (
	"log"
	"net/http"
	"os"

	"github.com/trustchain/runtime-agent/internal/agent"
)

func main() {
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}
	
	publisher := agent.NewPublisher([]string{kafkaBrokers}, "trustchain.runtime.observations")
	defer publisher.Close()

	api := agent.NewAPI(publisher)

	http.HandleFunc("/simulate", api.HandleSimulate)

	log.Println("Runtime Agent (Mock) starting on :8084")
	if err := http.ListenAndServe(":8084", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
