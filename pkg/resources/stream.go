package resources

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
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
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		Description:  "Specifies an identifier for the table the stream will monitor.",
		ExactlyOneOf: []string{"on_table", "on_view", "on_stage"},
	},
	"on_view": {
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		Description:  "Specifies an identifier for the view the stream will monitor.",
		ExactlyOneOf: []string{"on_table", "on_view", "on_stage"},
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
}

func Stream() *schema.Resource {
	return &schema.Resource{
		Create: CreateStream,
		Read:   ReadStream,
		Update: UpdateStream,
		Delete: DeleteStream,

		Schema: streamSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateStream implements schema.CreateFunc.
func CreateStream(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	appendOnly := d.Get("append_only").(bool)
	insertOnly := d.Get("insert_only").(bool)
	showInitialRows := d.Get("show_initial_rows").(bool)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	onTable, onTableSet := d.GetOk("on_table")
	onView, onViewSet := d.GetOk("on_view")
	onStage, onStageSet := d.GetOk("on_stage")

	switch {
	case onTableSet:
		tableId := helpers.DecodeSnowflakeID(onTable.(string)).(sdk.SchemaObjectIdentifier)

		tq := snowflake.NewTableBuilder(tableId.Name(), tableId.DatabaseName(), tableId.SchemaName()).Show()
		tableRow := snowflake.QueryRow(db, tq)
		t, err := snowflake.ScanTable(tableRow)
		if err != nil {
			return err
		}

		if t.IsExternal.String == "Y" {
			req := sdk.NewCreateStreamOnExternalTableRequest(id, tableId)
			if insertOnly {
				req.WithInsertOnly(sdk.Bool(true))
			}
			if v, ok := d.GetOk("comment"); ok {
				req.WithComment(sdk.String(v.(string)))
			}
			err := client.Streams.CreateOnExternalTable(ctx, req)
			if err != nil {
				return fmt.Errorf("error creating stream %v err = %w", name, err)
			}
		} else {
			req := sdk.NewCreateStreamOnTableRequest(id, tableId)
			if appendOnly {
				req.WithAppendOnly(sdk.Bool(true))
			}
			if showInitialRows {
				req.WithShowInitialRows(sdk.Bool(true))
			}
			if v, ok := d.GetOk("comment"); ok {
				req.WithComment(sdk.String(v.(string)))
			}
			err := client.Streams.CreateOnTable(ctx, req)
			if err != nil {
				return fmt.Errorf("error creating stream %v err = %w", name, err)
			}
		}
		break
	case onViewSet:
		viewId := helpers.DecodeSnowflakeID(onView.(string)).(sdk.SchemaObjectIdentifier)
		req := sdk.NewCreateStreamOnViewRequest(id, viewId)
		if appendOnly {
			req.WithAppendOnly(sdk.Bool(true))
		}
		if showInitialRows {
			req.WithShowInitialRows(sdk.Bool(true))
		}
		if v, ok := d.GetOk("comment"); ok {
			req.WithComment(sdk.String(v.(string)))
		}
		err := client.Streams.CreateOnView(ctx, req)
		if err != nil {
			return fmt.Errorf("error creating stream %v err = %w", name, err)
		}
		break
	case onStageSet:
		stageId := helpers.DecodeSnowflakeID(onStage.(string)).(sdk.SchemaObjectIdentifier)
		stageBuilder := snowflake.NewStageBuilder(stageId.Name(), stageId.DatabaseName(), stageId.SchemaName())
		sq := stageBuilder.Describe()
		stageDesc, err := snowflake.DescStage(db, sq)
		if err != nil {
			return err
		}
		if !strings.Contains(stageDesc.Directory, "ENABLE = true") {
			return fmt.Errorf("directory must be enabled on stage")
		}
		req := sdk.NewCreateStreamOnDirectoryTableRequest(id, stageId)
		if v, ok := d.GetOk("comment"); ok {
			req.WithComment(sdk.String(v.(string)))
		}
		err = client.Streams.CreateOnDirectoryTable(ctx, req)
		if err != nil {
			return fmt.Errorf("error creating stream %v err = %w", name, err)
		}
		break
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadStream(d, meta)
}

// ReadStream implements schema.ReadFunc.
func ReadStream(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	stream, err := client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(id))
	if err != nil {
		log.Printf("[DEBUG] stream (%s) not found", d.Id())
		d.SetId("")
		return nil
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
	case "Stage":
		if err := d.Set("on_stage", *stream.TableName); err != nil {
			return err
		}
	case "View":
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
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		if comment == "" {
			err := client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithUnsetComment(sdk.Bool(true)))
			if err != nil {
				return fmt.Errorf("error unsetting stream comment on %v", d.Id())
			}
		} else {
			err := client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithSetComment(sdk.String(comment)))
			if err != nil {
				return fmt.Errorf("error setting stream comment on %v", d.Id())
			}
		}
	}

	return ReadStream(d, meta)
}

// DeleteStream implements schema.DeleteFunc.
func DeleteStream(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	streamId := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(streamId))
	if err != nil {
		return fmt.Errorf("error deleting stream %v err = %w", d.Id(), err)
	}

	d.SetId("")

	return nil
}
