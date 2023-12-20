package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	EnableService bool `yaml:"enableService"`
	Database      struct {
		Psql struct {
			Url string `yaml:"url"`
		} `yaml:"psql"`
		Redis struct {
			Addr     string `yaml:"addr"`
			Password string `yaml:"password"`
			Db       int    `yaml:"db"`
		} `yaml:"redis"`
	} `yaml:"database"`
}

func Init(confPath string) Config {
	f, err := os.Open(confPath)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		processError(err)
	}

	return cfg
}

func processError(err error) {
	fmt.Println(err.Error())
	os.Exit(0)
}
