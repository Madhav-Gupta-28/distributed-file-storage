package peer2peer

import (
	"fmt"
	"log"
	"net"
	"sync"
)

// TCPPeer is a remote node over a established connection
type TCPPeer struct {
	// connection is the underlying connection of the peer
	connection net.Conn

	//if we dial a  connection  to a remote node, we are an outbound peer = true
	//if we accept a connection from a remote node, we are an inbound peer = false
	outbound bool
}

func NewTCPPeer(connection net.Conn, outbound bool) *TCPPeer {

	return &TCPPeer{
		connection: connection,
		outbound:   outbound,
	}
}

type TCPTransportOptions struct {
	ListenAddress string
	Handshakefunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOptions
	listener net.Listener
	mu       sync.RWMutex

	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOptions) *TCPTransport {

	return &TCPTransport{
		TCPTransportOptions: opts,
	}
}

func (t *TCPTransport) ListenAndAccept() error {

	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddress)

	if err != nil {
		log.Printf("Error starting listener: %v\n", err)
		return err
	}

	log.Printf("Listening on %s\n", t.ListenAddress)
	go t.startacceptLoop()

	return nil
}

func (t *TCPTransport) startacceptLoop() {

	for {
		connection, err := t.listener.Accept()

		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		log.Printf("Accepted new connection from %s\n", connection.RemoteAddr())
		go t.handleConnection(connection)

	}
}

type Temp struct{}

func (t *TCPTransport) handleConnection(connection net.Conn) {
	peer := NewTCPPeer(connection, false)

	if err := t.Handshakefunc(peer); err != nil {
		connection.Close()
		fmt.Printf("Error shaking hands with peer: %v\n", err)
		return

	}

	// Read loop
	msg := &Message{}

	for {

		if err := t.Decoder.Decode(connection, msg); err != nil {
			fmt.Printf("Error decoding message: %v\n", err)
			continue

		}

		msg.from = connection.RemoteAddr()

		fmt.Printf("Message : %v\n  ", msg)
	}

}
