package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/lib/pq"
	"github.com/mwojtyna/swift-api/internal/db"
	"github.com/mwojtyna/swift-api/internal/utils"
)

// NOTE: Return error in function only if status is 500!

func (s *APIServer) handleGetSwiftCodeV1(w http.ResponseWriter, r *http.Request) error {
	swiftCode := r.PathValue("swiftCode")
	if swiftCode == "" {
		WriteHTTPError(w, http.StatusBadRequest)
		return nil
	}

	bank, err := db.GetBank(s.db, swiftCode)
	if errors.Is(err, sql.ErrNoRows) {
		WriteHTTPError(w, http.StatusNotFound)
		return nil
	}
	if err != nil {
		return err
	}

	if bank.IsHq() {
		branchesRaw, err := db.GetBankBranches(s.db, swiftCode)
		if err != nil {
			return err
		}

		branches := utils.Map(branchesRaw, func(b db.Bank) GetSwiftCodeHqBranch {
			return GetSwiftCodeHqBranch{
				Address:       b.Address,
				BankName:      b.BankName,
				CountryISO2:   b.CountryISO2Code,
				IsHeadquarter: b.IsHq(),
				SwiftCode:     b.SwiftCode,
			}
		})

		res := GetSwiftCodeHqRes{
			Address:       bank.Address,
			BankName:      bank.BankName,
			CountryISO2:   bank.CountryISO2Code,
			CountryName:   bank.CountryName,
			IsHeadquarter: bank.IsHq(),
			SwiftCode:     bank.SwiftCode,
			Branches:      branches,
		}

		err = WriteJSON(w, http.StatusOK, res)
		if err != nil {
			return err
		}
	} else {
		res := GetSwiftCodeBranchRes{
			Address:       bank.Address,
			BankName:      bank.BankName,
			CountryISO2:   bank.CountryISO2Code,
			CountryName:   bank.CountryName,
			IsHeadquarter: bank.IsHq(),
			SwiftCode:     bank.SwiftCode,
		}

		err = WriteJSON(w, http.StatusOK, res)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *APIServer) handleGetSwiftCodesForCountryV1(w http.ResponseWriter, r *http.Request) error {
	countryCode := r.PathValue("countryISO2code")

	banks, err := db.GetBanksInCountry(s.db, countryCode)
	if len(banks) == 0 {
		WriteHTTPError(w, http.StatusNotFound)
		return nil
	}
	if err != nil {
		return err
	}

	codes := utils.Map(banks, func(b db.Bank) GetSwiftCodesForCountrySwiftCode {
		return GetSwiftCodesForCountrySwiftCode{
			Address:       b.Address,
			BankName:      b.BankName,
			CountryISO2:   b.CountryISO2Code,
			IsHeadquarter: b.IsHq(),
			SwiftCode:     b.SwiftCode,
		}
	})
	res := GetSwiftCodesForCountryRes{
		// Since all banks are from the same country, just get the country data from any bank so we don't have to query the DB
		CountryISO2: banks[0].CountryISO2Code,
		CountryName: banks[0].CountryName,
		SwiftCodes:  codes,
	}

	err = WriteJSON(w, http.StatusOK, res)
	if err != nil {
		return err
	}

	return nil
}

func (s *APIServer) handleAddSwiftCodeV1(w http.ResponseWriter, r *http.Request) error {
	var req AddSwiftCodeReq

	err := ReadJSON(w, r, &req)
	if err != nil {
		res := MessageRes{Message: err.Error()}
		WriteJSON(w, http.StatusBadRequest, res)
		return nil
	}

	// HQ handling
	isHq, hqCode := utils.IsSwiftCodeHq(req.SwiftCode)
	if (isHq && !req.IsHeadquarter) || (!isHq && req.IsHeadquarter) {
		res := MessageRes{Message: "isHeadquarter disagrees with swiftCode"}
		WriteJSON(w, http.StatusUnprocessableEntity, res)
		return nil
	}

	dbHqCode := sql.NullString{}
	if !isHq {
		exists, err := db.CheckBankHqExists(s.db, hqCode)
		if err != nil {
			return err
		}

		if exists {
			dbHqCode = sql.NullString{
				String: hqCode,
				Valid:  true,
			}
		}
	}

	bank := db.Bank{
		SwiftCode:       req.SwiftCode,
		HqSwiftCode:     dbHqCode,
		BankName:        req.BankName,
		Address:         req.Address,
		CountryISO2Code: req.CountryISO2,
		CountryName:     req.CountryName,
	}

	pgErr, isPgErr := db.InsertBank(s.db, bank).(*pq.Error)
	if isPgErr && pgErr.Code == db.UniqueViolationErrorCode {
		WriteHTTPError(w, http.StatusConflict)
		return nil
	} else if pgErr != nil {
		return pgErr
	}

	res := MessageRes{Message: "SWIFT code added successfully"}
	err = WriteJSON(w, http.StatusCreated, res)
	if err != nil {
		return err
	}

	return nil
}
