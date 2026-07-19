package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trustchain/verification/internal/api"
	"github.com/trustchain/verification/internal/event"
	"github.com/trustchain/verification/internal/storage"
	"github.com/trustchain/verification/internal/verifier"
)

func main() {
	// 1. Init DB connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://trustchain:trustchain_password@localhost:5432/trustchain?sslmode=disable"
	}
	
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	// 2. Init Repositories and Publisher
	dbRepo := storage.NewPostgresRepository(pool)
	
	kafkaBrokers := []string{"localhost:9092"}
	if envBrokers := os.Getenv("KAFKA_BROKERS"); envBrokers != "" {
		kafkaBrokers = []string{envBrokers}
	}
	publisher := event.NewKafkaPublisher(kafkaBrokers, "trustchain.verification.result")
	defer publisher.Close()

	// 3. Init Verifiers
	cosignVerifier := verifier.NewCosignVerifier()
	slsaVerifier := verifier.NewSLSAVerifier()

	// 4. Init API Handlers & Router
	h := api.NewVerificationHandler(dbRepo, publisher, cosignVerifier, slsaVerifier)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	h.RegisterRoutes(r)

	// 5. Start Server
	srv := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	go func() {
		log.Println("Verification service listening on :8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
