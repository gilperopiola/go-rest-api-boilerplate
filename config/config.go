package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type ConfigActions interface {
	Setup(env string)
}

type MyConfig struct {
	PORT     string
	DEBUG    bool
	DATABASE struct {
		TYPE     string
		USERNAME string
		PASSWORD string
		HOSTNAME string
		PORT     string
		SCHEMA   string

		CREATE_SCHEMA bool
		PURGE         bool
		CREATE_ADMIN  bool
	}
	JWT struct {
		SECRET           string
		SESSION_DURATION int
	}
}

func (config *MyConfig) Setup(env string) {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.SetConfigName("config_" + env) //config filename without the .JSON or .YAML extension

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("error unmarshaling config: %v", err)
	}

	os.Setenv("PORT", config.PORT)
}
