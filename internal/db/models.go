package db

type Bank struct {
	SwiftCode       string
	HqSwiftCode     string
	CountryISO2Code string
	BankName        string
	Address         string
}

type Country struct {
	ISO2Code    string
	CountryName string
	TimeZone    string
}
