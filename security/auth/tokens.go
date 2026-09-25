package auth

import (
	"errors"
	"time"

	"atonixcorp/security/crypto"
)

// Token contains identity + expiry + signature.
type Token struct {
	Identity  Identity
	Expires   time.Time
	Signature []byte
}

// TokenService issues and verifies tokens.
type TokenService struct {
	KM crypto.KeyManager
}

// NewTokenService creates a token service using your crypto keys.
func NewTokenService(km crypto.KeyManager) *TokenService {
	return &TokenService{KM: km}
}

// IssueToken creates a signed token for an identity.
func (ts *TokenService) IssueToken(id Identity, ttl time.Duration) (Token, error) {
	t := Token{
		Identity: id,
		Expires:  time.Now().UTC().Add(ttl),
	}

	// Serialize minimal token payload
	payload := []byte(id.ID + "|" + id.ProjectID + "|" + t.Expires.Format(time.RFC3339Nano))

	// Sign payload
	sig, err := crypto.Sign(id.ID, payload, ts.KM)
	if err != nil {
		return Token{}, err
	}

	t.Signature = sig
	return t, nil
}

// VerifyToken checks signature + expiry.
func (ts *TokenService) VerifyToken(t Token) error {
	if time.Now().UTC().After(t.Expires) {
		return errors.New("token expired")
	}

	payload := []byte(t.Identity.ID + "|" + t.Identity.ProjectID + "|" + t.Expires.Format(time.RFC3339Nano))

	return crypto.Verify(t.Identity.ID, payload, t.Signature, ts.KM)
}
