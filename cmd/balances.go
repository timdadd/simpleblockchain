package cmd

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
		state, err := dao.NewStateFromDisk()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer state.Close()

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
		//
		//fmt.Println("Account: balance:")
		//fmt.Println("__________________")
		//for account, balance := range state.Balances {
		//	fmt.Println(account(string) + strings.Repeat(" ", maxAccountLen - len(account)) +":" + balance)
		//	fmt.Println(fmt.Sprintf("%s: %d", account, balance))
		//}
	},
}

var balancesStateCmd = &cobra.Command{
	Use:   "state",
	Short: "State of balances (json).",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := dao.NewStateFromDisk()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer state.Close()

		json, err := json.MarshalIndent(state, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(string(json))
	},
}
