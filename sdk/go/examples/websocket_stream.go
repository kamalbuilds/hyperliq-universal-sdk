package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/hyperliquid-labs/hyperliquid-go-sdk/client"
	"github.com/hyperliquid-labs/hyperliquid-go-sdk/types"
	"github.com/hyperliquid-labs/hyperliquid-go-sdk/websocket"
)

func main() {
	fmt.Println("🔌 Hyperliquid WebSocket Streaming Example")
	fmt.Println("=========================================")
	
	// Create WebSocket manager
	ws := websocket.NewManager(client.MainnetWS)
	
	// Connect to WebSocket
	ctx := context.Background()
	fmt.Println("\nConnecting to WebSocket...")
	if err := ws.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer ws.Disconnect()
	
	fmt.Println("✅ Connected successfully!")
	
	// Subscribe to multiple data streams
	fmt.Println("\n📊 Setting up subscriptions...")
	
	// 1. Subscribe to all mid prices
	midSub, err := ws.SubscribeToAllMids(func(data types.AllMidsData) error {
		fmt.Printf("💹 Mid prices update: %d markets\n", len(data.Mids))
		if btc, ok := data.Mids["BTC"]; ok {
			fmt.Printf("   BTC: %s\n", btc)
		}
		return nil
	})
	if err != nil {
		log.Printf("Failed to subscribe to mids: %v", err)
	} else {
		fmt.Printf("✅ Subscribed to all mid prices (ID: %s)\n", midSub)
	}
	
	// 2. Subscribe to BTC trades
	tradeSub, err := ws.SubscribeToTrades("BTC", func(trades []types.TradeData) error {
		for _, trade := range trades {
			emoji := "🟢"
			if trade.Side == "SELL" {
				emoji = "🔴"
			}
			fmt.Printf("%s Trade: %s %s BTC @ %s\n", 
				emoji, trade.Side, trade.Sz, trade.Px)
		}
		return nil
	})
	if err != nil {
		log.Printf("Failed to subscribe to trades: %v", err)
	} else {
		fmt.Printf("✅ Subscribed to BTC trades (ID: %s)\n", tradeSub)
	}
	
	// 3. Subscribe to BTC order book
	bookSub, err := ws.SubscribeToL2Book("BTC", func(book types.L2BookData) error {
		fmt.Printf("📖 Order book update: %s\n", book.Coin)
		fmt.Printf("   Levels: %d\n", len(book.Levels))
		return nil
	})
	if err != nil {
		log.Printf("Failed to subscribe to order book: %v", err)
	} else {
		fmt.Printf("✅ Subscribed to BTC order book (ID: %s)\n", bookSub)
	}
	
	// 4. Subscribe to 1-hour candles
	candleSub, err := ws.SubscribeToCandles("BTC", "1h", func(candle types.CandleData) error {
		fmt.Printf("🕯️ Candle: O:%s H:%s L:%s C:%s V:%s\n",
			candle.O, candle.H, candle.L, candle.C, candle.V)
		return nil
	})
	if err != nil {
		log.Printf("Failed to subscribe to candles: %v", err)
	} else {
		fmt.Printf("✅ Subscribed to BTC 1h candles (ID: %s)\n", candleSub)
	}
	
	// Monitor connection stats
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		
		for range ticker.C {
			stats := ws.GetStats()
			fmt.Printf("\n📈 Connection Stats:\n")
			fmt.Printf("   Connected: %v\n", stats.Connected)
			fmt.Printf("   Subscriptions: %d\n", stats.Subscriptions)
			fmt.Printf("   Reconnects: %d\n", stats.Reconnects)
		}
	}()
	
	// Wait for interrupt signal
	fmt.Println("\n🎯 Streaming live data... Press Ctrl+C to stop")
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	
	// Clean shutdown
	fmt.Println("\n\n🛑 Shutting down...")
	
	// Unsubscribe from all
	fmt.Println("Unsubscribing from streams...")
	ws.Unsubscribe(midSub)
	ws.Unsubscribe(tradeSub)
	ws.Unsubscribe(bookSub)
	ws.Unsubscribe(candleSub)
	
	fmt.Println("✅ Shutdown complete!")
}