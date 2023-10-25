// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/internal/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestAPIIntegration(t *testing.T) {
	r := require.New(t)
	err := resources.APIIntegration().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestAPIIntegrationCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":                 "test_api_integration",
		"api_allowed_prefixes": []interface{}{"https://123456.execute-api.us-west-2.amazonaws.com/prod/"},
		"api_provider":         "aws_api_gateway",
		"api_aws_role_arn":     "arn:aws:iam::000000000001:/role/test",
		"api_key":              "12345",
	}

	in2 := map[string]interface{}{
		"name":                 "test_gov_api_integration",
		"api_allowed_prefixes": []interface{}{"https://123456.execute-api.us-gov-west-1.amazonaws.com/prod/"},
		"api_provider":         "aws_gov_api_gateway",
		"api_aws_role_arn":     "arn:aws:iam::000000000001:/role/test",
		"api_key":              "12345",
	}

	d := schema.TestResourceDataRaw(t, resources.APIIntegration().Schema, in)
	d2 := schema.TestResourceDataRaw(t, resources.APIIntegration().Schema, in2)

	r.NotNil(d)
	r.NotNil(d2)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE API INTEGRATION "test_api_integration" API_PROVIDER=aws_api_gateway API_AWS_ROLE_ARN='arn:aws:iam::000000000001:/role/test' API_KEY='12345' API_ALLOWED_PREFIXES=\('https://123456.execute-api.us-west-2.amazonaws.com/prod/'\) ENABLED=true$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadAPIIntegration(mock)

		err := resources.CreateAPIIntegration(d, db)
		r.NoError(err)
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE API INTEGRATION "test_gov_api_integration" API_PROVIDER=aws_gov_api_gateway API_AWS_ROLE_ARN='arn:aws:iam::000000000001:/role/test' API_KEY='12345' API_ALLOWED_PREFIXES=\('https://123456.execute-api.us-gov-west-1.amazonaws.com/prod/'\) ENABLED=true$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadGovAPIIntegration(mock)

		err := resources.CreateAPIIntegration(d2, db)
		r.NoError(err)
	})
}

func TestAPIIntegrationRead(t *testing.T) {
	r := require.New(t)

	d := apiIntegration(t, "test_api_integration", map[string]interface{}{"name": "test_api_integration"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadAPIIntegration(mock)

		err := resources.ReadAPIIntegration(d, db)
		r.NoError(err)
	})
}

func TestAPIIntegrationDelete(t *testing.T) {
	r := require.New(t)

	d := apiIntegration(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP API INTEGRATION "drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteAPIIntegration(d, db)
		r.NoError(err)
	})
}

func expectReadAPIIntegration(mock sqlmock.Sqlmock) {
	showRows := sqlmock.NewRows([]string{
		"name", "type", "category", "enabled", "created_on",
	},
	).AddRow("test_api_integration", "EXTERNAL_API", "API", true, "now")
	mock.ExpectQuery(`^SHOW API INTEGRATIONS LIKE 'test_api_integration'$`).WillReturnRows(showRows)

	descRows := sqlmock.NewRows([]string{
		"property", "property_type", "property_value", "property_default",
	}).AddRow("ENABLED", "Boolean", true, false).
		AddRow("API_KEY", "String", "12345", nil).
		AddRow("API_ALLOWED_PREFIXES", "List", "https://123456.execute-api.us-west-2.amazonaws.com/prod/,https://123456.execute-api.us-west-2.amazonaws.com/staging/", nil).
		AddRow("API_AWS_IAM_USER_ARN", "String", "arn:aws:iam::000000000000:/user/test", nil).
		AddRow("API_AWS_ROLE_ARN", "String", "arn:aws:iam::000000000001:/role/test", nil).
		AddRow("API_AWS_EXTERNAL_ID", "String", "AGreatExternalID", nil)

	mock.ExpectQuery(`DESCRIBE API INTEGRATION "test_api_integration"$`).WillReturnRows(descRows)
}

func expectReadGovAPIIntegration(mock sqlmock.Sqlmock) {
	showRows := sqlmock.NewRows([]string{
		"name", "type", "category", "enabled", "created_on",
	},
	).AddRow("test_gov_api_integration", "EXTERNAL_API", "API", true, "now")
	mock.ExpectQuery(`^SHOW API INTEGRATIONS LIKE 'test_gov_api_integration'$`).WillReturnRows(showRows)

	descRows := sqlmock.NewRows([]string{
		"property", "property_type", "property_value", "property_default",
	}).AddRow("ENABLED", "Boolean", true, false).
		AddRow("API_KEY", "String", "12345", nil).
		AddRow("API_ALLOWED_PREFIXES", "List", "https://123456.execute-api.us-gov-west-1.amazonaws.com/prod/,https://123456.execute-api.us-gov-west-1.amazonaws.com/staging/", nil).
		AddRow("API_AWS_IAM_USER_ARN", "String", "arn:aws:iam::000000000000:/user/test", nil).
		AddRow("API_AWS_ROLE_ARN", "String", "arn:aws:iam::000000000001:/role/test", nil).
		AddRow("API_AWS_EXTERNAL_ID", "String", "AGreatExternalID", nil)

	mock.ExpectQuery(`DESCRIBE API INTEGRATION "test_gov_api_integration"$`).WillReturnRows(descRows)
}
