package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/mwojtyna/swift-api/internal/db"
	"github.com/mwojtyna/swift-api/internal/utils"
)

func (s *APIServer) getSwiftCodeDetailsV1(w http.ResponseWriter, r *http.Request) {
	swiftCode := r.PathValue("swiftCode")
	if swiftCode == "" {
		WriteError(w, 400)
		return
	}

	bank, err := db.GetBank(s.db, swiftCode)
	if errors.Is(err, sql.ErrNoRows) {
		WriteError(w, 404)
		return
	}
	if err != nil {
		WriteError(w, 500)
		return
	}

	if bank.IsHq() {
		branchesRaw, err := db.GetBankBranches(s.db, swiftCode)
		if err != nil {
			WriteError(w, 500)
			return
		}

		branches := utils.Map(branchesRaw, func(b db.Bank) SwiftCodeDetailsHqBranch {
			return SwiftCodeDetailsHqBranch{
				Address:       b.Address,
				BankName:      b.BankName,
				CountryISO2:   b.CountryISO2Code,
				IsHeadquarter: b.IsHq(),
				SwiftCode:     b.SwiftCode,
			}
		})

		res := SwiftCodeDetailsHqResponse{
			Address:       bank.Address,
			BankName:      bank.BankName,
			CountryISO2:   bank.CountryISO2Code,
			CountryName:   bank.CountryName,
			IsHeadquarter: bank.IsHq(),
			SwiftCode:     bank.SwiftCode,
			Branches:      branches,
		}

		err = WriteJSON(w, 200, res)
		if err != nil {
			WriteError(w, 500)
			return
		}
	} else {
		res := SwiftCodeDetailsBranchResponse{
			Address:       bank.Address,
			BankName:      bank.BankName,
			CountryISO2:   bank.CountryISO2Code,
			CountryName:   bank.CountryName,
			IsHeadquarter: bank.IsHq(),
			SwiftCode:     bank.SwiftCode,
		}

		err = WriteJSON(w, 200, res)
		if err != nil {
			WriteError(w, 500)
			return
		}
	}
}
