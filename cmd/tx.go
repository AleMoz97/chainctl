package cmd

import (
	"context"

	outputpkg "chainctl/internal/output"
	rpcpkg "chainctl/internal/rpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var txCmd = &cobra.Command{
	Use:   "tx [hash]",
	Short: "Mostra dati base di una transazione",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hash := common.HexToHash(args[0])

		ctx := context.Background()
		client, err := rpcpkg.New(ctx, cfg)
		if err != nil {
			return err
		}
		defer client.Close()

		tx, pending, err := client.TransactionByHash(ctx, hash)
		if err != nil {
			return err
		}

		to := "contract creation"
		if tx.To() != nil {
			to = tx.To().Hex()
		}

		return outputpkg.JSON(map[string]any{
			"hash":      tx.Hash().Hex(),
			"nonce":     tx.Nonce(),
			"to":        to,
			"value":     tx.Value().String(),
			"gas":       tx.Gas(),
			"gas_price": tx.GasPrice().String(),
			"pending":   pending,
			"data":      common.Bytes2Hex(tx.Data()),
		})
	},
}

func init() {
	rootCmd.AddCommand(txCmd)
}
