package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env        string `envconfig:"ENV"`
	Delay      int    `envconfig:"DELAY"`
	KeyVtex    string `envconfig:"KEY_VTEX"`
	TokenVtex  string `envconfig:"TOKEN_VTEX"`
	Account    string `envconfig:"ACCOUNT_VTEX"`
	SkuId      string `envconfig:"SKUID"`
	HostDb     string `envconfig:"HOSTDB"`
	PortDb     int    `envconfig:"PORTDB"`
	UserDb     string `envconfig:"USERDB"`
	PasswordDb string `envconfig:"PASSWORDDB"`
	DbName     string `envconfig:"DBNAME"`
	TokenGer   string `envconfig:"TOKEN_GER"`
}

var Env Config

func New(cfg interface{}) (err error) {
	err = godotenv.Load()
	if err != nil {
		return
	}
	err = envconfig.Process("", cfg)
	if err != nil {
		return
	}
	return
}
