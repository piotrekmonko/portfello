package cmd

import (
	"fmt"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

const baseVersion = "v1.0.0"

var (
	cfgFile     string
	buildNumber = "dev"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "portfello",
	Short:   "Backend services for PortfelloApp",
	Long:    `PortfelloApp is an opensource project for managing your household budget.`,
	Version: "v1.0.0-dev",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		lvl, err := cmd.Flags().GetString("level")
		if err != nil {
			return fmt.Errorf("cannot read log level: %w", err)
		}
		return logz.ParseFlag(lvl)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Version = fmt.Sprintf("%s-%s", baseVersion, buildNumber)
	logz.SetVer(rootCmd.Version)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.portfello.yaml)")
	rootCmd.PersistentFlags().String("level", zap.InfoLevel.String(),
		"set logger level (one of: debug, info, warn, error, fatal)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cobra.CheckErr(conf.InitConfig(cfgFile, "PORTFELLO"))
}
