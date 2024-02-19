package twofish

import (
	"crypto-app/internal/bits"
)

func g0(X uint32, key *Key) uint32 {
	return key.S[0][bits.B0(X)] ^ key.S[1][bits.B1(X)] ^ key.S[2][bits.B2(X)] ^ key.S[3][bits.B3(X)]
}

func g1(X uint32, key *Key) uint32 {
	return key.S[0][bits.B3(X)] ^ key.S[1][bits.B0(X)] ^ key.S[2][bits.B1(X)] ^ key.S[3][bits.B2(X)]
}

func encryptRnd(A, B, C, D, T0, T1 *uint32, key *Key, r int) {
	*T0 = g0(*A, key)
	*T1 = g1(*B, key)
	*C ^= *T0 + *T1 + key.K[8+2*r]
	*C = bits.Ror32(*C, 1)
	*D = bits.Rol32(*D, 1)
	*D ^= *T0 + 2*(*T1) + key.K[8+2*r+1]
}

func encryptCycle(A, B, C, D, T0, T1 *uint32, key *Key, r int) {
	encryptRnd(A, B, C, D, T0, T1, key, 2*r)
	encryptRnd(C, D, A, B, T0, T1, key, 2*r+1)
}

func encrypt(A, B, C, D *uint32, key *Key) {
	var T0, T1 uint32

	for i := 0; i < 8; i++ {
		encryptCycle(A, B, C, D, &T0, &T1, key, i)
	}
}

func decryptRnd(A, B, C, D, T0, T1 *uint32, key *Key, r int) {
	*T0 = g0(*A, key)
	*T1 = g1(*B, key)
	*C = bits.Rol32(*C, 1)
	*C ^= *T0 + *T1 + key.K[8+2*r]
	*D ^= *T0 + 2*(*T1) + key.K[8+2*r+1]
	*D = bits.Ror32(*D, 1)
}

func decryptCycle(A, B, C, D, T0, T1 *uint32, key *Key, r int) {
	decryptRnd(A, B, C, D, T0, T1, key, 2*r+1)
	decryptRnd(C, D, A, B, T0, T1, key, 2*r)
}

func decrypt(A, B, C, D *uint32, key *Key) {
	var T0, T1 uint32

	for i := 7; i >= 0; i-- {
		decryptCycle(A, B, C, D, &T0, &T1, key, i)
	}
}

func inputWhitening(src []byte, A, B, C, D *uint32, key *Key, koff int) {
	*A = bits.Get32(src[:]) ^ key.K[koff]
	*B = bits.Get32(src[4:]) ^ key.K[1+koff]
	*C = bits.Get32(src[8:]) ^ key.K[2+koff]
	*D = bits.Get32(src[12:]) ^ key.K[3+koff]
}

func outputWhitening(dst []byte, A, B, C, D *uint32, key *Key, koff int) {
	*A ^= key.K[koff]
	*B ^= key.K[1+koff]
	*C ^= key.K[2+koff]
	*D ^= key.K[3+koff]

	bits.Put32(*A, dst[:])
	bits.Put32(*B, dst[4:])
	bits.Put32(*C, dst[8:])
	bits.Put32(*D, dst[12:])
}
