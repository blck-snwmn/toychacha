package toychacha

import (
	"math/big"
)

var (
	p       *big.Int
	clamper *big.Int
	mask130 *big.Int
	five    *big.Int
)

func init() {
	i := new(big.Int)
	shifted := i.Lsh(big.NewInt(1), 130)
	p = new(big.Int).Sub(shifted, big.NewInt(5))

	// mask130 = 2^130 - 1
	mask130 = new(big.Int).Sub(shifted, big.NewInt(1))

	five = big.NewInt(5)

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

// modP computes x mod (2^130 - 5) using fast reduction.
// Since 2^130 ≡ 5 (mod p), we can reduce by splitting x into high and low parts.
// high and low are pre-allocated workspace provided by the caller.
func modP(x, high, low *big.Int) *big.Int {
	for x.BitLen() > 130 {
		// high = x >> 130 (upper bits)
		high.Rsh(x, 130)
		// low = x & mask130 (lower 130 bits)
		low.And(x, mask130)
		// x = high * 5 + low
		x.Mul(high, five)
		x.Add(x, low)
	}

	// Final adjustment: if x >= p, subtract p
	if x.Cmp(p) >= 0 {
		x.Sub(x, p)
	}
	return x
}

// macWithStdMod is an implementation using standard Mod for benchmarking comparison.
func macWithStdMod(msg, key []byte) [16]byte {
	r := leBytesToNum(key[0:16])
	clamp(r)

	s := leBytesToNum(key[16:32])

	a := big.NewInt(0)
	nn := make([]byte, 17)
	for len(msg) > 0 {
		l := min(16, len(msg))

		copy(nn[0:l], msg[0:l])
		nn[l] = 0x01
		nn = nn[:l+1]
		block := leBytesToNum(nn)

		a = a.Add(a, block)
		a = a.Mul(a, r)
		a = a.Mod(a, p)

		msg = msg[l:]
	}
	result := a.Add(a, s)

	return numTo16LeBytes(result)
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

func numTo16LeBytes(n *big.Int) [16]byte {
	b := n.Bytes()
	// Use only the last 16 bytes if b length > 16
	bstart := max(0, len(b)-16)
	b = b[bstart:] // b's lenth <= 16

	var result [16]byte

	// pack values from the end of the array.
	start := 16 - len(b)
	copy(result[start:], b)
	convertLittleEndian(result[:])
	return result
}

func mac(msg, key []byte) [16]byte {
	r := leBytesToNum(key[0:16])
	clamp(r)

	s := leBytesToNum(key[16:32])

	a := big.NewInt(0)
	nn := make([]byte, 17)

	// Workspace for modP (reused across loop iterations)
	high := new(big.Int)
	low := new(big.Int)

	for len(msg) > 0 {
		l := min(16, len(msg))

		copy(nn[0:l], msg[0:l])
		nn[l] = 0x01  // Index is almost 16. Index is len(msg) only once
		nn = nn[:l+1] // Range is almost [0:17]. Range is [0:len(msg)+1] only once
		block := leBytesToNum(nn)

		a = a.Add(a, block)
		a = a.Mul(a, r)
		a = modP(a, high, low)

		msg = msg[l:]
	}
	result := a.Add(a, s)

	return numTo16LeBytes(result)
}
