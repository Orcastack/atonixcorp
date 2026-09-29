package utils

import (
	"io"
	"os"
	"path/filepath"
)

func EnsureDir(path string) error {
	if path == "" {
		return New(ErrCodeValidation, "empty directory path")
	}
	return os.MkdirAll(path, 0o755)
}

func WriteFile(path string, r io.Reader) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return Wrap(ErrCodeInternal, "create file failed", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return Wrap(ErrCodeInternal, "write file failed", err)
	}

	return nil
}

func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, Wrap(ErrCodeInternal, "read file failed", err)
	}
	return data, nil
}
