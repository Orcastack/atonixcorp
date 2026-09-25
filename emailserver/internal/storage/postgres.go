package storage

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

// Save incoming SMTP message
func (s *PostgresStore) SaveIncoming(msg *Message) error {
	msg.ID = uuid.NewString()
	if msg.Date.IsZero() {
		msg.Date = time.Now()
	}
	msg.Size = len(msg.Raw)

	_, err := s.db.Exec(`
        INSERT INTO messages (id, user_id, mailbox, subject, from_addr, to_addr, date, body, raw, size)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
    `,
		msg.ID,
		msg.UserID,
		msg.Mailbox,
		msg.Subject,
		msg.From,
		msg.To,
		msg.Date,
		msg.Body,
		msg.Raw,
		msg.Size,
	)

	return err
}

// List messages for IMAP/POP3/API
func (s *PostgresStore) ListMessages(userID, mailbox string) ([]Message, error) {
	rows, err := s.db.Query(`
        SELECT id, user_id, mailbox, subject, from_addr, to_addr, date, body, raw, size
        FROM messages
        WHERE user_id=$1 AND mailbox=$2
        ORDER BY date DESC
    `, userID, mailbox)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.Mailbox,
			&m.Subject,
			&m.From,
			&m.To,
			&m.Date,
			&m.Body,
			&m.Raw,
			&m.Size,
		); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

func (s *PostgresStore) GetMessage(userID, mailbox, id string) (*Message, error) {
	var m Message
	err := s.db.QueryRow(`
        SELECT id, user_id, mailbox, subject, from_addr, to_addr, date, body, raw, size
        FROM messages
        WHERE user_id=$1 AND mailbox=$2 AND id=$3
    `, userID, mailbox, id).Scan(
		&m.ID,
		&m.UserID,
		&m.Mailbox,
		&m.Subject,
		&m.From,
		&m.To,
		&m.Date,
		&m.Body,
		&m.Raw,
		&m.Size,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *PostgresStore) DeleteMessage(userID, mailbox, id string) error {
	res, err := s.db.Exec(`
        DELETE FROM messages
        WHERE user_id=$1 AND mailbox=$2 AND id=$3
    `, userID, mailbox, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("not found")
	}
	return nil
}

func (s *PostgresStore) ListMailboxes(userID string) ([]Mailbox, error) {
	rows, err := s.db.Query(`
        SELECT name
        FROM mailboxes
        WHERE user_id=$1
        ORDER BY name
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Mailbox
	for rows.Next() {
		var m Mailbox
		if err := rows.Scan(&m.Name); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

func (s *PostgresStore) CreateMailbox(userID, name string) error {
	_, err := s.db.Exec(`
        INSERT INTO mailboxes (user_id, name)
        VALUES ($1, $2)
    `, userID, name)
	return err
}

func (s *PostgresStore) DeleteMailbox(userID, name string) error {
	_, err := s.db.Exec(`
        DELETE FROM mailboxes
        WHERE user_id=$1 AND name=$2
    `, userID, name)
	return err
}
