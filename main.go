package main

import (
	"distributed-file-storage/peer2peer"
	"fmt"
	"log"
)

func onpeer(peer peer2peer.Peer) error {

	peer.Close()
	return fmt.Errorf("failed the onpeer func")
}

func main() {

	tcpopts := peer2peer.TCPTransportOptions{
		ListenAddress: ":3000",
		Handshakefunc: peer2peer.NOPhandshakeFunc,
		Decoder:       peer2peer.DefaultDeocoder{},
		OnPeer:        onpeer,
	}

	tr := peer2peer.NewTCPTransport(tcpopts)

	go func() {

		for rpc := range tr.Consume() {
			fmt.Printf("rpc : %+v\n  ", rpc)
		}
	}()

	err := tr.ListenAndAccept()

	if err != nil {
		log.Fatal(err)
	}

	select {}

	// fmt.Println("Hello, World!")
}
