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
		apiPort string
		wantErr bool
	}{
		{"correct", "user", "password", "name", "1234", false},
		{"variable not set", "user", "password", "", "1234", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DB_USER", tt.dbUser)
			t.Setenv("DB_PASS", tt.dbPass)
			t.Setenv("DB_NAME", tt.dbName)
			t.Setenv("API_PORT", tt.apiPort)

			got, err := LoadEnv()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.dbUser, got.DB_USER)
				assert.Equal(t, tt.dbPass, got.DB_PASS)
				assert.Equal(t, tt.dbName, got.DB_NAME)
				assert.Equal(t, tt.apiPort, got.API_PORT)
			}
		})
	}
}
