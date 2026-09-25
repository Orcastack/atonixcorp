package storage

type Mailbox struct {
	Name string
}

type FolderStore interface {
	ListMailboxes(userID string) ([]Mailbox, error)
	CreateMailbox(userID, name string) error
	DeleteMailbox(userID, name string) error
}
