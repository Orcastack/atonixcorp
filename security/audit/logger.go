package audit

import (
	"time"

	"atonixcorp/security/ledger"
)

// Auditor is the interface used across the entire platform.
type Auditor interface {
	Info(identity string, event string, fields map[string]any)
	Warn(identity string, event string, fields map[string]any)
	Error(identity string, event string, fields map[string]any)
}

// LedgerAuditor writes audit events into the ledger chain.
type LedgerAuditor struct {
	chain *ledger.Chain
}

func NewLedgerAuditor(chain *ledger.Chain) *LedgerAuditor {
	return &LedgerAuditor{chain: chain}
}

func (a *LedgerAuditor) Info(identity string, event string, fields map[string]any) {
	a.append("INFO", identity, event, fields)
}

func (a *LedgerAuditor) Warn(identity string, event string, fields map[string]any) {
	a.append("WARN", identity, event, fields)
}

func (a *LedgerAuditor) Error(identity string, event string, fields map[string]any) {
	a.append("ERROR", identity, event, fields)
}

func (a *LedgerAuditor) append(level, identity, event string, fields map[string]any) {
	payload := map[string]any{
		"level":     level,
		"event":     event,
		"fields":    fields,
		"identity":  identity,
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
	}

	// Serialize payload (simple JSON)
	data, _ := ledger.MarshalJSON(payload)

	// Append block to ledger
	a.chain.Append(data, identity)
}
