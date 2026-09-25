package rpc

import "fmt"

type GRPCRPC struct{}

func NewGRPCRPC() *GRPCRPC {
	return &GRPCRPC{}
}

func (g *GRPCRPC) Send(msg RPCMessage) error {
	fmt.Println("grpc: send", msg.Type)
	return nil
}

func (g *GRPCRPC) Receive() (RPCMessage, error) {
	fmt.Println("grpc: receive")
	return RPCMessage{Type: "grpc.message"}, nil
}
