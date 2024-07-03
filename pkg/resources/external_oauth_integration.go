package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"
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

var oauthExternalIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the name of the External Oath integration. This name follows the rules for Object Identifiers. The name should be unique among security integrations in your account.",
	},
	"external_oauth_type": {
		Type:         schema.TypeString,
		Required:     true,
		Description:  fmt.Sprintf("Specifies the OAuth 2.0 authorization server to be Okta, Microsoft Azure AD, Ping Identity PingFederate, or a Custom OAuth 2.0 authorization server. Valid options are: %v", sdk.AllExternalOauthSecurityIntegrationTypes),
		ValidateFunc: validation.StringInSlice(sdk.AsStringList(sdk.AllExternalOauthSecurityIntegrationTypes), true),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalize := func(s string) string {
				return strings.ToUpper(strings.ReplaceAll(s, "-", ""))
			}
			return normalize(old) == normalize(new)
		},
	},
	"enabled": {
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Specifies whether to initiate operation of the integration or suspend it.",
	},
	"external_oauth_issuer": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the URL to define the OAuth 2.0 authorization server.",
	},
	"external_oauth_token_user_mapping_claim": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: "Specifies the access token claim or claims that can be used to map the access token to a Snowflake user record.",
	},
	"external_oauth_snowflake_user_mapping_attribute": {
		Type:         schema.TypeString,
		Required:     true,
		Description:  fmt.Sprintf("Indicates which Snowflake user record attribute should be used to map the access token to a Snowflake user record. Valid options are: %v", sdk.AllExternalOauthSecurityIntegrationSnowflakeUserMappingAttributes),
		ValidateFunc: validation.StringInSlice(sdk.AsStringList(sdk.AllExternalOauthSecurityIntegrationSnowflakeUserMappingAttributes), true),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalize := func(s string) string {
				return strings.ToUpper(strings.ReplaceAll(s, "-", ""))
			}
			return normalize(old) == normalize(new)
		},
	},
	"external_oauth_jws_keys_url": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		MaxItems:      3,
		Optional:      true,
		ConflictsWith: []string{"external_oauth_rsa_public_key", "external_oauth_rsa_public_key_2"},
		Description:   "Specifies the endpoint or a list of endpoints from which to download public keys or certificates to validate an External OAuth access token. The maximum number of URLs that can be specified in the list is 3.",
	},
	"external_oauth_rsa_public_key": {
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{"external_oauth_jws_keys_url"},
		Description:   "Specifies a Base64-encoded RSA public key, without the -----BEGIN PUBLIC KEY----- and -----END PUBLIC KEY----- headers.",
	},
	"external_oauth_rsa_public_key_2": {
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "Specifies a second RSA public key, without the -----BEGIN PUBLIC KEY----- and -----END PUBLIC KEY----- headers. Used for key rotation.",
		ConflictsWith: []string{"external_oauth_jws_keys_url"},
	},
	"external_oauth_blocked_roles_list": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies the list of roles that a client cannot set as the primary role. By default, this list includes the ACCOUNTADMIN, ORGADMIN, and SECURITYADMIN roles. To remove these privileged roles from the list, use the ALTER ACCOUNT command to set the EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST account parameter to FALSE.",
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			value := d.Get("external_oauth_add_privileged_roles_to_blocked_list").(bool)
			if !value {
				return old == new
			}
			return old == "ACCOUNTADMIN" || old == "SECURITYADMIN"
		},
		ConflictsWith: []string{"external_oauth_allowed_roles_list"},
	},
	"external_oauth_allowed_roles_list": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		Description:   "Specifies the list of roles that the client can set as the primary role.",
		ConflictsWith: []string{"external_oauth_blocked_roles_list"},
	},
	"external_oauth_audience_list": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies additional values that can be used for the access token's audience validation on top of using the Customer's Snowflake Account URL ",
	},
	"external_oauth_any_role_mode": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  fmt.Sprintf("Specifies whether the OAuth client or user can use a role that is not defined in the OAuth access token. Valid options are: %v", sdk.AsStringList(sdk.AllExternalOauthSecurityIntegrationAnyRoleModes)),
		ValidateFunc: validation.StringInSlice(sdk.AsStringList(sdk.AllExternalOauthSecurityIntegrationAnyRoleModes), true),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalize := func(s string) string {
				return strings.ToUpper(strings.ReplaceAll(s, "-", ""))
			}
			return normalize(old) == normalize(new)
		},
	},
	"external_oauth_scope_delimiter": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the scope delimiter in the authorization token.",
	},
	"external_oauth_scope_mapping_attribute": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the access token claim to map the access token to an account role.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the OAuth integration.",
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
			Schema: schemas.DescribeExternalOauthSecurityIntegrationSchema,
		},
	},
}

func ExternalOauthIntegration() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextExternalOauthIntegration,
		ReadContext:   ReadContextExternalOauthIntegration(true),
		UpdateContext: UpdateContextExternalOauthIntegration,
		DeleteContext: DeleteContextExternalOauthIntegration,
		Schema:        oauthExternalIntegrationSchema,
		CustomizeDiff: customdiff.All(
			ForceNewIfChangeToEmptyString("external_oauth_snowflake_user_mapping_attribute"),
			ForceNewIfChangeToEmptyString("external_oauth_rsa_public_key"),
			ForceNewIfChangeToEmptyString("external_oauth_rsa_public_key_2"),
			ForceNewIfChangeToEmptyString("external_oauth_scope_mapping_attribute"),
			ForceNewIfChangeToEmptySet[any]("external_oauth_jws_keys_url"),
			// ModifyStateIfParameterSet("external_oauth_blocked_roles_list", "EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST", func(d *schema.ResourceDiff) error {
			// 	allowed := d.Get("external_oauth_allowed_roles_list").(*schema.Set)
			// 	if allowed.Len() > 0 {
			// 		return nil
			// 	}
			// 	set := d.Get("external_oauth_blocked_roles_list").(*schema.Set)
			// 	set.Add("ACCOUNTADMIN")
			// 	set.Add("SECURITYADMIN")
			// 	return d.SetNew("external_oauth_blocked_roles_list", set)
			// }),
			ComputedIfAnyAttributeChanged(ShowOutputAttributeName, "enabled", "external_oauth_type", "comment"),
			ComputedIfAnyAttributeChanged(DescribeOutputAttributeName, "enabled", "external_oauth_issuer", "external_oauth_jws_keys_url", "external_oauth_any_role_mode",
				"external_oauth_rsa_public_key", "external_oauth_rsa_public_key_2", "external_oauth_blocked_roles_list", "external_oauth_allowed_roles_list",
				"external_oauth_audience_list", "external_oauth_token_user_mapping_claim", "external_oauth_snowflake_user_mapping_attribute", "external_oauth_scope_delimiter",
				"comment"),
		),
		Importer: &schema.ResourceImporter{
			StateContext: ImportExternalOauthIntegration,
		},
	}
}

func ImportExternalOauthIntegration(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	logging.DebugLogger.Printf("[DEBUG] Starting external oauth integration import")
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
	if oauthType, err := integration.SubType(); err == nil {
		if err = d.Set("external_oauth_type", oauthType); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "EXTERNAL_OAUTH_ISSUER" }); err == nil {
		if err = d.Set("external_oauth_issuer", prop.Value); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_JWS_KEYS_URL"
	}); err == nil {
		if err = d.Set("external_oauth_jws_keys_url", listValueToSlice(prop.Value, true)); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_ANY_ROLE_MODE"
	}); err == nil {
		if err = d.Set("external_oauth_any_role_mode", prop.Value); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_RSA_PUBLIC_KEY"
	}); err == nil {
		if err = d.Set("external_oauth_rsa_public_key", prop.Value); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2"
	}); err == nil {
		if err = d.Set("external_oauth_rsa_public_key_2", prop.Value); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_BLOCKED_ROLES_LIST"
	}); err == nil {
		if err = d.Set("external_oauth_blocked_roles_list", listValueToSlice(prop.Value, false)); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_ALLOWED_ROLES_LIST"
	}); err == nil {
		if err = d.Set("external_oauth_allowed_roles_list", listValueToSlice(prop.Value, false)); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_AUDIENCE_LIST"
	}); err == nil {
		if err = d.Set("external_oauth_audience_list", listValueToSlice(prop.Value, true)); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM"
	}); err == nil {
		if err = d.Set("external_oauth_token_user_mapping_claim", listValueToSlice(prop.Value, true)); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE"
	}); err == nil {
		if err = d.Set("external_oauth_snowflake_user_mapping_attribute", prop.Value); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_SCOPE_DELIMITER"
	}); err == nil {
		if err = d.Set("external_oauth_scope_delimiter", prop.Value); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "COMMENT" }); err == nil {
		if err = d.Set("comment", prop.Value); err != nil {
			return nil, err
		}
	}
	return []*schema.ResourceData{d}, nil
}

func CreateContextExternalOauthIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	enabled := d.Get("enabled").(bool)
	externalOauthIssuer := d.Get("external_oauth_issuer").(string)
	externalOauthSnowflakeUserMappingAttributeRaw := d.Get("external_oauth_snowflake_user_mapping_attribute").(string)
	externalOauthTokenUserMappingClaimRaw := expandStringList(d.Get("external_oauth_token_user_mapping_claim").(*schema.Set).List())
	name := d.Get("name").(string)
	integrationTypeRaw := d.Get("external_oauth_type").(string)
	integrationType, err := sdk.ToExternalOauthSecurityIntegrationTypeOption(integrationTypeRaw)
	if err != nil {
		return diag.FromErr(err)
	}
	externalOauthSnowflakeUserMappingAttribute, err := sdk.ToExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption(externalOauthSnowflakeUserMappingAttributeRaw)
	if err != nil {
		return diag.FromErr(err)
	}
	externalOauthTokenUserMappingClaim := make([]sdk.TokenUserMappingClaim, 0, len(externalOauthTokenUserMappingClaimRaw))
	for _, v := range externalOauthTokenUserMappingClaimRaw {
		externalOauthTokenUserMappingClaim = append(externalOauthTokenUserMappingClaim, sdk.TokenUserMappingClaim{Claim: v})
	}
	id := sdk.NewAccountObjectIdentifier(name)
	req := sdk.NewCreateExternalOauthSecurityIntegrationRequest(id, enabled, integrationType, externalOauthIssuer, externalOauthTokenUserMappingClaim, externalOauthSnowflakeUserMappingAttribute)

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	if v, ok := d.GetOk("external_oauth_allowed_roles_list"); ok {
		req.WithExternalOauthAllowedRolesList(sdk.AllowedRolesListRequest{AllowedRolesList: expandObjectIdentifierList(v.(*schema.Set).List())})
	}

	if v, ok := d.GetOk("external_oauth_any_role_mode"); ok {
		valueRaw := v.(string)
		value, err := sdk.ToExternalOauthSecurityIntegrationAnyRoleModeOption(valueRaw)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithExternalOauthAnyRoleMode(value)
	}

	if v, ok := d.GetOk("external_oauth_audience_list"); ok {
		elems := expandStringList(v.(*schema.Set).List())
		audienceUrls := make([]sdk.AudienceListItem, len(elems))
		for i := range elems {
			audienceUrls[i] = sdk.AudienceListItem{Item: elems[i]}
		}
		req.WithExternalOauthAudienceList(sdk.AudienceListRequest{AudienceList: audienceUrls})
	}

	if v, ok := d.GetOk("external_oauth_blocked_roles_list"); ok {
		if _, okAllowed := d.GetOk("external_oauth_allowed_roles_list"); !okAllowed {
			req.WithExternalOauthBlockedRolesList(sdk.BlockedRolesListRequest{BlockedRolesList: expandObjectIdentifierList(v.(*schema.Set).List())})
		}
	}

	if v, ok := d.GetOk("external_oauth_jws_keys_url"); ok {
		elems := expandStringList(v.(*schema.Set).List())
		urls := make([]sdk.JwsKeysUrl, len(elems))
		for i := range elems {
			urls[i] = sdk.JwsKeysUrl{JwsKeyUrl: elems[i]}
		}
		req.WithExternalOauthJwsKeysUrl(urls)
	}

	if v, ok := d.GetOk("external_oauth_rsa_public_key"); ok {
		req.WithExternalOauthRsaPublicKey(v.(string))
	}

	if v, ok := d.GetOk("external_oauth_rsa_public_key_2"); ok {
		req.WithExternalOauthRsaPublicKey2(v.(string))
	}

	if v, ok := d.GetOk("external_oauth_scope_delimiter"); ok {
		req.WithExternalOauthScopeDelimiter(v.(string))
	}

	if v, ok := d.GetOk("external_oauth_scope_mapping_attribute"); ok {
		req.WithExternalOauthScopeMappingAttribute(v.(string))
	}

	if err := client.SecurityIntegrations.CreateExternalOauth(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadContextExternalOauthIntegration(false)(ctx, d, meta)
}

func ReadContextExternalOauthIntegration(withExternalChangesMarking bool) schema.ReadContextFunc {
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
		if err := d.Set("name", integration.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("comment", integration.Comment); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("enabled", integration.Enabled); err != nil {
			return diag.FromErr(err)
		}
		subType, err := integration.SubType()
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("external_oauth_type", subType); err != nil {
			return diag.FromErr(err)
		}
		if withExternalChangesMarking {
			externalOauthIssuer, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "EXTERNAL_OAUTH_ISSUER" })
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthJwsKeysUrl, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_JWS_KEYS_URL"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthAnyRoleMode, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_ANY_ROLE_MODE"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthRsaPublicKey, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_RSA_PUBLIC_KEY"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthRsaPublicKey2, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthBlockedRolesList, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_BLOCKED_ROLES_LIST"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthAllowedRolesList, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_ALLOWED_ROLES_LIST"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthAudienceList, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_AUDIENCE_LIST"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthTokenUserMappingClaim, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthSnowflakeUserMappingAttribute, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthScopeDelimiter, err := collections.FindOne(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_SCOPE_DELIMITER"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			if err = handleExternalChangesToObjectInDescribe(d,
				describeMapping{"external_oauth_issuer", "external_oauth_issuer", externalOauthIssuer.Value, externalOauthIssuer.Value, nil},
				describeMapping{"external_oauth_jws_keys_url", "external_oauth_jws_keys_url", externalOauthJwsKeysUrl.Value, externalOauthJwsKeysUrl.Value, nil},
				describeMapping{"external_oauth_any_role_mode", "external_oauth_any_role_mode", externalOauthAnyRoleMode.Value, externalOauthAnyRoleMode.Value, nil},
				describeMapping{"external_oauth_rsa_public_key", "external_oauth_rsa_public_key", externalOauthRsaPublicKey.Value, externalOauthRsaPublicKey.Value, nil},
				describeMapping{"external_oauth_rsa_public_key_2", "external_oauth_rsa_public_key_2", externalOauthRsaPublicKey2.Value, externalOauthRsaPublicKey2.Value, nil},
				describeMapping{"external_oauth_blocked_roles_list", "external_oauth_blocked_roles_list", externalOauthBlockedRolesList.Value, externalOauthBlockedRolesList.Value, nil},
				describeMapping{"external_oauth_allowed_roles_list", "external_oauth_allowed_roles_list", externalOauthAllowedRolesList.Value, externalOauthAllowedRolesList.Value, nil},
				describeMapping{"external_oauth_audience_list", "external_oauth_audience_list", externalOauthAudienceList.Value, externalOauthAudienceList.Value, nil},
				describeMapping{"external_oauth_token_user_mapping_claim", "external_oauth_token_user_mapping_claim", externalOauthTokenUserMappingClaim.Value, externalOauthTokenUserMappingClaim.Value, nil},
				describeMapping{"external_oauth_snowflake_user_mapping_attribute", "external_oauth_snowflake_user_mapping_attribute", externalOauthSnowflakeUserMappingAttribute.Value, externalOauthSnowflakeUserMappingAttribute.Value, nil},
				describeMapping{"external_oauth_scope_delimiter", "external_oauth_scope_delimiter", externalOauthScopeDelimiter.Value, externalOauthScopeDelimiter.Value, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}
		if !d.GetRawConfig().IsNull() {
			if v := d.GetRawConfig().AsValueMap()["external_oauth_jws_keys_url"]; !v.IsNull() {
				if err = d.Set("external_oauth_jws_keys_url", ctyValToSliceString(v.AsValueSlice())); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["external_oauth_rsa_public_key"]; !v.IsNull() {
				if err = d.Set("external_oauth_rsa_public_key", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["external_oauth_rsa_public_key_2"]; !v.IsNull() {
				if err = d.Set("external_oauth_rsa_public_key_2", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["external_oauth_blocked_roles_list"]; !v.IsNull() {
				if err = d.Set("external_oauth_blocked_roles_list", ctyValToSliceString(v.AsValueSlice())); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["external_oauth_allowed_roles_list"]; !v.IsNull() {
				if err = d.Set("external_oauth_allowed_roles_list", ctyValToSliceString(v.AsValueSlice())); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["external_oauth_audience_list"]; !v.IsNull() {
				if err = d.Set("external_oauth_audience_list", ctyValToSliceString(v.AsValueSlice())); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["external_oauth_any_role_mode"]; !v.IsNull() {
				if err = d.Set("external_oauth_any_role_mode", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["external_oauth_scope_delimiter"]; !v.IsNull() {
				if err = d.Set("external_oauth_scope_delimiter", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["external_oauth_scope_mapping_attribute"]; !v.IsNull() {
				if err = d.Set("external_oauth_scope_mapping_attribute", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["comment"]; !v.IsNull() {
				if err = d.Set("comment", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.SecurityIntegrationToSchema(integration)}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ExternalOauthSecurityIntegrationPropertiesToSchema(integrationProperties)}); err != nil {
			return diag.FromErr(err)
		}
		// param := "EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST"
		// params, err := client.Parameters.ShowParameters(ctx, &sdk.ShowParametersOptions{
		// 	Like: &sdk.Like{
		// 		Pattern: sdk.Pointer(param),
		// 	},
		// 	In: &sdk.ParametersIn{
		// 		Account: sdk.Pointer(true),
		// 	},
		// })
		// if err != nil {
		// 	return diag.FromErr(err)
		// }
		// var found *sdk.Parameter
		// for _, v := range params {
		// 	if v.Key == param {
		// 		found = v
		// 		break
		// 	}
		// }
		// if found == nil {
		// 	return diag.FromErr(fmt.Errorf("parameter %s not found", param))
		// }
		// paramVal := helpers.StringToBool(found.Value)
		// if !paramVal {
		// 	return nil
		// }
		// if err := d.Set("external_oauth_add_privileged_roles_to_blocked_list", paramVal); err != nil {
		// 	return diag.FromErr(err)
		// }
		return nil
	}
}

func UpdateContextExternalOauthIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	set, unset := sdk.NewExternalOauthIntegrationSetRequest(), sdk.NewExternalOauthIntegrationUnsetRequest()

	if d.HasChange("comment") {
		// TODO(this pr) add to sdk
		// set.WithComment(sdk.StringAllowEmpty{Value: d.Get("comment").(string)})
	}

	if d.HasChange("enabled") {
		// this field is required
		set.WithEnabled(d.Get("enabled").(bool))
	}

	if d.HasChange("external_oauth_allowed_roles_list") {
		v := d.Get("external_oauth_allowed_roles_list").([]any)
		allowedRoles := make([]sdk.AccountObjectIdentifier, len(v))
		for i := range v {
			allowedRoles[i] = sdk.NewAccountObjectIdentifier(v[i].(string))
		}
		set.WithExternalOauthAllowedRolesList(sdk.AllowedRolesListRequest{AllowedRolesList: allowedRoles})
	}

	if d.HasChange("external_oauth_any_role_mode") {
		v := d.Get("external_oauth_any_role_mode").(string)
		if len(v) > 0 {
			value, err := sdk.ToExternalOauthSecurityIntegrationAnyRoleModeOption(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithExternalOauthAnyRoleMode(value)
		} else {
			// TODO(SNOW-1515781): use UNSET
			set.WithExternalOauthAnyRoleMode(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)
		}
	}

	if d.HasChange("external_oauth_audience_list") {
		v := d.Get("external_oauth_audience_list").([]any)
		if len(v) > 0 {
			audienceList := make([]sdk.AudienceListItem, len(v))
			for i := range v {
				audienceList[i] = sdk.AudienceListItem{Item: v[i].(string)}
			}
			set.WithExternalOauthAudienceList(sdk.AudienceListRequest{AudienceList: audienceList})
		} else {
			unset.WithExternalOauthAudienceList(true)
		}
	}

	if d.HasChange("external_oauth_blocked_roles_list") {
		v := d.Get("external_oauth_blocked_roles_list").([]any)
		blockedRoles := make([]sdk.AccountObjectIdentifier, len(v))
		for i := range v {
			blockedRoles[i] = sdk.NewAccountObjectIdentifier(v[i].(string))
		}
		set.WithExternalOauthBlockedRolesList(sdk.BlockedRolesListRequest{BlockedRolesList: blockedRoles})
	}

	if d.HasChange("external_oauth_issuer") {
		// this field is required
		set.WithExternalOauthIssuer(d.Get("external_oauth_issuer").(string))
	}

	if d.HasChange("external_oauth_jws_keys_url") {
		v := d.Get("external_oauth_jws_keys_url").([]any)
		if len(v) > 0 {
			urls := make([]sdk.JwsKeysUrl, len(v))
			for i := range v {
				urls[i] = sdk.JwsKeysUrl{JwsKeyUrl: v[i].(string)}
			}
			set.WithExternalOauthJwsKeysUrl(urls)
		}
	}

	if d.HasChange("external_oauth_rsa_public_key") {
		set.WithExternalOauthRsaPublicKey(d.Get("external_oauth_rsa_public_key").(string))
		if v, ok := d.GetOk("external_oauth_rsa_public_key"); ok {
			set.WithExternalOauthRsaPublicKey2(v.(string))
		}
		// else: force new
	}

	if d.HasChange("external_oauth_rsa_public_key_2") {
		if v, ok := d.GetOk("external_oauth_rsa_public_key_2"); ok {
			set.WithExternalOauthRsaPublicKey2(v.(string))
		}
		// else: force new
	}

	if d.HasChange("external_oauth_scope_delimiter") {
		if v, ok := d.GetOk("external_oauth_scope_delimiter"); ok {
			set.WithExternalOauthScopeDelimiter(v.(string))
		} else {
			// TODO(SNOW-1515781): use UNSET
			set.WithExternalOauthScopeDelimiter(",")
		}
	}

	if d.HasChange("external_oauth_scope_mapping_attribute") {
		// TODO(this pr) add to sdk
		// set.WithExternalOauthScopeMappingAttribute(d.Get("external_oauth_scope_mapping_attribute").(string))
	}

	if d.HasChange("external_oauth_snowflake_user_mapping_attribute") {
		if v, ok := d.GetOk("external_oauth_snowflake_user_mapping_attribute"); ok {
			value, err := sdk.ToExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption(v.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithExternalOauthSnowflakeUserMappingAttribute(value)
		}
		// else: force new
	}

	if d.HasChange("external_oauth_token_user_mapping_claim") {
		v := d.Get("external_oauth_token_user_mapping_claim").([]any)
		claims := make([]sdk.TokenUserMappingClaim, len(v))
		for i := range v {
			claims[i] = sdk.TokenUserMappingClaim{
				Claim: v[i].(string),
			}
		}
		set.WithExternalOauthTokenUserMappingClaim(claims)
	}

	if d.HasChange("external_oauth_type") {
		if v, ok := d.GetOk("external_oauth_type"); ok {
			value, err := sdk.ToExternalOauthSecurityIntegrationTypeOption(v.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithExternalOauthType(value)
		}
	}

	if !reflect.DeepEqual(*set, sdk.ExternalOauthIntegrationSetRequest{}) {
		if err := client.SecurityIntegrations.AlterExternalOauth(ctx, sdk.NewAlterExternalOauthSecurityIntegrationRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}
	if !reflect.DeepEqual(*unset, sdk.ExternalOauthIntegrationUnsetRequest{}) {
		if err := client.SecurityIntegrations.AlterExternalOauth(ctx, sdk.NewAlterExternalOauthSecurityIntegrationRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}
	return ReadContextExternalOauthIntegration(false)(ctx, d, meta)
}

func DeleteContextExternalOauthIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
