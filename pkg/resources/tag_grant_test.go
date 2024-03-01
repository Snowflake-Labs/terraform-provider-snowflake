package resources_test

import (
	"database/sql"
	"testing"
	"time"

	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestTagGrant(t *testing.T) {
	r := require.New(t)
	err := resources.TagGrant().Resource.InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestTagGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"tag_name":      "test_tag",
		"schema_name":   "PUBLIC",
		"database_name": "test-db",
		"privilege":     "APPLY",
		"roles":         []interface{}{"test-role-1", "test-role-2"},
	}
	d := schema.TestResourceDataRaw(t, resources.TagGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT APPLY ON TAG "test-db"."PUBLIC"."test_tag" TO ROLE "test-role-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT APPLY ON TAG "test-db"."PUBLIC"."test_tag" TO ROLE "test-role-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadTagGrant(mock)
		err := resources.CreateTagGrant(d, &internalprovider.Context{
			Client: sdk.NewClientFromDB(db),
		})
		r.NoError(err)
	})
}

func TestTagGrantRead(t *testing.T) {
	r := require.New(t)

	d := tagGrant(t, "test-db|PUBLIC|test_tag|APPLY||false", map[string]interface{}{
		"tag_name":          "test_tag",
		"schema_name":       "PUBLIC",
		"database_name":     "test-db",
		"privilege":         "APPLY",
		"roles":             []interface{}{},
		"with_grant_option": false,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadTagGrant(mock)
		err := resources.ReadTagGrant(d, &internalprovider.Context{
			Client: sdk.NewClientFromDB(db),
		})
		r.NoError(err)
	})

	roles := d.Get("roles").(*schema.Set)
	r.True(roles.Contains("test-role-1"))
	r.True(roles.Contains("test-role-2"))
	r.Equal(2, roles.Len())
}

func expectReadTagGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "APPLY", "TAG", "test_tag", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "APPLY", "TAG", "test_tag", "ROLE", "test-role-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON TAG "test-db"."PUBLIC"."test_tag"$`).WillReturnRows(rows)
}
