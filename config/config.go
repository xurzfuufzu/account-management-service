package config

import (
	"account-management/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Server Server `env-required:"true" yaml:"server"`
	//Log    Log    `env-required:"true" yaml:"log"`
	DB  DB  `env-required:"true" yaml:"db"`
	JWT JWT `env-required:"true" yaml:"jwt"`
}

type Server struct {
	Host string `env-required:"true" yaml:"host"`
	Port string `env-required:"true" yaml:"port"`
}

//type Log struct {
//	Level string `env-required:"true" yaml:"level"`
//}

type DB struct {
	Host     string `env-required:"true" yaml:"host"`
	Port     string `env-required:"true" yaml:"port"`
	Username string `env-required:"true" yaml:"username"`
	Password string `env-required:"true" yaml:"password"`
	Database string `env-required:"true" yaml:"database"`
	URL      string `yaml:"-"`
}

type JWT struct {
	SecretKey string        `env-required:"true" yaml:"sign_key"`
	TokenTTL  time.Duration `env-required:"true" yaml:"token_ttl"`
}

var instance *Config

func NewConfig() *Config {
	logger := logging.GetLogger()
	logger.Info("read application config")

	instance := &Config{}
	if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
		help, _ := cleanenv.GetDescription(instance, nil)
		logger.Info(help)
		logger.Fatal(err)
	}

	return instance
}
