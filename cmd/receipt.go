package cmd

import (
	"context"

	outputpkg "chainctl/internal/output"
	rpcpkg "chainctl/internal/rpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var receiptCmd = &cobra.Command{
	Use:   "receipt [hash]",
	Short: "Mostra la receipt di una transazione",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hash := common.HexToHash(args[0])

		ctx := context.Background()
		client, err := rpcpkg.New(ctx, cfg)
		if err != nil {
			return err
		}
		defer client.Close()

		receipt, err := client.TransactionReceipt(ctx, hash)
		if err != nil {
			return err
		}

		return outputpkg.JSON(map[string]any{
			"tx_hash":             receipt.TxHash.Hex(),
			"block_hash":          receipt.BlockHash.Hex(),
			"block_number":        receipt.BlockNumber.String(),
			"status":              receipt.Status,
			"gas_used":            receipt.GasUsed,
			"cumulative_gas_used": receipt.CumulativeGasUsed,
			"contract_address":    receipt.ContractAddress.Hex(),
			"logs_count":          len(receipt.Logs),
		})
	},
}

func init() {
	rootCmd.AddCommand(receiptCmd)
}
