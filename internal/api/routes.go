package api

import (
	"net/http"

	"github.com/mwojtyna/swift-api/internal/db"
	"github.com/mwojtyna/swift-api/internal/utils"
)

func (s *APIServer) getBankBySwiftCodeV1(w http.ResponseWriter, r *http.Request) {
	swiftCode := r.PathValue("swiftCode")
	if swiftCode == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bank, err := db.GetBank(s.db, swiftCode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	branchesRaw, err := db.GetBankBranches(s.db, swiftCode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	branches := utils.Map(branchesRaw, func(b db.Bank) BankBranch {
		return BankBranch{
			Address:       bank.Address,
			BankName:      bank.BankName,
			CountryISO2:   bank.CountryISO2Code,
			CountryName:   "CHUJ",
			IsHeadquarter: bank.HqSwiftCode.Valid == false, // If HQ code is NULL, then this bank is HQ
			SwiftCode:     bank.SwiftCode,
		}
	})

	country, err := db.GetCountry(s.db, bank.CountryISO2Code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := GetSwiftCodeResponse{
		Address:       bank.Address,
		BankName:      bank.BankName,
		CountryISO2:   bank.CountryISO2Code,
		CountryName:   country.CountryName,
		IsHeadquarter: bank.HqSwiftCode.Valid == false, // If HQ code is NULL, then this bank is HQ
		SwiftCode:     bank.SwiftCode,
		Branches:      branches,
	}

	err = WriteJSON(w, 200, res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
