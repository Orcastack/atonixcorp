package domain

import (
	"database/sql"
)

type Domain struct {
	Name      string
	MX        string
	CreatedAt string
}

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) RegisterDomain(name string) error {
	_, err := s.db.Exec(`
        INSERT INTO domains (name, created_at)
        VALUES ($1, now())
        ON CONFLICT (name) DO NOTHING
    `, name)
	return err
}

func (s *Service) ListDomains() ([]Domain, error) {
	rows, err := s.db.Query(`
        SELECT name, mx, created_at
        FROM domains
        ORDER BY created_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Domain
	for rows.Next() {
		var d Domain
		if err := rows.Scan(&d.Name, &d.MX, &d.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}

// ValidateRecipientDomain is used by SMTP to check if we host the domain.
func (s *Service) ValidateRecipientDomain(email string) error {
	// very simple: split at '@' and check domain exists
	var domain string
	for i := len(email) - 1; i >= 0; i-- {
		if email[i] == '@' {
			domain = email[i+1:]
			break
		}
	}
	if domain == "" {
		return sql.ErrNoRows
	}

	var name string
	err := s.db.QueryRow(`SELECT name FROM domains WHERE name = $1`, domain).Scan(&name)
	return err
}
