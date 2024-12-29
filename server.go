package main

import (
	"bytes"
	"distributed-file-storage/peer2peer"
	"encoding/gob"
	"io"
	"log"
	"sync"
	"time"
)

type FileServerOptions struct {
	ListenAddress     string
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         peer2peer.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOptions
	store    *Store
	quitchan chan struct{}
	peers    map[string]peer2peer.Peer
	Peerlock sync.Mutex
}

func NewFileServer(fileServerOptions FileServerOptions) *FileServer {
	storeOpts := StoreOpts{
		Root:              fileServerOptions.StorageRoot,
		PathTransformFunc: fileServerOptions.PathTransformFunc,
	}
	return &FileServer{
		FileServerOptions: fileServerOptions,
		store:             NewStore(storeOpts),
		quitchan:          make(chan struct{}),
		peers:             make(map[string]peer2peer.Peer),
	}
}

func (fs *FileServer) Stop() {

	close(fs.quitchan)

}

func (fs *FileServer) loop() {

	defer func() {
		log.Println("FileServer loop exited user quirt action")
		fs.Transport.Close()
	}()

	for {
		select {
		case rpc := <-fs.Transport.Consume():
			var p Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&p); err != nil {
				log.Fatal("Failed to decode message", err)
			}

			peer, ok := fs.peers[rpc.From]
			if !ok {
				log.Println("Peer not found", rpc.From)
			}

			buf := make([]byte, 1024)
			if _, err := peer.Read(buf); err != nil {
				panic(err)
			}
			peer.(*peer2peer.TCPPeer).Wg.Done()

			// err := fs.handleMessage(&p)
			// if err != nil {
			// 	log.Println("Failed to handle message", err)
			// }
		case <-fs.quitchan:
			return
		}
	}
}

type Message struct {
	Payload any
}

// func (fs *FileServer) handleMessage(msg *Message) error {

// 	switch v := msg.Payload.(type) {
// 	case *MessageData:
// 		fmt.Println("Received MessageData", v)
// 	}

// 	return nil
// }

func (fs *FileServer) bootstrapNetwork() error {
	for _, addr := range fs.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func() {
			if err := fs.Transport.Dial(addr); err != nil {
				log.Println("Failed to dial bootstrap node", err)
			}
		}()
	}
	return nil
}

func (fs *FileServer) Start() error {

	if err := fs.Transport.ListenAndAccept(); err != nil {
		return err
	}

	fs.bootstrapNetwork()
	fs.loop()
	return nil

}

func (fs *FileServer) onPeer(peer peer2peer.Peer) error {

	fs.Peerlock.Lock()
	defer fs.Peerlock.Unlock()

	fs.peers[peer.RemoteAddr().String()] = peer

	log.Println("New peer connected", peer)

	return nil
}

func (fs *FileServer) broadcast(p *Message) error {
	peers := []io.Writer{}
	for _, peer := range fs.peers {
		peers = append(peers, peer)
	}
	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(p)
}

func (fs *FileServer) StoreData(key string, data io.Reader) error {

	// // creating a buffer to store the data and a tee reader to read the data and then passing it to the MessageData
	buf := new(bytes.Buffer)
	// tee := io.TeeReader(data, buf)

	msg := &Message{
		Payload: []byte("Hello"),
	}

	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	for _, peer := range fs.peers {
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}

	time.Sleep(3 * time.Second)

	payload := []byte("Large File ------")

	for _, peer := range fs.peers {
		if err := peer.Send(payload); err != nil {
			return err
		}

	}

	return nil

	// // Store the data to the
	// if err := fs.store.Write(key, tee); err != nil {
	// 	return err
	// }
	// p := &MessageData{key: key, Data: buf.Bytes()}

	// // fmt.Println("MessageData", p)

	// return fs.broadcast(&Message{
	// 	From:    fs.ListenAddress,
	// 	Payload: p,
	// })
}
