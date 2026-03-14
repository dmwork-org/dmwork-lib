package util

import (
	"testing"
)

func TestGetRandomName(t *testing.T) {
	// Test that all names can be selected (including the last one)
	// Run multiple iterations to statistically verify the last element is reachable
	lastElement := names[len(names)-1]
	foundLast := false

	// Run enough iterations to have high probability of hitting the last element
	for i := 0; i < 10000; i++ {
		name := GetRandomName()
		if name == lastElement {
			foundLast = true
			break
		}
	}

	if !foundLast {
		t.Errorf("GetRandomName() never returned the last element %q after 10000 iterations", lastElement)
	}
}

func TestGetRandomNameReturnsValidName(t *testing.T) {
	for i := 0; i < 100; i++ {
		name := GetRandomName()
		found := false
		for _, n := range names {
			if n == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetRandomName() returned invalid name: %q", name)
		}
	}
}

func TestGenerUUID(t *testing.T) {
	uuid, err := GenerUUID()
	if err != nil {
		t.Fatalf("GenerUUID() returned error: %v", err)
	}

	// UUID without dashes should be 32 characters
	if len(uuid) != 32 {
		t.Errorf("GenerUUID() returned UUID with length %d, expected 32", len(uuid))
	}

	// Should not contain dashes
	for _, c := range uuid {
		if c == '-' {
			t.Error("GenerUUID() returned UUID containing dash")
			break
		}
	}
}

func TestGenerUUIDUniqueness(t *testing.T) {
	uuids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		uuid, err := GenerUUID()
		if err != nil {
			t.Fatalf("GenerUUID() returned error on iteration %d: %v", i, err)
		}
		if uuids[uuid] {
			t.Errorf("GenerUUID() returned duplicate UUID on iteration %d", i)
		}
		uuids[uuid] = true
	}
}
