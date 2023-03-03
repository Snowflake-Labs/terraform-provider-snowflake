package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
)

func TestResourceMonitorGrant(t *testing.T) {
	r := require.New(t)
	err := resources.ResourceMonitorGrant().Resource.InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestResourceMonitorGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"monitor_name":      "test-monitor",
		"privilege":         "MONITOR",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": true,
	}
	d := schema.TestResourceDataRaw(t, resources.ResourceMonitorGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT MONITOR ON RESOURCE MONITOR "test-monitor" TO ROLE "test-role-1" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT MONITOR ON RESOURCE MONITOR "test-monitor" TO ROLE "test-role-2" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadResourceMonitorGrant(mock)
		err := resources.CreateResourceMonitorGrant(d, db)
		r.NoError(err)
	})
}

func TestResourceMonitorGrantRead(t *testing.T) {
	r := require.New(t)

	d := resourceMonitorGrant(t, "test-monitor❄️MONITOR❄️❄️false", map[string]interface{}{
		"monitor_name":      "test-monitor",
		"privilege":         "MONITOR",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": false,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadResourceMonitorGrant(mock)
		err := resources.ReadResourceMonitorGrant(d, db)
		r.NoError(err)
	})
}

func expectReadResourceMonitorGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "MONITOR", "RESOURCE MONITOR", "test-monitor", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "MONITOR", "RESOURCE MONITOR", "test-monitor", "ROLE", "test-role-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON RESOURCE MONITOR "test-monitor"$`).WillReturnRows(rows)
}

func TestParseResourceMonitorGrantID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseResourceMonitorGrantID("test-rm|MONITOR|true|role1,role2")
	r.NoError(err)
	r.Equal("test-rm", grantID.ObjectName)
	r.Equal("MONITOR", grantID.Privilege)
	r.Equal(true, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
}

func TestParseResourceMonitorGrantEmojiID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseResourceMonitorGrantID("test-rm❄️MONITOR❄️true❄️role1,role2")
	r.NoError(err)
	r.Equal("test-rm", grantID.ObjectName)
	r.Equal("MONITOR", grantID.Privilege)
	r.Equal(true, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
}

func TestParseResourceMonitorGrantOldID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseResourceMonitorGrantID("test-rm|||MONITOR|role1,role2|true")
	r.NoError(err)
	r.Equal("test-rm", grantID.ObjectName)
	r.Equal("MONITOR", grantID.Privilege)
	r.Equal(true, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
}
