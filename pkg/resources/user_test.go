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

func TestUser(t *testing.T) {
	resources.User().InternalValidate(provider.Provider().Schema, false)
}

func TestUserCreate(t *testing.T) {
	a := assert.New(t)

	in := map[string]interface{}{
		"name":    "good_name",
		"comment": "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.User().Schema, in)
	a.NotNil(d)

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec("CREATE USER good_name COMMENT='great comment").WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadUser(mock)
		err := resources.CreateUser(d, db)
		a.NoError(err)
	})
}

func expectReadUser(mock sqlmock.Sqlmock) {
	mock.ExpectExec("USE WAREHOUSE good_name").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("SHOW USERS LIKE 'good_name'").WillReturnResult(sqlmock.NewResult(1, 1))
	rows := sqlmock.NewRows([]string{"name", "comment"}).AddRow("good_name", "mock comment")
	mock.ExpectQuery(`select "name", "comment" from table\(result_scan\(last_query_id\(\)\)\)`).WillReturnRows(rows)
}

func TestUserRead(t *testing.T) {
	a := assert.New(t)

	d := warehouse(t, "good_name", map[string]interface{}{"name": "good_name"})

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadUser(mock)
		err := resources.ReadUser(d, db)
		a.NoError(err)
		a.Equal("mock comment", d.Get("comment").(string))
	})
}

func TestUserDelete(t *testing.T) {
	a := assert.New(t)

	d := warehouse(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec("DROP USER drop_it").WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteUser(d, db)
		a.NoError(err)
	})
}
