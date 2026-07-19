package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trustchain/export/internal/stix"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://trustchain:trustchain_password@localhost:5432/trustchain?sslmode=disable"
	}
	
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	mapper := stix.NewMapper(pool)

	http.HandleFunc("/api/v1/stix/incidents", func(w http.ResponseWriter, r *http.Request) {
		incidents, err := mapper.ExportIncidents(r.Context())
		if err != nil {
			http.Error(w, "failed to query incidents", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/stix+json;version=2.1")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"type": "bundle",
			"id":   "bundle--trustchain-export",
			"objects": incidents,
		})
	})

	log.Println("STIX/TAXII Export Service running on :8085")
	if err := http.ListenAndServe(":8085", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
