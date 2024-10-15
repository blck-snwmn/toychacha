package toychacha

import (
	"crypto/cipher"
	"crypto/subtle"
	"encoding/binary"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
)

func numTo8LeBytes(in []byte) []byte {
	l := len(in)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(l))
	return b
}

// paddedSize returns the size padded to an integral multiple of 16
func paddedSize(d []byte) int {
	return ((len(d) + 15) / 16) * 16
}

func constructMacData(macdata, aad, ciphertext []byte) {
	aadsize := paddedSize(aad)
	ciphertextsize := paddedSize(ciphertext)

	// macData := make([]byte, aadsize+ciphertextsize+8+8)
	head := macdata

	copy(head[:aadsize], aad)
	head = head[aadsize:]

	copy(head[:ciphertextsize], ciphertext)
	head = head[ciphertextsize:]

	copy(head[:8], numTo8LeBytes(aad))
	head = head[8:]
	copy(head[:8], numTo8LeBytes(ciphertext))
}

func aeadEncrpt(dst, aad, key, iv, constant, plaintext []byte) ([]byte, []byte) {
	nonce := append(constant, iv...)
	otk := genMacKey(key, nonce)

	ciphertext := make([]byte, len(plaintext))
	ciphertext = encrypt(ciphertext, key, nonce, plaintext, 1)

	aadsize := paddedSize(aad)
	ciphertextsize := paddedSize(ciphertext)

	macData := make([]byte, aadsize+ciphertextsize+8+8)

	constructMacData(macData, aad, ciphertext)
	tag := mac(macData, otk)
	return ciphertext, tag[:]
}

func aeadDecrypt(dst, aad, key, iv, constant, ciphertext, tag []byte) ([]byte, error) {
	nonce := append(constant, iv...)
	otk := genMacKey(key, nonce)

	aadsize := paddedSize(aad)
	ciphertextsize := paddedSize(ciphertext)

	macData := make([]byte, aadsize+ciphertextsize+8+8)

	constructMacData(macData, aad, ciphertext)

	calcTag := mac(macData, otk)
	if subtle.ConstantTimeCompare(calcTag[:], tag) == 0 {
		return nil, errors.New("invalid tag")
	}
	if len(dst) < len(ciphertext) {
		dst = make([]byte, len(ciphertext))
	}
	return encrypt(dst, key, nonce, ciphertext, 1), nil
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
	return aeadDecrypt(dst, additionalData, tc.key, nonce, nil, ciphertext[:len(ciphertext)-tc.Overhead()], ciphertext[len(ciphertext)-tc.Overhead():])
}

// Overhead implements cipher.AEAD
func (*toyChacha20Poly1305) Overhead() int { return chacha20poly1305.Overhead }

// Seal implements cipher.AEAD
func (tc *toyChacha20Poly1305) Seal(dst []byte, nonce []byte, plaintext []byte, additionalData []byte) []byte {
	if len(dst) < len(plaintext)+tc.Overhead() {
		dst = make([]byte, len(plaintext)+tc.Overhead())
	}
	e, a := aeadEncrpt(dst, additionalData, tc.key, nonce, nil, plaintext)
	copy(dst, e)
	copy(dst[len(dst)-tc.Overhead():], a)
	return dst
}
