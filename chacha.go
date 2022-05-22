package gochacha

import (
	"encoding/binary"
	"fmt"
)

func rotationN(n uint32, shift uint) uint32 {
	return n>>(32-shift) | n<<shift
}

func quarterRound(a, b, c, d uint32) (uint32, uint32, uint32, uint32) {
	a += b
	d ^= a
	d = rotationN(d, 16)

	c += d
	b ^= c
	b = rotationN(b, 12)

	a += b
	d ^= a
	d = rotationN(d, 8)

	c += d
	b ^= c
	b = rotationN(b, 7)
	return a, b, c, d
}

type state [][]uint32

func (s state) quarterRound(x, y, z, w uint) {
	getIndex := func(v uint) (uint, uint) {
		return v / 4, v % 4
	}
	getValue := func(v uint) uint32 {
		i, j := getIndex(v)
		return s[i][j]
	}
	setValue := func(v uint, value uint32) {
		i, j := getIndex(v)
		s[i][j] = value
	}
	a := getValue(x)
	b := getValue(y)
	c := getValue(z)
	d := getValue(w)

	a, b, c, d = quarterRound(a, b, c, d)

	setValue(x, a)
	setValue(y, b)
	setValue(z, c)
	setValue(w, d)
}

type keySizeError int

func (k keySizeError) Error() string {
	return fmt.Sprintf("invalid key length. got=%d, want=%d", k, 32)
}

type nonceSizeError int

func (n nonceSizeError) Error() string {
	return fmt.Sprintf("invalid nonce length. got=%d, want=%d", n, 12)
}
func NewState(key, nonce []byte) (state, error) {
	if len(key) != 32 {
		return nil, keySizeError(len(key))
	}
	if len(nonce) != 12 {
		return nil, nonceSizeError(len(nonce))
	}
	s := make(state, 0, 4)
	// magic
	s = append(s, []uint32{0x61707865, 0x3320646e, 0x79622d32, 0x6b206574})

	s = append(s, []uint32{
		binary.LittleEndian.Uint32(key[28:32]),
		binary.LittleEndian.Uint32(key[24:28]),
		binary.LittleEndian.Uint32(key[20:24]),
		binary.LittleEndian.Uint32(key[16:20]),
	})

	s = append(s, []uint32{
		binary.LittleEndian.Uint32(key[12:16]),
		binary.LittleEndian.Uint32(key[8:12]),
		binary.LittleEndian.Uint32(key[4:8]),
		binary.LittleEndian.Uint32(key[0:4]),
	})
	s = append(s, []uint32{
		1,
		binary.LittleEndian.Uint32(nonce[8:12]),
		binary.LittleEndian.Uint32(nonce[4:8]),
		binary.LittleEndian.Uint32(nonce[0:4]),
	})

	return s, nil
}
