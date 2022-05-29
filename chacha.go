package gochacha

import (
	"encoding/binary"
	"fmt"
)

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
	s := make(state, 16)

	// magic
	s[0] = 0x61707865
	s[1] = 0x3320646e
	s[2] = 0x79622d32
	s[3] = 0x6b206574

	s[4] = binary.LittleEndian.Uint32(key[0:4])
	s[5] = binary.LittleEndian.Uint32(key[4:8])
	s[6] = binary.LittleEndian.Uint32(key[8:12])
	s[7] = binary.LittleEndian.Uint32(key[12:16])

	s[8] = binary.LittleEndian.Uint32(key[16:20])
	s[9] = binary.LittleEndian.Uint32(key[20:24])
	s[10] = binary.LittleEndian.Uint32(key[24:28])
	s[11] = binary.LittleEndian.Uint32(key[28:32])

	s[12] = counter
	s[13] = binary.LittleEndian.Uint32(nonce[0:4])
	s[14] = binary.LittleEndian.Uint32(nonce[4:8])
	s[15] = binary.LittleEndian.Uint32(nonce[8:12])

	return s, nil
}

type state []uint32

func (s state) quarterRound(x, y, z, w uint) {
	a := s[x]
	b := s[y]
	c := s[z]
	d := s[w]

	a, b, c, d = quarterRound(a, b, c, d)

	s[x] = a
	s[y] = b
	s[z] = c
	s[w] = d
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
	newS := make(state, 16)
	copy(newS, s)
	return newS
}

func (s state) add(other state) {
	for i := 0; i < 16; i++ {
		s[i] += other[i]
	}
}

func (s state) serialize() []byte {
	// state is 4*4 size
	// state'element is uint32(4byte)
	serialized := make([]byte, 4*4*4)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			index := (i*4 + j)
			offset := index * 4
			binary.LittleEndian.PutUint32(serialized[offset:offset+4], s[index])
		}
	}
	return serialized
}

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
	l := 64
	for len(plaintext) > 0 {
		// Process in 64-byte units
		if len(plaintext) < l {
			//  Use the remaining plaintext, if the plaintext is less than 64 bytes in length
			l = len(plaintext)
		}

		keyStream := block(key, nonce, counter)
		copy(header[0:l], plaintext[0:l])
		xor(header[0:l], keyStream)

		counter++
		plaintext, header = plaintext[l:], header[l:]
	}

	return encrypted
}
