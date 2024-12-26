package peer2peer

import (
	"encoding/gob"
	"fmt"
	"io"
)

// Decoder is an interface that decodes a message
type Decoder interface {
	Decode(io.Reader, *Message) error
}

type GOBDecoder struct {
}

func (dec GOBDecoder) Decode(r io.Reader, msg *Message) error {

	return gob.NewDecoder(r).Decode(msg)

}

type DefaultDeocoder struct{}

func (dec DefaultDeocoder) Decode(r io.Reader, msg *Message) error {
	buff := make([]byte, 1028)

	n, err := r.Read(buff)

	if err != nil {
		fmt.Printf("Error reading from reader: %v\n", err)
	}

	msg.Payload = buff[:n]

	fmt.Printf("Read %v bytes from reader\n", string(msg.Payload))

	return nil
}
