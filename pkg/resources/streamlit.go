package resources

import (
	"context"
	"errors"
	"fmt"
	"path"

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
		Description: "String that specifies the identifier (i.e. name) for the streamlit; must be unique in your account.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the streamlit",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the streamlit.",
		ForceNew:    true,
	},
	"stage": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "The stage in which streamlit files are located.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInDescribe("root_location")),
	},
	"directory_location": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies the full path to the named stage containing the Streamlit Python files, media files, and the environment.yml file.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("root_location"),
	},
	"main_file": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the filename of the Streamlit Python application. This filename is relative to the value of `root_location`",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("main_file"),
	},
	"query_warehouse": {
		Type:             schema.TypeString,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		Optional:         true,
		Description:      "Specifies the warehouse where SQL queries issued by the Streamlit application are run.",
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInShow("query_warehouse")),
	},
	"external_access_integrations": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		},
		Optional:         true,
		Description:      "External access integrations connected to the Streamlit.",
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInDescribe("external_access_integrations")),
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
		Description: "Outputs the result of `SHOW STREAMLIT` for the given streamlit.",
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
		ReadContext:   ReadContextStreamlit,
		UpdateContext: UpdateContextStreamlit,
		DeleteContext: DeleteContextStreamlit,
		Description:   "Resource used to manage streamlits objects. For more information, check [streamlit documentation](https://docs.snowflake.com/en/sql-reference/commands-streamlit).",

		Schema: streamlitSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportStreamlit,
		},

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(ShowOutputAttributeName, "name", "title", "comment", "query_warehouse"),
			ComputedIfAnyAttributeChanged(DescribeOutputAttributeName, "name", "title", "comment", "root_location", "main_file", "query_warehouse", "external_access_integrations"),
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
	if err = d.Set("database", streamlit.DatabaseName); err != nil {
		return nil, err
	}
	if err = d.Set("schema", streamlit.SchemaName); err != nil {
		return nil, err
	}
	stageId, location, err := helpers.ParseRootLocation(streamlitDetails.RootLocation)
	if err != nil {
		return nil, err
	}
	if err := d.Set("stage", stageId.FullyQualifiedName()); err != nil {
		return nil, err
	}
	if err := d.Set("directory_location", location); err != nil {
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

func CreateContextStreamlit(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)
	stageId := sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(d.Get("stage").(string))
	rootLocation := fmt.Sprintf("@%s", stageId.FullyQualifiedName())
	if v, ok := d.GetOk("directory_location"); ok {
		rootLocation = path.Join(rootLocation, v.(string))
	}
	req := sdk.NewCreateStreamlitRequest(id, rootLocation, d.Get("main_file").(string))

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

	return ReadContextStreamlit(ctx, d, meta)
}

func ReadContextStreamlit(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		return diag.FromErr(err)
	}
	if err := d.Set("name", streamlit.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("schema", streamlit.SchemaName); err != nil {
		return diag.FromErr(err)
	}
	stageId, location, err := helpers.ParseRootLocation(streamlitDetails.RootLocation)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("stage", stageId.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("directory_location", location); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("main_file", streamlitDetails.MainFile); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("query_warehouse", streamlit.QueryWarehouse); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("external_access_integrations", streamlitDetails.ExternalAccessIntegrations); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("title", streamlit.Title); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("comment", streamlit.Comment); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.StreamlitToSchema(streamlit)}); err != nil {
		return diag.FromErr(err)
	}
	schemaDetails, err := schemas.StreamlitPropertiesToSchema(*streamlitDetails)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set(DescribeOutputAttributeName, []map[string]any{schemaDetails}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateContextStreamlit(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	set, unset := sdk.NewStreamlitSetRequest(), sdk.NewStreamlitUnsetRequest()

	if d.HasChange("name") {
		databaseName := d.Get("database").(string)
		schemaName := d.Get("schema").(string)
		name := d.Get("name").(string)
		newId := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

		err := client.Streamlits.Alter(ctx, sdk.NewAlterStreamlitRequest(id).WithRenameTo(newId))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeSnowflakeID(newId))
		id = newId
	}

	if d.HasChange("stage") || d.HasChange("directory_location") {
		stageId := sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(d.Get("stage").(string))
		rootLocation := fmt.Sprintf("@%s", stageId.FullyQualifiedName())
		if v, ok := d.GetOk("directory_location"); ok {
			rootLocation = path.Join(rootLocation, v.(string))
		}
		set.WithRootLocation(rootLocation)
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
	if d.HasChange("external_access_integrations") {
		raw := expandStringList(d.Get("external_access_integrations").(*schema.Set).List())
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

	return ReadContextStreamlit(ctx, d, meta)
}

func DeleteContextStreamlit(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
