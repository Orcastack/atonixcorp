package policy

import (
	"errors"

	"atonixcorp/security/auth"
)

// Engine evaluates permissions and enforces security rules.
type Engine struct{}

// NewEngine creates a new policy engine.
func NewEngine() *Engine {
	return &Engine{}
}

// CheckPermission verifies that an identity is allowed to perform an action.
func (e *Engine) CheckPermission(id auth.Identity, action Action) error {
	// Check roles
	for _, role := range id.Roles {
		if HasPermission(role, action) {
			return nil
		}
	}
	return errors.New("policy: permission denied")
}

// CheckProjectIsolation ensures identity cannot access another project.
func (e *Engine) CheckProjectIsolation(id auth.Identity, targetProject string) error {
	if id.ProjectID != targetProject {
		return errors.New("policy: cross-project access denied")
	}
	return nil
}

// Authorize combines permission + project isolation checks.
func (e *Engine) Authorize(id auth.Identity, action Action, targetProject string) error {
	if err := e.CheckProjectIsolation(id, targetProject); err != nil {
		return err
	}
	if err := e.CheckPermission(id, action); err != nil {
		return err
	}
	return nil
}
