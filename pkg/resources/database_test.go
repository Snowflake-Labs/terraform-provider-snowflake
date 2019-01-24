package resources_test

import (
	"database/sql"
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestDatabase(t *testing.T) {
	resources.Database().InternalValidate(provider.Provider().Schema, false)
}

func TestDatabaseCreate(t *testing.T) {
	a := assert.New(t)

	in := map[string]interface{}{
		"name":    "good_name",
		"comment": "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, in)
	a.NotNil(d)

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec("CREATE DATABASE good_name COMMENT='great comment").WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.CreateDatabase(d, db)
		a.NoError(err)
	})
}

func TestDatabaseRead(t *testing.T) {
	a := assert.New(t)

	d := database(t, "good_name", map[string]interface{}{"name": "good_name"})

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"created_on", "name", "is_default", "is_current", "origin", "owner", "comment", "options", "retentionTime"}).AddRow("created_on", "good_name", "is_default", "is_current", "origin", "owner", "mock comment", "options", "1")
		mock.ExpectQuery("SHOW DATABASES LIKE 'good_name'").WillReturnRows(rows)
		err := resources.ReadDatabase(d, db)
		a.NoError(err)
		a.Equal("good_name", d.Get("name").(string))
		a.Equal("mock comment", d.Get("comment").(string))
		a.Equal(1, d.Get("data_retention_time_in_days").(int))
	})
}

func TestDatabaseDelete(t *testing.T) {
	a := assert.New(t)

	d := database(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec("DROP DATABASE drop_it").WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteDatabase(d, db)
		a.NoError(err)
	})
}
