package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
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
		Description: blocklistedCharactersFieldDescription("Specifies the name of the OAuth integration. This name follows the rules for Object Identifiers. The name should be unique among security integrations in your account."),
	},
	"oauth_client": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      fmt.Sprintf("Creates an OAuth interface between Snowflake and a partner application. Valid options are: %v.", possibleValuesListed(sdk.AllOauthSecurityIntegrationClients)),
		ValidateDiagFunc: sdkValidation(sdk.ToOauthSecurityIntegrationClientOption),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToOauthSecurityIntegrationClientOption),
	},
	"oauth_redirect_uri": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: externalChangesNotDetectedFieldDescription("Specifies the client URI. After a user is authenticated, the web browser is redirected to this URI. The field should be only set when OAUTH_CLIENT = LOOKER. In any other case the field should be left out empty."),
	},
	"enabled": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("enabled"),
		Description:      booleanStringFieldDescription("Specifies whether this OAuth integration is enabled or disabled."),
	},
	"oauth_issue_refresh_tokens": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("oauth_issue_refresh_tokens"),
		Description:      booleanStringFieldDescription("Specifies whether to allow the client to exchange a refresh token for an access token when the current access token has expired."),
	},
	"oauth_refresh_token_validity": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("oauth_refresh_token_validity"),
		Description:      "Specifies how long refresh tokens should be valid (in seconds). OAUTH_ISSUE_REFRESH_TOKENS must be set to TRUE.",
	},
	"oauth_use_secondary_roles": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      fmt.Sprintf("Specifies whether default secondary roles set in the user properties are activated by default in the session being opened. Valid options are: %v.", possibleValuesListed(sdk.AllOauthSecurityIntegrationUseSecondaryRoles)),
		ValidateDiagFunc: sdkValidation(sdk.ToOauthSecurityIntegrationUseSecondaryRolesOption),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToOauthSecurityIntegrationUseSecondaryRolesOption), IgnoreChangeToCurrentSnowflakeListValueInDescribe("oauth_use_secondary_roles")),
	},
	"blocked_roles_list": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		},
		// TODO(SNOW-1517937): Check if can make optional
		Required:         true,
		Description:      "A set of Snowflake roles that a user cannot explicitly consent to using after authenticating.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("blocked_roles_list"),
	},
	"comment": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies a comment for the OAuth integration.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("comment"),
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW SECURITY INTEGRATION` for the given integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowSecurityIntegrationSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE SECURITY INTEGRATION` for the given integration.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeOauthIntegrationForPartnerApplications,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func OauthIntegrationForPartnerApplications() *schema.Resource {
	return &schema.Resource{
		Schema: oauthIntegrationForPartnerApplicationsSchema,

		CreateContext: CreateContextOauthIntegrationForPartnerApplications,
		ReadContext:   ReadContextOauthIntegrationForPartnerApplications(true),
		UpdateContext: UpdateContextOauthIntegrationForPartnerApplications,
		DeleteContext: DeleteContextSecurityIntegration,
		Description:   "Resource used to manage oauth security integration for partner applications objects. For more information, check [security integrations documentation](https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-oauth-snowflake).",

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(
				ShowOutputAttributeName,
				"name",
				"enabled",
				"comment",
			),
			ComputedIfAnyAttributeChanged(
				DescribeOutputAttributeName,
				"oauth_client",
				"enabled",
				"oauth_issue_refresh_tokens",
				"oauth_refresh_token_validity",
				"oauth_use_secondary_roles",
				"blocked_roles_list",
				"comment",
			),
		),

		Importer: &schema.ResourceImporter{
			StateContext: ImportOauthForPartnerApplicationIntegration,
		},
	}
}

func ImportOauthForPartnerApplicationIntegration(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	logging.DebugLogger.Printf("[DEBUG] Starting oauth integration for partner applications import")
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

	if err = d.Set("enabled", booleanStringFromBool(integration.Enabled)); err != nil {
		return nil, err
	}

	if issueRefreshTokens, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_ISSUE_REFRESH_TOKENS"
	}); err == nil {
		if err = d.Set("oauth_issue_refresh_tokens", issueRefreshTokens.Value); err != nil {
			return nil, err
		}
	}

	if refreshTokenValidity, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_REFRESH_TOKEN_VALIDITY"
	}); err == nil {
		refreshTokenValidityValue, err := strconv.ParseInt(refreshTokenValidity.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		if err = d.Set("oauth_refresh_token_validity", refreshTokenValidityValue); err != nil {
			return nil, err
		}
	}

	if oauthUseSecondaryRoles, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_USE_SECONDARY_ROLES"
	}); err == nil {
		oauthUseSecondaryRolesValue, err := sdk.ToOauthSecurityIntegrationUseSecondaryRolesOption(oauthUseSecondaryRoles.Value)
		if err != nil {
			return nil, err
		}
		if err = d.Set("oauth_use_secondary_roles", oauthUseSecondaryRolesValue); err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
}

func CreateContextOauthIntegrationForPartnerApplications(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := sdk.NewAccountObjectIdentifier(d.Get("name").(string))
	oauthClient, err := sdk.ToOauthSecurityIntegrationClientOption(d.Get("oauth_client").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	req := sdk.NewCreateOauthForPartnerApplicationsSecurityIntegrationRequest(id, oauthClient)

	if v, ok := d.GetOk("oauth_redirect_uri"); ok {
		req.WithOauthRedirectUri(v.(string))
	}

	if v := d.Get("enabled").(string); v != BooleanDefault {
		parsedBool, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithEnabled(parsedBool)
	}

	if v := d.Get("oauth_issue_refresh_tokens").(string); v != BooleanDefault {
		parsedBool, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithOauthIssueRefreshTokens(parsedBool)
	}

	if v := d.Get("oauth_refresh_token_validity").(int); v != IntDefault {
		req.WithOauthRefreshTokenValidity(v)
	}

	if v, ok := d.GetOk("oauth_use_secondary_roles"); ok {
		useSecondaryRolesOption, err := sdk.ToOauthSecurityIntegrationUseSecondaryRolesOption(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithOauthUseSecondaryRoles(useSecondaryRolesOption)
	}

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

	if err := client.SecurityIntegrations.CreateOauthForPartnerApplications(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadContextOauthIntegrationForPartnerApplications(false)(ctx, d, meta)
}

func ReadContextOauthIntegrationForPartnerApplications(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			return diag.FromErr(err)
		}

		if c := integration.Category; c != sdk.SecurityIntegrationCategory {
			return diag.FromErr(fmt.Errorf("expected %v to be a %s integration, got %v", id, sdk.SecurityIntegrationCategory, c))
		}
		if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("name", sdk.NewAccountObjectIdentifier(integration.Name).Name()); err != nil {
			return diag.FromErr(err)
		}

		oauthClient, err := integration.SubType()
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("oauth_client", oauthClient); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("comment", integration.Comment); err != nil {
			return diag.FromErr(err)
		}

		blockedRolesList, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
			return property.Name == "BLOCKED_ROLES_LIST"
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to find pre authorized roles list, err = %w", err))
		}
		var blockedRoles []string
		if len(blockedRolesList.Value) > 0 {
			blockedRoles = strings.Split(blockedRolesList.Value, ",")
		}
		if err := d.Set("blocked_roles_list", blockedRoles); err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				showMapping{"enabled", "enabled", integration.Enabled, booleanStringFromBool(integration.Enabled), nil},
			); err != nil {
				return diag.FromErr(err)
			}

			oauthIssueRefreshTokens, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "OAUTH_ISSUE_REFRESH_TOKENS"
			})
			if err != nil {
				return diag.FromErr(err)
			}

			oauthRefreshTokenValidity, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "OAUTH_REFRESH_TOKEN_VALIDITY"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			oauthRefreshTokenValidityValue, err := strconv.ParseInt(oauthRefreshTokenValidity.Value, 10, 64)
			if err != nil {
				return diag.FromErr(err)
			}

			oauthUseSecondaryRoles, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "OAUTH_USE_SECONDARY_ROLES"
			})
			if err != nil {
				return diag.FromErr(err)
			}

			if err = handleExternalChangesToObjectInDescribe(d,
				describeMapping{"oauth_issue_refresh_tokens", "oauth_issue_refresh_tokens", oauthIssueRefreshTokens.Value, oauthIssueRefreshTokens.Value, nil},
				describeMapping{"oauth_refresh_token_validity", "oauth_refresh_token_validity", oauthRefreshTokenValidity.Value, oauthRefreshTokenValidityValue, nil},
				describeMapping{"oauth_use_secondary_roles", "oauth_use_secondary_roles", oauthUseSecondaryRoles.Value, oauthUseSecondaryRoles.Value, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, oauthIntegrationForPartnerApplicationsSchema, []string{
			"enabled",
			"oauth_issue_refresh_tokens",
			"oauth_refresh_token_validity",
			"oauth_use_secondary_roles",
		}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.SecurityIntegrationToSchema(integration)}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(DescribeOutputAttributeName, []map[string]any{schemas.DescribeOauthIntegrationForPartnerApplicationsToSchema(integrationProperties)}); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
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
		// can call SET with an empty list
	}

	if d.HasChange("comment") {
		set.WithComment(d.Get("comment").(string))
		// TODO(SNOW-1515781): No UNSET
	}

	if d.HasChange("enabled") {
		if v := d.Get("oauth_issue_refresh_tokens").(string); v != BooleanDefault {
			parsedBool, err := booleanStringToBool(d.Get("enabled").(string))
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithEnabled(parsedBool)
		} else {
			unset.WithEnabled(true)
		}
	}

	if d.HasChange("oauth_issue_refresh_tokens") {
		if v := d.Get("oauth_issue_refresh_tokens").(string); v != BooleanDefault {
			parsedBool, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithOauthIssueRefreshTokens(parsedBool)
		} else {
			// TODO(SNOW-1515781): No UNSET
			set.WithOauthIssueRefreshTokens(true)
		}
	}

	if d.HasChange("oauth_redirect_uri") {
		// Field can only be set when oauth_client = LOOKER and is required (shouldn't be UNSET in those cases).
		// With any other case oauth_client, the field shouldn't be set.
		set.WithOauthRedirectUri(d.Get("oauth_redirect_uri").(string))
	}

	if d.HasChange("oauth_refresh_token_validity") {
		if v := d.Get("oauth_refresh_token_validity").(int); v != -1 {
			set.WithOauthRefreshTokenValidity(v)
		} else {
			// TODO(SNOW-1515781): No UNSET
			set.WithOauthRefreshTokenValidity(7776000)
		}
	}

	if d.HasChange("oauth_use_secondary_roles") {
		if v, ok := d.GetOk("oauth_use_secondary_roles"); ok {
			value, err := sdk.ToOauthSecurityIntegrationUseSecondaryRolesOption(v.(string))
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

	return ReadContextOauthIntegrationForPartnerApplications(false)(ctx, d, meta)
}

func DeleteContextSecurityIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	client := meta.(*provider.Context).Client

	err := client.SecurityIntegrations.Drop(ctx, sdk.NewDropSecurityIntegrationRequest(sdk.NewAccountObjectIdentifier(id.Name())).WithIfExists(true))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error deleting oauth integration for partner applications",
				Detail:   fmt.Sprintf("id %v err = %v", id.Name(), err),
			},
		}
	}

	d.SetId("")
	return nil
}
