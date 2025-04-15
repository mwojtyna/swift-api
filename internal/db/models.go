package db

import "database/sql"

type Bank struct {
	SwiftCode       string         `db:"swift_code"`
	HqSwiftCode     sql.NullString `db:"hq_swift_code"`
	CountryISO2Code string         `db:"country_iso2_code"`
	BankName        string         `db:"bank_name"`
	Address         string         `db:"address"`
}

type Country struct {
	ISO2Code    string `db:"iso2_code"`
	CountryName string `db:"country_name"`
	TimeZone    string `db:"time_zone"`
}
