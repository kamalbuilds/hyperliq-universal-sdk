package types

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

// OrderRequest represents a request to place an order
type OrderRequest struct {
	Asset       string          `json:"coin"`
	IsBuy       bool            `json:"is_buy"`
	LimitPx     decimal.Decimal `json:"limit_px"`
	Sz          decimal.Decimal `json:"sz"`
	ReduceOnly  bool            `json:"reduce_only"`
	OrderType   OrderType       `json:"order_type"`
	Cloid       *string         `json:"cloid,omitempty"`
}

// OrderType represents the type of order
type OrderType struct {
	Limit   *LimitOrderType   `json:"limit,omitempty"`
	Trigger *TriggerOrderType `json:"trigger,omitempty"`
}

// LimitOrderType represents a limit order configuration
type LimitOrderType struct {
	Tif string `json:"tif"` // "Gtc", "Ioc", "Alo"
}

// TriggerOrderType represents a trigger order configuration
type TriggerOrderType struct {
	TriggerPx decimal.Decimal `json:"trigger_px"`
	IsMarket  bool            `json:"is_market"`
	TpSl      string          `json:"tp_sl"` // "tp" or "sl"
}

// CancelRequest represents a request to cancel orders
type CancelRequest struct {
	Asset string  `json:"coin"`
	Oid   *int64  `json:"oid,omitempty"`
	Cloid *string `json:"cloid,omitempty"`
}

// ModifyRequest represents a request to modify an order
type ModifyRequest struct {
	Asset   string          `json:"coin"`
	Oid     int64           `json:"oid"`
	IsBuy   bool            `json:"is_buy"`
	LimitPx decimal.Decimal `json:"limit_px"`
	Sz      decimal.Decimal `json:"sz"`
}

// UserState represents a user's account state
type UserState struct {
	MarginSummary  MarginSummary   `json:"marginSummary"`
	AssetPositions []AssetPosition `json:"assetPositions"`
	CrossMarginSummary CrossMarginSummary `json:"crossMarginSummary"`
}

// MarginSummary represents margin information
type MarginSummary struct {
	AccountValue    decimal.Decimal `json:"accountValue"`
	TotalMarginUsed decimal.Decimal `json:"totalMarginUsed"`
	TotalNtlPos     decimal.Decimal `json:"totalNtlPos"`
	TotalRawUsd     decimal.Decimal `json:"totalRawUsd"`
	WithdrawableUsd decimal.Decimal `json:"withdrawableUsd"`
}

// CrossMarginSummary represents cross margin information
type CrossMarginSummary struct {
	AccountValue    decimal.Decimal `json:"accountValue"`
	TotalMarginUsed decimal.Decimal `json:"totalMarginUsed"`
	TotalNtlPos     decimal.Decimal `json:"totalNtlPos"`
	TotalRawUsd     decimal.Decimal `json:"totalRawUsd"`
}

// AssetPosition represents a position in an asset
type AssetPosition struct {
	Position Position `json:"position"`
	Type     string   `json:"type"`
}

// Position represents position details
type Position struct {
	Coin           string          `json:"coin"`
	EntryPx        decimal.Decimal `json:"entryPx"`
	Szi            decimal.Decimal `json:"szi"`
	Leverage       Leverage        `json:"leverage"`
	UnrealizedPnl  decimal.Decimal `json:"unrealizedPnl"`
	RealizedPnl    decimal.Decimal `json:"realizedPnl"`
	CumFunding     CumFunding      `json:"cumFunding"`
	PositionValue  decimal.Decimal `json:"positionValue"`
	MaxTradeSz     decimal.Decimal `json:"maxTradeSz"`
	MarginUsed     decimal.Decimal `json:"marginUsed"`
}

// Leverage represents leverage information
type Leverage struct {
	Type    string          `json:"type"`
	Value   decimal.Decimal `json:"value"`
	RawUsd  decimal.Decimal `json:"rawUsd"`
}

// CumFunding represents cumulative funding
type CumFunding struct {
	AllTime    decimal.Decimal `json:"allTime"`
	SinceOpen  decimal.Decimal `json:"sinceOpen"`
	SinceChange decimal.Decimal `json:"sinceChange"`
}

// L2Book represents a level 2 order book
type L2Book struct {
	Coin   string          `json:"coin"`
	Time   int64           `json:"time"`
	Levels [][]interface{} `json:"levels"`
}

// Level represents an order book level
type Level struct {
	Px   decimal.Decimal `json:"px"`
	Sz   decimal.Decimal `json:"sz"`
	N    int             `json:"n"`
}

// Trade represents a trade
type Trade struct {
	Coin  string          `json:"coin"`
	Side  string          `json:"side"`
	Px    decimal.Decimal `json:"px"`
	Sz    decimal.Decimal `json:"sz"`
	Time  int64           `json:"time"`
	Hash  string          `json:"hash"`
	Tid   int64           `json:"tid"`
}

// Candle represents a candlestick
type Candle struct {
	T int64           `json:"t"`
	O decimal.Decimal `json:"o"`
	H decimal.Decimal `json:"h"`
	L decimal.Decimal `json:"l"`
	C decimal.Decimal `json:"c"`
	V decimal.Decimal `json:"v"`
	N int             `json:"n"`
}

// OpenOrder represents an open order
type OpenOrder struct {
	Coin        string          `json:"coin"`
	LimitPx     decimal.Decimal `json:"limitPx"`
	Oid         int64           `json:"oid"`
	Side        string          `json:"side"`
	Sz          decimal.Decimal `json:"sz"`
	Timestamp   int64           `json:"timestamp"`
	OrigSz      decimal.Decimal `json:"origSz"`
	Cloid       *string         `json:"cloid,omitempty"`
	ReduceOnly  bool            `json:"reduceOnly"`
	OrderType   string          `json:"orderType"`
}

// Fill represents a filled order
type Fill struct {
	Coin         string          `json:"coin"`
	Px           decimal.Decimal `json:"px"`
	Sz           decimal.Decimal `json:"sz"`
	Side         string          `json:"side"`
	Time         int64           `json:"time"`
	StartPosition decimal.Decimal `json:"startPosition"`
	Dir          string          `json:"dir"`
	ClosedPnl    decimal.Decimal `json:"closedPnl"`
	Hash         string          `json:"hash"`
	Oid          int64           `json:"oid"`
	Crossed      bool            `json:"crossed"`
	Fee          decimal.Decimal `json:"fee"`
	Tid          int64           `json:"tid"`
	FeeToken     string          `json:"feeToken"`
}

// FundingHistory represents funding payment history
type FundingHistory struct {
	Coin    string          `json:"coin"`
	FundingRate decimal.Decimal `json:"fundingRate"`
	Szi     decimal.Decimal `json:"szi"`
	Type    string          `json:"type"`
	Time    int64           `json:"time"`
	Usdc    decimal.Decimal `json:"usdc"`
}

// SpotMeta represents spot token metadata
type SpotMeta struct {
	Universe []SpotToken `json:"universe"`
}

// SpotToken represents a spot token
type SpotToken struct {
	Name     string `json:"name"`
	SzDecimals int    `json:"szDecimals"`
	TokenId  int    `json:"tokenId"`
}

// Meta represents market metadata
type Meta struct {
	Universe []Asset `json:"universe"`
}

// Asset represents an asset/market
type Asset struct {
	Name        string          `json:"name"`
	SzDecimals  int             `json:"szDecimals"`
	MaxLeverage int             `json:"maxLeverage"`
	OnlyIsolated bool           `json:"onlyIsolated"`
}

// TransferRequest represents a transfer request
type TransferRequest struct {
	Destination string          `json:"destination"`
	Amount      decimal.Decimal `json:"amount"`
	Asset       string          `json:"asset"`
}

// WithdrawRequest represents a withdrawal request  
type WithdrawRequest struct {
	Destination string          `json:"destination"`
	Amount      decimal.Decimal `json:"amount"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Status   string          `json:"status"`
	Response json.RawMessage `json:"response"`
	Error    *string         `json:"error,omitempty"`
}

// OrderResponse represents a response from order placement
type OrderResponse struct {
	Status string `json:"status"`
	Response struct {
		Type string `json:"type"`
		Data struct {
			Statuses []OrderStatus `json:"statuses"`
		} `json:"data"`
	} `json:"response"`
}

// OrderStatus represents the status of an order
type OrderStatus struct {
	Resting     *RestingOrder `json:"resting,omitempty"`
	Filled      *FilledOrder  `json:"filled,omitempty"`
	Error       *string       `json:"error,omitempty"`
	Frontend    *string       `json:"frontend,omitempty"`
}

// RestingOrder represents a resting order
type RestingOrder struct {
	Oid   int64  `json:"oid"`
	Cloid string `json:"cloid,omitempty"`
}

// FilledOrder represents a filled order  
type FilledOrder struct {
	TotalSz decimal.Decimal `json:"totalSz"`
	AvgPx   decimal.Decimal `json:"avgPx"`
	Oid     int64           `json:"oid"`
}