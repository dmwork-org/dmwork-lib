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
			name:    "valid padding 1 byte",
			input:   []byte{0x01},
			want:    []byte{},
			wantErr: false,
		},
		{
			name:    "valid padding 3 bytes",
			input:   []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x03, 0x03, 0x03},
			want:    []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
			wantErr: false,
		},
		{
			name:    "invalid padding - zero",
			input:   []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid padding - exceeds length",
			input:   []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x10},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid padding - inconsistent bytes",
			input:   []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x02, 0x03, 0x03},
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
			name:    "valid padding 1 byte",
			input:   []byte{0x01},
			want:    []byte{},
			wantErr: false,
		},
		{
			name:    "valid padding 4 bytes",
			input:   []byte{0x68, 0x65, 0x6c, 0x6c, 0x04, 0x04, 0x04, 0x04},
			want:    []byte{0x68, 0x65, 0x6c, 0x6c},
			wantErr: false,
		},
		{
			name:    "invalid padding - zero",
			input:   []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid padding - exceeds length",
			input:   []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid padding - inconsistent bytes",
			input:   []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x01, 0x03, 0x03},
			want:    nil,
			wantErr: true,
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

func TestAesEncryptDecryptRoundTrip(t *testing.T) {
	key := []byte("1234567890123456") // 16 bytes for AES-128
	iv := []byte("1234567890123456")

	tests := []struct {
		name      string
		plaintext []byte
	}{
		{
			name:      "simple text",
			plaintext: []byte("hello world"),
		},
		{
			name:      "exact block size",
			plaintext: []byte("1234567890123456"),
		},
		{
			name:      "multiple blocks",
			plaintext: []byte("this is a longer message that spans multiple blocks"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" PKCS5", func(t *testing.T) {
			encrypted, err := AesEncryptPkcs5(tt.plaintext, key, iv)
			if err != nil {
				t.Fatalf("AesEncryptPkcs5() error = %v", err)
			}
			decrypted, err := AesDecryptPkcs5(encrypted, key, iv)
			if err != nil {
				t.Fatalf("AesDecryptPkcs5() error = %v", err)
			}
			if !bytes.Equal(decrypted, tt.plaintext) {
				t.Errorf("round trip failed: got %v, want %v", decrypted, tt.plaintext)
			}
		})

		t.Run(tt.name+" PKCS7", func(t *testing.T) {
			encrypted, err := AesEncryptPkcs7(tt.plaintext, key, iv)
			if err != nil {
				t.Fatalf("AesEncryptPkcs7() error = %v", err)
			}
			decrypted, err := AesDecryptPkcs7(encrypted, key, iv)
			if err != nil {
				t.Fatalf("AesDecryptPkcs7() error = %v", err)
			}
			if !bytes.Equal(decrypted, tt.plaintext) {
				t.Errorf("round trip failed: got %v, want %v", decrypted, tt.plaintext)
			}
		})
	}
}
