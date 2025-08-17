package config

import (
	"github.com/spf13/viper"
)

var (
	KakaoCfg  *KaKaoConfig
	GoogleCfg *GoogleConfig
	DbCfg     *DBConfig
	RedisCfg  *RedisDBConfig
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
	GoogleCfg = &GoogleConfig{}
	RedisCfg = &RedisDBConfig{}

	if err := LoadConfigs(KakaoCfg, DbCfg, GoogleCfg, RedisCfg); err != nil {
		return err
	}

	return nil
}

func (k *KaKaoConfig) Load() error {
	//In Yaml
	k.Domain = viper.GetString("kakao.domain")
	k.RedirectURI = viper.GetString("kakao.redirectURI")
	k.ApiHost = viper.GetString("kakao.apiHost")

	//In Env
	if val := viper.GetString("KAKAO_CLIENT_ID"); val != "" {
		k.ClientId = val
	}
	if val := viper.GetString("KAKAO_CLIENT_SECRET"); val != "" {
		k.ClientSecret = val
	}
	return nil
}

func (d *DBConfig) Load() error {
	//In yaml
	d.Host = viper.GetString("db.host")
	d.Port = viper.GetString("db.port")
	d.Database = viper.GetString("db.database")

	//In Env
	if val := viper.GetString("DB_USER"); val != "" {
		d.User = val
	}
	if val := viper.GetString("DB_PASSWORD"); val != "" {
		d.Password = val
	}

	return nil
}

func (d *RedisDBConfig) Load() error {
	//In yaml
	d.Host = viper.GetString("rdb.host")
	d.Port = viper.GetString("rdb.port")
	d.Database = viper.GetString("rdb.database")

	//In Env
	if val := viper.GetString("REDIS_DB_PASSWORD"); val != "" {
		d.Password = val
	}

	return nil
}

func (g *GoogleConfig) Load() error {
	//In Yaml
	g.Domain = viper.GetString("google.domain")
	g.RedirectURI = viper.GetString("google.redirectURI")
	g.ApiHost = viper.GetString("google.apiHost")

	//In Env
	if val := viper.GetString("GOOGLE_CLIENT_ID"); val != "" {
		g.ClientId = val
	}
	if val := viper.GetString("GOOGLE_CLIENT_SECRET"); val != "" {
		g.ClientSecret = val
	}
	return nil
}
