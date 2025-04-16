package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/mwojtyna/swift-api/internal/db"
	"github.com/mwojtyna/swift-api/internal/utils"
)

func (s *APIServer) getSwiftCodeDetailsV1(w http.ResponseWriter, r *http.Request) error {
	swiftCode := r.PathValue("swiftCode")
	if swiftCode == "" {
		WriteError(w, http.StatusBadRequest)
		return nil
	}

	bank, err := db.GetBank(s.db, swiftCode)
	if errors.Is(err, sql.ErrNoRows) {
		WriteError(w, http.StatusNotFound)
		return nil
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError)
		return err
	}

	if bank.IsHQ() {
		branchesRaw, err := db.GetBankBranches(s.db, swiftCode)
		if err != nil {
			WriteError(w, http.StatusInternalServerError)
			return err
		}

		branches := utils.Map(branchesRaw, func(b db.Bank) SwiftCodeDetailsHqBranch {
			return SwiftCodeDetailsHqBranch{
				Address:       b.Address,
				BankName:      b.BankName,
				CountryISO2:   b.CountryISO2Code,
				IsHeadquarter: b.IsHQ(),
				SwiftCode:     b.SwiftCode,
			}
		})

		res := SwiftCodeDetailsHqResponse{
			Address:       bank.Address,
			BankName:      bank.BankName,
			CountryISO2:   bank.CountryISO2Code,
			CountryName:   bank.CountryName,
			IsHeadquarter: bank.IsHQ(),
			SwiftCode:     bank.SwiftCode,
			Branches:      branches,
		}

		err = WriteJSON(w, http.StatusOK, res)
		if err != nil {
			WriteError(w, http.StatusInternalServerError)
			return err
		}
	} else {
		res := SwiftCodeDetailsBranchResponse{
			Address:       bank.Address,
			BankName:      bank.BankName,
			CountryISO2:   bank.CountryISO2Code,
			CountryName:   bank.CountryName,
			IsHeadquarter: bank.IsHQ(),
			SwiftCode:     bank.SwiftCode,
		}

		err = WriteJSON(w, http.StatusOK, res)
		if err != nil {
			WriteError(w, http.StatusInternalServerError)
			return err
		}
	}

	return nil
}

func (s *APIServer) getSwiftCodesForCountryV1(w http.ResponseWriter, r *http.Request) error {
	countryCode := r.PathValue("countryISO2code")

	banks, err := db.GetBanksInCountry(s.db, countryCode)
	if len(banks) == 0 {
		WriteError(w, http.StatusNotFound)
		return nil
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError)
		return err
	}

	codes := utils.Map(banks, func(b db.Bank) SwiftCodesForCountrySwiftCode {
		return SwiftCodesForCountrySwiftCode{
			Address:       b.Address,
			BankName:      b.BankName,
			CountryISO2:   b.CountryISO2Code,
			IsHeadquarter: b.IsHQ(),
			SwiftCode:     b.SwiftCode,
		}
	})
	res := SwiftCodesForCountryResponse{
		// Since all banks are from the same country, just get the country data from any bank so we don't have to query the DB
		CountryISO2: banks[0].CountryISO2Code,
		CountryName: banks[0].CountryName,
		SwiftCodes:  codes,
	}

	err = WriteJSON(w, http.StatusOK, res)
	if err != nil {
		WriteError(w, http.StatusInternalServerError)
		return err
	}

	return nil
}
