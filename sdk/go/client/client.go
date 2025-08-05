package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

const (
	// API endpoints
	MainnetAPI = "https://api.hyperliquid.xyz"
	TestnetAPI = "https://api.hyperliquid-testnet.xyz"

	// WebSocket endpoints
	MainnetWS = "wss://api.hyperliquid.xyz/ws"
	TestnetWS = "wss://api.hyperliquid-testnet.xyz/ws"
)

// Client is the main Hyperliquid client
type Client struct {
	baseURL     string
	wsURL       string
	httpClient  *http.Client
	rateLimiter *rate.Limiter
	privateKey  string
	address     string
}

type InfoClient struct {
	client *Client
}

// NewClient creates a new Hyperliquid client
func NewClient(baseURL, wsURL, privateKey string) *Client {
	return &Client{
		baseURL:    baseURL,
		wsURL:      wsURL,
		privateKey: privateKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		// Rate limit: 1200 requests per minute (20 per second)
		rateLimiter: rate.NewLimiter(rate.Every(50*time.Millisecond), 20),
	}
}

// NewMainnetClient creates a client for mainnet
func NewMainnetClient(privateKey string) *Client {
	return NewClient(MainnetAPI, MainnetWS, privateKey)
}

// NewTestnetClient creates a client for testnet
func NewTestnetClient(privateKey string) *Client {
	return NewClient(TestnetAPI, TestnetWS, privateKey)
}

// SetAddress sets the client's address (derived from private key or provided)
func (c *Client) SetAddress(address string) {
	c.address = address
}

// GetAddress returns the client's address
func (c *Client) GetAddress() string {
	return c.address
}

// request performs an HTTP request with rate limiting
func (c *Client) request(ctx context.Context, endpoint string, payload interface{}) ([]byte, error) {
	// Apply rate limiting
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	// Marshal payload
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Info returns an InfoClient for market data queries
func (c *Client) Info() *InfoClient {
	return &InfoClient{client: c}
}

// Exchange returns an ExchangeClient for trading operations
func (c *Client) Exchange() *ExchangeClient {
	return &ExchangeClient{client: c}
}