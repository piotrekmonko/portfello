package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"conf"},
	Short:   "Config verifies configuration is complete",
	RunE: func(cmd *cobra.Command, _ []string) error {
		conf := conf.New()
		err := conf.Validate()
		if err != nil {
			return err
		}

		if showConf, _ := cmd.Flags().GetBool("show-json"); showConf {
			pretty, err := json.MarshalIndent(conf, "", "  ")
			if err != nil {
				return err
			}

			fmt.Printf("Parsed configuration values:\n%s\n", string(pretty))
		}

		if showConf, _ := cmd.Flags().GetBool("show-yaml"); showConf {
			pretty, err := yaml.Marshal(conf)
			if err != nil {
				return err
			}

			fmt.Printf("Parsed configuration values:\n\n%s\n", string(pretty))
		}

		if showRoutes, _ := cmd.Flags().GetBool("routes"); showRoutes {
			router, closer, err := initializeRouter(cmd.Context(), conf)
			if err != nil {
				return err
			}
			defer closer()

			pretty, err := json.MarshalIndent(router, "", "  ")
			if err != nil {
				return err
			}

			fmt.Printf("Available routes:\n%s\n", string(pretty))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolP("show-json", "s", false, "Show the parsed config values as json")
	configCmd.Flags().BoolP("show-yaml", "y", false, "Show the parsed config values as yaml")
	configCmd.Flags().BoolP("routes", "r", false, "List available routes")
}
