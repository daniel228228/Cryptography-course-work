package pkcs7

import (
	"bytes"
)

func Padding(b []byte, blockSize int) []byte {
	padding := blockSize - len(b)%blockSize
	return append(b, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

func Unpadding(b []byte) []byte {
	return b[:len(b)-int(b[len(b)-1])]
}
