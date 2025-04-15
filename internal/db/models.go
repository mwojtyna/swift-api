package db

import "database/sql"

type Bank struct {
	SwiftCode       string
	HqSwiftCode     sql.NullString
	CountryISO2Code string
	BankName        string
	Address         string
}

type Country struct {
	ISO2Code    string
	CountryName string
	TimeZone    string
}
