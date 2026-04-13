package cmd

import (
	"context"

	outputpkg "chainctl/internal/output"
	rpcpkg "chainctl/internal/rpc"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Mostra stato base del nodo",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := rpcpkg.New(ctx, cfg)
		if err != nil {
			return err
		}
		defer client.Close()

		chainID, err := client.ChainID(ctx)
		if err != nil {
			return err
		}
		blockNumber, err := client.BlockNumber(ctx)
		if err != nil {
			return err
		}

		return outputpkg.JSON(map[string]any{
			"rpc_url":      cfg.RPCURL,
			"chain_id":     chainID.String(),
			"block_number": blockNumber,
		})
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
