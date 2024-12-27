package peer2peer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {

	opts := TCPTransportOptions{
		ListenAddress: ":3000",
		Decoder:       GOBDecoder{},
		Handshakefunc: NOPhandshakeFunc,
	}

	tr := NewTCPTransport(opts)

	assert.Equal(t, tr.ListenAddress, ":3000")

	// Server

	assert.Equal(t, tr.ListenAndAccept(), nil)

}
