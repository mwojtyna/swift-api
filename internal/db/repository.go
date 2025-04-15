package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func GetCountry(db *sqlx.DB, iso2Code string) (Country, error) {
	var country []Country
	err := db.Select(&country, "SELECT * FROM country WHERE iso2_code=$1;", iso2Code)
	if err != nil {
		return Country{}, err
	}

	return country[0], nil
}

func GetBank(db *sqlx.DB, swiftCode string) (Bank, error) {
	var bank []Bank
	err := db.Select(&bank, "SELECT * FROM bank WHERE swift_code=$1;", swiftCode)
	if err != nil {
		return Bank{}, err
	}

	return bank[0], nil
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

func InsertCountries(db *sqlx.DB, countries []Country) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO country (iso2_code, country_name, time_zone) VALUES ($1, $2, $3);")
	if err != nil {
		return err
	}

	for _, c := range countries {
		_, err := stmt.Exec(c.ISO2Code, c.CountryName, c.TimeZone)
		if err != nil {
			return fmt.Errorf(`%s - %s`, c, err.Error())
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func InsertBanks(db *sqlx.DB, banks []Bank) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO bank (swift_code, hq_swift_code, country_iso2_code, bank_name, address) VALUES ($1, $2, $3, $4, $5);")
	if err != nil {
		return err
	}

	for _, b := range banks {
		_, err := stmt.Exec(b.SwiftCode, b.HqSwiftCode, b.CountryISO2Code, b.BankName, b.Address)
		if err != nil {
			return fmt.Errorf("%+v - %s", b, err.Error())
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
