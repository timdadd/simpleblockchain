package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"simpleblockchain/dao"
)

var state *dao.State

// Establish the current state by starting at genesis
// and applying any existing transactions
func setState(cmd *cobra.Command, args []string, deferEarly bool) {
	var err error
	state, err = dao.NewStateFromDisk(dataDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if deferEarly {
		defer state.Close()
	}
}
