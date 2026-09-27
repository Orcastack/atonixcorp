package internal

import "fmt"

var (
	ErrNotFound = fmt.Errorf("not found")
)

func Wrap(msg string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}
