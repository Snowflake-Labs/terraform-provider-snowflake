package snowflake_test

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestCurrentRoleSelect(t *testing.T) {
	r := require.New(t)
	r.Equal(`SELECT CURRENT_ROLE() AS "currentRole";`, snowflake.SelectCurrentRole())
}

func TestCurrentRoleRead(t *testing.T) {
	type testCaseEntry struct {
		currentRole string
	}

	testCases := map[string]testCaseEntry{
		"name": {
			"sys_admin",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			r := require.New(t)
			mockDB, mock, err := sqlmock.New()
			r.NoError(err)
			defer mockDB.Close()
			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

			rows := sqlmock.NewRows([]string{"currentRole"}).AddRow(testCase.currentRole)
			mock.ExpectQuery(`SELECT CURRENT_ROLE\(\) AS "currentRole";`).WillReturnRows(rows)

			acc, err := snowflake.ReadCurrentRole(sqlxDB.DB)
			r.NoError(err)
			r.Equal(testCase.currentRole, acc.Role)
		})
	}
}
