package parser

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/mwojtyna/swift-api/internal/db"
)

func ParseCsv(r io.Reader) ([]db.Bank, error) {
	reader := csv.NewReader(r)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var banks, hqBanks []db.Bank
	const hqPartLen = 8
	hqBankCodes := make(map[string]struct{}) // Dumb hack because Go doesn't have sets

	if len(records) < 2 || (len(records) > 1 && len(records[0]) != 8) {
		return nil, fmt.Errorf("Invalid CSV")
	}

	for i, record := range records[1:] { // Skip header row
		countryCode := strings.ToUpper(record[0])
		swiftCode := record[1]
		// Skip index 2 (CODE TYPE) - "Redundant columns in the file may be omitted."
		bankName := record[3]
		bankAddress := record[4]
		townName := record[5] // Read town name in case the address is empty
		countryName := strings.ToUpper(record[6])
		// Skip index 7 (TIME ZONE) - "Redundant columns in the file may be omitted."

		if len(countryCode) != 2 {
			return nil, fmt.Errorf(`Invalid row %d with invalid country code "%s" in "%s"`, i, countryCode, record)
		}
		if len(swiftCode) != 11 {
			return nil, fmt.Errorf(`Invalid row %d with invalid SWIFT code "%s" in "%s"`, i, swiftCode, record)
		}

		// EDGE CASE: Set address to "town_name" if it is empty
		var address string
		if strings.TrimSpace(bankAddress) == "" {
			address = townName
		} else {
			address = bankAddress
		}
		bank := db.Bank{
			SwiftCode:       swiftCode,
			HqSwiftCode:     sql.NullString{},
			BankName:        bankName,
			Address:         address,
			CountryISO2Code: countryCode,
			CountryName:     countryName,
		}

		// If swift code doesn't end with XXX, then the first 8 characters are the swift code for this bank's HQ (plus XXX)
		// We assume that this HQ exists, later we remove ones that don't (we use a set to keep track of HQs that exist)
		// I think this is the best way to minimize iterating through the banks
		isHq, hqCode := IsSwiftCodeHq(swiftCode)
		if !isHq {
			bank.HqSwiftCode = sql.NullString{
				String: hqCode,
				Valid:  true,
			}
			banks = append(banks, bank)
		} else {
			// If it *does* end with XXX, then it is an HQ code and this bank is the HQ
			hqBankCodes[swiftCode] = struct{}{} // Add to set
			hqBanks = append(hqBanks, bank)
		}
	}

	// EDGE CASE: Remove bank's hq_code if HQ doesn't exist (e.g. ALBPPLP1XXX doesn't exist, but ALBPPLP1BMW does)
	for i := range banks {
		b := &banks[i]
		if !b.HqSwiftCode.Valid {
			continue
		}

		_, ok := hqBankCodes[b.HqSwiftCode.String]
		if !ok {
			b.HqSwiftCode = sql.NullString{}
		}
	}

	// Make HQ banks appear first in array to prevent foreign key errors
	sortedBanks := append(hqBanks, banks...)

	return sortedBanks, nil
}

const hqPartLen = 8

// Returns whether the bank is the headquarters, if not - returns the bank's headquarters code assuming they exist.
// Assumes the given swift code is valid.
func IsSwiftCodeHq(code string) (bool, string) {
	if code[hqPartLen:] == "XXX" {
		return true, ""
	} else {
		return false, code[:hqPartLen] + "XXX"
	}
}
