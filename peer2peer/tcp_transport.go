package peer2peer

import (
	"fmt"
	"log"
	"net"
)

// TCPPeer is a remote node over a established connection
type TCPPeer struct {
	// connection is the underlying connection of the peer
	connection net.Conn

	//if we dial a  connection  to a remote node, we are an outbound peer = true
	//if we accept a connection from a remote node, we are an inbound peer = false
	outbound bool
}

func (p *TCPPeer) Close() error {
	return p.connection.Close()
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
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOptions
	listener net.Listener
	rpcchan  chan RPC
}

func NewTCPTransport(opts TCPTransportOptions) *TCPTransport {

	return &TCPTransport{
		TCPTransportOptions: opts,
		rpcchan:             make(chan RPC),
	}
}

func (t *TCPTransport) Consume() <-chan RPC {

	return t.rpcchan

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

func (t *TCPTransport) handleConnection(connection net.Conn) {

	var err error

	defer func() {

		fmt.Printf("Closing connection from %s\n", err)
		connection.Close()

	}()

	peer := NewTCPPeer(connection, false)

	if err := t.Handshakefunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {

		if err = t.OnPeer(peer); err != nil {
			return
		}

	}

	// Read loop
	rpc := RPC{}

	for {

		err = t.Decoder.Decode(connection, &rpc)

		if err != nil {
			return
		}

		rpc.from = connection.RemoteAddr()

		t.rpcchan <- rpc

		// fmt.Printf("rpc : %v\n  ", rpc)
	}

}
