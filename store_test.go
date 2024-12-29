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

func newStore() *Store {
	opts := &StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
		Root:              defaultRootFolderName,
	}
	s := NewStore(*opts)
	return s
}

func TestStore(t *testing.T) {

	store := newStore()
	defer store.ClearRoot()

	for i := 0; i < 5; i++ {

		key := fmt.Sprintf("madhav_%d", i)
		data := []byte("fuck the world")

		err := store.writeStream(key, bytes.NewReader(data))
		if err != nil {
			t.Fatalf("failed to write stream: %v", err)
		}

		if ok := store.Has(key); !ok {
			t.Errorf("Expected to have an key ")
		}

		reader, err := store.Read(key)
		fmt.Printf("reader %+v", reader)
		if err != nil {
			t.Fatalf("failed to read stream: %v", err)
		}
		b, _ := ioutil.ReadAll(reader)
		if string(b) != string(data) {
			t.Errorf("data mismatch -> want %s , got %s", string(data), string(b))
		}

	}
}
