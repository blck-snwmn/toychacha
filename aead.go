package gochacha

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
)

func pad16(x []byte) []byte {
	if len(x)%16 == 0 {
		return nil
	}
	return bytes.Repeat([]byte{0}, 16-(len(x)%16))
}

func numTo8LeBytes(l int) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(l))
	return b
}

func size(l int) int {
	return l + 16 - l%16
}

func constructMacData(aad, ciphertext []byte) []byte {
	aadsize := size(len(aad))
	ciphertextsize := size(len(ciphertext))

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
	if !reflect.DeepEqual(calcTag, tag) {
		return nil, errors.New("invalid tag")
	}

	return encrypt(key, nonce, ciphertext, 1), nil
}
