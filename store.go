package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const defaultRootFolderName = "madhavnetwork"

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

var DefualtPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Filename: key,
	}
}

type PathKey struct {
	Pathname string
	Filename string
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s%s", p.Pathname, p.Filename)

}

type StoreOpts struct {
	// Root is the folder name of the root containing all the files and folders of the system
	Root              string
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefualtPathTransformFunc
	}

	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}

	return &Store{StoreOpts: opts}
}

func (s *Store) Has(key string) bool {
	pathkey := s.PathTransformFunc(key)
	pathkeywithfoldername := s.Root + "/" + pathkey.FullPath()
	_, err := os.Stat(pathkeywithfoldername)
	return !errors.Is(err, os.ErrNotExist)
}

func (p *PathKey) GetPathFolderName() string {

	paths := strings.Split(p.Pathname, "/")

	if len(paths) == 0 {
		return ""
	}

	return paths[0]
}

func (s *Store) Delete(key string) error {

	pathkey := s.PathTransformFunc(key)
	return os.RemoveAll(s.Root + "/" + pathkey.GetPathFolderName())
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
	pathKeyWithRoot := s.Root + "/" + pathkey.FullPath()
	return os.Open(pathKeyWithRoot)
}

func (s *Store) writeStream(key string, r io.Reader) error {

	pathKey := s.PathTransformFunc(key)
	if err := os.MkdirAll(s.Root+"/"+pathKey.Pathname, os.ModePerm); err != nil {
		return err
	}
	pathFilename := pathKey.FullPath()

	f, err := os.Create(s.Root + "/" + pathFilename)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	fmt.Printf("wrote %d bytes to disk : %s\n ", n, s.Root+"/"+pathFilename)
	return nil
}
