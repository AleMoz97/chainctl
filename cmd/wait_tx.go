package cmd

import (
	"context"
	"errors"
	"time"

	outputpkg "chainctl/internal/output"
	rpcpkg "chainctl/internal/rpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var waitTxTimeoutSec int

var waitTxCmd = &cobra.Command{
	Use:   "wait-tx [hash]",
	Short: "Attende finché una tx non viene minata",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hash := common.HexToHash(args[0])
		ctx := context.Background()
		client, err := rpcpkg.New(ctx, cfg)
		if err != nil {
			return err
		}
		defer client.Close()

		deadline := time.Now().Add(time.Duration(waitTxTimeoutSec) * time.Second)
		for {
			receipt, err := client.TransactionReceipt(ctx, hash)
			if err == nil && receipt != nil {
				return outputpkg.JSON(map[string]any{
					"tx_hash":      receipt.TxHash.Hex(),
					"block_number": receipt.BlockNumber.String(),
					"status":       receipt.Status,
					"gas_used":     receipt.GasUsed,
				})
			}

			if time.Now().After(deadline) {
				return errors.New("timeout waiting for receipt")
			}
			time.Sleep(time.Duration(cfg.PollIntervalSeconds) * time.Second)
		}
	},
}

func init() {
	waitTxCmd.Flags().IntVar(&waitTxTimeoutSec, "timeout", 120, "timeout in seconds")
	rootCmd.AddCommand(waitTxCmd)
}
