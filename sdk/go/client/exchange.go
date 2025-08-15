package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperliquid-labs/hyperliquid-go-sdk/types"
	"github.com/hyperliquid-labs/hyperliquid-go-sdk/utils"
)

// ExchangeClient methods for trading operations that require authentication

// PlaceOrder places a new order
func (e *ExchangeClient) PlaceOrder(ctx context.Context, order types.OrderRequest) (*types.OrderResponse, error) {
	action := map[string]interface{}{
		"type":       "order",
		"orders":     []types.OrderRequest{order},
		"grouping":   "na",
	}

	payload, err := e.createSignedRequest(action)
	if err != nil {
		return nil, fmt.Errorf("failed to create signed request: %w", err)
	}

	resp, err := e.client.request(ctx, "/exchange", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to place order: %w", err)
	}

	var orderResp types.OrderResponse
	if err := json.Unmarshal(resp, &orderResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order response: %w", err)
	}

	return &orderResp, nil
}

// PlaceOrders places multiple orders atomically
func (e *ExchangeClient) PlaceOrders(ctx context.Context, orders []types.OrderRequest, grouping string) (*types.OrderResponse, error) {
	if grouping == "" {
		grouping = "na"
	}

	action := map[string]interface{}{
		"type":     "order",
		"orders":   orders,
		"grouping": grouping,
	}

	payload, err := e.createSignedRequest(action)
	if err != nil {
		return nil, fmt.Errorf("failed to create signed request: %w", err)
	}

	resp, err := e.client.request(ctx, "/exchange", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to place orders: %w", err)
	}

	var orderResp types.OrderResponse
	if err := json.Unmarshal(resp, &orderResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order response: %w", err)
	}

	return &orderResp, nil
}

// CancelOrder cancels an order by ID or client order ID
func (e *ExchangeClient) CancelOrder(ctx context.Context, cancel types.CancelRequest) (*types.APIResponse, error) {
	action := map[string]interface{}{
		"type":    "cancel",
		"cancels": []types.CancelRequest{cancel},
	}

	payload, err := e.createSignedRequest(action)
	if err != nil {
		return nil, fmt.Errorf("failed to create signed request: %w", err)
	}

	resp, err := e.client.request(ctx, "/exchange", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}

	var apiResp types.APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cancel response: %w", err)
	}

	return &apiResp, nil
}

// CancelOrders cancels multiple orders
func (e *ExchangeClient) CancelOrders(ctx context.Context, cancels []types.CancelRequest) (*types.APIResponse, error) {
	action := map[string]interface{}{
		"type":    "cancel",
		"cancels": cancels,
	}

	payload, err := e.createSignedRequest(action)
	if err != nil {
		return nil, fmt.Errorf("failed to create signed request: %w", err)
	}

	resp, err := e.client.request(ctx, "/exchange", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel orders: %w", err)
	}

	var apiResp types.APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cancel response: %w", err)
	}

	return &apiResp, nil
}

// CancelAllOrders cancels all orders for an asset
func (e *ExchangeClient) CancelAllOrders(ctx context.Context, asset string) (*types.APIResponse, error) {
	action := map[string]interface{}{
		"type": "cancelByCloid",
		"cancels": []map[string]interface{}{
			{
				"asset": asset,
			},
		},
	}

	payload, err := e.createSignedRequest(action)
	if err != nil {
		return nil, fmt.Errorf("failed to create signed request: %w", err)
	}

	resp, err := e.client.request(ctx, "/exchange", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel all orders: %w", err)
	}

	var apiResp types.APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cancel response: %w", err)
	}

	return &apiResp, nil
}

// ModifyOrder modifies an existing order
func (e *ExchangeClient) ModifyOrder(ctx context.Context, modify types.ModifyRequest) (*types.APIResponse, error) {
	action := map[string]interface{}{
		"type":   "modify",
		"oid":    modify.Oid,
		"order":  modify,
	}

	payload, err := e.createSignedRequest(action)
	if err != nil {
		return nil, fmt.Errorf("failed to create signed request: %w", err)
	}

	resp, err := e.client.request(ctx, "/exchange", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to modify order: %w", err)
	}

	var apiResp types.APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal modify response: %w", err)
	}

	return &apiResp, nil
}

// UpdateLeverage updates leverage for an asset
func (e *ExchangeClient) UpdateLeverage(ctx context.Context, asset string, leverageMode string, leverage int) (*types.APIResponse, error) {
	action := map[string]interface{}{
		"type":     "updateLeverage",
		"asset":    asset,
		"isCross":  leverageMode == "cross",
		"leverage": leverage,
	}

	payload, err := e.createSignedRequest(action)
	if err != nil {
		return nil, fmt.Errorf("failed to create signed request: %w", err)
	}

	resp, err := e.client.request(ctx, "/exchange", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update leverage: %w", err)
	}

	var apiResp types.APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal leverage response: %w", err)
	}

	return &apiResp, nil
}

// UpdateIsolatedMargin adds or removes isolated margin
func (e *ExchangeClient) UpdateIsolatedMargin(ctx context.Context, asset string, amount float64) (*types.APIResponse, error) {
	action := map[string]interface{}{
		"type":     "updateIsolatedMargin",
		"asset":    asset,
		"isBuy":    amount > 0,
		"ntli":     amount,
	}

	payload, err := e.createSignedRequest(action)
	if err != nil {
		return nil, fmt.Errorf("failed to create signed request: %w", err)
	}

	resp, err := e.client.request(ctx, "/exchange", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update isolated margin: %w", err)
	}

	var apiResp types.APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal margin response: %w", err)
	}

	return &apiResp, nil
}

// Transfer performs a USDC transfer
func (e *ExchangeClient) Transfer(ctx context.Context, transfer types.TransferRequest) (*types.APIResponse, error) {
	action := map[string]interface{}{
		"type":        "usdSend",
		"signatureChainId": "0xa4b1",
		"hyperliquidChain": "Mainnet",
		"destination": transfer.Destination,
		"amount":      transfer.Amount.String(),
		"time":        time.Now().UnixMilli(),
	}

	payload, err := e.createSignedRequest(action)
	if err != nil {
		return nil, fmt.Errorf("failed to create signed request: %w", err)
	}

	resp, err := e.client.request(ctx, "/exchange", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to transfer: %w", err)
	}

	var apiResp types.APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transfer response: %w", err)
	}

	return &apiResp, nil
}

// Withdraw withdraws USDC to L1
func (e *ExchangeClient) Withdraw(ctx context.Context, withdraw types.WithdrawRequest) (*types.APIResponse, error) {
	action := map[string]interface{}{
		"type":        "withdraw3",
		"signatureChainId": "0xa4b1",
		"hyperliquidChain": "Mainnet",
		"destination": withdraw.Destination,
		"amount":      withdraw.Amount.String(),
		"time":        time.Now().UnixMilli(),
	}

	payload, err := e.createSignedRequest(action)
	if err != nil {
		return nil, fmt.Errorf("failed to create signed request: %w", err)
	}

	resp, err := e.client.request(ctx, "/exchange", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to withdraw: %w", err)
	}

	var apiResp types.APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal withdraw response: %w", err)
	}

	return &apiResp, nil
}

// SetReferrer sets a referral code
func (e *ExchangeClient) SetReferrer(ctx context.Context, code string) (*types.APIResponse, error) {
	action := map[string]interface{}{
		"type": "setReferrer",
		"code": code,
	}

	payload, err := e.createSignedRequest(action)
	if err != nil {
		return nil, fmt.Errorf("failed to create signed request: %w", err)
	}

	resp, err := e.client.request(ctx, "/exchange", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to set referrer: %w", err)
	}

	var apiResp types.APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal referrer response: %w", err)
	}

	return &apiResp, nil
}

// createSignedRequest creates a signed request payload
func (e *ExchangeClient) createSignedRequest(action interface{}) (map[string]interface{}, error) {
	if e.client.privateKey == "" {
		return nil, fmt.Errorf("private key not set")
	}

	nonce := time.Now().UnixMilli()
	
	// Create the signature
	signature, err := utils.SignAction(action, e.client.privateKey, nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to sign action: %w", err)
	}

	payload := map[string]interface{}{
		"action":    action,
		"nonce":     nonce,
		"signature": signature,
	}

	return payload, nil
}