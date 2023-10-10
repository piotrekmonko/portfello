package config

import "github.com/spf13/viper"

type Config struct {
	DatabaseDSN       string
	GraphqlPort       string
	GraphqlPlayground bool
}

func New() *Config {
	return &Config{
		DatabaseDSN:       viper.GetString("database.dsn"),
		GraphqlPort:       viper.GetString("graphql.port"),
		GraphqlPlayground: viper.GetBool("graphql.playground"),
	}
}
