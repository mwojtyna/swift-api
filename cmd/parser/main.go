package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"

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

	_, _, err = parseCSV(csv_name)
	if err != nil {
		logger.Fatalf(`Error parsing file and inserting data to db '%s': "%s"`, csv_name, err.Error())
	}
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
	for i, record := range records[1:] { // Skip header row
		if len(record) != 8 {
			logger.Printf(`Skipping row %d with unexpected length "%s"`, i, record)
			continue
		}

		countryCode := strings.ToUpper(record[0])
		swiftCode := record[1]
		// Skip index 2 (CODE TYPE) - "Redundant columns in the file may be omitted."
		bankName := record[3]
		bankAddress := record[4]
		// Skip index 5 (TOWN NAME) - "Redundant columns in the file may be omitted."
		countryName := strings.ToUpper(record[6])
		countryTimeZone := record[7]

		if len(countryCode) != 2 {
			logger.Printf(`Skipping row %d with invalid country code "%s" in "%s"`, i, countryCode, record)
			continue
		}
		if len(swiftCode) != 11 {
			logger.Printf(`Skipping row %d with invalid SWIFT code "%s" in "%s"`, i, swiftCode, record)
			continue
		}

		hqSwiftCode := ""
		const hqPartLen = 8
		if swiftCode[hqPartLen:] != "XXX" {
			hqSwiftCode = swiftCode[:hqPartLen] + "XXX"
		}

		bank := db.Bank{
			SwiftCode:       swiftCode,
			HqSwiftCode:     hqSwiftCode,
			CountryISO2Code: countryCode,
			BankName:        bankName,
			Address:         bankAddress,
		}
		// Countries will be deduplicated while inserting into db
		country := db.Country{
			ISO2Code:    countryCode,
			CountryName: countryName,
			TimeZone:    countryTimeZone,
		}

		banks = append(banks, bank)
		countries = append(countries, country)
	}

	return banks, countries, nil
}
