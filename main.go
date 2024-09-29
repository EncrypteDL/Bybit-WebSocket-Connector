package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/EncrypteDL/Bybit-WebSocket-Connector/internal"
)

func main() {
	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configuration for Bybit WebSocket
	cfg := internal.Confirguration{
		// Url:    "wss://stream.bybit.com/realtime_public", // Bybit public WebSocket URL
		// Key:    "your_api_key",                           // Your API key (if needed for private data)
		// Secret: "your_api_secret",                        // Your API secret (if needed for private data)
	}

	// Create a channel to receive WebSocket data
	ch := make(chan internal.Responce)

	// Define channels to subscribe to (e.g., order book, trades, etc.)
	channels := []string{internal.BybitWSOrderBookL2_25, internal.BybitWSTrade}

	// Define symbols to subscribe to (e.g., BTCUSD, ETHUSD, etc.)
	symbols := []string{"BTCUSD", "ETHUSD"}

	// Connect to the Bybit WebSocket
	go func() {
		if err := internal.Connect(ctx, ch, channels, symbols, cfg); err != nil {
			log.Fatalf("Error connecting to Bybit WebSocket: %v", err)
		}
	}()

	// Process WebSocket responses
	go func() {
		for {
			select {
			case res := <-ch:
				// Handle different message types
				switch res.Type {
				case internal.ORDERBOOK:
					fmt.Printf("Received order book update for %s: %+v\n", res.Symbol, res.Orderbook)
				case internal.TRADES:
					fmt.Printf("Received trades for %s: %+v\n", res.Symbol, res.Trades)
				case internal.TICKER:
					fmt.Printf("Received ticker for %s: %+v\n", res.Symbol, res.Ticker)
				case internal.ERROR:
					fmt.Printf("Error: %v\n", res.Result)
				default:
					fmt.Printf("Undefined message type: %v\n", res.Result)
				}
			case <-ctx.Done():
				log.Println("Context canceled, stopping WebSocket handler.")
				return
			}
		}
	}()

	// Keep the main function running (or do other work)
	time.Sleep(10 * time.Minute)
	log.Println("Finished execution")
}
