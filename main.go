package main

import (
	"filestorage/p2p"
	"fmt"
	"log"
)

func onPeer(peer p2p.Peer) error {
	peer.Close()
	return nil
}

func main() {

	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":5555",
		HandShakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefualtDecoder{},
		OnPeer:        onPeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("msg : %+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
