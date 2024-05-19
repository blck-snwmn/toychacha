package toychacha

import (
	"encoding/binary"
	"math/rand/v2"
)

type Chacha8 struct {
	states []state
	flip   bool
	i      uint32
}

var _ rand.Source = (*Chacha8)(nil)

func NewChaCha8(seed [32]byte) *Chacha8 {
	ss := make([]state, 0, 4)
	for i := 0; i < 4; i++ {
		nonce := make([]byte, 12)
		s, _ := newState(seed[:], nonce, uint32(i))
		init := s.clone()
		for i := 0; i < 4; i++ { // 4 iterations = 8 rounds
			s.innerBlock()
		}
		s[4] += init[4]
		s[5] += init[5]
		s[6] += init[6]
		s[7] += init[7]
		s[8] += init[8]
		s[9] += init[9]
		s[10] += init[10]
		s[11] += init[11]

		ss = append(ss, s)
	}
	return &Chacha8{
		states: ss,
		flip:   true,
	}
}

func (c *Chacha8) Uint64() uint64 {
	data := make([]byte, 8)

	index := c.i

	l, r := 0, 1
	if !c.flip {
		l, r = 2, 3
		// if flip false, increment i
		c.i++
	}
	c.flip = !c.flip
	// TODO bigendian?
	binary.LittleEndian.PutUint32(data[:], c.states[l][index])
	binary.LittleEndian.PutUint32(data[4:], c.states[r][index])
	return binary.LittleEndian.Uint64(data[:])
}
