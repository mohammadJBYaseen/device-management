package config

import (
	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
)

type (
	AppConfig struct {
		Server ServerConfig
	}
	ServerConfig struct {
		AppName string `env:"APP_NAME"`
		Port    int    `env:"SERVER_PORT"`
	}
)

// SetUp Post config initialization
func (appConf *AppConfig) SetUp() error {
	return nil
}

// Load config initialization
func Load() AppConfig {
	cfg := AppConfig{}
	//TODO make it more dynamic
	envFeeder := feeder.DotEnv{Path: ".env"}
	jsonFeeder := feeder.Json{Path: "config/config.json"}
	if err := config.New().AddFeeder(jsonFeeder).AddFeeder(envFeeder).AddStruct(&cfg).Feed(); err != nil {
		panic(err)
	}
	return cfg
}
