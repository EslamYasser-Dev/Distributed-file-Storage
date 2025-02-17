package p2p

import "net"

// the remote node
type Peer interface {
	net.Conn
	Send([]byte) error
}

// handles communication
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
