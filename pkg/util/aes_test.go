package util

import (
	"bytes"
	"testing"
)

func TestPKCS5UnPadding(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "empty input",
			input:   []byte{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "valid padding 1",
			input:   []byte{0x01},
			want:    []byte{},
			wantErr: false,
		},
		{
			name:    "valid padding 3",
			input:   []byte{'h', 'e', 'l', 'l', 'o', 0x03, 0x03, 0x03},
			want:    []byte{'h', 'e', 'l', 'l', 'o'},
			wantErr: false,
		},
		{
			name:    "padding value zero",
			input:   []byte{'h', 'e', 'l', 'l', 'o', 0x00},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "padding value larger than length",
			input:   []byte{'h', 'e', 'l', 'l', 'o', 0x10},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "inconsistent padding bytes",
			input:   []byte{'h', 'e', 'l', 'l', 'o', 0x02, 0x03, 0x03},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PKCS5UnPadding(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PKCS5UnPadding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("PKCS5UnPadding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPKCS7UnPadding(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "empty input",
			input:   []byte{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "valid padding 1",
			input:   []byte{0x01},
			want:    []byte{},
			wantErr: false,
		},
		{
			name:    "valid padding 4",
			input:   []byte{'t', 'e', 's', 't', 0x04, 0x04, 0x04, 0x04},
			want:    []byte{'t', 'e', 's', 't'},
			wantErr: false,
		},
		{
			name:    "padding value zero",
			input:   []byte{'t', 'e', 's', 't', 0x00},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "padding value larger than length",
			input:   []byte{'t', 'e', 's', 't', 0x20},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "inconsistent padding bytes",
			input:   []byte{'t', 'e', 's', 't', 0x01, 0x04, 0x04, 0x04},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "full block padding",
			input:   []byte{0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10},
			want:    []byte{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PKCS7UnPadding(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PKCS7UnPadding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("PKCS7UnPadding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAesEncryptDecryptPkcs5(t *testing.T) {
	key := []byte("1234567890123456") // 16 bytes for AES-128
	iv := []byte("1234567890123456")

	tests := []struct {
		name      string
		plaintext []byte
	}{
		{"empty", []byte{}},
		{"short", []byte("hello")},
		{"exact block", []byte("1234567890123456")},
		{"long", []byte("this is a longer test message for encryption")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := AesEncryptPkcs5(tt.plaintext, key, iv)
			if err != nil {
				t.Fatalf("AesEncryptPkcs5() error = %v", err)
			}

			decrypted, err := AesDecryptPkcs5(encrypted, key, iv)
			if err != nil {
				t.Fatalf("AesDecryptPkcs5() error = %v", err)
			}

			if !bytes.Equal(decrypted, tt.plaintext) {
				t.Errorf("roundtrip failed: got %v, want %v", decrypted, tt.plaintext)
			}
		})
	}
}

func TestAesEncryptDecryptPkcs7(t *testing.T) {
	key := []byte("1234567890123456") // 16 bytes for AES-128
	iv := []byte("1234567890123456")

	tests := []struct {
		name      string
		plaintext []byte
	}{
		{"empty", []byte{}},
		{"short", []byte("hello")},
		{"exact block", []byte("1234567890123456")},
		{"long", []byte("this is a longer test message for encryption")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := AesEncryptPkcs7(tt.plaintext, key, iv)
			if err != nil {
				t.Fatalf("AesEncryptPkcs7() error = %v", err)
			}

			decrypted, err := AesDecryptPkcs7(encrypted, key, iv)
			if err != nil {
				t.Fatalf("AesDecryptPkcs7() error = %v", err)
			}

			if !bytes.Equal(decrypted, tt.plaintext) {
				t.Errorf("roundtrip failed: got %v, want %v", decrypted, tt.plaintext)
			}
		})
	}
}

func TestAesDecryptWithInvalidPadding(t *testing.T) {
	key := []byte("1234567890123456")
	iv := []byte("1234567890123456")

	// Create ciphertext with invalid padding
	encrypted, err := AesEncryptPkcs7([]byte("test"), key, iv)
	if err != nil {
		t.Fatalf("AesEncryptPkcs7() error = %v", err)
	}

	// Corrupt the last byte
	encrypted[len(encrypted)-1] ^= 0xFF

	_, err = AesDecryptPkcs7(encrypted, key, iv)
	if err == nil {
		t.Error("expected error for corrupted ciphertext, got nil")
	}
}
