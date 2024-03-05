package resources_test

import (
	"database/sql"
	"testing"

	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestRoleOwnershipGrant(t *testing.T) {
	r := require.New(t)
	err := resources.RoleOwnershipGrant().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestRoleOwnershipGrantCreate(t *testing.T) {
	r := require.New(t)

	d := roleOwnershipGrant(t, "good_name", map[string]interface{}{
		"on_role_name":   "good_name",
		"to_role_name":   "other_good_name",
		"current_grants": "COPY",
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`GRANT OWNERSHIP ON ROLE "good_name" TO ROLE "other_good_name" COPY CURRENT GRANTS`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadRoleOwnershipGrant(mock)
		err := resources.CreateRoleOwnershipGrant(d, &internalprovider.Context{
			Client: sdk.NewClientFromDB(db),
		})
		r.NoError(err)
	})
}

func TestRoleOwnershipGrantRead(t *testing.T) {
	r := require.New(t)

	d := roleOwnershipGrant(t, "good_name|other_good_name|COPY", map[string]interface{}{
		"on_role_name":   "good_name",
		"to_role_name":   "other_good_name",
		"current_grants": "COPY",
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadRoleOwnershipGrant(mock)
		err := resources.ReadRoleOwnershipGrant(d, &internalprovider.Context{
			Client: sdk.NewClientFromDB(db),
		})
		r.NoError(err)
	})
}

func expectReadRoleOwnershipGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on",
		"name",
		"is_default",
		"is_current",
		"is_inherited",
		"assigned_to_users",
		"granted_to_roles",
		"granted_roles",
		"owner",
		"comment",
	}).AddRow("_", "good_name", "", "", "", "", "", "", "other_good_name", "")
	mock.ExpectQuery(`SHOW ROLES LIKE 'good_name'`).WillReturnRows(rows)
}

func TestRoleOwnershipGrantDelete(t *testing.T) {
	r := require.New(t)

	d := roleOwnershipGrant(t, "good_name|other_good_name|COPY", map[string]interface{}{
		"on_role_name":   "good_name",
		"to_role_name":   "other_good_name",
		"current_grants": "COPY",
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`GRANT OWNERSHIP ON ROLE "good_name" TO ROLE "ACCOUNTADMIN" COPY CURRENT GRANTS`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteRoleOwnershipGrant(d, &internalprovider.Context{
			Client: sdk.NewClientFromDB(db),
		})
		r.NoError(err)
	})
}
