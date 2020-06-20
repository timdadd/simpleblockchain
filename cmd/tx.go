package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"simpleblockchain/dao"
)

const flagFrom = "from"
const flagTo = "to"
const flagValue = "value"
const flagData = "data"

func TxCmd() *cobra.Command {
	var txsCmd = &cobra.Command{
		Use:   "tx",
		Short: "Interact with txs (add...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return IncorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	txsCmd.AddCommand(txAddCmd())

	return txsCmd
}

func txAddCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "add",
		Short: "Adds new TX to database.",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString(flagFrom)
			to, _ := cmd.Flags().GetString(flagTo)
			value, _ := cmd.Flags().GetUint(flagValue)
			data, _ := cmd.Flags().GetString(flagData)

			tx := dao.NewTx(dao.NewAccount(from), dao.NewAccount(to), value, data)

			state, err := dao.NewStateFromDisk()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()

			err = state.Add(tx)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			_, err = state.Persist()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			fmt.Println("TX successfully added to the ledger.")
		},
	}

	cmd.Flags().String(flagFrom, "", "From what account to send tokens")
	cmd.MarkFlagRequired(flagFrom)

	cmd.Flags().String(flagTo, "", "To what account to send tokens")
	cmd.MarkFlagRequired(flagTo)

	cmd.Flags().Uint(flagValue, 0, "How many tokens to send")
	cmd.MarkFlagRequired(flagValue)

	cmd.Flags().String(flagData, "", "e.g.: 'reward', 'services',' vodka' ...")

	return cmd
}
