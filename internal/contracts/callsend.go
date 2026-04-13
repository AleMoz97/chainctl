package contracts

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	rpcpkg "chainctl/internal/rpc"
	walletpkg "chainctl/internal/wallet"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func parseArg(input abi.Argument, raw string) (any, error) {
	switch input.Type.T {
	case abi.StringTy:
		return raw, nil
	case abi.BoolTy:
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, fmt.Errorf("parse bool %q: %w", raw, err)
		}
		return v, nil
	case abi.AddressTy:
		if !common.IsHexAddress(raw) {
			return nil, fmt.Errorf("invalid address: %s", raw)
		}
		return common.HexToAddress(raw), nil
	case abi.UintTy, abi.IntTy:
		z, ok := new(big.Int).SetString(raw, 10)
		if !ok {
			return nil, fmt.Errorf("invalid integer: %s", raw)
		}
		return z, nil
	default:
		return nil, fmt.Errorf("unsupported ABI arg type: %s", input.Type.String())
	}
}

func PackMethod(a *abi.ABI, method string, rawArgs []string) ([]byte, error) {
	m, ok := a.Methods[method]
	if !ok {
		return nil, fmt.Errorf("method %s not found in ABI", method)
	}
	if len(rawArgs) != len(m.Inputs) {
		return nil, fmt.Errorf("method %s expects %d args, got %d", method, len(m.Inputs), len(rawArgs))
	}

	args := make([]any, 0, len(rawArgs))
	for i, in := range m.Inputs {
		val, err := parseArg(in, strings.TrimSpace(rawArgs[i]))
		if err != nil {
			return nil, fmt.Errorf("arg %d (%s): %w", i, in.Name, err)
		}
		args = append(args, val)
	}

	data, err := a.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("abi pack: %w", err)
	}
	return data, nil
}

func CallContract(ctx context.Context, client *rpcpkg.Client, contract common.Address, data []byte) ([]any, error) {
	callCtx, cancel := client.WithTimeout(ctx)
	defer cancel()

	res, err := client.Eth.CallContract(callCtx, ethereum.CallMsg{
		To:   &contract,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("eth_call: %w", err)
	}
	return []any{res}, nil
}

func SendContractTx(
	ctx context.Context,
	client *rpcpkg.Client,
	signer *walletpkg.Signer,
	contract common.Address,
	value *big.Int,
	data []byte,
	gasLimit uint64,
) (common.Hash, error) {
	signed, err := signer.BuildAndSignLegacyTx(ctx, client, &contract, value, data, gasLimit, nil)
	if err != nil {
		return common.Hash{}, err
	}

	sendCtx, cancel := client.WithTimeout(ctx)
	defer cancel()
	if err := client.Eth.SendTransaction(sendCtx, signed); err != nil {
		return common.Hash{}, fmt.Errorf("send transaction: %w", err)
	}
	return signed.Hash(), nil
}
