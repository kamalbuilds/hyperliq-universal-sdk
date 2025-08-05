package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/hyperliquid-labs/hyperliquid-go-sdk/types"
)

// SubscribeToTrades subscribes to trade updates for a specific coin
func (m *Manager) SubscribeToTrades(coin string) (*Subscription, error) {
	params := map[string]interface{}{
		"type": "trades",
		"coin": coin,
	}
	return m.Subscribe("trades", params)
}

// SubscribeToAllMids subscribes to mid price updates for all coins
func (m *Manager) SubscribeToAllMids() (*Subscription, error) {
	params := map[string]interface{}{
		"type": "allMids",
	}
	return m.Subscribe("allMids", params)
}

// SubscribeToL2Book subscribes to order book updates
func (m *Manager) SubscribeToL2Book(coin string) (*Subscription, error) {
	params := map[string]interface{}{
		"type": "l2Book",
		"coin": coin,
	}
	return m.Subscribe("l2Book", params)
}

// SubscribeToCandles subscribes to candlestick updates
func (m *Manager) SubscribeToCandles(coin, interval string) (*Subscription, error) {
	params := map[string]interface{}{
		"type":     "candle",
		"coin":     coin,
		"interval": interval,
	}
	return m.Subscribe("candle", params)
}

// SubscribeToOrderUpdates subscribes to user's order updates
func (m *Manager) SubscribeToOrderUpdates(user string) (*Subscription, error) {
	params := map[string]interface{}{
		"type": "orderUpdates",
		"user": user,
	}
	return m.Subscribe("orderUpdates", params)
}

// SubscribeToUserEvents subscribes to user events (fills, fundings, liquidations)
func (m *Manager) SubscribeToUserEvents(user string) (*Subscription, error) {
	params := map[string]interface{}{
		"type": "userEvents",
		"user": user,
	}
	return m.Subscribe("userEvents", params)
}

// SubscribeToUserFills subscribes to user's fills only
func (m *Manager) SubscribeToUserFills(user string) (*Subscription, error) {
	params := map[string]interface{}{
		"type": "userFills",
		"user": user,
	}
	return m.Subscribe("userFills", params)
}

// SubscribeToUserFundings subscribes to user's funding payments
func (m *Manager) SubscribeToUserFundings(user string) (*Subscription, error) {
	params := map[string]interface{}{
		"type": "userFundings",
		"user": user,
	}
	return m.Subscribe("userFundings", params)
}

// SubscribeToLiquidations subscribes to liquidation events
func (m *Manager) SubscribeToLiquidations(user *string) (*Subscription, error) {
	params := map[string]interface{}{
		"type": "liquidations",
	}
	if user != nil {
		params["user"] = *user
	}
	return m.Subscribe("liquidations", params)
}

// ParseTradeUpdate parses a trade update message
func ParseTradeUpdate(data json.RawMessage) (*types.Trade, error) {
	var trade types.Trade
	if err := json.Unmarshal(data, &trade); err != nil {
		return nil, fmt.Errorf("failed to parse trade update: %w", err)
	}
	return &trade, nil
}

// ParseL2Update parses an L2 book update message
func ParseL2Update(data json.RawMessage) (*types.L2Book, error) {
	var book types.L2Book
	if err := json.Unmarshal(data, &book); err != nil {
		return nil, fmt.Errorf("failed to parse L2 update: %w", err)
	}
	return &book, nil
}

// ParseCandleUpdate parses a candle update message
func ParseCandleUpdate(data json.RawMessage) (*types.Candle, error) {
	var candle types.Candle
	if err := json.Unmarshal(data, &candle); err != nil {
		return nil, fmt.Errorf("failed to parse candle update: %w", err)
	}
	return &candle, nil
}

// ParseAllMidsUpdate parses an all mids update message
func ParseAllMidsUpdate(data json.RawMessage) (map[string]string, error) {
	var mids map[string]string
	if err := json.Unmarshal(data, &mids); err != nil {
		return nil, fmt.Errorf("failed to parse all mids update: %w", err)
	}
	return mids, nil
}

// ParseOrderUpdate parses an order update message
func ParseOrderUpdate(data json.RawMessage) (map[string]interface{}, error) {
	var update map[string]interface{}
	if err := json.Unmarshal(data, &update); err != nil {
		return nil, fmt.Errorf("failed to parse order update: %w", err)
	}
	return update, nil
}

// ParseUserFillUpdate parses a user fill update
func ParseUserFillUpdate(data json.RawMessage) (*types.Fill, error) {
	var fill types.Fill
	if err := json.Unmarshal(data, &fill); err != nil {
		return nil, fmt.Errorf("failed to parse fill update: %w", err)
	}
	return &fill, nil
}

// UserEvent represents a user event (fill, funding, liquidation)
type UserEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// ParseUserEvent parses a user event
func ParseUserEvent(data json.RawMessage) (*UserEvent, error) {
	var event UserEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("failed to parse user event: %w", err)
	}
	return &event, nil
}