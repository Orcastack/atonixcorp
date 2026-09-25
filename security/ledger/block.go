package ledger

type Block struct {
	Index     int64
	PrevHash  []byte
	Hash      []byte
	Payload   []byte
	Signature []byte
	Identity  string // who wrote this block
}
