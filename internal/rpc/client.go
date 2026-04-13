package rpc

import (
	"context"
	"fmt"
	"math/big"
	"time"

	cfgpkg "chainctl/internal/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	Eth     *ethclient.Client
	RPC     *gethrpc.Client
	Timeout time.Duration
	Cfg     *cfgpkg.Config
}

func New(ctx context.Context, cfg *cfgpkg.Config) (*Client, error) {
	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	dialCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	raw, err := gethrpc.DialContext(dialCtx, cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("dial rpc: %w", err)
	}

	return &Client{
		Eth:     ethclient.NewClient(raw),
		RPC:     raw,
		Timeout: timeout,
		Cfg:     cfg,
	}, nil
}

func (c *Client) Close() {
	if c.RPC != nil {
		c.RPC.Close()
	}
}

func (c *Client) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, c.Timeout)
}

func (c *Client) ChainID(ctx context.Context) (*big.Int, error) {
	return c.Eth.ChainID(ctx)
}

func (c *Client) BlockNumber(ctx context.Context) (uint64, error) {
	return c.Eth.BlockNumber(ctx)
}

func (c *Client) BalanceAt(ctx context.Context, addr common.Address) (*big.Int, error) {
	return c.Eth.BalanceAt(ctx, addr, nil)
}

func (c *Client) TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	return c.Eth.TransactionByHash(ctx, hash)
}

func (c *Client) TransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	return c.Eth.TransactionReceipt(ctx, hash)
}
