package gochacha

import (
	"math/big"
)

func clamp(r []byte) {
	t := []byte{
		0x0f, 0xff, 0xff, 0xfc,
		0x0f, 0xff, 0xff, 0xfc,
		0x0f, 0xff, 0xff, 0xfc,
		0x0f, 0xff, 0xff, 0xff,
	}
	for i := range r {
		r[i] &= t[i]
	}
	// xor(r, []byte{
	// 	0x0f, 0xff, 0xff, 0xfc,
	// 	0x0f, 0xff, 0xff, 0xfc,
	// 	0x0f, 0xff, 0xff, 0xfc,
	// 	0x0f, 0xff, 0xff, 0xff,
	// })
}

var p *big.Int

func init() {
	i := new(big.Int)
	shifted := i.Lsh(big.NewInt(1), 130)
	p = i.Sub(shifted, big.NewInt(5))
}

func leBytesToNum(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func mac(msg, key []byte) []byte {

	rr := leBytesToNum(key[0:16])
	clamp(rr)
	r := new(big.Int).SetBytes(rr)

	ss := leBytesToNum(key[16:32])
	s := new(big.Int).SetBytes(ss)

	a := big.NewInt(0)
	for len(msg) > 0 {
		l := 16
		if len(msg) < l {
			l = len(msg)
		}
		nn := make([]byte, l, l+1)
		copy(nn, msg[0:l])
		nn = append(nn, 0x01)
		nn = leBytesToNum(nn)
		block := new(big.Int).SetBytes(nn)

		a = a.Add(a, block)
		a = a.Mul(a, r)
		a = a.Mod(a, p)

		msg = msg[l:]
	}
	result := a.Add(a, s)
	b := result.Bytes()
	// 128 least significant bits
	tag := leBytesToNum(b[len(b)-16:])
	return tag
}
