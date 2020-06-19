package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"simpleblockchain/cmd"
)

func main() {
	var tbbCmd = &cobra.Command{
		Use:   "tbb",
		Short: "The Blockchain Bar CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	tbbCmd.AddCommand(cmd.VersionCmd)
	tbbCmd.AddCommand(cmd.BalancesCmd())
	tbbCmd.AddCommand(cmd.TxCmd())

	err := tbbCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
