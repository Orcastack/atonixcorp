package security

import "errors"

// Common security errors used across the platform.
var (
	ErrUnauthorized       = errors.New("security: unauthorized")
	ErrForbidden          = errors.New("security: forbidden")
	ErrInvalidToken       = errors.New("security: invalid token")
	ErrExpiredToken       = errors.New("security: token expired")
	ErrInvalidSignature   = errors.New("security: invalid signature")
	ErrKeyNotFound        = errors.New("security: key not found")
	ErrSecretNotFound     = errors.New("security: secret not found")
	ErrLedgerTampered     = errors.New("security: ledger integrity failure")
	ErrCrossProjectAccess = errors.New("security: cross-project access denied")
	ErrPermissionDenied   = errors.New("security: permission denied")
	ErrAnomalyDetected    = errors.New("security: anomaly detected")
)

// WrapError adds context to a security error.
func WrapError(base error, msg string) error {
	return errors.New(msg + ": " + base.Error())
}
