package cmd

import (
	"context"
	"fmt"

	contractspkg "chainctl/internal/contracts"
	outputpkg "chainctl/internal/output"
	rpcpkg "chainctl/internal/rpc"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var (
	contractABICall     string
	contractAddressCall string
	contractMethodCall  string
	contractArgsCall    []string
)

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Interazioni con smart contract via ABI",
}

var contractCallCmd = &cobra.Command{
	Use:   "call",
	Short: "Chiama una funzione read-only del contratto",
	RunE: func(cmd *cobra.Command, args []string) error {
		abiDef, err := contractspkg.LoadABI(contractABICall)
		if err != nil {
			return err
		}
		if !common.IsHexAddress(contractAddressCall) {
			return fmt.Errorf("invalid contract address: %s", contractAddressCall)
		}
		contractAddr := common.HexToAddress(contractAddressCall)

		data, err := contractspkg.PackMethod(abiDef, contractMethodCall, contractArgsCall)
		if err != nil {
			return err
		}

		ctx := context.Background()
		client, err := rpcpkg.New(ctx, cfg)
		if err != nil {
			return err
		}
		defer client.Close()

		callCtx, cancel := client.WithTimeout(ctx)
		defer cancel()
		raw, err := client.Eth.CallContract(callCtx, ethereum.CallMsg{To: &contractAddr, Data: data}, nil)
		if err != nil {
			return err
		}

		method, ok := abiDef.Methods[contractMethodCall]
		if !ok {
			return fmt.Errorf("method %s not found", contractMethodCall)
		}
		unpacked, err := method.Outputs.Unpack(raw)
		if err != nil {
			return outputpkg.JSON(map[string]any{"raw": common.Bytes2Hex(raw)})
		}
		return outputpkg.JSON(map[string]any{
			"contract": contractAddr.Hex(),
			"method":   contractMethodCall,
			"outputs":  unpacked,
			"raw":      common.Bytes2Hex(raw),
		})
	},
}

func init() {
	contractCallCmd.Flags().StringVar(&contractABICall, "abi", "", "path file ABI JSON")
	contractCallCmd.Flags().StringVar(&contractAddressCall, "address", "", "contract address")
	contractCallCmd.Flags().StringVar(&contractMethodCall, "method", "", "method name")
	contractCallCmd.Flags().StringSliceVar(&contractArgsCall, "args", nil, "method args, es. --args a,b,c")
	_ = contractCallCmd.MarkFlagRequired("abi")
	_ = contractCallCmd.MarkFlagRequired("address")
	_ = contractCallCmd.MarkFlagRequired("method")

	contractCmd.AddCommand(contractCallCmd)
	rootCmd.AddCommand(contractCmd)
}
