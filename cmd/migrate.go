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
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/piotrekmonko/portfello/dbschema"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	migrateCmd.AddCommand(upCmd)
	migrateCmd.AddCommand(downCmd)
	migrateCmd.AddCommand(dropCmd)
	rootCmd.AddCommand(migrateCmd)
}

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "A tool to migrate your database",
}

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply missing migrations",
	Run: func(cmd *cobra.Command, args []string) {
		conf := conf.New()
		migrator, err := dbschema.NewMigrator(conf.DatabaseDSN)
		if err != nil {
			logFatalMigrate(err)
		}

		err = migrator.Up()
		if err != nil {
			logFatalMigrate(err)
		}
	},
}

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Revert one last migration",
	Run: func(cmd *cobra.Command, args []string) {
		conf := conf.New()
		migrator, err := dbschema.NewMigrator(conf.DatabaseDSN)
		if err != nil {
			logFatalMigrate(err)
		}

		err = migrator.Steps(-1)
		if err != nil {
			logFatalMigrate(err)
		}
	},
}

// dropCmd represents the drop command
var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop database. WARNING: This will delete all data!",
	Run: func(cmd *cobra.Command, args []string) {
		conf := conf.New()
		migrator, err := dbschema.NewMigrator(conf.DatabaseDSN)
		if err != nil {
			logFatalMigrate(err)
		}

		err = migrator.Drop()
		if err != nil {
			logFatalMigrate(err)
		}
	},
}

func logFatalMigrate(err error) {
	if err == nil {
		return
	}

	if errors.Is(err, migrate.ErrNoChange) {
		// ErrNoChange is not an error, log it as info
		log.Printf("database: %s", err.Error())
		return
	}

	log.Fatalln(err)
}
