package rpc

type RabbitRPC struct {
	ch    *amqp.Channel
	queue string
}

func NewRabbitRPC(url, queue string) (*RabbitRPC, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = ch.QueueDeclare(queue, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return &RabbitRPC{ch: ch, queue: queue}, nil
}

func (r *RabbitRPC) Send(msg RPCMessage) error {
	return r.ch.Publish("", r.queue, false, false,
		amqp.Publishing{Body: []byte(msg.Type)})
}

func (r *RabbitRPC) Receive() (RPCMessage, error) {
	msgs, err := r.ch.Consume(r.queue, "", true, false, false, false, nil)
	if err != nil {
		return RPCMessage{}, err
	}
	m := <-msgs
	return RPCMessage{Type: string(m.Body)}, nil
}
