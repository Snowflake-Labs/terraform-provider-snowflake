package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
)

func TestResourceMonitor(t *testing.T) {
	r := require.New(t)
	err := resources.ResourceMonitor().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestResourceMonitorCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":                      "good_name",
		"notify_users":              []interface{}{"USERONE", "USERTWO"},
		"credit_quota":              100,
		"notify_triggers":           []interface{}{75, 88},
		"suspend_trigger":           99,
		"suspend_immediate_trigger": 105,
		"set_for_account":           true,
	}

	d := schema.TestResourceDataRaw(t, resources.ResourceMonitor().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE RESOURCE MONITOR "good_name" CREDIT_QUOTA=100 NOTIFY_USERS=\('USERTWO', 'USERONE'\) TRIGGERS ON 99 PERCENT DO SUSPEND ON 105 PERCENT DO SUSPEND_IMMEDIATE ON 88 PERCENT DO NOTIFY ON 75 PERCENT DO NOTIFY$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^ALTER ACCOUNT SET RESOURCE_MONITOR = "good_name"$`).WillReturnResult(sqlmock.NewResult(1, 1))

		expectReadResourceMonitor(mock)
		err := resources.CreateResourceMonitor(d, db)
		r.NoError(err)
	})
}

func expectReadResourceMonitor(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"name", "credit_quota", "used_credits", "remaining_credits", "level",
		"frequency", "start_time", "end_time", "notify_at", "suspend_at",
		"suspend_immediately_at", "created_on", "owner", "comment", "notify_users",
	}).AddRow(
		"good_name", 100.00, 0.00, 100.00, "ACCOUNT", "MONTHLY", "2001-01-01 00:00:00.000 -0700",
		"", "75%,88%", "99%", "105%", "2001-01-01 00:00:00.000 -0700", "ACCOUNTADMIN", "", "USERONE, USERTWO")
	mock.ExpectQuery(`^SHOW RESOURCE MONITORS LIKE 'good_name'$`).WillReturnRows(rows)
}

func TestResourceMonitorDelete(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name": "good_name",
	}

	d := schema.TestResourceDataRaw(t, resources.ResourceMonitor().Schema, in)
	d.SetId("good_name")

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^DROP RESOURCE MONITOR "good_name"$`).WillReturnResult(sqlmock.NewResult(1, 1))

		err := resources.DeleteResourceMonitor(d, db)
		r.NoError(err)
	})
}

func TestResourceMonitorRead(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name": "good_name",
	}

	d := resourceMonitor(t, "good_name", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		// Test when resource is not found, checking if state will be empty
		r.NotEmpty(d.State())
		q := snowflake.NewResourceMonitorBuilder(d.Id()).Show()
		mock.ExpectQuery(q).WillReturnError(sql.ErrNoRows)
		err := resources.ReadResourceMonitor(d, db)
		r.Empty(d.State())
		r.Nil(err)
	})
}
