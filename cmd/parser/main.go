package main

import (
	"database/sql"
	"encoding/csv"
	"log"
	"maps"
	"os"
	"slices"
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
		logger.Fatalf(`Error reading envs: "%s"`, err.Error())
	}
	logger.Println("Read envs")

	pg, err := db.Connect(env.DB_USER, env.DB_PASS, env.DB_NAME)
	if err != nil {
		logger.Fatalf(`Error connecting to db: "%s"`, err.Error())
	}
	logger.Println("Connected to db")

	banks, countries, err := parseCSV(csv_name)
	if err != nil {
		logger.Fatalf(`Error parsing file and inserting data to db '%s': "%s"`, csv_name, err.Error())
	}
	logger.Printf("Parsed %d banks, %d countries", len(banks), len(countries))

	err = db.InsertCountries(pg, countries)
	if err != nil {
		logger.Fatalf(`Error inserting countries: "%s"`, err.Error())
	}
	logger.Println("Inserted countries")

	err = db.InsertBanks(pg, banks)
	if err != nil {
		logger.Fatalf(`Error inserting banks: "%s"`, err.Error())
	}
	logger.Println("Inserted banks")
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

	var banks, hqBanks []db.Bank
	const hqPartLen = 8
	hqBankCodes := make(map[string]struct{}) // Dumb hack because Go doesn't have sets
	countriesMap := make(map[db.Country]struct{})

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

		bank := db.Bank{
			SwiftCode:       swiftCode,
			HqSwiftCode:     sql.NullString{},
			CountryISO2Code: countryCode,
			BankName:        bankName,
			Address:         bankAddress,
		}
		country := db.Country{
			ISO2Code:    countryCode,
			CountryName: countryName,
			TimeZone:    countryTimeZone,
		}

		// If swift code doesn't end with XXX, then the first 8 characters are the swift code for this bank's HQ (plus XXX)
		// We assume that this HQ exists, later we remove ones that don't (we use a set to keep track of HQs that exist)
		// I think this is the best way to minimize iterating through the banks
		const hqPartLen = 8
		if swiftCode[hqPartLen:] != "XXX" {
			bank.HqSwiftCode = sql.NullString{
				String: swiftCode[:hqPartLen] + "XXX",
				Valid:  true,
			}
			banks = append(banks, bank)
		} else {
			// If it *does* end with XXX, then it is an HQ code and this bank is the HQ
			hqBankCodes[swiftCode] = struct{}{} // Add to set
			hqBanks = append(hqBanks, bank)
		}

		countriesMap[country] = struct{}{} // Add to set
	}

	// EDGE CASE: Remove bank's hq_code if HQ doesn't exist (e.g. ALBPPLP1XXX doesn't exist, but ALBPPLP1BMW does)
	for i := range banks {
		b := &banks[i]
		if !b.HqSwiftCode.Valid {
			continue
		}

		_, ok := hqBankCodes[b.HqSwiftCode.String]
		if !ok {
			// logger.Printf("Removed hq code %s from bank %+v", b.HqSwiftCode.String, b)
			b.HqSwiftCode = sql.NullString{}
		}
	}

	// Make HQ banks appear first in array to prevent foreign key errors
	sortedBanks := append(hqBanks, banks...)

	// Convert back to array
	countries := slices.Collect(maps.Keys(countriesMap))

	return sortedBanks, countries, nil
}
