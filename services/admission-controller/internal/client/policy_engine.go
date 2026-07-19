package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type PolicyEngineClient struct {
	baseURL    string
	httpClient *http.Client
	failOpen   bool
}

func NewPolicyEngineClient(baseURL string, failOpen bool) *PolicyEngineClient {
	return &PolicyEngineClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		failOpen: failOpen,
	}
}

type PolicyRequest struct {
	Digest string `json:"digest"`
}

type PolicyResponse struct {
	Digest  string `json:"digest"`
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}

// Evaluate calls the policy engine to check if a specific digest is allowed
func (c *PolicyEngineClient) Evaluate(ctx context.Context, digest string) (bool, string, error) {
	reqBody, _ := json.Marshal(PolicyRequest{Digest: digest})
	
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/policy/evaluate", bytes.NewBuffer(reqBody))
	if err != nil {
		return c.failOpen, "failed to construct request", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return c.failOpen, fmt.Sprintf("policy engine unreachable (failOpen=%v)", c.failOpen), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.failOpen, fmt.Sprintf("policy engine returned status %d (failOpen=%v)", resp.StatusCode, c.failOpen), fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var polResp PolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&polResp); err != nil {
		return c.failOpen, "failed to decode policy engine response", err
	}

	return polResp.Allowed, polResp.Reason, nil
}
