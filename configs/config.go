package configs

type (
	Application struct {
		Name      string `mapstructure:"APP_NAME"`
		SecretKey string `mapstructure:"SECRET_KEY"`
	}
	Server struct {
		Host         string `mapstructure:"APP_HOST"`
		Port         int    `mapstructure:"APP_PORT"`
		IsProduction bool   `mapstructure:"IS_PRODUCTION"`
		IsLogging    bool   `mapstructure:"IS_LOGGING"`
	}
	Database struct {
		Host     string `mapstructure:"DB_HOST"`
		Name     string `mapstructure:"DB_NAME"`
		Username string `mapstructure:"DB_USERNAME"`
		Password string `mapstructure:"DB_PASSWORD"`
		Port     int    `mapstructure:"DB_PORT"`
	}
	Redis struct {
		Host          string `mapstructure:"REDIS_HOST"`
		Database      int    `mapstructure:"REDIS_DB"`
		Password      string `mapstructure:"REDIS_PASSWORD"`
		DefaultExpire int    `mapstructure:"REDIS_DEFAULT_EXPIRE"`
	}
	Config struct {
		Application
		Server
		Database
		Redis
	}
)
