package secrets

import "errors"

// VaultClient represents a connection to an external secret store (Vault, KMS, etc).
type VaultClient interface {
	Read(path string) ([]byte, error)
	Write(path string, value []byte) error
	Rotate(path string) ([]byte, error)
}

// VaultManager implements SecretManager using a VaultClient.
type VaultManager struct {
	client VaultClient
	prefix string
}

func NewVaultManager(client VaultClient, prefix string) *VaultManager {
	return &VaultManager{
		client: client,
		prefix: prefix,
	}
}

func (v *VaultManager) fullPath(name string) string {
	return v.prefix + "/" + name
}

func (v *VaultManager) Get(name string) ([]byte, error) {
	return v.client.Read(v.fullPath(name))
}

func (v *VaultManager) Set(name string, value []byte) error {
	return v.client.Write(v.fullPath(name), value)
}

func (v *VaultManager) Rotate(name string) ([]byte, error) {
	newValue, err := v.client.Rotate(v.fullPath(name))
	if err != nil {
		return nil, errors.New("vault rotation failed: " + err.Error())
	}
	return newValue, nil
}
