package config

import (
	"os"

	"github.com/spf13/viper"
)

var Config *ServerConfig

type ServerConfig struct {
	DB    DBConfig
	CORS  CORSConfig
	Auth0 Auth0Config

	Prefix string
}

type DBConfig struct {
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

type Auth0Config struct {
	Audience []string
	Issuer   string
	JWKSURI  string `mapstructure:"jwks_uri"`
	UserInfo string `mapstructure:"user_info"`
	Secret   string
}

func init() {
	configFile := os.Getenv("LATTETALK_CONFIG")
	if configFile == "" {
		panic("please set LATTETALK_CONFIG file to the config file location.")
	}
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	Config = &ServerConfig{
		Prefix: "/",
	}
	if err := viper.Unmarshal(Config); err != nil {
		panic(err)
	}
}
