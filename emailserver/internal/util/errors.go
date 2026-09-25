package util

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	// Generic
	ErrUnknown  ErrorCode = "unknown_error"
	ErrInternal ErrorCode = "internal_error"

	// Auth
	ErrAuthFailed   ErrorCode = "auth_failed"
	ErrUnauthorized ErrorCode = "unauthorized"

	// SMTP
	ErrInvalidSender    ErrorCode = "invalid_sender"
	ErrInvalidRecipient ErrorCode = "invalid_recipient"
	ErrSpamDetected     ErrorCode = "spam_detected"

	// Storage
	ErrNotFound     ErrorCode = "not_found"
	ErrSaveFailed   ErrorCode = "save_failed"
	ErrDeleteFailed ErrorCode = "delete_failed"

	// Queue
	ErrDeliveryFailed ErrorCode = "delivery_failed"
	ErrRetryExceeded  ErrorCode = "retry_exceeded"

	// Rate limit
	ErrRateLimited ErrorCode = "rate_limited"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// ─────────────────────────────────────────────
// Constructors
// ─────────────────────────────────────────────

func New(code ErrorCode, msg string) *AppError {
	return &AppError{Code: code, Message: msg}
}

func Wrap(code ErrorCode, msg string, err error) *AppError {
	return &AppError{Code: code, Message: msg, Err: err}
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

func Is(err error, code ErrorCode) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}
