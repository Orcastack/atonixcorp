package smtp

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"strings"
)

type Config struct {
	ListenAddr string
	TLSCert    string
	TLSKey     string
}

type Logger interface {
	Info(...any)
	Error(...any)
}

type Store interface {
	SaveIncomingMessage(msg *Message) error
}

type Queue interface {
	EnqueueOutgoing(msg *Message) error
}

type SpamFilter interface {
	Check(msg *Message) (bool, error)
}

type DomainService interface {
	ValidateRecipientDomain(email string) error
}

type Handler interface {
	BeforeAccept(msg *Message) error
	AfterAccept(msg *Message)
}

type Message struct {
	From    string
	To      []string
	Subject string
	Body    []byte
	Raw     []byte
}

type Server struct {
	cfg     Config
	logger  Logger
	store   Store
	queue   Queue
	spam    SpamFilter
	domain  DomainService
	handler Handler
}

func NewServer(cfg Config, logger Logger, store Store, queue Queue, spam SpamFilter, domain DomainService, handler Handler) *Server {
	return &Server{
		cfg:     cfg,
		logger:  logger,
		store:   store,
		queue:   queue,
		spam:    spam,
		domain:  domain,
		handler: handler,
	}
}

func (s *Server) ListenAndServe() error {
	cert, err := tls.LoadX509KeyPair(s.cfg.TLSCert, s.cfg.TLSKey)
	if err != nil {
		return fmt.Errorf("load tls cert: %w", err)
	}

	tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	ln, err := tls.Listen("tcp", s.cfg.ListenAddr, tlsCfg)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	s.logger.Info("SMTP listening on ", s.cfg.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			s.logger.Error("accept:", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)

	writeLine(w, "220 atonixcorp.emailserver ESMTP ready")

	var session Session
	session.Reset()

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			continue
		}

		cmd, arg := parseCmd(line)

		switch strings.ToUpper(cmd) {
		case "EHLO", "HELO":
			writeLine(w, "250-atonixtower")
			writeLine(w, "250-SIZE 52428800")
			writeLine(w, "250-8BITMIME")
			writeLine(w, "250-PIPELINING")
			writeLine(w, "250 OK")

		case "MAIL":
			if !strings.HasPrefix(strings.ToUpper(arg), "FROM:") {
				writeLine(w, "501 Syntax: MAIL FROM:<address>")
				continue
			}
			from := extractAddress(arg[5:])
			session.From = from
			writeLine(w, "250 OK")

		case "RCPT":
			if !strings.HasPrefix(strings.ToUpper(arg), "TO:") {
				writeLine(w, "501 Syntax: RCPT TO:<address>")
				continue
			}
			to := extractAddress(arg[3:])
			if err := s.domain.ValidateRecipientDomain(to); err != nil {
				writeLine(w, "550 Invalid recipient domain")
				continue
			}
			session.To = append(session.To, to)
			writeLine(w, "250 OK")

		case "DATA":
			if session.From == "" || len(session.To) == 0 {
				writeLine(w, "503 Need MAIL FROM and RCPT TO first")
				continue
			}
			writeLine(w, "354 End data with <CR><LF>.<CR><LF>")
			raw, body, subject := readData(r)
			msg := &Message{
				From:    session.From,
				To:      session.To,
				Subject: subject,
				Body:    body,
				Raw:     raw,
			}

			if s.handler != nil {
				if err := s.handler.BeforeAccept(msg); err != nil {
					writeLine(w, "550 Message rejected")
					session.Reset()
					continue
				}
			}

			isSpam, err := s.spam.Check(msg)
			if err != nil {
				s.logger.Error("spam check:", err)
			}
			if isSpam {
				writeLine(w, "550 Message rejected as spam")
				session.Reset()
				continue
			}

			if err := s.store.SaveIncomingMessage(msg); err != nil {
				s.logger.Error("store message:", err)
				writeLine(w, "451 Temporary local problem")
				session.Reset()
				continue
			}

			if err := s.queue.EnqueueOutgoing(msg); err != nil {
				s.logger.Error("enqueue outgoing:", err)
			}

			if s.handler != nil {
				s.handler.AfterAccept(msg)
			}

			writeLine(w, "250 Message accepted for delivery")
			session.Reset()

		case "QUIT":
			writeLine(w, "221 Bye")
			return

		default:
			writeLine(w, "502 Command not implemented")
		}
	}
}

type Session struct {
	From string
	To   []string
}

func (s *Session) Reset() {
	s.From = ""
	s.To = nil
}

func writeLine(w *bufio.Writer, line string) {
	w.WriteString(line + "\r\n")
	w.Flush()
}

func parseCmd(line string) (string, string) {
	parts := strings.SplitN(line, " ", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

func extractAddress(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "<")
	s = strings.TrimSuffix(s, ">")
	return s
}

func readData(r *bufio.Reader) (raw []byte, body []byte, subject string) {
	var lines []string
	for {
		l, err := r.ReadString('\n')
		if err != nil {
			break
		}
		if l == ".\r\n" || l == ".\n" {
			break
		}
		lines = append(lines, strings.TrimRight(l, "\r\n"))
	}
	rawStr := strings.Join(lines, "\r\n")
	raw = []byte(rawStr)
	body = []byte(rawStr)
	// TODO: parse headers to extract real Subject
	subject = ""
	return
}
