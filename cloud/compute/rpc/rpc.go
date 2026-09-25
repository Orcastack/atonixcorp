package rpc

import "fmt"

type RPC interface {
	Send(msg RPCMessage) error
	Receive() (RPCMessage, error)
}

type LocalRPC struct{}

func NewLocalRPC() *LocalRPC {
	return &LocalRPC{}
}

func (l *LocalRPC) Send(msg RPCMessage) error {
	fmt.Println("RPC send:", msg.Type)
	return nil
}

func (l *LocalRPC) Receive() (RPCMessage, error) {
	fmt.Println("RPC receive")
	return RPCMessage{Type: "noop"}, nil
}
