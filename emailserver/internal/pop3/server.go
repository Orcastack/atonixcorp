package pop3

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Config struct {
	ListenAddr string
}

type Logger interface {
	Info(...any)
	Error(...any)
}

type Message struct {
	ID   string
	Size int
	Body []byte
}

type Store interface {
	ListPOP3Messages(userID string) ([]Message, error)
	GetPOP3Message(userID, id string) ([]byte, error)
	DeletePOP3Message(userID, id string) error
}

type AuthService interface {
	Authenticate(username, password string) (string, error)
}

type Server struct {
	cfg    Config
	logger Logger
	store  Store
	auth   AuthService
}

func NewServer(cfg Config, logger Logger, store Store, auth AuthService) *Server {
	return &Server{cfg, logger, store, auth}
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.cfg.ListenAddr)
	if err != nil {
		return err
	}

	s.logger.Info("POP3 listening on ", s.cfg.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			s.logger.Error("POP3 accept:", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)

	write(w, "+OK AtonixCorp POP3 Server Ready")

	session := &Session{
		Authenticated: false,
		UserID:        "",
		Messages:      nil,
		Deleted:       map[int]bool{},
	}

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		cmd, args := parse(line)

		switch strings.ToUpper(cmd) {

		case "USER":
			if len(args) < 1 {
				write(w, "-ERR USER requires username")
				continue
			}
			session.Username = args[0]
			write(w, "+OK User accepted")

		case "PASS":
			if session.Username == "" {
				write(w, "-ERR USER required before PASS")
				continue
			}
			if len(args) < 1 {
				write(w, "-ERR PASS requires password")
				continue
			}
			userID, err := s.auth.Authenticate(session.Username, args[0])
			if err != nil {
				write(w, "-ERR Authentication failed")
				continue
			}
			session.Authenticated = true
			session.UserID = userID

			msgs, err := s.store.ListPOP3Messages(userID)
			if err != nil {
				write(w, "-ERR Cannot load messages")
				continue
			}
			session.Messages = msgs

			write(w, "+OK Mailbox locked and ready")

		case "STAT":
			if !session.Authenticated {
				write(w, "-ERR Authenticate first")
				continue
			}
			count, size := session.stats()
			write(w, fmt.Sprintf("+OK %d %d", count, size))

		case "LIST":
			if !session.Authenticated {
				write(w, "-ERR Authenticate first")
				continue
			}
			if len(args) == 0 {
				// list all
				write(w, "+OK scan listing follows")
				for i, m := range session.Messages {
					if !session.Deleted[i+1] {
						write(w, fmt.Sprintf("%d %d", i+1, m.Size))
					}
				}
				write(w, ".")
			} else {
				// list single
				idx := atoi(args[0])
				if idx < 1 || idx > len(session.Messages) || session.Deleted[idx] {
					write(w, "-ERR No such message")
					continue
				}
				m := session.Messages[idx-1]
				write(w, fmt.Sprintf("+OK %d %d", idx, m.Size))
			}

		case "RETR":
			if !session.Authenticated {
				write(w, "-ERR Authenticate first")
				continue
			}
			idx := atoi(args[0])
			if idx < 1 || idx > len(session.Messages) || session.Deleted[idx] {
				write(w, "-ERR No such message")
				continue
			}
			body, err := s.store.GetPOP3Message(session.UserID, session.Messages[idx-1].ID)
			if err != nil {
				write(w, "-ERR Cannot retrieve message")
				continue
			}
			write(w, fmt.Sprintf("+OK %d octets", len(body)))
			writeRaw(w, body)
			write(w, ".")

		case "DELE":
			if !session.Authenticated {
				write(w, "-ERR Authenticate first")
				continue
			}
			idx := atoi(args[0])
			if idx < 1 || idx > len(session.Messages) || session.Deleted[idx] {
				write(w, "-ERR No such message")
				continue
			}
			session.Deleted[idx] = true
			write(w, "+OK Message marked for deletion")

		case "RSET":
			session.Deleted = map[int]bool{}
			write(w, "+OK Deletions undone")

		case "QUIT":
			if session.Authenticated {
				for idx := range session.Deleted {
					m := session.Messages[idx-1]
					s.store.DeletePOP3Message(session.UserID, m.ID)
				}
			}
			write(w, "+OK AtonixCorp POP3 Server signing off")
			return

		default:
			write(w, "-ERR Unknown command")
		}
	}
}

type Session struct {
	Authenticated bool
	Username      string
	UserID        string
	Messages      []Message
	Deleted       map[int]bool
}

func (s *Session) stats() (count int, size int) {
	for i, m := range s.Messages {
		if !s.Deleted[i+1] {
			count++
			size += m.Size
		}
	}
	return
}

func write(w *bufio.Writer, msg string) {
	w.WriteString(msg + "\r\n")
	w.Flush()
}

func writeRaw(w *bufio.Writer, body []byte) {
	w.Write(body)
	w.WriteString("\r\n")
	w.Flush()
}

func parse(line string) (cmd string, args []string) {
	parts := strings.Split(line, " ")
	cmd = parts[0]
	if len(parts) > 1 {
		args = parts[1:]
	}
	return
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
