package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path"
)

const flagDataDir = "datadir"

var dataDir string

func AddGlobalFlags(cmd *cobra.Command) {
	// get current working directory or error
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Determine the default data filepath
	defaultDataDir := path.Join(cwd, "data")
	// Persistent flag is global
	cmd.PersistentFlags().StringVarP(&dataDir, "datadir", "f", defaultDataDir, "Absolute path to the node data dir where the DB will be/is stored")
}
