package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
)

func TestResourceMonitorGrant(t *testing.T) {
	r := require.New(t)
	err := resources.ResourceMonitorGrant().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestResourceMonitorGrantCreate(t *testing.T) {
	a := assert.New(t)

	in := map[string]interface{}{
		"monitor_name": "test-monitor",
		"privilege":    "MONITOR",
		"roles":        []interface{}{"test-role-1", "test-role-2"},
	}
	d := schema.TestResourceDataRaw(t, resources.ResourceMonitorGrant().Schema, in)
	a.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT MONITOR ON RESOURCE MONITOR "test-monitor" TO ROLE "test-role-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT MONITOR ON RESOURCE MONITOR "test-monitor" TO ROLE "test-role-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadResourceMonitorGrant(mock)
		err := resources.CreateResourceMonitorGrant(d, db)
		a.NoError(err)
	})
}

func TestResourceMonitorGrantRead(t *testing.T) {
	a := assert.New(t)

	d := resourceMonitorGrant(t, "test-monitor|||MONITOR", map[string]interface{}{
		"monitor_name": "test-monitor",
		"privilege":    "MONITOR",
		"roles":        []interface{}{"test-role-1", "test-role-2"},
	})

	a.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadResourceMonitorGrant(mock)
		err := resources.ReadResourceMonitorGrant(d, db)
		a.NoError(err)
	})
}

func expectReadResourceMonitorGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "MONITOR", "RESOURCE MONITOR", "test-monitor", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "MONITOR", "RESOURCE MONITOR", "test-monitor", "ROLE", "test-role-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON RESOURCE MONITOR "test-monitor"$`).WillReturnRows(rows)
}
