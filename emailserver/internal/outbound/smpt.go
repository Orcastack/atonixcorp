package outbound

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"

	"atonixcorp/emailserver/internal/queue"
)

type Logger interface {
	Info(...any)
	Error(...any)
}

type SMTPSender struct {
	Logger Logger
}

func NewSMTPSender(logger Logger) *SMTPSender {
	return &SMTPSender{Logger: logger}
}

func (s *SMTPSender) Send(msg queue.OutgoingMessage) error {
	if len(msg.To) == 0 {
		return fmt.Errorf("no recipients")
	}

	for _, rcpt := range msg.To {
		if err := s.sendToRecipient(msg, rcpt); err != nil {
			return err
		}
	}
	return nil
}

func (s *SMTPSender) sendToRecipient(msg queue.OutgoingMessage, rcpt string) error {
	domain := extractDomain(rcpt)
	if domain == "" {
		return fmt.Errorf("invalid recipient: %s", rcpt)
	}

	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return fmt.Errorf("no MX for domain %s", domain)
	}

	// sort by preference (lowest first)
	// simple: assume returned already sorted

	for _, mx := range mxRecords {
		host := strings.TrimSuffix(mx.Host, ".")
		addr := net.JoinHostPort(host, "25")

		s.Logger.Info("Outbound: trying", addr)

		// DANE (very simplified): try TLS, but real DANE needs TLSA lookup + pinning
		tlsCfg := &tls.Config{
			ServerName: host,
			MinVersion: tls.VersionTLS12,
		}

		conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second}, "tcp", addr, tlsCfg)
		if err != nil {
			s.Logger.Error("Outbound: TLS connect failed:", err)
			continue
		}
		defer conn.Close()

		if err := s.smtpSession(conn, msg, rcpt); err != nil {
			s.Logger.Error("Outbound: SMTP session failed:", err)
			continue
		}

		return nil
	}

	return fmt.Errorf("all MX attempts failed for %s", domain)
}

func (s *SMTPSender) smtpSession(conn net.Conn, msg queue.OutgoingMessage, rcpt string) error {
	readBuf := make([]byte, 4096)

	// read banner
	if _, err := conn.Read(readBuf); err != nil {
		return err
	}

	writeLine(conn, "EHLO atonixcorp.com")
	if _, err := conn.Read(readBuf); err != nil {
		return err
	}

	writeLine(conn, fmt.Sprintf("MAIL FROM:<%s>", msg.From))
	if _, err := conn.Read(readBuf); err != nil {
		return err
	}

	writeLine(conn, fmt.Sprintf("RCPT TO:<%s>", rcpt))
	if _, err := conn.Read(readBuf); err != nil {
		return err
	}

	writeLine(conn, "DATA")
	if _, err := conn.Read(readBuf); err != nil {
		return err
	}

	// send raw message (already DKIM-signed)
	conn.Write(msg.Raw)
	writeLine(conn, "\r\n.")
	if _, err := conn.Read(readBuf); err != nil {
		return err
	}

	writeLine(conn, "QUIT")
	return nil
}

func writeLine(conn net.Conn, line string) {
	conn.Write([]byte(line + "\r\n"))
}

func extractDomain(addr string) string {
	parts := strings.Split(addr, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}
