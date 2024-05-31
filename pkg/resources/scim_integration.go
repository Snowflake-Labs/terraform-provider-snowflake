package resources

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
		Description: "Specify whether the security integration is enabled.",
	},
	"scim_client": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the client type for the scim integration",
		ValidateFunc: validation.StringInSlice(sdk.AsStringList([]sdk.ScimSecurityIntegrationScimClientOption{
			sdk.ScimSecurityIntegrationScimClientOkta, sdk.ScimSecurityIntegrationScimClientAzure, sdk.ScimSecurityIntegrationScimClientGeneric,
		}), true),
		DiffSuppressFunc: ignoreTrimSpaceSuppressFunc,
	},
	"run_as_role": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specify the SCIM role in Snowflake that owns any users and roles that are imported from the identity provider into Snowflake using SCIM.",
		ValidateFunc: validation.StringInSlice(sdk.AsStringList([]sdk.ScimSecurityIntegrationRunAsRoleOption{
			sdk.ScimSecurityIntegrationRunAsRoleOktaProvisioner, sdk.ScimSecurityIntegrationRunAsRoleAadProvisioner, sdk.ScimSecurityIntegrationRunAsRoleGenericScimProvisioner,
		}), true),
		DiffSuppressFunc: ignoreTrimSpaceSuppressFunc,
	},
	"network_policy": {
		Type:             schema.TypeString,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		Optional:         true,
		Description:      "Specifies an existing network policy that controls SCIM network traffic.",
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

// SCIMIntegration returns a pointer to the resource representing a network policy.
func SCIMIntegration() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextSCIMIntegration,
		ReadContext:   ReadContextSCIMIntegration,
		UpdateContext: UpdateContextSCIMIntegration,
		DeleteContext: DeleteContextSCIMIntegration,

		Schema: scimIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateContextSCIMIntegration implements schema.CreateFunc.
func CreateContextSCIMIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	enabled := d.Get("enabled").(bool)
	scimClient := d.Get("scim_client").(string)
	runAsRole := d.Get("run_as_role").(string)

	req := sdk.NewCreateScimSecurityIntegrationRequest(id, enabled, sdk.ScimSecurityIntegrationScimClientOption(scimClient), sdk.ScimSecurityIntegrationRunAsRoleOption(runAsRole))

	// Set optionals
	if v, ok := d.GetOk("network_policy"); ok {
		req.WithNetworkPolicy(sdk.NewAccountObjectIdentifier(v.(string)))
	}
	req.WithSyncPassword(d.Get("sync_password").(bool))
	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}
	if err := client.SecurityIntegrations.CreateScim(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadContextSCIMIntegration(ctx, d, meta)
}

// ReadContextSCIMIntegration implements schema.ReadFunc.
func ReadContextSCIMIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
	if err != nil {
		log.Printf("[DEBUG] notification integration (%s) not found", d.Id())
		d.SetId("")
		return diag.FromErr(err)
	}

	if c := integration.Category; c != "SECURITY" {
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
	typeParts := strings.Split(integration.IntegrationType, "-")
	if len(typeParts) < 2 {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid security integration type format.",
				Detail:   fmt.Sprintf("expected \"SCIM - <scim_client>\", got: %s", integration.IntegrationType),
			},
		}
	}
	scimClient := strings.TrimSpace(typeParts[1])
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
			if err := d.Set("network_policy", value); err != nil {
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

// UpdateContextSCIMIntegration implements schema.UpdateFunc.
func UpdateContextSCIMIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	var runSet, runUnset bool
	set, unset := sdk.NewScimIntegrationSetRequest(), sdk.NewScimIntegrationUnsetRequest()
	if d.HasChange("enabled") {
		runSet = true
		set.WithEnabled(d.Get("enabled").(bool))
	}
	if d.HasChange("network_policy") {
		networkPolicyID := sdk.NewAccountObjectIdentifier(d.Get("network_policy").(string))
		if networkPolicyID.Name() != "" {
			runSet = true
			set.WithNetworkPolicy(networkPolicyID)
		} else {
			runUnset = true
			unset.WithNetworkPolicy(true)
		}
	}
	if d.GetRawConfig().AsValueMap()["sync_password"].IsNull() {
		runUnset = true
		unset.WithSyncPassword(true)
	} else if d.HasChange("sync_password") {
		runSet = true
		set.WithSyncPassword(d.Get("sync_password").(bool))
	}

	if d.HasChange("comment") {
		runSet = true
		set.WithComment(sdk.StringAllowEmptyRequest{Value: d.Get("comment").(string)})
	}
	if runSet {
		if err := client.SecurityIntegrations.AlterScim(ctx, sdk.NewAlterScimSecurityIntegrationRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}
	if runUnset {
		if err := client.SecurityIntegrations.AlterScim(ctx, sdk.NewAlterScimSecurityIntegrationRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}
	return ReadContextSCIMIntegration(ctx, d, meta)
}

// DeleteContextSCIMIntegration implements schema.DeleteFunc.
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
