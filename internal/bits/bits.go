package bits

const Uint32Mask = (2 << 31) - 1

func Read(b byte, pos int) byte {
	return b >> (7 - pos) & 1
}

func Set(b *byte, val byte, pos int) {
	if val&1 == 1 {
		*b |= 1 << (7 - pos)
	} else {
		*b &^= 1 << (7 - pos)
	}
}

func Rol32(x, n uint32) uint32 {
	return (x << n) | ((x & Uint32Mask) >> (32 - n))
}

func Ror32(x, n uint32) uint32 {
	return (x >> n) | ((x & Uint32Mask) << (32 - n))
}

func Get32(p []byte) uint32 {
	return uint32(p[0]) | uint32(p[1])<<8 | uint32(p[2])<<16 | uint32(p[3])<<24
}

func Put32(v uint32, p []byte) {
	p[0] = byte(v & 0xff)
	p[1] = byte((v >> 8) & 0xff)
	p[2] = byte((v >> 16) & 0xff)
	p[3] = byte((v >> 24) & 0xff)
}

func B0(X uint32) uint32 { return X & 0xff }
func B1(X uint32) uint32 { return (X >> 8) & 0xff }
func B2(X uint32) uint32 { return (X >> 16) & 0xff }
func B3(X uint32) uint32 { return (X >> 24) & 0xff }

func H02(y byte, l []byte, mds [4][256]uint32, q [2][256]byte) uint32 {
	return mds[0][q[0][q[0][y]^l[8]]^l[0]]
}
func H12(y byte, l []byte, mds [4][256]uint32, q [2][256]byte) uint32 {
	return mds[1][q[0][q[1][y]^l[9]]^l[1]]
}
func H22(y byte, l []byte, mds [4][256]uint32, q [2][256]byte) uint32 {
	return mds[2][q[1][q[0][y]^l[10]]^l[2]]
}
func H32(y byte, l []byte, mds [4][256]uint32, q [2][256]byte) uint32 {
	return mds[3][q[1][q[1][y]^l[11]]^l[3]]
}
func H03(y byte, l []byte, mds [4][256]uint32, q [2][256]byte) uint32 {
	return H02(q[1][y]^l[16], l, mds, q)
}
func H13(y byte, l []byte, mds [4][256]uint32, q [2][256]byte) uint32 {
	return H12(q[1][y]^l[17], l, mds, q)
}
func H23(y byte, l []byte, mds [4][256]uint32, q [2][256]byte) uint32 {
	return H22(q[0][y]^l[18], l, mds, q)
}
func H33(y byte, l []byte, mds [4][256]uint32, q [2][256]byte) uint32 {
	return H32(q[0][y]^l[19], l, mds, q)
}
func H04(y byte, l []byte, mds [4][256]uint32, q [2][256]byte) uint32 {
	return H03(q[1][y]^l[24], l, mds, q)
}
func H14(y byte, l []byte, mds [4][256]uint32, q [2][256]byte) uint32 {
	return H23(q[0][y]^l[25], l, mds, q)
}
func H24(y byte, l []byte, mds [4][256]uint32, q [2][256]byte) uint32 {
	return H23(q[0][y]^l[26], l, mds, q)
}
func H34(y byte, l []byte, mds [4][256]uint32, q [2][256]byte) uint32 {
	return H33(q[1][y]^l[27], l, mds, q)
}
