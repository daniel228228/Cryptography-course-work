package gf

func Mult(a, b byte, pol int) (res byte) {
	for i := 0; i < 8; i++ {
		if b&1 == 1 {
			res ^= a
		}

		if int(a) > int(a)^pol {
			a = byte(int(a) ^ pol)
		}

		a <<= 1
		b >>= 1
	}

	return
}

func MultMatrix(vector, matrix []byte, pol int) (res []byte) {
	res = make([]byte, 4)

	for i := range res {
		for j := 0; j < 4; j++ {
			res[i] ^= Mult(vector[j], matrix[i*8+j], pol)
		}
	}

	return
}
