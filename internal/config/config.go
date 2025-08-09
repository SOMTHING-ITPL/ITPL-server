package config

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func GetDBConfig() DBConfig {
	return DBConfig{
		User:     "root",
		Password: "534f4b4354", //check this field
		Host:     "localhost",
		Port:     "3306",
		Database: "itpl_db",
	}
}
