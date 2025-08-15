package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/hyperliquid-labs/hyperliquid-go-sdk/client"
    "github.com/hyperliquid-labs/hyperliquid-go-sdk/types"
    "github.com/hyperliquid-labs/hyperliquid-go-sdk/websocket"
)

func main() {
    fmt.Println("🚀 Testing Hyperliquid Go SDK")
    fmt.Println("================================")
    
    // Create testnet client (no real private key needed for info endpoints)
    c := client.NewTestnetClient("")
    ctx := context.Background()
    
    // Test 1: Get all mid prices
    fmt.Println("\n📊 Test 1: Fetching Market Prices...")
    mids, err := c.Info().GetAllMids(ctx)
    if err != nil {
        log.Printf("Error getting mids: %v", err)
    } else {
        fmt.Printf("✅ Successfully fetched %d market prices\n", len(mids))
        if btcPrice, ok := mids["BTC"]; ok {
            fmt.Printf("   BTC Price: %s\n", btcPrice)
        }
    }
    
    // Test 2: Get order book
    fmt.Println("\n📖 Test 2: Fetching Order Book...")
    book, err := c.Info().GetL2Book(ctx, "BTC")
    if err != nil {
        log.Printf("Error getting order book: %v", err)
    } else {
        fmt.Printf("✅ Order book fetched for %s\n", book.Coin)
        fmt.Printf("   Levels: %d\n", len(book.Levels))
    }
    
    // Test 3: WebSocket connection
    fmt.Println("\n🔌 Test 3: Testing WebSocket Connection...")
    ws := websocket.NewManager(client.TestnetWS)
    
    err = ws.Connect(ctx)
    if err != nil {
        log.Printf("Error connecting to WebSocket: %v", err)
    } else {
        fmt.Println("✅ WebSocket connected successfully")
        
        // Subscribe to trades
        subID, err := ws.SubscribeToTrades("BTC", func(trades []types.TradeData) error {
            fmt.Printf("   Trade: %+v\n", trades[0])
            return nil
        })
        
        if err == nil {
            fmt.Printf("✅ Subscribed to BTC trades (ID: %s)\n", subID)
            
            // Listen for 5 seconds
            time.Sleep(5 * time.Second)
            
            // Unsubscribe
            ws.Unsubscribe(subID)
            fmt.Println("✅ Unsubscribed from trades")
        }
        
        ws.Disconnect()
    }
    
    fmt.Println("\n✨ Go SDK tests completed!")
}