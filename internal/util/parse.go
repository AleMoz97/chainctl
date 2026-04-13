package util

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

func MustAddress(s string) (common.Address, error) {
	if !common.IsHexAddress(s) {
		return common.Address{}, fmt.Errorf("invalid address: %s", s)
	}
	return common.HexToAddress(s), nil
}

func ParseBigInt(s string) (*big.Int, error) {
	z, ok := new(big.Int).SetString(strings.TrimSpace(s), 10)
	if !ok {
		return nil, fmt.Errorf("invalid integer: %s", s)
	}
	return z, nil
}

func WeiFromEtherString(s string) (*big.Int, error) {
	parts := strings.SplitN(strings.TrimSpace(s), ".", 2)
	if len(parts) == 1 {
		i, err := ParseBigInt(parts[0])
		if err != nil {
			return nil, err
		}
		return new(big.Int).Mul(i, big.NewInt(1e18)), nil
	}

	whole, err := ParseBigInt(parts[0])
	if err != nil {
		return nil, err
	}

	frac := parts[1]
	if len(frac) > 18 {
		return nil, fmt.Errorf("too many decimal places: max 18")
	}
	frac = frac + strings.Repeat("0", 18-len(frac))
	fracInt, err := ParseBigInt(frac)
	if err != nil {
		return nil, err
	}

	result := new(big.Int).Mul(whole, big.NewInt(1e18))
	result.Add(result, fracInt)
	return result, nil
}
