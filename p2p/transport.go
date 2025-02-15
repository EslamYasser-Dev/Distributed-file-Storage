package p2p

// the remote node
type Peer interface {
	Close() error
}

// handles communication
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
