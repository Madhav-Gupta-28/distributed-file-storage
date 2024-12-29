package peer2peer

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

// TCPPeer is a remote node over a established connection
type TCPPeer struct {
	// the underlying connection of the peer which is this care is tcp connection
	net.Conn
	//if we dial a  connection  to a remote node, we are an outbound peer = true
	//if we accept a connection from a remote node, we are an inbound peer = false
	outbound bool
	Wg       *sync.WaitGroup
}

func (p *TCPPeer) Send(data []byte) error {
	_, err := p.Conn.Write(data)
	return err
}

func NewTCPPeer(connection net.Conn, outbound bool) *TCPPeer {

	return &TCPPeer{
		Conn:     connection,
		outbound: outbound,
		Wg:       &sync.WaitGroup{},
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

func (t *TCPTransport) ListenAddr() string {
	return t.ListenAddress
}

func (t *TCPTransport) ListenAndAccept() error {

	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddress)

	if err != nil {
		log.Printf("Error starting listener: %v\n", err)
		return err
	}

	go t.startacceptLoop()

	log.Printf(" TCP Transport Listening on %s\n", t.ListenAddress)

	return nil
}

func (t *TCPTransport) startacceptLoop() {

	for {
		connection, err := t.listener.Accept()

		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
		}
		go t.handleConnection(connection, false)
	}
}

func (t *TCPTransport) handleConnection(connection net.Conn, outbound bool) {
	var err error
	defer func() {
		fmt.Printf("Closing connection from %s\n", err)
		connection.Close()
	}()
	peer := NewTCPPeer(connection, outbound)
	if err := t.Handshakefunc(peer); err != nil {
		return
	}
	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}
	// Read loop
	for {
		rpc := RPC{}
		err = t.Decoder.Decode(connection, &rpc)
		if err != nil {
			log.Printf("Decoding failed: %v\n", err)
			return
		}
		rpc.From = connection.RemoteAddr().String()
		peer.Wg.Add(1)
		fmt.Println("Waiting till stream is done ")

		t.rpcchan <- rpc
		peer.Wg.Wait()
		fmt.Println("stream done continuing normal real loop")
	}
}

// Close Function Implements the Transport interface
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial Function Implements the Transport interface
func (t *TCPTransport) Dial(addr string) error {

	connection, err := net.Dial("tcp", addr)

	if err != nil {
		return err
	}

	go t.handleConnection(connection, true)

	return nil
}
