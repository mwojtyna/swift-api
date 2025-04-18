package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/mwojtyna/swift-api/internal/db"
	"github.com/mwojtyna/swift-api/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testApiArgs struct {
	router *http.ServeMux
	db     *sqlx.DB
}

func testApi(f func(testApiArgs)) {
	utils.TestWithPostgres(func(args utils.TestWithPostgresArgs) {
		pg, err := db.Connect(args.Env.DB_USER, args.Env.DB_PASS, args.Env.DB_NAME, args.Env.DB_PORT)
		if err != nil {
			log.Fatalln("failed to connect to db")
		}

		var logBuf bytes.Buffer
		logger := log.New(&logBuf, "", 0)
		api := NewApiServer(":"+args.Env.API_PORT, pg, logger)
		router := api.NewRouter()

		f(testApiArgs{router: router, db: pg})
	})
}

var (
	hqBank = db.Bank{
		SwiftCode:       "ABCDEFGHXXX",
		HqSwiftCode:     sql.NullString{},
		BankName:        "HQ Bank",
		Address:         "456 HQ Street",
		CountryISO2Code: "GB",
		CountryName:     "UNITED KINGDOM",
	}
	branchBank1 = db.Bank{
		SwiftCode:       "ABCDEFGH001",
		HqSwiftCode:     sql.NullString{String: hqBank.SwiftCode, Valid: true},
		BankName:        "Branch Bank",
		Address:         "456 Branch Street",
		CountryISO2Code: "GB",
		CountryName:     "UNITED KINGDOM",
	}
	branchBank2 = db.Bank{
		SwiftCode:       "ABCDEFGH002",
		HqSwiftCode:     sql.NullString{String: hqBank.SwiftCode, Valid: true},
		BankName:        "Branch Bank",
		Address:         "456 Branch Street",
		CountryISO2Code: "GB",
		CountryName:     "UNITED KINGDOM",
	}
)

func TestHandleGetSwiftCodeV1(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		swiftCode  string
		setup      func(pg *sqlx.DB) error
		expected   any
		wantErr    bool
		statusCode int
	}{
		{
			name:      "hq response schema",
			swiftCode: hqBank.SwiftCode,
			setup: func(pg *sqlx.DB) error {
				return db.InsertBanks(pg, []db.Bank{hqBank, branchBank1, branchBank2})
			},
			expected: GetSwiftCodeHqRes{
				Address:       hqBank.Address,
				BankName:      hqBank.BankName,
				CountryISO2:   hqBank.CountryISO2Code,
				CountryName:   hqBank.CountryName,
				IsHeadquarter: hqBank.IsHq(),
				SwiftCode:     hqBank.SwiftCode,
				Branches: []GetSwiftCodeHqBranch{
					{
						Address:       branchBank1.Address,
						BankName:      branchBank1.BankName,
						CountryISO2:   branchBank1.CountryISO2Code,
						SwiftCode:     branchBank1.SwiftCode,
						IsHeadquarter: branchBank1.IsHq(),
					},
					{
						Address:       branchBank2.Address,
						BankName:      branchBank2.BankName,
						CountryISO2:   branchBank2.CountryISO2Code,
						SwiftCode:     branchBank2.SwiftCode,
						IsHeadquarter: branchBank2.IsHq(),
					},
				},
			},
			statusCode: http.StatusOK,
		},
		{
			name:      "branch response schema",
			swiftCode: branchBank1.SwiftCode,
			setup: func(pg *sqlx.DB) error {
				return db.InsertBanks(pg, []db.Bank{hqBank, branchBank1})
			},
			expected: GetSwiftCodeBranchRes{
				Address:       branchBank1.Address,
				BankName:      branchBank1.BankName,
				CountryName:   branchBank1.CountryName,
				CountryISO2:   branchBank1.CountryISO2Code,
				IsHeadquarter: branchBank1.IsHq(),
				SwiftCode:     branchBank1.SwiftCode,
			},
			statusCode: http.StatusOK,
		},
		{
			name:       "not found",
			swiftCode:  "MISSING",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	testApi(func(args testApiArgs) {
		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				if tt.setup != nil {
					require.NoError(t, tt.setup(args.db))
				}
				t.Cleanup(func() {
					args.db.Exec("TRUNCATE bank")
				})

				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/v1/swift-codes/"+tt.swiftCode, nil)

				args.router.ServeHTTP(w, r)

				res := w.Result()
				assert.Equal(t, tt.statusCode, res.StatusCode)

				if !tt.wantErr {
					var actualResponse any
					require.NoError(t, json.NewDecoder(res.Body).Decode(&actualResponse))

					expectedJSON, err := json.Marshal(tt.expected)
					require.NoError(t, err)

					var expectedResponse any
					require.NoError(t, json.Unmarshal(expectedJSON, &expectedResponse))

					assert.Equal(t, expectedResponse, actualResponse, "response mismatch")
				}
			})
		}
	})
}

func TestHandleGetSwiftCodesForCountryV1(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		countryCode string
		statusCode  int
		setup       func(pg *sqlx.DB) error
		expected    any
		wantErr     bool
	}{
		{
			name:        "success",
			countryCode: hqBank.CountryISO2Code,
			statusCode:  http.StatusOK,
			setup: func(pg *sqlx.DB) error {
				return db.InsertBanks(pg, []db.Bank{hqBank, branchBank1})
			},
			expected: GetSwiftCodesForCountryRes{
				CountryISO2: hqBank.CountryISO2Code,
				CountryName: hqBank.CountryName,
				SwiftCodes: []GetSwiftCodesForCountrySwiftCode{
					{
						Address:       hqBank.Address,
						BankName:      hqBank.BankName,
						CountryISO2:   hqBank.CountryISO2Code,
						IsHeadquarter: hqBank.IsHq(),
						SwiftCode:     hqBank.SwiftCode,
					},
					{
						Address:       branchBank1.Address,
						BankName:      branchBank1.BankName,
						CountryISO2:   branchBank1.CountryISO2Code,
						IsHeadquarter: branchBank1.IsHq(),
						SwiftCode:     branchBank1.SwiftCode,
					},
				},
			},
		},
		{
			name:        "not found",
			countryCode: "ABC",
			statusCode:  http.StatusNotFound,
			wantErr:     true,
		},
	}

	testApi(func(args testApiArgs) {
		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				if tt.setup != nil {
					require.NoError(t, tt.setup(args.db))
				}
				t.Cleanup(func() {
					args.db.Exec("TRUNCATE bank")
				})

				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/v1/swift-codes/country/"+tt.countryCode, nil)

				args.router.ServeHTTP(w, r)

				res := w.Result()
				assert.Equal(t, tt.statusCode, res.StatusCode)

				if !tt.wantErr {
					var actualResponse any
					require.NoError(t, json.NewDecoder(res.Body).Decode(&actualResponse))

					expectedJSON, err := json.Marshal(tt.expected)
					require.NoError(t, err)

					var expectedResponse any
					require.NoError(t, json.Unmarshal(expectedJSON, &expectedResponse))

					assert.Equal(t, expectedResponse, actualResponse, "response mismatch")
				}
			})
		}
	})
}

func TestHandleAddSwiftCodeV1(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		requestBody    any
		setup          func(pg *sqlx.DB) error
		expectedStatus int
		expectedBody   any
		wantErr        bool
	}{
		{
			name: "successful HQ bank creation",
			requestBody: AddSwiftCodeReq{
				Address:       hqBank.Address,
				BankName:      hqBank.BankName,
				CountryISO2:   hqBank.CountryISO2Code,
				CountryName:   hqBank.CountryName,
				IsHeadquarter: hqBank.IsHq(),
				SwiftCode:     hqBank.SwiftCode,
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   MessageRes{Message: "Added bank with SWIFT code " + hqBank.SwiftCode},
		},
		{
			name: "successful branch bank creation",
			requestBody: AddSwiftCodeReq{
				Address:       branchBank1.Address,
				BankName:      branchBank1.BankName,
				CountryISO2:   branchBank1.CountryISO2Code,
				CountryName:   branchBank1.CountryName,
				IsHeadquarter: branchBank1.IsHq(),
				SwiftCode:     branchBank1.SwiftCode,
			},
			setup: func(pg *sqlx.DB) error {
				return db.InsertBank(pg, hqBank)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   MessageRes{Message: "Added bank with SWIFT code " + branchBank1.SwiftCode},
		},
		{
			name:           "unparsable JSON",
			requestBody:    "unparsable json",
			expectedStatus: http.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "parsable json, but not valid request body",
			requestBody: AddSwiftCodeReq{
				SwiftCode: "TESTBANK", // missing other required fields
			},
			expectedStatus: http.StatusUnprocessableEntity,
			wantErr:        true,
		},
		{
			name: "hq flag mismatch with swift code",
			requestBody: AddSwiftCodeReq{
				SwiftCode:     "HQTESTBANK", // looks like HQ code
				BankName:      "Test Bank",
				Address:       "123 Test Street",
				CountryISO2:   "US",
				CountryName:   "United States",
				IsHeadquarter: false, // but marked as not HQ
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   MessageRes{Message: "isHeadquarter disagrees with swiftCode"},
			wantErr:        true,
		},
		{
			name: "branch without existing hq",
			requestBody: AddSwiftCodeReq{
				SwiftCode:     "BRANCHBANK",
				BankName:      "Branch Bank",
				Address:       "456 Branch Street",
				CountryISO2:   "US",
				CountryName:   "United States",
				IsHeadquarter: false,
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   MessageRes{Message: "branch must have an existing HQ bank"},
			wantErr:        true,
		},
		{
			name: "duplicate swift code",
			requestBody: AddSwiftCodeReq{
				Address:       hqBank.Address,
				BankName:      hqBank.BankName,
				CountryISO2:   hqBank.CountryISO2Code,
				CountryName:   hqBank.CountryName,
				IsHeadquarter: hqBank.IsHq(),
				SwiftCode:     hqBank.SwiftCode,
			},
			setup: func(pg *sqlx.DB) error {
				// Insert the same bank first
				return db.InsertBank(pg, hqBank)
			},
			expectedStatus: http.StatusConflict,
			wantErr:        true,
		},
	}

	testApi(func(args testApiArgs) {
		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				if tt.setup != nil {
					require.NoError(t, tt.setup(args.db))
				}
				t.Cleanup(func() {
					args.db.Exec("TRUNCATE bank")
				})

				// Create request body
				var reqBody []byte
				switch v := tt.requestBody.(type) {
				case string:
					reqBody = []byte(v)
				default:
					var err error
					reqBody, err = json.Marshal(v)
					require.NoError(t, err)
				}

				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/v1/swift-codes", bytes.NewReader(reqBody))
				r.Header.Set("Content-Type", "application/json")

				args.router.ServeHTTP(w, r)

				res := w.Result()
				defer res.Body.Close()

				assert.Equal(t, tt.expectedStatus, res.StatusCode, "status code mismatch")

				if !tt.wantErr {
					var actualBody MessageRes
					require.NoError(t, json.NewDecoder(res.Body).Decode(&actualBody))
					assert.Equal(t, tt.expectedBody, actualBody)

					// Verify the bank was actually created in DB
					if tt.expectedStatus == http.StatusCreated {
						var count int
						err := args.db.Get(&count, "SELECT COUNT(*) FROM bank WHERE swift_code = $1", tt.requestBody.(AddSwiftCodeReq).SwiftCode)
						require.NoError(t, err)
						assert.Equal(t, 1, count)
					}
				}
			})
		}
	})
}

func TestHandleDeleteSwiftCodeV1(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		swiftCode      string
		setup          func(pg *sqlx.DB) error
		expectedStatus int
		expectedBody   any
		wantErr        bool
	}{
		{
			name:      "delete branch bank",
			swiftCode: hqBank.SwiftCode,
			setup: func(pg *sqlx.DB) error {
				return db.InsertBanks(pg, []db.Bank{hqBank, branchBank1, branchBank2})
			},
			expectedStatus: http.StatusOK,
			expectedBody:   MessageRes{Message: "Deleted bank with SWIFT code " + hqBank.SwiftCode},
		},
		{
			name:      "delete branch bank",
			swiftCode: branchBank1.SwiftCode,
			setup: func(pg *sqlx.DB) error {
				return db.InsertBanks(pg, []db.Bank{hqBank, branchBank1})
			},
			expectedStatus: http.StatusOK,
			expectedBody:   MessageRes{Message: "Deleted bank with SWIFT code " + branchBank1.SwiftCode},
		},
		{
			name:           "non-existent bank",
			swiftCode:      "MISSINGBANK",
			expectedStatus: http.StatusNotFound,
			wantErr:        true,
		},
	}

	testApi(func(args testApiArgs) {
		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				if tt.setup != nil {
					require.NoError(t, tt.setup(args.db))
				}
				t.Cleanup(func() {
					args.db.Exec("TRUNCATE bank")
				})

				w := httptest.NewRecorder()
				r := httptest.NewRequest("DELETE", "/v1/swift-codes/"+tt.swiftCode, nil)

				args.router.ServeHTTP(w, r)

				res := w.Result()
				defer res.Body.Close()

				assert.Equal(t, tt.expectedStatus, res.StatusCode, "status code mismatch")

				if !tt.wantErr {
					var actualBody MessageRes
					require.NoError(t, json.NewDecoder(res.Body).Decode(&actualBody))
					assert.Equal(t, tt.expectedBody, actualBody)

					// Verify the bank was actually deleted from DB
					var count int
					err := args.db.Get(&count, "SELECT COUNT(*) FROM bank WHERE swift_code = $1", tt.swiftCode)
					require.NoError(t, err)
					assert.Equal(t, 0, count)

					var hqCount int
					err = args.db.Get(&hqCount, "SELECT COUNT(*) FROM bank WHERE hq_swift_code = $1", tt.swiftCode)
					require.NoError(t, err)
					assert.Equal(t, 0, hqCount)
				}
			})
		}
	})
}
