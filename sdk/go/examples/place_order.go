package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hyperliquid-labs/hyperliquid-go-sdk/client"
	"github.com/hyperliquid-labs/hyperliquid-go-sdk/types"
	"github.com/hyperliquid-labs/hyperliquid-go-sdk/utils"
	"github.com/shopspring/decimal"
)

func main() {
	// Get private key from environment variable
	privateKey := os.Getenv("HYPERLIQUID_PRIVATE_KEY")
	if privateKey == "" {
		log.Fatal("HYPERLIQUID_PRIVATE_KEY environment variable not set")
	}

	// Create client
	c := client.NewMainnetClient(privateKey)
	
	// Get address from private key
	address, err := utils.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		log.Fatal("Failed to get address from private key:", err)
	}
	c.SetAddress(address)

	ctx := context.Background()

	// Get account state
	fmt.Println("Fetching account state...")
	userState, err := c.Info().GetUserState(ctx, address)
	if err != nil {
		log.Fatal("Failed to get user state:", err)
	}

	fmt.Printf("Account Value: %s\n", userState.MarginSummary.AccountValue)
	fmt.Printf("Withdrawable: %s\n", userState.MarginSummary.WithdrawableUsd)
	fmt.Println()

	// Get current BTC price
	fmt.Println("Fetching current prices...")
	mids, err := c.Info().GetAllMids(ctx)
	if err != nil {
		log.Fatal("Failed to get mid prices:", err)
	}

	btcPrice := mids["BTC"]
	fmt.Printf("Current BTC price: %s\n", btcPrice)

	// Calculate order price (1% below current price for a limit buy)
	currentPrice, err := decimal.NewFromString(btcPrice)
	if err != nil {
		log.Fatal("Failed to parse BTC price:", err)
	}
	
	orderPrice := currentPrice.Mul(decimal.NewFromFloat(0.99))
	orderSize := decimal.NewFromFloat(0.001) // 0.001 BTC

	// Create limit order
	order := types.OrderRequest{
		Asset:      "BTC",
		IsBuy:      true,
		LimitPx:    orderPrice,
		Sz:         orderSize,
		ReduceOnly: false,
		OrderType: types.OrderType{
			Limit: &types.LimitOrderType{
				Tif: "Gtc", // Good till cancelled
			},
		},
	}

	fmt.Printf("\nPlacing limit buy order:\n")
	fmt.Printf("Asset: %s\n", order.Asset)
	fmt.Printf("Price: %s\n", order.LimitPx)
	fmt.Printf("Size: %s\n", order.Sz)
	fmt.Printf("Type: Limit GTC\n")

	// Place the order
	orderResp, err := c.Exchange().PlaceOrder(ctx, order)
	if err != nil {
		log.Fatal("Failed to place order:", err)
	}

	// Check order status
	if orderResp.Status == "success" {
		fmt.Println("\nOrder placed successfully!")
		
		// Print order details
		for _, status := range orderResp.Response.Data.Statuses {
			if status.Resting != nil {
				fmt.Printf("Order ID: %d\n", status.Resting.Oid)
				if status.Resting.Cloid != "" {
					fmt.Printf("Client Order ID: %s\n", status.Resting.Cloid)
				}
			} else if status.Filled != nil {
				fmt.Printf("Order filled!\n")
				fmt.Printf("Filled size: %s\n", status.Filled.TotalSz)
				fmt.Printf("Average price: %s\n", status.Filled.AvgPx)
			} else if status.Error != nil {
				fmt.Printf("Order error: %s\n", *status.Error)
			}
		}
	} else {
		fmt.Printf("Order failed: %v\n", orderResp)
	}

	// Check open orders
	fmt.Println("\nChecking open orders...")
	openOrders, err := c.Info().GetOpenOrders(ctx, address)
	if err != nil {
		log.Fatal("Failed to get open orders:", err)
	}

	fmt.Printf("Open orders: %d\n", len(openOrders))
	for _, order := range openOrders {
		fmt.Printf("- %s %s: %s @ %s (OID: %d)\n", 
			order.Coin, 
			order.Side, 
			order.Sz, 
			order.LimitPx, 
			order.Oid,
		)
	}
}