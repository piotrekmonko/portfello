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
