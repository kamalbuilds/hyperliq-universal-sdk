package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hyperliquid-labs/hyperliquid-go-sdk/client"
	"github.com/hyperliquid-labs/hyperliquid-go-sdk/websocket"
)

func main() {
	// Create WebSocket manager
	ws := websocket.NewManager(client.MainnetWS)

	// Connect to WebSocket
	fmt.Println("Connecting to Hyperliquid WebSocket...")
	if err := ws.Connect(); err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer ws.Close()

	fmt.Println("Connected successfully!")

	// Subscribe to BTC trades
	fmt.Println("\nSubscribing to BTC trades...")
	tradesSub, err := ws.SubscribeToTrades("BTC")
	if err != nil {
		log.Fatal("Failed to subscribe to trades:", err)
	}

	// Subscribe to BTC order book
	fmt.Println("Subscribing to BTC order book...")
	bookSub, err := ws.SubscribeToL2Book("BTC")
	if err != nil {
		log.Fatal("Failed to subscribe to order book:", err)
	}

	// Subscribe to all mid prices
	fmt.Println("Subscribing to all mid prices...")
	midsSub, err := ws.SubscribeToAllMids()
	if err != nil {
		log.Fatal("Failed to subscribe to mids:", err)
	}

	// Subscribe to 1m candles
	fmt.Println("Subscribing to BTC 1m candles...")
	candlesSub, err := ws.SubscribeToCandles("BTC", "1m")
	if err != nil {
		log.Fatal("Failed to subscribe to candles:", err)
	}

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("\nListening for updates... (Press Ctrl+C to exit)")

	// Process updates
	for {
		select {
		case <-sigChan:
			fmt.Println("\nShutting down...")
			return

		case msg := <-tradesSub.Channel:
			trade, err := websocket.ParseTradeUpdate(msg)
			if err != nil {
				fmt.Printf("Error parsing trade: %v\n", err)
				continue
			}
			fmt.Printf("Trade: %s %s %s @ %s\n", 
				trade.Coin, 
				trade.Side, 
				trade.Sz, 
				trade.Px,
			)

		case msg := <-bookSub.Channel:
			book, err := websocket.ParseL2Update(msg)
			if err != nil {
				fmt.Printf("Error parsing book: %v\n", err)
				continue
			}
			
			// Show top of book
			if len(book.Levels) >= 2 {
				// Levels are [bids, asks]
				if bids, ok := book.Levels[0].([]interface{}); ok && len(bids) > 0 {
					if bid, ok := bids[0].([]interface{}); ok && len(bid) >= 2 {
						fmt.Printf("Best Bid: %v @ %v\n", bid[1], bid[0])
					}
				}
				if asks, ok := book.Levels[1].([]interface{}); ok && len(asks) > 0 {
					if ask, ok := asks[0].([]interface{}); ok && len(ask) >= 2 {
						fmt.Printf("Best Ask: %v @ %v\n", ask[1], ask[0])
					}
				}
			}

		case msg := <-midsSub.Channel:
			mids, err := websocket.ParseAllMidsUpdate(msg)
			if err != nil {
				fmt.Printf("Error parsing mids: %v\n", err)
				continue
			}
			
			// Show a few key assets
			for _, asset := range []string{"BTC", "ETH", "SOL"} {
				if price, ok := mids[asset]; ok {
					fmt.Printf("%s: %s  ", asset, price)
				}
			}
			fmt.Println()

		case msg := <-candlesSub.Channel:
			candle, err := websocket.ParseCandleUpdate(msg)
			if err != nil {
				fmt.Printf("Error parsing candle: %v\n", err)
				continue
			}
			fmt.Printf("Candle: O:%s H:%s L:%s C:%s V:%s\n",
				candle.O,
				candle.H,
				candle.L,
				candle.C,
				candle.V,
			)
		}
	}
}