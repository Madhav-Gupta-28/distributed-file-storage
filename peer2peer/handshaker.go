package peer2peer

// HandshakeFunc is a function that performs a handshake with a peer
type HandshakeFunc func(Peer) error

// NOPhandshakeFunc is a no-op handshake function
func NOPhandshakeFunc(Peer) error { return nil }
