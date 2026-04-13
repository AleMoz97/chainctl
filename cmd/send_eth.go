package cmd

import (
	"context"

	outputpkg "chainctl/internal/output"
	rpcpkg "chainctl/internal/rpc"
	utilpkg "chainctl/internal/util"
	walletpkg "chainctl/internal/wallet"

	"github.com/spf13/cobra"
)

var (
	sendEthTo       string
	sendEthValue    string
	sendEthGasLimit uint64
)

var sendEthCmd = &cobra.Command{
	Use:   "send-eth",
	Short: "Invia una transazione ETH semplice",
	RunE: func(cmd *cobra.Command, args []string) error {
		to, err := utilpkg.MustAddress(sendEthTo)
		if err != nil {
			return err
		}
		valueWei, err := utilpkg.WeiFromEtherString(sendEthValue)
		if err != nil {
			return err
		}

		ctx := context.Background()
		client, err := rpcpkg.New(ctx, cfg)
		if err != nil {
			return err
		}
		defer client.Close()

		signer, err := walletpkg.LoadSigner(cfg)
		if err != nil {
			return err
		}

		tx, err := signer.BuildAndSignLegacyTx(ctx, client, &to, valueWei, nil, sendEthGasLimit, nil)
		if err != nil {
			return err
		}

		sendCtx, cancel := client.WithTimeout(ctx)
		defer cancel()
		if err := client.Eth.SendTransaction(sendCtx, tx); err != nil {
			return err
		}

		return outputpkg.JSON(map[string]any{
			"from":      signer.From.Hex(),
			"to":        to.Hex(),
			"value_wei": valueWei.String(),
			"tx_hash":   tx.Hash().Hex(),
			"gas":       tx.Gas(),
			"gas_price": tx.GasPrice().String(),
		})
	},
}

func init() {
	sendEthCmd.Flags().StringVar(&sendEthTo, "to", "", "destination address")
	sendEthCmd.Flags().StringVar(&sendEthValue, "value", "0", "amount in ether, es. 0.1")
	sendEthCmd.Flags().Uint64Var(&sendEthGasLimit, "gas-limit", 0, "optional gas limit")
	_ = sendEthCmd.MarkFlagRequired("to")
	_ = sendEthCmd.MarkFlagRequired("value")
	rootCmd.AddCommand(sendEthCmd)
}
