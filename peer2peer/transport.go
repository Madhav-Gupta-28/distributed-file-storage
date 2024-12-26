package peer2peer

// Peer is an interface that represents the remote node
type Peer interface {
}

// Handles the communication between nodes in the network
type Transport interface {
	ListenAndAccept() error
}
