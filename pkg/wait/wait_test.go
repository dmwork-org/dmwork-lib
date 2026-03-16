package wait

import (
	"errors"
	"testing"
)

func TestRegister(t *testing.T) {
	w := New()
	id := uint64(123)

	ch, err := w.Register(id)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if ch == nil {
		t.Fatal("Register returned nil channel")
	}
	if !w.IsRegistered(id) {
		t.Fatal("ID should be registered")
	}
}

func TestRegisterDuplicateID(t *testing.T) {
	w := New()
	id := uint64(456)

	_, err := w.Register(id)
	if err != nil {
		t.Fatalf("First Register failed: %v", err)
	}

	ch, err := w.Register(id)
	if err == nil {
		t.Fatal("Expected error for duplicate ID registration")
	}
	if ch != nil {
		t.Fatal("Expected nil channel for duplicate ID registration")
	}
	if !errors.Is(err, ErrDuplicateID) {
		t.Fatalf("Expected ErrDuplicateID, got: %v", err)
	}
}

func TestRegisterDuplicateIDMultiple(t *testing.T) {
	w := New()
	ids := []uint64{1, 2, 3, 100, 1000}

	for _, id := range ids {
		_, err := w.Register(id)
		if err != nil {
			t.Fatalf("Register(%d) failed: %v", id, err)
		}
	}

	for _, id := range ids {
		_, err := w.Register(id)
		if err == nil {
			t.Fatalf("Expected error for duplicate ID %d", id)
		}
		if !errors.Is(err, ErrDuplicateID) {
			t.Fatalf("Expected ErrDuplicateID for ID %d, got: %v", id, err)
		}
	}
}

func TestTrigger(t *testing.T) {
	w := New()
	id := uint64(789)
	expected := "test value"

	ch, err := w.Register(id)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	go func() {
		w.Trigger(id, expected)
	}()

	result := <-ch
	if result != expected {
		t.Fatalf("Expected %v, got %v", expected, result)
	}
}

func TestTriggerUnregisters(t *testing.T) {
	w := New()
	id := uint64(999)

	_, err := w.Register(id)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if !w.IsRegistered(id) {
		t.Fatal("ID should be registered before trigger")
	}

	w.Trigger(id, nil)

	if w.IsRegistered(id) {
		t.Fatal("ID should not be registered after trigger")
	}
}

func TestRegisterAfterTrigger(t *testing.T) {
	w := New()
	id := uint64(111)

	ch1, err := w.Register(id)
	if err != nil {
		t.Fatalf("First Register failed: %v", err)
	}

	w.Trigger(id, "first")
	<-ch1

	ch2, err := w.Register(id)
	if err != nil {
		t.Fatalf("Second Register after Trigger failed: %v", err)
	}
	if ch2 == nil {
		t.Fatal("Expected non-nil channel for re-registration after trigger")
	}
}

func TestIsRegistered(t *testing.T) {
	w := New()
	id := uint64(222)

	if w.IsRegistered(id) {
		t.Fatal("ID should not be registered initially")
	}

	_, err := w.Register(id)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if !w.IsRegistered(id) {
		t.Fatal("ID should be registered after Register")
	}
}
