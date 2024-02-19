package twofish

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
)

var ErrIncorrectKey = errors.New("incorrect key")

func ReadEnc(filename string) ([]byte, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	return io.ReadAll(fp)
}

func NewContextFromData(data []byte) (*Context, string, error) {
	c := &Context{}

	if len(data) < 1 {
		return nil, "", ErrIncorrectKey
	}

	keySize := int(data[0])

	if len(data) < 1+2*keySize+3 {
		return nil, "", ErrIncorrectKey
	}

	key := make([]byte, keySize)
	if n := copy(key, data[1:]); n < keySize {
		return nil, "", ErrIncorrectKey
	}

	iv := make([]byte, keySize)
	if n := copy(iv, data[1+keySize:]); n < keySize {
		return nil, "", ErrIncorrectKey
	}

	m := make([]byte, 3)
	if n := copy(m, data[1+2*keySize:]); n < 3 {
		return nil, "", ErrIncorrectKey
	}

	mode, err := ParseMode(string(m))
	if err != nil {
		return nil, "", err
	}

	scheduledKey, err := keySchedule(key)
	if err != nil {
		return nil, "", err
	}

	c.scheduledKey = scheduledKey
	c.Key = key
	c.Iv = iv
	c.Mode = mode

	return c, string(data[1+2*keySize+3:]), nil
}

func (c *Context) WriteEnc(data []byte, filename string, perms uint32) error {
	return save(bytes.NewReader(data), filename, perms)
}

func (c *Context) GetKey(inputFilename string) []byte {
	key := []byte{byte(BlockSize)}
	key = append(key, c.Key...)
	key = append(key, c.Iv...)
	key = append(key, c.Mode.String()...)
	key = append(key, inputFilename...)

	return key
}

func save(data io.Reader, filename string, perms uint32) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fs.FileMode(perms))
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
