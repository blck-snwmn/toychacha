package gochacha

import (
	"math/big"
)

var (
	p       *big.Int
	p_      *big.Int
	r       *big.Int
	r2      *big.Int
	mask    *big.Int
	clamper *big.Int
)

func init() {
	r = new(big.Int).Lsh(big.NewInt(1), 130)
	mask = new(big.Int).Sub(r, big.NewInt(1))
	p = new(big.Int).Sub(r, big.NewInt(5))
	r2 = new(big.Int)
	r2 = r2.Mul(r, r).Mod(r2, p)

	clamper = new(big.Int).SetBytes([]byte{
		0x0f, 0xff, 0xff, 0xfc,
		0x0f, 0xff, 0xff, 0xfc,
		0x0f, 0xff, 0xff, 0xfc,
		0x0f, 0xff, 0xff, 0xff,
	})

	var (
		result = big.NewInt(0)
		t      = big.NewInt(0)
		rr     = new(big.Int).Set(r)
		i      = big.NewInt(1)
	)

	for rr.Cmp(big.NewInt(1)) > 0 {
		if t.Bit(0) == 0 {
			t.Add(t, p)
			result.Add(result, i)
		}
		t.Rsh(t, 1)
		rr.Rsh(rr, 1)
		i.Lsh(i, 1)
	}
	p_ = result
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

// mac_ calc mac using big.Int w/ montgomety multiplication
func mac(msg, key []byte) []byte {
	r := leBytesToNum(key[0:16])
	clamp(r)
	rr := montgomeryRepresentation(r)

	s := leBytesToNum(key[16:32])

	a := big.NewInt(0)
	a = montgomeryRepresentation(a)
	nn := make([]byte, 17)
	for len(msg) > 0 {
		l := 16
		if len(msg) < l {
			l = len(msg)
		}
		copy(nn[0:l], msg[0:l])
		nn[l] = 0x01  // Index is almost 16. Index is len(msg) only once
		nn = nn[:l+1] // Range is almost [0:17]. Range is [0:len(msg)+1] only once
		block := leBytesToNum(nn)
		block = montgomeryRepresentation(block)

		a = mr(a.Add(a, block).Mul(a, rr))

		msg = msg[l:]
	}
	// result := a.Add(a, ss)
	// result = mr(result)
	// 加算処理でズレが生じてしまうため、 mr 後に通常の加算を行う
	tttt := mr(a)
	result := tttt.Add(tttt, s)
	return numTo16LeBytes(result)
}

func montgomeryRepresentation(t *big.Int) *big.Int {
	tmp := new(big.Int).Set(t)
	return mr(tmp.Mul(tmp, r2))
}

// mr do montgomery reduction
func mr(t *big.Int) *big.Int {
	tmp := new(big.Int)
	tmp = tmp.
		Mul(t, p_).
		And(tmp, mask).
		Mul(tmp, p).
		Add(tmp, t).
		Rsh(tmp, uint(r.BitLen()-1))
	switch tmp.Cmp(p) {
	case -1: // tmp < p
		return tmp
	default: // tmp >= p
		return tmp.Sub(tmp, p)
	}
}

// mac_ calc mac using big.Int
func mac_(msg, key []byte) []byte {
	r := leBytesToNum(key[0:16])
	clamp(r)

	s := leBytesToNum(key[16:32])

	a := big.NewInt(0)
	nn := make([]byte, 17)
	for len(msg) > 0 {
		l := 16
		if len(msg) < l {
			l = len(msg)
		}
		copy(nn[0:l], msg[0:l])
		nn[l] = 0x01  // Index is almost 16. Index is len(msg) only once
		nn = nn[:l+1] // Range is almost [0:17]. Range is [0:len(msg)+1] only once
		block := leBytesToNum(nn)

		a = a.Add(a, block)
		a = a.Mul(a, r)
		a = a.Mod(a, p)

		msg = msg[l:]
	}
	result := a.Add(a, s)

	return numTo16LeBytes(result)
}
