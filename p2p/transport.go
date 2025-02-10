package p2p

// the remote node
type Peer interface {
}

// handles communication
type Transport interface {
	ListenAndAccept() error
}
