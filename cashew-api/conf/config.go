package conf

import (
	"os"

	"gopkg.in/yaml.v3"
)

type OAuthConfig struct {
	ClientKey    string
	ClientSecret string
}

func LoadOAuthConfig() *OAuthConfig {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	var config OAuthConfig

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return &config
}
