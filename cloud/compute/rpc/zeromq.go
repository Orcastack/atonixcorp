package rpc

import "fmt"

type ZeroMQRPC struct{}

func NewZeroMQRPC() *ZeroMQRPC {
	return &ZeroMQRPC{}
}

func (z *ZeroMQRPC) Send(msg RPCMessage) error {
	fmt.Println("zeromq: send", msg.Type)
	return nil
}

func (z *ZeroMQRPC) Receive() (RPCMessage, error) {
	fmt.Println("zeromq: receive")
	return RPCMessage{Type: "zeromq.message"}, nil
}
