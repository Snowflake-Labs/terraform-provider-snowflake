package resources

import (
	"context"
	"fmt"
	"log"
	"strings"

	providerresources "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var streamSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier for the stream; must be unique for the database and schema in which the stream is created.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema in which to create the stream.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database in which to create the stream.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the stream.",
	},
	"on_table": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "Specifies an identifier for the table the stream will monitor.",
		ExactlyOneOf:     []string{"on_table", "on_view", "on_stage"},
		DiffSuppressFunc: suppressIdentifierQuoting,
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
	"on_view": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "Specifies an identifier for the view the stream will monitor.",
		ExactlyOneOf:     []string{"on_table", "on_view", "on_stage"},
		DiffSuppressFunc: suppressIdentifierQuoting,
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
	"on_stage": {
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		Description:  "Specifies an identifier for the stage the stream will monitor.",
		ExactlyOneOf: []string{"on_table", "on_view", "on_stage"},
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// Suppress diff if the stage name is the same, even if database and schema are not specified
			return strings.Trim(strings.Split(old, ".")[len(strings.Split(old, "."))-1], "\"") == strings.Trim(strings.Split(new, ".")[len(strings.Split(new, "."))-1], "\"")
		},
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
	"append_only": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     false,
		Description: "Type of the stream that will be created.",
	},
	"insert_only": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     false,
		Description: "Create an insert only stream type.",
	},
	"show_initial_rows": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     false,
		Description: "Specifies whether to return all existing rows in the source table as row inserts the first time the stream is consumed.",
	},
	"owner": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the role that owns the stream.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func Stream() *schema.Resource {
	return &schema.Resource{
		Create: CreateStream,
		Read:   ReadStream,
		Update: UpdateStream,
		Delete: DeleteStream,
		DeprecationMessage: deprecatedResourceDescription(
			string(providerresources.StreamOnDirectoryTable),
			string(providerresources.StreamOnExternalTable),
			string(providerresources.StreamOnTable),
			string(providerresources.StreamOnView),
		),

		Schema: streamSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateStream implements schema.CreateFunc.
func CreateStream(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	appendOnly := d.Get("append_only").(bool)
	insertOnly := d.Get("insert_only").(bool)
	showInitialRows := d.Get("show_initial_rows").(bool)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	ctx := context.Background()

	onTable, onTableSet := d.GetOk("on_table")
	onView, onViewSet := d.GetOk("on_view")
	onStage, onStageSet := d.GetOk("on_stage")

	switch {
	case onTableSet:
		tableObjectIdentifier, err := helpers.DecodeSnowflakeParameterID(onTable.(string))
		if err != nil {
			return err
		}
		tableId := tableObjectIdentifier.(sdk.SchemaObjectIdentifier)

		table, err := client.Tables.ShowByID(ctx, tableId)
		if err != nil {
			return err
		}

		if table.IsExternal {
			req := sdk.NewCreateOnExternalTableStreamRequest(id, tableId)
			if insertOnly {
				req.WithInsertOnly(true)
			}
			if v, ok := d.GetOk("comment"); ok {
				req.WithComment(v.(string))
			}
			err := client.Streams.CreateOnExternalTable(ctx, req)
			if err != nil {
				return fmt.Errorf("error creating stream %v err = %w", name, err)
			}
		} else {
			req := sdk.NewCreateOnTableStreamRequest(id, tableId)
			if appendOnly {
				req.WithAppendOnly(true)
			}
			if showInitialRows {
				req.WithShowInitialRows(true)
			}
			if v, ok := d.GetOk("comment"); ok {
				req.WithComment(v.(string))
			}
			err := client.Streams.CreateOnTable(ctx, req)
			if err != nil {
				return fmt.Errorf("error creating stream %v err = %w", name, err)
			}
		}
	case onViewSet:
		viewObjectIdentifier, err := helpers.DecodeSnowflakeParameterID(onView.(string))
		viewId := viewObjectIdentifier.(sdk.SchemaObjectIdentifier)
		if err != nil {
			return err
		}

		_, err = client.Views.ShowByID(ctx, viewId)
		if err != nil {
			return err
		}

		req := sdk.NewCreateOnViewStreamRequest(id, viewId)
		if appendOnly {
			req.WithAppendOnly(true)
		}
		if showInitialRows {
			req.WithShowInitialRows(true)
		}
		if v, ok := d.GetOk("comment"); ok {
			req.WithComment(v.(string))
		}
		err = client.Streams.CreateOnView(ctx, req)
		if err != nil {
			return fmt.Errorf("error creating stream %v err = %w", name, err)
		}
	case onStageSet:
		stageObjectIdentifier, err := helpers.DecodeSnowflakeParameterID(onStage.(string))
		stageId := stageObjectIdentifier.(sdk.SchemaObjectIdentifier)
		if err != nil {
			return err
		}
		stageProperties, err := client.Stages.Describe(ctx, stageId)
		if err != nil {
			return err
		}
		if findStagePropertyValueByName(stageProperties, "ENABLE") != "true" {
			return fmt.Errorf("directory must be enabled on stage")
		}
		req := sdk.NewCreateOnDirectoryTableStreamRequest(id, stageId)
		if v, ok := d.GetOk("comment"); ok {
			req.WithComment(v.(string))
		}
		err = client.Streams.CreateOnDirectoryTable(ctx, req)
		if err != nil {
			return fmt.Errorf("error creating stream %v err = %w", name, err)
		}
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadStream(d, meta)
}

// ReadStream implements schema.ReadFunc.
func ReadStream(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	stream, err := client.Streams.ShowByID(ctx, id)
	if err != nil {
		log.Printf("[DEBUG] stream (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return err
	}
	if err := d.Set("name", stream.Name); err != nil {
		return err
	}
	if err := d.Set("database", stream.DatabaseName); err != nil {
		return err
	}
	if err := d.Set("schema", stream.SchemaName); err != nil {
		return err
	}
	switch *stream.SourceType {
	case sdk.StreamSourceTypeStage:
		if err := d.Set("on_stage", *stream.TableName); err != nil {
			return err
		}
	case sdk.StreamSourceTypeView:
		if err := d.Set("on_view", *stream.TableName); err != nil {
			return err
		}
	default:
		if err := d.Set("on_table", *stream.TableName); err != nil {
			return err
		}
	}
	if err := d.Set("append_only", *stream.Mode == "APPEND_ONLY"); err != nil {
		return err
	}
	if err := d.Set("insert_only", *stream.Mode == "INSERT_ONLY"); err != nil {
		return err
	}
	// TODO: SHOW STREAMS doesn't return that value right now (I'm not sure if it ever did), but probably we can assume
	// 	the customers got 'false' every time and hardcode it (it's only on create thing, so it's not necessary
	//	to track its value after creation).
	if err := d.Set("show_initial_rows", false); err != nil {
		return err
	}
	if err := d.Set("comment", *stream.Comment); err != nil {
		return err
	}
	if err := d.Set("owner", *stream.Owner); err != nil {
		return err
	}
	return nil
}

// UpdateStream implements schema.UpdateFunc.
func UpdateStream(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		if comment == "" {
			err := client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithUnsetComment(true))
			if err != nil {
				return fmt.Errorf("error unsetting stream comment on %v", d.Id())
			}
		} else {
			err := client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithSetComment(comment))
			if err != nil {
				return fmt.Errorf("error setting stream comment on %v", d.Id())
			}
		}
	}

	return ReadStream(d, meta)
}

// DeleteStream implements schema.DeleteFunc.
func DeleteStream(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	streamId := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(streamId))
	if err != nil {
		return fmt.Errorf("error deleting stream %v err = %w", d.Id(), err)
	}

	d.SetId("")

	return nil
}
