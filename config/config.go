package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

var Config configuration

type configuration struct {
	MatrixAPI    string `toml:"matrix_api"`
	InstanceName string `toml:"instance_name"`
	SharedSecret string `toml:"shared_secret"`
	ClientURL    string `toml:"client_url"`

	Debug bool `toml:"debug"`
}

func LoadConfig() {

	_, err := toml.DecodeFile("config.toml", &Config)
	if err != nil {
		log.Fatalf("failed to load config: %v\n", err)
	}

}
