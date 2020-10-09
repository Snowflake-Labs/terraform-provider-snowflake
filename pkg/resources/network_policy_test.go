package resources_test

import (
	"database/sql"
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
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

	// FIXME: this doesn't play nicely with the big post-processing SHOW/DESC sql
	//WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
	//	mock.ExpectExec(`^CREATE OR REPLACE NETWORK POLICY "test-network-policy" ALLOWED_IP_LIST=\('192\.168\.1\.0/24'\) BLOCKED_IP_LIST=\('155\.548\.2\.98'\) COMMENT="great comment"$`).WillReturnResult(sqlmock.NewResult(1, 1))
	//	mock.ExpectExec(`DESC NETWORK POLICY "test-network-policy"`).WillReturnResult(sqlmock.NewResult(1, 1))
	//	mock.ExpectExec(`SHOW NETWORK POLICIES`).WillReturnResult(sqlmock.NewResult(1, 1))
	//	err := resources.CreateNetworkPolicy(d, db)
	//	r.NoError(err)
	//})
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
