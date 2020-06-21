package dao

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Balances map[Account]uint

type State struct {
	Balances  Balances // The current balances
	txMempool []Tx     // The transactions that are executed but not in the tx.dao file

	blockDbFile     *os.File // The handler to the transaction file
	latestBlockHash Hash     // Hash code associated with the current block
}

func NewStateFromDisk(dataDir string) (*State, error) {
	err := initDataDirIfNotExists(dataDir)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialise the data: %w", err)
	}

	// Load the genesis file
	gen, err := loadGenesis(getGenesisJsonFilePath(dataDir))
	if err != nil {
		return nil, fmt.Errorf("Failed to load the genesis file: %w", err)
	}

	// Load the balances for each account
	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	// Now open the file for append as well as read
	blockDbFile, err := os.OpenFile(getBlocksDbFilePath(dataDir), os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("Cannot open block file: %w", err)
	}

	// Create the baseline state
	state := &State{balances, make([]Tx, 0), blockDbFile, Hash{}}

	// Read each line separately - each line is a transaction
	scanner := bufio.NewScanner(blockDbFile)

	// Iterate over each the existing transactions in the tx.dao file
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		// Read next json message within the block
		blockFsJson := scanner.Bytes()
		var blockFs BlockFS
		err = json.Unmarshal(blockFsJson, &blockFs)
		if err != nil {
			return nil, err
		}

		err = state.applyBlock(blockFs.Value)
		if err != nil {
			return nil, err
		}
		// Now we update the state with the hash of the last block
		state.latestBlockHash = blockFs.Key

	}

	if err != nil {
		return nil, err
	}

	return state, nil
}

func (s *State) LatestBlockHash() Hash {
	return s.latestBlockHash
}

// This iterates through each transaction within the block
// and adds each tx using AddTx
func (s *State) AddBlock(b Block) error {
	for _, tx := range b.TXs {
		if err := s.AddTx(tx); err != nil {
			return err
		}
	}
	return nil
}

// This applies a transaction and then remembers it in the memory pool
func (s *State) AddTx(tx Tx) error {
	if err := s.applyTx(tx); err != nil {
		return err
	}
	s.txMempool = append(s.txMempool, tx)

	return nil
}

// This saves the latest state of the block chain
// It then takes the whole file and generates a hash
func (s *State) Persist() (Hash, error) {
	// First of all create the block
	block := NewBlock(s.latestBlockHash, uint64(time.Now().Unix()), s.txMempool)
	// Now determine the hash for the whole block
	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, err
	}

	// Create the block with the hash
	blockFs := BlockFS{blockHash, block}

	// Encode it into a JSON string
	blockFsJson, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, err
	}

	fmt.Printf("Persisting new Block to disk:\n")
	fmt.Printf("\t%s\n", blockFsJson)

	// Write it to the DB file on a new line
	_, err = s.blockDbFile.Write(append(blockFsJson, '\n'))
	if err != nil {
		return Hash{}, err
	}
	s.latestBlockHash = blockHash

	// Clear the mempool
	s.txMempool = []Tx{}

	return s.latestBlockHash, nil
}

func (s *State) Close() {
	s.blockDbFile.Close()
}

// Apply all the transactions in the block
func (s *State) applyBlock(b Block) error {
	for _, tx := range b.TXs {
		if err := s.applyTx(tx); err != nil {
			return err
		}
	}

	return nil
}

// This applies a transaction to the balances
// It does not store the transaction in the mempool
func (s *State) applyTx(tx Tx) error {
	// If this is a reward then just increase the value of the whole pot
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if s.Balances[tx.From] < tx.Value {
		return fmt.Errorf("insufficient balance")
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}

//// take the content of all the transactions and use it to generate a sha256 hash
//func (s *State) doSnapshot() error {
//	// Re-read the whole file from the first byte
//	_, err := s.blockDbFile.Seek(0, 0)
//	if err != nil {
//		return err
//	}
//
//	txsData, err := ioutil.ReadAll(s.txDbFile)
//	if err != nil {
//		return err
//	}
//	s.snapshot = sha256.Sum256(txsData)
//
//	return nil
//}
