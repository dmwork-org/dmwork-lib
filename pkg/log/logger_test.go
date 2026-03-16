package log

import (
	"sync"
	"testing"

	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	resetForTesting()

	opts := NewOptions()
	opts.Level = zap.DebugLevel
	opts.LineNum = true
	Configure(opts)

	Info("this is info")
	Debug("this is debug")
	Error("this is error", zap.String("key", "value"))
}

func TestLoggerConcurrent(t *testing.T) {
	resetForTesting()

	var wg sync.WaitGroup
	numGoroutines := 100

	// Start many goroutines that all try to log concurrently
	// This tests that the sync.Once properly handles concurrent access
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			Info("concurrent info", zap.Int("goroutine", id))
			Debug("concurrent debug", zap.Int("goroutine", id))
			Error("concurrent error", zap.Int("goroutine", id))
			Warn("concurrent warn", zap.Int("goroutine", id))
		}(i)
	}

	wg.Wait()
}

func TestLoggerConcurrentConfigure(t *testing.T) {
	resetForTesting()

	var wg sync.WaitGroup
	numGoroutines := 50

	// Start many goroutines that all try to configure and log concurrently
	// This tests the race condition where multiple goroutines try to Configure()
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Some goroutines try to configure
			if id%2 == 0 {
				opts := NewOptions()
				opts.Level = zap.DebugLevel
				Configure(opts)
			}
			// All goroutines try to log
			Info("concurrent configure test", zap.Int("goroutine", id))
		}(i)
	}

	wg.Wait()
}

func TestLoggerEnsureConfigured(t *testing.T) {
	resetForTesting()

	// Call logging functions without explicit Configure()
	// This tests that ensureConfigured() properly initializes the logger
	Info("auto configured info")
	Debug("auto configured debug")
	Error("auto configured error")
	Warn("auto configured warn")
}

func TestTLogConcurrent(t *testing.T) {
	resetForTesting()

	var wg sync.WaitGroup
	numGoroutines := 50

	tlog := NewTLog("test-prefix")

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			tlog.Info("tlog concurrent info", zap.Int("goroutine", id))
			tlog.Debug("tlog concurrent debug", zap.Int("goroutine", id))
			tlog.Error("tlog concurrent error", zap.Int("goroutine", id))
			tlog.Warn("tlog concurrent warn", zap.Int("goroutine", id))
		}(i)
	}

	wg.Wait()
}
