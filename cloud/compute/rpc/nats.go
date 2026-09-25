package rpc

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

type NatsRPC struct {
	nc      *nats.Conn
	subject string
}

func NewNatsRPC(url, subject string) (*NatsRPC, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsRPC{nc: nc, subject: subject}, nil
}

func (n *NatsRPC) Send(msg RPCMessage) error {
	data := fmt.Sprintf("%s", msg.Type)
	return n.nc.Publish(n.subject, []byte(data))
}

func (n *NatsRPC) Receive() (RPCMessage, error) {
	sub, err := n.nc.SubscribeSync(n.subject)
	if err != nil {
		return RPCMessage{}, err
	}
	m, err := sub.NextMsg(0)
	if err != nil {
		return RPCMessage{}, err
	}
	return RPCMessage{Type: string(m.Data)}, nil
}
