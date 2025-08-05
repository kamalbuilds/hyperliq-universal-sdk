# Hyperliquid Go SDK

A comprehensive Go SDK for interacting with the Hyperliquid DEX, providing easy access to trading, market data, and WebSocket streams.

## Features

- 🚀 Full API coverage (Info & Exchange endpoints)
- 📊 Real-time WebSocket subscriptions
- 🔐 EIP-712 signature authentication
- ⚡ Rate limiting and retry logic
- 🛡️ Type-safe interfaces
- 📝 Comprehensive examples

## Installation

```bash
go get github.com/hyperliquid-labs/hyperliquid-go-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/hyperliquid-labs/hyperliquid-go-sdk/client"
    "github.com/hyperliquid-labs/hyperliquid-go-sdk/types"
    "github.com/shopspring/decimal"
)

func main() {
    // Create client
    c := client.NewMainnetClient("your_private_key")
    
    // Get market prices
    ctx := context.Background()
    mids, err := c.Info().GetAllMids(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("BTC Price: %s\n", mids["BTC"])
    
    // Place an order
    order := types.OrderRequest{
        Asset:   "BTC",
        IsBuy:   true,
        LimitPx: decimal.NewFromFloat(50000),
        Sz:      decimal.NewFromFloat(0.01),
        OrderType: types.OrderType{
            Limit: &types.LimitOrderType{Tif: "Gtc"},
        },
    }
    
    resp, err := c.Exchange().PlaceOrder(ctx, order)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Order placed: %+v\n", resp)
}
```

## Documentation

### Client Initialization

```go
// Mainnet
client := client.NewMainnetClient(privateKey)

// Testnet
client := client.NewTestnetClient(privateKey)

// Custom endpoint
client := client.NewClient(baseURL, wsURL, privateKey)
```

### Info API (Public Data)

The Info API provides access to public market data without authentication:

```go
// Get user state
userState, err := client.Info().GetUserState(ctx, address)

// Get open orders
orders, err := client.Info().GetOpenOrders(ctx, address)

// Get L2 order book
book, err := client.Info().GetL2Book(ctx, "BTC")

// Get recent trades
trades, err := client.Info().GetUserFills(ctx, address, nil, nil)

// Get candles
candles, err := client.Info().GetCandles(ctx, "BTC", "1h", startTime, endTime)

// Get all mid prices
mids, err := client.Info().GetAllMids(ctx)

// Get metadata
meta, err := client.Info().GetMeta(ctx)
```

### Exchange API (Authenticated)

The Exchange API requires authentication and allows trading operations:

```go
// Place a single order
order := types.OrderRequest{
    Asset:      "BTC",
    IsBuy:      true,
    LimitPx:    decimal.NewFromFloat(50000),
    Sz:         decimal.NewFromFloat(0.01),
    ReduceOnly: false,
    OrderType: types.OrderType{
        Limit: &types.LimitOrderType{Tif: "Gtc"},
    },
}
resp, err := client.Exchange().PlaceOrder(ctx, order)

// Place multiple orders
orders := []types.OrderRequest{order1, order2}
resp, err := client.Exchange().PlaceOrders(ctx, orders, "na")

// Cancel order
cancel := types.CancelRequest{
    Asset: "BTC",
    Oid:   &orderID,
}
resp, err := client.Exchange().CancelOrder(ctx, cancel)

// Modify order
modify := types.ModifyRequest{
    Asset:   "BTC",
    Oid:     orderID,
    IsBuy:   true,
    LimitPx: decimal.NewFromFloat(51000),
    Sz:      decimal.NewFromFloat(0.01),
}
resp, err := client.Exchange().ModifyOrder(ctx, modify)

// Update leverage
resp, err := client.Exchange().UpdateLeverage(ctx, "BTC", "cross", 10)

// Transfer USDC
transfer := types.TransferRequest{
    Destination: "0x...",
    Amount:      decimal.NewFromFloat(100),
}
resp, err := client.Exchange().Transfer(ctx, transfer)
```

### WebSocket Subscriptions

```go
// Create WebSocket manager
ws := websocket.NewManager(client.MainnetWS)

// Connect
err := ws.Connect()
defer ws.Close()

// Subscribe to trades
sub, err := ws.SubscribeToTrades("BTC")

// Process updates
for msg := range sub.Channel {
    trade, err := websocket.ParseTradeUpdate(msg)
    if err != nil {
        continue
    }
    fmt.Printf("Trade: %+v\n", trade)
}

// Available subscriptions:
// - SubscribeToTrades(coin)
// - SubscribeToL2Book(coin)
// - SubscribeToCandles(coin, interval)
// - SubscribeToAllMids()
// - SubscribeToOrderUpdates(user)
// - SubscribeToUserFills(user)
// - SubscribeToLiquidations(user)
```

### Order Types

```go
// Market order
order := types.OrderRequest{
    Asset: "BTC",
    IsBuy: true,
    Sz:    decimal.NewFromFloat(0.01),
    OrderType: types.OrderType{
        Limit: &types.LimitOrderType{Tif: "Ioc"},
    },
}

// Limit order (Good Till Canceled)
order := types.OrderRequest{
    Asset:   "BTC",
    IsBuy:   true,
    LimitPx: decimal.NewFromFloat(50000),
    Sz:      decimal.NewFromFloat(0.01),
    OrderType: types.OrderType{
        Limit: &types.LimitOrderType{Tif: "Gtc"},
    },
}

// Stop loss order
order := types.OrderRequest{
    Asset: "BTC",
    IsBuy: false,
    Sz:    decimal.NewFromFloat(0.01),
    OrderType: types.OrderType{
        Trigger: &types.TriggerOrderType{
            TriggerPx: decimal.NewFromFloat(45000),
            IsMarket:  true,
            TpSl:      "sl",
        },
    },
}

// Take profit order
order := types.OrderRequest{
    Asset: "BTC",
    IsBuy: false,
    Sz:    decimal.NewFromFloat(0.01),
    OrderType: types.OrderType{
        Trigger: &types.TriggerOrderType{
            TriggerPx: decimal.NewFromFloat(55000),
            IsMarket:  true,
            TpSl:      "tp",
        },
    },
}
```

### Error Handling

```go
resp, err := client.Exchange().PlaceOrder(ctx, order)
if err != nil {
    // Network or request error
    log.Printf("Request failed: %v", err)
    return
}

if resp.Status != "success" {
    // API error
    log.Printf("Order failed: %s", resp.Error)
    return
}

// Check individual order statuses
for _, status := range resp.Response.Data.Statuses {
    if status.Error != nil {
        log.Printf("Order error: %s", *status.Error)
    }
}
```

## Examples

Check the `/examples` directory for complete examples:

- `place_order.go` - Placing and managing orders
- `websocket_stream.go` - Real-time market data streaming
- `market_maker.go` - Simple market making bot
- `portfolio_tracker.go` - Track positions and P&L

## Testing

Run the test suite:

```bash
go test ./...
```

Run specific tests:

```bash
go test -run TestPlaceOrder
```

## Rate Limits

The SDK implements automatic rate limiting:
- Default: 20 requests per second (1200 per minute)
- WebSocket: Automatic reconnection with exponential backoff
- All API calls respect rate limits automatically

## Security

- Private keys are never logged or transmitted except for signing
- All connections use TLS
- EIP-712 signatures for authentication
- No external dependencies beyond standard crypto libraries

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see LICENSE file for details

## Support

- Documentation: https://hyperliquid.gitbook.io/hyperliquid-docs
- Discord: https://discord.gg/hyperliquid
- Issues: https://github.com/hyperliquid-labs/hyperliquid-go-sdk/issues