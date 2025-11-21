package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database Database `yaml:"database"`
	Http     Http     `yaml:"http"`
}

type Database struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	DBName      string `yaml:"db_name"`
	MaxOpenConn int    `yaml:"max_open_conn"`
	MaxIdleConn int    `yaml:"max_idle_conn"`
}

type Http struct {
	Port string `yaml:"port"`
}

func LoadConfig(cfgFile string) (*Config, error) {
	var cfg Config

	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("error read config: %w", err)
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal yaml error: %w", err)
	}

	return &cfg, nil
}
