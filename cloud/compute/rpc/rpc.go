package rpc

type RPC interface {
	Send(msg RPCMessage) error
	Receive() (RPCMessage, error)
}
