package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

//lintignore:AT003
func TestAccountGrant(t *testing.T) {
	r := require.New(t)
	err := resources.AccountGrant().Resource.InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

//lintignore:AT003
func TestAccountGrantCreate(t *testing.T) { //lintignore:AT003
	r := require.New(t)

	in := map[string]interface{}{
		"privilege":         "CREATE DATABASE",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": true,
	}
	d := schema.TestResourceDataRaw(t, resources.AccountGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT CREATE DATABASE ON ACCOUNT TO ROLE "test-role-1" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT CREATE DATABASE ON ACCOUNT TO ROLE "test-role-2" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadAccountGrant(mock)
		err := resources.CreateAccountGrant(d, db)
		r.NoError(err)
	})
}

//lintignore:AT003
func TestAccountGrantRead(t *testing.T) {
	r := require.New(t)

	d := accountGrant(t, "ACCOUNT|||MANAGE GRANTS||true", map[string]interface{}{
		"privilege":         "MANAGE GRANTS",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": true,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadAccountGrant(mock)
		err := resources.ReadAccountGrant(d, db)
		r.NoError(err)
	})
}

func TestMonitorExecution(t *testing.T) {
	r := require.New(t)

	d := accountGrant(t, "ACCOUNT|||MONITOR EXECUTION||true", map[string]interface{}{
		"privilege":         "MONITOR EXECUTION",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": true,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadAccountGrant(mock)
		err := resources.ReadAccountGrant(d, db)
		r.NoError(err)
	})
}

func TestExecuteTask(t *testing.T) {
	r := require.New(t)

	d := accountGrant(t, "ACCOUNT|||EXECUTE TASK||false", map[string]interface{}{
		"privilege":         "EXECUTE TASK",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": false,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadAccountGrant(mock)
		err := resources.ReadAccountGrant(d, db)
		r.NoError(err)
	})
}

func expectReadAccountGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "MANAGE GRANTS", "ACCOUNT", "", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "MANAGE GRANTS", "ACCOUNT", "", "ROLE", "test-role-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON ACCOUNT$`).WillReturnRows(rows)
}

func TestApplyMaskingPolicy(t *testing.T) {
	r := require.New(t)

	d := accountGrant(t, "ACCOUNT|||APPLY MASKING POLICY||true", map[string]interface{}{
		"privilege":         "APPLY MASKING POLICY",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": true,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadAccountGrant(mock)
		err := resources.ReadAccountGrant(d, db)
		r.NoError(err)
	})
}

func expectApplyMaskingPolicy(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "APPLY MASKING POLICY", "ACCOUNT", "", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "APPLY MASKING POLICY", "ACCOUNT", "", "ROLE", "test-role-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON ACCOUNT$`).WillReturnRows(rows)
}
