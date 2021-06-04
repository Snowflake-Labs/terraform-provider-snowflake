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

func TestNotificationIntegration(t *testing.T) {
	r := require.New(t)
	err := resources.NotificationIntegration().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestNotificationIntegrationCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":                            "test_notification_integration",
		"comment":                         "great comment",
		"azure_storage_queue_primary_uri": "azure://great-bucket/great-path/",
		"azure_tenant_id":                 "some-guid",
	}
	d := schema.TestResourceDataRaw(t, resources.NotificationIntegration().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE NOTIFICATION INTEGRATION "test_notification_integration" AZURE_STORAGE_QUEUE_PRIMARY_URI='azure://great-bucket/great-path/' AZURE_TENANT_ID='some-guid' COMMENT='great comment' NOTIFICATION_PROVIDER='AZURE_STORAGE_QUEUE' TYPE='QUEUE' ENABLED=true$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadNotificationIntegration(mock)

		err := resources.CreateNotificationIntegration(d, db)
		r.NoError(err)
	})
}

func TestNotificationIntegrationRead(t *testing.T) {
	r := require.New(t)

	d := notificationIntegration(t, "test_notification_integration", map[string]interface{}{"name": "test_notification_integration"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadNotificationIntegration(mock)

		err := resources.ReadNotificationIntegration(d, db)
		r.NoError(err)
	})
}

func TestNotificationIntegrationDelete(t *testing.T) {
	r := require.New(t)

	d := notificationIntegration(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP NOTIFICATION INTEGRATION "drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteNotificationIntegration(d, db)
		r.NoError(err)
	})
}

func expectReadNotificationIntegration(mock sqlmock.Sqlmock) {
	showRows := sqlmock.NewRows([]string{
		"name", "type", "category", "enabled", "created_on"},
	).AddRow("test_notification_integration", "QUEUE", "NOTIFICATION", true, "now")
	mock.ExpectQuery(`^SHOW NOTIFICATION INTEGRATIONS LIKE 'test_notification_integration'$`).WillReturnRows(showRows)

	descRows := sqlmock.NewRows([]string{
		"property", "property_type", "property_value", "property_default",
	}).AddRow("ENABLED", "Boolean", true, false).
		AddRow("NOTIFICATION_PROVIDER", "String", "AZURE_STORAGE_QUEUE", nil).
		AddRow("AZURE_STORAGE_QUEUE_PRIMARY_URI", "String", "azure://great-bucket/great-path/", nil).
		AddRow("AZURE_TENANT_ID", "String", "some-guid", nil)

	mock.ExpectQuery(`DESCRIBE NOTIFICATION INTEGRATION "test_notification_integration"$`).WillReturnRows(descRows)
}
