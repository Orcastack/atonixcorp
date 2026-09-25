package storage

import "time"

// Canonical message used by SMTP, POP3, API, Queue, Storage
type Message struct {
	ID      string
	UserID  string
	Mailbox string

	Subject string
	From    string
	To      string

	Date time.Time
	Body []byte
	Raw  []byte
	Size int
}
