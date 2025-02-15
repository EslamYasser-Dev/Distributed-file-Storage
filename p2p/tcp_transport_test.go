package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr:    ":5556",
		HandShakeFunc: NOPHandShakeFunc,
		Decoder:       DefualtDecoder{},
	}
	tr := NewTCPTransport(opts)
	assert.Equal(t, tr.ListenAddr, ":5556")

	assert.Nil(t, tr.ListenAndAccept())
}
