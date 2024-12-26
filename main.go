package main

import (
	"distributed-file-storage/peer2peer"
	"log"
)

func main() {

	tcpopts := peer2peer.TCPTransportOptions{
		ListenAddress: ":3000",
		Handshakefunc: peer2peer.NOPhandshakeFunc,
		Decoder:       peer2peer.DefaultDeocoder{},
	}

	tr := peer2peer.NewTCPTransport(tcpopts)

	err := tr.ListenAndAccept()

	if err != nil {
		log.Fatal(err)
	}

	select {}

	// fmt.Println("Hello, World!")
}
