package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"simpleblockchain/node"
)

var (
	ip   string
	port uint64
)

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

			// This is the configuration of the bootstrap peer
			// Myst have one bootstrap in every peer-to-peer syste,
			// Everyone registers with bootstrap
			bootstrap := node.NewPeerNode("127.0.0.1", 8080, true, false)

			n := node.New(state, ip, port, bootstrap)
			err := n.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	runCmd.Flags().Uint64VarP(&port, "port", "p", node.DefaultHTTPort, "exposed HTTP port for communication with peers")
	runCmd.Flags().StringVar(&ip, "ip", node.DefaultIP, "exposed IP for communication with peers")

	return runCmd
}
