package resources

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var scimIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "String that specifies the identifier (i.e. name) for the integration; must be unique in your account.",
	},
	"enabled": {
		Type:             schema.TypeBool,
		Required:         true,
		Description:      "Specify whether the security integration is enabled.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("enabled"),
	},
	"scim_client": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      fmt.Sprintf("Specifies the client type for the scim integration. Valid options are: %v.", sdk.AsStringList(sdk.AllScimSecurityIntegrationScimClients)),
		ValidateDiagFunc: StringInSlice(sdk.AsStringList(sdk.AllScimSecurityIntegrationScimClients), true),
		DiffSuppressFunc: ignoreCaseAndTrimSpaceSuppressFunc,
	},
	"run_as_role": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		Description: fmt.Sprintf("Specify the SCIM role in Snowflake that owns any users and roles that are imported from the identity provider into Snowflake using SCIM."+
			" Provider assumes that the specified role is already provided. Valid options are: %v.", sdk.AllScimSecurityIntegrationRunAsRoles),
		ValidateDiagFunc: StringInSlice(sdk.AsStringList(sdk.AllScimSecurityIntegrationRunAsRoles), true),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalize := func(s string) string {
				return strings.ToUpper(strings.ReplaceAll(s, "-", ""))
			}
			return normalize(old) == normalize(new)
		},
	},
	"network_policy": {
		Type:             schema.TypeString,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		Optional:         true,
		Description:      "Specifies an existing network policy that controls SCIM network traffic.",
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeListValueInDescribe("network_policy")),
	},
	"sync_password": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("sync_password"),
		Description:      booleanStringFieldDescription("Specifies whether to enable or disable the synchronization of a user password from an Okta SCIM client as part of the API request to Snowflake."),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the integration.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW SECURITY INTEGRATIONS` for the given security integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowSecurityIntegrationSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE SECURITY INTEGRATIONS` for the given security integration.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeScimSecurityIntegrationSchema,
		},
	},
}

func SCIMIntegration() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: CreateContextSCIMIntegration,
		ReadContext:   ReadContextSCIMIntegration(true),
		UpdateContext: UpdateContextSCIMIntegration,
		DeleteContext: DeleteContextSCIMIntegration,

		Schema: scimIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportScimIntegration,
		},

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(ShowOutputAttributeName, "enabled", "scim_client", "comment"),
			ComputedIfAnyAttributeChanged(DescribeOutputAttributeName, "enabled", "comment", "network_policy", "run_as_role", "sync_password"),
		),

		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v092ScimIntegrationStateUpgrader,
			},
		},
	}
}

func ImportScimIntegration(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
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

func CreateContextSCIMIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := sdk.NewAccountObjectIdentifier(d.Get("name").(string))

	scimClient, err := sdk.ToScimSecurityIntegrationScimClientOption(d.Get("scim_client").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	runAsRole, err := sdk.ToScimSecurityIntegrationRunAsRoleOption(d.Get("run_as_role").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	req := sdk.NewCreateScimSecurityIntegrationRequest(id, scimClient, runAsRole).WithEnabled(d.Get("enabled").(bool))

	if v, ok := d.GetOk("network_policy"); ok {
		req.WithNetworkPolicy(sdk.NewAccountObjectIdentifier(v.(string)))
	}

	if v := d.Get("sync_password").(string); v != BooleanDefault {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return diag.FromErr(err)
		}

		req.WithSyncPassword(parsed)
	}

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	if err := client.SecurityIntegrations.CreateScim(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadContextSCIMIntegration(false)(ctx, d, meta)
}

func ReadContextSCIMIntegration(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

		integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query security integration. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Security integration name: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		integrationProperties, err := client.SecurityIntegrations.Describe(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query security integration properties. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Security integration name: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		if c := integration.Category; c != sdk.SecurityIntegrationCategory {
			return diag.FromErr(fmt.Errorf("expected %v to be a SECURITY integration, got %v", id, c))
		}

		if err := d.Set("name", integration.Name); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("enabled", integration.Enabled); err != nil {
			return diag.FromErr(err)
		}

		scimClient, err := integration.SubType()
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("scim_client", scimClient); err != nil {
			return diag.FromErr(err)
		}

		runAsRoleProperty, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "RUN_AS_ROLE" })
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("run_as_role", runAsRoleProperty.Value); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set("comment", integration.Comment); err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			networkPolicyProperty, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "NETWORK_POLICY" })
			if err != nil {
				return diag.FromErr(err)
			}

			syncPasswordProperty, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "SYNC_PASSWORD" })
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

		if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.SecurityIntegrationToSchema(integration)}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ScimSecurityIntegrationPropertiesToSchema(integrationProperties)}); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func UpdateContextSCIMIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	set, unset := sdk.NewScimIntegrationSetRequest(), sdk.NewScimIntegrationUnsetRequest()

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

	if (*set != sdk.ScimIntegrationSetRequest{}) {
		if err := client.SecurityIntegrations.AlterScim(ctx, sdk.NewAlterScimSecurityIntegrationRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	if (*unset != sdk.ScimIntegrationUnsetRequest{}) {
		if err := client.SecurityIntegrations.AlterScim(ctx, sdk.NewAlterScimSecurityIntegrationRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextSCIMIntegration(false)(ctx, d, meta)
}

func DeleteContextSCIMIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	client := meta.(*provider.Context).Client

	err := client.SecurityIntegrations.Drop(ctx, sdk.NewDropSecurityIntegrationRequest(sdk.NewAccountObjectIdentifier(id.Name())).WithIfExists(true))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error deleting integration",
				Detail:   fmt.Sprintf("id %v err = %v", id.Name(), err),
			},
		}
	}

	d.SetId("")
	return nil
}
