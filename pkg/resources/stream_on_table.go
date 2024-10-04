package resources

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var streamOnTableSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the stream; must be unique for the database and schema in which the stream is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the stream."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the stream."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"copy_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Retains the access permissions from the original stream when a new stream is created using the OR REPLACE clause. Use only if the resource is already managed by Terraform. Otherwise, this field is skipped.",
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			return oldValue != "" && oldValue != newValue
		},
	},
	"table": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies an identifier for the table the stream will monitor."),
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInShow("table_name")),
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
	"append_only": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShowWithMapping("mode", func(x any) any {
			return x.(string) == string(sdk.StreamModeAppendOnly)
		}),
		Description: booleanStringFieldDescription("Specifies whether this is an append-only stream."),
	},
	"show_initial_rows": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		Description:      externalChangesNotDetectedFieldDescription(booleanStringFieldDescription("Specifies whether to return all existing rows in the source table as row inserts the first time the stream is consumed.")),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the stream.",
	},
	AtAttributeName:     atSchema,
	BeforeAttributeName: beforeSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW STREAMS` for the given stream.",
		Elem: &schema.Resource{
			Schema: schemas.ShowStreamSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE STREAM` for the given stream.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeStreamSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func StreamOnTable() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateStreamOnTable(false),
		ReadContext:   ReadStreamOnTable(true),
		UpdateContext: UpdateStreamOnTable,
		DeleteContext: DeleteStreamContext,
		Description:   "Resource used to manage streams on tables. For more information, check [stream documentation](https://docs.snowflake.com/en/sql-reference/sql/create-stream).",

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(streamOnTableSchema, ShowOutputAttributeName, "table", "append_only", "comment"),
			ComputedIfAnyAttributeChanged(streamOnTableSchema, DescribeOutputAttributeName, "table", "append_only", "comment"),
		),

		Schema: streamOnTableSchema,

		Importer: &schema.ResourceImporter{
			StateContext: ImportStreamOnTable,
		},
	}
}

func ImportStreamOnTable(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] Starting stream import")
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	v, err := client.Streams.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := d.Set("name", id.Name()); err != nil {
		return nil, err
	}
	if err := d.Set("database", id.DatabaseName()); err != nil {
		return nil, err
	}
	if err := d.Set("schema", id.SchemaName()); err != nil {
		return nil, err
	}
	if err := d.Set("append_only", booleanStringFromBool(v.IsAppendOnly())); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func CreateStreamOnTable(orReplace bool) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		databaseName := d.Get("database").(string)
		schemaName := d.Get("schema").(string)
		name := d.Get("name").(string)
		id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

		tableIdRaw := d.Get("table").(string)
		tableId, err := sdk.ParseSchemaObjectIdentifier(tableIdRaw)
		if err != nil {
			return diag.FromErr(err)
		}

		req := sdk.NewCreateOnTableStreamRequest(id, tableId)
		if orReplace {
			req.WithOrReplace(true)
			if d.Get("copy_grants").(bool) {
				req.WithCopyGrants(true)
			}
		}

		err = booleanStringAttributeCreate(d, "append_only", &req.AppendOnly)
		if err != nil {
			return diag.FromErr(err)
		}

		err = booleanStringAttributeCreate(d, "show_initial_rows", &req.ShowInitialRows)
		if err != nil {
			return diag.FromErr(err)
		}

		streamTimeTravelReq := handleStreamTimeTravel(d)
		if streamTimeTravelReq != nil {
			req.WithOn(*streamTimeTravelReq)
		}

		if v, ok := d.GetOk("comment"); ok {
			req.WithComment(v.(string))
		}

		err = client.Streams.CreateOnTable(ctx, req)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(helpers.EncodeResourceIdentifier(id))

		return ReadStreamOnTable(false)(ctx, d, meta)
	}
}

func ReadStreamOnTable(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		stream, err := client.Streams.ShowByID(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to stream. Marking the resource as removed.",
						Detail:   fmt.Sprintf("stream name: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}
		if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
			return diag.FromErr(err)
		}
		tableId, err := sdk.ParseSchemaObjectIdentifier(*stream.TableName)
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to parse table ID in Read.",
					Detail:   fmt.Sprintf("stream name: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		if err := d.Set("table", tableId.FullyQualifiedName()); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("comment", *stream.Comment); err != nil {
			return diag.FromErr(err)
		}
		streamDescription, err := client.Streams.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}
		if withExternalChangesMarking {
			var mode sdk.StreamMode
			if stream.Mode != nil {
				mode = *stream.Mode
			}
			if err = handleExternalChangesToObjectInShow(d,
				showMapping{"mode", "append_only", string(mode), booleanStringFromBool(stream.IsAppendOnly()), nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, streamOnTableSchema, []string{
			"append_only",
		}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.StreamToSchema(stream)}); err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set(DescribeOutputAttributeName, []map[string]any{schemas.StreamDescriptionToSchema(*streamDescription)}); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
}

func UpdateStreamOnTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// change on these fields can not be ForceNew because then the object is dropped explicitly and copying grants does not have effect
	if keys := changedKeys(d, "table", "append_only", "at", "before", "show_initial_rows"); len(keys) > 0 {
		log.Printf("[DEBUG] Detected change on %q, recreating...", keys)
		return CreateStreamOnTable(true)(ctx, d, meta)
	}

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		if comment == "" {
			err := client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithUnsetComment(true))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			err := client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithSetComment(comment))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return ReadStreamOnTable(false)(ctx, d, meta)
}
