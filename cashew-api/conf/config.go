package conf

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type OAuthConfig struct {
	ClientKey    string
	ClientSecret string
}

func LoadOAuthConfig() *OAuthConfig {
	cwd, _ := os.Getwd()
	data, err := os.ReadFile("conf/config.yaml")
	if err != nil {
		log.Printf("failing here")
		log.Printf("current working directory:%v", cwd)
		panic(err)
	}

	var config OAuthConfig

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return &config
}
