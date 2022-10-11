package resources

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		Description: "If true, skips validation of the tag association. It can take up to an hour for the SNOWFLAKE.TAG_REFERENCES table to update, and also requires ACCOUNT_ADMIN role to read from. https://docs.snowflake.com/en/sql-reference/account-usage/tag_references.html",
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

	skipValidate := d.Get("skip_validation").(bool)
	if !skipValidate {
		log.Println("[DEBUG] validating tag creation")

		err = resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate)-time.Minute, func() *resource.RetryError {

			resp, err := snowflake.ListTagAssociations(builder, db)

			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("error: %s", err))
			}

			// if length of response is zero, tag association was not found. retry for up to 70 minutes
			if len(resp) == 0 {
				return resource.RetryableError(fmt.Errorf("expected tag association to be created but not yet created"))
			}
			return nil
		})
		if err != nil {
			return errors.Wrap(err, "error validating tag association")
		}
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

	tagName := d.Get("tag_id").(string)
	objectName := d.Get("object_name").(string)
	objectType := d.Get("object_type").(string)
	tagValue := d.Get("tag_value").(string)
	skipValidate := d.Get("skip_validation").(bool)
	if skipValidate {
		log.Println("[DEBUG] skipping read for tag association that has skip_validation enabled")
		return nil
	}

	builder := snowflake.TagAssociation(tagName).WithObjectName(objectName).WithObjectType(objectType).WithTagValue(tagValue)
	_, err := snowflake.ListTagAssociations(builder, db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] tag associations (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		//return err
		return errors.Wrapf(err, "error listing tags, error is x")
	}

	return nil
}

func UpdateTagAssociation(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("skip_validation") {
		old, new := d.GetChange("skip_validation")
		log.Printf("[DEBUG] skip_validation changed from %v to %v", old, new)
	}

	return ReadTagAssociation(d, meta)
}

// DeleteSchema implements schema.DeleteFunc.
func DeleteTagAssociation(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	tagName := d.Get("tag_id").(string)
	objectName := d.Get("object_name").(string)
	objectType := d.Get("object_type").(string)

	q := snowflake.TagAssociation(tagName).WithObjectName(objectName).WithObjectType(objectType).Drop()

	err := snowflake.Exec(db, q)
	if err != nil {
		log.Printf("[DEBUG] error is %v", err.Error())
		return errors.Wrapf(err, "error deleting tag association for object [%v]", objectName)
	}

	d.SetId("")

	return nil
}
