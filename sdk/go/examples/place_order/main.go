package main

import (
	"context"
	"fmt"
	"log"
	"os"
	
	"github.com/hyperliquid-labs/hyperliquid-go-sdk/client"
	"github.com/hyperliquid-labs/hyperliquid-go-sdk/types"
	"github.com/shopspring/decimal"
)

func main() {
	// Get private key from environment
	privateKey := os.Getenv("HYPERLIQUID_PRIVATE_KEY")
	if privateKey == "" {
		log.Fatal("Please set HYPERLIQUID_PRIVATE_KEY environment variable")
	}
	
	// Create client (use testnet for safety)
	c := client.NewTestnetClient(privateKey)
	ctx := context.Background()
	
	fmt.Println("🚀 Hyperliquid Order Placement Example")
	fmt.Println("=====================================")
	
	// Get current market price
	fmt.Println("\n📊 Fetching current market prices...")
	mids, err := c.Info().GetAllMids(ctx)
	if err != nil {
		log.Fatalf("Failed to get market prices: %v", err)
	}
	
	btcPrice := mids["BTC"]
	fmt.Printf("Current BTC price: %s\n", btcPrice)
	
	// Calculate order price (1% below market for a limit buy)
	currentPrice, _ := decimal.NewFromString(btcPrice)
	orderPrice := currentPrice.Mul(decimal.NewFromFloat(0.99))
	
	// Create order request
	order := types.OrderRequest{
		Asset:      "BTC",
		IsBuy:      true,
		LimitPx:    orderPrice,
		Sz:         decimal.NewFromFloat(0.001), // Small size for testing
		ReduceOnly: false,
		OrderType: types.OrderType{
			Limit: &types.LimitOrderType{
				Tif: "Gtc", // Good till canceled
			},
		},
	}
	
	fmt.Printf("\n📝 Placing order:\n")
	fmt.Printf("   Asset: %s\n", order.Asset)
	fmt.Printf("   Side: BUY\n")
	fmt.Printf("   Price: %s\n", order.LimitPx)
	fmt.Printf("   Size: %s\n", order.Sz)
	
	// Place the order
	resp, err := c.Exchange().PlaceOrder(ctx, order)
	if err != nil {
		log.Fatalf("Failed to place order: %v", err)
	}
	
	if resp.Status == "success" {
		fmt.Println("\n✅ Order placed successfully!")
		
		// Check order status
		for i, status := range resp.Response.Data.Statuses {
			if status.Resting != nil {
				fmt.Printf("   Order %d ID: %d\n", i+1, status.Resting.Oid)
			} else if status.Filled != nil {
				fmt.Printf("   Order %d filled: %s @ %s\n", 
					i+1, status.Filled.TotalSz, status.Filled.AvgPx)
			} else if status.Error != nil {
				fmt.Printf("   Order %d error: %s\n", i+1, *status.Error)
			}
		}
		
		// Get open orders
		fmt.Println("\n📋 Fetching open orders...")
		openOrders, err := c.Info().GetOpenOrders(ctx, c.GetAddress())
		if err == nil {
			fmt.Printf("Found %d open orders\n", len(openOrders))
			for _, o := range openOrders {
				fmt.Printf("   Order: %s %s %s @ %s\n", 
					o.Coin, o.Side, o.Sz, o.LimitPx)
			}
		}
		
	} else {
		fmt.Printf("\n❌ Order failed: %v\n", resp.Error)
	}
}