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
			Address:       b.Address,
			BankName:      b.BankName,
			CountryISO2:   b.CountryISO2Code,
			CountryName:   b.CountryName,
			IsHeadquarter: b.HqSwiftCode.Valid == false, // If HQ code is NULL, then this bank is HQ
			SwiftCode:     b.SwiftCode,
		}
	})

	res := GetSwiftCodeResponse{
		Address:       bank.Address,
		BankName:      bank.BankName,
		CountryISO2:   bank.CountryISO2Code,
		CountryName:   bank.CountryName,
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
