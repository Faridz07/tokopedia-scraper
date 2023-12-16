package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App            AppConfig            `json:"app"`
	Infrastructure InfrastructureConfig `json:"infrastructure"`
	Scrape         ScrapeConfig         `json:"scrape"`
}

type AppConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Env     string `json:"env"`
}

type InfrastructureConfig struct {
	Database DatabaseConfig `json:"database"`
	Log      LogConfig      `json:"log"`
}

type LogConfig struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Stdout   bool   `json:"stdout"`
}

type DatabaseConfig struct {
	Dialect       string `json:"dialect"`
	Host          string `json:"host"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	DBName        string `json:"dbName"`
	MigrationPath string `json:"migrationPath"`
}

type ScrapeConfig struct {
	Worker int          `json:"worker"`
	Web    WebConfig    `json:"web"`
	Output OutputConfig `json:"output"`
}

type WebConfig struct {
	BaseUrl     string `json:"baseUrl"`
	PathToScrap string `json:"pathToScrap"`
}

type OutputConfig struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
}

func New() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./resource")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("fatal error config file: config file not found")
		}
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
