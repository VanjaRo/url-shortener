package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
	Options struct {
		Schema string `json:"schema"`
		Prefix string `json:"prefix"`
	} `json:"options"`
	Redis struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"redis"`
}

func FromFile(path string) (*Config, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(f, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
