package cli

import (
	"fmt"
	"os"
	"simpleblockchain/dao"
)

var state *dao.State

// Establish the current state by starting at genesis
// and applying any existing transactions
func openState() {
	var err error
	state, err = dao.LoadStateFromDisk(dataDir)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func closeState() {
	err := state.Close()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
