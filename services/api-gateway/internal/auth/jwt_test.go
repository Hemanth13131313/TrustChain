package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJWTAuth_Middleware(t *testing.T) {
	auth := NewJWTAuth("test-secret")
	token, _ := auth.GenerateToken("test-user")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{"Valid Token", "Bearer " + token, http.StatusOK},
		{"Missing Header", "", http.StatusUnauthorized},
		{"Invalid Format", "Token " + token, http.StatusUnauthorized},
		{"Invalid Token", "Bearer fake.token.here", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			auth.Middleware(nextHandler).ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
