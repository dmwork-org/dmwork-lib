package register

import (
	"sync"
	"testing"
)

func TestGetModulesConcurrent(t *testing.T) {
	// Reset state for testing
	once = sync.Once{}
	moduleList = nil
	modules = nil

	// Add a test module
	AddModule(func(ctx interface{}) Module {
		return Module{
			Name:    "test-module",
			Service: "test-service",
		}
	})

	var wg sync.WaitGroup
	const goroutines = 100

	// Launch multiple goroutines calling GetModules concurrently
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mods := GetModules(nil)
			if len(mods) != 1 {
				t.Errorf("expected 1 module, got %d", len(mods))
			}
		}()
	}

	wg.Wait()
}

func TestGetModuleByNameConcurrent(t *testing.T) {
	// Reset state for testing
	once = sync.Once{}
	moduleList = nil
	modules = nil

	// Add test modules
	AddModule(func(ctx interface{}) Module {
		return Module{
			Name:    "module-a",
			Service: "service-a",
		}
	})
	AddModule(func(ctx interface{}) Module {
		return Module{
			Name:    "module-b",
			Service: "service-b",
		}
	})

	var wg sync.WaitGroup
	const goroutines = 100

	// Launch multiple goroutines calling GetModuleByName and GetModules concurrently
	for i := 0; i < goroutines; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			m := GetModuleByName("module-a", nil)
			if m.Name != "module-a" {
				t.Errorf("expected module-a, got %s", m.Name)
			}
		}()
		go func() {
			defer wg.Done()
			GetModules(nil)
		}()
	}

	wg.Wait()
}

func TestGetServiceConcurrent(t *testing.T) {
	// Reset state for testing
	once = sync.Once{}
	moduleList = nil
	modules = nil

	// Add a test module
	AddModule(func(ctx interface{}) Module {
		return Module{
			Name:    "service-module",
			Service: "my-service",
		}
	})

	// Initialize modules first
	GetModules(nil)

	var wg sync.WaitGroup
	const goroutines = 100

	// Launch multiple goroutines calling GetService concurrently
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			svc := GetService("service-module")
			if svc != "my-service" {
				t.Errorf("expected my-service, got %v", svc)
			}
		}()
	}

	wg.Wait()
}

func TestGetModuleByNameNotFound(t *testing.T) {
	// Reset state for testing
	once = sync.Once{}
	moduleList = nil
	modules = nil

	AddModule(func(ctx interface{}) Module {
		return Module{Name: "existing"}
	})

	m := GetModuleByName("nonexistent", nil)
	if m.Name != "" {
		t.Errorf("expected empty module, got %s", m.Name)
	}
}

func TestGetServiceNotFound(t *testing.T) {
	// Reset state for testing
	once = sync.Once{}
	moduleList = nil
	modules = nil

	AddModule(func(ctx interface{}) Module {
		return Module{Name: "existing", Service: "svc"}
	})
	GetModules(nil)

	svc := GetService("nonexistent")
	if svc != nil {
		t.Errorf("expected nil, got %v", svc)
	}
}
