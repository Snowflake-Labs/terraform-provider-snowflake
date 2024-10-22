package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_BasicListFlow(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)
	acc.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: objectRenamingConfigList([]map[string]any{
					{"name": "", "string": "111", "int": 111},
					{"name": "", "string": "222", "int": 222},
					{"name": "", "string": "333", "int": 333},
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
					{"name": "", "string": "222", "int": 222},
					{"name": "", "string": "444", "int": 444},
					{"name": "", "string": "333", "int": 333},
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
					{"name": "", "string": "444", "int": 444},
					{"name": "", "string": "333", "int": 333},
					{"name": "", "string": "111", "int": 111},
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
					{"name": "", "string": "222", "int": 222},
					{"name": "", "string": "333", "int": 333},
					{"name": "", "string": "444", "int": 444},
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
					{"name": "", "string": "444", "int": 444},
					{"name": "", "string": "555", "int": 555},
					{"name": "", "string": "333", "int": 333},
					{"name": "", "string": "222", "int": 222},
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
					{"name": "", "string": "444", "int": 444},
					{"name": "", "string": "555", "int": 555},
					{"name": "", "string": "333", "int": 333},
					{"name": "", "string": "222", "int": 222},
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
			// Add an item that is identical to another one (currently, failing)
			// {
			//	Config: objectRenamingConfigList([]map[string]any{
			//		{"name": "", "string": "444", "int": 444},
			//		{"name": "", "string": "555", "int": 555},
			//		{"name": "", "string": "333", "int": 333},
			//		{"name": "", "string": "222", "int": 222},
			//		{"name": "", "string": "222", "int": 222},
			//	}),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		assert.HasListItemsOrderIndependent("snowflake_object_renaming.test", "list", []map[string]string{
			//			{"name": "", "string": "444", "int": "444"},
			//			{"name": "", "string": "555", "int": "555"},
			//			{"name": "", "string": "333", "int": "333"},
			//			{"name": "", "string": "222", "int": "222"},
			//			{"name": "", "string": "222", "int": "222"},
			//		}),
			//	),
			// },
		},
	})
}

// This test researches the possibility of performing update instead of remove + add item
func TestAcc_ListNameUpdate(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)
	acc.TestAccPreCheck(t)

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
			// It's hard to handle reorder + rename with this approach,
			// because without any additional metadata, we cannot identify a given list item.
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

func TestAcc_ListsWithDuplicatedItems(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)
	acc.TestAccPreCheck(t)

	// Fails, because the SuppressDiffFunc works on the hash of individual items.
	// To correctly suppress such changes, the number of repeated hashes should be counted.
	t.Skip("Currently failing, because duplicated hashes are not supported.")

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

func TestAcc_BasicManuallyOrderedListFlow(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableObjectRenamingTest)
	acc.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "name1", "order": 1},
					{"name": "name2", "order": 2},
					{"name": "name3", "order": 3},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.0.name", "name1"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.1.name", "name2"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.2.name", "name3"),
				),
			},
			// Change values
			{
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "name2", "order": 1},
					{"name": "name2", "order": 2},
					{"name": "name4", "order": 3},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.0.name", "name2"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.1.name", "name2"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.2.name", "name4"),
				),
			},
			// Change the order
			{
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "name4", "order": 3},
					{"name": "name2", "order": 2},
					{"name": "name2", "order": 1},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.0.name", "name2"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.1.name", "name2"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.2.name", "name4"),
				),
			},
			// change order and values
			{
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "name1", "order": 2},
					{"name": "name2", "order": 3},
					{"name": "name3", "order": 1},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.0.name", "name3"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.1.name", "name1"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.2.name", "name2"),
				),
			},
			// Add items
			{
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "name4", "order": 4}, // added
					{"name": "name1", "order": 2},
					{"name": "name5", "order": 5}, // added
					{"name": "name2", "order": 3},
					{"name": "name3", "order": 1},
					{"name": "name6", "order": 6}, // added
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.0.name", "name3"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.1.name", "name1"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.2.name", "name2"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.3.name", "name4"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.4.name", "name5"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.5.name", "name6"),
				),
			},
			// Remove items
			// The removal is kind of remove + update when you remove from the middle of the list,
			// because some items have to jump into place of removed ones.
			// With the Terraform SDKv2, it would be hard to achieve "shift after removal" functionality without assuming that,
			// e.g., none of the items were changed between the states. The topic is further discussed in the research documentation.
			{
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "name4", "order": 2},
					{"name": "name5", "order": 3},
					{"name": "name2", "order": 1},
					{"name": "name6", "order": 4},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.0.name", "name2"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.1.name", "name4"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.2.name", "name5"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.3.name", "name6"),
				),
			},
			// Change externally
			{
				PreConfig: func() {
					resources.ObjectRenamingDatabaseInstance.ManuallyOrderedList[0].Name = "name changed externally"
					resources.ObjectRenamingDatabaseInstance.ManuallyOrderedList[1].Name = "name changed externally"
					resources.ObjectRenamingDatabaseInstance.ManuallyOrderedList[2].Name = "name changed externally"
					resources.ObjectRenamingDatabaseInstance.ManuallyOrderedList[3].Name = "name changed externally"
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "name4", "order": 2},
					{"name": "name5", "order": 3},
					{"name": "name2", "order": 1},
					{"name": "name6", "order": 4},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.0.name", "name2"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.1.name", "name4"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.2.name", "name5"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.3.name", "name6"),
				),
			},
			// Add externally
			{
				PreConfig: func() {
					resources.ObjectRenamingDatabaseInstance.ManuallyOrderedList = append(resources.ObjectRenamingDatabaseInstance.ManuallyOrderedList, resources.ObjectRenamingDatabaseManuallyOrderedListItem{Name: "name added externally"})
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "name4", "order": 2},
					{"name": "name5", "order": 3},
					{"name": "name2", "order": 1},
					{"name": "name6", "order": 4},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.0.name", "name2"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.1.name", "name4"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.2.name", "name5"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.3.name", "name6"),
				),
			},
			// Removed externally
			{
				PreConfig: func() {
					resources.ObjectRenamingDatabaseInstance.ManuallyOrderedList = resources.ObjectRenamingDatabaseInstance.ManuallyOrderedList[1:]
				},
				Config: objectRenamingConfigManuallyOrderedList([]map[string]any{
					{"name": "name4", "order": 2},
					{"name": "name5", "order": 3},
					{"name": "name2", "order": 1},
					{"name": "name6", "order": 4},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.0.name", "name2"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.1.name", "name4"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.2.name", "name5"),
					resource.TestCheckResourceAttr("snowflake_object_renaming.test", "manually_ordered_list.3.name", "name6"),
				),
			},
		},
	})
}

func objectRenamingConfigManuallyOrderedList(listItems []map[string]any) string {
	generateListItem := func(name string, order int) string {
		return fmt.Sprintf(`
manually_ordered_list {
	name = "%s"
	order = %d
}
`, name, order)
	}

	generatedListItems := ""
	for _, item := range listItems {
		generatedListItems += generateListItem(item["name"].(string), item["order"].(int))
	}

	return fmt.Sprintf(`
resource "snowflake_object_renaming" "test" {
	%s
}
`, generatedListItems)
}
