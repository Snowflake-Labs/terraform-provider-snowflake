package resources_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestTagIdentifierAndObjectIdentifier(t *testing.T) {
	t.Run("account object identifier", func(t *testing.T) {
		in := map[string]interface{}{
			"tag_id":      "\"test_db\".\"test_schema\".\"test_tag\"",
			"object_type": "DATABASE",
			"object_identifier": []interface{}{map[string]interface{}{
				"name": "test_db",
			}},
		}
		d := schema.TestResourceDataRaw(t, resources.TagAssociation().Schema, in)
		tid, identifiers, objectType := resources.TagIdentifierAndObjectIdentifier(d)
		assert.Equal(t, sdk.NewSchemaObjectIdentifier("test_db", "test_schema", "test_tag"), tid)
		assert.Len(t, identifiers, 1)
		assert.Equal(t, "\"test_db\"", identifiers[0].FullyQualifiedName())
		assert.Equal(t, sdk.ObjectTypeDatabase, objectType)
	})

	t.Run("database object identifier", func(t *testing.T) {
		in := map[string]interface{}{
			"tag_id":      "\"test_db\".\"test_schema\".\"test_tag\"",
			"object_type": "SCHEMA",
			"object_identifier": []interface{}{map[string]interface{}{
				"name":     "test_schema",
				"database": "test_db",
			}},
		}
		d := schema.TestResourceDataRaw(t, resources.TagAssociation().Schema, in)
		tid, identifiers, objectType := resources.TagIdentifierAndObjectIdentifier(d)
		assert.Equal(t, sdk.NewSchemaObjectIdentifier("test_db", "test_schema", "test_tag"), tid)
		assert.Len(t, identifiers, 1)
		assert.Equal(t, "\"test_db\".\"test_schema\"", identifiers[0].FullyQualifiedName())
		assert.Equal(t, sdk.ObjectTypeSchema, objectType)
	})

	t.Run("schema object identifier", func(t *testing.T) {
		in := map[string]interface{}{
			"tag_id":      "\"test_db\".\"test_schema\".\"test_tag\"",
			"object_type": "TABLE",
			"object_identifier": []interface{}{map[string]interface{}{
				"name":     "test_table",
				"database": "test_db",
				"schema":   "test_schema",
			}},
		}
		d := schema.TestResourceDataRaw(t, resources.TagAssociation().Schema, in)
		tid, identifiers, objectType := resources.TagIdentifierAndObjectIdentifier(d)
		assert.Equal(t, sdk.NewSchemaObjectIdentifier("test_db", "test_schema", "test_tag"), tid)
		assert.Len(t, identifiers, 1)
		assert.Equal(t, "\"test_db\".\"test_schema\".\"test_table\"", identifiers[0].FullyQualifiedName())
		assert.Equal(t, sdk.ObjectTypeTable, objectType)
	})

	t.Run("column object identifier", func(t *testing.T) {
		in := map[string]interface{}{
			"tag_id":      "\"test_db\".\"test_schema\".\"test_tag\"",
			"object_type": "COLUMN",
			"object_identifier": []interface{}{map[string]interface{}{
				"name":     "test_table.test_column",
				"database": "test_db",
				"schema":   "test_schema",
			}},
		}
		d := schema.TestResourceDataRaw(t, resources.TagAssociation().Schema, in)
		tid, identifiers, objectType := resources.TagIdentifierAndObjectIdentifier(d)
		assert.Equal(t, sdk.NewSchemaObjectIdentifier("test_db", "test_schema", "test_tag"), tid)
		assert.Len(t, identifiers, 1)
		assert.Equal(t, "\"test_db\".\"test_schema\".\"test_table\".\"test_column\"", identifiers[0].FullyQualifiedName())
		assert.Equal(t, sdk.ObjectTypeColumn, objectType)
	})
}
