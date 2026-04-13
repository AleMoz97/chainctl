package cmd

import (
	"context"
	"math/big"

	outputpkg "chainctl/internal/output"
	rpcpkg "chainctl/internal/rpc"
	utilpkg "chainctl/internal/util"

	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use:   "balance [address]",
	Short: "Legge il balance di un address",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		addr, err := utilpkg.MustAddress(args[0])
		if err != nil {
			return err
		}

		ctx := context.Background()
		client, err := rpcpkg.New(ctx, cfg)
		if err != nil {
			return err
		}
		defer client.Close()

		bal, err := client.BalanceAt(ctx, addr)
		if err != nil {
			return err
		}

		ether := new(big.Rat).SetFrac(bal, big.NewInt(1e18))
		return outputpkg.JSON(map[string]any{
			"address": addr.Hex(),
			"wei":     bal.String(),
			"ether":   ether.FloatString(18),
		})
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}
