package config

import (
	"log"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
)

const FileName = "config.yml"
const EnvVarSmtpHost = "SMTP_HOST"
const EnvVarSmtpUsername = "SMTP_USERNAME"
const EnvVarSmtpPassword = "SMTP_PASSWORD"

func ReadConfig() (*Config, error) {
	yml, err := os.ReadFile(FileName)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(yml, &config); err != nil {
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

	return &config, nil
}

func readSmtpConfig() SmtpConfig {
	return SmtpConfig{
		Host:     os.Getenv(EnvVarSmtpHost),
		Username: os.Getenv(EnvVarSmtpUsername),
		Password: os.Getenv(EnvVarSmtpPassword),
	}
}
