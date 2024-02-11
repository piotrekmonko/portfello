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
	"encoding/json"
	"fmt"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/spf13/cobra"
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

		if showConf, _ := cmd.Flags().GetBool("show"); showConf {
			pretty, err := json.MarshalIndent(conf, "", "  ")
			if err != nil {
				return err
			}

			fmt.Printf("Parsed configuration values:\n%s\n", string(pretty))
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
	configCmd.Flags().BoolP("show", "s", false, "Show the parsed config values")
	configCmd.Flags().BoolP("routes", "r", false, "List available routes")
}
