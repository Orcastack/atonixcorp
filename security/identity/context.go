package identity

import (
	"context"

	"atonixcorp/security/auth"
)

type contextKey string

const identityKey contextKey = "atonixcorp.identity"

// WithIdentity attaches an Identity to a context.
func WithIdentity(ctx context.Context, id auth.Identity) context.Context {
	return context.WithValue(ctx, identityKey, id)
}

// FromContext retrieves the Identity from a context.
func FromContext(ctx context.Context) (auth.Identity, bool) {
	val := ctx.Value(identityKey)
	if val == nil {
		return auth.Identity{}, false
	}
	id, ok := val.(auth.Identity)
	return id, ok
}
