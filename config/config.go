package config

import (
	"flag"
	"os"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
)

type Config struct {
	IsDevelopment bool   `yaml:"is_development"`
	Token         string `yaml:"token"`
	AdminServerID string `yaml:"admin_server_id"`
	WhitelistPath string `yaml:"whitelist_path"`
}

var configPath = flag.String("config-path", "", "The path to a file containing the bot's configuration. Defaults to $XDG_CONFIG_HOME/CIDBot/config.yml")

func Parse() (*Config, error) {
	flag.Parse()

	if *configPath == "" {
		cfgPath, err := xdg.ConfigFile("CIDBot/config.yml")
		if err != nil {
			return nil, err
		}

		*configPath = cfgPath
	}

	fileContents, err := os.ReadFile(*configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(fileContents, &config)
	if err != nil {
		return nil, err
	}

	return &config, err
}
