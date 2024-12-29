package peer2peer

// Peer is an interface that represents the remote node
type Peer interface {
	Close() error
}

// Handles the communication between nodes in the network
type Transport interface {
	ListenAndAccept() error
	Dial(addr string) error
	Consume() <-chan RPC
	Close() error
}
