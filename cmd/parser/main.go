package main

import (
	"log"
	"os"

	"github.com/mwojtyna/swift-api/config"
	"github.com/mwojtyna/swift-api/internal/db"
	"github.com/mwojtyna/swift-api/internal/parser"
)

var logger = log.New(os.Stderr, "[CSV PARSER] ", log.LstdFlags|log.Lshortfile)

func main() {
	if len(os.Args) != 2 {
		logger.Fatal("Error: CSV filename not specified")
	}
	csvName := os.Args[1]

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

	empty, err := db.IsEmpty(pg)
	if err != nil {
		logger.Fatalf(`ERROR checking if DB is empty: "%s"`, err.Error())
	}
	if !empty {
		logger.Fatalln("Error: DB isn't empty")
	}

	file, err := os.Open(csvName)
	if err != nil {
		logger.Fatalf("ERROR opening csv file '%s'", csvName)
	}
	defer file.Close()

	banks, err := parser.ParseCsv(file)
	if err != nil {
		logger.Fatalf(`ERROR parsing file and inserting data to db '%s': "%s"`, csvName, err.Error())
	}
	logger.Printf("Parsed %d banks", len(banks))

	err = db.InsertBanks(pg, banks)
	if err != nil {
		logger.Fatalf(`ERROR inserting banks: "%s"`, err.Error())
	}
	logger.Println("Inserted banks")

	logger.Println("Done!")
}
