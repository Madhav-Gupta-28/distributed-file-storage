package main

import (
	"bytes"
	"distributed-file-storage/peer2peer"
	"encoding/gob"
	"fmt"
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
			reader := bytes.NewReader(rpc.Payload)
			if err := gob.NewDecoder(reader).Decode(&p); err != nil {
				log.Printf("Failed to decode message %v", err)
				continue
			}

			if err := fs.handleMessage(rpc.From, &p); err != nil {
				log.Printf("Failed to handle message %v", err)
				continue
			}
		case <-fs.quitchan:
			return
		}
	}
}

type Message struct {
	Payload any
}

type MessageStoreFile struct {
	Key  string
	Size int64
}

func (fs *FileServer) handleMessage(from string, p *Message) error {
	switch v := p.Payload.(type) {
	case MessageStoreFile:
		return fs.handleMessageStoreFile(from, v)
	case MsgGetFile:
		return fs.handleMessageGetFile(from, v)

	default:
		fmt.Printf("Received unknown message: %+v\n", v)
	}
	return nil
}

func (fs *FileServer) handleMessageGetFile(from string, p MsgGetFile) error {
	if !fs.store.Has(p.Key) {
		return fmt.Errorf("file not found")
	}
	r, err := fs.store.Read(p.Key)
	if err != nil {
		return err
	}

	peer, ok := fs.peers[from]
	if !ok {
		return fmt.Errorf("peer not found")
	}

	n, err := io.Copy(peer, r)
	if err != nil {
		return err
	}
	fmt.Println("written to the the network ", n)
	return nil
}

func (fs *FileServer) handleMessageStoreFile(from string, p MessageStoreFile) error {

	peer, ok := fs.peers[from]
	if !ok {
		return fmt.Errorf("peer not found")
	}
	if n, err := fs.store.Write(p.Key, io.LimitReader(peer, p.Size)); err != nil {
		fmt.Println("received and written to the disk", n)
		return err
	}
	peer.(*peer2peer.TCPPeer).Wg.Done()
	return nil
}

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

func (fs *FileServer) stream(p *Message) error {
	peers := []io.Writer{}
	for _, peer := range fs.peers {
		peers = append(peers, peer)
	}
	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(p)
}

func (fs *FileServer) broadcast(msg *Message) error {

	Buf := new(bytes.Buffer)
	if err := gob.NewEncoder(Buf).Encode(msg); err != nil {
		return err
	}
	for _, peer := range fs.peers {
		peer.Send([]byte{peer2peer.IncomingMessage})
		if err := peer.Send(Buf.Bytes()); err != nil {
			return err
		}
	}

	return nil

}

type MsgGetFile struct {
	Key string
}

func (fs *FileServer) Get(key string) (io.Reader, error) {

	if fs.store.Has(key) {
		return fs.store.Read(key)
	}

	msg := Message{
		Payload: MsgGetFile{
			Key: key,
		},
	}

	if err := fs.broadcast(&msg); err != nil {
		return nil, err
	}

	time.Sleep(1 * time.Millisecond)

	for _, peer := range fs.peers {
		filebuf := new(bytes.Buffer)

		n, err := io.CopyN(filebuf, peer, 10)
		if err != nil {
			return nil, err
		}
		fmt.Println("received and written to the disk", n)
	}

	select {}

	return nil, nil
}

func (fs *FileServer) Store(key string, data io.Reader) error {
	fileBuffer := new(bytes.Buffer)
	tee := io.TeeReader(data, fileBuffer)
	n, err := fs.store.Write(key, tee)
	if err != nil {
		return err
	}
	msg := Message{
		Payload: MessageStoreFile{
			Key:  key,
			Size: n,
		},
	}

	if err := fs.broadcast(&msg); err != nil {
		return err
	}

	time.Sleep(500 * time.Millisecond)
	// Sending a large file to the peers
	for _, peer := range fs.peers {
		peer.Send([]byte{peer2peer.IncomingStream})
		n, err := io.Copy(peer, fileBuffer)
		if err != nil {
			return err
		}
		fmt.Println("received and written to the disk", n)
	}
	return nil
}

func init() {

	gob.Register(MessageStoreFile{})
	gob.Register(MsgGetFile{})
}
