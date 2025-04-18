package parser

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/mwojtyna/swift-api/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestParseCsv(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []db.Bank
		wantErr bool
	}{
		{
			name: "valid with hq and branch",
			input: `COUNTRY,SWIFT CODE,CODE TYPE,BANK NAME,BANK ADDRESS,TOWN NAME,COUNTRY NAME,TIME ZONE
PL,BPHKPLPKXXX,BIC11,BANK BPH SA,"UL. CYPRIANA KAMILA NORWIDA 1  GDANSK, POMORSKIE, 80-280",GDANSK,POLAND,Europe/Warsaw
PL,BPHKPLPKCUS,BIC11,BANK BPH SA,"UL. CYPRIANA KAMILA NORWIDA 1  GDANSK, POMORSKIE, 80-280",GDANSK,POLAND,Europe/Warsaw`,
			want: []db.Bank{
				{
					SwiftCode:       "BPHKPLPKXXX",
					HqSwiftCode:     sql.NullString{},
					IsHeadquarter:   true,
					BankName:        "BANK BPH SA",
					Address:         "UL. CYPRIANA KAMILA NORWIDA 1  GDANSK, POMORSKIE, 80-280",
					CountryISO2Code: "PL",
					CountryName:     "POLAND",
				},
				{
					SwiftCode:       "BPHKPLPKCUS",
					HqSwiftCode:     sql.NullString{String: "BPHKPLPKXXX", Valid: true},
					IsHeadquarter:   false,
					BankName:        "BANK BPH SA",
					Address:         "UL. CYPRIANA KAMILA NORWIDA 1  GDANSK, POMORSKIE, 80-280",
					CountryISO2Code: "PL",
					CountryName:     "POLAND",
				},
			},
		},
		{
			name: "valid with empty address",
			input: `COUNTRY,SWIFT CODE,CODE TYPE,BANK NAME,BANK ADDRESS,TOWN NAME,COUNTRY NAME,TIME ZONE
PL,BPHKPLPKXXX,BIC11,BANK BPH SA,,GDANSK,POLAND,Europe/Warsaw`,
			want: []db.Bank{
				{
					SwiftCode:       "BPHKPLPKXXX",
					HqSwiftCode:     sql.NullString{},
					IsHeadquarter:   true,
					BankName:        "BANK BPH SA",
					Address:         "GDANSK",
					CountryISO2Code: "PL",
					CountryName:     "POLAND",
				},
			},
		},
		{
			name: "invalid country code",
			input: `COUNTRY,SWIFT CODE,CODE TYPE,BANK NAME,BANK ADDRESS,TOWN NAME,COUNTRY NAME,TIME ZONE
INVALIDCOUNTRYCODE,BPHKPLPKXXX,BIC11,BANK BPH SA,"UL. CYPRIANA KAMILA NORWIDA 1  GDANSK, POMORSKIE, 80-280",GDANSK,POLAND,Europe/Warsaw`,
			wantErr: true,
		},
		{
			name: "invalid swift code",
			input: `COUNTRY,SWIFT CODE,CODE TYPE,BANK NAME,BANK ADDRESS,TOWN NAME,COUNTRY NAME,TIME ZONE
PL,SWIFTINVALID,BIC11,BANK BPH SA,"UL. CYPRIANA KAMILA NORWIDA 1  GDANSK, POMORSKIE, 80-280",GDANSK,POLAND,Europe/Warsaw`,
			wantErr: true,
		},
		{
			name: "invalid number of columns",
			input: `COUNTRY,SWIFT CODE,CODE TYPE,BANK NAME,BANK ADDRESS,TOWN NAME,COUNTRY NAME
PL,BPHKPLPKXXX,BIC11,BANK BPH SA,"UL. CYPRIANA KAMILA NORWIDA 1  GDANSK, POMORSKIE, 80-280",GDANSK,POLAND,Europe/Warsaw,Additional column`,
			wantErr: true,
		},
		{
			name: "branch without hq, also trims whitespace",
			input: `COUNTRY,SWIFT CODE,CODE TYPE,BANK NAME,BANK ADDRESS,TOWN NAME,COUNTRY NAME,TIME ZONE
PL,ALBPPLP1BMW,BIC11,ALIOR BANK SPOLKA AKCYJNA,"  WARSZAWA, MAZOWIECKIE",WARSZAWA,POLAND,Europe/Warsaw`,
			want: []db.Bank{
				{
					SwiftCode:       "ALBPPLP1BMW",
					HqSwiftCode:     sql.NullString{},
					IsHeadquarter:   false,
					BankName:        "ALIOR BANK SPOLKA AKCYJNA",
					Address:         "WARSZAWA, MAZOWIECKIE",
					CountryISO2Code: "PL",
					CountryName:     "POLAND",
				},
			},
		},
		{
			name:    "invalid csv",
			input:   "not a csv file\nnot a csv file",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			got, err := ParseCsv(r)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestIsSwiftCodeHq(t *testing.T) {
	var tests = []struct {
		name   string
		code   string
		isHq   bool
		hqCode string
	}{
		{"hq code", "ABCDEFGHXXX", true, ""},
		{"not hq code", "ABCDEFGHIJKL", false, "ABCDEFGHXXX"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isHq, hqCode := IsSwiftCodeHq(tt.code)
			assert.Equal(t, tt.isHq, isHq)
			assert.Equal(t, tt.hqCode, hqCode)
		})
	}
}
