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

func TestMaskingPolicyGrant(t *testing.T) {
	r := require.New(t)
	err := resources.MaskingPolicyGrant().Resource.InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestMaskingPolicyGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"masking_policy_name": "test-masking-policy",
		"schema_name":         "PUBLIC",
		"database_name":       "test-db",
		"privilege":           "APPLY",
		"roles":               []interface{}{"test-role-1", "test-role-2"},
	}
	d := schema.TestResourceDataRaw(t, resources.MaskingPolicyGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT APPLY ON MASKING POLICY "test-db"."PUBLIC"."test-masking-policy" TO ROLE "test-role-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT APPLY ON MASKING POLICY "test-db"."PUBLIC"."test-masking-policy" TO ROLE "test-role-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadMaskingPolicyGrant(mock)
		err := resources.CreateMaskingPolicyGrant(d, db)
		r.NoError(err)
	})
}

func TestMaskingPolicyGrantRead(t *testing.T) {
	r := require.New(t)

	d := maskingPolicyGrant(t, "test-db|PUBLIC|test-masking-policy|APPLY||false", map[string]interface{}{
		"masking_policy_name": "test-masking-policy",
		"schema_name":         "PUBLIC",
		"database_name":       "test-db",
		"privilege":           "APPLY",
		"roles":               []interface{}{},
		"with_grant_option":   false,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadMaskingPolicyGrant(mock)
		err := resources.ReadMaskingPolicyGrant(d, db)
		r.NoError(err)
	})

	roles := d.Get("roles").(*schema.Set)
	r.True(roles.Contains("test-role-1"))
	r.True(roles.Contains("test-role-2"))
	r.Equal(2, roles.Len())
}

func expectReadMaskingPolicyGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "APPLY", "MASKING_POLICY", "test-masking-policy", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "APPLY", "MASKING_POLICY", "test-masking-policy", "ROLE", "test-role-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON MASKING POLICY "test-db"."PUBLIC"."test-masking-policy"$`).WillReturnRows(rows)
}

func TestParseMaskingPolicyGrantID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseMaskingPolicyGrantID("test-db|PUBLIC|test-masking-policy|APPLY|true|role1,role2")
	r.NoError(err)
	r.Equal("test-db", grantID.DatabaseName)
	r.Equal("PUBLIC", grantID.SchemaName)
	r.Equal("test-masking-policy", grantID.ObjectName)
	r.Equal("APPLY", grantID.Privilege)
	r.Equal(true, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
}

func TestParseMaskingPolicyGrantEmojiID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseMaskingPolicyGrantID("test-db❄️PUBLIC❄️test-masking-policy❄️APPLY❄️true❄️role1,role2")
	r.NoError(err)
	r.Equal("test-db", grantID.DatabaseName)
	r.Equal("PUBLIC", grantID.SchemaName)
	r.Equal("test-masking-policy", grantID.ObjectName)
	r.Equal("APPLY", grantID.Privilege)
	r.Equal(true, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
}

func TestParseMaskingPolicyGrantOldID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseMaskingPolicyGrantID("test-db|PUBLIC|test-materialized-view|SELECT|role1,role2|true")
	r.NoError(err)
	r.Equal("test-db", grantID.DatabaseName)
	r.Equal("PUBLIC", grantID.SchemaName)
	r.Equal("test-materialized-view", grantID.ObjectName)
	r.Equal("SELECT", grantID.Privilege)
	r.Equal(true, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
}
