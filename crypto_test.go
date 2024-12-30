package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCopyEncrypt(t *testing.T) {
	buf := new(bytes.Buffer)
	dst := bytes.NewReader([]byte("hello world"))
	key := NewEncryptKey()

	_, err := CopyEncrypt(key, dst, buf)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(buf.Bytes())

}
