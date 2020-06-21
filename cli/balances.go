package cli

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"simpleblockchain/dao"
	"strings"
)

func BalancesCmd() *cobra.Command {
	var balancesCmd = &cobra.Command{
		Use:   "balances",
		Short: "Interact with balances (list...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ErrIncorrectUsage
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			openState()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			closeState()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	balancesCmd.AddCommand(balancesListCmd)
	balancesCmd.AddCommand(balancesStateCmd)

	return balancesCmd
}

var balancesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all balances.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Accounts balances at %x:\n", state.LatestBlockHash())
		var maxAccountLen int = 7
		for account := range state.Balances {
			if len(account) > maxAccountLen {
				maxAccountLen = len(account)
			}
		}
		fmt.Println("Account" + strings.Repeat(" ", maxAccountLen-7) + " : Balance")

		for account, balance := range state.Balances {
			tidyFmt := "%s" + strings.Repeat(" ", maxAccountLen-len(account)) + " : %d"
			fmt.Println(fmt.Sprintf(tidyFmt, account, balance))
		}
	},
}

var balancesStateCmd = &cobra.Command{
	Use:   "state",
	Short: "State of balances (json).",
	Run: func(cmd *cobra.Command, args []string) {
		js := struct {
			ParentHash string
			State      *dao.State
		}{fmt.Sprintf("%x", state.LatestBlockHash()), state}
		json, err := json.MarshalIndent(js, "", "  ")
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(string(json))
	},
}
