package main

import (
	"bytes"
	"encoding/gob"
	"filestorage/p2p"
	"fmt"
	"io"
	"log"
	"sync"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootStrapNodes    []string
}

type FileServer struct {
	FileServerOpts
	peerLock sync.Mutex
	peers    map[string]p2p.Peer
	store    *Store
	quitch   chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	StoreOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		peers:          make(map[string]p2p.Peer),
		store:          NewStore(StoreOpts),
		quitch:         make(chan struct{}),
	}
}

type Message struct {
	Payload any
}

func (s *FileServer) broadcast(msg *Message) error {
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}
	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)
}
func (s *FileServer) StoreData(key string, r io.Reader) error {
	// buff := new(bytes.Buffer)

	// tee := io.TeeReader(r, buff)

	// if err := s.store.Write(key, tee); err != nil {
	// 	return err
	// }

	// p := &DataMessage{
	// 	Key:  key,
	// 	Data: buff.Bytes(),
	// }
	// return s.broadcast(&Message{
	// 	From:    "todo",
	// 	Payload: p,
	// })

	buff := new(bytes.Buffer)
	msg := &Message{
		Payload: []byte("some bytes to read"),
	}

	if err := gob.NewEncoder(buff).Encode(msg); err != nil {
		return err
	}
	for _, peer := range s.peers {
		if err := peer.Send(buff.Bytes()); err != nil {
			return err
		}
	}
	payload := []byte("large file")
	for _, peer := range s.peers {
		if err := peer.Send(payload); err != nil {
			return err
		}
	}
	return nil
}
func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[p.RemoteAddr().String()] = p
	log.Printf("Connected with remote: %s", p.RemoteAddr())
	return nil
}
func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due to user quit action")
		s.Transport.Close()
	}()
	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Fatal(err)
			}
			peer, ok := s.peers[rpc.From]
			if !ok {
				panic("peer not found")
			}
			b := make([]byte, 1024)
			if _, err := peer.Read(b); err != nil {
				panic(err)
			}
			panic("err ")
			fmt.Printf("recv Message : %s", string(msg.Payload.([]byte)))
			// if err := s.handleMessage(&m); err != nil {
			// 	log.Println("hand le message error", err)
			// }
		case <-s.quitch:
			return
		}
	}
}

//	func (s *FileServer) handleMessage(msg *Message) error {
//		switch v := msg.Payload.(type) {
//		case *DataMessage:
//			fmt.Printf("received msg from %+v\n", v)
//		}
//		return nil
//	}
func (s *FileServer) BootStrapNetwork() error {

	for _, addr := range s.BootStrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			fmt.Println("Attempting to connect with remote: ", addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error", err)

			}
		}(addr)
	}

	return nil
}
func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	if len(s.BootStrapNodes) != 0 {
		s.BootStrapNetwork()
	}
	s.loop()
	return nil
}
