package resources_test

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUserPublicKeys(t *testing.T) {
	r := require.New(t)
	err := resources.UserPublicKeys().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func rowsFromMap(in map[string]string) *sqlmock.Rows {
	cols := []string{}
	vals := []driver.Value{}
	for col, val := range in {
		cols = append(cols, col)
		vals = append(vals, val)
	}
	rows := sqlmock.NewRows(cols)
	rows.AddRow(vals...)
	return rows
}

func TestUserPublicKeysCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":             "good_name",
		"rsa_public_key":   "asdf",
		"rsa_public_key_2": "asdf2",
	}
	d := schema.TestResourceDataRaw(t, resources.UserPublicKeys().Schema, in)
	r.NotNil(d)

	rows := map[string]string{
		"name": "good_name",
	}

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`ALTER USER "good_name" SET rsa_public_key = 'asdf'`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`ALTER USER "good_name" SET rsa_public_key_2 = 'asdf2'`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectQuery(`SHOW USERS LIKE 'good_name'`).WillReturnRows(rowsFromMap(rows))
		err := resources.CreateUserPublicKeys(d, db)
		r.NoError(err)

		r.Equal(in["name"], d.Id())
	})
}

func TestUsePublicKeysrDelete(t *testing.T) {
	r := require.New(t)
	in := map[string]interface{}{
		"name":             "good_name",
		"rsa_public_key":   "asdf",
		"rsa_public_key_2": "asdf2",
	}
	d := schema.TestResourceDataRaw(t, resources.UserPublicKeys().Schema, in)
	d.SetId(in["name"].(string))

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`ALTER USER "good_name" UNSET rsa_public_key`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`ALTER USER "good_name" UNSET rsa_public_key_2`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteUserPublicKeys(d, db)
		r.NoError(err)
	})
}
