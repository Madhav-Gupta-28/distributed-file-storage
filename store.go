package main

import (
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

func (s *Store) ClearRoot() error {
	return os.RemoveAll(s.Root)
}

// Public Function of wirte stream
func (s *Store) Write(key string, r io.Reader) (int64, error) {
	return s.writeStream(key, r)
}

// Public function of read stream
func (s *Store) Read(key string) (int64, io.Reader, error) {
	return s.readStream(key)
}

func (s *Store) readStream(key string) (int64, io.ReadCloser, error) {

	pathkey := s.PathTransformFunc(key)
	pathKeyWithRoot := s.Root + "/" + pathkey.FullPath()

	file, err := os.Open(pathKeyWithRoot)

	if err != nil {
		return 0, nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return 0, nil, err
	}

	return fi.Size(), file, nil
}

func (s *Store) writeStream(key string, r io.Reader) (int64, error) {

	pathKey := s.PathTransformFunc(key)
	if err := os.MkdirAll(s.Root+"/"+pathKey.Pathname, os.ModePerm); err != nil {
		return 0, err
	}
	pathFilename := pathKey.FullPath()

	f, err := os.Create(s.Root + "/" + pathFilename)
	if err != nil {
		return 0, err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		return 0, err
	}
	return n, nil
}
