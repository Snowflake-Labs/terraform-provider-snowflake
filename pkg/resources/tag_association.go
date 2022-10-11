package resources

import (
	"database/sql"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/pkg/errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	snowflakeValidation "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/validation"
)

var tagAssociationSchema = map[string]*schema.Schema{
	"object_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the object identifier for the tag association.",
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
		ForceNew: true,
	},
	"tag_id": {
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Specifies the identifier for the tag. Note: format must follow: \"databaseName\".\"schemaName\".\"tagName\" or \"databaseName.schemaName.tagName\" or \"databaseName|schemaName.tagName\" (snowflake_tag.tag.id)",
		ValidateFunc: snowflakeValidation.ValidateFullyQualifiedObjectID,
		ForceNew:     true,
	},
	"tag_value": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the value of the tag, (e.g. 'finance' or 'engineering')",
		ForceNew:    true,
	},
	"skip_validation": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "If true, skips validation of the tag association.",
		Deprecated:  "Tag associations are now always validated without latency using the SYSTEM$GET_TAG function.",
		Default:     false,
	},
}

// Schema returns a pointer to the resource representing a schema.
func TagAssociation() *schema.Resource {
	return &schema.Resource{
		Create: CreateTagAssociation,
		Read:   ReadTagAssociation,
		Update: UpdateTagAssociation,
		Delete: DeleteTagAssociation,

		Schema: tagAssociationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Minute),
		},
	}
}

// CreateSchema implements schema.CreateFunc.
func CreateTagAssociation(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	tagID := d.Get("tag_id").(string)
	objectName := d.Get("object_name").(string)
	objectType := d.Get("object_type").(string)
	tagValue := d.Get("tag_value").(string)
	builder := snowflake.TagAssociation(tagID).WithObjectName(objectName).WithObjectType(objectType).WithTagValue(tagValue)

	q := builder.Create()
	err := snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error associating tag to object: [%v] with command: [%v], tag_id [%v]", objectName, q, tagID)
	}

	_, err = snowflake.ListTagAssociations(builder, db)
	if err != nil {
		return errors.Wrap(err, "error validating tag association")
	}

	t := &TagID{
		DatabaseName: builder.GetTagDatabase(),
		SchemaName:   builder.GetTagSchema(),
		TagName:      builder.GetTagName(),
	}
	dataIDInput, err := t.String()
	if err != nil {
		return errors.Wrap(err, "error creating tag id")
	}
	d.SetId(dataIDInput)

	return ReadTagAssociation(d, meta)
}

// ReadSchema implements schema.ReadFunc.
func ReadTagAssociation(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	tagID := d.Get("tag_id").(string)
	objectName := d.Get("object_name").(string)
	objectType := d.Get("object_type").(string)

	q := snowflake.TagAssociation(tagID).WithObjectName(objectName).WithObjectType(objectType).Show()
	row := snowflake.QueryRow(db, q)

	ta, err := snowflake.ScanTagAssociation(row)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] tag association (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		//return err
		return errors.Wrap(err, "error listing tag associations")
	}

	err = d.Set("tag_value", ta.TagValue.String)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTagAssociation(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	tagID := d.Get("tag_id").(string)
	objectName := d.Get("object_name").(string)
	objectType := d.Get("object_type").(string)

	builder := snowflake.TagAssociation(tagID).WithObjectName(objectName).WithObjectType(objectType)

	if d.HasChange("skip_validation") {
		old, new := d.GetChange("skip_validation")
		log.Printf("[DEBUG] skip_validation changed from %v to %v", old, new)
	}

	if d.HasChange("tag_value") {
		tagValue, ok := d.GetOk("tag_value")
		var q string
		if ok {
			q = builder.WithTagValue(tagValue.(string)).Create()
		} else {
			q = builder.WithTagValue("").Create()
		}
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating tag association value for object [%v]", objectName)
		}
	}

	return ReadTagAssociation(d, meta)
}

// DeleteSchema implements schema.DeleteFunc.
func DeleteTagAssociation(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	tagID := d.Get("tag_id").(string)
	objectName := d.Get("object_name").(string)
	objectType := d.Get("object_type").(string)

	q := snowflake.TagAssociation(tagID).WithObjectName(objectName).WithObjectType(objectType).Drop()

	err := snowflake.Exec(db, q)
	if err != nil {
		log.Printf("[DEBUG] error is %v", err.Error())
		return errors.Wrapf(err, "error deleting tag association for object [%v]", objectName)
	}

	d.SetId("")

	return nil
}
