package rpc

import "fmt"

// LocalRPC is a simple in-memory RPC backend.
// It is NOT distributed — only for testing.
type LocalRPC struct {
	last RPCMessage
}

func NewLocalRPC() *LocalRPC {
	return &LocalRPC{}
}

func (l *LocalRPC) Send(msg RPCMessage) error {
	fmt.Println("local rpc: send", msg.Type)
	l.last = msg
	return nil
}

func (l *LocalRPC) Receive() (RPCMessage, error) {
	fmt.Println("local rpc: receive", l.last.Type)
	return l.last, nil
}
