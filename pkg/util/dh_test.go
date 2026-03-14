package util

import (
	"testing"
)

func TestGetCurve25519KeypPair(t *testing.T) {
	priv, pub, err := GetCurve25519KeypPair()
	if err != nil {
		t.Fatalf("GetCurve25519KeypPair() returned error: %v", err)
	}

	// Check that keys are not all zeros
	var zeroKey [32]byte
	if priv == zeroKey {
		t.Error("GetCurve25519KeypPair() returned zero private key")
	}
	if pub == zeroKey {
		t.Error("GetCurve25519KeypPair() returned zero public key")
	}

	// Verify keys are different
	if priv == pub {
		t.Error("GetCurve25519KeypPair() returned identical private and public keys")
	}
}

func TestGetCurve25519KeypPairUniqueness(t *testing.T) {
	// Generate multiple key pairs and verify they're unique
	pairs := make(map[[32]byte]bool)
	for i := 0; i < 10; i++ {
		priv, _, err := GetCurve25519KeypPair()
		if err != nil {
			t.Fatalf("GetCurve25519KeypPair() returned error on iteration %d: %v", i, err)
		}
		if pairs[priv] {
			t.Errorf("GetCurve25519KeypPair() returned duplicate private key on iteration %d", i)
		}
		pairs[priv] = true
	}
}

func TestGetCurve25519Key(t *testing.T) {
	priv, pub, err := GetCurve25519KeypPair()
	if err != nil {
		t.Fatalf("GetCurve25519KeypPair() returned error: %v", err)
	}

	key := GetCurve25519Key(priv, pub)

	// Check that key is not all zeros
	var zeroKey [32]byte
	if key == zeroKey {
		t.Error("GetCurve25519Key() returned zero key")
	}
}
