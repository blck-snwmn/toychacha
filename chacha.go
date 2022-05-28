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

func (s state) innerBlock() {
	// column rounds
	s.quarterRound(0, 4, 8, 12)
	s.quarterRound(1, 5, 9, 13)
	s.quarterRound(2, 6, 10, 14)
	s.quarterRound(3, 7, 11, 15)

	// diagonal rounds
	s.quarterRound(0, 5, 10, 15)
	s.quarterRound(1, 6, 11, 12)
	s.quarterRound(2, 7, 8, 13)
	s.quarterRound(3, 4, 9, 14)
}

func (s state) clone() state {
	newS := make(state, 4)
	for i := 0; i < 4; i++ {
		newS[i] = append(newS[i], s[i]...)
	}
	return newS
}

func (s state) add(other state) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			s[i][j] += other[i][j]
		}
	}
}

func (s state) serialize() []byte {
	// state is 4*4 size
	// state'element is uint32(4byte)
	serialized := make([]byte, 4*4*4)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			offset := (i*4 + j) * 4
			binary.LittleEndian.PutUint32(serialized[offset:offset+4], s[i][j])
		}
	}
	return serialized
}

func block(key, nonce []byte, counter uint32) []byte {
	s, _ := newState(key, nonce, counter)
	init := s.clone()
	for i := 0; i < 10; i++ {
		s.innerBlock()
	}
	s.add(init)
	return s.serialize()
}

func xor(l, r []byte) []byte {
	if len(l) > len(r) {
		l, r = r, l
	}

	for i := 0; i < len(l); i++ {
		l[i] ^= r[i]
	}
	return l
}

func encrypt(key, nonce, plaintext []byte, counter uint32) []byte {
	encrypted := make([]byte, len(plaintext))

	header := encrypted
	// TODO one loop
	for len(plaintext) >= 64 {
		keyStream := block(key, nonce, counter)
		b, h := plaintext[0:64], header[0:64]
		copy(h, b)
		xor(h, keyStream)

		counter++ // for next loop
		plaintext = plaintext[64:]
		header = header[64:]
	}
	if len(plaintext) > 0 {
		keyStream := block(key, nonce, counter)
		b, h := plaintext, header
		copy(h, b)
		xor(h, keyStream)
	}
	return encrypted
}

type keySizeError int

func (k keySizeError) Error() string {
	return fmt.Sprintf("invalid key length. got=%d, want=%d", k, 32)
}

type nonceSizeError int

func (n nonceSizeError) Error() string {
	return fmt.Sprintf("invalid nonce length. got=%d, want=%d", n, 12)
}
func newState(key, nonce []byte, counter uint32) (state, error) {
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
		binary.LittleEndian.Uint32(key[0:4]),
		binary.LittleEndian.Uint32(key[4:8]),
		binary.LittleEndian.Uint32(key[8:12]),
		binary.LittleEndian.Uint32(key[12:16]),
	})

	s = append(s, []uint32{
		binary.LittleEndian.Uint32(key[16:20]),
		binary.LittleEndian.Uint32(key[20:24]),
		binary.LittleEndian.Uint32(key[24:28]),
		binary.LittleEndian.Uint32(key[28:32]),
	})

	s = append(s, []uint32{
		counter,
		binary.LittleEndian.Uint32(nonce[0:4]),
		binary.LittleEndian.Uint32(nonce[4:8]),
		binary.LittleEndian.Uint32(nonce[8:12]),
	})

	return s, nil
}
