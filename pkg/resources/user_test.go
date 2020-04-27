package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	r := require.New(t)
	err := resources.User().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestUserCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":                 "good_name",
		"comment":              "great comment",
		"password":             "awesomepassword",
		"login_name":           "gname",
		"disabled":             true,
		"default_warehouse":    "mywarehouse",
		"default_namespace":    "mynamespace",
		"default_role":         "bestrole",
		"rsa_public_key":       "asdf",
		"rsa_public_key_2":     "asdf2",
		"must_change_password": true,
	}
	d := schema.TestResourceDataRaw(t, resources.User().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^CREATE USER "good_name" COMMENT='great comment' DEFAULT_NAMESPACE='mynamespace' DEFAULT_ROLE='bestrole' DEFAULT_WAREHOUSE='mywarehouse' LOGIN_NAME='gname' PASSWORD='awesomepassword' RSA_PUBLIC_KEY='asdf' RSA_PUBLIC_KEY_2='asdf2' DISABLED=true MUST_CHANGE_PASSWORD=true$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadUser(mock)
		err := resources.CreateUser(d, db)
		r.NoError(err)
	})
}

func expectReadUser(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"name", "created_on", "login_name", "display_name", "first_name", "last_name", "email", "mins_to_unlock",
		"days_to_expiry", "comment", "disabled", "must_change_password", "snowflake_lock", "default_warehouse",
		"default_namespace", "default_role", "ext_authn_duo", "ext_authn_uid", "mins_to_bypass_mfa", "owner",
		"last_success_login", "expires_at_time", "locked_until_time", "has_password", "has_rsa_public_key"},
	).AddRow("good_name", "created_on", "myloginname", "display_name", "first_name", "last_name", "email", "mins_to_unlock", "days_to_expiry", "mock comment", false, true, "snowflake_lock", "default_warehouse", "default_namespace", "default_role", "ext_authn_duo", "ext_authn_uid", "mins_to_bypass_mfa", "owner", "last_success_login", "expires_at_time", "locked_until_time", "has_password", false)
	mock.ExpectQuery(`^SHOW USERS LIKE 'good_name'$`).WillReturnRows(rows)
}

func TestUserRead(t *testing.T) {
	r := require.New(t)

	d := user(t, "good_name", map[string]interface{}{"name": "good_name"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadUser(mock)
		err := resources.ReadUser(d, db)
		r.NoError(err)
		r.Equal("mock comment", d.Get("comment").(string))
		r.Equal("myloginname", d.Get("login_name").(string))
		r.Equal(false, d.Get("disabled").(bool))
	})
}

func TestUserExists(t *testing.T) {
	r := require.New(t)

	d := user(t, "good_name", map[string]interface{}{"name": "good_name"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadUser(mock)
		b, err := resources.UserExists(d, db)
		r.NoError(err)
		r.True(b)
	})
}

func TestUserDelete(t *testing.T) {
	r := require.New(t)

	d := user(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^DROP USER "drop_it"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteUser(d, db)
		r.NoError(err)
	})
}
