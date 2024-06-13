package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var scimIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "String that specifies the identifier (i.e. name) for the integration; must be unique in your account.",
	},
	"enabled": {
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Specify whether the security integration is enabled. ",
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			return d.GetRawConfig().AsValueMap()["enabled"].IsNull()
		},
	},
	"scim_client": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      fmt.Sprintf("Specifies the client type for the scim integration. Valid options are: %v.", sdk.AsStringList(sdk.AllScimSecurityIntegrationScimClients)),
		ValidateFunc:     validation.StringInSlice(sdk.AsStringList(sdk.AllScimSecurityIntegrationScimClients), true),
		DiffSuppressFunc: ignoreCaseAndTrimSpaceSuppressFunc,
	},
	"run_as_role": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		Description: fmt.Sprintf("Specify the SCIM role in Snowflake that owns any users and roles that are imported from the identity provider into Snowflake using SCIM."+
			" Provider assumes that the specified role is already provided. Valid options are: %v.", sdk.AllScimSecurityIntegrationRunAsRoles),
		ValidateFunc: validation.StringInSlice(sdk.AsStringList(sdk.AllScimSecurityIntegrationRunAsRoles), true),
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
		DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
			return sdk.NewAccountObjectIdentifierFromFullyQualifiedName(old) == sdk.NewAccountObjectIdentifierFromFullyQualifiedName(new)
		},
	},
	"sync_password": {
		Type:        schema.TypeBool,
		Optional:    true,
		Computed:    true,
		Description: "Specifies whether to enable or disable the synchronization of a user password from an Okta SCIM client as part of the API request to Snowflake.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the integration.",
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the SCIM integration was created.",
	},
}

func SCIMIntegration() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: CreateContextSCIMIntegration,
		ReadContext:   ReadContextSCIMIntegration,
		UpdateContext: UpdateContextSCIMIntegration,
		DeleteContext: DeleteContextSCIMIntegration,
		CustomizeDiff: customdiff.All(
			BoolComputedIf("sync_password", func(client *sdk.Client, id sdk.AccountObjectIdentifier) (string, error) {
				props, err := client.SecurityIntegrations.Describe(context.Background(), id)
				if err != nil {
					return "", err
				}
				for _, prop := range props {
					if prop.GetName() == "SYNC_PASSWORD" {
						return prop.GetDefault(), nil
					}
				}
				return "", fmt.Errorf("")
			}),
		),
		Schema: scimIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v091ScimIntegrationStateUpgrader,
			},
		},
	}
}

func CreateContextSCIMIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	scimClientRaw := d.Get("scim_client").(string)
	runAsRoleRaw := d.Get("run_as_role").(string)
	scimClient, err := sdk.ToScimSecurityIntegrationScimClientOption(scimClientRaw)
	if err != nil {
		return diag.FromErr(err)
	}
	runAsRole, err := sdk.ToScimSecurityIntegrationRunAsRoleOption(runAsRoleRaw)
	if err != nil {
		return diag.FromErr(err)
	}
	req := sdk.NewCreateScimSecurityIntegrationRequest(id, scimClient, runAsRole).WithEnabled(d.Get("enabled").(bool))

	if v, ok := d.GetOk("network_policy"); ok {
		req.WithNetworkPolicy(sdk.NewAccountObjectIdentifier(v.(string)))
	}
	if !d.GetRawConfig().AsValueMap()["sync_password"].IsNull() {
		req.WithSyncPassword(d.Get("sync_password").(bool))
	}
	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}
	if err := client.SecurityIntegrations.CreateScim(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadContextSCIMIntegration(ctx, d, meta)
}

func ReadContextSCIMIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if c := integration.Category; c != sdk.SecurityIntegrationCategory {
		return diag.FromErr(fmt.Errorf("expected %v to be a SECURITY integration, got %v", id, c))
	}

	if err := d.Set("name", integration.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", integration.Comment); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_on", integration.CreatedOn.String()); err != nil {
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
	integrationProperties, err := client.SecurityIntegrations.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, property := range integrationProperties {
		name := property.Name
		value := property.Value
		switch name {
		case "ENABLED", "COMMENT":
			// We set this using the SHOW INTEGRATION call so let's ignore it here
		case "NETWORK_POLICY":
			networkPolicyID := sdk.NewAccountObjectIdentifier(value)
			if err := d.Set("network_policy", networkPolicyID.FullyQualifiedName()); err != nil {
				return diag.FromErr(err)
			}
		case "SYNC_PASSWORD":
			if err := d.Set("sync_password", helpers.StringToBool(value)); err != nil {
				return diag.FromErr(err)
			}
		case "RUN_AS_ROLE":
			if err := d.Set("run_as_role", value); err != nil {
				return diag.FromErr(err)
			}
		default:
			log.Printf("[WARN] unexpected property %v returned from Snowflake", name)
		}
	}

	return nil
}

func UpdateContextSCIMIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	set, unset := sdk.NewScimIntegrationSetRequest(), sdk.NewScimIntegrationUnsetRequest()
	if d.HasChange("enabled") {
		set.WithEnabled(d.Get("enabled").(bool))
	}
	if d.HasChange("network_policy") {
		networkPolicyID := sdk.NewAccountObjectIdentifier(d.Get("network_policy").(string))
		if networkPolicyID.Name() != "" {
			set.WithNetworkPolicy(networkPolicyID)
		} else {
			unset.WithNetworkPolicy(true)
		}
	}
	if !d.GetRawConfig().AsValueMap()["sync_password"].IsNull() {
		set.WithSyncPassword(d.Get("sync_password").(bool))
	} else {
		unset.WithSyncPassword(true)
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
	return ReadContextSCIMIntegration(ctx, d, meta)
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
