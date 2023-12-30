/*
Copyright © 2023 Piotr Mońko <piotrek.monko@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"go.uber.org/zap"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.portfello.yaml)")
	rootCmd.PersistentFlags().String("level", zap.InfoLevel.String(),
		"set logger level (one of: debug, info, warn, error, fatal)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			cobra.CheckErr(fmt.Errorf("fatal error config file @ %s: %w", cfgFile, err))
		}
	}

	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Search config in home directory with name ".portfello" (without extension).
	viper.AddConfigPath(home)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv() // read in environment variables that match

	// Loop through found config files until all are parsed
	for _, configFile := range []string{".portfello", ".portfello-local"} {
		viper.SetConfigName(configFile)
		if err := viper.MergeInConfig(); err == nil {
			_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
