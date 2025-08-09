package config

import (
	"github.com/spf13/viper"
)

type ServerConfig struct {
	KaKaoConfig *KaKaoConfig
	DBConfig    *DBConfig
}

func LoadConfigs(configs ...Config) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	for _, cfg := range configs {
		if err := cfg.Load(); err != nil {
			return err
		}
	}

	return nil
}

func InitConfigs() (*ServerConfig, error) {
	kakaoCfg := &KaKaoConfig{}
	dbCfg := &DBConfig{}

	if err := LoadConfigs(kakaoCfg, dbCfg); err != nil {
		return nil, err
	}

	return &ServerConfig{
		KaKaoConfig: kakaoCfg,
		DBConfig:    dbCfg,
	}, nil
}
