package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
)

func TestDatabaseGrant(t *testing.T) {
	r := require.New(t)
	err := resources.DatabaseGrant().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestDatabaseGrantCreate(t *testing.T) {
	a := assert.New(t)

	in := map[string]interface{}{
		"database_name": "test-database",
		"privilege":     "USAGE",
		"roles":         []string{"test-role-1", "test-role-2"},
		"shares":        []string{"test-share-1", "test-share-2"},
	}
	d := schema.TestResourceDataRaw(t, resources.DatabaseGrant().Schema, in)
	a.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT USAGE ON DATABASE "test-database" TO ROLE "test-role-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON DATABASE "test-database" TO ROLE "test-role-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON DATABASE "test-database" TO SHARE "test-share-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON DATABASE "test-database" TO SHARE "test-share-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		
		err := resources.CreateDatabaseGrant(d, db)
		a.NoError(err)
	})
}
