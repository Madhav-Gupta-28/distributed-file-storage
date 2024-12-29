package main

import (
	"bytes"
	"distributed-file-storage/peer2peer"
	"log"
	"time"
)

func makeServer(addr string, root string, nodes ...string) *FileServer {

	tcptransportopts := peer2peer.TCPTransportOptions{
		ListenAddress: addr,
		Handshakefunc: peer2peer.NOPhandshakeFunc,
		Decoder:       peer2peer.DefaultDeocoder{},
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

	tcptransport.OnPeer = s.onPeer

	return s
}

func main() {

	s1 := makeServer(":3000", "")

	s2 := makeServer(":4000", "", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(1 * time.Second)

	go s2.Start()

	time.Sleep(1 * time.Second)
	data := bytes.NewReader([]byte("hello world I ambuildinng smth"))
	s2.StoreData("madhavgupta", data)

	select {}
}
