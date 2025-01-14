package cidbot

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Token         string `yaml:"token"`
	AdminServerID string `yaml:"admin_server_id"`
}

var configPath = flag.String("config-path", "./config.yml", "The path to a file containing the bot's configuration")

var Configuration *Config

func ParseConfig() (err error) {
	flag.Parse()

	fileContents, err := os.ReadFile(*configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(fileContents, &Configuration)
	if err != nil {
		return err
	}
	return nil
}
