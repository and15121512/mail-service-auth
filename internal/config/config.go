package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Auth    struct {
		Login        string `yaml:"login"`
		PasswordHash string `yaml:"password_hash"`
		Salt         string `yaml:"salt"`
		Secret       string `yaml:"secret"`
	}
	Listen struct {
		Port string `yaml:"port"`
	}
}

var instance *Config
var once sync.Once

func GetConfig(logger *zap.SugaredLogger) *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Errorf(help) // TODO: Use normal logger here !!!
		}
	})
	return instance
}
