package keylock

import (
	"sync"
	"testing"
	"time"
)

func TestKeyLock_LockUnlock(t *testing.T) {
	kl := NewKeyLock()

	kl.Lock("key1")
	kl.Unlock("key1")

	// Should not deadlock
	kl.Lock("key1")
	kl.Unlock("key1")
}

func TestKeyLock_ConcurrentLock(t *testing.T) {
	kl := NewKeyLock()
	var counter int
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			kl.Lock("key1")
			counter++
			kl.Unlock("key1")
		}()
	}

	wg.Wait()
	if counter != 100 {
		t.Errorf("expected counter to be 100, got %d", counter)
	}
}

func TestKeyLock_DifferentKeys(t *testing.T) {
	kl := NewKeyLock()
	var wg sync.WaitGroup

	// Different keys should not block each other
	wg.Add(2)

	go func() {
		defer wg.Done()
		kl.Lock("key1")
		time.Sleep(10 * time.Millisecond)
		kl.Unlock("key1")
	}()

	go func() {
		defer wg.Done()
		kl.Lock("key2")
		time.Sleep(10 * time.Millisecond)
		kl.Unlock("key2")
	}()

	wg.Wait()
}

func TestKeyLock_Clean(t *testing.T) {
	kl := NewKeyLock()

	kl.Lock("key1")
	kl.Unlock("key1")

	kl.Clean()

	kl.mutex.Lock()
	_, exists := kl.locks["key1"]
	kl.mutex.Unlock()

	if exists {
		t.Error("expected key1 to be cleaned up")
	}
}

func TestKeyLock_CleanDoesNotRemoveActiveLock(t *testing.T) {
	kl := NewKeyLock()

	kl.Lock("key1")
	kl.Clean()

	kl.mutex.Lock()
	_, exists := kl.locks["key1"]
	kl.mutex.Unlock()

	if !exists {
		t.Error("expected key1 to still exist while locked")
	}

	kl.Unlock("key1")
}

func TestKeyLock_StartStopCleanLoop(t *testing.T) {
	kl := NewKeyLock()
	kl.StartCleanLoop()
	kl.StopCleanLoop()
}

func TestKeyLock_StopCleanLoopMultipleTimes(t *testing.T) {
	kl := NewKeyLock()
	kl.StartCleanLoop()

	// Should not panic when called multiple times
	kl.StopCleanLoop()
	kl.StopCleanLoop()
	kl.StopCleanLoop()
}

func TestKeyLock_StopCleanLoopWithoutStart(t *testing.T) {
	kl := NewKeyLock()

	// Should not panic when called without StartCleanLoop
	kl.StopCleanLoop()
	kl.StopCleanLoop()
}

func TestKeyLock_StopCleanLoopConcurrent(t *testing.T) {
	kl := NewKeyLock()
	kl.StartCleanLoop()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			kl.StopCleanLoop()
		}()
	}
	wg.Wait()
}

func TestKeyLock_UnlockNonExistentKey(t *testing.T) {
	kl := NewKeyLock()

	// Should not panic
	kl.Unlock("nonexistent")
}
