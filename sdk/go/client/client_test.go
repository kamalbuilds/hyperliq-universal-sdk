package client

import (
	"context"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		wsURL   string
		wantErr bool
	}{
		{
			name:    "Mainnet client",
			baseURL: MainnetAPI,
			wsURL:   MainnetWS,
			wantErr: false,
		},
		{
			name:    "Testnet client",
			baseURL: TestnetAPI,
			wsURL:   TestnetWS,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.baseURL, tt.wsURL, "test_key")
			
			if client == nil {
				t.Error("Expected client to be created")
			}
			
			if client.baseURL != tt.baseURL {
				t.Errorf("Expected baseURL %s, got %s", tt.baseURL, client.baseURL)
			}
			
			if client.wsURL != tt.wsURL {
				t.Errorf("Expected wsURL %s, got %s", tt.wsURL, client.wsURL)
			}
		})
	}
}

func TestRateLimiter(t *testing.T) {
	client := NewTestnetClient("test_key")
	
	// Test that rate limiter allows bursts
	ctx := context.Background()
	start := time.Now()
	
	// Should be able to make 20 requests quickly
	for i := 0; i < 20; i++ {
		err := client.rateLimiter.Wait(ctx)
		if err != nil {
			t.Errorf("Rate limiter failed on request %d: %v", i, err)
		}
	}
	
	elapsed := time.Since(start)
	if elapsed > time.Second {
		t.Errorf("Expected burst of 20 requests to complete in < 1s, took %v", elapsed)
	}
	
	// 21st request should be rate limited
	start = time.Now()
	err := client.rateLimiter.Wait(ctx)
	if err != nil {
		t.Errorf("Rate limiter failed on 21st request: %v", err)
	}
	elapsed = time.Since(start)
	
	// Should have waited ~50ms for the next token
	if elapsed < 40*time.Millisecond {
		t.Errorf("Expected rate limiting delay, but only waited %v", elapsed)
	}
}

func TestSetAddress(t *testing.T) {
	client := NewTestnetClient("test_key")
	
	testAddress := "0x1234567890123456789012345678901234567890"
	client.SetAddress(testAddress)
	
	if client.GetAddress() != testAddress {
		t.Errorf("Expected address %s, got %s", testAddress, client.GetAddress())
	}
}