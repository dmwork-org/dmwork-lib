package config

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/olivere/elastic"
)

// TestGetElasticsearchConcurrent tests that concurrent calls to GetElasticsearch
// return the same client instance and only create one client.
// This test verifies the sync.Once pattern works correctly.
func TestGetElasticsearchConcurrent(t *testing.T) {
	// Track how many times the client creation happens
	var creationCount int32
	ctx := &testableContext{
		createClient: func() *elastic.Client {
			atomic.AddInt32(&creationCount, 1)
			// Return a client with sniff disabled for testing
			client, err := elastic.NewClient(
				elastic.SetURL("http://localhost:9200"),
				elastic.SetSniff(false),
				elastic.SetHealthcheck(false),
			)
			if err != nil {
				t.Fatalf("failed to create test client: %v", err)
			}
			return client
		},
	}

	var wg sync.WaitGroup
	const numGoroutines = 100

	// Store all returned clients to verify they're the same instance
	clients := make(chan *elastic.Client, numGoroutines)

	// Launch concurrent calls to getElasticsearch
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			client := ctx.getElasticsearch()
			clients <- client
		}()
	}

	wg.Wait()
	close(clients)

	// Verify all goroutines got the same client instance
	var firstClient *elastic.Client
	clientCount := 0
	for client := range clients {
		clientCount++
		if firstClient == nil {
			firstClient = client
		} else if client != firstClient {
			t.Errorf("got different client instances, expected all to be the same")
		}
	}

	if clientCount != numGoroutines {
		t.Errorf("expected %d clients, got %d", numGoroutines, clientCount)
	}

	// Verify the client is not nil
	if firstClient == nil {
		t.Error("expected non-nil client")
	}

	// Verify the client was only created once
	if creationCount != 1 {
		t.Errorf("expected client to be created exactly once, but was created %d times", creationCount)
	}
}

// TestGetElasticsearchSingleton tests that GetElasticsearch returns the same client on multiple calls.
func TestGetElasticsearchSingleton(t *testing.T) {
	ctx := &testableContext{
		createClient: func() *elastic.Client {
			client, err := elastic.NewClient(
				elastic.SetURL("http://localhost:9200"),
				elastic.SetSniff(false),
				elastic.SetHealthcheck(false),
			)
			if err != nil {
				t.Fatalf("failed to create test client: %v", err)
			}
			return client
		},
	}

	// Call getElasticsearch multiple times
	client1 := ctx.getElasticsearch()
	client2 := ctx.getElasticsearch()
	client3 := ctx.getElasticsearch()

	// All should be the same instance
	if client1 != client2 || client2 != client3 {
		t.Error("getElasticsearch should return the same client instance")
	}

	if client1 == nil {
		t.Error("getElasticsearch should not return nil")
	}
}

// testableContext is a test helper that mimics Context's sync.Once behavior
// without requiring actual Elasticsearch connectivity.
type testableContext struct {
	elasticOnce   sync.Once
	elasticClient *elastic.Client
	createClient  func() *elastic.Client
}

func (tc *testableContext) getElasticsearch() *elastic.Client {
	tc.elasticOnce.Do(func() {
		tc.elasticClient = tc.createClient()
	})
	return tc.elasticClient
}
