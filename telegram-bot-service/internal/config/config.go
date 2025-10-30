package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Tg              TgConfig              `yaml:"telegram"`
	ShortenerClient ShortenerClientConfig `yaml:"shortener_client"`
}

type TgConfig struct {
	Token string `yaml:"token"`
}

type ShortenerClientConfig struct {
	Addr string `yaml:"addr"`
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
	cfg.Tg.Token = viper.GetString("telegram.token")
	cfg.ShortenerClient.Addr = viper.GetString("shortener_client.addr")
	return cfg, nil
}
