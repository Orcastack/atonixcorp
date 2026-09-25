package queue

import (
	"time"
)

type Worker struct {
	queue  *Queue
	sender SMTPSender
	logger Logger
}

type SMTPSender interface {
	Send(msg OutgoingMessage) error
}

func NewWorker(queue *Queue, sender SMTPSender, logger Logger) *Worker {
	return &Worker{queue: queue, sender: sender, logger: logger}
}

func (w *Worker) Start() {
	w.logger.Info("Queue worker started")

	for {
		msgs, err := w.queue.store.ListPending(20)
		if err != nil {
			w.logger.Error("Queue worker: list pending:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, msg := range msgs {
			if time.Now().Before(msg.NextRetry) {
				continue
			}

			w.process(msg)
		}

		time.Sleep(2 * time.Second)
	}
}

func (w *Worker) process(msg OutgoingMessage) {
	w.logger.Info("Queue worker: delivering message ", msg.ID)

	err := w.sender.Send(msg)
	if err == nil {
		w.logger.Info("Queue worker: delivered ", msg.ID)
		w.queue.store.MarkDelivered(msg.ID)
		return
	}

	w.logger.Error("Queue worker: delivery failed: ", err)

	msg.Attempts++
	if msg.Attempts >= msg.MaxRetry {
		w.queue.store.MarkFailed(msg.ID, err)
		return
	}

	// exponential backoff
	backoff := time.Duration(msg.Attempts*msg.Attempts) * time.Minute
	next := time.Now().Add(backoff)

	w.queue.store.UpdateRetry(msg.ID, msg.Attempts, next)
}
