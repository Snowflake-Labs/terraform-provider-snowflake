package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestUserOwnershipGrant(t *testing.T) {
	r := require.New(t)
	err := resources.UserOwnershipGrant().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestUserOwnershipGrantCreate(t *testing.T) {
	r := require.New(t)

	d := userOwnershipGrant(t, "user1", map[string]interface{}{
		"on_user_name":   "user1",
		"to_role_name":   "role1",
		"current_grants": "COPY",
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`GRANT OWNERSHIP ON USER "user1" TO ROLE "role1" COPY CURRENT GRANTS`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadUserOwnershipGrant(mock)
		err := resources.CreateUserOwnershipGrant(d, db)
		r.NoError(err)
	})
}

func TestUserOwnershipGrantRead(t *testing.T) {
	r := require.New(t)

	d := userOwnershipGrant(t, "user1|role1|COPY", map[string]interface{}{
		"on_user_name":   "user1",
		"to_role_name":   "role1",
		"current_grants": "COPY",
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadUserOwnershipGrant(mock)
		err := resources.ReadUserOwnershipGrant(d, db)
		r.NoError(err)
	})
}

func expectReadUserOwnershipGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"name",
		"created_on",
		"login_name",
		"display_name",
		"first_name",
		"last_name",
		"email",
		"mins_to_unlock",
		"days_to_expiry",
		"comment",
		"disabled",
		"must_change_password",
		"snowflake_lock",
		"default_warehouse",
		"default_namespace",
		"default_role",
		"default_secondary_roles",
		"ext_authn_duo",
		"ext_authn_uid",
		"mins_to_bypass_mfa",
		"owner",
		"last_success_login",
		"expires_at_time",
		"locked_until_time",
		"has_password",
		"has_rsa_public_key",
	}).AddRow("user1", "_", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "role1", "", "", "", "", "")
	mock.ExpectQuery(`SHOW USERS LIKE 'user1'`).WillReturnRows(rows)
}

func TestUserOwnershipGrantDelete(t *testing.T) {
	r := require.New(t)

	d := userOwnershipGrant(t, "user1|role1|COPY", map[string]interface{}{
		"on_user_name":   "user1",
		"to_role_name":   "role1",
		"current_grants": "COPY",
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`GRANT OWNERSHIP ON USER "user1" TO ROLE "ACCOUNTADMIN" COPY CURRENT GRANTS`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteUserOwnershipGrant(d, db)
		r.NoError(err)
	})
}
