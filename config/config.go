package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DB_USER string
	DB_PASS string
	DB_NAME string
}

func LoadEnv() (Env, error) {
	err := godotenv.Load()
	if err != nil {
		return Env{}, err
	}

	config := Env{
		DB_USER: os.Getenv("DB_USER"),
		DB_PASS: os.Getenv("DB_PASS"),
		DB_NAME: os.Getenv("DB_NAME"),
	}

	return config, nil
}
