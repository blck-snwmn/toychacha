package main

import (
	"fmt"
	"os"

	gochacha "github.com/blck-snwmn/go-chacha"
	"golang.org/x/crypto/chacha20poly1305"
)

func main() {
	if len(os.Args) <= 1 {
		return
	}
	plaintext := os.Args[1]

	aad := []byte{
		0x50, 0x51, 0x52, 0x53,
		0xc0, 0xc1, 0xc2, 0xc3,
		0xc4, 0xc5, 0xc6, 0xc7,
	}
	key := []byte{
		0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87,
		0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
		0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97,
		0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
	}

	nonce := []byte{
		0x07, 0x00, 0x00, 0x00, 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47,
	}

	tcp, _ := gochacha.NewToyChacha20Poly1305(key)
	aead := tcp.Seal(nil, nonce, []byte(plaintext), aad)
	fmt.Printf("AEAD=%x\n", aead)

	cipher, _ := chacha20poly1305.New(key)
	ciphertext := cipher.Seal(nil, nonce, []byte(plaintext), aad)
	fmt.Printf("AEAD=%x\n", ciphertext)
}
