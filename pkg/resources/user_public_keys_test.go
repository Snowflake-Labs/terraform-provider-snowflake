package resources_test

import (
	"database/sql"
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

func TestUserPublicKeysCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":             "good_name",
		"rsa_public_key":   "asdf",
		"rsa_public_key_2": "asdf2",
	}
	d := schema.TestResourceDataRaw(t, resources.UserPublicKeys().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`ALTER USER "good_name" SET rsa_public_key = 'asdf'`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`ALTER USER "good_name" SET rsa_public_key_2 = 'asdf2'`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadUser(mock, "good_name")
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
