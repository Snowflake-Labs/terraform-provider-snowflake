package resources_test

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
)

// todo: remove the rest of these which are not used. also this file should be renamed for clarity to make it clear it is for testing only
// https://snowflakecomputing.atlassian.net/browse/SNOW-936093
type grantType int

const (
	normal grantType = iota
	onFuture
	onAll
)

func TestGetPropertyAsPointer(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"integer": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"string": {
			Type:     schema.TypeString,
			Required: true,
		},
		"boolean": {
			Type:     schema.TypeBool,
			Required: true,
		},
	}, map[string]interface{}{
		"integer": 123,
		"string":  "some string",
		"boolean": true,
		"invalid": true,
	})

	assert.Equal(t, 123, *resources.GetPropertyAsPointer[int](d, "integer"))
	assert.Equal(t, "some string", *resources.GetPropertyAsPointer[string](d, "string"))
	assert.Equal(t, true, *resources.GetPropertyAsPointer[bool](d, "boolean"))
	assert.Nil(t, resources.GetPropertyAsPointer[bool](d, "invalid"))
}

func database(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func databaseGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.DatabaseGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func schemaGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.SchemaGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func stageGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.StageGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func tableGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.TableGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func viewGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ViewGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func materializedViewGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.MaterializedViewGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func resourceMonitorGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ResourceMonitorGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func integrationGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.IntegrationGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func accountGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.AccountGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func managedAccount(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ManagedAccount().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func networkPolicy(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.NetworkPolicy().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func pipe(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Pipe().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func resourceMonitor(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ResourceMonitor().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func sequence(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Sequence().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func share(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Share().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func stage(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Stage().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func stream(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Stream().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func tag(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Tag().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func roleGrants(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.RoleGrants().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func userOwnershipGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.UserOwnershipGrant().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func roleOwnershipGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.RoleOwnershipGrant().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func apiIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.APIIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func samlIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.SAMLIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func scimIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.SCIMIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func oauthIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.OAuthIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func externalFunction(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ExternalFunction().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func procedure(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Procedure().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func storageIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.StorageIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func notificationIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.NotificationIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func table(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Table().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func externalTable(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ExternalTable().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func user(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.User().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func view(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.View().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func materializedView(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.MaterializedView().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func warehouse(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Warehouse().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func externalTableGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ExternalTableGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func fileFormatGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.FileFormatGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func sequenceGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.SequenceGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func streamGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.StreamGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func maskingPolicyGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.MaskingPolicyGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func pipeGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.PipeGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func taskGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.TaskGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func rowAccessPolicy(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.RowAccessPolicy().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func rowAccessPolicyGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.RowAccessPolicyGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func function(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Function().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func tagGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.TagGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func TestIsDataType(t *testing.T) {
	isDataType := resources.IsDataType()
	key := "tag"

	testCases := []struct {
		Name  string
		Value any
		Error string
	}{
		{
			Name:  "validation: correct DataType value",
			Value: "NUMBER",
		},
		{
			Name:  "validation: correct DataType value in lowercase",
			Value: "number",
		},
		{
			Name:  "validation: incorrect DataType value",
			Value: "invalid data type",
			Error: "expected tag to be one of",
		},
		{
			Name:  "validation: incorrect value type",
			Value: 123,
			Error: "expected type of tag to be string",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			_, errors := isDataType(tt.Value, key)
			if tt.Error != "" {
				assert.Len(t, errors, 1)
				assert.ErrorContains(t, errors[0], tt.Error)
			} else {
				assert.Len(t, errors, 0)
			}
		})
	}
}

func TestIsValidIdentifier(t *testing.T) {
	accountObjectIdentifierCheck := resources.IsValidIdentifier[sdk.AccountObjectIdentifier]()
	databaseObjectIdentifierCheck := resources.IsValidIdentifier[sdk.DatabaseObjectIdentifier]()
	schemaObjectIdentifierCheck := resources.IsValidIdentifier[sdk.SchemaObjectIdentifier]()
	externalObjectIdentifierCheck := resources.IsValidIdentifier[sdk.ExternalObjectIdentifier]()
	tableColumnIdentifierCheck := resources.IsValidIdentifier[sdk.TableColumnIdentifier]()

	testCases := []struct {
		Name       string
		Value      any
		Error      string
		CheckingFn schema.SchemaValidateDiagFunc
	}{
		{
			Name:       "validation: invalid value type",
			Value:      123,
			Error:      "Expected schema string type, but got: int",
			CheckingFn: accountObjectIdentifierCheck,
		},
		{
			Name:       "validation: invalid identifier representation",
			Value:      "",
			Error:      "Unable to parse the identifier: ",
			CheckingFn: accountObjectIdentifierCheck,
		},
		// Invalid form for different checkers (tests getExpectedIdentifierForm function)
		{
			Name:       "validation: incorrect form for account object identifier",
			Value:      "a.b",
			Error:      "<name>, but was <database_name>.<name>",
			CheckingFn: accountObjectIdentifierCheck,
		},
		{
			Name:       "validation: incorrect form for database object identifier",
			Value:      "a.b.c",
			Error:      "<database_name>.<name>, but was <database_name>.<schema_name>.<name>",
			CheckingFn: databaseObjectIdentifierCheck,
		},
		{
			Name:       "validation: incorrect form for schema object identifier",
			Value:      "a.b.c.d",
			Error:      "<database_name>.<schema_name>.<name>, but was <database_name>.<schema_name>.<table_name>.<column_name>",
			CheckingFn: schemaObjectIdentifierCheck,
		},
		{
			Name:       "validation: incorrect form for table column identifier",
			Value:      "a",
			Error:      "<database_name>.<schema_name>.<table_name>.<column_name>, but was <name>",
			CheckingFn: tableColumnIdentifierCheck,
		},
		{
			Name:       "validation: external object identifier is validated",
			Value:      "a",
			Error:      "Identifier validation is not available for sdk.ExternalObjectIdentifier type.",
			CheckingFn: externalObjectIdentifierCheck,
		},
		// Valid form for different checkers (tests getExpectedIdentifierForm function)
		{
			Name:       "correct form for account object identifier",
			Value:      "a",
			CheckingFn: accountObjectIdentifierCheck,
		},
		{
			Name:       "correct form for database object identifier",
			Value:      "a.b",
			CheckingFn: databaseObjectIdentifierCheck,
		},
		{
			Name:       "correct form for schema object identifier",
			Value:      "a.b.c",
			CheckingFn: schemaObjectIdentifierCheck,
		},
		{
			Name:       "correct form for table column identifier",
			Value:      "a.b.c.d",
			CheckingFn: tableColumnIdentifierCheck,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			diag := tt.CheckingFn(tt.Value, cty.IndexStringPath("path"))
			if tt.Error != "" {
				assert.Len(t, diag, 1)
				assert.Contains(t, diag[0].Detail, tt.Error)
			} else {
				assert.Len(t, diag, 0)
			}
		})
	}
}
