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
	Password string
	Addr     string
	Database int
}

type GoogleConfig struct {
	ClientId     string
	ClientSecret string
	Domain       string
	RedirectURI  string
	ApiHost      string
}

type GoogleSMTPConfig struct {
	HostServer  string
	Port        string
	From        string
	AppPassword string
}

type KopisApiConfig struct {
	SecretKey string
}

type OpenAiConfig struct {
	SecretKey string
}
