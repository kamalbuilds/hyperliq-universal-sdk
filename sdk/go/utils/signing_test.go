package utils

import (
	"strings"
	"testing"
)

func TestGetAddressFromPrivateKey(t *testing.T) {
	// Test private key (DO NOT USE IN PRODUCTION)
	testPrivateKey := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	
	address, err := GetAddressFromPrivateKey(testPrivateKey)
	if err != nil {
		t.Fatalf("Failed to get address from private key: %v", err)
	}
	
	// Check that address is valid format
	if !strings.HasPrefix(address, "0x") {
		t.Error("Address should start with 0x")
	}
	
	if len(address) != 42 {
		t.Errorf("Address should be 42 characters, got %d", len(address))
	}
	
	// Test with invalid private key
	_, err = GetAddressFromPrivateKey("invalid_key")
	if err == nil {
		t.Error("Expected error for invalid private key")
	}
}

func TestRemoveHexPrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0x1234", "1234"},
		{"1234", "1234"},
		{"0X1234", "0X1234"}, // Only lowercase 0x is removed
		{"", ""},
		{"0x", ""},
	}
	
	for _, tt := range tests {
		result := removeHexPrefix(tt.input)
		if result != tt.expected {
			t.Errorf("removeHexPrefix(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestValidateAddress(t *testing.T) {
	tests := []struct {
		address string
		valid   bool
	}{
		{"0x1234567890123456789012345678901234567890", true},
		{"0x0000000000000000000000000000000000000000", true},
		{"1234567890123456789012345678901234567890", false}, // Missing 0x
		{"0x123", false}, // Too short
		{"0xGGGG567890123456789012345678901234567890", false}, // Invalid hex
		{"", false},
	}
	
	for _, tt := range tests {
		result := ValidateAddress(tt.address)
		if result != tt.valid {
			t.Errorf("ValidateAddress(%s) = %v, want %v", tt.address, result, tt.valid)
		}
	}
}

func TestNormalizeAddress(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"0x1234567890123456789012345678901234567890",
			"0x1234567890123456789012345678901234567890",
		},
		{
			"0xABCDEF1234567890123456789012345678901234",
			"0xAbCdEf1234567890123456789012345678901234", // EIP-55 checksum
		},
	}
	
	for _, tt := range tests {
		result := NormalizeAddress(tt.input)
		if strings.ToLower(result) != strings.ToLower(tt.expected) {
			t.Errorf("NormalizeAddress(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestHexConversion(t *testing.T) {
	testData := []byte("Hello, World!")
	
	// Convert to hex
	hexStr := BytesToHex(testData)
	if !strings.HasPrefix(hexStr, "0x") {
		t.Error("BytesToHex should add 0x prefix")
	}
	
	// Convert back to bytes
	bytes, err := HexToBytes(hexStr)
	if err != nil {
		t.Fatalf("HexToBytes failed: %v", err)
	}
	
	if string(bytes) != string(testData) {
		t.Errorf("Round trip conversion failed: got %s, want %s", string(bytes), string(testData))
	}
}

func TestHashMessage(t *testing.T) {
	message := []byte("test message")
	hash := HashMessage(message)
	
	if len(hash) != 32 {
		t.Errorf("Expected hash length 32, got %d", len(hash))
	}
	
	// Hash should be deterministic
	hash2 := HashMessage(message)
	if string(hash) != string(hash2) {
		t.Error("Hash should be deterministic")
	}
	
	// Different messages should have different hashes
	hash3 := HashMessage([]byte("different message"))
	if string(hash) == string(hash3) {
		t.Error("Different messages should have different hashes")
	}
}