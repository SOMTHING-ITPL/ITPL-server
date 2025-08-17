package config

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

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

type RedisDBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

type GoogleConfig struct {
	ClientId     string
	ClientSecret string
	Domain       string
	RedirectURI  string
	ApiHost      string
}
