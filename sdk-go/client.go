package trustchain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

type ClientOptions struct {
	Timeout time.Duration
}

func NewClient(baseURL string, opts *ClientOptions) *Client {
	timeout := 10 * time.Second
	if opts != nil && opts.Timeout != 0 {
		timeout = opts.Timeout
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Authenticate fetches a JWT token from the API Gateway.
func (c *Client) Authenticate(ctx context.Context, username, password string) error {
	payload := map[string]string{
		"username": username,
		"password": password,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/v1/auth/login", c.baseURL), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed with status: %d", resp.StatusCode)
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	c.token = result.Token
	return nil
}

// GetArtifacts fetches the list of known artifacts.
func (c *Client) GetArtifacts(ctx context.Context) ([]map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/v1/artifacts", c.baseURL), nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch artifacts: %d", resp.StatusCode)
	}

	var artifacts []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&artifacts); err != nil {
		return nil, err
	}

	return artifacts, nil
}
