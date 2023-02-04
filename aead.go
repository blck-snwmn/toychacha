package gochacha

import (
	"crypto/subtle"
	"encoding/binary"
	"errors"
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

func constructMacData(aad, ciphertext []byte) []byte {
	aadsize := paddedSize(aad)
	ciphertextsize := paddedSize(ciphertext)

	macData := make([]byte, aadsize+ciphertextsize+8+8)
	header := macData

	copy(header[:aadsize], aad)
	header = header[aadsize:]

	copy(header[:ciphertextsize], ciphertext)
	header = header[ciphertextsize:]

	copy(header[:8], numTo8LeBytes(len(aad)))
	header = header[8:]
	copy(header[:8], numTo8LeBytes(len(ciphertext)))
	return macData
}

func AeadEncrpt(aad, key, iv, constant, plaintext []byte) ([]byte, []byte) {
	return aeadEncrpt(aad, key, iv, constant, plaintext)
}

func aeadEncrpt(aad, key, iv, constant, plaintext []byte) ([]byte, []byte) {
	nonce := append(constant, iv...)
	otk := genMacKey(key, nonce)

	ciphertext := encrypt(key, nonce, plaintext, 1)

	macData := constructMacData(aad, ciphertext)
	tag := mac(macData, otk)
	return ciphertext, tag
}

func aeadDecrypt(aad, key, iv, constant, ciphertext, tag []byte) ([]byte, error) {
	nonce := append(constant, iv...)
	otk := genMacKey(key, nonce)

	macData := constructMacData(aad, ciphertext)

	calcTag := mac(macData, otk)
	if subtle.ConstantTimeCompare(calcTag, tag) == 0 {
		return nil, errors.New("invalid tag")
	}

	return encrypt(key, nonce, ciphertext, 1), nil
}
