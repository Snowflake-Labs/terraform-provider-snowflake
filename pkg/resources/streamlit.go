package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var streamlitSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "String that specifies the identifier (i.e. name) for the streamlit; must be unique in your account.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the Cortex search service.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the streamlit.",
		ForceNew:    true,
	},
	"root_location": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the full path to the named stage containing the Streamlit Python files, media files, and the environment.yml file.",
		DiffSuppressFunc: SuppressIfAny(suppressLocationQuoting, IgnoreChangeToCurrentSnowflakeValueInOutput(DescribeOutputAttributeName, "root_location")),
	},
	"main_file": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the filename of the Streamlit Python application. This filename is relative to the value of `root_location`",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInOutput(DescribeOutputAttributeName, "main_file"),
	},
	"query_warehouse": {
		Type:             schema.TypeString,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		Optional:         true,
		Description:      "Specifies the warehouse where SQL queries issued by the Streamlit application are run.",
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInDescribe("query_warehouse")),
	},
	"external_access_integrations": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:         true,
		Description:      "External access integrations connected to the Streamlit.",
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInOutput(DescribeOutputAttributeName, "external_access_integrations")),
	},
	"title": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a title for the Streamlit app to display in Snowsight.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the streamlit.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW STREAMLIT` for the given streamli.",
		Elem: &schema.Resource{
			Schema: schemas.ShowStreamlitSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE STREAMLIT` for the given streamlit.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeStreamlitSchema,
		},
	},
}

func Streamlit() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextStreamlit,
		ReadContext:   ReadContextStreamlit(true),
		UpdateContext: UpdateContextStreamlit,
		DeleteContext: DeleteContextStreamlit,

		Schema: streamlitSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportStreamlit,
		},

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(ShowOutputAttributeName, "enabled", "scim_client", "comment"),
			ComputedIfAnyAttributeChanged(DescribeOutputAttributeName, "enabled", "comment", "network_policy", "run_as_role", "sync_password"),
		),
	}
}

func ImportStreamlit(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	logging.DebugLogger.Printf("[DEBUG] Starting streamlit import")
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	streamlit, err := client.Streamlits.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	streamlitDetails, err := client.Streamlits.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	if err = d.Set("name", streamlit.Name); err != nil {
		return nil, err
	}
	if err = d.Set("schema", streamlit.SchemaName); err != nil {
		return nil, err
	}
	if err = d.Set("root_location", streamlitDetails.RootLocation); err != nil {
		return nil, err
	}
	if err = d.Set("main_file", streamlitDetails.MainFile); err != nil {
		return nil, err
	}
	if err = d.Set("query_warehouse", streamlit.QueryWarehouse); err != nil {
		return nil, err
	}
	if err = d.Set("external_access_integrations", streamlitDetails.ExternalAccessIntegrations); err != nil {
		return nil, err
	}
	if err = d.Set("title", streamlit.Title); err != nil {
		return nil, err
	}
	if err = d.Set("comment", streamlit.Comment); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func CreateContextStreamlit(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)
	req := sdk.NewCreateStreamlitRequest(id, d.Get("root_location").(string), d.Get("main_file").(string))

	if v, ok := d.GetOk("query_warehouse"); ok {
		req.WithQueryWarehouse(sdk.NewAccountObjectIdentifier(v.(string)))
	}
	if v, ok := d.GetOk("title"); ok {
		req.WithTitle(v.(string))
	}
	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}
	if v, ok := d.GetOk("external_access_integrations"); ok {
		raw := expandStringList(v.(*schema.Set).List())
		integrations := make([]sdk.AccountObjectIdentifier, len(raw))
		for i, v := range raw {
			integrations[i] = sdk.NewAccountObjectIdentifier(v)
		}
		req.WithExternalAccessIntegrations(sdk.ExternalAccessIntegrationsRequest{
			ExternalAccessIntegrations: integrations,
		})
	}
	if err := client.Streamlits.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadContextStreamlit(false)(ctx, d, meta)
}

func ReadContextStreamlit(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

		streamlit, err := client.Streamlits.ShowByID(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query streamlit. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Streamlit name: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		streamlitDetails, err := client.Streamlits.Describe(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query streamlit properties. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Streamlit name: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		if err := d.Set("name", streamlit.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("schema", streamlit.SchemaName); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("root_location", streamlitDetails.RootLocation); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("main_file", streamlitDetails.MainFile); err != nil {
			return diag.FromErr(err)
		}
		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				showMapping{"query_warehouse", "query_warehouse", streamlit.QueryWarehouse, streamlit.QueryWarehouse, nil},
				showMapping{"external_access_integrations", "external_access_integrations", streamlitDetails.ExternalAccessIntegrations, streamlitDetails.ExternalAccessIntegrations, nil},
				showMapping{"title", "title", streamlit.Title, streamlit.Title, nil},
				showMapping{"comment", "comment", streamlit.Comment, streamlit.Comment, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, streamlitSchema, []string{
			"query_warehouse",
			"external_access_integrations",
			"title",
			"comment",
		}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.StreamlitToSchema(streamlit)}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(DescribeOutputAttributeName, []map[string]any{schemas.StreamlitPropertiesToSchema(*streamlitDetails)}); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func UpdateContextStreamlit(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	set, unset := sdk.NewStreamlitSetRequest(), sdk.NewStreamlitUnsetRequest()

	if d.HasChange("root_location") {
		// required field
		set.WithRootLocation(d.Get("root_location").(string))
	}

	if d.HasChange("main_file") {
		// required field
		set.WithMainFile(d.Get("main_file").(string))
	}

	if d.HasChange("title") {
		if v, ok := d.GetOk("title"); ok {
			set.WithTitle(v.(string))
		} else {
			unset.WithTitle(true)
		}
	}

	if d.HasChange("query_warehouse") {
		if v, ok := d.GetOk("query_warehouse"); ok {
			set.WithQueryWarehouse(sdk.NewAccountObjectIdentifier(v.(string)))
		} else {
			unset.WithQueryWarehouse(true)
		}
	}

	if d.HasChange("title") {
		if v, ok := d.GetOk("title"); ok {
			set.WithTitle(v.(string))
		} else {
			unset.WithTitle(true)
		}
	}

	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			set.WithComment(v.(string))
		} else {
			unset.WithComment(true)
		}
	}
	if v, ok := d.GetOk("external_access_integrations"); ok {
		raw := expandStringList(v.(*schema.Set).List())
		integrations := make([]sdk.AccountObjectIdentifier, len(raw))
		for i, v := range raw {
			integrations[i] = sdk.NewAccountObjectIdentifier(v)
		}
		set.WithExternalAccessIntegrations(sdk.ExternalAccessIntegrationsRequest{
			ExternalAccessIntegrations: integrations,
		})
	}

	if (*set != sdk.StreamlitSetRequest{}) {
		if err := client.Streamlits.Alter(ctx, sdk.NewAlterStreamlitRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	if (*unset != sdk.StreamlitUnsetRequest{}) {
		if err := client.Streamlits.Alter(ctx, sdk.NewAlterStreamlitRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextStreamlit(false)(ctx, d, meta)
}

func DeleteContextStreamlit(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	client := meta.(*provider.Context).Client

	err := client.Streamlits.Drop(ctx, sdk.NewDropStreamlitRequest(id).WithIfExists(true))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error deleting streamlit",
				Detail:   fmt.Sprintf("id %v err = %v", id.Name(), err),
			},
		}
	}

	d.SetId("")
	return nil
}
