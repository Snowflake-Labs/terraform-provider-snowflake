package resources

import (
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

var tagAttachmentSchema = map[string]*schema.Schema{
	"resourceId": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the resource identifier for the tag attachment.",
		ForceNew:    true,
	},
	"objectType": {
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
	"tagName": {
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Specifies the identifier for the tag. Note: format must follow: 'database.schema.tagId'",
		ValidateFunc: snowflakeValidation.ValidateFullyQualifiedTagPath,
		ForceNew:     true,
	},
	"tagValue": {
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Specifies the value of the tag",
		ValidateFunc: snowflakeValidation.ValidateFullyQualifiedTagPath,
		ForceNew:     true,
	},
}

// Schema returns a pointer to the resource representing a schema
func TagAttachment() *schema.Resource {
	return &schema.Resource{
		Create: CreateTagAttachment,
		Read:   ReadTagAttachment,
		Delete: DeleteTagAttachment,

		Schema: tagAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Minute),
		},
	}
}

// CreateSchema implements schema.CreateFunc
func CreateTagAttachment(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	tagName := d.Get("tag").(string)
	resourceId := d.Get("resourceId").(string)
	objectType := d.Get("objectType").(string)
	tagValue := d.Get("tagValue").(string)
	builder := snowflake.TagAttachment(tagName).WithResourceId(resourceId).WithObjectType(objectType).WithTagValue(tagValue)

	q := builder.Create()
	err := snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error attaching tag to resource: [%v]", resourceId)
	}
	// retry read until it works. add max timeout
	return resource.Retry(d.Timeout(schema.TimeoutCreate)-time.Minute, func() *resource.RetryError {
		resp, err := snowflake.ListTagAttachments(builder, db)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error listing tags: %s", err))
		}

		if resp == nil {
			return resource.RetryableError(fmt.Errorf("expected tag to be created but not yet created"))
		}

		err = ReadTagAttachment(d, meta)
		if err != nil {
			return resource.NonRetryableError(err)
		} else {
			return nil
		}
	})
}

// ReadSchema implements schema.ReadFunc
func ReadTagAttachment(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	tagName := d.Get("tag").(string)
	resourceId := d.Get("resourceId").(string)
	objectType := d.Get("objectType").(string)
	tagValue := d.Get("tagValue").(string)

	builder := snowflake.TagAttachment(tagName).WithResourceId(resourceId).WithObjectType(objectType).WithTagValue(tagValue)
	_, err := snowflake.ListTagAttachments(builder, db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] tag (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

// DeleteSchema implements schema.DeleteFunc
func DeleteTagAttachment(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	tagId := d.Get("tag").(string)
	resourceId := d.Get("resourceId").(string)
	objectType := d.Get("objectType").(string)

	q := snowflake.TagAttachment(tagId).WithResourceId(resourceId).WithObjectType(objectType).Drop()

	err := snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting tag attachment for resource [%v]", resourceId)
	}

	d.SetId("")

	return nil
}
