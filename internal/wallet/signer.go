package wallet

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"strings"

	cfgpkg "chainctl/internal/config"
	rpcpkg "chainctl/internal/rpc"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Signer struct {
	PrivateKey *ecdsa.PrivateKey
	From       common.Address
}

func LoadSigner(cfg *cfgpkg.Config) (*Signer, error) {
	raw := strings.TrimSpace(cfg.PrivateKey)
	if raw == "" {
		raw = strings.TrimSpace(os.Getenv(cfg.PrivateKeyEnv))
	}
	if raw == "" {
		return nil, fmt.Errorf("private key is empty: pass --private-key or set %s", cfg.PrivateKeyEnv)
	}
	raw = strings.TrimPrefix(raw, "0x")

	pk, err := crypto.HexToECDSA(raw)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	pubKey, ok := pk.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid public key type")
	}
	from := crypto.PubkeyToAddress(*pubKey)

	if cfg.FromAddress != "" && !strings.EqualFold(from.Hex(), cfg.FromAddress) {
		return nil, fmt.Errorf("private key address %s does not match configured from_address %s", from.Hex(), cfg.FromAddress)
	}

	return &Signer{PrivateKey: pk, From: from}, nil
}

func (s *Signer) BuildAndSignLegacyTx(
	ctx context.Context,
	client *rpcpkg.Client,
	to *common.Address,
	value *big.Int,
	data []byte,
	gasLimit uint64,
	gasPrice *big.Int,
) (*types.Transaction, error) {
	nonceCtx, cancel := client.WithTimeout(ctx)
	defer cancel()

	nonce, err := client.Eth.PendingNonceAt(nonceCtx, s.From)
	if err != nil {
		return nil, fmt.Errorf("get nonce: %w", err)
	}

	if gasPrice == nil {
		gpCtx, gpCancel := client.WithTimeout(ctx)
		defer gpCancel()
		gasPrice, err = client.Eth.SuggestGasPrice(gpCtx)
		if err != nil {
			return nil, fmt.Errorf("suggest gas price: %w", err)
		}
	}

	if gasLimit == 0 {
		estCtx, estCancel := client.WithTimeout(ctx)
		defer estCancel()
		msg := ethereum.CallMsg{
			From:     s.From,
			To:       to,
			GasPrice: gasPrice,
			Value:    value,
			Data:     data,
		}
		gasLimit, err = client.Eth.EstimateGas(estCtx, msg)
		if err != nil {
			return nil, fmt.Errorf("estimate gas: %w", err)
		}
	}

	chainID := big.NewInt(client.Cfg.ChainID)
	tx := types.NewTransaction(nonce, derefAddr(to), value, gasLimit, gasPrice, data)
	signed, err := types.SignTx(tx, types.NewEIP155Signer(chainID), s.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("sign tx: %w", err)
	}
	return signed, nil
}

func derefAddr(a *common.Address) common.Address {
	if a == nil {
		return common.Address{}
	}
	return *a
}
