package toychacha

import (
	"crypto/cipher"
	"crypto/subtle"
	"encoding/binary"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
)

func numTo8LeBytes(l int) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(l))
	return b
}

// paddedSize returns the size padded to an integral multiple of 16
func paddedSize(d []byte) int {
	c := len(d) % 16
	if c == 0 {
		return len(d)
	}
	return len(d) + 16 - c
}

func constructMacData(macdata, aad, ciphertext []byte) {
	aadsize := paddedSize(aad)
	ciphertextsize := paddedSize(ciphertext)

	// macData := make([]byte, aadsize+ciphertextsize+8+8)
	header := macdata

	copy(header[:aadsize], aad)
	header = header[aadsize:]

	copy(header[:ciphertextsize], ciphertext)
	header = header[ciphertextsize:]

	copy(header[:8], numTo8LeBytes(len(aad)))
	header = header[8:]
	copy(header[:8], numTo8LeBytes(len(ciphertext)))
}

func aeadEncrpt(aad, key, iv, constant, plaintext []byte) ([]byte, []byte) {
	nonce := append(constant, iv...)
	otk := genMacKey(key, nonce)

	ciphertext := encrypt(key, nonce, plaintext, 1)

	aadsize := paddedSize(aad)
	ciphertextsize := paddedSize(ciphertext)

	macData := make([]byte, aadsize+ciphertextsize+8+8)

	constructMacData(macData, aad, ciphertext)
	tag := mac(macData, otk)
	return ciphertext, tag
}

func aeadDecrypt(aad, key, iv, constant, ciphertext, tag []byte) ([]byte, error) {
	nonce := append(constant, iv...)
	otk := genMacKey(key, nonce)

	aadsize := paddedSize(aad)
	ciphertextsize := paddedSize(ciphertext)

	macData := make([]byte, aadsize+ciphertextsize+8+8)

	constructMacData(macData, aad, ciphertext)

	calcTag := mac(macData, otk)
	if subtle.ConstantTimeCompare(calcTag, tag) == 0 {
		return nil, errors.New("invalid tag")
	}

	return encrypt(key, nonce, ciphertext, 1), nil
}

func New(key []byte) (cipher.AEAD, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, errors.New("chacha20poly1305: bad key length")
	}
	return &toyChacha20Poly1305{key: key}, nil
}

func NewToyChacha20Poly1305(key []byte) (*toyChacha20Poly1305, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, errors.New("chacha20poly1305: bad key length")
	}
	return &toyChacha20Poly1305{key: key}, nil
}

type toyChacha20Poly1305 struct {
	key []byte
}

// NonceSize implements cipher.AEAD
func (*toyChacha20Poly1305) NonceSize() int { return chacha20poly1305.NonceSize }

// Open implements cipher.AEAD
func (tc *toyChacha20Poly1305) Open(dst []byte, nonce []byte, ciphertext []byte, additionalData []byte) ([]byte, error) {
	return aeadDecrypt(additionalData, tc.key, nonce, nil, ciphertext[:len(ciphertext)-16], ciphertext[len(ciphertext)-16:])
}

// Overhead implements cipher.AEAD
func (*toyChacha20Poly1305) Overhead() int { return chacha20poly1305.Overhead }

// Seal implements cipher.AEAD
func (tc *toyChacha20Poly1305) Seal(dst []byte, nonce []byte, plaintext []byte, additionalData []byte) []byte {
	if len(dst) < len(plaintext)+16 {
		dst = make([]byte, len(plaintext)+16)
	}
	e, a := aeadEncrpt(additionalData, tc.key, nonce, nil, plaintext)
	copy(dst, e)
	copy(dst[len(dst)-16:], a)
	return dst
}
