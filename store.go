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
	Root string

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

func (s *Store) Has(id string, key string) bool {
	pathkey := s.PathTransformFunc(key)
	pathkeywithfoldername := s.Root + "/" + id + "/" + pathkey.FullPath()
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

func (s *Store) Delete(id string, key string) error {
	pathkey := s.PathTransformFunc(key)
	return os.RemoveAll(s.Root + "/" + id + "/" + pathkey.GetPathFolderName())
}

func (s *Store) ClearRoot() error {
	return os.RemoveAll(s.Root)
}

// Public Function of wirte stream
func (s *Store) Write(id string, key string, r io.Reader) (int64, error) {
	return s.writeStream(id, key, r)
}

// Public function of read stream
func (s *Store) Read(id string, key string) (int64, io.Reader, error) {
	return s.readStream(id, key)
}

func (s *Store) readStream(id string, key string) (int64, io.ReadCloser, error) {

	pathkey := s.PathTransformFunc(key)
	pathKeyWithRoot := s.Root + "/" + id + "/" + pathkey.FullPath()

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

func (s *Store) WriteDecrypt(id string, encKey []byte, key string, r io.Reader) (int64, error) {

	f, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err
	}
	n, err := CopyDecrypt(encKey, r, f)
	if err != nil {
		return 0, err
	}
	return int64(n), nil

}

func (s *Store) openFileForWriting(id string, key string) (*os.File, error) {
	pathKey := s.PathTransformFunc(key)
	if err := os.MkdirAll(s.Root+"/"+id+"/"+pathKey.Pathname, os.ModePerm); err != nil {
		return nil, err
	}
	pathFilename := pathKey.FullPath()

	return os.Create(s.Root + "/" + id + "/" + pathFilename)

}

func (s *Store) writeStream(id string, key string, r io.Reader) (int64, error) {
	f, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err
	}
	return io.Copy(f, r)

}
