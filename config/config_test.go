package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadEnv(t *testing.T) {
	tests := []struct {
		name    string
		dbUser  string
		dbPass  string
		dbName  string
		dbHost  string
		apiPort string
		wantErr bool
	}{
		{
			name:    "correct",
			dbUser:  "user",
			dbPass:  "password",
			dbName:  "name",
			dbHost:  "host",
			apiPort: "5678",
			wantErr: false,
		},
		{
			name:    "variable not set",
			dbUser:  "user",
			dbPass:  "password",
			dbName:  "",
			dbHost:  "host",
			apiPort: "5678",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DB_USER", tt.dbUser)
			t.Setenv("DB_PASS", tt.dbPass)
			t.Setenv("DB_NAME", tt.dbName)
			t.Setenv("DB_HOST", tt.dbHost)
			t.Setenv("API_PORT", tt.apiPort)

			got, err := LoadEnv()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.dbUser, got.DB_USER)
				assert.Equal(t, tt.dbPass, got.DB_PASS)
				assert.Equal(t, tt.dbName, got.DB_NAME)
				assert.Equal(t, tt.dbHost, got.DB_HOST)
				assert.Equal(t, tt.apiPort, got.API_PORT)
			}
		})
	}
}
