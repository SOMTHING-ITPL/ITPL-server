package config

import (
	"github.com/spf13/viper"
)

var (
	KakaoCfg  *KaKaoConfig
	GoogleCfg *GoogleConfig
	DbCfg     *DBConfig
)

func LoadConfigs(configs ...Config) error {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	viper.AutomaticEnv()

	for _, cfg := range configs {
		if err := cfg.Load(); err != nil {
			return err
		}
	}

	return nil
}

func InitConfigs() error {
	KakaoCfg = &KaKaoConfig{}
	DbCfg = &DBConfig{}

	if err := LoadConfigs(KakaoCfg, DbCfg); err != nil {
		return err
	}

	return nil
}
