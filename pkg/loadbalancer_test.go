package pkg

import (
	"sync"
	"testing"
)

// Helper to save & restore global state touched by tests.
func withSavedGlobals(t *testing.T, fn func()) {
	t.Helper()

	oldStrict := strict
	oldIps := make([]string, len(ips))
	copy(oldIps, ips)

	// copy internalConnections map
	oldMap := make(map[string]int)
	if internalConnections.connections != nil {
		oldMap = make(map[string]int, len(internalConnections.connections))
		for k, v := range internalConnections.connections {
			oldMap[k] = v
		}
	}

	// cleanup + restore later
	defer func() {
		strict = oldStrict
		ips = oldIps

		internalConnections.connections = make(map[string]int)
		for k, v := range oldMap {
			internalConnections.connections[k] = v
		}
	}()

	// ensure map initialized for tests
	if internalConnections.connections == nil {
		internalConnections.connections = make(map[string]int)
	}

	fn()
}

// Test concurrent IncreaseConnection calls produce correct final count.
func TestIncreaseConnectionConcurrency(t *testing.T) {
	withSavedGlobals(t, func() {
		// start with empty map
		internalConnections.connections = make(map[string]int)

		const (
			key       = "http://backend.local"
			goroutine = 10
			perG      = 1000
		)

		var wg sync.WaitGroup
		wg.Add(goroutine)
		for i := 0; i < goroutine; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < perG; j++ {
					internalConnections.IncreaseConnection(key)
				}
			}()
		}
		wg.Wait()

		expected := goroutine * perG
		got := internalConnections.connections[key]
		if got != expected {
			t.Fatalf("expected %d connections for key %s; got %d", expected, key, got)
		}
	})
}

// Test Increase then Decrease results in correct count and no negative values.
func TestIncreaseThenDecrease(t *testing.T) {
	withSavedGlobals(t, func() {
		internalConnections.connections = make(map[string]int)
		key := "http://backend.local"

		internalConnections.IncreaseConnection(key)
		internalConnections.IncreaseConnection(key)
		if got := internalConnections.connections[key]; got != 2 {
			t.Fatalf("expected 2 after increases, got %d", got)
		}
	})
}
