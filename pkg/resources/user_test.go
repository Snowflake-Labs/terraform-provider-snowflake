package resources_test

import (
	"database/sql"
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestUser(t *testing.T) {
	resources.User().InternalValidate(provider.Provider().Schema, false)
}

func TestUserCreate(t *testing.T) {
	a := assert.New(t)

	in := map[string]interface{}{
		"name":     "good_name",
		"comment":  "great comment",
		"password": "awesomepassword",
	}
	d := schema.TestResourceDataRaw(t, resources.User().Schema, in)
	a.NotNil(d)

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec("CREATE USER good_name COMMENT='great comment' PASSWORD='awesomepassword'").WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadUser(mock)
		err := resources.CreateUser(d, db)
		a.NoError(err)
	})
}

func expectReadUser(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"name", "created_on", "login_name", "display_name", "first_name", "last_name", "email", "mins_to_unlock",
		"days_to_expiry", "comment", "disabled", "must_change_password", "snowflake_lock", "default_warehouse",
		"default_namespace", "default_role", "ext_authn_duo", "ext_authn_uid", "mins_to_bypass_mfa", "owner",
		"last_success_login", "expires_at_time", "locked_until_time", "has_password", "has_rsa_public_key"},
	).AddRow("good_name", "created_on", "login_name", "display_name", "first_name", "last_name", "email", "mins_to_unlock", "days_to_expiry", "mock comment", "disabled", "must_change_password", "snowflake_lock", "default_warehouse", "default_namespace", "default_role", "ext_authn_duo", "ext_authn_uid", "mins_to_bypass_mfa", "owner", "last_success_login", "expires_at_time", "locked_until_time", "has_password", "has_rsa_public_key")
	mock.ExpectQuery(`SHOW USERS LIKE 'good_name'`).WillReturnRows(rows)
}

func TestUserRead(t *testing.T) {
	a := assert.New(t)

	d := user(t, "good_name", map[string]interface{}{"name": "good_name"})

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadUser(mock)
		err := resources.ReadUser(d, db)
		a.NoError(err)
		a.Equal("mock comment", d.Get("comment").(string))
	})
}

func TestUserDelete(t *testing.T) {
	a := assert.New(t)

	d := user(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec("DROP USER drop_it").WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteUser(d, db)
		a.NoError(err)
	})
}
