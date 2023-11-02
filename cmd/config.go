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
	"github.com/piotrekmonko/portfello/pkg/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"conf"},
	Short:   "Config verifies configuration is complete",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.New()
		err := conf.Validate()
		cobra.CheckErr(err)

		if showConf, _ := cmd.Flags().GetBool("show"); showConf {
			pretty, err := json.MarshalIndent(conf, "", "  ")
			cobra.CheckErr(err)
			fmt.Printf("Parsed configuration values:\n%s\n", string(pretty))
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolP("show", "s", false, "Show the parsed config values")
}
