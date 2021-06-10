package snowflake_test

import (
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCurrentAccountSelect(t *testing.T) {
	r := require.New(t)
	r.Equal(`SELECT CURRENT_ACCOUNT() AS "account", CURRENT_REGION() AS "region";`, snowflake.SelectCurrentAccount())
}

func TestCurrentAccountRead(t *testing.T) {
	type testCaseEntry struct {
		account string
		region  string
		url     string
	}

	testCases := map[string]testCaseEntry{
		"aws oregon": {
			"ab1234",
			"AWS_US_WEST_2",
			"https://ab1234.snowflakecomputing.com",
		},
		"aws n virginia": {
			"cd5678",
			"AWS_US_EAST_1",
			"https://cd5678.us-east-1.snowflakecomputing.com",
		},
		"aws canada central": {
			"ef9012",
			"AWS_CA_CENTRAL_1",
			"https://ef9012.ca-central-1.aws.snowflakecomputing.com",
		},
		"gcp canada central": {
			"gh3456",
			"gcp_us_central1",
			"https://gh3456.us-central1.gcp.snowflakecomputing.com",
		},
		"azure washington": {
			"ij7890",
			"azure_westus2",
			"https://ij7890.west-us-2.azure.snowflakecomputing.com",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			r := require.New(t)
			mockDB, mock, err := sqlmock.New()
			r.NoError(err)
			defer mockDB.Close()
			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

			rows := sqlmock.NewRows([]string{"account", "region"}).AddRow(testCase.account, testCase.region)
			mock.ExpectQuery(`SELECT CURRENT_ACCOUNT\(\) AS "account", CURRENT_REGION\(\) AS "region";`).WillReturnRows(rows)

			acc, err := snowflake.ReadCurrentAccount(sqlxDB.DB)
			r.NoError(err)
			r.Equal(testCase.account, acc.Account)
			r.Equal(testCase.region, acc.Region)
			url, err := acc.AccountURL()
			r.NoError(err)
			r.Equal(testCase.url, url)
		})
	}
}
