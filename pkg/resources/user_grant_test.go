package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUserGrant(t *testing.T) {
	r := require.New(t)
	err := resources.UserGrant().Resource.InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestUserGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"user_name": "test-user",
		"privilege": "MONITOR",
		"roles":     []interface{}{"test-role-1", "test-role-2"},
	}
	d := schema.TestResourceDataRaw(t, resources.UserGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT MONITOR ON USER "test-user" TO ROLE "test-role-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT MONITOR ON USER "test-user" TO ROLE "test-role-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadUserGrant(mock)
		err := resources.CreateUserGrant(d, db)
		r.NoError(err)
	})
}

func expectReadUserGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "MONITOR", "USER", "test-user", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "MONITOR", "USER", "test-user", "ROLE", "test-role-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON USER "test-user"$`).WillReturnRows(rows)
}
