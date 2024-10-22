package resources

import (
	"context"
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
	"log"
	"strconv"
	"strings"
)

type objectRenamingDatabaseListItem struct {
	Name   string
	String string
	Int    int
}

func mapObjectRenamingDatabaseListItemFromValue(items []cty.Value) []objectRenamingDatabaseListItem {
	return collections.Map(items, func(item cty.Value) objectRenamingDatabaseListItem {
		intValue, _ := item.AsValueMap()["int"].AsBigFloat().Int64()
		var name string
		if nameValue, ok := item.AsValueMap()["name"]; ok && !nameValue.IsNull() {
			name = nameValue.AsString()
		}
		return objectRenamingDatabaseListItem{
			Name:   name,
			String: item.AsValueMap()["string"].AsString(),
			Int:    int(intValue),
		}
	})
}

type objectRenamingDatabaseOrderedListItem struct {
	Name  string
	Order string
}

func mapObjectRenamingDatabaseOrderedListItemFromValue(items []cty.Value) []objectRenamingDatabaseOrderedListItem {
	return collections.Map(items, func(item cty.Value) objectRenamingDatabaseOrderedListItem {
		var order string
		if orderValue, ok := item.AsValueMap()["order"]; ok && !orderValue.IsNull() {
			order = orderValue.AsString()
		}
		var name string
		if nameValue, ok := item.AsValueMap()["name"]; ok && !nameValue.IsNull() {
			name = nameValue.AsString()
		}
		return objectRenamingDatabaseOrderedListItem{
			Name:  name,
			Order: order,
		}
	})
}

func objectRenamingDatabaseOrderedListFromSchema(list []any) []objectRenamingDatabaseOrderedListItem {
	objectRenamingDatabaseOrderedListItems := make([]objectRenamingDatabaseOrderedListItem, len(list))
	for index, item := range list {
		var name string
		if nameValue, ok := item.(map[string]any)["name"]; ok {
			name = nameValue.(string)
		}
		objectRenamingDatabaseOrderedListItems[index] = objectRenamingDatabaseOrderedListItem{
			Name:  name,
			Order: strconv.Itoa(index),
		}
	}
	return objectRenamingDatabaseOrderedListItems
}

type ObjectRenamingDatabaseManuallyOrderedListItem struct {
	Name string
}

func objectRenamingDatabaseManuallyOrderedListFromSchema(list []any) []ObjectRenamingDatabaseManuallyOrderedListItem {
	objectRenamingDatabaseOrderedListItems := make([]ObjectRenamingDatabaseManuallyOrderedListItem, len(list))
	for index, item := range list {
		var name string
		if nameValue, ok := item.(map[string]any)["name"]; ok {
			name = nameValue.(string)
		}
		objectRenamingDatabaseOrderedListItems[index] = ObjectRenamingDatabaseManuallyOrderedListItem{
			Name: name,
		}
	}
	return objectRenamingDatabaseOrderedListItems
}

type objectRenamingDatabase struct {
	List                []objectRenamingDatabaseListItem
	OrderedList         []objectRenamingDatabaseOrderedListItem
	ManuallyOrderedList []ObjectRenamingDatabaseManuallyOrderedListItem
}

var ObjectRenamingDatabaseInstance = new(objectRenamingDatabase)

var objectRenamingListsAndSetsSchema = map[string]*schema.Schema{
	// The list field was tested to be used in places where the order of the items should be ignored.
	// It was ignored by comparing hashes of the items to see if any changes were made on the items themselves
	// (if the hashes before and after were the same, we know that nothing was changed, only the order).
	// Also, it doesn't fully support repeating items. This is because they have the same hash and to fully support it,
	// hash counting could be added (counting if the same hash occurs in state and config the number of times, otherwise cause update).
	// Modifications of the items will still cause remove/add behavior.
	"list": {
		Optional:         true,
		Type:             schema.TypeList,
		DiffSuppressFunc: ignoreListOrderAfterFirstApply("list"),
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"string": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"int": {
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	},
	// The manually_ordered_list focuses on providing both aspects:
	//  - Immunity to item reordering after create.
	//  - Handling updates for changed items instead of removing the old item and adding a new one.
	// It does it by providing the required order field that represents what should be the actual order of items
	// on the Snowflake side. The order is ignored on the DiffSuppressFunc level, and the item update (renaming the item)
	// is handled in the resource update function. It is assumed the order starts from 1 (list index + 1) and there are no gaps in order,
	// meaning max_order == len(list).
	"manually_ordered_list": {
		Optional:         true,
		Type:             schema.TypeList,
		DiffSuppressFunc: ignoreOrderAfterFirstApplyWithManuallyOrderedList("manually_ordered_list"),
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"order": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	},
	// The ordered_list field was an attempt of making manual work done in manually_ordered_list automatic by making the order field computed.
	// It didn't work because in DiffSuppressFunc it's hard to get the computed value in the "after" state to compare against.
	// Possibly (but with very low probability), the solution could work by introducing a Computed + Optional list that would be managed by a custom diff function.
	// Due to increased complexity, it was left as is and more research was dedicated to manually_ordered_list.
	"ordered_list": {
		Optional:         true,
		Type:             schema.TypeList,
		DiffSuppressFunc: ignoreOrderAfterFirstApplyWithOrderedList("ordered_list"),
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"order": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	},
}

func ObjectRenamingListsAndSets() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateObjectRenamingListsAndSets,
		UpdateContext: UpdateObjectRenamingListsAndSets,
		ReadContext:   ReadObjectRenamingListsAndSets,
		DeleteContext: DeleteObjectRenamingListsAndSets,

		Schema: objectRenamingListsAndSetsSchema,
	}
}

func CreateObjectRenamingListsAndSets(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ObjectRenamingDatabaseInstance.List = collections.Map(d.Get("list").([]any), func(item any) objectRenamingDatabaseListItem {
		var name string
		if nameValue, ok := item.(map[string]any)["name"]; ok {
			name = nameValue.(string)
		}
		return objectRenamingDatabaseListItem{
			Name:   name,
			String: item.(map[string]any)["string"].(string),
			Int:    item.(map[string]any)["int"].(int),
		}
	})

	ObjectRenamingDatabaseInstance.OrderedList = objectRenamingDatabaseOrderedListFromSchema(d.Get("ordered_list").([]any))
	ObjectRenamingDatabaseInstance.ManuallyOrderedList = objectRenamingDatabaseManuallyOrderedListFromSchema(d.Get("manually_ordered_list").([]any))

	d.SetId("identifier")

	return ReadObjectRenamingListsAndSets(ctx, d, meta)
}

func UpdateObjectRenamingListsAndSets(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if d.HasChange("list") {
		// It wasn't working with d.getChange(). It was returning null elements in one of the test case's steps.
		oldList := d.GetRawState().AsValueMap()["list"].AsValueSlice()
		newList := d.GetRawConfig().AsValueMap()["list"].AsValueSlice()

		oldListMapped := mapObjectRenamingDatabaseListItemFromValue(oldList)
		newListMapped := mapObjectRenamingDatabaseListItemFromValue(newList)

		addedItems, removedItems := ListDiff(oldListMapped, newListMapped)

		for _, removedItem := range removedItems {
			ObjectRenamingDatabaseInstance.List = slices.DeleteFunc(ObjectRenamingDatabaseInstance.List, func(item objectRenamingDatabaseListItem) bool {
				return item == removedItem
			})
		}

		for _, addedItem := range addedItems {
			ObjectRenamingDatabaseInstance.List = append(ObjectRenamingDatabaseInstance.List, addedItem)
		}
	}

	if d.HasChange("ordered_list") {
		oldOrderedList := d.GetRawState().AsValueMap()["ordered_list"].AsValueSlice()
		newOrderedList := d.GetRawConfig().AsValueMap()["ordered_list"].AsValueSlice()
		oldOrderedListMapped := mapObjectRenamingDatabaseOrderedListItemFromValue(oldOrderedList)
		newOrderedListMapped := mapObjectRenamingDatabaseOrderedListItemFromValue(newOrderedList)

		itemsToAdd, itemsToRemove := ListDiff(oldOrderedListMapped, newOrderedListMapped)
		for _, removedItem := range itemsToRemove {
			ObjectRenamingDatabaseInstance.OrderedList = slices.DeleteFunc(ObjectRenamingDatabaseInstance.OrderedList, func(item objectRenamingDatabaseOrderedListItem) bool {
				return item == removedItem
			})
		}

		for _, addedItem := range itemsToAdd {
			addedItem.Order = ""
			ObjectRenamingDatabaseInstance.OrderedList = append(ObjectRenamingDatabaseInstance.OrderedList, addedItem)
		}

		// The implementation is not complete due to mentioned issues with computed order in DiffSuppressFunc
	}

	if d.HasChange("manually_ordered_list") {
		_, newManuallyOrderedList := d.GetChange("manually_ordered_list")
		hasNewItems := len(newManuallyOrderedList.([]any)) > len(ObjectRenamingDatabaseInstance.ManuallyOrderedList)

		// Copy the value from external source and compare with the new state to create the final state that will be saved
		finalState := ObjectRenamingDatabaseInstance.ManuallyOrderedList

		for index, item := range ObjectRenamingDatabaseInstance.ManuallyOrderedList {
			order := index + 1
			newValue, err := collections.FindFirst(newManuallyOrderedList.([]any), func(item any) bool {
				return item.(map[string]any)["order"].(int) == order
			})

			// Items were removed
			if errors.Is(err, collections.ErrObjectNotFound) {
				finalState = finalState[:index]
				break
			}

			newValueName := (*newValue).(map[string]any)["name"].(string)
			// Rename item (handle item change)
			if newValueName != item.Name {
				log.Printf("[DEBUG] Renaming item at order %d from name %s to %s", order, item.Name, newValueName)
				ObjectRenamingDatabaseInstance.ManuallyOrderedList[index].Name = newValueName
			}
		}

		// Add remaining items
		if hasNewItems {
			maxOrder := len(finalState)
			newMaxOrder := len(newManuallyOrderedList.([]any))
			for order := maxOrder + 1; order <= newMaxOrder; order++ {
				item, err := collections.FindFirst(newManuallyOrderedList.([]any), func(item any) bool {
					return item.(map[string]any)["order"].(int) == order
				})
				if err != nil {
					return diag.FromErr(fmt.Errorf("couldn't find item to add with order %d", order))
				}

				finalState = append(finalState, ObjectRenamingDatabaseManuallyOrderedListItem{
					Name: (*item).(map[string]any)["name"].(string),
				})
			}
		}

		ObjectRenamingDatabaseInstance.ManuallyOrderedList = finalState
	}

	return ReadObjectRenamingListsAndSets(ctx, d, meta)
}

func ReadObjectRenamingListsAndSets(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	list := collections.Map(ObjectRenamingDatabaseInstance.List, func(t objectRenamingDatabaseListItem) map[string]any {
		return map[string]any{
			"name":   t.Name,
			"string": t.String,
			"int":    t.Int,
		}
	})
	if err := d.Set("list", list); err != nil {
		return diag.FromErr(err)
	}

	orderedList := make([]map[string]any, len(ObjectRenamingDatabaseInstance.OrderedList))
	for index, item := range ObjectRenamingDatabaseInstance.OrderedList {
		orderedList[index] = map[string]any{
			"name":  item.Name,
			"order": strconv.Itoa(index),
		}
	}
	if err := d.Set("ordered_list", orderedList); err != nil {
		return diag.FromErr(err)
	}

	manuallyOrderedList := make([]map[string]any, len(ObjectRenamingDatabaseInstance.ManuallyOrderedList))
	for index, item := range ObjectRenamingDatabaseInstance.ManuallyOrderedList {
		manuallyOrderedList[index] = map[string]any{
			"name":  item.Name,
			"order": index + 1, // We start from 1, because there may be some implications with 0 and the fact it's treated as unset value in Terraform
		}
	}
	if err := d.Set("manually_ordered_list", manuallyOrderedList); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func DeleteObjectRenamingListsAndSets(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ObjectRenamingDatabaseInstance.List = nil
	ObjectRenamingDatabaseInstance.OrderedList = nil
	ObjectRenamingDatabaseInstance.ManuallyOrderedList = nil
	d.SetId("")

	return nil
}

func ignoreListOrderAfterFirstApply(parentKey string) schema.SchemaDiffSuppressFunc {
	return func(key string, oldValue string, newValue string, d *schema.ResourceData) bool {
		if strings.HasSuffix(key, ".#") {
			return false
		}

		// Raw state is not null after first apply
		if !d.GetRawState().IsNull() && oldValue != newValue {
			// Parse item index from the key
			keyParts := strings.Split(strings.TrimLeft(key, parentKey+"."), ".")
			if len(keyParts) >= 2 {
				index, err := strconv.Atoi(keyParts[0])
				if err != nil {
					log.Println("[DEBUG] Failed to convert list item index: ", err)
					return false
				}

				// Get the hash of the whole item from config (because it represents new value)
				newItemHash := d.GetRawConfig().AsValueMap()[parentKey].AsValueSlice()[index].Hash()

				newItemWasAlreadyPresent := false

				// Try to find the same hash in the state; if found, the new item was already present, and it only changed place in the list
				for _, oldItem := range d.GetRawState().AsValueMap()[parentKey].AsValueSlice() {
					// Matching hashes indicate the order changed, but the item stayed in the config, so suppress the change
					if oldItem.Hash() == newItemHash {
						newItemWasAlreadyPresent = true
					}
				}

				oldItemIsStillPresent := false

				// Sizes of config and state may not be the same
				if len(d.GetRawState().AsValueMap()[parentKey].AsValueSlice()) > index {
					// Get the hash of the whole item from state (because it represents old value)
					oldItemHash := d.GetRawState().AsValueMap()[parentKey].AsValueSlice()[index].Hash()

					// Try to find the same hash in the config; if found, the old item still exists, but changed its place in the list
					for _, newItem := range d.GetRawConfig().AsValueMap()[parentKey].AsValueSlice() {
						if newItem.Hash() == oldItemHash {
							oldItemIsStillPresent = true
						}
					}
				} else if newItemWasAlreadyPresent {
					// Happens in cases where there's a new item at the end of the list, but it was already present, so do nothing
					return true
				}

				if newItemWasAlreadyPresent && oldItemIsStillPresent {
					return true
				}
			}
		}

		return false
	}
}

func ignoreOrderAfterFirstApplyWithOrderedList(parentKey string) schema.SchemaDiffSuppressFunc {
	return func(key string, oldValue string, newValue string, d *schema.ResourceData) bool {
		if strings.HasSuffix(key, ".#") {
			return false
		}

		// Raw state is not null after first apply
		if !d.GetRawState().IsNull() && oldValue != newValue {
			// Parse item index from the key
			keyParts := strings.Split(strings.TrimLeft(key, parentKey+"."), ".")
			if len(keyParts) >= 2 {
				index, err := strconv.Atoi(keyParts[0])
				if err != nil {
					log.Println("[DEBUG] Failed to convert list item index: ", err)
					return false
				}

				newItem := d.GetRawConfig().AsValueMap()[parentKey].AsValueSlice()[index]

				newItemOrder := -1
				// The new order value cannot be retrieved because it's not set on config level.
				// There's also no other way (most likely) to get the newly computed order value for a given item,
				// making this approach not possible.
				newItemOrderValue := newItem.AsValueMap()["order"]
				if !newItemOrderValue.IsNull() {
					newItemOrder, _ = strconv.Atoi(newItemOrderValue.AsString())
				}

				_ = newItemOrder
			}
		}

		return false
	}
}

func ignoreOrderAfterFirstApplyWithManuallyOrderedList(parentKey string) schema.SchemaDiffSuppressFunc {
	return func(key string, oldValue string, newValue string, d *schema.ResourceData) bool {
		if strings.HasSuffix(key, ".#") {
			return false
		}

		// Waw state is not null after first apply
		if !d.GetRawState().IsNull() && oldValue != newValue {
			// Parse item index from the key
			keyParts := strings.Split(strings.TrimLeft(key, parentKey+"."), ".")
			if len(keyParts) >= 2 {
				index, err := strconv.Atoi(keyParts[0])
				if err != nil {
					log.Println("[DEBUG] Failed to convert list item index: ", err)
					return false
				}

				newItems := d.GetRawConfig().AsValueMap()[parentKey].AsValueSlice()
				if len(newItems) <= index {
					// item was removed
					return false
				}

				newItem := newItems[index]
				itemWasAlreadyPresent := false
				itemIsStillPresent := false

				var newItemOrder int64 = -1
				newItemOrderValue := newItem.AsValueMap()["order"]
				if !newItemOrderValue.IsNull() {
					newItemOrder, _ = newItemOrderValue.AsBigFloat().Int64()
				}

				// That's a new item
				if newItemOrder == -1 {
					return false
				} else { // Else it was already present, but we need to check the hash
					for _, oldItem := range d.GetRawState().AsValueMap()[parentKey].AsValueSlice() {
						oldItemOrder, _ := oldItem.AsValueMap()["order"].AsBigFloat().Int64()
						if oldItemOrder == newItemOrder {
							if oldItem.Hash() != newItem.Hash() {
								// The item has the same order, but the values in other fields changed (different hash)
								return false
							} else {
								itemWasAlreadyPresent = true
							}
						}
					}
				}

				// Check if a new item is indexable (with new items added at the end, it's not possible to index state value for them, because it doesn't exist yet)
				if len(d.GetRawState().AsValueMap()[parentKey].AsValueSlice()) > index {
					oldItem := d.GetRawState().AsValueMap()[parentKey].AsValueSlice()[index]
					oldItemOrder, _ := oldItem.AsValueMap()["order"].AsBigFloat().Int64()

					// Check if this order is still present
					for _, newItem := range d.GetRawConfig().AsValueMap()[parentKey].AsValueSlice() {
						newItemOrder, _ := newItem.AsValueMap()["order"].AsBigFloat().Int64()
						if oldItemOrder == newItemOrder {
							if oldItem.Hash() != newItem.Hash() {
								// The order is still present, but the values in other fields changed (different hash)
								return false
							} else {
								itemIsStillPresent = true
							}
						}
					}
				}

				if itemWasAlreadyPresent && itemIsStillPresent {
					return true
				}
			}
		}

		return false
	}
}
