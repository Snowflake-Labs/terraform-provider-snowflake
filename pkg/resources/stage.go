package resources

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var stageSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the stage; must be unique for the database and schema in which the stage is created.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the stage.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the stage.",
		ForceNew:    true,
	},
	"url": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the URL for the stage.",
	},
	"credentials": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the credentials for the stage.",
		Sensitive:   true,
	},
	"storage_integration": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the name of the storage integration used to delegate authentication responsibility for external cloud storage to a Snowflake identity and access management (IAM) entity.",
	},
	"file_format": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the file format for the stage.",
	},
	"copy_options": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the copy options for the stage.",
	},
	"encryption": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the encryption settings for the stage.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the stage.",
	},
	"directory": {
		Type:        schema.TypeString,
		ForceNew:    true,
		Optional:    true,
		Description: "Specifies the directory settings for the stage.",
	},
	"aws_external_id": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"snowflake_iam_user": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"tag": tagReferenceSchema,
}

// TODO (SNOW-1019005): Remove snowflake package that is used in Create and Update operations
func Stage() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateStage,
		ReadContext:   ReadStage,
		UpdateContext: UpdateStage,
		DeleteContext: DeleteStage,

		Schema: stageSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateStage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)

	builder := snowflake.NewStageBuilder(name, database, schema)

	if v, ok := d.GetOk("url"); ok {
		builder.WithURL(v.(string))
	}

	if v, ok := d.GetOk("credentials"); ok {
		builder.WithCredentials(v.(string))
	}

	if v, ok := d.GetOk("storage_integration"); ok {
		builder.WithStorageIntegration(v.(string))
	}

	if v, ok := d.GetOk("file_format"); ok {
		builder.WithFileFormat(v.(string))
	}

	if v, ok := d.GetOk("copy_options"); ok {
		builder.WithCopyOptions(v.(string))
	}

	if v, ok := d.GetOk("directory"); ok {
		builder.WithDirectory(v.(string))
	}

	if v, ok := d.GetOk("encryption"); ok {
		builder.WithEncryption(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := d.GetOk("tag"); ok {
		tags := getTags(v)
		builder.WithTags(tags.toSnowflakeTagValues())
	}

	q := builder.Create()

	if err := snowflake.Exec(db, q); err != nil {
		return diag.Errorf("error creating stage %v", name)
	}

	d.SetId(helpers.EncodeSnowflakeID(database, schema, name))

	return ReadStage(ctx, d, meta)
}

func ReadStage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	db := meta.(*sql.DB)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	client := sdk.NewClientFromDB(db)

	properties, err := client.Stages.Describe(ctx, id)
	if err != nil {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to describe stage",
				Detail:   fmt.Sprintf("Id: %s, Err: %s", d.Id(), err),
			},
		}
	}

	stage, err := client.Stages.ShowByID(ctx, id)
	if err != nil {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to show stage by id",
				Detail:   fmt.Sprintf("Id: %s, Err: %s", d.Id(), err),
			},
		}
	}

	if err := d.Set("name", stage.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database", stage.DatabaseName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("schema", stage.SchemaName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("url", strings.Trim(findStagePropertyValueByName(properties, "URL"), "[\"]")); err != nil {
		return diag.FromErr(err)
	}

	fileFormat := make([]string, 0)
	for _, property := range properties {
		if property.Parent == "STAGE_FILE_FORMAT" && property.Value != property.Default {
			fileFormat = append(fileFormat, fmt.Sprintf("%s = %s", property.Name, property.Value))
		}
	}
	if err := d.Set("file_format", strings.Join(fileFormat, " ")); err != nil {
		return diag.FromErr(err)
	}

	copyOptions := make([]string, 0)
	for _, property := range properties {
		if property.Parent == "STAGE_COPY_OPTIONS" && property.Value != property.Default {
			copyOptions = append(copyOptions, fmt.Sprintf("%s = %s", property.Name, property.Value))
		}
	}
	if err := d.Set("copy_options", strings.Join(copyOptions, " ")); err != nil {
		return diag.FromErr(err)
	}

	directory := make([]string, 0)
	for _, property := range properties {
		if property.Parent == "DIRECTORY" && property.Value != property.Default && property.Name != "LAST_REFRESHED_ON" {
			directory = append(directory, fmt.Sprintf("%s = %s", property.Name, property.Value))
		}
	}
	if err := d.Set("directory", strings.Join(directory, " ")); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("storage_integration", stage.StorageIntegration); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", stage.Comment); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("aws_external_id", findStagePropertyValueByName(properties, "AWS_EXTERNAL_ID")); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("snowflake_iam_user", findStagePropertyValueByName(properties, "SNOWFLAKE_IAM_USER")); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateStage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	builder := snowflake.NewStageBuilder(id.Name(), id.DatabaseName(), id.SchemaName())

	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	if d.HasChange("credentials") {
		credentials := d.Get("credentials")
		q := builder.ChangeCredentials(credentials.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return diag.Errorf("error updating stage credentials on %v", d.Id())
		}
	}

	if d.HasChange("storage_integration") && d.HasChange("url") {
		si := d.Get("storage_integration")
		url := d.Get("url")
		q := builder.ChangeStorageIntegrationAndUrl(si.(string), url.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return diag.Errorf("error updating stage storage integration and url on %v", d.Id())
		}
	} else {
		if d.HasChange("storage_integration") {
			si := d.Get("storage_integration")
			q := builder.ChangeStorageIntegration(si.(string))
			if err := snowflake.Exec(db, q); err != nil {
				return diag.Errorf("error updating stage storage integration on %v", d.Id())
			}
		}

		if d.HasChange("url") {
			url := d.Get("url")
			q := builder.ChangeURL(url.(string))
			if err := snowflake.Exec(db, q); err != nil {
				return diag.Errorf("error updating stage url on %v", d.Id())
			}
		}
	}

	if d.HasChange("encryption") {
		encryption := d.Get("encryption")
		q := builder.ChangeEncryption(encryption.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return diag.Errorf("error updating stage encryption on %v", d.Id())
		}
	}

	if d.HasChange("file_format") {
		fileFormat := d.Get("file_format")
		q := builder.ChangeFileFormat(fileFormat.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return diag.Errorf("error updating stage file format on %v", d.Id())
		}
	}

	if d.HasChange("copy_options") {
		copyOptions := d.Get("copy_options")
		q := builder.ChangeCopyOptions(copyOptions.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return diag.Errorf("error updating stage copy options on %v", d.Id())
		}
	}

	if d.HasChange("comment") {
		comment := d.Get("comment")
		q := builder.ChangeComment(comment.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return diag.Errorf("error updating stage comment on %v", d.Id())
		}
	}

	if d.HasChange("tag") {
		unsetTags, setTags := GetTagsDiff(d, "tag")

		if len(unsetTags) > 0 {
			err := client.Stages.Alter(ctx, sdk.NewAlterStageRequest(id).WithUnsetTags(unsetTags))
			if err != nil {
				return diag.Errorf("error occurred when dropping tags on stage with id: %v, err = %s", d.Id(), err)
			}
		}

		if len(setTags) > 0 {
			err := client.Stages.Alter(ctx, sdk.NewAlterStageRequest(id).WithSetTags(setTags))
			if err != nil {
				return diag.Errorf("error occurred when setting tags on stage with id: %v, err = %s", d.Id(), err)
			}
		}
	}

	return ReadStage(ctx, d, meta)
}

func DeleteStage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	db := meta.(*sql.DB)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	client := sdk.NewClientFromDB(db)

	err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to drop stage",
				Detail:   fmt.Sprintf("Id: %s, Err: %s", d.Id(), err),
			},
		}
	}

	d.SetId("")

	return nil
}

func findStagePropertyValueByName(properties []sdk.StageProperty, name string) string {
	for _, property := range properties {
		if property.Name == name {
			return property.Value
		}
	}
	return ""
}
