package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DatabaseDSN string  `yaml:"database_dsn"`
	Graph       GraphQL `yaml:"graphql"`
	Auth        Auth0   `yaml:"auth"`
}

type GraphQL struct {
	Port             string `yaml:"port"`
	EnablePlayground bool   `yaml:"enable_playground"`
}

type Auth0 struct {
	Domain       string `yaml:"domain"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Audience     string `yaml:"audience"`
	ConnectionID string `yaml:"connection_id"`
}

func New() *Config {
	c := &Config{}

	if err := viper.Unmarshal(c, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "yaml"
	}); err != nil {
		log.Fatal(err)
	}

	return c
}

func (c *Config) Validate() error {
	if c.DatabaseDSN == "" {
		return fmt.Errorf("database_dsn is required")
	}
	return nil
}
