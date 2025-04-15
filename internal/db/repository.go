package db

import (
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

func InsertBanks(db *sqlx.DB, banks []Bank) error {
	_, err := db.NamedExec(`INSERT INTO bank (swift_code, hq_swift_code, bank_name, address, country_iso2_code, country_name) 
		VALUES (:swift_code, :hq_swift_code, :bank_name, :address, :country_iso2_code, :country_name);`, banks)
	if err != nil {
		return err
	}

	return nil
}
