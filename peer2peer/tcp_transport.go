package peer2peer

import (
	"net"
	"sync"
)

type TCPTransport struct {
	listenAddress string
	listener      net.Listener

	mu sync.RWMutex

	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddress string) *TCPTransport {

	return &TCPTransport{
		listenAddress: listenAddress,
	}
}
