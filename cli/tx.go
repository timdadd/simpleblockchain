package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"simpleblockchain/dao"
	"simpleblockchain/node"
	"time"
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
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			openState()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			closeState()
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
			var txAddRes node.TxAddRes

			// If we don't have a connection to the server then
			// we directly call the blockchain routines
			if conn == nil {
				block := dao.NewBlock(
					state.LatestBlockHash(),
					state.NextBlockNumber(),
					uint64(time.Now().Unix()),
					[]dao.Tx{tx},
				)
				var err error
				txAddRes.Hash, err = state.AddBlock(block)
				if err != nil {
					_, _ = fmt.Fprintln(os.Stderr, err)
					return
				}
			} else {
				// Send the request to the server
				url := fmt.Sprintf("http://%s%s", thisPeerNode.TcpAddress(), node.EndpointTxAdd)
				jsonTx, err := json.Marshal(tx)
				if err != nil {
					_, _ = fmt.Fprintln(os.Stderr, err)
					return
				}
				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonTx))
				if err != nil {
					_, _ = fmt.Fprintln(os.Stderr, err)
					return
				}
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req) // Send the post and hopefully get a response
				if err != nil {
					_, _ = fmt.Fprintln(os.Stderr, err)
					return
				}
				defer resp.Body.Close()
				err = json.NewDecoder(resp.Body).Decode(&txAddRes)
				if err != nil {
					_, _ = fmt.Fprintln(os.Stderr, err)
					return
				}

				//fmt.Println("response Status:", resp.Status)
				//fmt.Println("response Headers:", resp.Header)
				//body, _ := ioutil.ReadAll(resp.Body)
				//err = json.Unmarshal(body, &txAddRes)
				//fmt.Println("response Body:", string(body))
				//if err != nil {
				//	_, _ = fmt.Fprintln(os.Stderr, err)
				//	return
				//}

			}

			fmt.Printf("TX successfully added to the ledger. %v\n", txAddRes.Hash.Hex())
		},
	}

	cmd.Flags().String(flagFrom, "", "From what account to send tokens")
	_ = cmd.MarkFlagRequired(flagFrom)

	cmd.Flags().String(flagTo, "", "To what account to send tokens")
	_ = cmd.MarkFlagRequired(flagTo)

	cmd.Flags().Uint(flagValue, 0, "How many tokens to send")
	_ = cmd.MarkFlagRequired(flagValue)

	_ = cmd.Flags().String(flagData, "", "e.g.: 'reward', 'services',' vodka' ...")

	return cmd
}
