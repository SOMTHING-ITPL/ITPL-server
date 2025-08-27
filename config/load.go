package config

import (
	"github.com/spf13/viper"
)

// TODO : prod 환겨엥서는 환경변수 없으면 걍  panic

var (
	KakaoCfg  *KaKaoConfig
	GoogleCfg *GoogleConfig
	DbCfg     *DBConfig
	RedisCfg  *RedisDBConfig
	SmtpCfg   *GoogleSMTPConfig
	KopisCfg  *KopisApiConfig
	OpenAICfg *OpenAiConfig
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
	SmtpCfg = &GoogleSMTPConfig{}
	KopisCfg = &KopisApiConfig{}
	OpenAICfg = &OpenAiConfig{}

	if err := LoadConfigs(KakaoCfg, DbCfg, GoogleCfg, RedisCfg, SmtpCfg, KopisCfg, OpenAICfg); err != nil {
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
	d.Addr = viper.GetString("rdb.addr")
	d.Database = viper.GetInt("rdb.database")

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

func (gsmtp *GoogleSMTPConfig) Load() error {
	//In yaml
	gsmtp.HostServer = viper.GetString("google_smtp.smtpHost")
	gsmtp.Port = viper.GetString("google_smtp.smtpPort")

	//In Env
	if val := viper.GetString("GOOGLE_SMTP_HOST"); val != "" {
		gsmtp.From = val
	}
	if val := viper.GetString("GOOGLE_SMTP_APP_PASSWORD"); val != "" {
		gsmtp.AppPassword = val
	}
	return nil
}

func (k *KopisApiConfig) Load() error {
	//In Env
	if val := viper.GetString("KOPIS_API_SECRET_KEY"); val != "" {
		k.SecretKey = val
	}

	return nil
}

func (o *OpenAiConfig) Load() error {
	//In Env
	if val := viper.GetString("OPEN_AI_SECRET_KEY"); val != "" {
		o.SecretKey = val
	}

	return nil
}

func (s *S3Config) Load() {
	//In Env
	if val := viper.GetString("AWS_ACCESS_KEY_ID"); val != "" {
		s.AccessKey = val
	}
	if val := viper.GetString("AWS_SECRET_ACCESS_KEY"); val != "" {
		s.SecretKey = val
	}
	if val := viper.GetString("BUCKET_NAME"); val != "" {
		s.BucketName = val
	}
}
