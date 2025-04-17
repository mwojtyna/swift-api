package config

import (
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type envType string

const (
	SwiftApiEnvDevelopment envType = "development"
	SwiftApiEnvTesting     envType = "testing"
	SwiftApiEnvProduction  envType = "production"
)

type Env struct {
	DB_USER      string  `validate:"required"`
	DB_PASS      string  `validate:"required"`
	DB_NAME      string  `validate:"required"`
	API_PORT     string  `validate:"required"`
	SWIFTAPI_ENV envType `validate:"required"`
}

func LoadEnv() (Env, error) {
	env := envType(os.Getenv("SWIFTAPI_ENV"))
	if env == "" {
		env = SwiftApiEnvDevelopment
	}
	if testing.Testing() {
		env = SwiftApiEnvTesting
	}

	var err error
	switch env {
	case SwiftApiEnvProduction:
		err = godotenv.Load()
	case SwiftApiEnvDevelopment:
		err = godotenv.Load(".env.development.local")
	}
	if err != nil {
		return Env{}, err
	}

	config := Env{
		DB_USER:      os.Getenv("DB_USER"),
		DB_PASS:      os.Getenv("DB_PASS"),
		DB_NAME:      os.Getenv("DB_NAME"),
		API_PORT:     os.Getenv("API_PORT"),
		SWIFTAPI_ENV: env,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(config)
	if err != nil {
		return Env{}, err
	}

	return config, nil
}
