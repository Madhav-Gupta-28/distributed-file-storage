package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

func hashKey(key string) string {
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

func NewEncryptKey() []byte {
	key := make([]byte, 32)
	io.ReadFull(rand.Reader, key)
	return key
}

func generateId() string {
	buf := make([]byte, 32)
	io.ReadFull(rand.Reader, buf)
	return hex.EncodeToString(buf)
}

func CopyEncrypt(key []byte, r io.Reader, w io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	iv := make([]byte, aes.BlockSize)

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return 0, err
	}

	// prepand the vi to the file
	if _, err := w.Write(iv); err != nil {
		return 0, err
	}

	buf := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
	nw := block.BlockSize()

	for {
		n, err := r.Read(buf)

		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			nn, err := w.Write(buf[:n])
			if err != nil {
				return nw, err
			}
			nw += nn
		}

		if err == io.EOF {
			break
		}
	}

	return nw, nil
}

func CopyDecrypt(key []byte, r io.Reader, w io.Writer) (int, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	iv := make([]byte, aes.BlockSize)

	if _, err := r.Read(iv); err != nil {
		return 0, err
	}

	buf := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
	nw := block.BlockSize()

	for {
		n, err := r.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			nn, err := w.Write(buf[:n])
			if err != nil {
				return nn, err
			}

			nw += nn
		}

		if err == io.EOF {
			break
		}
	}

	return nw, nil
}
