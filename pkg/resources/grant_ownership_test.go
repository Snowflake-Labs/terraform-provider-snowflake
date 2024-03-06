package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetOnObjectIdentifier(t *testing.T) {
	testCases := []struct {
		Name       string
		ObjectType sdk.ObjectType
		ObjectName string
		Expected   sdk.ObjectIdentifier
		Error      string
	}{
		{
			Name:       "database - account object identifier",
			ObjectType: sdk.ObjectTypeDatabase,
			ObjectName: "test_database",
			Expected:   sdk.NewAccountObjectIdentifier("test_database"),
		},
		{
			Name:       "database - account object identifier - quoted",
			ObjectType: sdk.ObjectTypeDatabase,
			ObjectName: "\"test_database\"",
			Expected:   sdk.NewAccountObjectIdentifier("test_database"),
		},
		{
			Name:       "schema - database object identifier",
			ObjectType: sdk.ObjectTypeSchema,
			ObjectName: "test_database.test_schema",
			Expected:   sdk.NewDatabaseObjectIdentifier("test_database", "test_schema"),
		},
		{
			Name:       "schema - database object identifier - quoted",
			ObjectType: sdk.ObjectTypeSchema,
			ObjectName: "\"test_database\".\"test_schema\"",
			Expected:   sdk.NewDatabaseObjectIdentifier("test_database", "test_schema"),
		},
		{
			Name:       "table - schema object identifier",
			ObjectType: sdk.ObjectTypeTable,
			ObjectName: "test_database.test_schema.test_table",
			Expected:   sdk.NewSchemaObjectIdentifier("test_database", "test_schema", "test_table"),
		},
		{
			Name:       "table - schema object identifier - quoted",
			ObjectType: sdk.ObjectTypeTable,
			ObjectName: "\"test_database\".\"test_schema\".\"test_table\"",
			Expected:   sdk.NewSchemaObjectIdentifier("test_database", "test_schema", "test_table"),
		},
		{
			Name:       "validation - valid identifier",
			ObjectType: sdk.ObjectTypeDatabase,
			ObjectName: "to.many.parts.in.this.identifier",
			Error:      "unable to classify identifier",
		},
		{
			Name:       "validation - unsupported type",
			ObjectType: sdk.ObjectTypeShare,
			ObjectName: "some_share",
			Error:      "object_type SHARE is not supported",
		},
		{
			Name:       "validation - invalid database account object identifier",
			ObjectType: sdk.ObjectTypeDatabase,
			ObjectName: "test_database.test_schema",
			Error:      "invalid object_name test_database.test_schema expected account object identifier",
		},
		{
			Name:       "validation - invalid database account object identifier",
			ObjectType: sdk.ObjectTypeSchema,
			ObjectName: "test_database.test_schema.test_table",
		},
		{
			Name:       "table - schema object identifier",
			ObjectType: sdk.ObjectTypeTable,
			ObjectName: "test_database.test_schema.test_table",
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			id, err := getOnObjectIdentifier(tt.ObjectType, tt.ObjectName)
			if tt.Error == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.Expected, id)
			} else {
				assert.ErrorContains(t, err, tt.Error)
			}
		})
	}
}
