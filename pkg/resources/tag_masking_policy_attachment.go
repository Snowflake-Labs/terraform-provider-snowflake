package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	snowflakeValidation "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/validation"
)

const (
	tagAttachmentPolicyIDDelimiter = '|'
)

var mpAttachmentPolicySchema = map[string]*schema.Schema{
	"tag_id": {
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Specifies the identifier for the tag. Note: format must follow: \"databaseName\".\"schemaName\".\"tagName\" or \"databaseName.schemaName.tagName\" or \"databaseName|schemaName.tagName\" (snowflake_tag.tag.id)",
		ValidateFunc: snowflakeValidation.ValidateFullyQualifiedTagID,
		ForceNew:     true,
	},
	"masking_policy_database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database where the masking policy is located.",
	},
	"masking_policy_schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema where the masking policy is located",
	},
	"masking_policy_name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The name of the masking policy to attach.",
	},
}

type attachmentID struct {
	TagDatabaseName           string
	TagSchemaName             string
	TagName                   string
	MaskingPolicyDatabaseName string
	MaskingPolicySchemaName   string
	MaskingPolicyName         string
}

// String() takes in a schemaID object and returns a pipe-delimited string:
// TagDatabaseName | TagSchemaName | TagName | MaskingPolicyDatabaseName | MaskingPolicySchemaName | MaskingPolicyName
func (ti *attachmentID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = tagAttachmentPolicyIDDelimiter
	dataIdentifiers := [][]string{{ti.TagDatabaseName, ti.TagSchemaName, ti.TagName, ti.MaskingPolicyDatabaseName, ti.MaskingPolicySchemaName, ti.MaskingPolicyName}}
	err := csvWriter.WriteAll(dataIdentifiers)
	if err != nil {
		return "", err
	}
	strTagID := strings.TrimSpace(buf.String())
	return strTagID, nil
}

// attachedPolicyIDFromString() takes in a pipe-delimited string: TagDatabaseName | TagSchemaName | TagName | MaskingPolicyDatabaseName | MaskingPolicySchemaName | MaskingPolicyName
// and returns a attachmentID object
func attachedPolicyIDFromString(stringID string) (*attachmentID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = tagAttachmentPolicyIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per schema")
	}
	if len(lines[0]) != 6 {
		return nil, fmt.Errorf("6 fields allowed")
	}

	attachmentResult := &attachmentID{
		TagDatabaseName:           lines[0][0],
		TagSchemaName:             lines[0][1],
		TagName:                   lines[0][2],
		MaskingPolicyDatabaseName: lines[0][3],
		MaskingPolicySchemaName:   lines[0][4],
		MaskingPolicyName:         lines[0][5],
	}
	return attachmentResult, nil
}

// Schema returns a pointer to the resource representing a schema
func TagMaskingPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: CreateTagMaskingPolicyAttachemt,
		Read:   ReadTagMaskingPolicyAttachemt,
		Delete: DeleteTagMaskingPolicyAttachemt,

		Schema: mpAttachmentPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateTagMaskingPolicyAttachemt implements schema.CreateFunc
func CreateTagMaskingPolicyAttachemt(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	tagId := d.Get("tag_id").(string)
	tagIdStruct, idErr := tagIDFromString(tagId)
	if idErr != nil {
		return idErr
	}
	tagDb := tagIdStruct.DatabaseName
	tagSchema := tagIdStruct.SchemaName
	tagName := tagIdStruct.TagName
	mpDb := d.Get("masking_policy_database").(string)
	mpSchema := d.Get("masking_policy_schema").(string)
	mpName := d.Get("masking_policy_name").(string)

	mP := snowflake.MaskingPolicy(mpName, mpDb, mpSchema)
	builder := snowflake.Tag(tagName).WithDB(tagDb).WithSchema(tagSchema).WithMaskingPolicy(mP)

	q := builder.AddMaskingPolicy()

	err := snowflake.Exec(db, q)

	if err != nil {
		return errors.Wrapf(err, "error attaching masking policy %v to tag %v", mpName, tagName)
	}

	mpID := &attachmentID{
		TagDatabaseName:           tagDb,
		TagSchemaName:             tagSchema,
		TagName:                   tagName,
		MaskingPolicyDatabaseName: mpDb,
		MaskingPolicySchemaName:   mpSchema,
		MaskingPolicyName:         mpName,
	}
	dataIDInput, err := mpID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadTagMaskingPolicyAttachemt(d, meta)
}

// ReadTagTagMaskingPolicyAttachemt implements schema.ReadFunc
func ReadTagMaskingPolicyAttachemt(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	attachementID, err := attachedPolicyIDFromString(d.Id())
	if err != nil {
		return err
	}

	tagDbName := attachementID.TagDatabaseName
	tagSchemaName := attachementID.TagSchemaName
	tagName := attachementID.TagName
	mpDbName := attachementID.MaskingPolicyDatabaseName
	mpSchameName := attachementID.MaskingPolicySchemaName
	mpName := attachementID.MaskingPolicyName

	mP := snowflake.MaskingPolicy(mpName, mpDbName, mpSchameName)
	builder := snowflake.Tag(tagName).WithDB(tagDbName).WithSchema(tagSchemaName).WithMaskingPolicy(mP)

	row := snowflake.QueryRow(db, builder.ShowAttachedPolicy())
	t, err := snowflake.ScanTagPolicy(row)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] attached policy (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	tagId := TagID{
		DatabaseName: t.RefDb.String,
		SchemaName:   t.RefSchema.String,
		TagName:      t.RefEntity.String,
	}

	tagIdString, err := tagId.String()
	if err != nil {
		return err
	}
	d.Set("tag_id", tagIdString)
	d.Set("masking_policy_database", t.PolicyDb.String)
	d.Set("masking_policy_schema", t.PolicySchema.String)
	d.Set("masking_policy_name", t.PolicyName.String)

	return nil
}

// DeleteTagMaskingPolicyAttachemt implements schema.DeleteFunc
func DeleteTagMaskingPolicyAttachemt(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	attachementID, err := attachedPolicyIDFromString(d.Id())
	if err != nil {
		return err
	}

	tagDbName := attachementID.TagDatabaseName
	tagSchemaName := attachementID.TagSchemaName
	tagName := attachementID.TagName
	mpDbName := attachementID.MaskingPolicyDatabaseName
	mpSchameName := attachementID.MaskingPolicySchemaName
	mpName := attachementID.MaskingPolicyName

	mP := snowflake.MaskingPolicy(mpName, mpDbName, mpSchameName)

	builder := snowflake.Tag(tagName).WithDB(tagDbName).WithSchema(tagSchemaName).WithMaskingPolicy(mP)

	err = snowflake.Exec(db, builder.RemoveMaskingPolicy())
	if err != nil {
		return errors.Wrapf(err, "error unattaching masking policy for %v", d.Id())
	}

	d.SetId("")

	return nil
}
