package db

type Bank struct {
	CountryISO2Code string
	SwiftCode       string
	CodeType        string
	Name            string
	Address         string
	TownName        string
}

type Country struct {
	ISO2Code string
	Name     string
	TimeZone string
}
