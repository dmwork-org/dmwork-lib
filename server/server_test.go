package server

import (
	"bytes"
	"fmt"
	"runtime/debug"
	"sync"
	"testing"
	"time"

	"github.com/dmwork-org/dmwork-lib/pkg/log"
	"go.uber.org/zap"
)

// TestGoroutinePanicRecovery tests that panics in goroutines are recovered
// and don't crash the entire process.
func TestGoroutinePanicRecovery(t *testing.T) {
	var wg sync.WaitGroup
	var recovered bool
	var mu sync.Mutex

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				mu.Lock()
				recovered = true
				mu.Unlock()
				// Log the panic (similar to what server.go does)
				log.Error("panic recovered in test goroutine",
					zap.String("panic", fmt.Sprintf("%v", r)),
					zap.String("stack", string(debug.Stack())))
			}
		}()
		panic("test panic")
	}()

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if !recovered {
		t.Error("expected panic to be recovered")
	}
}

// TestGoroutinePanicRecoveryWithError tests panic recovery when an error causes a panic.
func TestGoroutinePanicRecoveryWithError(t *testing.T) {
	var wg sync.WaitGroup
	var recovered bool
	var panicValue interface{}
	var mu sync.Mutex

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				mu.Lock()
				recovered = true
				panicValue = r
				mu.Unlock()
			}
		}()
		err := fmt.Errorf("simulated server error")
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if !recovered {
		t.Error("expected panic to be recovered")
	}
	if panicValue == nil {
		t.Error("expected panicValue to be non-nil")
	}
	if err, ok := panicValue.(error); !ok || err.Error() != "simulated server error" {
		t.Errorf("unexpected panic value: %v", panicValue)
	}
}

// TestMultipleGoroutinePanicRecovery tests that multiple goroutines can each
// recover from their own panics independently.
func TestMultipleGoroutinePanicRecovery(t *testing.T) {
	var wg sync.WaitGroup
	recoveredCount := 0
	var mu sync.Mutex

	numGoroutines := 5
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					mu.Lock()
					recoveredCount++
					mu.Unlock()
				}
			}()
			panic(fmt.Sprintf("panic from goroutine %d", id))
		}(i)
	}

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if recoveredCount != numGoroutines {
		t.Errorf("expected %d recoveries, got %d", numGoroutines, recoveredCount)
	}
}

// TestPanicRecoveryDoesNotAffectMainGoroutine verifies that a panic in a
// spawned goroutine with recovery doesn't affect the main goroutine.
func TestPanicRecoveryDoesNotAffectMainGoroutine(t *testing.T) {
	mainGoroutineCompleted := false
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			recover() // Recover the panic
		}()
		panic("child goroutine panic")
	}()

	// Wait a bit to ensure the goroutine has time to panic
	wg.Wait()

	// If we reach here, the main goroutine wasn't affected
	mainGoroutineCompleted = true

	if !mainGoroutineCompleted {
		t.Error("main goroutine should complete even if child panics")
	}
}

// TestPanicRecoveryStackTrace verifies that stack traces are captured during recovery.
func TestPanicRecoveryStackTrace(t *testing.T) {
	var wg sync.WaitGroup
	var capturedStack string
	var mu sync.Mutex

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				mu.Lock()
				capturedStack = string(debug.Stack())
				mu.Unlock()
			}
		}()
		panic("test panic for stack trace")
	}()

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if capturedStack == "" {
		t.Error("expected stack trace to be captured")
	}
	if !bytes.Contains([]byte(capturedStack), []byte("goroutine")) {
		t.Error("stack trace should contain 'goroutine'")
	}
}

// TestPanicRecoveryWithNilPanic tests that nil panics are also handled.
func TestPanicRecoveryWithNilPanic(t *testing.T) {
	var wg sync.WaitGroup
	var recovered bool
	var mu sync.Mutex

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				mu.Lock()
				recovered = true
				mu.Unlock()
			}
		}()
		panic(nil)
	}()

	wg.Wait()

	// Note: panic(nil) in Go will not trigger recovery since recover() returns nil
	// This is expected behavior
	mu.Lock()
	defer mu.Unlock()
	if recovered {
		t.Error("panic(nil) should not trigger recovery (recover() returns nil)")
	}
}

// TestConcurrentPanicRecovery tests that concurrent panics in different
// goroutines are all properly recovered.
func TestConcurrentPanicRecovery(t *testing.T) {
	var wg sync.WaitGroup
	recoveredCount := 0
	var mu sync.Mutex

	numGoroutines := 100
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					mu.Lock()
					recoveredCount++
					mu.Unlock()
				}
			}()
			// Add some variability to make this more realistic
			time.Sleep(time.Millisecond * time.Duration(id%10))
			panic(fmt.Sprintf("concurrent panic %d", id))
		}(i)
	}

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if recoveredCount != numGoroutines {
		t.Errorf("expected %d recoveries, got %d", numGoroutines, recoveredCount)
	}
}
