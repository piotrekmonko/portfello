package conf

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
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
	AuthProviderMock  = "mock"
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

// NewTestConfig returns configuration suitable for testing.
func NewTestConfig() *Config {
	return &Config{
		DatabaseDSN: "",
		Graph:       GraphQL{},
		Auth:        Auth0{},
		Logging:     Logging{},
	}
}

// InitConfig reads in config file and ENV variables if set. Should be called once at app startup.
func InitConfig(cfgFilePath string) error {
	if cfgFilePath != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFilePath)
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("fatal error config file @ %s: %w", cfgFilePath, err)
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

	return nil
}
