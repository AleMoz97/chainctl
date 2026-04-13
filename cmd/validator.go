package cmd

import (
	"context"

	outputpkg "chainctl/internal/output"
	rpcpkg "chainctl/internal/rpc"
	validatorpkg "chainctl/internal/validator"

	"github.com/spf13/cobra"
)

var validatorCmd = &cobra.Command{
	Use:   "validator",
	Short: "Operazioni sui validator via RPC",
}

var validatorListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista validator correnti",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := rpcpkg.New(ctx, cfg)
		if err != nil {
			return err
		}
		defer client.Close()

		vals, err := validatorpkg.List(ctx, client, cfg.Validator.ListMethod, "latest")
		if err != nil {
			return err
		}
		return outputpkg.JSON(map[string]any{
			"method":     cfg.Validator.ListMethod,
			"validators": vals,
		})
	},
}
var validatorProposeAddCmd = &cobra.Command{
	Use:   "propose-add [address]",
	Short: "Propone l'aggiunta di un validator",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := rpcpkg.New(ctx, cfg)
		if err != nil {
			return err
		}
		defer client.Close()

		if err := validatorpkg.Propose(ctx, client, cfg.Validator.ProposeMethod, args[0], true); err != nil {
			return err
		}
		return outputpkg.JSON(map[string]any{
			"method":  cfg.Validator.ProposeMethod,
			"address": args[0],
			"vote":    "add",
			"result":  "submitted",
		})
	},
}

var validatorProposeRemoveCmd = &cobra.Command{
	Use:   "propose-remove [address]",
	Short: "Propone la rimozione di un validator",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := rpcpkg.New(ctx, cfg)
		if err != nil {
			return err
		}
		defer client.Close()

		if err := validatorpkg.Propose(ctx, client, cfg.Validator.ProposeMethod, args[0], false); err != nil {
			return err
		}
		return outputpkg.JSON(map[string]any{
			"method":  cfg.Validator.ProposeMethod,
			"address": args[0],
			"vote":    "remove",
			"result":  "submitted",
		})
	},
}

func init() {
	validatorCmd.AddCommand(validatorListCmd)
	validatorCmd.AddCommand(validatorProposeAddCmd)
	validatorCmd.AddCommand(validatorProposeRemoveCmd)
	rootCmd.AddCommand(validatorCmd)
}
