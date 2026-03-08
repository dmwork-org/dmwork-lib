package util

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/curve25519"
)

// GetCurve25519KeypPair generates a Curve25519 key pair.
func GetCurve25519KeypPair() (Aprivate, Apublic [32]byte) {
	if _, err := io.ReadFull(rand.Reader, Aprivate[:]); err != nil {
		panic(err)
	}
	pub, err := curve25519.X25519(Aprivate[:], curve25519.Basepoint)
	if err != nil {
		panic(err)
	}
	copy(Apublic[:], pub)
	return
}

// GetCurve25519Key computes a shared secret from a private and public key.
func GetCurve25519Key(private, public [32]byte) (Key [32]byte) {
	shared, err := curve25519.X25519(private[:], public[:])
	if err != nil {
		panic(err)
	}
	copy(Key[:], shared)
	return
}
