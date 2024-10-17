package resources

import (
	"context"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
	"log"
	"strconv"
	"strings"
)

var objectRenamingListsAndSetsSchema = map[string]*schema.Schema{
	"list": {
		Optional:         true,
		Type:             schema.TypeList,
		DiffSuppressFunc: ignoreOrderAfterFirstApply("list"),
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
				// TODO: Add more if needed (e.g. with hashing or smth)
			},
		},
	},
	"ordered_list": {
		Computed:         true,
		Optional:         true,
		Type:             schema.TypeList,
		DiffSuppressFunc: ignoreOrderAfterFirstApplyWithOrderedList("ordered_list"),
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
				"order": {
					Type:     schema.TypeString, // because 0 is a valid index
					Computed: true,
				},
			},
		},
	},
	//"set": {},
}

// TODO: This may fail for lists where items with the same values are allowed
//   - This could be potentially solved by comparing the number of hash repetitions in state vs config
func ignoreOrderAfterFirstApply(parentKey string) schema.SchemaDiffSuppressFunc {
	return func(key string, oldValue string, newValue string, d *schema.ResourceData) bool {
		if strings.HasSuffix(key, ".#") {
			return false
		}

		// raw state is not null after first apply
		if !d.GetRawState().IsNull() && oldValue != newValue {
			// parse item index from the key
			keyParts := strings.Split(strings.TrimLeft(key, parentKey+"."), ".")
			if len(keyParts) >= 2 {
				index, err := strconv.Atoi(keyParts[0])
				if err != nil {
					log.Println("[DEBUG] Failed to convert list item index: ", err)
					return false
				}

				// get the hash of the whole item from config (because it represents new value)
				newItemHash := d.GetRawConfig().AsValueMap()[parentKey].AsValueSlice()[index].Hash()

				newItemWasAlreadyPresent := false

				// try to find the same hash in the state; if found, the new item was already present, and it only changed place in the list
				for _, oldItem := range d.GetRawState().AsValueMap()[parentKey].AsValueSlice() {
					// matching hashes indicate the order changed, but the item stayed in the config, so suppress the change
					if oldItem.Hash() == newItemHash {
						newItemWasAlreadyPresent = true
					}
				}

				oldItemIsStillPresent := false

				// sizes of config and state may not be the same
				if len(d.GetRawState().AsValueMap()[parentKey].AsValueSlice()) > index {
					// get the hash of the whole item from state (because it represents old value)
					oldItemHash := d.GetRawState().AsValueMap()[parentKey].AsValueSlice()[index].Hash()

					// try to find the same hash in the config; if found, the old item still exists, but changed its place in the list
					for _, newItem := range d.GetRawConfig().AsValueMap()[parentKey].AsValueSlice() {
						if newItem.Hash() == oldItemHash {
							oldItemIsStillPresent = true
						}
					}
				} else if newItemWasAlreadyPresent {
					// happens in cases where there's a new item at the end of the list, but it was already present, so do nothing
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

		// raw state is not null after first apply
		if !d.GetRawState().IsNull() && oldValue != newValue {
			// parse item index from the key
			keyParts := strings.Split(strings.TrimLeft(key, parentKey+"."), ".")
			if len(keyParts) >= 2 {
				index, err := strconv.Atoi(keyParts[0])
				if err != nil {
					log.Println("[DEBUG] Failed to convert list item index: ", err)
					return false
				}

				newItem := d.GetRawConfig().AsValueMap()[parentKey].AsValueSlice()[index]

				//// TODO: (?) item is being removed
				//if newItem.AsValueMap()["order"].IsNull() {
				//
				//}

				itemWasAlreadyPresent := false
				itemIsStillPresent := false

				newItemOrder := -1
				newItemOrderValue := newItem.AsValueMap()["order"]
				if !newItemOrderValue.IsNull() {
					newItemOrder, _ = strconv.Atoi(newItemOrderValue.AsString())
				}

				// that's a new item
				if newItemOrder == -1 {
					return false
				} else { // else it was already present, but we need to check the hash
					for _, oldItem := range d.GetRawState().AsValueMap()[parentKey].AsValueSlice() {
						oldItemOrder, _ := strconv.Atoi(oldItem.AsValueMap()["order"].AsString())
						if oldItemOrder == newItemOrder {
							if oldItem.Hash() != newItem.Hash() {
								// the item has the same order, but the values in other fields changed (different hash)
								return false
							} else {
								itemWasAlreadyPresent = true
							}
						}
					}
				}

				// check if a new item is indexable (with new items added at the end, it's not possible to index state value for them, because it doesn't exist yet)
				if len(d.GetRawState().AsValueMap()[parentKey].AsValueSlice()) > index {
					oldItem := d.GetRawState().AsValueMap()[parentKey].AsValueSlice()[index]
					oldItemOrder, _ := strconv.Atoi(oldItem.AsValueMap()["order"].AsString())

					// check if this order is still present
					for _, newItem := range d.GetRawConfig().AsValueMap()[parentKey].AsValueSlice() {
						newItemOrder, _ := strconv.Atoi(newItem.AsValueMap()["order"].AsString())
						if oldItemOrder == newItemOrder {
							if oldItem.Hash() != newItem.Hash() {
								// the order is still present, but the values in other fields changed (different hash)
								return false
							} else {
								itemIsStillPresent = true
							}
						}
					}
				} // this else should be handled by newItemOrder == -1 earlier

				if itemWasAlreadyPresent && itemIsStillPresent { // and items are with the same hashes
					return true
				}
			}
		}

		return false
	}
}

func ObjectRenamingListsAndSets() *schema.Resource {
	return &schema.Resource{
		Description:   "TODO",
		CreateContext: CreateObjectRenamingListsAndSets,
		UpdateContext: UpdateObjectRenamingListsAndSets,
		ReadContext:   ReadObjectRenamingListsAndSets,
		DeleteContext: DeleteObjectRenamingListsAndSets,

		Schema: objectRenamingListsAndSetsSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
			if d.GetRawState().IsNull() {
				return nil
			}

			//if _, ok := d.GetOk("ordered_list.0.order"); ok {
			//	if err := d.SetNew("ordered_list.0.order", 123); err != nil {
			//		return err
			//	}
			//}
			if d.HasChange("ordered_list") {
				oldL, newL := d.GetChange("ordered_list")
				_, _ = oldL, newL
			}

			return nil
		},
	}
}

type objectRenamingDatabaseListItem struct {
	Name   string
	String string
	Int    int
}

type objectRenamingDatabaseOrderedListItem struct {
	Name   string
	String string
	Order  string
}

type objectRenamingDatabase struct {
	List        []objectRenamingDatabaseListItem
	OrderedList []objectRenamingDatabaseOrderedListItem
}

var objectRenamingDatabaseInstance = new(objectRenamingDatabase)

func CreateObjectRenamingListsAndSets(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	objectRenamingDatabaseInstance.List = collections.Map(d.Get("list").([]any), func(item any) objectRenamingDatabaseListItem {
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

	orderedList := d.Get("ordered_list").([]any)
	objectRenamingDatabaseOrderedListItems := make([]objectRenamingDatabaseOrderedListItem, len(orderedList))

	for index, item := range d.Get("ordered_list").([]any) {
		var name string
		if nameValue, ok := item.(map[string]any)["name"]; ok {
			name = nameValue.(string)
		}
		objectRenamingDatabaseOrderedListItems[index] = objectRenamingDatabaseOrderedListItem{
			Name:   name,
			String: item.(map[string]any)["string"].(string),
			Order:  strconv.Itoa(index),
		}
	}
	objectRenamingDatabaseInstance.OrderedList = objectRenamingDatabaseOrderedListItems

	d.SetId("identifier")

	return ReadObjectRenamingListsAndSets(ctx, d, meta)
}

// TODO: For now, replicating columns from table
func UpdateObjectRenamingListsAndSets(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if d.HasChange("list") {
		// TODO: It wasn't working with d.getChange()
		oldList := d.GetRawState().AsValueMap()["list"].AsValueSlice()
		newList := d.GetRawConfig().AsValueMap()["list"].AsValueSlice()
		oldListMapped := collections.Map(oldList, func(item cty.Value) objectRenamingDatabaseListItem {
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
		newListMapped := collections.Map(newList, func(item cty.Value) objectRenamingDatabaseListItem {
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

		addedItems, removedItems := ListDiff(oldListMapped, newListMapped)

		for _, removedItem := range removedItems {
			objectRenamingDatabaseInstance.List = slices.DeleteFunc(objectRenamingDatabaseInstance.List, func(item objectRenamingDatabaseListItem) bool {
				return item == removedItem
			})
		}

		for _, addedItem := range addedItems {
			objectRenamingDatabaseInstance.List = append(objectRenamingDatabaseInstance.List, addedItem)
		}
	}

	if d.HasChange("ordered_list") {
		// TODO: try with d.GetChange()
		//oldOrderedList, newOrderedList := d.GetChange("ordered_list")
		oldOrderedList := d.GetRawState().AsValueMap()["ordered_list"].AsValueSlice()
		newOrderedList := d.GetRawConfig().AsValueMap()["ordered_list"].AsValueSlice()
		oldOrderedListMapped := collections.Map(oldOrderedList, func(item cty.Value) objectRenamingDatabaseOrderedListItem {
			var order string
			if orderValue, ok := item.AsValueMap()["order"]; ok && !orderValue.IsNull() {
				order = orderValue.AsString()
			}
			var name string
			if nameValue, ok := item.AsValueMap()["name"]; ok && !nameValue.IsNull() {
				name = nameValue.AsString()
			}
			return objectRenamingDatabaseOrderedListItem{
				Name:   name,
				String: item.AsValueMap()["string"].AsString(),
				Order:  order,
			}
		})
		newOrderedListMapped := collections.Map(newOrderedList, func(item cty.Value) objectRenamingDatabaseOrderedListItem {
			var order string
			if orderValue, ok := item.AsValueMap()["order"]; ok && !orderValue.IsNull() {
				order = orderValue.AsString()
			}
			var name string
			if nameValue, ok := item.AsValueMap()["name"]; ok && !nameValue.IsNull() {
				name = nameValue.AsString()
			}
			return objectRenamingDatabaseOrderedListItem{
				Name:   name,
				String: item.AsValueMap()["string"].AsString(),
				Order:  order,
			}
		})

		itemsToAdd, itemsToRemove := ListDiff(oldOrderedListMapped, newOrderedListMapped)
		for _, removedItem := range itemsToRemove {
			objectRenamingDatabaseInstance.OrderedList = slices.DeleteFunc(objectRenamingDatabaseInstance.OrderedList, func(item objectRenamingDatabaseOrderedListItem) bool {
				return item == removedItem
			})
		}

		for _, addedItem := range itemsToAdd {
			addedItem.Order = ""
			objectRenamingDatabaseInstance.OrderedList = append(objectRenamingDatabaseInstance.OrderedList, addedItem)
		}

		// TODO: Can this be done? (what about reordering)
		// After changes, recompute order
		//for i := range d.Get("ordered_list").([]any) {
		//	if err := d.Set(fmt.Sprintf("ordered_list.%d.order", i), i); err != nil {
		//
		//	}
		//}
	}

	return ReadObjectRenamingListsAndSets(ctx, d, meta)
}

func ReadObjectRenamingListsAndSets(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	list := collections.Map(objectRenamingDatabaseInstance.List, func(t objectRenamingDatabaseListItem) map[string]any {
		return map[string]any{
			"name":   t.Name,
			"string": t.String,
			"int":    t.Int,
		}
	})
	if err := d.Set("list", list); err != nil {
		return diag.FromErr(err)
	}

	orderedList := make([]map[string]any, len(objectRenamingDatabaseInstance.OrderedList))
	for index, item := range objectRenamingDatabaseInstance.OrderedList {
		orderedList[index] = map[string]any{
			"name":   item.Name,
			"string": item.String,
			"order":  strconv.Itoa(index),
		}
	}
	if err := d.Set("ordered_list", orderedList); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func DeleteObjectRenamingListsAndSets(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	objectRenamingDatabaseInstance.List = nil
	objectRenamingDatabaseInstance.OrderedList = nil
	d.SetId("")

	return nil
}
