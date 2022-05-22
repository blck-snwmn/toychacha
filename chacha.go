package gochacha

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

