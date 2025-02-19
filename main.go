package main

import (
	"bytes"
	"filestorage/p2p"
	"log"
	"time"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandShakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefualtDecoder{},
		//onPeerFunc
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)
	opts := FileServerOpts{
		StorageRoot:       (listenAddr + "_network"),
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootStrapNodes:    nodes,
	}
	server := NewFileServer(opts)
	tcpTransport.OnPeer = server.OnPeer
	return server
}
func main() {

	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")
	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(time.Second * 4)

	go s2.Start()
	time.Sleep(time.Second * 4)

	data := bytes.NewReader([]byte("imagine that is a big file"))
	s2.StoreData("myData", data)
	select {}
}
