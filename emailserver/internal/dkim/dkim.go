package dkim

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
	"time"
)

type Signer struct {
	Domain   string
	Selector string
	Key      *rsa.PrivateKey
	Hash     crypto.Hash
}

func NewSigner(domain, selector string, pemKey []byte) (*Signer, error) {
	block, _ := pem.Decode(pemKey)
	if block == nil {
		return nil, fmt.Errorf("dkim: invalid PEM")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return &Signer{
		Domain:   domain,
		Selector: selector,
		Key:      key,
		Hash:     crypto.SHA256,
	}, nil
}

// raw is full RFC 5322 message (headers + body)
func (s *Signer) Sign(raw []byte) ([]byte, error) {
	headers, body := splitHeadersBody(raw)

	// canonicalize body (simple)
	bodyCanon := canonicalizeBody(body)

	// body hash
	h := s.Hash.New()
	h.Write(bodyCanon)
	bodyHash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// DKIM-Signature header (without b=)
	dkimHeader := s.buildHeader(bodyHash)

	// headers to sign: DKIM-Signature + From + To + Subject + Date
	var buf bytes.Buffer
	buf.WriteString(canonicalizeHeader("DKIM-Signature", dkimHeader))
	for _, name := range []string{"From", "To", "Subject", "Date"} {
		if v := findHeader(headers, name); v != "" {
			buf.WriteString(canonicalizeHeader(name, v))
		}
	}

	// sign
	sigHash := s.Hash.New()
	sigHash.Write(buf.Bytes())
	digest := sigHash.Sum(nil)

	sig, err := s.Key.Sign(rand.Reader, digest, s.Hash)
	if err != nil {
		return nil, err
	}
	sigB64 := base64.StdEncoding.EncodeToString(sig)

	// final DKIM header
	finalHeader := dkimHeader + "; b=" + sigB64

	// prepend DKIM-Signature to raw message
	var out bytes.Buffer
	out.WriteString("DKIM-Signature: " + finalHeader + "\r\n")
	out.Write(headers)
	out.WriteString("\r\n")
	out.Write(body)
	return out.Bytes(), nil
}

func (s *Signer) buildHeader(bodyHash string) string {
	now := time.Now().Unix()
	h := "from:to:subject:date"
	return fmt.Sprintf(
		"v=1; a=rsa-sha256; c=simple/simple; d=%s; s=%s; t=%d; h=%s; bh=%s",
		s.Domain, s.Selector, now, h, bodyHash,
	)
}

func splitHeadersBody(raw []byte) (headers []byte, body []byte) {
	parts := bytes.SplitN(raw, []byte("\r\n\r\n"), 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return raw, nil
}

func canonicalizeHeader(name, value string) string {
	name = strings.TrimSpace(name)
	value = strings.TrimSpace(value)
	return fmt.Sprintf("%s: %s\r\n", name, value)
}

func canonicalizeBody(body []byte) []byte {
	// simple: ensure CRLF, strip trailing empty lines
	s := string(body)
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")

	// remove trailing empty lines
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return []byte(strings.Join(lines, "\r\n") + "\r\n")
}

func findHeader(headers []byte, name string) string {
	lines := strings.Split(string(headers), "\r\n")
	prefix := strings.ToLower(name) + ":"
	for _, l := range lines {
		if strings.HasPrefix(strings.ToLower(l), prefix) {
			return strings.TrimSpace(l[len(name)+1:])
		}
	}
	return ""
}
