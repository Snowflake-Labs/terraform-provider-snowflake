package resources

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	snowflakeValidation "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/validation"
)

var tagAssociationSchema = map[string]*schema.Schema{
	"object_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the object identifier for the tag association.",
		Deprecated:  "Use `object_identifier` instead",
		ForceNew:    true,
	},
	"object_identifier": {
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		Description: "Specifies the object identifier for the tag association.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Name of the object to associate the tag with.",
				},
				"database": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Name of the database that the object was created in.",
				},
				"schema": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Name of the schema that the object was created in.",
				},
			},
		},
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
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the tag. Note: format must follow: \"databaseName\".\"schemaName\".\"tagName\" or \"databaseName.schemaName.tagName\" or \"databaseName|schemaName.tagName\" (snowflake_tag.tag.id)",
		ForceNew:    true,
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
		Default:     true,
	},
}

// TagAssociation returns a pointer to the resource representing a schema.
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

// CreateTagAssociation implements schema.CreateFunc.
func CreateTagAssociation(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	tagID := d.Get("tag_id").(string)
	objectType := d.Get("object_type").(string)
	tagValue := d.Get("tag_value").(string)
	objectDatabase, objectSchema, objectName := expandObjectIdentifier(d.Get("object_identifier"))
	fullyQualifierObjectIdentifier := snowflakeValidation.FormatFullyQualifiedObjectID(objectDatabase, objectSchema, objectName)

	builder := snowflake.NewTagAssociationBuilder(tagID).WithObjectIdentifier(fullyQualifierObjectIdentifier).WithObjectType(objectType).WithTagValue(tagValue)

	q := builder.Create()
	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error associating tag to object: [%v] with command: [%v], tag_id [%v]", fullyQualifierObjectIdentifier, q, tagID)
	}

	skipValidate := d.Get("skip_validation").(bool)
	if !skipValidate {
		log.Println("[DEBUG] validating tag creation")

		if err := retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate)-time.Minute, func() *retry.RetryError {
			resp, err := snowflake.ListTagAssociations(builder, db)
			if err != nil {
				return retry.NonRetryableError(fmt.Errorf("error: %w", err))
			}

			// if length of response is zero, tag association was not found. retry for up to 70 minutes
			if len(resp) == 0 {
				return retry.RetryableError(fmt.Errorf("expected tag association to be created but not yet created"))
			}
			return nil
		}); err != nil {
			return fmt.Errorf("error validating tag association")
		}
	}

	d.SetId(helpers.EncodeSnowflakeID(builder.GetTagDatabase(), builder.GetTagSchema(), builder.GetTagName()))
	return ReadTagAssociation(d, meta)
}

// ReadTagAssociation implements schema.ReadFunc.
func ReadTagAssociation(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	tagID := d.Get("tag_id").(string)
	objectType := d.Get("object_type").(string)
	objectDatabase, objectSchema, objectName := expandObjectIdentifier(d.Get("object_identifier"))
	fullyQualifierObjectIdentifier := snowflakeValidation.FormatFullyQualifiedObjectID(objectDatabase, objectSchema, objectName)

	q := snowflake.NewTagAssociationBuilder(tagID).WithObjectIdentifier(fullyQualifierObjectIdentifier).WithObjectType(objectType).Show()
	row := snowflake.QueryRow(db, q)

	ta, err := snowflake.ScanTagAssociation(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] tag association (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		// return err
		return fmt.Errorf("error listing tag associations")
	}

	if err := d.Set("tag_value", ta.TagValue.String); err != nil {
		return err
	}
	return nil
}

func UpdateTagAssociation(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	tagID := d.Get("tag_id").(string)
	objectType := d.Get("object_type").(string)
	objectDatabase, objectSchema, objectName := expandObjectIdentifier(d.Get("object_identifier"))
	fullyQualifierObjectIdentifier := snowflakeValidation.FormatFullyQualifiedObjectID(objectDatabase, objectSchema, objectName)

	builder := snowflake.NewTagAssociationBuilder(tagID).WithObjectIdentifier(fullyQualifierObjectIdentifier).WithObjectType(objectType)

	if d.HasChange("skip_validation") {
		o, n := d.GetChange("skip_validation")
		log.Printf("[DEBUG] skip_validation changed from %v to %v", o, n)
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
			return fmt.Errorf("error updating tag association value for object [%v]", fullyQualifierObjectIdentifier)
		}
	}

	return ReadTagAssociation(d, meta)
}

// DeleteTagAssociation implements schema.DeleteFunc.
func DeleteTagAssociation(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	tagID := d.Get("tag_id").(string)
	objectType := d.Get("object_type").(string)
	objectDatabase, objectSchema, objectName := expandObjectIdentifier(d.Get("object_identifier"))
	fullyQualifierObjectIdentifier := snowflakeValidation.FormatFullyQualifiedObjectID(objectDatabase, objectSchema, objectName)
	q := snowflake.NewTagAssociationBuilder(tagID).WithObjectIdentifier(fullyQualifierObjectIdentifier).WithObjectType(objectType).Drop()

	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error deleting tag association for object id [%s]: %w", tagID, err)
	}

	d.SetId("")

	return nil
}
