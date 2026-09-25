package ledger

// Store abstracts persistence for the ledger blocks.
type Store interface {
	AppendBlock(b Block) error
	GetBlock(index int64) (Block, error)
	GetLastBlock() (Block, error)
	ListBlocks() ([]Block, error)
}
