package main

import (
	"bytes"
	"testing"

	"github.com/blck-snwmn/toychacha"
	"golang.org/x/crypto/chacha20poly1305"
)

func TestSealMessageRoundTrip(t *testing.T) {
	key := make([]byte, chacha20poly1305.KeySize)
	aead, err := toychacha.New(key)
	if err != nil {
		t.Fatal(err)
	}
	standard, err := chacha20poly1305.New(key)
	if err != nil {
		t.Fatal(err)
	}
	plaintext := []byte("message with its own nonce")
	aad := []byte("associated data")

	message, err := sealMessage(aead, plaintext, aad)
	if err != nil {
		t.Fatal(err)
	}
	if want := aead.NonceSize() + len(plaintext) + aead.Overhead(); len(message) != want {
		t.Fatalf("message length = %d, want %d", len(message), want)
	}
	got, err := openMessage(aead, message, aad)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Fatalf("opened plaintext = %q, want %q", got, plaintext)
	}

	// The serialized nonce and ciphertext also work with the standard AEAD.
	got, err = standard.Open(nil, message[:aead.NonceSize()], message[aead.NonceSize():], aad)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Fatalf("standard AEAD opened plaintext = %q, want %q", got, plaintext)
	}
}

func TestOpenMessageRejectsInvalidInput(t *testing.T) {
	aead, err := toychacha.New(make([]byte, chacha20poly1305.KeySize))
	if err != nil {
		t.Fatal(err)
	}
	message, err := sealMessage(aead, []byte("plaintext"), []byte("aad"))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		message []byte
		aad     []byte
	}{
		{name: "truncated", message: message[:aead.NonceSize()+aead.Overhead()-1], aad: []byte("aad")},
		{name: "modified nonce", message: append([]byte(nil), message...), aad: []byte("aad")},
		{name: "modified aad", message: message, aad: []byte("other")},
	}
	tests[1].message[0] ^= 1
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := openMessage(aead, tt.message, tt.aad); err == nil {
				t.Fatal("openMessage accepted invalid input")
			}
		})
	}
}
