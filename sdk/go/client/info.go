package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hyperliquid-labs/hyperliquid-go-sdk/types"
	"github.com/shopspring/decimal"
)

// InfoClient methods for market data and read-only operations

// GetUserState retrieves the user's account state
func (i *InfoClient) GetUserState(ctx context.Context, user string) (*types.UserState, error) {
	payload := map[string]interface{}{
		"type": "clearinghouseState",
		"user": user,
	}

	resp, err := i.client.request(ctx, "/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get user state: %w", err)
	}

	var userState types.UserState
	if err := json.Unmarshal(resp, &userState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user state: %w", err)
	}

	return &userState, nil
}

// GetOpenOrders retrieves user's open orders
func (i *InfoClient) GetOpenOrders(ctx context.Context, user string) ([]types.OpenOrder, error) {
	payload := map[string]interface{}{
		"type": "openOrders",
		"user": user,
	}

	resp, err := i.client.request(ctx, "/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get open orders: %w", err)
	}

	var orders []types.OpenOrder
	if err := json.Unmarshal(resp, &orders); err != nil {
		return nil, fmt.Errorf("failed to unmarshal open orders: %w", err)
	}

	return orders, nil
}

// GetAllMids retrieves mid prices for all assets
func (i *InfoClient) GetAllMids(ctx context.Context) (map[string]string, error) {
	payload := map[string]interface{}{
		"type": "allMids",
	}

	resp, err := i.client.request(ctx, "/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get all mids: %w", err)
	}

	var mids map[string]string
	if err := json.Unmarshal(resp, &mids); err != nil {
		return nil, fmt.Errorf("failed to unmarshal mids: %w", err)
	}

	return mids, nil
}

// GetL2Book retrieves the level 2 order book for an asset
func (i *InfoClient) GetL2Book(ctx context.Context, coin string) (*types.L2Book, error) {
	payload := map[string]interface{}{
		"type": "l2Book",
		"coin": coin,
	}

	resp, err := i.client.request(ctx, "/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get L2 book: %w", err)
	}

	var book types.L2Book
	if err := json.Unmarshal(resp, &book); err != nil {
		return nil, fmt.Errorf("failed to unmarshal L2 book: %w", err)
	}

	return &book, nil
}

// GetCandles retrieves candlestick data
func (i *InfoClient) GetCandles(ctx context.Context, coin, interval string, startTime, endTime int64) ([]types.Candle, error) {
	payload := map[string]interface{}{
		"type":      "candles",
		"req": map[string]interface{}{
			"coin":      coin,
			"interval":  interval,
			"startTime": startTime,
			"endTime":   endTime,
		},
	}

	resp, err := i.client.request(ctx, "/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get candles: %w", err)
	}

	var candles []types.Candle
	if err := json.Unmarshal(resp, &candles); err != nil {
		return nil, fmt.Errorf("failed to unmarshal candles: %w", err)
	}

	return candles, nil
}

// GetUserFills retrieves user's trade fills
func (i *InfoClient) GetUserFills(ctx context.Context, user string, startTime, endTime *int64) ([]types.Fill, error) {
	payload := map[string]interface{}{
		"type": "userFills",
		"user": user,
	}

	if startTime != nil {
		payload["startTime"] = *startTime
	}
	if endTime != nil {
		payload["endTime"] = *endTime
	}

	resp, err := i.client.request(ctx, "/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get user fills: %w", err)
	}

	var fills []types.Fill
	if err := json.Unmarshal(resp, &fills); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fills: %w", err)
	}

	return fills, nil
}

// GetUserFunding retrieves user's funding history
func (i *InfoClient) GetUserFunding(ctx context.Context, user string, startTime, endTime *int64) ([]types.FundingHistory, error) {
	payload := map[string]interface{}{
		"type": "userFunding",
		"user": user,
	}

	if startTime != nil {
		payload["startTime"] = *startTime
	}
	if endTime != nil {
		payload["endTime"] = *endTime
	}

	resp, err := i.client.request(ctx, "/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get user funding: %w", err)
	}

	var funding []types.FundingHistory
	if err := json.Unmarshal(resp, &funding); err != nil {
		return nil, fmt.Errorf("failed to unmarshal funding: %w", err)
	}

	return funding, nil
}

// GetFundingHistory retrieves funding rate history
func (i *InfoClient) GetFundingHistory(ctx context.Context, coin string, startTime, endTime *int64) ([]map[string]interface{}, error) {
	payload := map[string]interface{}{
		"type": "fundingHistory",
		"coin": coin,
	}

	if startTime != nil {
		payload["startTime"] = *startTime
	}
	if endTime != nil {
		payload["endTime"] = *endTime
	}

	resp, err := i.client.request(ctx, "/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get funding history: %w", err)
	}

	var history []map[string]interface{}
	if err := json.Unmarshal(resp, &history); err != nil {
		return nil, fmt.Errorf("failed to unmarshal funding history: %w", err)
	}

	return history, nil
}

// GetMeta retrieves market metadata
func (i *InfoClient) GetMeta(ctx context.Context) (*types.Meta, error) {
	payload := map[string]interface{}{
		"type": "meta",
	}

	resp, err := i.client.request(ctx, "/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get meta: %w", err)
	}

	var meta types.Meta
	if err := json.Unmarshal(resp, &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal meta: %w", err)
	}

	return &meta, nil
}

// GetSpotMeta retrieves spot market metadata
func (i *InfoClient) GetSpotMeta(ctx context.Context) (*types.SpotMeta, error) {
	payload := map[string]interface{}{
		"type": "spotMeta",
	}

	resp, err := i.client.request(ctx, "/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get spot meta: %w", err)
	}

	var spotMeta types.SpotMeta
	if err := json.Unmarshal(resp, &spotMeta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal spot meta: %w", err)
	}

	return &spotMeta, nil
}

// GetOrderStatus retrieves order status by order ID or client order ID
func (i *InfoClient) GetOrderStatus(ctx context.Context, user string, oid *int64, cloid *string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"type": "orderStatus",
		"user": user,
	}

	if oid != nil {
		payload["oid"] = *oid
	}
	if cloid != nil {
		payload["cloid"] = *cloid
	}

	resp, err := i.client.request(ctx, "/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get order status: %w", err)
	}

	var status map[string]interface{}
	if err := json.Unmarshal(resp, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order status: %w", err)
	}

	return status, nil
}

// GetLiquidations retrieves recent liquidations
func (i *InfoClient) GetLiquidations(ctx context.Context, startTime, endTime *int64) ([]map[string]interface{}, error) {
	payload := map[string]interface{}{
		"type": "liquidations",
	}

	if startTime != nil {
		payload["startTime"] = *startTime
	}
	if endTime != nil {
		payload["endTime"] = *endTime
	}

	resp, err := i.client.request(ctx, "/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get liquidations: %w", err)
	}

	var liquidations []map[string]interface{}
	if err := json.Unmarshal(resp, &liquidations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal liquidations: %w", err)
	}

	return liquidations, nil
}

// GetHistoricalOrders retrieves historical orders for a user
func (i *InfoClient) GetHistoricalOrders(ctx context.Context, user string, startTime, endTime *int64) ([]types.OpenOrder, error) {
	payload := map[string]interface{}{
		"type": "historicalOrders",
		"user": user,
	}

	if startTime != nil {
		payload["startTime"] = *startTime
	}
	if endTime != nil {
		payload["endTime"] = *endTime
	}

	resp, err := i.client.request(ctx, "/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical orders: %w", err)
	}

	var orders []types.OpenOrder
	if err := json.Unmarshal(resp, &orders); err != nil {
		return nil, fmt.Errorf("failed to unmarshal historical orders: %w", err)
	}

	return orders, nil
}

// GetTradeVolume retrieves recent trade volume
func (i *InfoClient) GetTradeVolume(ctx context.Context, window string) (map[string]decimal.Decimal, error) {
	payload := map[string]interface{}{
		"type": "vaultDetails",
		"req": map[string]interface{}{
			"type":   "tradeVolume",
			"window": window, // "1d", "7d", "30d"
		},
	}

	resp, err := i.client.request(ctx, "/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get trade volume: %w", err)
	}

	var volume map[string]decimal.Decimal
	if err := json.Unmarshal(resp, &volume); err != nil {
		return nil, fmt.Errorf("failed to unmarshal trade volume: %w", err)
	}

	return volume, nil
}