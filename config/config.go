package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Configurations struct {
	JwtKey           string `default:"my-secret-key" envconfig:"JWT_SECRET_KEY"`
	DbHost           string `default:"localhost" envconfig:"DB_HOST"`
	DbPort           string `default:"3306" envconfig:"DB_PORT"`
	DbUser           string `default:"root" envconfig:"DB_USER"`
	DbPassword       string `default:"root123" envconfig:"DB_PASSWORD"`
	DbName           string `default:"oh_db" envconfig:"DB_NAME"`
	DefaultPageLimit string `default:"10" envconfig:"DEFAULT_PAGE_LIMIT"`
}

var Conf Configurations

func LoadConfig() error {

	errEnv := envconfig.Process("", &Conf)
	if errEnv != nil {
		fmt.Print(errEnv.Error())
	}

	return nil
}
