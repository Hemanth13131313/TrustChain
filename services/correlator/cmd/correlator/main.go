package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trustchain/correlator/internal/correlator"
)

func main() {
	// Initialize DB connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://trustchain:trustchain_password@localhost:5432/trustchain?sslmode=disable"
	}
	
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	analyzer := correlator.NewAnalyzer(pool)

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}

	consumer := correlator.NewConsumer(
		[]string{kafkaBrokers},
		"trustchain.runtime.observations",
		"trustchain-drift-correlator-group",
		analyzer,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go consumer.Start(ctx)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down correlator...")
	cancel()
	time.Sleep(1 * time.Second) // Give consumer time to exit
	consumer.Close()
	log.Println("Correlator exiting")
}
