package dao

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type Hash [32]byte

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

func (h *Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return err
}

func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) IsEmpty() bool {
	emptyHash := Hash{}
	return bytes.Equal(emptyHash[:], h[:])
}

type Block struct {
	Header BlockHeader `json:"header"`  // metadata (parent block hash + time)
	TXs    []Tx        `json:"payload"` // new transactions only (payload)
}

// This is the header of the block
type BlockHeader struct {
	Parent      Hash   `json:"parent"` // The hash of the previous block
	BlockNumber uint64 `json:"number"` // A sequence number for the block, "block height"
	Nonce       uint32 `json:"nonce"`  // Adding a bit of randomness to the block hash
	Time        uint64 `json:"time"`   // The time this block was completed
}

// This is what's written to the filesystem
// The hash here is the header of the block
// This hash is the hash in the next block
type BlockFS struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
}

// A block is made up of a header and transactions
// A block header has the time, sequence number and the hash of the previous block
func NewBlock(parent Hash, blockNumber uint64, nonce uint32, time uint64, txs []Tx) Block {
	return Block{BlockHeader{parent, blockNumber, nonce, time}, txs}
}

// This generates a hash for a block
func (b Block) Hash() (Hash, error) {
	blockJson, err := json.Marshal(b)
	if err != nil {
		return Hash{}, err
	}
	return sha256.Sum256(blockJson), nil
}

func IsBlockHashValid(hash Hash) bool {
	// 184648311	         6.36 ns/op
	return bytes.Equal(hash[:2], []byte{0, 0})
	// 8908318	       118 ns/op
	//return hash.Hex()[0:4] == "0000"
	// 4187757	       252 ns/op
	//return fmt.Sprintf("%x", hash[0]) == "0" &&
	//	fmt.Sprintf("%x", hash[1]) == "0" &&
	//	fmt.Sprintf("%x", hash[2]) == "0" &&
	//	fmt.Sprintf("%x", hash[3]) != "0"
}

// This returns all the blocks after a specific hash
func GetBlocksAfter(blockHash Hash, s *State) ([]Block, error) {
	f, err := os.OpenFile(getBlocksDbFilePath(s.dataDir), os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("Could not open the local blocks file: %w", err)
	}

	blocks := make([]Block, 0)
	shouldStartCollecting := false

	if reflect.DeepEqual(blockHash, Hash{}) {
		shouldStartCollecting = true
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("Cannot read line from block: %w", err)
		}

		// Read next json message within the block
		blockFsJson := scanner.Bytes()
		var blockFs BlockFS
		err = json.Unmarshal(blockFsJson, &blockFs)
		if err != nil {
			return nil, fmt.Errorf("Cannot interpret json %v: %w", blockFsJson, err)
		}

		// Are we starting to collect blocks?
		if shouldStartCollecting {
			blocks = append(blocks, blockFs.Value)
		} else if blockHash == blockFs.Key { // Should we start collecting blocks?
			shouldStartCollecting = true
		}
	}

	return blocks, nil
}
