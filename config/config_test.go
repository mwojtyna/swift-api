package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnv(t *testing.T) {
	var tests = []struct {
		dbUser        string
		dbPass        string
		dbName        string
		apiPort       string
		errorExpected bool
	}{
		{"user", "password", "name", "1234", false},
		{"user", "password", "", "1234", true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			t.Setenv("DB_USER", tt.dbUser)
			t.Setenv("DB_PASS", tt.dbPass)
			t.Setenv("DB_NAME", tt.dbName)
			t.Setenv("API_PORT", tt.apiPort)

			env, err := LoadEnv()
			if tt.errorExpected {
				assert.Error(t, err)
			} else {
				if !assert.NoError(t, err) {
					t.FailNow()
				}

				assert.Equal(t, tt.dbUser, env.DB_USER)
				assert.Equal(t, tt.dbPass, env.DB_PASS)
				assert.Equal(t, tt.dbName, env.DB_NAME)
				assert.Equal(t, tt.apiPort, env.API_PORT)
			}
		})
	}
}
