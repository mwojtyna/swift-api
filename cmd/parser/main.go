package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/mwojtyna/swift-api/config"
	"github.com/mwojtyna/swift-api/internal/db"
)

var logger = log.New(os.Stderr, "[CSV PARSER] ", log.LstdFlags|log.Lshortfile)

func main() {
	if len(os.Args) != 2 {
		logger.Fatal("Error: CSV filename not specified")
	}
	csv_name := os.Args[1]

	env, err := config.LoadEnv()
	if err != nil {
		logger.Fatalf(`Error reading config: "%s"`, err.Error())
	}
	logger.Println("Read envs")

	_, err = db.Connect(env.DB_USER, env.DB_PASS, env.DB_NAME)
	if err != nil {
		logger.Fatalf(`Error connecting to db: "%s"`, err.Error())
	}
	logger.Println("Connected to db")

	banks, countries, err := parseCSV(csv_name)
	if err != nil {
		logger.Fatalf(`Error parsing file and inserting data to db '%s': "%s"`, csv_name, err.Error())
	}
	logger.Println(banks)
	logger.Println(countries)
}

func parseCSV(filename string) ([]db.Bank, []db.Country, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	var banks []db.Bank
	var countries []db.Country
	for _, row := range records[1:] { // Skip header row
		var bank db.Bank
		var country db.Country

		for i, cell := range row {
			switch i {
			case 0:
				bank.CountryISO2Code = cell
				country.ISO2Code = cell
			case 1:
				bank.SwiftCode = cell
			// Skip index 2 (CODE TYPE) - "Redundant columns in the file may be omitted."
			case 3:
				bank.BankName = cell
			case 4:
				bank.Address = cell
			// Skip index 5 (TOWN NAME) - "Redundant columns in the file may be omitted."
			case 6:
				country.CountryName = cell
			case 7:
				country.TimeZone = cell
			default:
				logger.Printf(`Unexpected cell "%s" with index=%d in row "%s"`, cell, i, row)
			}

			banks = append(banks, bank)
			countries = append(countries, country)
		}
	}

	return banks, countries, nil
}
