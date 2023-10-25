package resources

import (
	"database/sql"
	"errors"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var roleSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ValidateFunc: func(val interface{}, key string) ([]string, []error) {
			additionalCharsToIgnoreValidation := []string{".", " ", ":", "(", ")"}
			return sdk.ValidateIdentifier(val, additionalCharsToIgnoreValidation)
		},
	},
	"comment": {
		Type:     schema.TypeString,
		Optional: true,
		// TODO validation
	},
	"tag": tagReferenceSchema,
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
	name := d.Get("name").(string)
	db := meta.(*sql.DB)
	builder := snowflake.NewRoleBuilder(db, name)
	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}
	if v, ok := d.GetOk("tag"); ok {
		tags := getTags(v)
		builder.WithTags(tags.toSnowflakeTagValues())
	}
	err := builder.Create()
	if err != nil {
		return err
	}
	d.SetId(name)
	return ReadRole(d, meta)
}

func ReadRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()
	// If the name is not set (such as during import) then use the id
	name := d.Get("name").(string)
	if name == "" {
		name = id
	}

	builder := snowflake.NewRoleBuilder(db, name)
	role, err := builder.Show()
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("[WARN] role (%s) not found", name)
		d.SetId("")
		return nil
	} else if err != nil {
		return err
	}
	if err := d.Set("name", role.Name.String); err != nil {
		return err
	}
	if err := d.Set("comment", role.Comment.String); err != nil {
		return err
	}
	return nil
}

func UpdateRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	builder := snowflake.NewRoleBuilder(db, name)

	if d.HasChange("name") {
		o, n := d.GetChange("name")
		builder.WithName(o.(string))
		err := builder.Rename(n.(string))
		if err != nil {
			return err
		}
		builder.WithName(n.(string))
	}

	if d.HasChange("comment") {
		o, n := d.GetChange("comment")
		if n == nil || n.(string) == "" {
			builder.WithComment(o.(string))
			err := builder.UnsetComment()
			if err != nil {
				return err
			}
		} else {
			err := builder.SetComment(n.(string))
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("tag") {
		o, n := d.GetChange("tag")
		removed, added, changed := getTags(o).diffs(getTags(n))
		for _, tA := range removed {
			err := builder.UnsetTag(tA.toSnowflakeTagValue())
			if err != nil {
				return err
			}
		}
		for _, tA := range added {
			err := builder.SetTag(tA.toSnowflakeTagValue())
			if err != nil {
				return err
			}
		}
		for _, tA := range changed {
			err := builder.ChangeTag(tA.toSnowflakeTagValue())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func DeleteRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	builder := snowflake.NewRoleBuilder(db, name)
	err := builder.Drop()
	return err
}
