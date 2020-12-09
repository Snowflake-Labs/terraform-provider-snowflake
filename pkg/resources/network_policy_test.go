package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestNetworkPolicy(t *testing.T) {
	r := require.New(t)
	err := resources.NetworkPolicy().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestNetworkPolicyCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":            "test-network-policy",
		"allowed_ip_list": []interface{}{"192.168.1.0/24"},
		"blocked_ip_list": []interface{}{"155.548.2.98"},
		"comment":         "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.NetworkPolicy().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^CREATE NETWORK POLICY "test-network-policy" ALLOWED_IP_LIST=\('192\.168\.1\.0/24'\) BLOCKED_IP_LIST=\('155\.548\.2\.98'\) COMMENT="great comment"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadNetworkPolicy(mock)
		err := resources.CreateNetworkPolicy(d, db)
		r.NoError(err)
	})
}

func expectReadNetworkPolicy(mock sqlmock.Sqlmock) {
	showRows := sqlmock.NewRows([]string{
		"created_on", "name", "comment", "entries_in_allowed_ip_list", "entries_in_blocked_ip_list",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "test-network-policy", "this is a comment", 2, 1,
	)
	mock.ExpectQuery(`^SHOW NETWORK POLICIES$`).WillReturnRows(showRows)

	descRows := sqlmock.NewRows([]string{
		"name", "value",
	}).AddRow(
		"ALLOWED_IP_LIST", "192.168.0.100,192.168.0.200/18",
	).AddRow(
		"BLOCKED_IP_LIST", "192.168.0.101",
	)
	mock.ExpectQuery(`^DESC NETWORK POLICY "test-network-policy"$`).WillReturnRows(descRows)
}

func TestNetworkPolicyDelete(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":            "test-network-policy",
		"allowed_ip_list": []interface{}{"192.168.1.0/24"},
		"blocked_ip_list": []interface{}{"155.548.2.98"},
		"comment":         "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.NetworkPolicy().Schema, in)
	d.SetId("test-network-policy")
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^DROP NETWORK POLICY "test-network-policy"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteNetworkPolicy(d, db)
		r.NoError(err)
	})
}

func TestIpListToString(t *testing.T) {
	r := require.New(t)

	in := []string{"192.168.0.100/24", "29.254.123.20"}
	out := snowflake.IpListToString(in)

	r.Equal("('192.168.0.100/24', '29.254.123.20')", out)
}

func TestNetworkPolicyRead(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":            "good-network-policy",
		"allowed_ip_list": []interface{}{"192.168.1.0/24"},
		"blocked_ip_list": []interface{}{"155.548.2.98"},
	}

	d := networkPolicy(t, "good-network-policy", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		// Test when resource is not found, checking if state will be empty
		r.NotEmpty(d.State())
		q := snowflake.NetworkPolicy(d.Id()).ShowAllNetworkPolicies()
		mock.ExpectQuery(q).WillReturnError(sql.ErrNoRows)
		err1 := resources.ReadNetworkPolicy(d, db)
		r.Empty(d.State())

		rows := sqlmock.NewRows([]string{
			"created_on", "name", "comment", "entries_in_allowed_ip_list", "entries_in_blocked_ip_list",
		}).AddRow(
			time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "bad-network-policy", "this is a comment", 2, 1,
		)
		mock.ExpectQuery(q).WillReturnRows(rows)
		err2 := resources.ReadNetworkPolicy(d, db)

		r.Nil(err1)
		r.Nil(err2)
	})
}
