package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"simpleblockchain/cli"
)

func main() {
	var tbbCmd = &cobra.Command{
		Use:   "tbb",
		Short: "The Blockchain Bar CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	cli.AddGlobalFlags(tbbCmd)

	tbbCmd.AddCommand(cli.VersionCmd)
	tbbCmd.AddCommand(cli.BalancesCmd())
	tbbCmd.AddCommand(cli.RunCmd())
	tbbCmd.AddCommand(cli.TxCmd())

	err := tbbCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
