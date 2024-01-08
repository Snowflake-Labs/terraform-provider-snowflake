package resources

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var externalTableSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier for the external table; must be unique for the database and schema in which the externalTable is created.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema in which to create the external table.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database in which to create the external table.",
	},
	// TODO: Could be a string that we would validate as always "delta" (could be easy to add another type if snowflake introduces one)
	"table_format_delta": {
		Type:         schema.TypeBool,
		Required:     true,
		ForceNew:     true,
		Description:  `Identifies the external table as referencing a Delta Lake on the cloud storage location. A Delta Lake on Amazon S3, Google Cloud Storage, or Microsoft Azure cloud storage is supported.`,
		RequiredWith: []string{"user_specified_partitions"},
	},
	"column": {
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		ForceNew:    true,
		Description: "Definitions of a column to create in the external table. Minimum one required.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Column name",
					ForceNew:    true,
				},
				"type": {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "Column type, e.g. VARIANT",
					ForceNew:     true,
					ValidateFunc: IsDataType(),
				},
				"as": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "String that specifies the expression for the column. When queried, the column returns results derived from this expression.",
					ForceNew:    true,
				},
			},
		},
	},
	"location": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies a location for the external table.",
	},
	"file_format": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the file format for the external table.",
	},
	"pattern": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the file names and/or paths on the external stage to match.",
	},
	"aws_sns_topic": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the aws sns topic for the external table.",
	},
	"user_specified_partitions": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Description: "Enables to manage partitions manually and perform updates instead of recreating table on partition_by change.",
	},
	"partition_by": {
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
		//ForceNew:    true,
		// TODO: Update on user_specified_partitions = true and force new on false
		Description: "Specifies any partition columns to evaluate for the external table.",
	},
	"refresh_on_create": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies weather to refresh when an external table is created.",
		Default:     true,
		ForceNew:    true,
	},
	"auto_refresh": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether to automatically refresh the external table metadata once, immediately after the external table is created.",
		Default:     true,
		ForceNew:    true,
	},
	"copy_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies to retain the access permissions from the original table when an external table is recreated using the CREATE OR REPLACE TABLE variant",
		Default:     false,
		ForceNew:    true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies a comment for the external table.",
	},
	"owner": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the role that owns the external table.",
	},
	"tag": tagReferenceSchema,
}

func ExternalTable() *schema.Resource {
	return &schema.Resource{
		Create: CreateExternalTable,
		Read:   ReadExternalTable,
		Update: UpdateExternalTable,
		Delete: DeleteExternalTable,

		Schema: externalTableSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateExternalTable implements schema.CreateFunc.
func CreateExternalTable(d *schema.ResourceData, meta any) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)

	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schema, name)
	location := d.Get("location").(string)
	fileFormat := d.Get("file_format").(string)

	tableColumns := d.Get("column").([]any)
	columnRequests := make([]*sdk.ExternalTableColumnRequest, len(tableColumns))
	for i, col := range tableColumns {
		columnDef := map[string]string{}
		for key, val := range col.(map[string]any) {
			columnDef[key] = val.(string)
		}

		name := columnDef["name"]
		dataTypeString := columnDef["type"]
		dataType, err := sdk.ToDataType(dataTypeString)
		if err != nil {
			return fmt.Errorf(`failed to parse datatype: %s`, dataTypeString)
		}
		as := columnDef["as"]
		columnRequests[i] = sdk.NewExternalTableColumnRequest(name, dataType, as)
	}
	autoRefresh := sdk.Bool(d.Get("auto_refresh").(bool))
	refreshOnCreate := sdk.Bool(d.Get("refresh_on_create").(bool))
	copyGrants := sdk.Bool(d.Get("copy_grants").(bool))

	var partitionBy []string
	if v, ok := d.GetOk("partition_by"); ok {
		partitionBy = expandStringList(v.([]any))
	}

	var pattern *string
	if v, ok := d.GetOk("pattern"); ok {
		pattern = sdk.String(v.(string))
	}

	var awsSnsTopic *string
	if v, ok := d.GetOk("aws_sns_topic"); ok {
		awsSnsTopic = sdk.String(v.(string))
	}

	var comment *string
	if v, ok := d.GetOk("comment"); ok {
		comment = sdk.String(v.(string))
	}

	var tagAssociationRequests []*sdk.TagAssociationRequest
	if _, ok := d.GetOk("tag"); ok {
		tagAssociations := getPropertyTags(d, "tag")
		tagAssociationRequests = make([]*sdk.TagAssociationRequest, len(tagAssociations))
		for i, t := range tagAssociations {
			tagAssociationRequests[i] = sdk.NewTagAssociationRequest(t.Name, t.Value)
		}
	}

	switch {
	case d.Get("table_format_delta").(bool):
		err := client.ExternalTables.CreateDeltaLake(
			ctx,
			sdk.NewCreateDeltaLakeExternalTableRequest(id, location).
				WithRawFileFormat(&fileFormat).
				WithColumns(columnRequests).
				WithPartitionBy(partitionBy).
				WithRefreshOnCreate(refreshOnCreate).
				WithAutoRefresh(autoRefresh).
				WithCopyGrants(copyGrants).
				WithComment(comment).
				WithTag(tagAssociationRequests),
		)
		if err != nil {
			return err
		}
	case d.Get("user_specified_partitions").(bool):
		err := client.ExternalTables.CreateWithManualPartitioning(
			ctx,
			sdk.NewCreateWithManualPartitioningExternalTableRequest(id, location).
				WithRawFileFormat(&fileFormat).
				WithColumns(columnRequests).
				WithPartitionBy(partitionBy).
				WithRawFileFormat(sdk.String(fileFormat)).
				WithCopyGrants(copyGrants).
				WithComment(comment).
				WithTag(tagAssociationRequests),
		)
		if err != nil {
			return err
		}
	default:
		err := client.ExternalTables.Create(
			ctx,
			sdk.NewCreateExternalTableRequest(id, location).
				WithRawFileFormat(&fileFormat).
				WithColumns(columnRequests).
				WithPartitionBy(partitionBy).
				WithRefreshOnCreate(refreshOnCreate).
				WithAutoRefresh(autoRefresh).
				WithPattern(pattern).
				WithRawFileFormat(sdk.String(fileFormat)).
				WithAwsSnsTopic(awsSnsTopic).
				WithCopyGrants(copyGrants).
				WithComment(comment).
				WithTag(tagAssociationRequests),
		)
		if err != nil {
			return err
		}
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadExternalTable(d, meta)
}

// ReadExternalTable implements schema.ReadFunc.
func ReadExternalTable(d *schema.ResourceData, meta any) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	externalTable, err := client.ExternalTables.ShowByID(ctx, sdk.NewShowExternalTableByIDRequest(id))
	if err != nil {
		log.Printf("[DEBUG] external table (%s) not found", d.Id())
		d.SetId("")
		return err
	}

	if err := d.Set("name", externalTable.Name); err != nil {
		return err
	}

	if err := d.Set("owner", externalTable.Owner); err != nil {
		return err
	}

	return nil
}

// UpdateExternalTable implements schema.UpdateFunc.
func UpdateExternalTable(d *schema.ResourceData, meta any) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	if d.HasChange("tag") {
		unsetTags, setTags := GetTagsDiff(d, "tag")

		err := client.ExternalTables.Alter(ctx, sdk.NewAlterExternalTableRequest(id).WithUnsetTag(unsetTags))
		if err != nil {
			return fmt.Errorf("error setting tags on %v, err = %w", d.Id(), err)
		}

		tagAssociationRequests := make([]*sdk.TagAssociationRequest, len(setTags))
		for i, t := range setTags {
			tagAssociationRequests[i] = sdk.NewTagAssociationRequest(t.Name, t.Value)
		}
		err = client.ExternalTables.Alter(ctx, sdk.NewAlterExternalTableRequest(id).WithSetTag(tagAssociationRequests))
		if err != nil {
			return fmt.Errorf("error setting tags on %v, err = %w", d.Id(), err)
		}
	}

	return ReadExternalTable(d, meta)
}

// DeleteExternalTable implements schema.DeleteFunc.
func DeleteExternalTable(d *schema.ResourceData, meta any) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	err := client.ExternalTables.Drop(ctx, sdk.NewDropExternalTableRequest(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
