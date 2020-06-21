package dao

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type Balances map[Account]uint

type State struct {
	Balances  Balances // The current balances
	txMempool []Tx     // The transactions that are executed but not in the tx.dao file

	dataDir     string
	blockDbFile *os.File // The handler to the transaction file

	latestBlock     Block // The latest block
	latestBlockHash Hash  // Hash code associated with the current block
	hasGenesisBlock bool
}

func (s *State) LatestBlock() Block {
	return s.latestBlock
}

func (s *State) LatestBlockHash() Hash {
	return s.latestBlockHash
}

func (s *State) DataDir() string {
	return s.dataDir
}

func LoadStateFromDisk(dataDir string) (*State, error) {
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
	state := &State{Balances: balances,
		txMempool:       make([]Tx, 0),
		dataDir:         dataDir,
		blockDbFile:     blockDbFile,
		latestBlock:     Block{},
		latestBlockHash: Hash{},
		hasGenesisBlock: false}

	// Read each line separately - each line is a transaction
	scanner := bufio.NewScanner(blockDbFile)

	// Iterate over each the existing transactions in the tx.dao file
	// Surely everything is in the last line being the last block?
	//var blockFs BlockFS
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

		err = state.applyBlock(blockFs.Value)
		if err != nil {
			return nil, fmt.Errorf("Cannot apply block %v: %w", blockFs.Value, err)
		}
		// Now we update the state with the hash of the last block
		state.latestBlock = blockFs.Value
		state.latestBlockHash = blockFs.Key
		state.hasGenesisBlock = true
	}

	return state, nil
}

func (s *State) AddBlocks(blocks []Block) error {
	for _, b := range blocks {
		_, err := s.AddBlock(b)
		if err != nil {
			return fmt.Errorf("Could not add blocks %v: %w", b, err)
		}
	}

	return nil
}

// This iterates through each transaction within the block
// and adds each tx using AddTx
//// This saves the latest state of the block chain
//// It then takes the whole file and generates a hash
func (s *State) AddBlock(b Block) (Hash, error) {
	pendingState := s.copy()

	err := pendingState.applyBlock(b)
	if err != nil {
		return Hash{}, fmt.Errorf("Cannot apply block %v: %w", b, err)
	}

	blockHash, err := b.Hash()
	if err != nil {
		return Hash{}, fmt.Errorf("Cannot hash block %v: %w", b, err)
	}

	blockFs := BlockFS{blockHash, b}

	blockFsJson, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, fmt.Errorf("Cannot convert to json %v: %w", blockFs.Value, err)
	}

	fmt.Printf("Persisting new Block to disk:\n")
	fmt.Printf("\t%s\n", blockFsJson)

	_, err = s.blockDbFile.Write(append(blockFsJson, '\n'))
	if err != nil {
		return Hash{}, fmt.Errorf("Cannot append json to file %v: %w", blockFsJson, err)
	}

	s.Balances = pendingState.Balances
	s.latestBlockHash = blockHash
	s.latestBlock = b
	s.hasGenesisBlock = true

	return blockHash, nil
}

func (s *State) NextBlockNumber() uint64 {
	if !s.hasGenesisBlock {
		return uint64(0)
	}

	return s.LatestBlock().Header.BlockNumber + 1
}

// This applies a transaction and then remembers it in the memory pool
func (s *State) AddTx(tx Tx) error {
	if err := s.applyTx(tx); err != nil {
		return err
	}
	s.txMempool = append(s.txMempool, tx)

	return nil
}

func (s *State) Close() error {
	err := s.blockDbFile.Close()
	if err != nil {
		return fmt.Errorf("Could not close the block file: %w", err)
	}
	return nil
}

// applyBlock verifies if block can be added to the blockchain.
//
// Block metadata are verified as well as transactions within (sufficient balances, etc).
func (s *State) applyBlock(b Block) error {
	nextExpectedBlockNumber := s.latestBlock.Header.BlockNumber + 1

	if s.hasGenesisBlock && b.Header.BlockNumber != nextExpectedBlockNumber {
		return fmt.Errorf("next expected block must be '%d' not '%d'", nextExpectedBlockNumber, b.Header.BlockNumber)
	}

	if s.hasGenesisBlock && s.latestBlock.Header.BlockNumber > 0 && !reflect.DeepEqual(b.Header.Parent, s.latestBlockHash) {
		return fmt.Errorf("next block parent hash must be '%x' not '%x'", s.latestBlockHash, b.Header.Parent)
	}

	return s.applyTXs(b.TXs)
}

func (s *State) applyTXs(txs []Tx) error {
	for _, tx := range txs {
		err := s.applyTx(tx)
		if err != nil {
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

// Make a copy of the current state
func (s *State) copy() *State {
	c := State{}
	c.hasGenesisBlock = s.hasGenesisBlock
	c.latestBlock = s.latestBlock
	c.latestBlockHash = s.latestBlockHash
	c.txMempool = make([]Tx, len(s.txMempool))
	c.Balances = make(map[Account]uint)

	for acc, balance := range s.Balances {
		c.Balances[acc] = balance
	}

	for _, tx := range s.txMempool {
		c.txMempool = append(c.txMempool, tx)
	}

	return &c
}

//// This saves the latest state of the block chain
//// It then takes the whole file and generates a hash
//func (s *State) Persist() (Hash, error) {
//	// First of all create the block
//	block := NewBlock(
//		s.latestBlockHash,
//		s.latestBlockHash.Header.Number + 1, // increase height
//		uint64(time.Now().Unix()),
//		s.txMempool,
//	)
//	block := NewBlock(s.latestBlockHash, uint64(time.Now().Unix()), s.txMempool)
//	// Now determine the hash for the whole block
//	blockHash, err := block.Hash()
//	if err != nil {
//		return Hash{}, err
//	}
//
//	// Create the block with the hash
//	blockFs := BlockFS{blockHash, block}
//
//	// Encode it into a JSON string
//	blockFsJson, err := json.Marshal(blockFs)
//	if err != nil {
//		return Hash{}, err
//	}
//
//	//fmt.Printf("Persisting new Block to disk:\n")
//	//fmt.Printf("\t%s\n", blockFsJson)
//
//	// Write it to the DB file on a new line
//	_, err = s.blockDbFile.Write(append(blockFsJson, '\n'))
//	if err != nil {
//		return Hash{}, err
//	}
//	s.latestBlockHash = blockHash
//
//	// Clear the mempool
//	s.txMempool = []Tx{}
//
//	return s.latestBlockHash, nil
//}

//// Apply all the transactions in the block
//func (s *State) applyBlock(b Block) error {
//	for _, tx := range b.TXs {
//		if err := s.applyTx(tx); err != nil {
//			return fmt.Errorf("Cannot apply block: %w", err)
//		}
//	}
//	return nil
//}
