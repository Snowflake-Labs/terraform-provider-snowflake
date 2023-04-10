package snowflake_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	r := require.New(t)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{
			"organization_name", "account_name", "snowflake_region", "edition", "account_url", "created_on", "comment", "account_locator", "account_locator_url", "managed_accounts", "consumption_billing_entity_name", "marketplace_consumer_billing_entity_name", "marketplace_provider_billing_entity_name", "old_account_url", "is_org_admin",
		},
		).AddRow(
			"TEST_ORG", "TEST_ACC", "AWS_US_WEST_2", "STANDARD", "https://acc.snowflakecomputing.com", "2020-01-01 12:00:00.123 -0800", "SNOWFLAKE", "LOCATOR", "https://loc.snowflakecomputing.com", 0, "BillingEntity", "", "", "", true,
		)
		mock.ExpectQuery(`SHOW ORGANIZATION ACCOUNTS LIKE 'test_account'`).WillReturnRows(rows)
		a, err := snowflake.ShowAccount(db, "test_account")
		r.NoError(err)
		r.Equal("TEST_ORG", a.OrganizationName.String)
		r.Equal("TEST_ACC", a.AccountName.String)
		r.Equal("AWS_US_WEST_2", a.SnowflakeRegion.String)
		r.Equal(true, a.IsOrgAdmin.Bool)
	})
}

