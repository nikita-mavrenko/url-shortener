package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	GRPC  *gRPCConfig  `yaml:"grpc"`
	Db    *DbConfig    `yaml:"db"`
	Http  *HttpConfig  `yaml:"http"`
	Redis *RedisConfig `yaml:"redis"`
}

type gRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type DbConfig struct {
	Url string `yaml:"url"`
}

type HttpConfig struct {
	Address string `yaml:"address"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

func Load() (*Config, error) {
	const op = "config.Load"
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	var cfg *Config
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}
	return cfg, nil
}
