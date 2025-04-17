package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mwojtyna/swift-api/config"
	"github.com/mwojtyna/swift-api/internal/api"
	"github.com/mwojtyna/swift-api/internal/db"
)

var logger = log.New(os.Stderr, "[API] ", log.Ldate|log.Ltime)

func main() {
	env, err := config.LoadEnv()
	if err != nil {
		logger.Fatalf(`ERROR reading envs: "%s"`, err.Error())
	}
	logger.Println("Read envs")

	pg, err := db.Connect(env.DB_USER, env.DB_PASS, env.DB_NAME, env.DB_PORT)
	if err != nil {
		logger.Fatalf(`ERROR connecting to db: "%s"`, err.Error())
	}
	logger.Println("Connected to db")

	addr := fmt.Sprintf(":%s", env.API_PORT)
	server := api.NewApiServer(addr, pg, logger)

	logger.Printf("Server running on %s", addr)
	server.Run()
}
