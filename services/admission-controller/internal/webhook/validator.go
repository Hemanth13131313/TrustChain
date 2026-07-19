package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/trustchain/admission-controller/internal/client"
)

type Validator struct {
	policyClient *client.PolicyEngineClient
}

func NewValidator(policyClient *client.PolicyEngineClient) *Validator {
	return &Validator{policyClient: policyClient}
}

// HandleValidate responds to Kubernetes ValidatingWebhook requests
func (v *Validator) HandleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "reading request body", http.StatusBadRequest)
		return
	}

	var admissionReviewReq admissionv1.AdmissionReview
	if err := json.Unmarshal(body, &admissionReviewReq); err != nil {
		http.Error(w, "parsing admission review", http.StatusBadRequest)
		return
	}

	admissionResponse := v.validate(r.Context(), &admissionReviewReq)

	admissionReviewResp := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: admissionResponse,
	}
	if admissionReviewReq.Request != nil {
		admissionReviewResp.Response.UID = admissionReviewReq.Request.UID
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(admissionReviewResp)
}

func (v *Validator) validate(ctx context.Context, req *admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	if req.Request == nil {
		return &admissionv1.AdmissionResponse{Allowed: true} // Allow if unparseable or nil request
	}

	// We only care about Pods for this validator
	if req.Request.Kind.Kind != "Pod" {
		return &admissionv1.AdmissionResponse{Allowed: true}
	}

	var pod corev1.Pod
	if err := json.Unmarshal(req.Request.Object.Raw, &pod); err != nil {
		log.Printf("could not unmarshal pod: %v", err)
		return &admissionv1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	// Gather all unique images
	var images []string
	for _, container := range pod.Spec.Containers {
		images = append(images, container.Image)
	}
	for _, container := range pod.Spec.InitContainers {
		images = append(images, container.Image)
	}

	// Check each image against policy engine
	for _, img := range images {
		// MVP: Attempt to extract digest, fallback to raw string (which will fail in real engine if no digest)
		digest := extractDigest(img)
		
		allowed, reason, err := v.policyClient.Evaluate(ctx, digest)
		if err != nil {
			log.Printf("policy engine evaluation error for image %s: %v", img, err)
			// The client handles failOpen boolean injection into `allowed`
		}

		if !allowed {
			return &admissionv1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Message: fmt.Sprintf("Image %s denied by policy: %s", img, reason),
				},
			}
		}
	}

	return &admissionv1.AdmissionResponse{
		Allowed: true,
		Result: &metav1.Status{
			Message: "All images passed policy checks.",
		},
	}
}

// extractDigest is a placeholder to find the sha256 component.
// Real implementations should use "github.com/google/go-containerregistry/pkg/name"
func extractDigest(image string) string {
	parts := strings.Split(image, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	// No digest found, returning raw image string for MVP
	return image
}
