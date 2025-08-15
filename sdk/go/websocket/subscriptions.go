package websocket

import (
	"encoding/json"
	"github.com/hyperliquid-labs/hyperliquid-go-sdk/types"
)

// Subscription helper functions

func (m *Manager) SubscribeToAllMids(handler func(data types.AllMidsData) error) (string, error) {
	sub := types.WSSubscription{
		Type: "allMids",
	}

	return m.Subscribe(sub, func(raw json.RawMessage) error {
		var data types.AllMidsData
		if err := json.Unmarshal(raw, &data); err != nil {
			return err
		}
		return handler(data)
	})
}

func (m *Manager) SubscribeToL2Book(coin string, handler func(data types.L2BookData) error) (string, error) {
	sub := types.WSSubscription{
		Type: "l2Book",
		Coin: coin,
	}

	return m.Subscribe(sub, func(raw json.RawMessage) error {
		var data types.L2BookData
		if err := json.Unmarshal(raw, &data); err != nil {
			return err
		}
		return handler(data)
	})
}

func (m *Manager) SubscribeToTrades(coin string, handler func(data []types.TradeData) error) (string, error) {
	sub := types.WSSubscription{
		Type: "trades",
		Coin: coin,
	}

	return m.Subscribe(sub, func(raw json.RawMessage) error {
		var data []types.TradeData
		if err := json.Unmarshal(raw, &data); err != nil {
			return err
		}
		return handler(data)
	})
}

func (m *Manager) SubscribeToCandles(coin, interval string, handler func(data types.CandleData) error) (string, error) {
	sub := types.WSSubscription{
		Type:     "candle",
		Coin:     coin,
		Interval: interval,
	}

	return m.Subscribe(sub, func(raw json.RawMessage) error {
		var data types.CandleData
		if err := json.Unmarshal(raw, &data); err != nil {
			return err
		}
		return handler(data)
	})
}

func (m *Manager) SubscribeToUserEvents(user string, handler func(data types.UserEvent) error) (string, error) {
	sub := types.WSSubscription{
		Type: "userEvents",
		User: user,
	}

	return m.Subscribe(sub, func(raw json.RawMessage) error {
		var data types.UserEvent
		if err := json.Unmarshal(raw, &data); err != nil {
			return err
		}
		return handler(data)
	})
}

func (m *Manager) SubscribeToUserFills(user string, handler func(data types.UserFillData) error) (string, error) {
	sub := types.WSSubscription{
		Type: "userFills",
		User: user,
	}

	return m.Subscribe(sub, func(raw json.RawMessage) error {
		var data types.UserFillData
		if err := json.Unmarshal(raw, &data); err != nil {
			return err
		}
		return handler(data)
	})
}

func (m *Manager) SubscribeToOrderUpdates(user string, handler func(data types.OrderUpdate) error) (string, error) {
	sub := types.WSSubscription{
		Type: "orderUpdates",
		User: user,
	}

	return m.Subscribe(sub, func(raw json.RawMessage) error {
		var data types.OrderUpdate
		if err := json.Unmarshal(raw, &data); err != nil {
			return err
		}
		return handler(data)
	})
}

func (m *Manager) SubscribeToUserFunding(user string, handler func(data types.FundingData) error) (string, error) {
	sub := types.WSSubscription{
		Type: "userFundings",
		User: user,
	}

	return m.Subscribe(sub, func(raw json.RawMessage) error {
		var data types.FundingData
		if err := json.Unmarshal(raw, &data); err != nil {
			return err
		}
		return handler(data)
	})
}

func (m *Manager) SubscribeToBBO(coin string, handler func(data types.BboData) error) (string, error) {
	sub := types.WSSubscription{
		Type: "bbo",
		Coin: coin,
	}

	return m.Subscribe(sub, func(raw json.RawMessage) error {
		var data types.BboData
		if err := json.Unmarshal(raw, &data); err != nil {
			return err
		}
		return handler(data)
	})
}

func (m *Manager) SubscribeToActiveAssetCtx(coin string, handler func(data types.ActiveAssetCtxData) error) (string, error) {
	sub := types.WSSubscription{
		Type: "activeAssetCtx",
		Coin: coin,
	}

	return m.Subscribe(sub, func(raw json.RawMessage) error {
		var data types.ActiveAssetCtxData
		if err := json.Unmarshal(raw, &data); err != nil {
			return err
		}
		return handler(data)
	})
}

func (m *Manager) SubscribeToActiveAssetData(coin, user string, handler func(data types.ActiveAssetDataData) error) (string, error) {
	sub := types.WSSubscription{
		Type: "activeAssetData",
		Coin: coin,
		User: user,
	}

	return m.Subscribe(sub, func(raw json.RawMessage) error {
		var data types.ActiveAssetDataData
		if err := json.Unmarshal(raw, &data); err != nil {
			return err
		}
		return handler(data)
	})
}

func (m *Manager) SubscribeToWebData2(user string, handler func(data types.WebData2Data) error) (string, error) {
	sub := types.WSSubscription{
		Type: "webData2",
		User: user,
	}

	return m.Subscribe(sub, func(raw json.RawMessage) error {
		var data types.WebData2Data
		if err := json.Unmarshal(raw, &data); err != nil {
			return err
		}
		return handler(data)
	})
}