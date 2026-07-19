package webhook

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/trustchain/admission-controller/internal/client"
)

func TestValidator_Validate(t *testing.T) {
	// Mock Policy Engine server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req client.PolicyRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := client.PolicyResponse{Digest: req.Digest}
		if req.Digest == "sha256:good" {
			resp.Allowed = true
			resp.Reason = "allowed"
		} else {
			resp.Allowed = false
			resp.Reason = "denied by policy"
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer mockServer.Close()

	peClient := client.NewPolicyEngineClient(mockServer.URL, false)
	validator := NewValidator(peClient)

	tests := []struct {
		name          string
		image         string
		expectAllowed bool
	}{
		{"Good Image", "registry.local/app@sha256:good", true},
		{"Bad Image", "registry.local/app@sha256:bad", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pod := corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Image: tt.image},
					},
				},
			}
			podBytes, _ := json.Marshal(pod)

			reviewReq := &admissionv1.AdmissionReview{
				Request: &admissionv1.AdmissionRequest{
					Kind: metav1.GroupVersionKind{Kind: "Pod"},
					Object: runtime.RawExtension{
						Raw: podBytes,
					},
				},
			}

			resp := validator.validate(context.Background(), reviewReq)

			if resp.Allowed != tt.expectAllowed {
				t.Errorf("expected allowed=%v, got %v", tt.expectAllowed, resp.Allowed)
			}
		})
	}
}
