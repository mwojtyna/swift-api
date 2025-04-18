package api

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

type ApiServer struct {
	address  string
	db       *sqlx.DB
	logger   *log.Logger
	validate *validator.Validate
}

type MessageRes struct {
	Message string `json:"message"`
}

type GetSwiftCodeHqRes struct {
	Address       string                 `json:"address"`
	BankName      string                 `json:"bankName"`
	CountryISO2   string                 `json:"countryISO2"`
	CountryName   string                 `json:"countryName"`
	IsHeadquarter bool                   `json:"isHeadquarter"`
	SwiftCode     string                 `json:"swiftCode"`
	Branches      []GetSwiftCodeHqBranch `json:"branches"`
}

type GetSwiftCodeHqBranch struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type GetSwiftCodeBranchRes struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type GetSwiftCodesForCountryRes struct {
	CountryISO2 string                             `json:"countryISO2"`
	CountryName string                             `json:"countryName"`
	SwiftCodes  []GetSwiftCodesForCountrySwiftCode `json:"swiftCodes"`
}

type GetSwiftCodesForCountrySwiftCode struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type AddSwiftCodeReq struct {
	Address       string `json:"address" validate:"required"`
	BankName      string `json:"bankName" validate:"required"`
	CountryISO2   string `json:"countryISO2" validate:"required,uppercase,country_code"`
	CountryName   string `json:"countryName" validate:"required,uppercase"`
	IsHeadquarter bool   `json:"isHeadquarter"` // Can't validate:"required" because zero-value for bool is false, meaning a branch bank won't be accepted
	SwiftCode     string `json:"swiftCode" validate:"required,len=11"`
}
