package domain

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
)

// GenerateDKIMKeyPair creates and stores DKIM keys for a domain.
func (s *Service) GenerateDKIMKeyPair(domain string, basePath string) (publicTXT string, err error) {
	// 2048-bit RSA key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}

	privBytes := x509.MarshalPKCS1PrivateKey(key)
	privPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})

	pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return "", err
	}
	pubPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	// store private key on disk (configs/dkim/<domain>.key)
	privPath := filepath.Join(basePath, domain+".key")
	if err := os.WriteFile(privPath, privPem, 0600); err != nil {
		return "", err
	}

	// DKIM TXT record content (selector "atonix")
	// e.g. atonix._domainkey.example.com TXT "v=DKIM1; k=rsa; p=BASE64..."
	pubBlock, _ := pem.Decode(pubPem)
	if pubBlock == nil {
		return "", errors.New("failed to decode public key")
	}
	pubDER := pubBlock.Bytes

	// base64 of DER
	// we keep it simple: user will paste this into DNS
	publicTXT = "v=DKIM1; k=rsa; p=" + encodeBase64(pubDER)

	_, err = s.db.Exec(`
        UPDATE domains
        SET dkim_selector = 'atonix', dkim_public = $2
        WHERE name = $1
    `, domain, publicTXT)

	return publicTXT, err
}

// encodeBase64 is a tiny helper; you can move it to util if you like.
func encodeBase64(b []byte) string {
	const table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var out []byte
	for i := 0; i < len(b); i += 3 {
		var v uint32
		n := 0
		for j := 0; j < 3 && i+j < len(b); j++ {
			v <<= 8
			v |= uint32(b[i+j])
			n++
		}
		for j := 0; j < 4; j++ {
			if j > n {
				out = append(out, '=')
			} else {
				idx := (v >> uint(18-6*j)) & 0x3F
				out = append(out, table[idx])
			}
		}
	}
	return string(out)
}
