package queue

import (
	"time"
)

type Logger interface {
	Info(...any)
	Error(...any)
}

type OutgoingMessage struct {
	ID        string
	From      string
	To        []string
	Subject   string
	Body      []byte
	Raw       []byte
	Tenant    string
	Attempts  int
	MaxRetry  int
	NextRetry time.Time
}

type Store interface {
	SaveOutgoing(msg *OutgoingMessage) error
	ListPending(limit int) ([]OutgoingMessage, error)
	MarkDelivered(id string) error
	MarkFailed(id string, err error) error
	UpdateRetry(id string, attempts int, next time.Time) error
}

type Queue struct {
	store  Store
	logger Logger
}

func New(store Store, logger Logger) *Queue {
	return &Queue{store: store, logger: logger}
}

func (q *Queue) EnqueueOutgoing(msg *OutgoingMessage) error {
	msg.Attempts = 0
	msg.MaxRetry = 5
	msg.NextRetry = time.Now()

	q.logger.Info("Queue: enqueue outgoing message ", msg.ID)
	return q.store.SaveOutgoing(msg)
}
