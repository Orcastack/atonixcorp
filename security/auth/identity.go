package auth

// Identity represents a user, service, or machine identity inside AtonixCorp.
type Identity struct {
	ID        string   // unique identity ID (user, bot, service)
	ProjectID string   // project or tenant
	Roles     []string // RBAC roles
}

// HasRole checks if identity has a specific role.
func (i Identity) HasRole(role string) bool {
	for _, r := range i.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// IdentityResolver resolves identity from tokens, certificates, or metadata.
type IdentityResolver interface {
	Resolve(token string) (Identity, error)
}
