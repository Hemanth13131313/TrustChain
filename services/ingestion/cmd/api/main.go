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
	"github.com/trustchain/ingestion/internal/api"
	"github.com/trustchain/ingestion/internal/storage"
)

func main() {
	// 1. Init DB connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://trustchain:trustchain_password@localhost:5432/trustchain?sslmode=disable" // fallback for local dev
	}
	
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	// 2. Init Repositories
	dbRepo := storage.NewPostgresRepository(pool)
	
	blobDir := os.Getenv("BLOB_STORAGE_DIR")
	if blobDir == "" {
		blobDir = "/tmp/trustchain-blobs"
	}
	blobRepo, err := storage.NewLocalBlobRepository(blobDir)
	if err != nil {
		log.Fatalf("Unable to init blob storage: %v\n", err)
	}

	// 3. Init API Handlers & Router
	h := api.NewIngestionHandler(dbRepo, blobRepo)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	h.RegisterRoutes(r)

	// 4. Start Server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Ingestion service listening on :8080")
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
