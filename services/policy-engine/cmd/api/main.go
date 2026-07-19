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
	"github.com/trustchain/policy-engine/internal/api"
	"github.com/trustchain/policy-engine/internal/evaluator"
	"github.com/trustchain/policy-engine/internal/storage"
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

	// 2. Init Redis Cache
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}
	cache := storage.NewRedisCache(redisURL, "")

	// 3. Init OPA Evaluator
	policyPath := os.Getenv("OPA_POLICY_PATH")
	if policyPath == "" {
		policyPath = "../../policies/trustchain/rules.rego" // relative for local dev
	}
	opaEval, err := evaluator.NewOPAEvaluator(policyPath)
	if err != nil {
		log.Fatalf("Failed to initialize OPA: %v\n", err)
	}

	// 4. Init Repo & API Handlers
	dbRepo := storage.NewPostgresRepository(pool)
	h := api.NewPolicyHandler(dbRepo, cache, opaEval)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	h.RegisterRoutes(r)

	// 5. Start Server
	srv := &http.Server{
		Addr:    ":8082",
		Handler: r,
	}

	go func() {
		log.Println("Policy engine listening on :8082")
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
