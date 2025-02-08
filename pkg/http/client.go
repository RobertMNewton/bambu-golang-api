package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/RobertMNewton/bambu-golang-api/pkg/types/config"
)

const (
	defaultTimeout = 30 * time.Second
	baseCloudURL   = "https://api.bambulab.com"
)

// Client represents an HTTP client for interacting with Bambu Lab's cloud API
type Client struct {
	config     config.PrinterConfig
	httpClient *http.Client
	token      string
}

// NewClient creates a new HTTP client
func NewClient(config config.PrinterConfig) (*Client, error) {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}, nil
}

// Login authenticates with the cloud API
func (c *Client) Login(ctx context.Context) error {
	loginURL := fmt.Sprintf("%s/v1/user-service/user/login", baseCloudURL)

	reqBody := map[string]string{
		"username": c.config.GetUsername(),
		"password": c.config.GetPassword(),
	}

	var resp struct {
		Token string `json:"token"`
	}

	err := c.doRequest(ctx, "POST", loginURL, reqBody, &resp)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	c.token = resp.Token
	return nil
}

func (c *Client) doRequest(ctx context.Context, method, url string, body interface{}, response interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
