package gochacha

import (
	"math/big"
)

var (
	p       *big.Int
	clamper *big.Int
)

func init() {
	i := new(big.Int)
	shifted := i.Lsh(big.NewInt(1), 130)
	p = i.Sub(shifted, big.NewInt(5))

	clamper = new(big.Int).SetBytes([]byte{
		0x0f, 0xff, 0xff, 0xfc,
		0x0f, 0xff, 0xff, 0xfc,
		0x0f, 0xff, 0xff, 0xfc,
		0x0f, 0xff, 0xff, 0xff,
	})
}

func genMacKey(key, nonce []byte) []byte {
	var counter uint32 = 0
	block := block(key, nonce, counter)
	return block[0:32]
}

func clamp(n *big.Int) {
	n.And(n, clamper)
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

var zeros = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func numTo16LeBytes(n *big.Int) []byte {
	b := n.Bytes()
	// padding 0 if len(b) < 16
	b = append(zeros, b...)
	// 128 least significant bits
	b = b[len(b)-16:]
	convertLittleEndian(b)
	return b
}

func mac(msg, key []byte) []byte {
	r := leBytesToNum(key[0:16])
	clamp(r)

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
