package gochacha

import (
	"math/big"
)

var p *big.Int

func init() {
	i := new(big.Int)
	shifted := i.Lsh(big.NewInt(1), 130)
	p = i.Sub(shifted, big.NewInt(5))
}

func clamp(n *big.Int) *big.Int {
	t := new(big.Int).SetBytes([]byte{
		0x0f, 0xff, 0xff, 0xfc,
		0x0f, 0xff, 0xff, 0xfc,
		0x0f, 0xff, 0xff, 0xfc,
		0x0f, 0xff, 0xff, 0xff,
	})
	t.And(n, t)
	return t
}

func convertLittleEndian(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func leBytesToNum(b []byte) *big.Int {
	convertLittleEndian(b)
	return new(big.Int).SetBytes(b)
}

func numTo16LeBytes(n *big.Int) []byte {
	b := n.Bytes()
	// 128 least significant bits
	b = b[len(b)-16:]
	convertLittleEndian(b)
	return b
}

func mac(msg, key []byte) []byte {
	r := leBytesToNum(key[0:16])
	r = clamp(r)

	s := leBytesToNum(key[16:32])

	a := big.NewInt(0)
	for len(msg) > 0 {
		l := 16
		if len(msg) < l {
			l = len(msg)
		}
		nn := make([]byte, l, l+1)
		copy(nn, msg[0:l])
		nn = append(nn, 0x01)
		block := leBytesToNum(nn)

		a = a.Add(a, block)
		a = a.Mul(a, r)
		a = a.Mod(a, p)

		msg = msg[l:]
	}
	result := a.Add(a, s)

	return numTo16LeBytes(result)
}
