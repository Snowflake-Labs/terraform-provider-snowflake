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
	//"set": {},
}

// TODO: This may fail for lists where items with the same values are allowed
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

				// try to find the same hash in the state (because it represents old value)
				newItemWasAlreadyPresent := false
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

					// try to find the same hash in the config (because it represents new value)
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

/*
Requirements for List:
- Differentiate between add/modify/remove

*/

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
	}
}

type objectRenamingDatabaseListItem struct {
	String string
	Int    int
}

type objectRenamingDatabase struct {
	List []objectRenamingDatabaseListItem
}

var objectRenamingDatabaseInstance = new(objectRenamingDatabase)

func CreateObjectRenamingListsAndSets(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	objectRenamingDatabaseInstance.List = collections.Map(d.Get("list").([]any), func(item any) objectRenamingDatabaseListItem {
		return objectRenamingDatabaseListItem{
			String: item.(map[string]any)["string"].(string),
			Int:    item.(map[string]any)["int"].(int),
		}
	})

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
			return objectRenamingDatabaseListItem{
				String: item.AsValueMap()["string"].AsString(),
				Int:    int(intValue),
			}
		})
		newListMapped := collections.Map(newList, func(item cty.Value) objectRenamingDatabaseListItem {
			intValue, _ := item.AsValueMap()["int"].AsBigFloat().Int64()
			return objectRenamingDatabaseListItem{
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

	return ReadObjectRenamingListsAndSets(ctx, d, meta)
}

func ReadObjectRenamingListsAndSets(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	list := collections.Map(objectRenamingDatabaseInstance.List, func(t objectRenamingDatabaseListItem) map[string]any {
		return map[string]any{
			"string": t.String,
			"int":    t.Int,
		}
	})
	if err := d.Set("list", list); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func DeleteObjectRenamingListsAndSets(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	objectRenamingDatabaseInstance.List = nil
	d.SetId("")

	return nil
}
