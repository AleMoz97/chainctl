package cmd

import (
	"fmt"
	"os"

	cfgpkg "chainctl/internal/config"
	versionpkg "chainctl/internal/version"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *cfgpkg.Config

	rpcURLFlag                 string
	chainIDFlag                int64
	fromAddressFlag            string
	privateKeyFlag             string
	privateKeyEnvFlag          string
	timeoutSecondsFlag         int
	pollIntervalSecondsFlag    int
	validatorListMethodFlag    string
	validatorProposeMethodFlag string
)

var rootCmd = &cobra.Command{
	Use:   "chainctl",
	Short: "CLI per operazioni su nodo EVM/Quorum/Besu",
	Long: `chainctl legge la configurazione da piu fonti con questa precedenza:
1. flag da linea di comando
2. variabili d'ambiente e file .env
3. config.yaml
4. valori di default`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "version" {
			return nil
		}

		versionFlag, err := cmd.Flags().GetBool("version")
		if err == nil && versionFlag {
			return nil
		}

		overrides := cfgpkg.FlagOverrides(cmd.Flags(), cfgpkg.Overrides{
			RPCURL:                 &rpcURLFlag,
			ChainID:                &chainIDFlag,
			FromAddress:            &fromAddressFlag,
			PrivateKey:             &privateKeyFlag,
			PrivateKeyEnv:          &privateKeyEnvFlag,
			TimeoutSeconds:         &timeoutSecondsFlag,
			PollIntervalSeconds:    &pollIntervalSecondsFlag,
			ValidatorListMethod:    &validatorListMethodFlag,
			ValidatorProposeMethod: &validatorProposeMethodFlag,
		})

		loaded, err := cfgpkg.Load(cfgFile, overrides)
		if err != nil {
			return err
		}
		cfg = loaded
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		versionFlag, err := cmd.Flags().GetBool("version")
		if err == nil && versionFlag {
			fmt.Println(versionpkg.Value)
			return nil
		}
		return cmd.Help()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Mostra la versione installata di chainctl",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(versionpkg.Value)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().Bool("version", false, "mostra la versione installata")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path al file di configurazione")
	rootCmd.PersistentFlags().StringVar(&rpcURLFlag, "rpc-url", "", "RPC URL del nodo")
	rootCmd.PersistentFlags().Int64Var(&chainIDFlag, "chain-id", 0, "chain ID EVM")
	rootCmd.PersistentFlags().StringVar(&fromAddressFlag, "from-address", "", "address mittente atteso")
	rootCmd.PersistentFlags().StringVar(&privateKeyFlag, "private-key", "", "private key esadecimale diretta (sconsigliata in produzione)")
	rootCmd.PersistentFlags().StringVar(&privateKeyEnvFlag, "private-key-env", "", "nome della variabile ambiente da cui leggere la private key")
	rootCmd.PersistentFlags().IntVar(&timeoutSecondsFlag, "timeout-seconds", 0, "timeout RPC in secondi")
	rootCmd.PersistentFlags().IntVar(&pollIntervalSecondsFlag, "poll-interval-seconds", 0, "intervallo di polling in secondi")
	rootCmd.PersistentFlags().StringVar(&validatorListMethodFlag, "validator-list-method", "", "RPC method per listare i validator")
	rootCmd.PersistentFlags().StringVar(&validatorProposeMethodFlag, "validator-propose-method", "", "RPC method per proporre add/remove validator")
	rootCmd.AddCommand(versionCmd)
}
