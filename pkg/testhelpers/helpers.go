package testhelpers

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func WithMockDb(t *testing.T, f func(*sql.DB, sqlmock.Sqlmock)) {
	a := assert.New(t)
	db, mock, err := sqlmock.New()
	defer db.Close()
	a.NoError(err)

	// Because we are using TypeSet not TypeList, order is non-deterministic.
	mock.MatchExpectationsInOrder(false)

	f(db, mock)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
