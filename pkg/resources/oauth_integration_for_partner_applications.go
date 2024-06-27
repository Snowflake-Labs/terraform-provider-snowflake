package resources

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var oauthIntegrationForPartnerApplicationsSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the name of the SAML2 integration. This name follows the rules for Object Identifiers. The name should be unique among security integrations in your account.",
	},
	"oauth_client": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      fmt.Sprintf("Creates an OAuth interface between Snowflake and a partner application. Valid options are: %v", sdk.AllOauthSecurityIntegrationClients),
		ValidateFunc:     validation.StringInSlice(sdk.AsStringList(sdk.AllOauthSecurityIntegrationClients), false),
		DiffSuppressFunc: ignoreCaseAndTrimSpaceSuppressFunc,
	},
	"oauth_redirect_uri": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the client URI. After a user is authenticated, the web browser is redirected to this URI.",
	},
	"enabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether this OAuth integration is enabled or disabled.",
	},
	"oauth_issue_refresh_tokens": {
		Type:         schema.TypeString,
		ValidateFunc: validation.StringInSlice([]string{"true", "false"}, true),
		Default:      "unknown",
		Optional:     true,
		Description:  "Specifies whether to allow the client to exchange a refresh token for an access token when the current access token has expired.",
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			return oldValue == "true" && newValue == "unknown"
		},
	},
	"oauth_refresh_token_validity": {
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(1),
		Description:  "Specifies how long refresh tokens should be valid (in seconds). OAUTH_ISSUE_REFRESH_TOKENS must be set to TRUE.",
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			return d.Get(k).(int) == 7776000 && newValue == "0"
		},
	},
	"oauth_use_secondary_roles": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "Specifies whether default secondary roles set in the user properties are activated by default in the session being opened.",
		ValidateFunc: validation.StringInSlice(sdk.AsStringList(sdk.AllOauthSecurityIntegrationUseSecondaryRoles), false),
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			return strings.EqualFold(oldValue, newValue) || d.Get(k).(string) == string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone) && newValue == ""
		},
	},
	"blocked_roles_list": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Computed:    true,
		Description: "List of roles that a user cannot explicitly consent to using after authenticating. Do not include ACCOUNTADMIN, ORGADMIN or SECURITYADMIN as they are already implicitly enforced and will cause in-place updates.",
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			value := d.Get("oauth_add_privileged_roles_to_blocked_list").(bool)
			if !value {
				return old == new
			}
			return old == "ACCOUNTADMIN" || old == "SECURITYADMIN"
		},
	},
	"oauth_add_privileged_roles_to_blocked_list": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"oauth_authorization_endpoint": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Authorization endpoint for oauth.",
	},
	"oauth_token_endpoint": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Token endpoint for oauth.",
	},
	"oauth_allowed_authorization_endpoints": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Computed:    true,
		Description: "A list of allowed authorization endpoints for oauth.",
	},
	"oauth_allowed_token_endpoints": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Computed:    true,
		Description: "A list of allowed token endpoints for oauth.",
	},
	"oauth_client_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Oauth client ID.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the OAuth integration.",
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the OAuth integration was created.",
	},
}

func OauthIntegrationForPartnerApplications() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextOauthIntegrationForPartnerApplications,
		ReadContext:   ReadContextOauthIntegrationForPartnerApplications,
		UpdateContext: UpdateContextOauthIntegrationForPartnerApplications,
		DeleteContext: DeleteContextSecurityIntegration,
		Schema:        oauthIntegrationForPartnerApplicationsSchema,
		CustomizeDiff: customdiff.All(
		// SuppressIfParameterSet("blocked_roles_list", "OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST"),
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextOauthIntegrationForPartnerApplications(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	oauthClientRaw := d.Get("oauth_client").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	oauthClient, err := sdk.ToOauthSecurityIntegrationClientOption(oauthClientRaw)
	if err != nil {
		return diag.FromErr(err)
	}
	req := sdk.NewCreateOauthForPartnerApplicationsSecurityIntegrationRequest(id, oauthClient)

	if v, ok := d.GetOk("blocked_roles_list"); ok {
		elems := expandStringList(v.(*schema.Set).List())
		blockedRoles := make([]sdk.AccountObjectIdentifier, len(elems))
		for i := range elems {
			blockedRoles[i] = sdk.NewAccountObjectIdentifier(elems[i])
		}
		req.WithBlockedRolesList(sdk.BlockedRolesListRequest{BlockedRolesList: blockedRoles})
	}

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	if v, ok := d.GetOk("enabled"); ok {
		req.WithEnabled(v.(bool))
	}
	if v := d.Get("oauth_issue_refresh_tokens").(string); v != "unknown" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithOauthIssueRefreshTokens(parsed)
	}

	if v, ok := d.GetOk("oauth_redirect_uri"); ok {
		req.WithOauthRedirectUri(v.(string))
	}

	if v, ok := d.GetOk("oauth_refresh_token_validity"); ok {
		req.WithOauthRefreshTokenValidity(v.(int))
	}

	if v, ok := d.GetOk("oauth_use_secondary_roles"); ok {
		valueRaw := v.(string)
		value, err := sdk.ToOauthSecurityIntegrationUseSecondaryRolesOption(valueRaw)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithOauthUseSecondaryRoles(value)
	}

	if err := client.SecurityIntegrations.CreateOauthForPartnerApplications(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	return ReadContextOauthIntegrationForPartnerApplications(ctx, d, meta)
}

func ReadContextOauthIntegrationForPartnerApplications(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
	if err != nil {
		log.Printf("[DEBUG] OauthIntegrationForPartnerApplications (%s) not found", d.Id())
		d.SetId("")
		return diag.FromErr(err)
	}

	if c := integration.Category; c != sdk.SecurityIntegrationCategory {
		return diag.FromErr(fmt.Errorf("expected %v to be a %s integration, got %v", id, sdk.SecurityIntegrationCategory, c))
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
	oauthClient, err := integration.SubType()
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("oauth_client", oauthClient); err != nil {
		return diag.FromErr(err)
	}

	properties, err := client.SecurityIntegrations.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	defaults := make(map[string]string)
	for _, property := range properties {
		name := property.Name
		value := property.Value
		switch name {
		case "BLOCKED_ROLES_LIST":
			var blockedRoles []string
			if len(value) > 0 {
				blockedRoles = strings.Split(value, ",")
			}

			if err := d.Set("blocked_roles_list", blockedRoles); err != nil {
				return diag.FromErr(err)
			}
		case "COMMENT":
			if err := d.Set("comment", value); err != nil {
				return diag.FromErr(err)
			}
		case "CREATED_ON":
			if err := d.Set("created_on", value); err != nil {
				return diag.FromErr(err)
			}
		case "ENABLED":
			if err := d.Set("enabled", helpers.StringToBool(value)); err != nil {
				return diag.FromErr(err)
			}
		case "OAUTH_CLIENT":
			if err := d.Set("oauth_client", value); err != nil {
				return diag.FromErr(err)
			}
		case "OAUTH_ISSUE_REFRESH_TOKENS":
			defaults["OAUTH_ISSUE_REFRESH_TOKENS"] = property.Default
			if err := d.Set("oauth_issue_refresh_tokens", value); err != nil {
				return diag.FromErr(err)
			}
		case "OAUTH_REDIRECT_URI":
			if err := d.Set("oauth_redirect_uri", value); err != nil {
				return diag.FromErr(err)
			}
		case "OAUTH_REFRESH_TOKEN_VALIDITY":
			v, err := strconv.Atoi(value)
			if err != nil {
				return diag.FromErr(err)
			}
			defaults["OAUTH_REFRESH_TOKEN_VALIDITY"] = property.Default
			if err := d.Set("oauth_refresh_token_validity", v); err != nil {
				return diag.FromErr(err)
			}
		case "OAUTH_USE_SECONDARY_ROLES":
			if err := d.Set("oauth_use_secondary_roles", value); err != nil {
				return diag.FromErr(err)
			}
		case "OAUTH_AUTHORIZATION_ENDPOINT":
			if err := d.Set("oauth_authorization_endpoint", value); err != nil {
				return diag.FromErr(err)
			}
		case "OAUTH_TOKEN_ENDPOINT":
			if err := d.Set("oauth_token_endpoint", value); err != nil {
				return diag.FromErr(err)
			}
		case "OAUTH_ALLOWED_AUTHORIZATION_ENDPOINTS":
			var elems []string
			if len(value) > 0 {
				elems = strings.Split(value, ",")
			}

			if err := d.Set("oauth_allowed_authorization_endpoints", elems); err != nil {
				return diag.FromErr(err)
			}
		case "OAUTH_ALLOWED_TOKEN_ENDPOINTS":
			var elems []string
			if len(value) > 0 {
				elems = strings.Split(value, ",")
			}
			if err := d.Set("oauth_allowed_token_endpoints", elems); err != nil {
				return diag.FromErr(err)
			}
		case "OAUTH_CLIENT_ID":
			if err := d.Set("oauth_client_id", value); err != nil {
				return diag.FromErr(err)
			}

		default:
			log.Printf("[WARN] unexpected property %v returned from Snowflake", name)
		}
	}
	paramRaw, err := getParameterInAccount(ctx, client, "OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST")
	if err != nil {
		return nil
	}
	param := helpers.StringToBool(paramRaw)
	if err := d.Set("oauth_add_privileged_roles_to_blocked_list", param); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func UpdateContextOauthIntegrationForPartnerApplications(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	set, unset := sdk.NewOauthForPartnerApplicationsIntegrationSetRequest(), sdk.NewOauthForPartnerApplicationsIntegrationUnsetRequest()

	if d.HasChange("blocked_roles_list") {
		elems := expandStringList(d.Get("blocked_roles_list").(*schema.Set).List())
		blockedRoles := make([]sdk.AccountObjectIdentifier, len(elems))
		for i := range elems {
			blockedRoles[i] = sdk.NewAccountObjectIdentifier(elems[i])
		}
		set.WithBlockedRolesList(sdk.BlockedRolesListRequest{BlockedRolesList: blockedRoles})
	}

	if d.HasChange("comment") {
		set.WithComment(d.Get("comment").(string))
	}

	if d.HasChange("enabled") {
		set.WithEnabled(d.Get("enabled").(bool))
	}

	if d.HasChange("oauth_issue_refresh_tokens") {
		if v := d.Get("oauth_issue_refresh_tokens").(string); v != "unknown" {
			parsed, err := strconv.ParseBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithOauthIssueRefreshTokens(parsed)
		} else {
			// TODO: fix
			set.WithOauthIssueRefreshTokens(true)
		}
	}

	if d.HasChange("oauth_redirect_uri") {
		set.WithOauthRedirectUri(d.Get("oauth_redirect_uri").(string))
	}

	if d.HasChange("oauth_refresh_token_validity") {
		v := d.Get("oauth_refresh_token_validity").(int)
		if v > 0 {
			set.WithOauthRefreshTokenValidity(v)
		} else {
			// TODO: fix
			// TODO: better logic, like in docs
			set.WithOauthRefreshTokenValidity(7776000)
		}
	}

	if d.HasChange("oauth_use_secondary_roles") {
		v := d.Get("oauth_use_secondary_roles").(string)
		if len(v) > 0 {
			value, err := sdk.ToOauthSecurityIntegrationUseSecondaryRolesOption(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithOauthUseSecondaryRoles(value)
		} else {
			unset.WithOauthUseSecondaryRoles(true)
		}
	}
	if !reflect.DeepEqual(*set, sdk.OauthForPartnerApplicationsIntegrationSetRequest{}) {
		if err := client.SecurityIntegrations.AlterOauthForPartnerApplications(ctx, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}
	if !reflect.DeepEqual(*unset, sdk.OauthForPartnerApplicationsIntegrationUnsetRequest{}) {
		if err := client.SecurityIntegrations.AlterOauthForPartnerApplications(ctx, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}
	return ReadContextOauthIntegrationForPartnerApplications(ctx, d, meta)
}
