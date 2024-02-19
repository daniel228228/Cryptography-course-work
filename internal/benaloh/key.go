package benaloh

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"math/big"
	"os"
	"strings"
)

var ErrIncorrectKey = errors.New("incorrect key")

func (b *Benaloh) ReadPublicKey(filename string) error {
	fp, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fp.Close()

	content, err := io.ReadAll(fp)
	if err != nil {
		return err
	}

	parts := strings.Split(string(content), " ")
	if len(parts) != 3 {
		return ErrIncorrectKey
	}

	for i, part := range parts {
		val, ok := new(big.Int).SetString(part, 10)
		if !ok {
			return ErrIncorrectKey
		}

		b.Public[i] = val
	}

	return nil
}

func (b *Benaloh) ReadPrivateKey(filename string) error {
	fp, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fp.Close()

	content, err := io.ReadAll(fp)
	if err != nil {
		return err
	}

	parts := strings.Split(string(content), " ")
	if len(parts) != 2 {
		return ErrIncorrectKey
	}

	for i, part := range parts {
		val, ok := new(big.Int).SetString(part, 10)
		if !ok {
			return ErrIncorrectKey
		}

		b.Private[i] = val
	}

	return nil
}

func (b *Benaloh) WritePublicKey(dir, filename string, perms uint32) error {
	publicKey := &bytes.Buffer{}
	publicKey.WriteString(b.Public[0].String())
	publicKey.WriteRune(' ')
	publicKey.WriteString(b.Public[1].String())
	publicKey.WriteRune(' ')
	publicKey.WriteString(b.Public[2].String())

	return saveInDir(publicKey, dir, filename, perms)
}

func (b *Benaloh) WritePrivateKey(dir, filename string, perms uint32) error {
	privateKey := &bytes.Buffer{}
	privateKey.WriteString(b.Private[0].String())
	privateKey.WriteRune(' ')
	privateKey.WriteString(b.Private[1].String())

	return saveInDir(privateKey, dir, filename, perms)
}

func saveInDir(data io.Reader, dir, filename string, perms uint32) error {
	ex, _ := os.Executable()
	if err := os.MkdirAll(ex+"/../"+dir, 0700); err != nil {
		return err
	}

	f, err := os.OpenFile(ex+"/../"+dir+"/"+filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fs.FileMode(perms))
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, data); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
