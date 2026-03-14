package util

import (
	"testing"
)

func TestNewV1(t *testing.T) {
	uuid, err := NewV1()
	if err != nil {
		t.Fatalf("NewV1() returned error: %v", err)
	}

	if uuid == Nil {
		t.Error("NewV1() returned nil UUID")
	}

	if uuid.Version() != 1 {
		t.Errorf("NewV1() returned UUID with version %d, expected 1", uuid.Version())
	}

	if uuid.Variant() != VariantRFC4122 {
		t.Errorf("NewV1() returned UUID with variant %d, expected RFC4122", uuid.Variant())
	}
}

func TestNewV2(t *testing.T) {
	uuid, err := NewV2(DomainPerson)
	if err != nil {
		t.Fatalf("NewV2(DomainPerson) returned error: %v", err)
	}

	if uuid == Nil {
		t.Error("NewV2() returned nil UUID")
	}

	if uuid.Version() != 2 {
		t.Errorf("NewV2() returned UUID with version %d, expected 2", uuid.Version())
	}

	if uuid.Variant() != VariantRFC4122 {
		t.Errorf("NewV2() returned UUID with variant %d, expected RFC4122", uuid.Variant())
	}
}

func TestNewV4(t *testing.T) {
	uuid, err := NewV4()
	if err != nil {
		t.Fatalf("NewV4() returned error: %v", err)
	}

	if uuid == Nil {
		t.Error("NewV4() returned nil UUID")
	}

	if uuid.Version() != 4 {
		t.Errorf("NewV4() returned UUID with version %d, expected 4", uuid.Version())
	}

	if uuid.Variant() != VariantRFC4122 {
		t.Errorf("NewV4() returned UUID with variant %d, expected RFC4122", uuid.Variant())
	}
}

func TestNewV4Uniqueness(t *testing.T) {
	uuids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		uuid, err := NewV4()
		if err != nil {
			t.Fatalf("NewV4() returned error on iteration %d: %v", i, err)
		}
		s := uuid.String()
		if uuids[s] {
			t.Errorf("NewV4() returned duplicate UUID on iteration %d", i)
		}
		uuids[s] = true
	}
}

func TestNewV3(t *testing.T) {
	uuid := NewV3(NamespaceDNS, "example.com")

	if uuid == Nil {
		t.Error("NewV3() returned nil UUID")
	}

	if uuid.Version() != 3 {
		t.Errorf("NewV3() returned UUID with version %d, expected 3", uuid.Version())
	}

	// NewV3 should be deterministic
	uuid2 := NewV3(NamespaceDNS, "example.com")
	if uuid != uuid2 {
		t.Error("NewV3() is not deterministic")
	}
}

func TestNewV5(t *testing.T) {
	uuid := NewV5(NamespaceDNS, "example.com")

	if uuid == Nil {
		t.Error("NewV5() returned nil UUID")
	}

	if uuid.Version() != 5 {
		t.Errorf("NewV5() returned UUID with version %d, expected 5", uuid.Version())
	}

	// NewV5 should be deterministic
	uuid2 := NewV5(NamespaceDNS, "example.com")
	if uuid != uuid2 {
		t.Error("NewV5() is not deterministic")
	}
}
