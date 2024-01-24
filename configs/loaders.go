package configs

import (
	"github.com/spf13/viper"
)

func LoadConfig() (conf *Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return &Config{
		Application: Application{
			Name:      viper.GetString("APP_NAME"),
			SecretKey: viper.GetString("APP_SECRET_KEY"),
		},
		Server: Server{
			Host:         viper.GetString("APP_HOST"),
			Port:         viper.GetInt("APP_PORT"),
			IsProduction: viper.GetBool("APP_PRODUCTION"),
			IsLogging:    viper.GetBool("APP_LOGGING"),
		},
		Database: Database{
			Host:     viper.GetString("DB_HOST"),
			Name:     viper.GetString("DB_NAME"),
			Username: viper.GetString("DB_USERNAME"),
			Password: viper.GetString("DB_PASSWORD"),
			Port:     viper.GetInt("DB_PORT"),
		},
		Redis: Redis{
			Host:          viper.GetString("REDIS_HOST"),
			Password:      viper.GetString("REDIS_PASSWORD"),
			Database:      viper.GetInt("REDIS_DATABASE"),
			DefaultExpire: viper.GetInt("REDIS_DEFAULT_EXPIRE"),
		},
	}, nil
}
