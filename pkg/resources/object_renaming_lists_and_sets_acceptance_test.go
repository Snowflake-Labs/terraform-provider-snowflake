package resources_test

import (
	"fmt"
	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"testing"
)

func TestAcc_BasicListFlow(t *testing.T) {
	// TODO: _ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "111", "int": 111},
					{"string": "222", "int": 222},
					{"string": "333", "int": 333},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
						{"string": "111", "int": "111"},
						{"string": "222", "int": "222"},
						{"string": "333", "int": "333"},
					}),
				),
			},
			// Remove, shift, and add one item (in the middle)
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "222", "int": 222},
					{"string": "444", "int": 444},
					{"string": "333", "int": 333},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
						{"string": "222", "int": "222"},
						{"string": "444", "int": "444"},
						{"string": "333", "int": "333"},
					}),
				),
			},
			// Remove, shift, and add one item (at the end)
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "444", "int": 444},
					{"string": "333", "int": 333},
					{"string": "111", "int": 111},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
						{"string": "444", "int": "444"},
						{"string": "333", "int": "333"},
						{"string": "111", "int": "111"},
					}),
				),
			},
			// Remove, shift, and add one item (at the beginning)
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "222", "int": 222},
					{"string": "333", "int": 333},
					{"string": "444", "int": 444},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
						{"string": "222", "int": "222"},
						{"string": "333", "int": "333"},
						{"string": "444", "int": "444"},
					}),
				),
			},
			// Reorder items and add one
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "444", "int": 444},
					{"string": "555", "int": 555},
					{"string": "333", "int": 333},
					{"string": "222", "int": 222},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
						{"string": "444", "int": "444"},
						{"string": "555", "int": "555"},
						{"string": "333", "int": "333"},
						{"string": "222", "int": "222"},
					}),
				),
			},
			// Replace all items
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "444", "int": 444},
					{"string": "555", "int": 555},
					{"string": "333", "int": 333},
					{"string": "222", "int": 222},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
						{"string": "444", "int": "444"},
						{"string": "555", "int": "555"},
						{"string": "333", "int": "333"},
						{"string": "222", "int": "222"},
					}),
				),
			},
			//// Add another item at the beginning
			//{
			//	Config: objectRenamingConfigList([]map[string]any{
			//		{"string": "555", "int": 555},
			//		{"string": "222", "int": 222},
			//		{"string": "444", "int": 444},
			//		{"string": "333", "int": 333},
			//	}),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
			//			{"string": "555", "int": "555"},
			//			{"string": "222", "int": "222"},
			//			{"string": "444", "int": "444"},
			//			{"string": "333", "int": "333"},
			//		}),
			//	),
			//},
			//// Reorder items
			//{
			//	Config: objectRenamingConfigList([]map[string]any{
			//		{"string": "222", "int": 222},
			//		{"string": "555", "int": 555},
			//		{"string": "333", "int": 333},
			//		{"string": "444", "int": 444},
			//	}),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
			//			{"string": "222", "int": "222"},
			//			{"string": "555", "int": "555"},
			//			{"string": "333", "int": "333"},
			//			{"string": "444", "int": "444"},
			//		}),
			//	),
			//},
			//// Replace some items
			//{
			//	Config: objectRenamingConfigList([]map[string]any{
			//		{"string": "111", "int": 111},
			//		{"string": "555", "int": 555},
			//		{"string": "666", "int": 666},
			//		{"string": "444", "int": 444},
			//	}),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
			//			{"string": "111", "int": "111"},
			//			{"string": "555", "int": "555"},
			//			{"string": "666", "int": "666"},
			//			{"string": "444", "int": "444"},
			//		}),
			//	),
			//},
			//// Replace all items
			//{
			//	Config: objectRenamingConfigList([]map[string]any{
			//		{"string": "222", "int": 222},
			//		{"string": "333", "int": 333},
			//		{"string": "777", "int": 777},
			//		{"string": "888", "int": 888},
			//	}),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
			//			{"string": "222", "int": "222"},
			//			{"string": "333", "int": "333"},
			//			{"string": "777", "int": "777"},
			//			{"string": "888", "int": "888"},
			//		}),
			//	),
			//},
			//// Add new item at the end
			//{
			//	Config: objectRenamingConfigList([]map[string]any{
			//		{"string": "222", "int": 222},
			//		{"string": "333", "int": 333},
			//		{"string": "777", "int": 777},
			//		{"string": "888", "int": 888},
			//		{"string": "999", "int": 999},
			//	}),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
			//			{"string": "222", "int": "222"},
			//			{"string": "333", "int": "333"},
			//			{"string": "777", "int": "777"},
			//			{"string": "888", "int": "888"},
			//			{"string": "999", "int": "999"},
			//		}),
			//	),
			//},
		},
	})
}

func objectRenamingConfigList(listItems []map[string]any) string {
	generateListItem := func(s string, i int) string {
		return fmt.Sprintf(`
list {
	string = "%[1]s"
	int = %[2]d
}
`, s, i)
	}

	generatedListItems := ""
	for _, item := range listItems {
		generatedListItems += generateListItem(item["string"].(string), item["int"].(int))
	}

	return fmt.Sprintf(`
resource "snowflake_object_renaming" "test" {
	%s
}
`, generatedListItems)
}
