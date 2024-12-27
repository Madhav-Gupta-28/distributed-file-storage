package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))

	hashString := hex.EncodeToString(hash[:])
	blockSize := 5
	sliceLength := len(hashString) / blockSize
	paths := make([]string, sliceLength)

	for i := 0; i < sliceLength; i++ {
		paths[i] = hashString[i*blockSize : (i+1)*blockSize]
	}

	patkkey := PathKey{
		Pathname: strings.Join(paths, "/"),
		Filename: hashString,
	}

	return patkkey

}

type PathTransformFunc func(string) PathKey

var DefualtPathTransformFunc = func(key string) string {
	return key
}

type PathKey struct {
	Pathname string
	Filename string
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s%s", p.Pathname, p.Filename)

}

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{StoreOpts: opts}
}

func (s *Store) Read(key string) (io.Reader, error) {

	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {

	pathkey := s.PathTransformFunc(key)

	return os.Open(pathkey.FullPath())

}

func (s *Store) writeStream(key string, r io.Reader) error {

	pathKey := s.PathTransformFunc(key)
	if err := os.MkdirAll(pathKey.Pathname, os.ModePerm); err != nil {
		return err
	}
	// buf := new(bytes.Buffer)
	// io.Copy(buf, r)
	// filenameBytes := md5.Sum(buf.Bytes())
	// filename := hex.EncodeToString(filenameBytes[:])
	pathFilename := pathKey.FullPath()

	f, err := os.Create(pathFilename)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	fmt.Printf("wrote %d bytes to disk : %s\n ", n, pathFilename)
	return nil
}
