package db

import (
	"testing"

	"github.com/mwojtyna/swift-api/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	t.Parallel()
	utils.TestWithPostgres(func(args utils.TestWithPostgresArgs) {
		t.Run("successful connection with correct credentials", func(t *testing.T) {
			db, err := Connect(args.Env.DB_USER, args.Env.DB_PASS, args.Env.DB_NAME, args.Port)
			require.NoError(t, err)

			t.Cleanup(func() {
				err := db.Close()
				require.NoError(t, err)
			})
		})
		t.Run("connection fails with wrong credentials", func(t *testing.T) {
			db, err := Connect("wronguser", "wrongpass", "wrongdb", "-1")
			assert.Error(t, err)
			assert.Nil(t, db)
		})
	})
}

func TestIsEmpty(t *testing.T) {
	t.Parallel()
	utils.TestWithPostgres(func(args utils.TestWithPostgresArgs) {
		db, err := Connect(args.Env.DB_USER, args.Env.DB_PASS, args.Env.DB_NAME, args.Port)
		require.NoError(t, err)

		t.Run("empty db returns true", func(t *testing.T) {
			empty, err := IsEmpty(db)
			require.NoError(t, err)
			assert.True(t, empty)
		})
		t.Run("non-empty db returns false", func(t *testing.T) {
			// Arrange
			_, err := db.NamedExec(`INSERT INTO bank (swift_code, hq_swift_code, bank_name, address, country_iso2_code, country_name) 
		VALUES (:swift_code, :hq_swift_code, :bank_name, :address, :country_iso2_code, :country_name);`, Bank{})
			require.NoError(t, err)

			// Act
			empty, err := IsEmpty(db)

			// Assert
			require.NoError(t, err)
			assert.False(t, empty)

			t.Cleanup(func() {
				_, err := db.Exec("TRUNCATE TABLE bank")
				require.NoError(t, err)
			})
		})

		t.Cleanup(func() {
			err := db.Close()
			require.NoError(t, err)
		})
	})
}
