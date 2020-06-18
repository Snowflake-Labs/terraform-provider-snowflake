package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
)

func TestAccountGrant(t *testing.T) {
	r := require.New(t)
	err := resources.AccountGrant().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestAccountGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"privilege": "CREATE DATABASE",
		"roles":     []interface{}{"test-role-1", "test-role-2"},
	}
	d := schema.TestResourceDataRaw(t, resources.AccountGrant().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT CREATE DATABASE ON ACCOUNT TO ROLE "test-role-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT CREATE DATABASE ON ACCOUNT TO ROLE "test-role-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadAccountGrant(mock)
		err := resources.CreateAccountGrant(d, db)
		r.NoError(err)
	})
}

func TestAccountGrantRead(t *testing.T) {
	r := require.New(t)

	d := accountGrant(t, "ACCOUNT|||MANAGE GRANTS", map[string]interface{}{
		"privilege": "MANAGE GRANTS",
		"roles":     []interface{}{"test-role-1", "test-role-2"},
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

	d := accountGrant(t, "ACCOUNT|||MONITOR EXECUTION", map[string]interface{}{
		"privilege": "MONITOR EXECUTION",
		"roles":     []interface{}{"test-role-1", "test-role-2"},
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

	d := accountGrant(t, "ACCOUNT|||EXECUTE TASK", map[string]interface{}{
		"privilege": "EXECUTE TASK",
		"roles":     []interface{}{"test-role-1", "test-role-2"},
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
