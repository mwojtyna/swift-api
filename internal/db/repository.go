package db

import (
	"context"
	"database/sql"
	"fmt"
)

func InsertCountries(db *sql.DB, countries []Country) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO country (iso2_code, country_name, time_zone) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	for _, c := range countries {
		// Insert null if hq swift code is empty
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

func InsertBanks(db *sql.DB, banks []Bank) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO bank (swift_code, hq_swift_code, country_iso2_code, bank_name, address) VALUES ($1, $2, $3, $4, $5)")
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
