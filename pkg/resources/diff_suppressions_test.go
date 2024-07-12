package resources_test

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func Test_NormalizeAndCompare(t *testing.T) {
	genericNormalize := func(value string) (any, error) {
		switch value {
		case "ok", "ok1":
			return "ok", nil
		default:
			return nil, fmt.Errorf("incorrect value %s", value)
		}
	}

	t.Run("generic normalize", func(t *testing.T) {
		result := resources.NormalizeAndCompare(genericNormalize)("", "ok", "ok", nil)
		assert.True(t, result)

		result = resources.NormalizeAndCompare(genericNormalize)("", "ok", "ok1", nil)
		assert.True(t, result)

		result = resources.NormalizeAndCompare(genericNormalize)("", "ok", "nok", nil)
		assert.False(t, result)
	})

	t.Run("warehouse size", func(t *testing.T) {
		result := resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", string(sdk.WarehouseSizeX4Large), string(sdk.WarehouseSizeX4Large), nil)
		assert.True(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", string(sdk.WarehouseSizeX4Large), "4X-LARGE", nil)
		assert.True(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", string(sdk.WarehouseSizeX4Large), string(sdk.WarehouseSizeX5Large), nil)
		assert.False(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", string(sdk.WarehouseSizeX4Large), "invalid", nil)
		assert.False(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", string(sdk.WarehouseSizeX4Large), "", nil)
		assert.False(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", "invalid", string(sdk.WarehouseSizeX4Large), nil)
		assert.False(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", "", string(sdk.WarehouseSizeX4Large), nil)
		assert.False(t, result)
	})
}

func Test_IgnoreAfterCreation(t *testing.T) {
	testSchema := map[string]*schema.Schema{
		"value": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	t.Run("without id", func(t *testing.T) {
		in := map[string]any{}
		d := schema.TestResourceDataRaw(t, testSchema, in)

		result := resources.IgnoreAfterCreation("", "", "", d)
		assert.False(t, result)
	})

	t.Run("with id", func(t *testing.T) {
		in := map[string]any{}
		d := schema.TestResourceDataRaw(t, testSchema, in)
		d.SetId("something")

		result := resources.IgnoreAfterCreation("", "", "", d)
		assert.True(t, result)
	})
}

func Test_NormalizeAndCompareIdentifiersSet(t *testing.T) {
	rawDataWithValues := func(values []any) *schema.ResourceData {
		return schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"value": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		}, map[string]any{
			"value": values,
		})
	}
	emptyResourceData := rawDataWithValues([]any{})

	t.Run("validation: size key", func(t *testing.T) {
		assert.False(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.#", "1", "2", emptyResourceData))
	})

	t.Run("validation: case mismatch", func(t *testing.T) {
		resourceData := rawDataWithValues([]any{"SCHEMA.OBJECT.IDENTIFIER"})
		assert.False(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "schema.object.identifier", "", resourceData))
		// TODO: Cannot be tested with schema.TestResourceDataRaw because it doesn't populate raw state which is used in the cases below
		// assert.False(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "", `"schema"."object"."identifier"`, resourceData))
	})

	t.Run("validation: case mismatch quoted identifier in the state", func(t *testing.T) {
		resourceData := rawDataWithValues([]any{`"SCHEMA"."OBJECT"."IDENTIFIER"`})
		assert.False(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "schema.object.identifier", "", resourceData))
		// TODO: Cannot be tested with schema.TestResourceDataRaw because it doesn't populate raw state which is used in the cases below
		// assert.False(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "", `"schema"."object"."identifier"`, resourceData))
	})

	t.Run(`change detected from schema.object.identifier to "schema"."object"."identifier" with schema.object.identifier in state`, func(t *testing.T) {
		resourceData := rawDataWithValues([]any{"schema.object.identifier"})
		assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "schema.object.identifier", "", resourceData))
		// TODO: Cannot be tested with schema.TestResourceDataRaw because it doesn't populate raw state which is used in the cases below
		// assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "", `"schema"."object"."identifier"`, resourceData))
	})

	t.Run(`change detected from schema.object.identifier to "schema"."object"."identifier" with "schema"."object"."identifier" in state`, func(t *testing.T) {
		resourceData := rawDataWithValues([]any{`"schema"."object"."identifier"`})
		assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "schema.object.identifier", "", resourceData))
		// TODO: Cannot be tested with schema.TestResourceDataRaw because it doesn't populate raw state which is used in the cases below
		// assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "", `"schema"."object"."identifier"`, resourceData))
	})

	t.Run(`change detected from "schema"."object"."identifier" to schema.object.identifier with schema.object.identifier in state`, func(t *testing.T) {
		resourceData := rawDataWithValues([]any{"schema.object.identifier"})
		assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "schema.object.identifier", "", resourceData))
		// TODO: Cannot be tested with schema.TestResourceDataRaw because it doesn't populate raw state which is used in the cases below
		// assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "", `"schema"."object"."identifier"`, resourceData))
	})

	t.Run(`change detected from "schema"."object"."identifier" to schema.object.identifier with "schema"."object"."identifier" in state`, func(t *testing.T) {
		resourceData := rawDataWithValues([]any{`"schema"."object"."identifier"`})
		assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "schema.object.identifier", "", resourceData))
		// TODO: Cannot be tested with schema.TestResourceDataRaw because it doesn't populate raw state which is used in the cases below
		// assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "", `"schema"."object"."identifier"`, resourceData))
	})

	t.Run(`change detected from "schema"."object"."identifier" to schema.object.identifier with "schema"."object"."identifier" in state uppercased`, func(t *testing.T) {
		resourceData := rawDataWithValues([]any{`"SCHEMA"."OBJECT"."IDENTIFIER"`})
		assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "SCHEMA.OBJECT.IDENTIFIER", "", resourceData))
		// TODO: Cannot be tested with schema.TestResourceDataRaw because it doesn't populate raw state which is used in the cases below
		// assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "", `"SCHEMA"."OBJECT"."IDENTIFIER"`, resourceData))
	})
}
