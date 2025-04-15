package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DB_USER  string
	DB_PASS  string
	DB_NAME  string
	API_PORT string
}

func LoadEnv() (Env, error) {
	err := godotenv.Load()
	if err != nil {
		return Env{}, err
	}

	config := Env{
		DB_USER:  os.Getenv("DB_USER"),
		DB_PASS:  os.Getenv("DB_PASS"),
		DB_NAME:  os.Getenv("DB_NAME"),
		API_PORT: os.Getenv("API_PORT"),
	}

	if config.DB_USER == "" {
		return Env{}, fmt.Errorf("DB_USER env var not specified")
	}
	if config.DB_PASS == "" {
		return Env{}, fmt.Errorf("DB_PASS env var not specified")
	}
	if config.DB_NAME == "" {
		return Env{}, fmt.Errorf("DB_NAME env var not specified")
	}
	if config.API_PORT == "" {
		return Env{}, fmt.Errorf("API_PORT env var not specified")
	}

	return config, nil
}
