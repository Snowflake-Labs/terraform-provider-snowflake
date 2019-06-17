package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
)

func TestShare(t *testing.T) {
	r := require.New(t)
	err := resources.Share().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestShareCreate(t *testing.T) {
	a := assert.New(t)

	in := map[string]interface{}{
		"name":    "test-share",
		"comment": "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.Share().Schema, in)
	a.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^CREATE SHARE "test-share" COMMENT='great comment'$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadShare(mock)
		err := resources.CreateShare(d, db)
		a.NoError(err)
	})
}

func expectReadShare(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"name", "comment",
	}).AddRow("test-share", "great comment")
	mock.ExpectQuery(`^SHOW SHARES LIKE 'test-share'$`).WillReturnRows(rows)
}
