package testhelpers

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func WithMockDb(t *testing.T, f func(*sql.DB, sqlmock.Sqlmock)) {
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

func Providers() map[string]*schema.Provider {
	p := provider.Provider()
	return map[string]*schema.Provider{
		"snowflake": p,
	}
}
