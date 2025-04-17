package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	DB_USER         string  `validate:"required"`
	DB_PASS         string  `validate:"required"`
	DB_NAME         string  `validate:"required"`
	DB_PORT         string  `validate:"required"`
	API_PORT        string  `validate:"required"`
	SWIFTAPI_ENV    envType `validate:"required"`
	ProjectRootPath string  `validate:"required"`
}

func LoadEnv() (Env, error) {
	env := envType(os.Getenv("SWIFTAPI_ENV"))
	if env == "" {
		env = SwiftApiEnvDevelopment
	}
	if testing.Testing() {
		env = SwiftApiEnvTesting
	}

	root := findProjectRoot()

	var err error
	if env == SwiftApiEnvProduction {
		err = godotenv.Load(filepath.Join(root, ".env"))
	} else {
		err = godotenv.Load(filepath.Join(root, fmt.Sprintf(".env.%s.local", env)))
	}
	if err != nil {
		return Env{}, err
	}

	config := Env{
		DB_USER:         os.Getenv("DB_USER"),
		DB_PASS:         os.Getenv("DB_PASS"),
		DB_NAME:         os.Getenv("DB_NAME"),
		DB_PORT:         os.Getenv("DB_PORT"),
		API_PORT:        os.Getenv("API_PORT"),
		SWIFTAPI_ENV:    env,
		ProjectRootPath: root,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(config)
	if err != nil {
		return Env{}, err
	}

	return config, nil
}

const projectDirName = "swift-api"

func findProjectRoot() string {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := string(re.Find([]byte(cwd)))
	return rootPath
}
