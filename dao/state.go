package dao

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type State struct {
	Balances  map[Account]uint // The current balances
	txMempool []Tx             // The transactions that are executed but not in the tx.dao file

	txDbFile *os.File // The handler to the transaction file
}

const (
	databaseDir = "db"
)

func NewStateFromDisk() (*State, error) {
	// get current working directory or error
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	//Determine the genesis.json filepath
	genFilePath := filepath.Join(cwd, databaseDir, "genesis.json")
	gen, err := loadGenesis(genFilePath)
	if err != nil {
		return nil, err
	}

	// Load the balances for each account
	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	// Now apply the transactions in the transaction file which is a set of json messages
	txDbFilePath := filepath.Join(cwd, databaseDir, "tx.dao")
	// If the file doesn't exist then create it
	_, err = os.Stat(txDbFilePath)
	if os.IsNotExist(err) {
		_, err = os.Create(txDbFilePath)
	}
	if err != nil {
		return nil, err
	}

	// Now open the file for append as well as read
	txDbFile, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	// Create the baseline state
	state := &State{balances, make([]Tx, 0), txDbFile}

	// Read each line separately - each line is a transaction
	scanner := bufio.NewScanner(txDbFile)

	// Iterate over each the existing transactions in the tx.dao file
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		// Convert JSON encoded TX into an object (struct)
		var tx Tx
		json.Unmarshal(scanner.Bytes(), &tx)

		// Rebuild the state (user balances), by re-applying all the transactions
		// since genesis.  tx are events, dao is the state
		if err := state.apply(tx); err != nil {
			return nil, err
		}
	}

	return state, nil
}

// This applies a transaction and then remembers it in the memory pool
func (s *State) Add(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}

	s.txMempool = append(s.txMempool, tx)

	return nil
}

// This saves the transactions that achieved the state to the tx.dao file
func (s *State) Persist() error {
	// create a copy of the mempool slice because the s.mempool
	// slice will be emptied as each transaction is
	// written to the file
	mempool := make([]Tx, len(s.txMempool))
	copy(mempool, s.txMempool)

	for i := 0; i < len(mempool); i++ {
		// JSON encode the transaction
		txJson, err := json.Marshal(s.txMempool[i])
		if err != nil {
			return err
		}
		// Append the json to the already opened "tx.dao" file
		if _, err = s.txDbFile.Write(append(txJson, '\n')); err != nil {
			return err
		}

		// Remove the TX written to a file from the mempool
		// Yes... this particular Go syntax is a bit weird
		s.txMempool = append(s.txMempool[:i], s.txMempool[i+1:]...)
	}

	return nil
}

func (s *State) Close() {
	s.txDbFile.Close()
}

// This applies a transaction to the balances
// It does not store the transaction in the mempool
func (s *State) apply(tx Tx) error {
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
