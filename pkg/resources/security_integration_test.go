package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestSecurityIntegration(t *testing.T) {
	r := require.New(t)
	err := resources.SecurityIntegration().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestSecuritytIntegrationCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name": "test-security-integration",
	}

	d := schema.TestResourceDataRaw(t, resources.SecurityIntegration().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^CREATE SECURITY INTEGRATION "test-security-integration"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		// expectReadSecurityIntegration(mock)
	})
}

func expectReadSecurityIntegration(mock sqlmock.Sqlmock) {
	// rows := sqlmock.NewRows([]string{})
}
