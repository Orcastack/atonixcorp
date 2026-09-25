package ledger

import "errors"

// Chain represents an in-memory view of a ledger chain.
type Chain struct {
	blocks []Block
	hasher Hasher
	signer Signer
}

// NewChain creates a new chain with the given hasher and signer.
func NewChain(hasher Hasher, signer Signer) *Chain {
	return &Chain{
		blocks: make([]Block, 0),
		hasher: hasher,
		signer: signer,
	}
}

// Append creates a new block from payload and identity, links it, signs it, and adds it to the chain.
func (c *Chain) Append(payload []byte, identity string) (Block, error) {
	var index int64 = int64(len(c.blocks))
	var prevHash []byte
	if index > 0 {
		prevHash = c.blocks[index-1].Hash
	}

	// Compute hash over (prevHash || payload)
	data := append(prevHash, payload...)
	hash := c.hasher.Hash(data)

	// Sign hash
	sig, err := c.signer.Sign(identity, hash)
	if err != nil {
		return Block{}, err
	}

	b := Block{
		Index:     index,
		PrevHash:  prevHash,
		Hash:      hash,
		Payload:   payload,
		Signature: sig,
		Identity:  identity,
	}

	c.blocks = append(c.blocks, b)
	return b, nil
}

// VerifyIntegrity checks that all blocks are correctly hash-linked and signatures verify.
func (c *Chain) VerifyIntegrity() error {
	for i, b := range c.blocks {
		var expectedPrev []byte
		if i > 0 {
			expectedPrev = c.blocks[i-1].Hash
		}
		if !equalBytes(b.PrevHash, expectedPrev) {
			return errors.New("ledger: broken prevHash link at index " + string(rune(i)))
		}

		data := append(b.PrevHash, b.Payload...)
		hash := c.hasher.Hash(data)
		if !equalBytes(hash, b.Hash) {
			return errors.New("ledger: hash mismatch at index " + string(rune(i)))
		}

		if err := c.signer.Verify(b.Identity, b.Hash, b.Signature); err != nil {
			return err
		}
	}
	return nil
}

func (c *Chain) Blocks() []Block {
	return c.blocks
}

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
