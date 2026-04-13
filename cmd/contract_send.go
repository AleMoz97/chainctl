package cmd

import (
	"context"
	"fmt"
	"math/big"

	contractspkg "chainctl/internal/contracts"
	outputpkg "chainctl/internal/output"
	rpcpkg "chainctl/internal/rpc"
	utilpkg "chainctl/internal/util"
	walletpkg "chainctl/internal/wallet"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var (
	contractABISend     string
	contractAddressSend string
	contractMethodSend  string
	contractArgsSend    []string
	contractValueSend   string
	contractGasLimit    uint64
)

var contractSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Invia una transazione a una funzione write del contratto",
	RunE: func(cmd *cobra.Command, args []string) error {
		abiDef, err := contractspkg.LoadABI(contractABISend)
		if err != nil {
			return err
		}
		if !common.IsHexAddress(contractAddressSend) {
			return fmt.Errorf("invalid contract address: %s", contractAddressSend)
		}
		contractAddr := common.HexToAddress(contractAddressSend)

		valueWei := big.NewInt(0)
		if contractValueSend != "" {
			valueWei, err = utilpkg.WeiFromEtherString(contractValueSend)
			if err != nil {
				return err
			}
		}

		data, err := contractspkg.PackMethod(abiDef, contractMethodSend, contractArgsSend)
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

		txHash, err := contractspkg.SendContractTx(ctx, client, signer, contractAddr, valueWei, data, contractGasLimit)
		if err != nil {
			return err
		}

		return outputpkg.JSON(map[string]any{
			"from":     signer.From.Hex(),
			"contract": contractAddr.Hex(),
			"method":   contractMethodSend,
			"tx_hash":  txHash.Hex(),
		})
	},
}

func init() {
	contractSendCmd.Flags().StringVar(&contractABISend, "abi", "", "path file ABI JSON")
	contractSendCmd.Flags().StringVar(&contractAddressSend, "address", "", "contract address")
	contractSendCmd.Flags().StringVar(&contractMethodSend, "method", "", "method name")
	contractSendCmd.Flags().StringSliceVar(&contractArgsSend, "args", nil, "method args, es. --args a,b,c")
	contractSendCmd.Flags().StringVar(&contractValueSend, "value", "0", "ether da inviare insieme alla call")
	contractSendCmd.Flags().Uint64Var(&contractGasLimit, "gas-limit", 0, "optional gas limit")
	_ = contractSendCmd.MarkFlagRequired("abi")
	_ = contractSendCmd.MarkFlagRequired("address")
	_ = contractSendCmd.MarkFlagRequired("method")

	contractCmd.AddCommand(contractSendCmd)
}
