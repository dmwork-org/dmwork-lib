package keylock

import (
	"sync"
	"testing"
)

func TestKeyLock_ConcurrentLockUnlock(t *testing.T) {
	kl := NewKeyLock()
	const goroutines = 100
	const iterations = 100

	counter := make(map[string]int)
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				kl.Lock("test")
				counter["test"]++
				kl.Unlock("test")
			}
		}()
	}
	wg.Wait()

	if counter["test"] != goroutines*iterations {
		t.Errorf("expected %d, got %d", goroutines*iterations, counter["test"])
	}
}

func TestKeyLock_ConcurrentLockUnlockWithClean(t *testing.T) {
	kl := NewKeyLock()
	const goroutines = 50
	const iterations = 50

	counter := make(map[string]int)
	var wg sync.WaitGroup
	wg.Add(goroutines + 1)

	// Concurrent clean
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			kl.Clean()
		}
	}()

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				kl.Lock("test")
				counter["test"]++
				kl.Unlock("test")
			}
		}()
	}
	wg.Wait()

	if counter["test"] != goroutines*iterations {
		t.Errorf("expected %d, got %d", goroutines*iterations, counter["test"])
	}
}

func TestKeyLock_MultipleKeys(t *testing.T) {
	kl := NewKeyLock()
	var wg sync.WaitGroup

	keys := []string{"a", "b", "c"}
	counters := make(map[string]int)
	var mu sync.Mutex

	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				kl.Lock(k)
				mu.Lock()
				counters[k]++
				mu.Unlock()
				kl.Unlock(k)
			}
		}(key)
	}
	wg.Wait()

	for _, key := range keys {
		if counters[key] != 100 {
			t.Errorf("key %s: expected 100, got %d", key, counters[key])
		}
	}
}
