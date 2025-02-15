package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct {
}

func (dec GOBDecoder) Decode(r io.Reader, v *RPC) error {
	return gob.NewDecoder(r).Decode(v)
}

type DefualtDecoder struct {
}

func (dec DefualtDecoder) Decode(r io.Reader, msg *RPC) error {
	buff := make([]byte, 2048)
	n, err := r.Read(buff)
	if err != nil {
		return err
	}
	msg.Payload = buff[:n]
	return nil
}
