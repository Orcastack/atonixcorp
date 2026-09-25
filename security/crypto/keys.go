package crypto

import "crypto/rsa"

type KeyManager interface {
	GetPublicKey(identity string) (*rsa.PublicKey, error)
	GetPrivateKey(identity string) (*rsa.PrivateKey, error)
}
