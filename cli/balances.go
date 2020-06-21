package cli

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"simpleblockchain/node"
	"strings"
)

var balances node.BalancesRes

func BalancesCmd() *cobra.Command {
	var balancesCmd = &cobra.Command{
		Use:   "balances",
		Short: "Interact with balances (list...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ErrIncorrectUsage
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			openState()
			var err error
			balances, err = getBalances()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
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
		fmt.Printf("Accounts balances at %s:\n", balances.Hash.Hex())
		var maxAccountLen int = 7
		for account := range balances.Balances {
			if len(account) > maxAccountLen {
				maxAccountLen = len(account)
			}
		}
		fmt.Println("Account" + strings.Repeat(" ", maxAccountLen-7) + " : Balance")

		for account, balance := range balances.Balances {
			tidyFmt := "%s" + strings.Repeat(" ", maxAccountLen-len(account)) + " : %d"
			fmt.Println(fmt.Sprintf(tidyFmt, account, balance))
		}
	},
}

var balancesStateCmd = &cobra.Command{
	Use:   "state",
	Short: "State of balances (json).",
	Run: func(cmd *cobra.Command, args []string) {
		json, err := json.MarshalIndent(balances, "", "  ")
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(string(json))
	},
}

func getBalances() (node.BalancesRes, error) {
	// If we don't have a connection to the server then
	// we directly call the blockchain routines
	if conn == nil {
		return node.BalancesRes{
			Hash:     state.LatestBlockHash(),
			Balances: state.Balances,
		}, nil
	}
	// Get the balancees from the server
	url := fmt.Sprintf("http://%s%s", thisPeerNode.TcpAddress(), node.EndpointBalancesList)
	resp, err := http.Get(url)
	var b node.BalancesRes = node.BalancesRes{}
	if err != nil {
		return b, fmt.Errorf("Error requesting balances: %w", err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&b)
	if err != nil {
		return b, fmt.Errorf("Couldn't decode the json %v: %w", resp.Body, err)
	}
	return b, nil
}
