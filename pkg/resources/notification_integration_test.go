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
		"name":                         "test_notification_integration",
		"comment":                      "yet another great comment",
		"notification_provider":        "GCP_PUBSUB",
		"gcp_pubsub_subscription_name": "projects/github-sh1n/subscriptions/hire-me",
	}

	d := schema.TestResourceDataRaw(t, resources.NotificationIntegration().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE NOTIFICATION INTEGRATION "test_notification_integration" COMMENT='yet another great comment' GCP_PUBSUB_SUBSCRIPTION_NAME='projects/github-sh1n/subscriptions/hire-me' NOTIFICATION_PROVIDER='GCP_PUBSUB' TYPE='QUEUE' ENABLED=true$`,
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
	).AddRow("test_notification_integration", "QUEUE - GCP_PUBSUB", "NOTIFICATION", true, "now")
	mock.ExpectQuery(`^SHOW NOTIFICATION INTEGRATIONS LIKE 'test_notification_integration'$`).WillReturnRows(showRows)

	descRows := sqlmock.NewRows([]string{
		"property", "property_type", "property_value", "property_default",
	}).AddRow("ENABLED", "Boolean", true, false).
		AddRow("GCP_PUBSUB_SUBSCRIPTION_NAME", "String", "projects/github-sh1n/subscriptions/hire-me", nil).
		AddRow("GCP_PUBSUB_SERVICE_ACCOUNT", "String", "random@region-something.iam.google.gcp", nil)

	mock.ExpectQuery(`DESCRIBE NOTIFICATION INTEGRATION "test_notification_integration"$`).WillReturnRows(descRows)
}
