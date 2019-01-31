package resources_test

import (
	"database/sql"
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestRoleGrants(t *testing.T) {
	resources.RoleGrants().InternalValidate(provider.Provider().Schema, false)
}

func TestRoleGrantsCreate(t *testing.T) {
	a := assert.New(t)

	d := roleGrants(t, "good_name", map[string]interface{}{
		"name":      "fake name",
		"role_name": "good_name",
		"roles":     []string{"role1", "role2"},
		"users":     []string{"user1", "user2"},
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec("GRANT ROLE good_name TO ROLE role2").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("GRANT ROLE good_name TO ROLE role1").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("GRANT ROLE good_name TO USER user1").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("GRANT ROLE good_name TO USER user2").WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadRoleGrants(mock)
		err := resources.CreateRoleGrants(d, db)
		a.NoError(err)
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
	mock.ExpectQuery(`SHOW GRANTS OF ROLE good_name`).WillReturnRows(rows)
}

func TestRoleGrantsRead(t *testing.T) {
	a := assert.New(t)

	d := roleGrants(t, "good_name", map[string]interface{}{
		"name":      "fake name",
		"role_name": "good_name",
		"roles":     []string{"role1", "role2"},
		"users":     []string{"user1", "user2"},
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadRoleGrants(mock)
		err := resources.ReadRoleGrants(d, db)
		a.NoError(err)
		a.Len(d.Get("users").(*schema.Set).List(), 2)
		a.Len(d.Get("roles").(*schema.Set).List(), 2)
	})
}

func TestRoleGrantsDelete(t *testing.T) {
	a := assert.New(t)

	d := roleGrants(t, "drop_it", map[string]interface{}{
		"name":      "drop_it",
		"role_name": "drop_it",
		"roles":     []string{"role1", "role2"},
		"users":     []string{"user1", "user2"},
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {

		mock.ExpectExec("REVOKE ROLE drop_it FROM ROLE role1").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("REVOKE ROLE drop_it FROM ROLE role2").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("REVOKE ROLE drop_it FROM USER user1").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("REVOKE ROLE drop_it FROM USER user2").WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteRoleGrants(d, db)
		a.NoError(err)
	})
}
