package resources

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	RoleNameKey    = "name"
	RoleCommentKey = "comment"
	RoleTagKey     = "tag"
)

var roleSchema = map[string]*schema.Schema{
	RoleNameKey: {
		Type:     schema.TypeString,
		Required: true,
	},
	RoleCommentKey: {
		Type:     schema.TypeString,
		Optional: true,
	},
	RoleTagKey: tagReferenceSchema,
}

func Role() *schema.Resource {
	return &schema.Resource{
		Create: CreateRole,
		Read:   ReadRole,
		Delete: DeleteRole,
		Update: UpdateRole,

		Schema: roleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	ctx := context.Background()

	tagList := d.Get(RoleTagKey).([]interface{})
	tagAssociations := make([]sdk.TagAssociation, len(tagList))
	for i, tag := range tagList {
		t := tag.(map[string]interface{})
		tagAssociations[i] = sdk.TagAssociation{
			Name:  sdk.NewAccountObjectIdentifier(t["name"].(string)),
			Value: t["value"].(string),
		}
	}

	err := client.Roles.Create(ctx, objectIdentifier, &sdk.RoleCreateOptions{
		Comment: sdk.String(d.Get(RoleCommentKey).(string)),
		Tag:     tagAssociations,
	})
	if err != nil {
		return err
	}

	d.SetId(helpers.EncodeSnowflakeID(objectIdentifier))

	tags := make([]interface{}, len(tagAssociations))
	for i, t := range tagAssociations {
		tags[i] = map[string]interface{}{
			"name":  t.Name,
			"value": t.Value,
		}
	}

	return ReadRole(d, meta)
}

func ReadRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	ctx := context.Background()

	role, err := client.Roles.ShowByID(ctx, objectIdentifier)
	if err != nil {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] role (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	return errors.Join(
		d.Set(RoleNameKey, role.Name),
		d.Set(RoleCommentKey, role.Comment),
	)
}

func UpdateRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	ctx := context.Background()

	if d.HasChange(RoleNameKey) {
		oldName, newName := d.GetChange(RoleNameKey)
		err := client.Roles.Alter(ctx, objectIdentifier, &sdk.RoleAlterOptions{
			RenameTo: sdk.NewAccountObjectIdentifier(newName.(string)),
		})
		if err != nil {
			return fmt.Errorf("Failed to update role's name from %v to %v", oldName, newName)
		}
	}

	if d.HasChange(RoleCommentKey) {
		oldComment, newComment := d.GetChange(RoleCommentKey)
		var err error
		if newComment == nil || len(newComment.(string)) == 0 {
			err = client.Roles.Alter(ctx, objectIdentifier, &sdk.RoleAlterOptions{
				Unset: &sdk.RoleUnset{
					Comment: sdk.Bool(true),
				},
			})
		} else {
			err = client.Roles.Alter(ctx, objectIdentifier, &sdk.RoleAlterOptions{
				Set: &sdk.RoleSet{
					Comment: sdk.String(newComment.(string)),
				},
			})
		}
		if err != nil {
			return fmt.Errorf("Failed to update role's comment from %v to %v", oldComment, newComment)
		}
	}

	if d.HasChange(RoleTagKey) {
		oldTags, newTags := d.GetChange(RoleTagKey)
		removed, added, changed := getTags(oldTags).diffs(getTags(newTags))

		toSet := make([]tag, len(added)+len(changed))
		copy(toSet, added)
		copy(toSet[len(added):], changed)

		for _, tag := range removed {
			err := client.Roles.Alter(ctx, objectIdentifier, &sdk.RoleAlterOptions{
				Unset: &sdk.RoleUnset{
					Tag: []sdk.ObjectIdentifier{sdk.NewAccountObjectIdentifier(tag.name)},
				},
			})
			if err != nil {
				return err
			}
		}

		for _, tag := range toSet {
			err := client.Roles.Alter(ctx, objectIdentifier, &sdk.RoleAlterOptions{
				Set: &sdk.RoleSet{
					Tag: []sdk.TagAssociation{
						{
							Name:  sdk.NewAccountObjectIdentifier(tag.name),
							Value: tag.value,
						},
					},
				},
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func DeleteRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	ctx := context.Background()

	err := client.Roles.Drop(ctx, objectIdentifier, nil)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
