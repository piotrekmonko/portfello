package cmd

import (
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/spf13/cobra"
)

// provisionCmd represents the provision command
var provisionCmd = &cobra.Command{
	Use:     "provision",
	Aliases: []string{"prov"},
	Short:   "Add objects to systems",
	Run: func(cmd *cobra.Command, args []string) {
		c := conf.New()
		provisioner, cleanup, err := initializeProvisioner(cmd.Context(), c)
		cobra.CheckErr(err)
		defer cleanup()

		err = provisioner.HandleCommand(cmd)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(provisionCmd)

	provisionCmd.Flags().StringP("user", "u", "", "Add a new user")
	provisionCmd.Flags().BoolP("test", "t", false, "Add example test data to database")
}
