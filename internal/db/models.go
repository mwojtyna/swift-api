package db

import (
	"database/sql"
)

type Bank struct {
	SwiftCode       string         `db:"swift_code"`
	HqSwiftCode     sql.NullString `db:"hq_swift_code"`
	BankName        string         `db:"bank_name"`
	Address         string         `db:"address"`
	CountryISO2Code string         `db:"country_iso2_code"`
	CountryName     string         `db:"country_name"`
}

func (b *Bank) IsHq() bool {
	// If HQ code is NULL, then this bank is HQ
	return !b.HqSwiftCode.Valid
}
