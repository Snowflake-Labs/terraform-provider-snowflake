package testhelpers

import (
	"database/sql"
	"os"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func WithMockDb(t *testing.T, f func(*sql.DB, sqlmock.Sqlmock)) {
	t.Helper()
	r := require.New(t)
	db, mock, err := sqlmock.New()
	r.NoError(err)
	defer db.Close()

	// Because we are using TypeSet not TypeList, order is non-deterministic.
	mock.MatchExpectationsInOrder(false)

	f(db, mock)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFile(t *testing.T, filename string, data []byte) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), filename)
	require.NoError(t, err)

	err = os.WriteFile(f.Name(), data, 0o600)
	require.NoError(t, err)
	return f.Name()
}
