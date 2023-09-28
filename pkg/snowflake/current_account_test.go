package snowflake_test

import (
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestCurrentAccountSelect(t *testing.T) {
	r := require.New(t)
	r.Equal(`SELECT CURRENT_ACCOUNT() as "account",  CASE WHEN CONTAINS(CURRENT_REGION(), '.') THEN LEFT(CURRENT_REGION(), POSITION('.' IN CURRENT_REGION()) - 1) ELSE 'PUBLIC' END AS "region_group", CASE WHEN CONTAINS(CURRENT_REGION(), '.') THEN RIGHT(CURRENT_REGION(), LENGTH(CURRENT_REGION()) - POSITION('.' IN CURRENT_REGION())) ELSE CURRENT_REGION() END AS "region";`, snowflake.SelectCurrentAccount())
}

func TestCurrentAccountRead(t *testing.T) {
	type testCaseEntry struct {
		account string
		region_group string
		region  string
		url     string
	}

	testCases := map[string]testCaseEntry{
		"aws oregon": {
			"ab1234",
			"PUBLIC",
			"AWS_US_WEST_2",
			"https://ab1234.snowflakecomputing.com",
		},
		"aws n virginia": {
			"cd5678",
			"PUBLIC",
			"AWS_US_EAST_1",
			"https://cd5678.us-east-1.snowflakecomputing.com",
		},
		"aws canada central": {
			"ef9012",
			"PUBLIC",
			"AWS_CA_CENTRAL_1",
			"https://ef9012.ca-central-1.aws.snowflakecomputing.com",
		},
		"gcp canada central": {
			"gh3456",
			"PUBLIC",
			"gcp_us_central1",
			"https://gh3456.us-central1.gcp.snowflakecomputing.com",
		},
		"azure washington": {
			"ij7890",
			"PUBLIC",
			"azure_westus2",
			"https://ij7890.west-us-2.azure.snowflakecomputing.com",
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			r := require.New(t)
			mockDB, mock, err := sqlmock.New()
			r.NoError(err)
			defer mockDB.Close()
			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

			rows := sqlmock.NewRows([]string{"account", "region_group", "region"}).AddRow(tc.account, tc.region_group, tc.region)
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT CURRENT_ACCOUNT() as "account",  CASE WHEN CONTAINS(CURRENT_REGION(), '.') THEN LEFT(CURRENT_REGION(), POSITION('.' IN CURRENT_REGION()) - 1) ELSE 'PUBLIC' END AS "region_group", CASE WHEN CONTAINS(CURRENT_REGION(), '.') THEN RIGHT(CURRENT_REGION(), LENGTH(CURRENT_REGION()) - POSITION('.' IN CURRENT_REGION())) ELSE CURRENT_REGION() END AS "region";`)).WillReturnRows(rows)

			acc, err := snowflake.ReadCurrentAccount(sqlxDB.DB)
			r.NoError(err)
			r.Equal(tc.account, acc.Account)
			r.Equal(tc.region_group, acc.RegionGroup)
			r.Equal(tc.region, acc.Region)
			url, err := acc.AccountURL()
			r.NoError(err)
			r.Equal(tc.url, url)
		})
	}
}
