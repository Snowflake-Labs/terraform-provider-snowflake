package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestRoleGrants(t *testing.T) {
	r := require.New(t)
	err := resources.RoleGrants().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestRoleGrantsCreate(t *testing.T) {
	r := require.New(t)

	d := roleGrants(t, "good_name", map[string]interface{}{
		"role_name": "good_name",
		"roles":     []interface{}{"role1", "role2"},
		"users":     []interface{}{"user1", "user2"},
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`GRANT ROLE "good_name" TO ROLE "role2"`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`GRANT ROLE "good_name" TO ROLE "role1"`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`GRANT ROLE "good_name" TO USER "user1"`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`GRANT ROLE "good_name" TO USER "user2"`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadRoleGrants(mock)
		err := resources.CreateRoleGrants(d, db)
		r.NoError(err)
	})
}

func expectReadRoleGrants(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on",
		"role",
		"granted_to",
		"grantee_name",
		"granted_by",
	}).
		AddRow("_", "good_name", "ROLE", "role1", "").
		AddRow("_", "good_name", "ROLE", "role2", "").
		AddRow("_", "good_name", "USER", "user1", "").
		AddRow("_", "good_name", "USER", "user2", "")
	mock.ExpectQuery(`SHOW GRANTS OF ROLE "good_name"`).WillReturnRows(rows)
}

func TestRoleGrantsRead(t *testing.T) {
	r := require.New(t)

	d := roleGrants(t, "good_name||||role1,role2|false", map[string]interface{}{
		"role_name": "good_name",
		"roles":     []interface{}{"role1", "role2"},
		"users":     []interface{}{"user1", "user2"},
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadRoleGrants(mock)
		err := resources.ReadRoleGrants(d, db)
		r.NoError(err)
		r.Len(d.Get("users").(*schema.Set).List(), 2)
		r.Len(d.Get("roles").(*schema.Set).List(), 2)
	})
}

func TestRoleGrantsDelete(t *testing.T) {
	r := require.New(t)

	d := roleGrants(t, "drop_it||||role1,role2|false", map[string]interface{}{
		"role_name": "drop_it",
		"roles":     []interface{}{"role1", "role2"},
		"users":     []interface{}{"user1", "user2"},
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {

		mock.ExpectExec(`REVOKE ROLE "drop_it" FROM ROLE "role1"`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`REVOKE ROLE "drop_it" FROM ROLE "role2"`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`REVOKE ROLE "drop_it" FROM USER "user1"`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`REVOKE ROLE "drop_it" FROM USER "user2"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteRoleGrants(d, db)
		r.NoError(err)
	})
}
