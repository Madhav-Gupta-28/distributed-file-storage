package main

import (
	"distributed-file-storage/peer2peer"
	"fmt"
	"log"
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
		case msg := <-fs.Transport.Consume():
			fmt.Printf("Received message: %v\n", msg)
		case <-fs.quitchan:
			return
		}
	}
}

func (fs *FileServer) bootstrapNetwork() error {
	for _, addr := range fs.BootstrapNodes {
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
