package peer2peer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {

	tr := NewTCPTransport(":3000")

	assert.Equal(t, tr.listenAddress, ":3000")

	// Server

	assert.Equal(t, tr.ListenAndAccept(), nil)

	select {}
}
