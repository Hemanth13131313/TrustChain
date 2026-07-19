package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/trustchain/admission-controller/internal/client"
	"github.com/trustchain/admission-controller/internal/webhook"
)

func main() {
	peURL := os.Getenv("POLICY_ENGINE_URL")
	if peURL == "" {
		peURL = "http://localhost:8082" // Local fallback
	}
	
	failOpenStr := strings.ToLower(os.Getenv("FAILURE_POLICY"))
	failOpen := false
	if failOpenStr == "fail-open" || failOpenStr == "ignore" {
		failOpen = true
	}

	log.Printf("Starting admission-controller. failOpen=%v, policyEngine=%s", failOpen, peURL)

	peClient := client.NewPolicyEngineClient(peURL, failOpen)
	validator := webhook.NewValidator(peClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/validate", validator.HandleValidate)

	tlsCertPath := os.Getenv("TLS_CERT_PATH")
	tlsKeyPath := os.Getenv("TLS_KEY_PATH")

	addr := ":8443"
	
	if tlsCertPath == "" || tlsKeyPath == "" {
		log.Println("WARNING: TLS_CERT_PATH or TLS_KEY_PATH not set. Kubernetes requires HTTPS for admission webhooks.")
		log.Println("Starting HTTP server on :8080 for testing purposes only...")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	} else {
		log.Printf("Starting HTTPS server on %s", addr)
		if err := http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, mux); err != nil {
			log.Fatalf("HTTPS server failed: %v", err)
		}
	}
}
