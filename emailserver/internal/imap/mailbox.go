package imap

import "time"

// Mailbox represents an IMAP mailbox (folder)
type Mailbox struct {
	Name        string
	UIDValidity uint32
	UIDNext     uint32
}

// Message is the unified mail object used by SMTP, IMAP, POP3, API, and Queue
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

// MailboxStore defines storage operations for IMAP mailboxes
type MailboxStore interface {
	ListMailboxes(userID string) ([]Mailbox, error)
	GetMailbox(userID, name string) (*Mailbox, error)
}
