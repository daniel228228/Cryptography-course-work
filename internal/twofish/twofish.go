package twofish

import (
	"bufio"
	"io"
	"os"

	"crypto-app/internal/pkcs7"
)

type Context struct {
	scheduledKey *Key
	Key, Iv      []byte
	Mode         Mode
}

func NewContext(key, iv []byte, mode Mode) (*Context, error) {
	c := &Context{}

	var err error
	if c.scheduledKey, err = keySchedule(key); err != nil {
		return nil, err
	}

	c.Key = key
	c.Iv = iv
	c.Mode = mode

	return c, nil
}

func (c *Context) EncryptFile(filepath, outpath string, callback func(progress float64)) error {
	fp, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer fp.Close()

	out, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer out.Close()

	var buf [BlockSize]byte

	var iv [BlockSize]byte
	copy(iv[:], c.Iv[:])

	mode := c.Mode

	fileInfo, err := fp.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()
	var bytesRead int64
	var progress float64

	reader := bufio.NewReader(fp)
	writer := bufio.NewWriter(out)

	shortLast := false
	read := false

	for !read {
		readSize, err := reader.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				read = true
			} else {
				return err
			}
		}

		if read {
			if shortLast {
				break
			}

			copy(buf[:], pkcs7.Padding(make([]byte, BlockSize), BlockSize)[BlockSize:])
		} else if readSize < BlockSize {
			copy(buf[:], pkcs7.Padding(buf[:readSize], BlockSize))
			shortLast = true
		}

		switch mode {
		case ECB:
			encryptBlock(c.scheduledKey, buf, &buf)

		case CBC:
			for i := range buf {
				buf[i] ^= iv[i]
			}
			encryptBlock(c.scheduledKey, buf, &buf)
			copy(iv[:], buf[:])

		case CFB, OFB:
			encryptBlock(c.scheduledKey, iv, &iv)
			for i := range buf {
				buf[i] ^= iv[i]
			}
			if mode == CFB {
				copy(iv[:], buf[:])
			}
		}

		if _, err := writer.Write(buf[:]); err != nil {
			return err
		}

		bytesRead += int64(readSize)
		progress = float64(bytesRead) / float64(fileSize)
		callback(progress)
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

func (c *Context) DecryptFile(filepath, outpath string, callback func(progress float64)) error {
	fp, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer fp.Close()

	tmp, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer tmp.Close()

	var buf [BlockSize]byte
	var iv [BlockSize]byte

	copy(iv[:], c.Iv[:])
	mode := c.Mode

	fileInfo, err := fp.Stat()
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()
	var bytesRead int64
	var progress float64

	reader := bufio.NewReader(fp)
	writer := bufio.NewWriter(tmp)

	var prev [BlockSize]byte

	first := true
	read := false

	for !read {
		readSize, err := reader.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				read = true
			} else {
				return err
			}
		}

		if first {
			copy(prev[:], buf[:])
			first = false
			continue
		}

		switch mode {
		case ECB:
			decryptBlock(c.scheduledKey, prev, &prev)

		case CBC:
			var temp [BlockSize]byte
			copy(temp[:], prev[:])
			decryptBlock(c.scheduledKey, prev, &prev)
			for i := range prev {
				prev[i] ^= iv[i]
			}
			copy(iv[:], temp[:])

		case CFB, OFB:
			var temp [BlockSize]byte
			encryptBlock(c.scheduledKey, iv, &iv)
			copy(temp[:], prev[:])
			for i := range prev {
				prev[i] ^= iv[i]
			}
			if mode == CFB {
				copy(iv[:], temp[:])
			}
		}

		if !read {
			if _, err := writer.Write(prev[:]); err != nil {
				return err
			}
		} else {
			if _, err := writer.Write(pkcs7.Unpadding(prev[:])); err != nil {
				return err
			}
		}

		copy(prev[:], buf[:])

		bytesRead += int64(readSize)
		progress = float64(bytesRead) / float64(fileSize)
		callback(progress)
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}
