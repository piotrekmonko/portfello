package conf

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

type Config struct {
	DatabaseDSN string  `yaml:"database_dsn" mapstructure:"database_dsn"`
	HostName    string  `yaml:"host_name" mapstructure:"host_name"`
	Graph       GraphQL `yaml:"graphql" mapstructure:"graphql"`
	Auth        Auth0   `yaml:"auth" mapstructure:"auth"`
	Logging     Logging `yaml:"logging" mapstructure:"logging"`
}

type GraphQL struct {
	Port             string `yaml:"port" mapstructure:"port"`
	EnablePlayground bool   `yaml:"enable_playground" mapstructure:"enable_playground"`
}

const (
	AuthProviderLocal = "local"
	AuthProviderAuth0 = "auth0"
	AuthProviderMock  = "mock"
)

type Auth0 struct {
	Provider     string `yaml:"provider" mapstructure:"provider"`
	Domain       string `yaml:"domain" mapstructure:"domain"`
	ClientID     string `yaml:"client_id" mapstructure:"client_id"`
	ClientSecret string `yaml:"client_secret" mapstructure:"client_secret"`
	Audience     string `yaml:"audience" mapstructure:"audience"`
	ConnectionID string `yaml:"connection_id" mapstructure:"connection_id"`
}

type Logging struct {
	Level  string `yaml:"level" mapstructure:"level"`
	Format string `yaml:"format" mapstructure:"format"`
}

func New() *Config {
	c := &Config{}

	if err := viper.Unmarshal(c); err != nil {
		log.Fatal(err)
	}

	return c
}

func (c *Config) Validate() error {
	if c.HostName == "" {
		c.HostName = "portfello.app" // set explicit default
	}
	if c.DatabaseDSN == "" {
		return fmt.Errorf("database_dsn is required")
	}

	switch c.Auth.Provider {
	case AuthProviderAuth0:
		if c.Auth.ClientID == "" {
			return fmt.Errorf("auth0 is not configured")
		}
	case AuthProviderLocal:
		if c.Auth.ClientSecret == "" {
			return fmt.Errorf("local auth provider is not configured")
		}
	case AuthProviderMock:
		return nil
	default:
		return fmt.Errorf("invalid auth provider: %s", c.Auth.Provider)
	}

	return nil
}

// NewTestConfig returns configuration suitable for testing.
func NewTestConfig() *Config {
	return &Config{
		DatabaseDSN: "sqlite://:memory:",
		Graph:       GraphQL{},
		Auth: Auth0{
			ClientSecret: "secret-key",
		},
		Logging: Logging{},
	}
}

// InitConfig reads in config file and ENV variables if set. Should be called once at app startup.
func InitConfig(cfgFilePath, envPrefix string) error {
	viper.SetOptions(viper.ExperimentalBindStruct())
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match
	defer viper.WatchConfig()

	if cfgFilePath != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFilePath)
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("fatal error config file @ %s: %w", cfgFilePath, err)
		}

		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		return nil
	}

	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Search config in home directory with name ".portfello" (without extension).
	viper.AddConfigPath(home)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	// Loop through found config files until all are parsed
	for _, configFile := range []string{".portfello", ".portfello-local"} {
		viper.SetConfigName(configFile)
		if err := viper.MergeInConfig(); err == nil {
			_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}

	return nil
}
