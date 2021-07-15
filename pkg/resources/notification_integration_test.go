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
	testCases := []struct {
		notificationProvider string
		raw                  map[string]interface{}
		expectSQL            string
	}{
		{
			notificationProvider: "AZURE_STORAGE_QUEUE",
			raw: map[string]interface{}{
				"name":                            "test_notification_integration",
				"comment":                         "great comment",
				"notification_provider":           "AZURE_STORAGE_QUEUE",
				"azure_storage_queue_primary_uri": "azure://great-bucket/great-path/",
				"azure_tenant_id":                 "some-guid",
			},
			expectSQL: `^CREATE NOTIFICATION INTEGRATION "test_notification_integration" AZURE_STORAGE_QUEUE_PRIMARY_URI='azure://great-bucket/great-path/' AZURE_TENANT_ID='some-guid' COMMENT='great comment' NOTIFICATION_PROVIDER='AZURE_STORAGE_QUEUE' TYPE='QUEUE' ENABLED=true$`,
		},
		{
			notificationProvider: "AWS_SQS",
			raw: map[string]interface{}{
				"name":                  "test_notification_integration",
				"comment":               "great comment",
				"direction":             "OUTBOUND",
				"notification_provider": "AWS_SQS",
				"aws_sqs_arn":           "some-sqs-arn",
				"aws_sqs_role_arn":      "some-iam-role-arn",
			},
			expectSQL: `^CREATE NOTIFICATION INTEGRATION "test_notification_integration" AWS_SQS_ARN='some-sqs-arn' AWS_SQS_ROLE_ARN='some-iam-role-arn' COMMENT='great comment' DIRECTION='OUTBOUND' NOTIFICATION_PROVIDER='AWS_SQS' TYPE='QUEUE' ENABLED=true$`,
		},
		{
			notificationProvider: "GCP_PUBSUB",
			raw: map[string]interface{}{
				"name":                         "test_notification_integration",
				"comment":                      "great comment",
				"notification_provider":        "GCP_PUBSUB",
				"gcp_pubsub_subscription_name": "some-gcp-sub-name",
			},
			expectSQL: `^CREATE NOTIFICATION INTEGRATION "test_notification_integration" COMMENT='great comment' GCP_PUBSUB_SUBSCRIPTION_NAME='some-gcp-sub-name' NOTIFICATION_PROVIDER='GCP_PUBSUB' TYPE='QUEUE' ENABLED=true$`,
		},
	}
	for _, testCase := range testCases {
		r := require.New(t)
		d := schema.TestResourceDataRaw(t, resources.NotificationIntegration().Schema, testCase.raw)
		r.NotNil(d)

		WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
			mock.ExpectExec(testCase.expectSQL).WillReturnResult(sqlmock.NewResult(1, 1))
			expectReadNotificationIntegration(mock, testCase.notificationProvider)

			err := resources.CreateNotificationIntegration(d, db)
			r.NoError(err)
		})
	}
}

func TestNotificationIntegrationRead(t *testing.T) {
	testCases := []struct {
		notificationProvider string
	}{
		{
			notificationProvider: "AZURE_STORAGE_QUEUE",
		},
		{
			notificationProvider: "AWS_SQS",
		},
		{
			notificationProvider: "GCP_PUBSUB",
		},
	}
	for _, testCase := range testCases {
		r := require.New(t)

		d := notificationIntegration(t, "test_notification_integration", map[string]interface{}{"name": "test_notification_integration"})

		WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
			expectReadNotificationIntegration(mock, testCase.notificationProvider)

			err := resources.ReadNotificationIntegration(d, db)
			r.NoError(err)
		})
	}
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

func expectReadNotificationIntegration(mock sqlmock.Sqlmock, notificationProvider string) {
	showRows := sqlmock.NewRows([]string{
		"name", "type", "category", "enabled", "created_on"},
	).AddRow("test_notification_integration", "QUEUE", "NOTIFICATION", true, "now")
	mock.ExpectQuery(`^SHOW NOTIFICATION INTEGRATIONS LIKE 'test_notification_integration'$`).WillReturnRows(showRows)

	descRows := sqlmock.NewRows([]string{
		"property", "property_type", "property_value", "property_default",
	}).AddRow("ENABLED", "Boolean", true, false)

	switch notificationProvider {
	case "AZURE_STORAGE_QUEUE":
		descRows = descRows.
			AddRow("NOTIFICATION_PROVIDER", "String", notificationProvider, nil).
			AddRow("AZURE_STORAGE_QUEUE_PRIMARY_URI", "String", "azure://great-bucket/great-path/", nil).
			AddRow("AZURE_TENANT_ID", "String", "some-guid", nil)
	case "AWS_SQS":
		descRows = descRows.
			AddRow("NOTIFICATION_PROVIDER", "String", notificationProvider, nil).
			AddRow("DIRECTION", "String", "OUTBOUND", nil).
			AddRow("AWS_SQS_ARN", "String", "some-sqs-arn", nil).
			AddRow("AWS_SQS_ROLE_ARN", "String", "some-iam-role-arn", nil).
			AddRow("AWS_SQS_EXTERNAL_ID", "String", "AGreatExternalID", nil)
	case "GCP_PUBSUB":
		descRows = descRows.
			AddRow("NOTIFICATION_PROVIDER", "String", notificationProvider, nil).
			AddRow("GCP_PUBSUB_SUBSCRIPTION_NAME", "String", "some-gcp-sub-name", nil)
	}
	mock.ExpectQuery(`DESCRIBE NOTIFICATION INTEGRATION "test_notification_integration"$`).WillReturnRows(descRows)
}
