package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Pgdb_config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

var instance *Pgdb_config
var once sync.Once

func GetDbConfig() *Pgdb_config {
	once.Do(func() {
		instance = &Pgdb_config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			panic("Error")
		}
	})
	return instance
}
