package resources

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
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
		Description: "String that specifies the identifier (i.e. name) for the integration; must be unique in your account.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the streamlit.",
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
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInDescribe("query_warehouse")),
	},
	"external_access_integrations": {
		Type:             schema.TypeSet,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		Optional:         true,
		Description:      "External access integrations connected to the Streamlit.",
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInDescribe("query_warehouse")),
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
	logging.DebugLogger.Printf("[DEBUG] Starting scim integration import")
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	integrationProperties, err := client.SecurityIntegrations.Describe(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = d.Set("name", integration.Name); err != nil {
		return nil, err
	}
	if err = d.Set("enabled", integration.Enabled); err != nil {
		return nil, err
	}
	if scimClient, err := integration.SubType(); err == nil {
		if err = d.Set("scim_client", scimClient); err != nil {
			return nil, err
		}
	}
	if runAsRoleProperty, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "RUN_AS_ROLE" }); err == nil {
		if err = d.Set("run_as_role", runAsRoleProperty.Value); err != nil {
			return nil, err
		}
	}
	if networkPolicyProperty, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "NETWORK_POLICY" }); err == nil {
		if err = d.Set("network_policy", networkPolicyProperty.Value); err != nil {
			return nil, err
		}
	}
	if syncPasswordProperty, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "SYNC_PASSWORD" }); err == nil {
		if err = d.Set("sync_password", syncPasswordProperty.Value); err != nil {
			return nil, err
		}
	}
	if err = d.Set("comment", integration.Comment); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateContextStreamlit(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schemaName, name)
	req := sdk.NewCreateStreamlitRequest(id, d.Get("root_location").(string), d.Get("main_file").(string))

	if v, ok := d.GetOk("query_warehouse"); ok {
		req.WithWarehouse(sdk.NewAccountObjectIdentifier(v.(string)))
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

		integration, err := client.Streamlits.ShowByID(ctx, id)
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

		describeResult, err := client.Streamlits.Describe(ctx, id)
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

		if err := d.Set("name", integration.Name); err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			networkPolicyProperty, err := collections.FindOne(describeResult, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "NETWORK_POLICY" })
			if err != nil {
				return diag.FromErr(err)
			}

			syncPasswordProperty, err := collections.FindOne(describeResult, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "SYNC_PASSWORD" })
			if err != nil {
				return diag.FromErr(err)
			}

			if err = handleExternalChangesToObjectInDescribe(d,
				describeMapping{"network_policy", "network_policy", networkPolicyProperty.Value, networkPolicyProperty.Value, nil},
				describeMapping{"sync_password", "sync_password", syncPasswordProperty.Value, syncPasswordProperty.Value, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		// These are all identity sets, needed for the case where:
		// - previous config was empty (therefore Snowflake defaults had been used)
		// - new config have the same values that are already in SF
		if !d.GetRawConfig().IsNull() {
			if v := d.GetRawConfig().AsValueMap()["network_policy"]; !v.IsNull() {
				if err = d.Set("network_policy", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["sync_password"]; !v.IsNull() {
				if err = d.Set("sync_password", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.StreamlitToSchema(streamlit)}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ScimSecurityIntegrationPropertiesToSchema(describeResult)}); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func UpdateContextStreamlit(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	set, unset := sdk.NewStreamlitSetRequest(nil, nil), sdk.NewStreamlitUnsetRequest()

	if d.HasChange("enabled") {
		set.WithEnabled(d.Get("enabled").(bool))
	}

	if d.HasChange("network_policy") {
		if v := d.Get("network_policy").(string); v != "" {
			set.WithNetworkPolicy(sdk.NewAccountObjectIdentifier(v))
		} else {
			unset.WithNetworkPolicy(true)
		}
	}

	if d.HasChange("sync_password") {
		if v := d.Get("sync_password").(string); v != BooleanDefault {
			parsed, err := strconv.ParseBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithSyncPassword(parsed)
		} else {
			unset.WithSyncPassword(true)
		}
	}

	if d.HasChange("comment") {
		set.WithComment(sdk.StringAllowEmpty{Value: d.Get("comment").(string)})
	}

	if (*set != sdk.StreamlitSetRequest{}) {
		if err := client.SecurityIntegrations.AlterScim(ctx, sdk.NewAlterScimSecurityIntegrationRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	if (*unset != sdk.StreamlitUnsetRequest{}) {
		if err := client.SecurityIntegrations.AlterScim(ctx, sdk.NewAlterScimSecurityIntegrationRequest(id).WithUnset(*unset)); err != nil {
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
