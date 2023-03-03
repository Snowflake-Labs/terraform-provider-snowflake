package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestRowAccessPolicyGrant(t *testing.T) {
	r := require.New(t)
	err := resources.RowAccessPolicyGrant().Resource.InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestRowAccessPolicyGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"row_access_policy_name": "test-row-access-policy",
		"schema_name":            "PUBLIC",
		"database_name":          "test-db",
		"privilege":              "APPLY",
		"roles":                  []interface{}{"test-role-1", "test-role-2"},
	}
	d := schema.TestResourceDataRaw(t, resources.RowAccessPolicyGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT APPLY ON ROW ACCESS POLICY "test-db"."PUBLIC"."test-row-access-policy" TO ROLE "test-role-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT APPLY ON ROW ACCESS POLICY "test-db"."PUBLIC"."test-row-access-policy" TO ROLE "test-role-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadRowAccessPolicyGrant(mock)
		err := resources.CreateRowAccessPolicyGrant(d, db)
		r.NoError(err)
	})
}

func TestRowAccessPolicyGrantRead(t *testing.T) {
	r := require.New(t)

	d := rowAccessPolicyGrant(t, "test-db|PUBLIC|test-row-access-policy|APPLY||false", map[string]interface{}{
		"row_access_policy_name": "test-row-access-policy",
		"schema_name":            "PUBLIC",
		"database_name":          "test-db",
		"privilege":              "APPLY",
		"roles":                  []interface{}{},
		"with_grant_option":      false,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadRowAccessPolicyGrant(mock)
		err := resources.ReadRowAccessPolicyGrant(d, db)
		r.NoError(err)
	})

	roles := d.Get("roles").(*schema.Set)
	r.True(roles.Contains("test-role-1"))
	r.True(roles.Contains("test-role-2"))
	r.Equal(2, roles.Len())
}

func expectReadRowAccessPolicyGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "APPLY", "ROW_ACCESS_POLICY", "test-row-access-policy", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "APPLY", "ROW_ACCESS_POLICY", "test-row-access-policy", "ROLE", "test-role-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON ROW ACCESS POLICY "test-db"."PUBLIC"."test-row-access-policy"$`).WillReturnRows(rows)
}

func TestRowAccessPolicyGrantID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseRowAccessPolicyGrantID("test-db|test-schema|test-policy|APPLY|true|role1,role2")
	r.NoError(err)
	r.Equal("test-db", grantID.DatabaseName)
	r.Equal("test-schema", grantID.SchemaName)
	r.Equal("test-policy", grantID.ObjectName)
	r.Equal("APPLY", grantID.Privilege)
	r.Equal(true, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
}

func TestParseRowAccessPolicyGrantEmojiID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseRowAccessPolicyGrantID("test-db❄️test-schema❄️test-policy❄️APPLY❄️true❄️role1,role2")
	r.NoError(err)
	r.Equal("test-db", grantID.DatabaseName)
	r.Equal("test-schema", grantID.SchemaName)
	r.Equal("test-policy", grantID.ObjectName)
	r.Equal("APPLY", grantID.Privilege)
	r.Equal(true, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
}

func TestParseRowAccessPolicyGrantOldID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseRowAccessPolicyGrantID("test-db|test-schema|test-policy|APPLY|role1,role2|true")
	r.NoError(err)
	r.Equal("test-db", grantID.DatabaseName)
	r.Equal("test-schema", grantID.SchemaName)
	r.Equal("test-policy", grantID.ObjectName)
	r.Equal("APPLY", grantID.Privilege)
	r.Equal(true, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
}
