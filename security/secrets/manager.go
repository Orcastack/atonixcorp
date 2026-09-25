package secrets

import "errors"

// SecretManager defines the interface for retrieving and managing secrets.
type SecretManager interface {
	Get(name string) ([]byte, error)
	Set(name string, value []byte) error
	Rotate(name string) ([]byte, error)
}

// InMemoryManager is a simple implementation for testing and development.
type InMemoryManager struct {
	store map[string][]byte
}

func NewInMemoryManager() *InMemoryManager {
	return &InMemoryManager{
		store: make(map[string][]byte),
	}
}

func (m *InMemoryManager) Get(name string) ([]byte, error) {
	val, ok := m.store[name]
	if !ok {
		return nil, errors.New("secret not found: " + name)
	}
	return val, nil
}

func (m *InMemoryManager) Set(name string, value []byte) error {
	m.store[name] = value
	return nil
}

func (m *InMemoryManager) Rotate(name string) ([]byte, error) {
	// For real rotation, integrate with Vault or KMS.
	newValue := []byte("rotated-" + name)
	m.store[name] = newValue
	return newValue, nil
}
