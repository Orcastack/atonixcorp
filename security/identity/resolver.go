package identity

import (
	"errors"

	"atonixcorp/security/auth"
)

// Resolver resolves identity from a token or raw input.
type Resolver struct {
	TokenService *auth.TokenService
}

// NewResolver creates a new identity resolver.
func NewResolver(ts *auth.TokenService) *Resolver {
	return &Resolver{TokenService: ts}
}

// ResolveIdentity resolves identity from a raw token string.
func (r *Resolver) ResolveIdentity(token string) (auth.Identity, error) {
	// Parse token into struct
	t, err := r.TokenService.ParseRawToken(token)
	if err != nil {
		return auth.Identity{}, err
	}

	// Verify signature + expiry
	if err := r.TokenService.VerifyToken(t); err != nil {
		return auth.Identity{}, err
	}

	return t.Identity, nil
}

// ParseRawToken is implemented inside TokenService.
// We expose it here for clarity.
func (ts *auth.TokenService) ParseRawToken(raw string) (auth.Token, error) {
	// raw format: id|project|expiry|signature(base64)
	parts := auth.SplitToken(raw)
	if parts == nil {
		return auth.Token{}, errors.New("invalid token format")
	}

	id := auth.Identity{
		ID:        parts.ID,
		ProjectID: parts.Project,
		Roles:     parts.Roles,
	}

	return auth.Token{
		Identity:  id,
		Expires:   parts.Expires,
		Signature: parts.Signature,
	}, nil
}
