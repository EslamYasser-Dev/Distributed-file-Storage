package p2p

import (
	"fmt"
	"net"
)

type TCPPeer struct {
	conn net.Conn
	// if we dial and retreive connection => outbound == true
	//if accept and retrive connection  => outbound == false
	outBound bool
}

func NewTCPPeer(conn net.Conn, outBound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outBound: outBound,
	}
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandShakeFunc HandShakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}
type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcChan  chan RPC

	// mu    sync.RWMutex
	// peers map[*net.Addr]Peer
}

func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcChan
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcChan:          make(chan RPC),
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Println("tcp tranport Accept Error : ", err)
		}
		fmt.Printf("new incomming connection: %+v\n", conn)
		go t.handleConn(conn)
	}
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error
	defer func() {
		fmt.Printf("Dropping peer connection %s", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)
	if err = t.HandShakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}

	}

	rpc := RPC{}
	for {
		err := t.Decoder.Decode(conn, &rpc)
		if err != nil {
			return
		}
		rpc.From = conn.RemoteAddr()
		t.rpcChan <- rpc
	}
}
