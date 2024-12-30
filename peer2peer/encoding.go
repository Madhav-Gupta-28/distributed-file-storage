package peer2peer

import (
	"encoding/gob"
	"fmt"
	"io"
)

// Decoder is an interface that decodes a RPC
type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct{}

func (dec GOBDecoder) Decode(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}

type DefaultDeocoder struct{}

func (dec DefaultDeocoder) Decode(r io.Reader, msg *RPC) error {

	peekBuf := make([]byte, 1)

	if _, err := r.Read(peekBuf); err != nil {
		return err
	}

	// In case of the stream we are not decoding what is being sent over the network
	// we are just setting the stream flag to true and returning  p
	stream := peekBuf[0] == IncomingStream
	if stream {
		msg.Stream = true
		return nil
	}

	buff := make([]byte, 1028)
	n, err := r.Read(buff)

	if err != nil {
		fmt.Printf("Error reading from reader: %v\n", err)
	}

	msg.Payload = buff[:n]
	fmt.Printf("Read %v bytes from reader\n", string(msg.Payload))
	return nil
}
