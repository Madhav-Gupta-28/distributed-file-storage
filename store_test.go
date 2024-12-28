package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {

	key := "fuck"
	pathname := CASPathTransformFunc(key)
	fmt.Println(pathname)
}

func TestStore(t *testing.T) {

	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	store := NewStore(opts)

	data := []byte("fuck the world")

	err := store.writeStream("testingpujn", bytes.NewReader(data))

	if err != nil {
		t.Fatalf("failed to write stream: %v", err)
	}

	if ok := store.Has("testingpujn"); !ok {
		t.Errorf("Expected to have an key  testingpujn")
	}

	reader, err := store.Read("testingpujn")

	fmt.Printf("reader %+v", reader)

	if err != nil {
		t.Fatalf("failed to read stream: %v", err)
	}

	b, _ := ioutil.ReadAll(reader)

	if string(b) != string(data) {
		t.Errorf("data mismatch -> want %s , got %s", string(data), string(b))

	}

}

func TestDelete(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
		Root:              defaultRootFolderName,
	}
	store := NewStore(opts)
	store.Delete("testingpujn")
}
