package imap

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

type Store interface {
	ListMailboxes(userID string) ([]Mailbox, error)
	ListMessages(userID, mailbox string) ([]Message, error)
	GetMessageBody(userID, mailbox, id string) ([]byte, error)
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
	return &Server{
		cfg:    cfg,
		logger: logger,
		store:  store,
		auth:   auth,
	}
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.cfg.ListenAddr)
	if err != nil {
		return err
	}

	s.logger.Info("IMAP listening on ", s.cfg.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			s.logger.Error("IMAP accept:", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)

	write(w, "* OK AtonixCorp IMAP Server Ready")

	session := &Session{
		Authenticated: false,
		SelectedBox:   "",
		UserID:        "",
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

		tag, cmd, args := parse(line)

		switch strings.ToUpper(cmd) {

		case "LOGIN":
			if len(args) < 2 {
				write(w, tag+" BAD LOGIN requires username and password")
				continue
			}
			user, pass := args[0], args[1]
			userID, err := s.auth.Authenticate(user, pass)
			if err != nil {
				write(w, tag+" NO Authentication failed")
				continue
			}
			session.Authenticated = true
			session.UserID = userID
			write(w, tag+" OK LOGIN completed")

		case "LIST":
			if !session.Authenticated {
				write(w, tag+" NO Authenticate first")
				continue
			}
			boxes, err := s.store.ListMailboxes(session.UserID)
			if err != nil {
				write(w, tag+" NO Cannot list mailboxes")
				continue
			}
			for _, b := range boxes {
				write(w, fmt.Sprintf("* LIST () \"/\" \"%s\"", b.Name))
			}
			write(w, tag+" OK LIST completed")

		case "SELECT":
			if !session.Authenticated {
				write(w, tag+" NO Authenticate first")
				continue
			}
			if len(args) < 1 {
				write(w, tag+" BAD SELECT requires mailbox")
				continue
			}
			session.SelectedBox = args[0]
			write(w, "* OK [READ-WRITE] Mailbox selected")
			write(w, tag+" OK SELECT completed")

		case "FETCH":
			if !session.Authenticated {
				write(w, tag+" NO Authenticate first")
				continue
			}
			if session.SelectedBox == "" {
				write(w, tag+" NO No mailbox selected")
				continue
			}
			msgs, err := s.store.ListMessages(session.UserID, session.SelectedBox)
			if err != nil {
				write(w, tag+" NO Cannot fetch messages")
				continue
			}
			for i, m := range msgs {
				write(w, fmt.Sprintf("* %d FETCH (UID %s FLAGS () BODY[] {%d})", i+1, m.ID, len(m.Body)))
				writeRaw(w, m.Body)
			}
			write(w, tag+" OK FETCH completed")

		case "LOGOUT":
			write(w, "* BYE Logging out")
			write(w, tag+" OK LOGOUT completed")
			return

		default:
			write(w, tag+" BAD Unknown command")
		}
	}
}

type Session struct {
	Authenticated bool
	SelectedBox   string
	UserID        string
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

func parse(line string) (tag, cmd string, args []string) {
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return "*", parts[0], parts[1:]
	}
	tag = parts[0]
	cmd = parts[1]
	args = parts[2:]
	return
}
