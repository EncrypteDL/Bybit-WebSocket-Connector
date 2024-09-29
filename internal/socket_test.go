package internal

import (
	"context"
	"log"
	"testing"
	"time"

	"nhooyr.io/websocket"
)

// Mock configuration struct used for testing
type MockConfig struct {
	url    string
	key    string
	secret string
	l      *log.Logger
}

// TestConnect tests the basic connection process
func TestConnect(t *testing.T) {
	// Define mock configuration
	cfg := &Confirguration{
		url:    "wss://stream.bybit.com/realtime_public",
		key:    "your_key",
		secret: "your_secret",
	}

	// Create a context with a timeout to prevent the test from hanging indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a channel to receive responses
	ch := make(chan Responce, 1)

	// Define some example channels and symbols
	channels := []string{BybitWSOrderBookL2_25, BybitWSTrade}
	symbols := []string{"BTCUSD"}

	// Run the Connect function in a separate goroutine
	go func() {
		err := Connect(ctx, ch, channels, symbols, cfg)
		if err != nil {
			t.Errorf("failed to connect: %v", err)
		}
	}()

	// Wait for the context to timeout or a response to be received
	select {
	case <-ctx.Done():
		// Timeout occurred
		t.Fatal("Test timed out before receiving a response")
	case res := <-ch:
		// Check if response type is as expected
		if res.Type == ERROR {
			t.Errorf("received error response: %v", res.Result)
		}
	}
}

// TestPing tests the ping function
func TestPing(t *testing.T) {
	// Establish a mock WebSocket connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "wss://stream.bybit.com/realtime_public", nil)
	if err != nil {
		t.Fatalf("failed to dial websocket: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "normal closure")

	// Test the ping function
	err = ping(ctx, conn)
	if err != nil {
		t.Errorf("ping failed: %v", err)
	}
}
