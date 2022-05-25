package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestNetworkPolicyAttachment(t *testing.T) {
	r := require.New(t)
	err := resources.NetworkPolicyAttachment().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestNetworkPolicyAttachmentCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"network_policy_name": "test-network-policy",
		"set_for_account":     true,
		"users":               []interface{}{"test-user"},
	}
	d := schema.TestResourceDataRaw(t, resources.NetworkPolicyAttachment().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^ALTER ACCOUNT SET NETWORK_POLICY = "test-network-policy"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^DESCRIBE USER "test-user"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^ALTER USER "test-user" SET NETWORK_POLICY = "test-network-policy"$`).WillReturnResult(sqlmock.NewResult(1, 1))

		err := resources.CreateNetworkPolicyAttachment(d, db)
		r.NoError(err)
	})
}

func TestNetworkPolicyAttachmentSetOnAccountDelete(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"network_policy_name": "test-network-policy",
		"set_for_account":     true,
		"users":               []interface{}{"test-user"},
	}
	d := schema.TestResourceDataRaw(t, resources.NetworkPolicyAttachment().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^ALTER ACCOUNT UNSET NETWORK_POLICY$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^DESCRIBE USER "test-user"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^ALTER USER "test-user" UNSET NETWORK_POLICY$`).WillReturnResult(sqlmock.NewResult(1, 1))

		err := resources.DeleteNetworkPolicyAttachment(d, db)
		r.NoError(err)
	})
}

func TestNetworkPolicyAttachmentDelete(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"network_policy_name": "test-network-policy",
		"set_for_account":     false,
		"users":               []interface{}{"test-user"},
	}
	d := schema.TestResourceDataRaw(t, resources.NetworkPolicyAttachment().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^DESCRIBE USER "test-user"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^ALTER USER "test-user" UNSET NETWORK_POLICY$`).WillReturnResult(sqlmock.NewResult(1, 1))

		err := resources.DeleteNetworkPolicyAttachment(d, db)
		r.NoError(err)
	})
}
