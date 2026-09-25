package imap

// Message is the IMAP view of a message (minimal for FETCH BODY[])
type Message struct {
	ID   string
	Body []byte
}
