package config

import "github.com/spf13/viper"

type Config interface {
	Load() error
}

// KaKaoConfig
type KaKaoConfig struct {
	ClientId     string
	ClientSecret string
	Domain       string
	RedirectURI  string
	ApiHost      string
}

func (k *KaKaoConfig) Load() error {
	k.ClientId = viper.GetString("kakao.clientId")
	k.ClientSecret = viper.GetString("kakao.clientSecret")
	k.Domain = viper.GetString("kakao.domain")
	k.RedirectURI = viper.GetString("kakao.redirectURI")
	k.ApiHost = viper.GetString("kakao.apiHost")
	return nil
}

// DBConfig
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func (d *DBConfig) Load() error {
	d.User = viper.GetString("db.user")
	d.Password = viper.GetString("db.password")
	d.Host = viper.GetString("db.host")
	d.Port = viper.GetString("db.port")
	d.Database = viper.GetString("db.database")
	return nil
}
