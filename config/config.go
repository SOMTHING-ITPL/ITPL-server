package config

import (
	"github.com/spf13/viper"
)

type Config interface {
	Load() error
}

type KaKaoConfig struct {
	ClientId     string
	ClientSecret string
	Domain       string
	RedirectURI  string
	ApiHost      string
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

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
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

type GoogleConfig struct {
	ClientId     string
	ClientSecret string
	Domain       string
	RedirectURI  string
	ApiHost      string
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
