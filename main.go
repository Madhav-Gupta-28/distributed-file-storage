package main

import (
	"distributed-file-storage/peer2peer"
	"log"
)

func makeServer(addr string, root string, nodes ...string) *FileServer {

	tcptransportopts := peer2peer.TCPTransportOptions{
		ListenAddress: addr,
		Handshakefunc: peer2peer.NOPhandshakeFunc,
		Decoder:       peer2peer.DefaultDeocoder{},
		// OnPeer:        OnPeer,
	}

	tcptransport := peer2peer.NewTCPTransport(tcptransportopts)

	FileServeroptions := FileServerOptions{
		ListenAddress:     addr,
		StorageRoot:       root,
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcptransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(FileServeroptions)

	return s
}

func main() {

	s1 := makeServer(":3000", "")

	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	s2.Start()

}
