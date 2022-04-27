package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Dialect  string `toml:"Dialect"`
	Host     string `toml:"Host"`
	Port     int    `toml:"Port"`
	Username string `toml:"Username"`
	Password string `toml:"Password"`
	Name     string `toml:"Name"`
	Charset  string `toml:"Charset"`
}

func ReadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return &Config{DB: &DBConfig{}}, err
	}
	var conf DBConfig
	err = toml.Unmarshal(data, &conf)
	if err != nil {
		return &Config{DB: &DBConfig{}}, err
	}
	return &Config{DB: &conf}, nil
}
