package resources_test

import (
	"fmt"
	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"testing"
)

func TestAcc_BasicListFlow(t *testing.T) {
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
						{"name": "", "string": "111", "int": "111"},
						{"name": "", "string": "222", "int": "222"},
						{"name": "", "string": "333", "int": "333"},
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
						{"name": "", "string": "222", "int": "222"},
						{"name": "", "string": "444", "int": "444"},
						{"name": "", "string": "333", "int": "333"},
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
						{"name": "", "string": "444", "int": "444"},
						{"name": "", "string": "333", "int": "333"},
						{"name": "", "string": "111", "int": "111"},
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
						{"name": "", "string": "222", "int": "222"},
						{"name": "", "string": "333", "int": "333"},
						{"name": "", "string": "444", "int": "444"},
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
						{"name": "", "string": "444", "int": "444"},
						{"name": "", "string": "555", "int": "555"},
						{"name": "", "string": "333", "int": "333"},
						{"name": "", "string": "222", "int": "222"},
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
						{"name": "", "string": "444", "int": "444"},
						{"name": "", "string": "555", "int": "555"},
						{"name": "", "string": "333", "int": "333"},
						{"name": "", "string": "222", "int": "222"},
					}),
				),
			},
		},
	})
}

func TestAcc_ListUpdatesWithNullAttribute(t *testing.T) {
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
					{"name": nil, "string": "111", "int": 111},
					{"name": nil, "string": "222", "int": 222},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
						{"name": "", "string": "111", "int": "111"},
						{"name": "", "string": "222", "int": "222"},
					}),
				),
			},
			// Update one item
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"name": nil, "string": "111", "int": 111},
					{"name": nil, "string": "333", "int": 333},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
						{"name": "", "string": "111", "int": "111"},
						{"name": "", "string": "333", "int": "333"},
					}),
				),
			},
		},
	})
}

// This test researches the possibility of performing update instead of remove + add item
func TestAcc_ListNameUpdate(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	// TODO: See how reordering items will impact the update capabilities (because with diff suppress it may be hard to update)
	// TODO: Test with external changes

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"name": "column1", "string": "111", "int": 111},
					{"name": "column2", "string": "222", "int": 222},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
						{"name": "column1", "string": "111", "int": "111"},
						{"name": "column2", "string": "222", "int": "222"},
					}),
				),
			},
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"name": "column2", "string": "222", "int": 222},
					{"name": "column1", "string": "111", "int": 111},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
						{"name": "column2", "string": "222", "int": "222"},
						{"name": "column1", "string": "111", "int": "111"},
					}),
				),
			},
			// TODO: It's hard combo to handle (reorder + rename) with this approach (because without any additional metadata we cannot identify a given list item)
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"name": "column3", "string": "222", "int": 222},
					{"name": "column1", "string": "111", "int": 111},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
						{"name": "column3", "string": "222", "int": "222"},
						{"name": "column1", "string": "111", "int": "111"},
					}),
				),
			},
		},
	})
}

// Because the list diff suppressor works on the hash of individual items, this test checks how it behaves with duplicates that can occur in lists
// TODO: Currently, failing come back to this later
func TestAcc_ListsWithDuplicatedItems(t *testing.T) {
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
			// Introduce duplicates
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"string": "111", "int": 111},
					{"string": "111", "int": 111},
					{"string": "222", "int": 222},
					{"string": "222", "int": 222},
					{"string": "333", "int": 333},
					{"string": "333", "int": 333},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
						{"string": "111", "int": "111"},
						{"string": "111", "int": "111"},
						{"string": "222", "int": "222"},
						{"string": "222", "int": "222"},
						{"string": "333", "int": "333"},
						{"string": "333", "int": "333"},
					}),
				),
			},
		},
	})
}

func objectRenamingConfigList(listItems []map[string]any) string {
	generateListItem := func(name *string, s string, i int) string {
		var nameField string
		if name != nil {
			nameField = fmt.Sprintf(`name = "%s"`, *name)
		}
		return fmt.Sprintf(`
list {
	%[1]s
	string = "%[2]s"
	int = %[3]d
}
`, nameField, s, i)
	}

	generatedListItems := ""
	for _, item := range listItems {
		var name *string
		if nameValue, ok := item["name"]; ok && nameValue != nil {
			name = sdk.String(nameValue.(string))
		}
		generatedListItems += generateListItem(name, item["string"].(string), item["int"].(int))
	}

	return fmt.Sprintf(`
resource "snowflake_object_renaming" "test" {
	%s
}
`, generatedListItems)
}

func TestAcc_BasicOrderedListFlow(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: objectRenamingConfigOrderedList([]map[string]any{
					{"name": "name1", "string": "111"},
					{"name": "name2", "string": "222"},
					{"name": "name3", "string": "333"},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "ordered_list", []map[string]string{
						{"name": "name1", "string": "111", "order": "0"},
						{"name": "name2", "string": "222", "order": "1"},
						{"name": "name3", "string": "333", "order": "2"},
					}),
				),
			},
			// Change the order
			{
				Config: objectRenamingConfigOrderedList([]map[string]any{
					{"name": "name3", "string": "333"},
					{"name": "name1", "string": "111"},
					{"name": "name2", "string": "222"},
				}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_object_renaming.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "ordered_list", []map[string]string{
						{"name": "name3", "string": "333", "order": "0"},
						{"name": "name1", "string": "111", "order": "1"},
						{"name": "name2", "string": "222", "order": "2"},
					}),
				),
			},
		},
	})
}

func objectRenamingConfigOrderedList(listItems []map[string]any) string {
	generateListItem := func(name *string, s string) string {
		var nameField string
		if name != nil {
			nameField = fmt.Sprintf(`name = "%s"`, *name)
		}
		return fmt.Sprintf(`
ordered_list {
	%[1]s
	string = "%[2]s"
}
`, nameField, s)
	}

	generatedListItems := ""
	for _, item := range listItems {
		var name *string
		if nameValue, ok := item["name"]; ok && nameValue != nil {
			name = sdk.String(nameValue.(string))
		}
		generatedListItems += generateListItem(name, item["string"].(string))
	}

	return fmt.Sprintf(`
resource "snowflake_object_renaming" "test" {
	%s
}
`, generatedListItems)
}
