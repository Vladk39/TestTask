package config

import (
	"github.com/caarlos0/env/v11"
)

type DBConfig struct {
	DBconn string `env:"DBCON"`
}

type URLApiForReq struct {
	GetAgeURL      string `env:"getage"`
	GetGenderURL   string `env:"getgender"`
	GetNationalURL string `env:"getnational"`
}
type Config struct {
	DBConfig
	URLApiForReq
	ServerConfig
}

type ServerConfig struct {
	Port string `env:"serverport"`
}

func (c *Config) GetURLApiForReqConfig() *URLApiForReq {
	return &c.URLApiForReq
}

func (c *Config) GetDBConfig() *DBConfig {
	return &c.DBConfig
}

func GetConfig() (*Config, error) {
	conf := &Config{}

	err := env.Parse(conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
