package dao

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Hash [32]byte

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

func (h *Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return err
}

type Block struct {
	Header BlockHeader `json:"header"`  // metadata (parent block hash + time)
	TXs    []Tx        `json:"payload"` // new transactions only (payload)
}

// This is the header of the block
type BlockHeader struct {
	Parent Hash   `json:"parent"`
	Time   uint64 `json:"time"`
}

// This is what's written to the filesystem
// The hash here is the header of the block
// This hash is the hash in the next block
type BlockFS struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
}

// A block is made up of a header and transactions
// A block header has the time and the parent
func NewBlock(parent Hash, time uint64, txs []Tx) Block {
	return Block{BlockHeader{parent, time}, txs}
}

// This genereates a hash for a block
func (b Block) Hash() (Hash, error) {
	blockJson, err := json.Marshal(b)
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(blockJson), nil
}