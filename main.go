package main

import (
	"bytes"
	"distributed-file-storage/peer2peer"
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
		EncryptKey:        NewEncryptKey(),
	}

	s := NewFileServer(FileServeroptions)

	tcptransport.OnPeer = s.onPeer

	return s
}

func main() {

	s1 := makeServer(":3000", "madhavgupta", ":5000")

	s2 := makeServer(":4000", "madhavgupta2", ":3000")

	go s1.Start()

	time.Sleep(2 * time.Second)

	go s2.Start()

	time.Sleep(1 * time.Second)

	// for i := 0; i < 10; i++ {
	// 	s2.Store(fmt.Sprintf("madhavgupta%d", i), bytes.NewReader([]byte("hello world I ambuildinng smth")))
	// 	time.Sleep(5 * time.Millisecond)

	// }

	s2.Store("madhav", bytes.NewReader([]byte("hello world I ambuildinng smth")))
	time.Sleep(5 * time.Millisecond)

	// data, err := s2.Get("madhav")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// b, err := ioutil.ReadAll(data)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(b))
}
