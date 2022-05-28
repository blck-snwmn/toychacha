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

func aeadEncrpt(aad, key, iv, constant, plaintext []byte) ([]byte, []byte) {
	nonce := append(constant, iv...)
	otk := genMacKey(key, nonce)

	ciphertext := encrypt(key, nonce, plaintext, 1)

	// TODO 最初からサイズを決め打ちにして、あとからpaddingしなくてもよくする
	macData := append(aad, pad16(aad)...)

	macData = append(macData, ciphertext...)
	macData = append(macData, pad16(ciphertext)...)

	macData = append(macData, numTo8LeBytes(len(aad))...)
	macData = append(macData, numTo8LeBytes(len(ciphertext))...)

	tag := mac(macData, otk)
	return ciphertext, tag
}

func aeadDecrypt(aad, key, iv, constant, ciphertext, tag []byte) ([]byte, error) {
	nonce := append(constant, iv...)
	otk := genMacKey(key, nonce)

	macData := append(aad, pad16(aad)...)

	macData = append(macData, ciphertext...)
	macData = append(macData, pad16(ciphertext)...)

	macData = append(macData, numTo8LeBytes(len(aad))...)
	macData = append(macData, numTo8LeBytes(len(ciphertext))...)

	calcTag := mac(macData, otk)
	if !reflect.DeepEqual(calcTag, tag) {
		return nil, errors.New("invalid tag")
	}

	return encrypt(key, nonce, ciphertext, 1), nil
}
