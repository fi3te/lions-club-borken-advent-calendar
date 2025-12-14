package config

import (
	"log"
	"os"

	"github.com/fi3te/lions-club-borken-advent-calendar/pkg/common"
	"github.com/joho/godotenv"
)

func ReadConfig() (*Config, error) {
	config, err := common.ReadYaml[Config](FileName)
	if err != nil {
		return nil, err
	}

	err = godotenv.Load()
	if err != nil {
		log.Printf("Loading .env file failed: %v\n", err)
	}

	config.Email.Smtp = readSmtpConfig()

	err = config.validate()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func readSmtpConfig() SmtpConfig {
	return SmtpConfig{
		Host:     os.Getenv(EnvVarSmtpHost),
		Username: os.Getenv(EnvVarSmtpUsername),
		Password: os.Getenv(EnvVarSmtpPassword),
	}
}
