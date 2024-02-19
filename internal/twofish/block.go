package twofish

func encryptBlock(key *Key, p [BlockSize]byte, c *[BlockSize]byte) {
	var A, B, C, D uint32

	inputWhitening(p[:], &A, &B, &C, &D, key, 0)
	encrypt(&A, &B, &C, &D, key)
	outputWhitening(c[:], &C, &D, &A, &B, key, 4)
}

func decryptBlock(key *Key, c [BlockSize]byte, p *[BlockSize]byte) {
	var A, B, C, D uint32

	inputWhitening(c[:], &A, &B, &C, &D, key, 4)
	decrypt(&A, &B, &C, &D, key)
	outputWhitening(p[:], &C, &D, &A, &B, key, 0)
}
