package initializers

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost          string `mapstructure:"POSTGRES_HOST"`
	DBUserName      string `mapstructure:"POSTGRES_USER"`
	DBUserPassoword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName          string `mapstructure:"POSTGRES_DB"`
	DBPort          string `mapstructure:"POSTGRES_PORT"`

	ServerPort   string `mapstructure:"SERVER_PORT"`
	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`

	JwtSecret    string        `mapstructure:"JWT_SECRET"`
	JwtExpiredIn time.Duration `mapstructure:"JWT_EXPIRED_IN"`
	JwtMaxAge    int           `mapstructure:"JWT_MAXAGE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("fatal error config file: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return
}

/* DOCUMENTACIÓN DE VIPER. VER:
https://github.com/spf13/viper
*/
