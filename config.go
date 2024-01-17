package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

var cfg config

type config struct {
	MatrixAPI    string `toml:"matrix_api"`
	InstanceName string `toml:"instance_name"`
	SharedSecret string `toml:"shared_secret"`
	ClientURL    string `toml:"client_url"`
}

func loadConfig() {

	_, err := toml.DecodeFile("config.toml", &cfg)
	if err != nil {
		log.Fatalf("failed to load config: %v\n", err)
	}

}
