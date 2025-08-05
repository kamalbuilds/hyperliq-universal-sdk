package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// SignAction signs an action using EIP-712
func SignAction(action interface{}, privateKeyHex string, nonce int64) (string, error) {
	// Remove 0x prefix if present
	privateKeyHex = removeHexPrefix(privateKeyHex)
	
	// Parse private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create typed data
	typedData, err := createTypedData(action, nonce)
	if err != nil {
		return "", fmt.Errorf("failed to create typed data: %w", err)
	}

	// Sign the typed data
	signature, err := signTypedData(privateKey, typedData)
	if err != nil {
		return "", fmt.Errorf("failed to sign typed data: %w", err)
	}

	return signature, nil
}

// GetAddressFromPrivateKey derives the Ethereum address from a private key
func GetAddressFromPrivateKey(privateKeyHex string) (string, error) {
	// Remove 0x prefix if present
	privateKeyHex = removeHexPrefix(privateKeyHex)
	
	// Parse private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	// Get public key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("failed to cast public key to ECDSA")
	}

	// Get address
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address.Hex(), nil
}

// createTypedData creates EIP-712 typed data for signing
func createTypedData(action interface{}, nonce int64) (*apitypes.TypedData, error) {
	// Convert action to JSON to get deterministic ordering
	actionJSON, err := json.Marshal(action)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal action: %w", err)
	}

	var actionMap map[string]interface{}
	if err := json.Unmarshal(actionJSON, &actionMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal action: %w", err)
	}

	// Create domain
	domain := apitypes.TypedDataDomain{
		Name:              "Hyperliquid",
		Version:           "1",
		ChainId:           (*big.Int)(big.NewInt(42161)), // Arbitrum mainnet
		VerifyingContract: "0x0000000000000000000000000000000000000000",
	}

	// Define types based on action type
	types := apitypes.Types{
		"EIP712Domain": []apitypes.Type{
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
	}

	// Add types based on action type
	actionType, ok := actionMap["type"].(string)
	if !ok {
		return nil, fmt.Errorf("action type not found")
	}

	message := make(map[string]interface{})
	primaryType := ""

	switch actionType {
	case "order":
		types["Order"] = []apitypes.Type{
			{Name: "orders", Type: "OrderRequest[]"},
			{Name: "grouping", Type: "string"},
		}
		types["OrderRequest"] = []apitypes.Type{
			{Name: "coin", Type: "string"},
			{Name: "is_buy", Type: "bool"},
			{Name: "limit_px", Type: "string"},
			{Name: "sz", Type: "string"},
			{Name: "reduce_only", Type: "bool"},
			{Name: "order_type", Type: "string"},
		}
		primaryType = "Order"
		message = actionMap

	case "cancel":
		types["Cancel"] = []apitypes.Type{
			{Name: "cancels", Type: "CancelRequest[]"},
		}
		types["CancelRequest"] = []apitypes.Type{
			{Name: "coin", Type: "string"},
			{Name: "oid", Type: "uint64"},
		}
		primaryType = "Cancel"
		message = actionMap

	case "usdSend", "withdraw3":
		types["Transfer"] = []apitypes.Type{
			{Name: "destination", Type: "address"},
			{Name: "amount", Type: "string"},
			{Name: "time", Type: "uint64"},
		}
		primaryType = "Transfer"
		message = map[string]interface{}{
			"destination": actionMap["destination"],
			"amount":      actionMap["amount"],
			"time":        actionMap["time"],
		}

	default:
		// For other action types, create a generic structure
		types["Action"] = []apitypes.Type{
			{Name: "action", Type: "string"},
			{Name: "nonce", Type: "uint64"},
		}
		primaryType = "Action"
		message = map[string]interface{}{
			"action": normalizeAction(actionMap),
			"nonce":  nonce,
		}
	}

	typedData := &apitypes.TypedData{
		Domain:      domain,
		Types:       types,
		PrimaryType: primaryType,
		Message:     message,
	}

	return typedData, nil
}

// signTypedData signs EIP-712 typed data
func signTypedData(privateKey *ecdsa.PrivateKey, typedData *apitypes.TypedData) (string, error) {
	// Get the hash to sign
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return "", fmt.Errorf("failed to hash domain: %w", err)
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return "", fmt.Errorf("failed to hash message: %w", err)
	}

	// Create the final hash
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hash := crypto.Keccak256(rawData)

	// Sign the hash
	sig, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign hash: %w", err)
	}

	// Adjust V value for Ethereum compatibility
	sig[64] += 27

	return hexutil.Encode(sig), nil
}

// normalizeAction converts action map to deterministic string representation
func normalizeAction(action map[string]interface{}) string {
	// Sort keys for deterministic ordering
	keys := make([]string, 0, len(action))
	for k := range action {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build normalized string
	result := "{"
	for i, k := range keys {
		if i > 0 {
			result += ","
		}
		v := action[k]
		result += fmt.Sprintf(`"%s":`, k)
		
		switch val := v.(type) {
		case string:
			result += fmt.Sprintf(`"%s"`, val)
		case bool:
			result += fmt.Sprintf("%t", val)
		case float64:
			result += fmt.Sprintf("%g", val)
		case nil:
			result += "null"
		default:
			// For complex types, use JSON marshaling
			jsonVal, _ := json.Marshal(val)
			result += string(jsonVal)
		}
	}
	result += "}"
	
	return result
}

// removeHexPrefix removes 0x prefix from hex string if present
func removeHexPrefix(s string) string {
	if len(s) >= 2 && s[0:2] == "0x" {
		return s[2:]
	}
	return s
}

// ValidateAddress checks if an address is valid
func ValidateAddress(address string) bool {
	return common.IsHexAddress(address)
}

// NormalizeAddress normalizes an Ethereum address
func NormalizeAddress(address string) string {
	return common.HexToAddress(address).Hex()
}

// HashMessage creates a Keccak256 hash of a message
func HashMessage(message []byte) []byte {
	return crypto.Keccak256(message)
}

// HexToBytes converts a hex string to bytes
func HexToBytes(hexStr string) ([]byte, error) {
	hexStr = removeHexPrefix(hexStr)
	return hex.DecodeString(hexStr)
}

// BytesToHex converts bytes to hex string with 0x prefix
func BytesToHex(b []byte) string {
	return hexutil.Encode(b)
}