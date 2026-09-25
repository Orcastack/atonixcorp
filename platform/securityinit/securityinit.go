package securityinit

import (
	"atonixcorp/security"
	"atonixcorp/security/audit"
	"atonixcorp/security/auth"
	"atonixcorp/security/crypto"
	"atonixcorp/security/detect"
	"atonixcorp/security/ledger"
	"atonixcorp/security/policy"
	"atonixcorp/security/secrets"
)

type Stack struct {
	Config     security.Config
	Secrets    secrets.SecretManager
	KeyManager crypto.KeyManager
	Ledger     *ledger.Chain
	Auditor    audit.Auditor
	Detector   detect.Detector
	Policy     *policy.Engine
	Tokens     *auth.TokenService
}

func NewStack() *Stack {
	cfg := security.DefaultConfig()

	// secrets
	sm := secrets.NewInMemoryManager()

	// crypto key manager (you’ll implement this)
	km := crypto.NewSecretKeyManager(sm, cfg)

	// ledger
	hasher := ledger.SHA256Hasher{}
	signer := ledger.NewCryptoSigner(km)
	chain := ledger.NewChain(hasher, signer)

	// audit
	auditor := audit.NewLedgerAuditor(chain)

	// detect
	rules := []detect.Rule{
		detect.SuspiciousEventRule{EventName: "permission_denied"},
	}
	detector := detect.NewSimpleDetector(rules)

	// policy
	pol := policy.NewEngine()

	// tokens
	ts := auth.NewTokenService(km)

	return &Stack{
		Config:     cfg,
		Secrets:    sm,
		KeyManager: km,
		Ledger:     chain,
		Auditor:    auditor,
		Detector:   detector,
		Policy:     pol,
		Tokens:     ts,
	}
}
