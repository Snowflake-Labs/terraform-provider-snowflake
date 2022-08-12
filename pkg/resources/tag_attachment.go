package resources

import (
	"database/sql"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/pkg/errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
)

var tagAttachmentSchema = map[string]*schema.Schema{
	"resource_id": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the resource identifier for the tag attachment.",
		ForceNew:    true,
	},
	"object_type": {
		Type:     schema.TypeString,
		Required: true,
		Description: "Specifies the type of object to add a tag to. ex: 'ACCOUNT', 'COLUMN', 'DATABASE', etc. " +
			"For more information: https://docs.snowflake.com/en/user-guide/object-tagging.html#supported-objects",
		ValidateFunc: validation.StringInSlice([]string{
			"ACCOUNT", "COLUMN", "DATABASE", "INTEGRATION", "PIPE", "ROLE", "SCHEMA", "STREAM", "SHARE", "STAGE",
			"TABLE", "TASK", "USER", "VIEW", "WAREHOUSE",
		}, true),
	},
	"tag": tagReferenceSchema,
}

// Schema returns a pointer to the resource representing a schema
func TagAttachment() *schema.Resource {
	return &schema.Resource{
		Create: CreateTagAttachment,
		Read:   ReadTagAttachment,
		Update: UpdateTagAttachment,
		Delete: DeleteTagAttachment,

		Schema: tagAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateSchema implements schema.CreateFunc
func CreateTagAttachment(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	resource_id := d.Get("resource_id").(string)
	object_type := d.Get("resource_type").(string)
	builder := snowflake.TagAttachment(d.Get("tag")).WithResourceId(resource_id).WithObjectType(object_type)

	q := builder.Create()
	err := snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error attaching tag to resource: [%v]", resource_id)
	}

	return ReadTagAttachment(d, meta)
}

// ReadSchema implements schema.ReadFunc
func ReadTagAttachment(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	// query the resource and get list of tags associated with the resource
	// check that the tag is associated with the resource
	// if tag is missing, remove from statefile
	// otherwise set the tag, set the resource ID, set the object type
	q := snowflake.TagAttachment(tag).WithResourceId(dbName).WithObjectType(schemaName).Show()
	row := snowflake.QueryRow(db, q)
	t, err := snowflake.ScanTag(row)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] tag (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	err := d.Set("tag")
	if err != nil {
		return err
	}

	err = d.Set("resource_id")
	if err != nil {
		return err
	}

	err = d.Set("object_type")
	if err != nil {
		return err
	}
	return nil
}

// UpdateSchema implements schema.UpdateFunc
func UpdateTagAttachment(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	resource_id := d.Get("resource_id").(string)
	object_type := d.Get("resource_type").(string)
	builder := snowflake.TagAttachment(d.Get("tag")).WithResourceId(resource_id).WithObjectType(object_type)

	if d.HasChange("resource_id") {
		q := builder.Drop()
		err := snowflake.Exec(db, q)
		// add back tag to the new resource
		if err != nil {
			return errors.Wrapf(err, "error updating resource id for tag attachment on %v", d.Id())
		}
	}

	if d.HasChange("object_type") {
		q := builder.Drop()
		err := snowflake.Exec(db, q)
		// add back tag to the new resource
		if err != nil {
			return errors.Wrapf(err, "error updating object type for tag attachment on %v", d.Id())
		}
	}

	if d.HasChange("tag") {
		q := builder.Drop()
		err := snowflake.Exec(db, q)
		// add new tag to the resource
		if err != nil {
			return errors.Wrapf(err, "error updating tag for tag attachment on %v", d.Id())
		}
	}

	return ReadTagAttachment(d, meta)
}

// DeleteSchema implements schema.DeleteFunc
// Remove the tag from the resource. Return error if missing permission or unable to remove
func DeleteTagAttachment(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	resource_id := d.Get("resource_id").(string)
	object_type := d.Get("resource_type").(string)

	q := snowflake.TagAttachment(d.Get("tag")).WithResourceId(resource_id).WithObjectType(object_type).Drop()

	err := snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting tag attachment for resource [%v]", resource_id)
	}

	d.SetId("")

	return nil
}
