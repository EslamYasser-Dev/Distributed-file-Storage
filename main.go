package main

import (
	"filestorage/p2p"
	"log"
)

func main() {

	tr := p2p.NewTCPTransport(":5555")
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
