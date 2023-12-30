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
	Logging     Logging `yaml:"logging"`
}

type GraphQL struct {
	Port             string `yaml:"port"`
	EnablePlayground bool   `yaml:"enable_playground"`
}

const (
	AuthProviderLocal = "local"
	AuthProviderAuth0 = "auth0"
)

type Auth0 struct {
	Provider     string `yaml:"provider"`
	Domain       string `yaml:"domain"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Audience     string `yaml:"audience"`
	ConnectionID string `yaml:"connection_id"`
}

type Logging struct {
	Level  string `json:"level"`
	Format string `json:"format"`
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
	if c.Auth.Provider != AuthProviderAuth0 && c.Auth.Provider != AuthProviderLocal {
		return fmt.Errorf("auth.provider must be set to either '%s' or '%s'", AuthProviderLocal, AuthProviderAuth0)
	}

	return nil
}
