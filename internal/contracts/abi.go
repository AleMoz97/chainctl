package contracts

import (
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func LoadABI(path string) (*abi.ABI, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read abi file: %w", err)
	}
	parsed, err := abi.JSON(strings.NewReader(string(b)))
	if err != nil {
		return nil, fmt.Errorf("parse abi: %w", err)
	}
	return &parsed, nil
}
