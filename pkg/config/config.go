package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DatabaseDSN       string
	GraphqlPort       string
	GraphqlPlayground bool

	Auth Auth0
}

type Auth0 struct {
	Domain       string
	ClientID     string
	ClientSecret string
	Audience     string
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
