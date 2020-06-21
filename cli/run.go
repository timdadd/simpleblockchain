package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"simpleblockchain/node"
)

var httpPort int

func RunCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Launches the TBB node and its HTTP API.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			openState()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			closeState()
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Launching TBB node and its HTTP API...")
			err := node.Run(httpPort, state)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	runCmd.Flags().IntVarP(&httpPort, "port", "p", 8080, "IP Port to listen on")

	return runCmd
}
