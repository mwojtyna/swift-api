package api

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type APIServer struct {
	address string
	db      *sqlx.DB
	logger  *log.Logger
}

type BankBranch struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type GetSwiftCodeResponse struct {
	Address       string       `json:"address"`
	BankName      string       `json:"bankName"`
	CountryISO2   string       `json:"countryISO2"`
	CountryName   string       `json:"countryName"`
	IsHeadquarter bool         `json:"isHeadquarter"`
	SwiftCode     string       `json:"swiftCode"`
	Branches      []BankBranch `json:"branches"`
}
