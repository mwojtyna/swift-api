package db

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

func GetBank(db *sqlx.DB, swiftCode string) (Bank, error) {
	var bank Bank

	err := db.Get(&bank, "SELECT * FROM bank WHERE swift_code=$1;", swiftCode)
	if err != nil {
		return Bank{}, err
	}

	return bank, nil
}

func GetBankBranches(db *sqlx.DB, swiftCode string) ([]Bank, error) {
	var branches []Bank

	err := db.Select(&branches, `
		SELECT b2.* FROM bank AS b1 
		JOIN bank AS b2 ON b2.hq_swift_code=b1.swift_code
		WHERE b1.swift_code=$1;
		`, swiftCode)
	if err != nil {
		return nil, err
	}

	return branches, nil
}

func GetBanksInCountry(db *sqlx.DB, countryCode string) ([]Bank, error) {
	var banks []Bank

	err := db.Select(&banks, "SELECT * FROM bank WHERE country_iso2_code=$1;", countryCode)
	if err != nil {
		return nil, err
	}

	return banks, nil
}

func CheckBankHqExists(db *sqlx.DB, hqSwiftCode string) (bool, error) {
	_, err := GetBank(db, hqSwiftCode)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		// Error other than ErrNoRows occurred
		return false, err
	} else if errors.Is(err, sql.ErrNoRows) {
		// Bank's HQ not found
		return false, nil
	} else {
		// Bank's HQ found
		return true, nil
	}
}

func InsertBank(db *sqlx.DB, bank Bank) error {
	return InsertBanks(db, []Bank{bank})
}

func InsertBanks(db *sqlx.DB, banks []Bank) error {
	_, err := db.NamedExec(`INSERT INTO bank (swift_code, hq_swift_code, bank_name, address, country_iso2_code, country_name) 
		VALUES (:swift_code, :hq_swift_code, :bank_name, :address, :country_iso2_code, :country_name);`, banks)
	if err != nil {
		return err
	}

	return nil
}
