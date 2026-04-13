package validator

import (
	"context"
	"fmt"

	rpcpkg "chainctl/internal/rpc"
)

func List(ctx context.Context, client *rpcpkg.Client, method string, blockTag string) ([]string, error) {
	callCtx, cancel := client.WithTimeout(ctx)
	defer cancel()

	var out []string
	if err := client.RPC.CallContext(callCtx, &out, method, blockTag); err != nil {
		return nil, fmt.Errorf("validator list via %s: %w", method, err)
	}
	return out, nil
}

func Propose(ctx context.Context, client *rpcpkg.Client, method string, address string, authorize bool) error {
	callCtx, cancel := client.WithTimeout(ctx)
	defer cancel()

	var result any
	if err := client.RPC.CallContext(callCtx, &result, method, address, authorize); err != nil {
		return fmt.Errorf("validator propose via %s: %w", method, err)
	}
	return nil
}
