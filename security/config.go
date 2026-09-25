package security

// Config represents global security configuration for the AtonixCorp platform.
type Config struct {
	// Secret backend: "vault", "memory", "postgres", etc.
	SecretBackend string

	// Crypto settings
	AESKeyName     string // name of AES key in secret manager
	RSAKeyName     string // name of RSA private key
	Ed25519KeyName string // name of Ed25519 private key

	// Ledger settings
	LedgerStore string // "postgres", "file", "memory"
	LedgerPath  string // file path or table name

	// Policy settings
	EnforceProjectIsolation bool
	EnforceRBAC             bool

	// Detection settings
	EnableAnomalyDetection bool
	EnableSecurityRules    bool

	// Audit settings
	EnableAudit bool
}

// DefaultConfig returns a safe default configuration.
func DefaultConfig() Config {
	return Config{
		SecretBackend:           "memory",
		AESKeyName:              "aes-master",
		RSAKeyName:              "rsa-master",
		Ed25519KeyName:          "ed25519-master",
		LedgerStore:             "memory",
		LedgerPath:              "",
		EnforceProjectIsolation: true,
		EnforceRBAC:             true,
		EnableAnomalyDetection:  true,
		EnableSecurityRules:     true,
		EnableAudit:             true,
	}
}
