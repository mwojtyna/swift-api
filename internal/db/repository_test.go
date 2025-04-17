package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/lib/pq"
	"github.com/mwojtyna/swift-api/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBank(t *testing.T) {
	t.Parallel()

	utils.TestWithPostgres(func(args utils.TestWithPostgresArgs) {
		db, err := Connect(args.Env.DB_USER, args.Env.DB_PASS, args.Env.DB_NAME, args.Env.DB_PORT)
		require.NoError(t, err)
		t.Cleanup(func() {
			db.Close()
		})

		// Test data
		hqBank := Bank{
			SwiftCode:       "HQTESTBANK",
			HqSwiftCode:     sql.NullString{},
			BankName:        "HQ Bank",
			Address:         "456 HQ Street",
			CountryISO2Code: "GB",
			CountryName:     "United Kingdom",
		}

		branchBank := Bank{
			SwiftCode:       "BRANCHBANK",
			HqSwiftCode:     sql.NullString{String: "HQTESTBANK", Valid: true},
			BankName:        "Branch Bank",
			Address:         "789 Branch Ave",
			CountryISO2Code: "GB",
			CountryName:     "United Kingdom",
		}

		t.Run("successfully retrieves HQ bank", func(t *testing.T) {
			// Arrange
			err := InsertBank(db, hqBank)
			require.NoError(t, err)
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act
			result, err := GetBank(db, hqBank.SwiftCode)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, hqBank, result)
		})

		t.Run("successfully retrieves branch bank", func(t *testing.T) {
			// Arrange
			err := InsertBank(db, hqBank)
			require.NoError(t, err)
			err = InsertBank(db, branchBank)
			require.NoError(t, err)
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act
			result, err := GetBank(db, branchBank.SwiftCode)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, branchBank, result)
		})

		t.Run("returns error for non-existent bank", func(t *testing.T) {
			// Act
			result, err := GetBank(db, "NONEXISTENT")

			// Assert
			require.Error(t, err)
			assert.True(t, errors.Is(err, sql.ErrNoRows))
			assert.Equal(t, Bank{}, result)
		})

	})
}

func TestGetBankBranches(t *testing.T) {
	t.Parallel()
	utils.TestWithPostgres(func(args utils.TestWithPostgresArgs) {
		db, err := Connect(args.Env.DB_USER, args.Env.DB_PASS, args.Env.DB_NAME, args.Env.DB_PORT)
		require.NoError(t, err)
		t.Cleanup(func() {
			db.Close()
		})

		// Test data
		hqBank := Bank{
			SwiftCode:       "HQTESTBANK",
			HqSwiftCode:     sql.NullString{},
			BankName:        "HQ Bank",
			Address:         "456 HQ Street",
			CountryISO2Code: "GB",
			CountryName:     "United Kingdom",
		}

		branch1 := Bank{
			SwiftCode:       "BRANCH001",
			HqSwiftCode:     sql.NullString{String: "HQTESTBANK", Valid: true},
			BankName:        "Branch 1",
			Address:         "123 Branch St",
			CountryISO2Code: "GB",
			CountryName:     "United Kingdom",
		}

		branch2 := Bank{
			SwiftCode:       "BRANCH002",
			HqSwiftCode:     sql.NullString{String: "HQTESTBANK", Valid: true},
			BankName:        "Branch 2",
			Address:         "456 Branch Ave",
			CountryISO2Code: "GB",
			CountryName:     "United Kingdom",
		}

		otherBank := Bank{
			SwiftCode:       "OTHERBANK",
			HqSwiftCode:     sql.NullString{},
			BankName:        "Other Bank",
			Address:         "789 Other St",
			CountryISO2Code: "US",
			CountryName:     "United States",
		}

		t.Run("returns branches for valid HQ bank", func(t *testing.T) {
			// Arrange
			err := InsertBank(db, hqBank)
			require.NoError(t, err)
			err = InsertBank(db, branch1)
			require.NoError(t, err)
			err = InsertBank(db, branch2)
			require.NoError(t, err)
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act
			branches, err := GetBankBranches(db, hqBank.SwiftCode)

			// Assert
			require.NoError(t, err)
			require.Len(t, branches, 2)
			assert.Contains(t, branches, branch1)
			assert.Contains(t, branches, branch2)
		})

		t.Run("returns empty slice for HQ with no branches", func(t *testing.T) {
			// Arrange
			err := InsertBank(db, hqBank)
			require.NoError(t, err)
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act
			branches, err := GetBankBranches(db, hqBank.SwiftCode)

			// Assert
			require.NoError(t, err)
			assert.Empty(t, branches)
		})

		t.Run("returns empty slice for non-existent HQ", func(t *testing.T) {
			// Act
			branches, err := GetBankBranches(db, "NONEXISTENT")

			// Assert
			require.NoError(t, err)
			assert.Empty(t, branches)
		})

		t.Run("returns empty slice when querying a branch", func(t *testing.T) {
			// Arrange
			err := InsertBank(db, hqBank)
			require.NoError(t, err)
			err = InsertBank(db, branch1)
			require.NoError(t, err)
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act
			branches, err := GetBankBranches(db, branch1.SwiftCode)

			// Assert
			require.NoError(t, err)
			assert.Empty(t, branches)
		})

		t.Run("returns only branches for specified HQ", func(t *testing.T) {
			// Arrange
			err := InsertBank(db, hqBank)
			require.NoError(t, err)
			err = InsertBank(db, branch1)
			require.NoError(t, err)
			err = InsertBank(db, otherBank)
			require.NoError(t, err)
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act
			branches, err := GetBankBranches(db, hqBank.SwiftCode)

			// Assert
			require.NoError(t, err)
			require.Len(t, branches, 1)
			assert.Contains(t, branches, branch1)
			assert.NotContains(t, branches, otherBank)
		})
	})
}

func TestGetBanksInCountry(t *testing.T) {
	t.Parallel()
	utils.TestWithPostgres(func(args utils.TestWithPostgresArgs) {
		db, err := Connect(args.Env.DB_USER, args.Env.DB_PASS, args.Env.DB_NAME, args.Env.DB_PORT)
		require.NoError(t, err)
		t.Cleanup(func() {
			db.Close()
		})

		// Create test data
		usBank1 := Bank{
			SwiftCode:       "USBANK001",
			HqSwiftCode:     sql.NullString{},
			BankName:        "US Bank 1",
			Address:         "123 Main St",
			CountryISO2Code: "US",
			CountryName:     "United States",
		}

		usBank2 := Bank{
			SwiftCode:       "USBANK002",
			HqSwiftCode:     sql.NullString{String: "USBANK001", Valid: true},
			BankName:        "US Bank 2",
			Address:         "456 Oak Ave",
			CountryISO2Code: "US",
			CountryName:     "United States",
		}

		ukBank := Bank{
			SwiftCode:       "UKBANK001",
			HqSwiftCode:     sql.NullString{},
			BankName:        "UK Bank",
			Address:         "789 High St",
			CountryISO2Code: "GB",
			CountryName:     "United Kingdom",
		}

		t.Run("returns banks for valid country code", func(t *testing.T) {
			// Arrange
			_, err = db.NamedExec(`INSERT INTO bank VALUES 
			(:swift_code, :hq_swift_code, :bank_name, :address, :country_iso2_code, :country_name)`,
				usBank1)
			require.NoError(t, err)

			_, err = db.NamedExec(`INSERT INTO bank VALUES 
				(:swift_code, :hq_swift_code, :bank_name, :address, :country_iso2_code, :country_name)`,
				usBank2)
			require.NoError(t, err)

			_, err = db.NamedExec(`INSERT INTO bank VALUES 
				(:swift_code, :hq_swift_code, :bank_name, :address, :country_iso2_code, :country_name)`,
				ukBank)
			require.NoError(t, err)

			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act
			banks, err := GetBanksInCountry(db, "US")

			// Assert
			require.NoError(t, err)
			require.Len(t, banks, 2)
			assert.Contains(t, banks, usBank1)
			assert.Contains(t, banks, usBank2)
			assert.NotContains(t, banks, ukBank)
		})

		t.Run("returns empty slice for country with no banks", func(t *testing.T) {
			// Arrange
			_, err = db.NamedExec(`INSERT INTO bank VALUES 
			(:swift_code, :hq_swift_code, :bank_name, :address, :country_iso2_code, :country_name)`,
				usBank1)
			require.NoError(t, err)
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act
			banks, err := GetBanksInCountry(db, "FR")

			// Assert
			require.NoError(t, err)
			assert.Empty(t, banks)
		})
	})
}

func TestCheckBankHqExists(t *testing.T) {
	t.Parallel()
	utils.TestWithPostgres(func(args utils.TestWithPostgresArgs) {
		db, err := Connect(args.Env.DB_USER, args.Env.DB_PASS, args.Env.DB_NAME, args.Env.DB_PORT)
		require.NoError(t, err)
		t.Cleanup(func() {
			db.Close()
		})

		// Test data templates
		hqBank := Bank{
			SwiftCode:       "HQTESTBANK",
			HqSwiftCode:     sql.NullString{},
			BankName:        "HQ Bank",
			Address:         "456 HQ Street",
			CountryISO2Code: "GB",
			CountryName:     "United Kingdom",
		}

		branchBank := Bank{
			SwiftCode:       "BRANCHBANK",
			HqSwiftCode:     sql.NullString{String: "HQTESTBANK", Valid: true},
			BankName:        "Branch Bank",
			Address:         "789 Branch Ave",
			CountryISO2Code: "GB",
			CountryName:     "United Kingdom",
		}

		t.Run("returns false when HQ doesn't exist", func(t *testing.T) {
			// Arrange
			db.NamedExec(`INSERT INTO bank VALUES 
				(:swift_code, :hq_swift_code, :bank_name, :address, :country_iso2_code, :country_name)`,
				branchBank)
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act
			exists, err := CheckBankHqExists(db, branchBank.SwiftCode)

			// Assert
			require.NoError(t, err)
			assert.False(t, exists)
		})

		t.Run("returns true when branch's HQ exists", func(t *testing.T) {
			// Arrange
			_, err := db.NamedExec(`INSERT INTO bank VALUES 
				(:swift_code, :hq_swift_code, :bank_name, :address, :country_iso2_code, :country_name)`,
				hqBank)
			require.NoError(t, err)
			_, err = db.NamedExec(`INSERT INTO bank VALUES 
				(:swift_code, :hq_swift_code, :bank_name, :address, :country_iso2_code, :country_name)`,
				branchBank)
			require.NoError(t, err)
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act
			exists, err := CheckBankHqExists(db, branchBank.SwiftCode)

			// Assert
			require.NoError(t, err)
			assert.True(t, exists)
		})
	})
}

func TestInsertBanks(t *testing.T) {
	t.Parallel()
	utils.TestWithPostgres(func(args utils.TestWithPostgresArgs) {
		db, err := Connect(args.Env.DB_USER, args.Env.DB_PASS, args.Env.DB_NAME, args.Env.DB_PORT)
		require.NoError(t, err)
		t.Cleanup(func() {
			db.Close()
		})

		// Test data
		hqBank := Bank{
			SwiftCode:       "HQTESTBANK",
			HqSwiftCode:     sql.NullString{},
			BankName:        "HQ Bank",
			Address:         "456 HQ Street",
			CountryISO2Code: "GB",
			CountryName:     "United Kingdom",
		}

		branchBank := Bank{
			SwiftCode:       "BRANCHBANK",
			HqSwiftCode:     sql.NullString{String: "HQTESTBANK", Valid: true},
			BankName:        "Branch Bank",
			Address:         "789 Branch Ave",
			CountryISO2Code: "GB",
			CountryName:     "United Kingdom",
		}

		t.Run("successfully inserts multiple banks", func(t *testing.T) {
			// Arrange
			banks := []Bank{hqBank, branchBank}
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act
			err := InsertBanks(db, banks)

			// Assert
			require.NoError(t, err)

			// Verify data was inserted
			var count int
			err = db.Get(&count, "SELECT COUNT(*) FROM bank WHERE swift_code IN ($1, $2)",
				hqBank.SwiftCode, branchBank.SwiftCode)
			require.NoError(t, err)
			assert.Equal(t, 2, count)
		})

		t.Run("returns error for duplicate swift code", func(t *testing.T) {
			// Arrange - insert first bank
			_, err := db.NamedExec(`INSERT INTO bank VALUES 
				(:swift_code, :hq_swift_code, :bank_name, :address, :country_iso2_code, :country_name)`,
				hqBank)
			require.NoError(t, err)
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Try to insert duplicate
			banks := []Bank{hqBank}

			// Act
			err = InsertBanks(db, banks)

			// Assert
			require.Error(t, err)
			pgErr, ok := err.(*pq.Error)
			require.True(t, ok)
			assert.Equal(t, pgErr.Code, UniqueViolationErrorCode)
		})

		t.Run("returns error for invalid hq_swift_code reference", func(t *testing.T) {
			// Arrange
			invalidBranch := Bank{
				SwiftCode:       "BADBRANCH",
				HqSwiftCode:     sql.NullString{String: "NONEXISTENT", Valid: true},
				BankName:        "Invalid Branch",
				Address:         "123 Invalid St",
				CountryISO2Code: "GB",
				CountryName:     "United Kingdom",
			}
			banks := []Bank{invalidBranch}
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act
			err := InsertBanks(db, banks)

			// Assert
			require.Error(t, err)
			pgErr, ok := err.(*pq.Error)
			require.True(t, ok)
			assert.Equal(t, pgErr.Code, ForeignKeyViolationErrorCode)
		})
	})
}

func TestInsertBank(t *testing.T) {
	t.Parallel()
	utils.TestWithPostgres(func(args utils.TestWithPostgresArgs) {
		db, err := Connect(args.Env.DB_USER, args.Env.DB_PASS, args.Env.DB_NAME, args.Env.DB_PORT)
		require.NoError(t, err)
		t.Cleanup(func() {
			db.Close()
		})

		// Test data
		hqBank := Bank{
			SwiftCode:       "HQTESTBANK",
			HqSwiftCode:     sql.NullString{},
			BankName:        "HQ Bank",
			Address:         "456 HQ Street",
			CountryISO2Code: "GB",
			CountryName:     "United Kingdom",
		}

		branchBank := Bank{
			SwiftCode:       "BRANCHBANK",
			HqSwiftCode:     sql.NullString{String: "HQTESTBANK", Valid: true},
			BankName:        "Branch Bank",
			Address:         "789 Branch Ave",
			CountryISO2Code: "GB",
			CountryName:     "United Kingdom",
		}

		t.Run("successfully inserts HQ bank", func(t *testing.T) {
			// Arrange
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act
			err := InsertBank(db, hqBank)

			// Assert
			require.NoError(t, err)

			// Verify insertion
			var insertedBank Bank
			err = db.Get(&insertedBank, "SELECT * FROM bank WHERE swift_code = $1", hqBank.SwiftCode)
			require.NoError(t, err)
			assert.Equal(t, hqBank, insertedBank)
		})

		t.Run("successfully inserts branch bank", func(t *testing.T) {
			// Arrange - first insert HQ bank
			err := InsertBank(db, hqBank)
			require.NoError(t, err)
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act
			err = InsertBank(db, branchBank)

			// Assert
			require.NoError(t, err)

			// Verify insertion
			var insertedBank Bank
			err = db.Get(&insertedBank, "SELECT * FROM bank WHERE swift_code = $1", branchBank.SwiftCode)
			require.NoError(t, err)
			assert.Equal(t, branchBank, insertedBank)
		})

		t.Run("returns error for duplicate swift code", func(t *testing.T) {
			// Arrange
			err := InsertBank(db, hqBank)
			require.NoError(t, err)
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act - try to insert same bank again
			err = InsertBank(db, hqBank)

			// Assert
			require.Error(t, err)
			pgErr, ok := err.(*pq.Error)
			require.True(t, ok)
			assert.Equal(t, pgErr.Code, UniqueViolationErrorCode)
		})

		t.Run("returns error for invalid hq_swift_code reference", func(t *testing.T) {
			// Arrange
			invalidBranch := Bank{
				SwiftCode:       "BADBRANCH",
				HqSwiftCode:     sql.NullString{String: "NONEXISTENT", Valid: true},
				BankName:        "Invalid Branch",
				Address:         "123 Invalid St",
				CountryISO2Code: "GB",
				CountryName:     "United Kingdom",
			}
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act
			err := InsertBank(db, invalidBranch)

			// Assert
			require.Error(t, err)
			pgErr, ok := err.(*pq.Error)
			require.True(t, ok)
			assert.Equal(t, pgErr.Code, ForeignKeyViolationErrorCode)
		})
	})
}

func TestDeleteBank(t *testing.T) {
	t.Parallel()
	utils.TestWithPostgres(func(args utils.TestWithPostgresArgs) {
		db, err := Connect(args.Env.DB_USER, args.Env.DB_PASS, args.Env.DB_NAME, args.Env.DB_PORT)
		require.NoError(t, err)
		t.Cleanup(func() {
			db.Close()
		})

		// Test data
		hqBank := Bank{
			SwiftCode:       "HQTESTBANK",
			HqSwiftCode:     sql.NullString{},
			BankName:        "HQ Bank",
			Address:         "456 HQ Street",
			CountryISO2Code: "GB",
			CountryName:     "United Kingdom",
		}

		branchBank := Bank{
			SwiftCode:       "BRANCHBANK",
			HqSwiftCode:     sql.NullString{String: "HQTESTBANK", Valid: true},
			BankName:        "Branch Bank",
			Address:         "789 Branch Ave",
			CountryISO2Code: "GB",
			CountryName:     "United Kingdom",
		}

		t.Run("successfully deletes HQ bank and nullifies branches", func(t *testing.T) {
			// Arrange
			_, err := db.NamedExec(`INSERT INTO bank VALUES 
				(:swift_code, :hq_swift_code, :bank_name, :address, :country_iso2_code, :country_name)`,
				hqBank)
			require.NoError(t, err)
			_, err = db.NamedExec(`INSERT INTO bank VALUES 
				(:swift_code, :hq_swift_code, :bank_name, :address, :country_iso2_code, :country_name)`,
				branchBank)
			require.NoError(t, err)
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act
			err = DeleteBank(db, hqBank.SwiftCode)

			// Assert
			require.NoError(t, err)

			// Verify HQ bank was deleted
			var count int
			err = db.Get(&count, "SELECT COUNT(*) FROM bank WHERE swift_code = $1", hqBank.SwiftCode)
			require.NoError(t, err)
			assert.Equal(t, 0, count)

			// Verify branch's hq_swift_code was set to NULL
			var branchHqCode sql.NullString
			err = db.Get(&branchHqCode, "SELECT hq_swift_code FROM bank WHERE swift_code = $1", branchBank.SwiftCode)
			require.NoError(t, err)
			assert.False(t, branchHqCode.Valid)
		})

		t.Run("successfully deletes standalone bank", func(t *testing.T) {
			// Arrange
			_, err := db.NamedExec(`INSERT INTO bank VALUES 
				(:swift_code, :hq_swift_code, :bank_name, :address, :country_iso2_code, :country_name)`,
				hqBank)
			require.NoError(t, err)
			t.Cleanup(func() {
				db.Exec("TRUNCATE bank")
			})

			// Act
			err = DeleteBank(db, hqBank.SwiftCode)

			// Assert
			require.NoError(t, err)

			// Verify bank was deleted
			var count int
			err = db.Get(&count, "SELECT COUNT(*) FROM bank WHERE swift_code = $1", hqBank.SwiftCode)
			require.NoError(t, err)
			assert.Equal(t, 0, count)
		})

		t.Run("returns error when bank doesn't exist", func(t *testing.T) {
			// Act
			err := DeleteBank(db, "NONEXISTENT")

			// Assert
			require.Error(t, err)
			assert.True(t, errors.Is(err, sql.ErrNoRows))
		})
	})
}
