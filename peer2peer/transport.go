package peer2peer

import "net"

// Peer is an interface that represents the remote node
type Peer interface {
	Send([]byte) error
	net.Conn
	CloseStream()
}

// Handles the communication between nodes in the network
type Transport interface {
	Addr() string
	ListenAndAccept() error
	Dial(addr string) error
	Consume() <-chan RPC
	Close() error
	ListenAddr() string
}
