package twofish

import (
	"errors"

	"crypto-app/internal/bits"
)

type Key struct {
	S [4][256]uint32
	K [40]uint32
}

func keySchedule(key []byte) (*Key, error) {
	switch len(key) {
	case 128 / 8, 192 / 8, 256 / 8:
	default:
		return nil, errors.New("incorrect key length")
	}

	Key := &Key{}

	k := len(key) / 8
	if k < 2 {
		k = 2
	}

	K := make([]byte, 32+32+4)
	var A, B uint32

	copy(K, key)

	for i := 0; i < 40; i += 2 {
		A = h(i, K, k)
		B = h(i+1, K[4:], k)
		B = bits.Rol32(B, 8)

		A += B
		B += A

		Key.K[i] = A
		Key.K[i+1] = bits.Rol32(B, 9)
	}

	kptr := 8 * k
	sptr := 32

	for kptr > 0 {
		kptr -= 8

		for i := sptr; i < sptr+4; i++ {
			K[i] = 0
		}

		for i, j := sptr+4, 0; j < 8; i, j = i+1, j+1 {
			K[i] = K[kptr+j]
		}

		for t := sptr + 11; t > sptr+3; t-- {
			b := K[t]
			bx := byte(uint32(b<<1) ^ RsPol[b>>7])
			bxx := byte(uint32(b>>1) ^ RsPolDiv[b&1] ^ uint32(bx))

			K[t-1] ^= bxx
			K[t-2] ^= bx
			K[t-3] ^= bxx
			K[t-4] ^= b
		}

		sptr += 8
	}

	fillKeyedSboxes(K[32:], k, Key)

	for i := range K {
		K[i] = 0
	}

	return Key, nil
}

func h(K int, l []byte, k int) uint32 { // Twofish 4.3.2
	switch k {
	case 2:
		return bits.H02(byte(K), l, MDS, Q) ^ bits.H12(byte(K), l, MDS, Q) ^ bits.H22(byte(K), l, MDS, Q) ^ bits.H32(byte(K), l, MDS, Q)
	case 3:
		return bits.H03(byte(K), l, MDS, Q) ^ bits.H13(byte(K), l, MDS, Q) ^ bits.H23(byte(K), l, MDS, Q) ^ bits.H33(byte(K), l, MDS, Q)
	case 4:
		return bits.H04(byte(K), l, MDS, Q) ^ bits.H14(byte(K), l, MDS, Q) ^ bits.H24(byte(K), l, MDS, Q) ^ bits.H34(byte(K), l, MDS, Q)
	default:
		return 0
	}
}

func fillKeyedSboxes(s []byte, k int, key *Key) {
	switch k {
	case 2:
		for i := 0; i < 256; i++ {
			key.S[0][i] = bits.H02(byte(i), s, MDS, Q)
			key.S[1][i] = bits.H12(byte(i), s, MDS, Q)
			key.S[2][i] = bits.H22(byte(i), s, MDS, Q)
			key.S[3][i] = bits.H32(byte(i), s, MDS, Q)
		}
	case 3:
		for i := 0; i < 256; i++ {
			key.S[0][i] = bits.H03(byte(i), s, MDS, Q)
			key.S[1][i] = bits.H13(byte(i), s, MDS, Q)
			key.S[2][i] = bits.H23(byte(i), s, MDS, Q)
			key.S[3][i] = bits.H33(byte(i), s, MDS, Q)
		}
	case 4:
		for i := 0; i < 256; i++ {
			key.S[0][i] = bits.H04(byte(i), s, MDS, Q)
			key.S[1][i] = bits.H14(byte(i), s, MDS, Q)
			key.S[2][i] = bits.H24(byte(i), s, MDS, Q)
			key.S[3][i] = bits.H34(byte(i), s, MDS, Q)
		}
	default:
		break
	}
}
