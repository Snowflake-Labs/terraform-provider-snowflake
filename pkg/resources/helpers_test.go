package resources_test

import (
	"database/sql"
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func withMockDb(t *testing.T, f func(*sql.DB, sqlmock.Sqlmock)) {
	a := assert.New(t)
	db, mock, err := sqlmock.New()
	defer db.Close()
	a.NoError(err)

	f(db, mock)
}

func warehouse(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	a := assert.New(t)
	d := schema.TestResourceDataRaw(t, resources.Warehouse().Schema, params)
	a.NotNil(d)
	d.SetId(id)
	return d
}

func database(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	a := assert.New(t)
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, params)
	a.NotNil(d)
	d.SetId(id)
	return d
}
