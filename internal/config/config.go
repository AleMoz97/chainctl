package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

type ValidatorConfig struct {
	ListMethod    string `mapstructure:"list_method"`
	ProposeMethod string `mapstructure:"propose_method"`
}

type Config struct {
	RPCURL              string          `mapstructure:"rpc_url"`
	ChainID             int64           `mapstructure:"chain_id"`
	FromAddress         string          `mapstructure:"from_address"`
	PrivateKey          string          `mapstructure:"private_key"`
	PrivateKeyEnv       string          `mapstructure:"private_key_env"`
	TimeoutSeconds      int             `mapstructure:"timeout_seconds"`
	PollIntervalSeconds int             `mapstructure:"poll_interval_seconds"`
	Validator           ValidatorConfig `mapstructure:"validator"`
}

type Overrides struct {
	RPCURL                 *string
	ChainID                *int64
	FromAddress            *string
	PrivateKey             *string
	PrivateKeyEnv          *string
	TimeoutSeconds         *int
	PollIntervalSeconds    *int
	ValidatorListMethod    *string
	ValidatorProposeMethod *string
}

func Load(configFile string, overrides Overrides) (*Config, error) {
	loadDotEnv(configFile)

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetEnvPrefix("CHAINCTL")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("rpc_url", "http://127.0.0.1:8545")
	v.SetDefault("chain_id", 1337)
	v.SetDefault("private_key", "")
	v.SetDefault("private_key_env", "CHAINCTL_PRIVATE_KEY")
	v.SetDefault("timeout_seconds", 10)
	v.SetDefault("poll_interval_seconds", 3)
	v.SetDefault("validator.list_method", "qbft_getValidatorsByBlockNumber")
	v.SetDefault("validator.propose_method", "clique_propose")

	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.chainctl")
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if configFile != "" || (!os.IsNotExist(err) && !errors.As(err, &notFound)) {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	applyOverrides(v, overrides)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if cfg.RPCURL == "" {
		return nil, fmt.Errorf("rpc_url is required")
	}
	if cfg.FromAddress != "" && !common.IsHexAddress(cfg.FromAddress) {
		return nil, fmt.Errorf("invalid from_address: %s", cfg.FromAddress)
	}
	return &cfg, nil
}

func loadDotEnv(configFile string) {
	if configFile != "" {
		_ = gotenv.Load(filepath.Join(filepath.Dir(configFile), ".env"))
	}

	_ = gotenv.Load(".env")

	homeDir, err := os.UserHomeDir()
	if err == nil {
		_ = gotenv.Load(filepath.Join(homeDir, ".chainctl", ".env"))
	}
}

func applyOverrides(v *viper.Viper, overrides Overrides) {
	if overrides.RPCURL != nil {
		v.Set("rpc_url", *overrides.RPCURL)
	}
	if overrides.ChainID != nil {
		v.Set("chain_id", *overrides.ChainID)
	}
	if overrides.FromAddress != nil {
		v.Set("from_address", *overrides.FromAddress)
	}
	if overrides.PrivateKey != nil {
		v.Set("private_key", *overrides.PrivateKey)
	}
	if overrides.PrivateKeyEnv != nil {
		v.Set("private_key_env", *overrides.PrivateKeyEnv)
	}
	if overrides.TimeoutSeconds != nil {
		v.Set("timeout_seconds", *overrides.TimeoutSeconds)
	}
	if overrides.PollIntervalSeconds != nil {
		v.Set("poll_interval_seconds", *overrides.PollIntervalSeconds)
	}
	if overrides.ValidatorListMethod != nil {
		v.Set("validator.list_method", *overrides.ValidatorListMethod)
	}
	if overrides.ValidatorProposeMethod != nil {
		v.Set("validator.propose_method", *overrides.ValidatorProposeMethod)
	}
}

func FlagOverrides(flags *pflag.FlagSet, values Overrides) Overrides {
	out := Overrides{}

	if flagChanged(flags, "rpc-url") {
		out.RPCURL = values.RPCURL
	}
	if flagChanged(flags, "chain-id") {
		out.ChainID = values.ChainID
	}
	if flagChanged(flags, "from-address") {
		out.FromAddress = values.FromAddress
	}
	if flagChanged(flags, "private-key") {
		out.PrivateKey = values.PrivateKey
	}
	if flagChanged(flags, "private-key-env") {
		out.PrivateKeyEnv = values.PrivateKeyEnv
	}
	if flagChanged(flags, "timeout-seconds") {
		out.TimeoutSeconds = values.TimeoutSeconds
	}
	if flagChanged(flags, "poll-interval-seconds") {
		out.PollIntervalSeconds = values.PollIntervalSeconds
	}
	if flagChanged(flags, "validator-list-method") {
		out.ValidatorListMethod = values.ValidatorListMethod
	}
	if flagChanged(flags, "validator-propose-method") {
		out.ValidatorProposeMethod = values.ValidatorProposeMethod
	}

	return out
}

func flagChanged(flags *pflag.FlagSet, name string) bool {
	flag := flags.Lookup(name)
	return flag != nil && flag.Changed
}
