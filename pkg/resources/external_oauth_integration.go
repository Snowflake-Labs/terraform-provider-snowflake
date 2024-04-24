package resources

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
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
	"type": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the OAuth 2.0 authorization server to be Okta, Microsoft Azure AD, Ping Identity PingFederate, or a Custom OAuth 2.0 authorization server.",
		ValidateFunc: validation.StringInSlice([]string{
			string(snowflake.Okta),
			string(snowflake.Azure),
			string(snowflake.PingFederate),
			string(snowflake.Custom),
		}, true),
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
	"issuer": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the URL to define the OAuth 2.0 authorization server.",
	},
	"token_user_mapping_claims": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: "Specifies the access token claim or claims that can be used to map the access token to a Snowflake user record.",
	},
	"scope_mapping_attribute": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the access token claim to map the access token to an account role.",
	},
	"snowflake_user_mapping_attribute": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Indicates which Snowflake user record attribute should be used to map the access token to a Snowflake user record.",
		ValidateFunc: validation.StringInSlice([]string{
			string(snowflake.LoginName),
			string(snowflake.EmailAddress),
		}, true),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalize := func(s string) string {
				return strings.ToUpper(strings.ReplaceAll(s, "-", ""))
			}
			return normalize(old) == normalize(new)
		},
	},
	"jws_keys_urls": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		MaxItems:    3,
		Optional:    true,
		Description: "Specifies the endpoint or a list of endpoints from which to download public keys or certificates to validate an External OAuth access token. The maximum number of URLs that can be specified in the list is 3.",
	},
	"rsa_public_key": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a Base64-encoded RSA public key, without the -----BEGIN PUBLIC KEY----- and -----END PUBLIC KEY----- headers.",
	},
	"rsa_public_key_2": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a second RSA public key, without the -----BEGIN PUBLIC KEY----- and -----END PUBLIC KEY----- headers. Used for key rotation.",
	},
	"blocked_roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies the list of roles that a client cannot set as the primary role. Do not include ACCOUNTADMIN, ORGADMIN or SECURITYADMIN as they are already implicitly enforced and will cause in-place updates.",
	},
	"allowed_roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies the list of roles that the client can set as the primary role.",
	},
	"audience_urls": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies additional values that can be used for the access token's audience validation on top of using the Customer's Snowflake Account URL ",
	},
	"any_role_mode": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     string(snowflake.Disable),
		Description: "Specifies whether the OAuth client or user can use a role that is not defined in the OAuth access token.",
		ValidateFunc: validation.StringInSlice([]string{
			string(snowflake.Disable),
			string(snowflake.Enable),
			string(snowflake.EnableForPrivilege),
		}, true),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalize := func(s string) string {
				return strings.ToUpper(strings.ReplaceAll(s, "-", ""))
			}
			return normalize(old) == normalize(new)
		},
	},
	"scope_delimiter": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the scope delimiter in the authorization token.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the OAuth integration.",
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the External OAUTH integration was created.",
	},
}

// ExternalOauthIntegration returns a pointer to the resource representing a network policy.
func ExternalOauthIntegration() *schema.Resource {
	return &schema.Resource{
		Description: "An External OAuth security integration allows a client to use a third-party authorization server to obtain the access tokens needed to interact with Snowflake.",
		Create:      CreateExternalOauthIntegration,
		Read:        ReadExternalOauthIntegration,
		Update:      UpdateExternalOauthIntegration,
		Delete:      DeleteExternalOauthIntegration,

		Schema: oauthExternalIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateExternalOauthIntegration implements schema.CreateFunc.
func CreateExternalOauthIntegration(d *schema.ResourceData, meta interface{}) error {
	manager, err := snowflake.NewExternalOauthIntegration3Manager()
	if err != nil {
		return fmt.Errorf("couldn't create external oauth integration manager: %w", err)
	}

	input := &snowflake.ExternalOauthIntegration3CreateInput{
		ExternalOauthIntegration3: snowflake.ExternalOauthIntegration3{
			TopLevelIdentifier: snowflake.TopLevelIdentifier{
				Name: d.Get("name").(string),
			},

			Type:                                 "EXTERNAL_OAUTH",
			TypeOk:                               true,
			Enabled:                              d.Get("enabled").(bool),
			EnabledOk:                            isOk(d.GetOk("enabled")),
			ExternalOauthType:                    snowflake.ExternalOauthType(d.Get("type").(string)),
			ExternalOauthTypeOk:                  isOk(d.GetOk("type")),
			ExternalOauthIssuer:                  d.Get("issuer").(string),
			ExternalOauthIssuerOk:                isOk(d.GetOk("issuer")),
			ExternalOauthTokenUserMappingClaim:   expandStringList(d.Get("token_user_mapping_claims").(*schema.Set).List()),
			ExternalOauthTokenUserMappingClaimOk: isOk(d.GetOk("token_user_mapping_claims")),
			ExternalOauthSnowflakeUserMappingAttribute:   snowflake.SFUserMappingAttribute(d.Get("snowflake_user_mapping_attribute").(string)),
			ExternalOauthSnowflakeUserMappingAttributeOk: isOk(d.GetOk("snowflake_user_mapping_attribute")),
			ExternalOauthJwsKeysURL:                      expandStringList(d.Get("jws_keys_urls").(*schema.Set).List()),
			ExternalOauthJwsKeysURLOk:                    isOk(d.GetOk("jws_keys_urls")),
			ExternalOauthBlockedRolesList:                expandStringList(d.Get("blocked_roles").(*schema.Set).List()),
			ExternalOauthBlockedRolesListOk:              isOk(d.GetOk("blocked_roles")),
			ExternalOauthAllowedRolesList:                expandStringList(d.Get("allowed_roles").(*schema.Set).List()),
			ExternalOauthAllowedRolesListOk:              isOk(d.GetOk("allowed_roles")),
			ExternalOauthRsaPublicKey:                    d.Get("rsa_public_key").(string),
			ExternalOauthRsaPublicKeyOk:                  isOk(d.GetOk("rsa_public_key")),
			ExternalOauthRsaPublicKey2:                   d.Get("rsa_public_key_2").(string),
			ExternalOauthRsaPublicKey2Ok:                 isOk(d.GetOk("rsa_public_key_2")),
			ExternalOauthAudienceList:                    expandStringList(d.Get("audience_urls").(*schema.Set).List()),
			ExternalOauthAudienceListOk:                  isOk(d.GetOk("audience_urls")),
			ExternalOauthAnyRoleMode:                     snowflake.AnyRoleMode(d.Get("any_role_mode").(string)),
			ExternalOauthAnyRoleModeOk:                   isOk(d.GetOk("any_role_mode")),
			ExternalOauthScopeDelimiter:                  d.Get("scope_delimiter").(string),
			ExternalOauthScopeDelimiterOk:                isOk(d.GetOk("scope_delimiter")),
			ExternalOauthScopeMappingAttribute:           d.Get("scope_mapping_attribute").(string),
			ExternalOauthScopeMappingAttributeOk:         isOk(d.GetOk("scope_mapping_attribute")),

			Comment:   sql.NullString{String: d.Get("comment").(string)},
			CommentOk: isOk(d.GetOk("comment")),
		},
	}

	stmt, err := manager.Create(input)
	if err != nil {
		return fmt.Errorf("couldn't generate create statement: %w", err)
	}

	client := meta.(*provider.Context).Client
	db := client.GetConn().DB
	err = snowflake.Exec(db, stmt)
	if err != nil {
		return fmt.Errorf("error executing create statement: %w", err)
	}

	d.SetId(ExternalOauthIntegrationID(&input.ExternalOauthIntegration3))

	return ReadExternalOauthIntegration(d, meta)
}

// ReadExternalOauthIntegration implements schema.ReadFunc.
func ReadExternalOauthIntegration(d *schema.ResourceData, meta interface{}) error {
	manager, err := snowflake.NewExternalOauthIntegration3Manager()
	if err != nil {
		return fmt.Errorf("couldn't create external oauth integration builder: %w", err)
	}

	input := ExternalOauthIntegrationIdentifier(d.Id())

	client := meta.(*provider.Context).Client
	db := client.GetConn().DB

	// This resource needs a SHOW and a DESCRIBE

	// SHOW
	stmt, err := manager.ReadShow(input)
	if err != nil {
		return fmt.Errorf("couldn't generate show statement: %w", err)
	}

	row := snowflake.QueryRow(db, stmt)
	if err != nil {
		return fmt.Errorf("error querying external oauth integration: %w", err)
	}

	showOutput, err := manager.ParseShow(row)
	if err != nil {
		return fmt.Errorf("error parsing show result: %w", err)
	}

	if err := d.Set("type", strings.TrimPrefix(showOutput.Type, "EXTERNAL_OAUTH - ")); err != nil {
		return fmt.Errorf("error setting type: %w", err)
	}
	if err := d.Set("name", showOutput.Name); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}
	if err := d.Set("enabled", showOutput.Enabled); err != nil {
		return fmt.Errorf("error setting enabled: %w", err)
	}
	if err := d.Set("comment", showOutput.Comment.String); err != nil {
		return fmt.Errorf("error setting comment: %w", err)
	}
	// if err := d.Set("created_on", showOutput.CreatedOn.String); err != nil {
	// 	return fmt.Errorf("error setting created_on: %w", err)
	// }

	// DESCRIBE
	stmt, err = manager.ReadDescribe(input)
	if err != nil {
		return fmt.Errorf("couldn't generate describe statement: %w", err)
	}

	rows, err := snowflake.Query(db, stmt)
	if err != nil {
		return fmt.Errorf("error querying external oauth integration: %w", err)
	}

	defer rows.Close()
	describeOutput, err := manager.ParseDescribe(rows.Rows)
	if err != nil {
		return fmt.Errorf("failed to parse result of describe: %w", err)
	}

	if err := d.Set("issuer", describeOutput.ExternalOauthIssuer); err != nil {
		return fmt.Errorf("error setting issuer: %w", err)
	}
	if err := d.Set("jws_keys_urls", describeOutput.ExternalOauthJwsKeysURL); err != nil {
		return fmt.Errorf("error setting jws_keys_urls: %w", err)
	}
	if err := d.Set("any_role_mode", describeOutput.ExternalOauthAnyRoleMode); err != nil {
		return fmt.Errorf("error setting any_role_mode: %w", err)
	}
	if err := d.Set("rsa_public_key", describeOutput.ExternalOauthRsaPublicKey); err != nil {
		return fmt.Errorf("error setting rsa_public_key: %w", err)
	}
	if err := d.Set("rsa_public_key_2", describeOutput.ExternalOauthRsaPublicKey2); err != nil {
		return fmt.Errorf("error setting rsa_public_key_2: %w", err)
	}
	// Filter out default roles
	blockedRoles := []string{}
	for i := range describeOutput.ExternalOauthBlockedRolesList {
		role := describeOutput.ExternalOauthBlockedRolesList[i]
		if role != "ACCOUNTADMIN" && role != "SECURITYADMIN" {
			blockedRoles = append(blockedRoles, role)
		}
	}
	if err := d.Set("blocked_roles", blockedRoles); err != nil {
		return fmt.Errorf("error setting blocked_roles: %w", err)
	}
	if err := d.Set("allowed_roles", describeOutput.ExternalOauthAllowedRolesList); err != nil {
		return fmt.Errorf("error setting allowed_roles: %w", err)
	}
	if err := d.Set("audience_urls", describeOutput.ExternalOauthAudienceList); err != nil {
		return fmt.Errorf("error setting audience_urls: %w", err)
	}
	if err := d.Set("token_user_mapping_claims", describeOutput.ExternalOauthTokenUserMappingClaim); err != nil {
		return fmt.Errorf("error setting token_user_mapping_claims: %w", err)
	}
	if err := d.Set("snowflake_user_mapping_attribute", describeOutput.ExternalOauthSnowflakeUserMappingAttribute); err != nil {
		return fmt.Errorf("error setting snowflake_user_mapping_attribute: %w", err)
	}
	if err := d.Set("scope_mapping_attribute", describeOutput.ExternalOauthScopeMappingAttribute); err != nil {
		return fmt.Errorf("error setting scope_mapping_attribute: %w", err)
	}

	return err
}

// UpdateExternalOauthIntegration implements schema.UpdateFunc.
func UpdateExternalOauthIntegration(d *schema.ResourceData, meta interface{}) error {
	manager, err := snowflake.NewExternalOauthIntegration3Manager()
	if err != nil {
		return fmt.Errorf("couldn't create external oauth integration builder: %w", err)
	}

	runAlter := false
	alterInput := &snowflake.ExternalOauthIntegration3UpdateInput{
		ExternalOauthIntegration3: snowflake.ExternalOauthIntegration3{
			TopLevelIdentifier: snowflake.TopLevelIdentifier{
				Name: d.Get("name").(string),
			},
		},
	}
	runUnset := false
	unsetInput := &snowflake.ExternalOauthIntegration3UpdateInput{
		ExternalOauthIntegration3: snowflake.ExternalOauthIntegration3{
			TopLevelIdentifier: snowflake.TopLevelIdentifier{
				Name: d.Get("name").(string),
			},
		},
	}

	if d.HasChange("enabled") {
		val, ok := d.GetOk("enabled")
		if ok {
			alterInput.Enabled = val.(bool)
			alterInput.EnabledOk = true
			runAlter = true
		} else {
			unsetInput.EnabledOk = true
			runUnset = true
		}
	}
	if d.HasChange("type") {
		val, ok := d.GetOk("type")
		if ok {
			alterInput.Type = val.(string)
			alterInput.TypeOk = true
			runAlter = true
		} else {
			unsetInput.TypeOk = true
			runUnset = true
		}
	}
	if d.HasChange("issuer") {
		val, ok := d.GetOk("issuer")
		if ok {
			alterInput.ExternalOauthIssuer = val.(string)
			alterInput.ExternalOauthIssuerOk = true
			runAlter = true
		} else {
			unsetInput.ExternalOauthIssuerOk = true
			runUnset = true
		}
	}
	if d.HasChange("token_user_mapping_claims") {
		val, ok := d.GetOk("token_user_mapping_claims")
		if ok {
			alterInput.ExternalOauthTokenUserMappingClaim = expandStringList(val.(*schema.Set).List())
			alterInput.ExternalOauthTokenUserMappingClaimOk = true
			runAlter = true
		} else {
			unsetInput.ExternalOauthTokenUserMappingClaimOk = true
			runUnset = true
		}
	}
	if d.HasChange("snowflake_user_mapping_attribute") {
		val, ok := d.GetOk("snowflake_user_mapping_attribute")
		if ok {
			alterInput.ExternalOauthSnowflakeUserMappingAttribute = snowflake.SFUserMappingAttribute(val.(string))
			alterInput.ExternalOauthSnowflakeUserMappingAttributeOk = true
			runAlter = true
		} else {
			unsetInput.ExternalOauthSnowflakeUserMappingAttributeOk = true
			runUnset = true
		}
	}
	if d.HasChange("jws_keys_urls") {
		val, ok := d.GetOk("jws_keys_urls")
		if ok {
			alterInput.ExternalOauthJwsKeysURL = expandStringList(val.(*schema.Set).List())
			alterInput.ExternalOauthJwsKeysURLOk = true
			runAlter = true
		} else {
			unsetInput.ExternalOauthJwsKeysURLOk = true
			runUnset = true
		}
	}
	if d.HasChange("rsa_public_key") {
		val, ok := d.GetOk("rsa_public_key")
		if ok {
			alterInput.ExternalOauthRsaPublicKey = val.(string)
			alterInput.ExternalOauthRsaPublicKeyOk = true
			runAlter = true
		} else {
			unsetInput.ExternalOauthRsaPublicKeyOk = true
			runUnset = true
		}
	}
	if d.HasChange("rsa_public_key_2") {
		val, ok := d.GetOk("rsa_public_key_2")
		if ok {
			alterInput.ExternalOauthRsaPublicKey2 = val.(string)
			alterInput.ExternalOauthRsaPublicKey2Ok = true
			runAlter = true
		} else {
			unsetInput.ExternalOauthRsaPublicKey2Ok = true
			runUnset = true
		}
	}
	if d.HasChange("blocked_roles") {
		val, ok := d.GetOk("blocked_roles")
		if ok {
			alterInput.ExternalOauthBlockedRolesList = expandStringList(val.(*schema.Set).List())
			alterInput.ExternalOauthBlockedRolesListOk = true
			runAlter = true
		} else {
			unsetInput.ExternalOauthBlockedRolesListOk = true
			runUnset = true
		}
	}
	if d.HasChange("allowed_roles") {
		val, ok := d.GetOk("allowed_roles")
		if ok {
			alterInput.ExternalOauthAllowedRolesList = expandStringList(val.(*schema.Set).List())
			alterInput.ExternalOauthAllowedRolesListOk = true
			runAlter = true
		} else {
			unsetInput.ExternalOauthAllowedRolesListOk = true
			runUnset = true
		}
	}
	if d.HasChange("audience_urls") {
		val, ok := d.GetOk("audience_urls")
		if ok {
			alterInput.ExternalOauthAudienceList = expandStringList(val.(*schema.Set).List())
			alterInput.ExternalOauthAudienceListOk = true
			runAlter = true
		} else {
			unsetInput.ExternalOauthAudienceListOk = true
			runUnset = true
		}
	}
	if d.HasChange("any_role_mode") {
		val, ok := d.GetOk("any_role_mode")
		if ok {
			alterInput.ExternalOauthAnyRoleMode = snowflake.AnyRoleMode(val.(string))
			alterInput.ExternalOauthAnyRoleModeOk = true
			runAlter = true
		} else {
			unsetInput.ExternalOauthAnyRoleModeOk = true
			runUnset = true
		}
	}
	if d.HasChange("scope_delimiter") {
		val, ok := d.GetOk("scope_delimiter")
		if ok {
			alterInput.ExternalOauthScopeDelimiter = val.(string)
			alterInput.ExternalOauthScopeDelimiterOk = true
			runAlter = true
		} else {
			unsetInput.ExternalOauthScopeDelimiterOk = true
			runUnset = true
		}
	}
	if d.HasChange("scope_mapping_attribute") {
		val, ok := d.GetOk("scope_mapping_attribute")
		if ok {
			alterInput.ExternalOauthScopeMappingAttribute = val.(string)
			alterInput.ExternalOauthScopeMappingAttributeOk = true
			runAlter = true
		} else {
			unsetInput.ExternalOauthScopeMappingAttributeOk = true
			runUnset = true
		}
	}
	if d.HasChange("comment") {
		val, ok := d.GetOk("comment")
		if ok {
			alterInput.Comment.String = val.(string)
			alterInput.CommentOk = true
			runAlter = true
		} else {
			unsetInput.CommentOk = true
			runUnset = true
		}
	}

	client := meta.(*provider.Context).Client
	db := client.GetConn().DB

	if runAlter {
		stmt, err := manager.Update(alterInput)
		if err != nil {
			return fmt.Errorf("couldn't generate alter statement for external oauth integration: %w", err)
		}

		err = snowflake.Exec(db, stmt)
		if err != nil {
			return fmt.Errorf("error executing alter statement: %w", err)
		}
	}

	if runUnset {
		stmt, err := manager.Unset(unsetInput)
		if err != nil {
			return fmt.Errorf("couldn't generate unset statement for external oauth integration: %w", err)
		}

		err = snowflake.Exec(db, stmt)
		if err != nil {
			return fmt.Errorf("error executing unset statement: %w", err)
		}
	}

	return ReadExternalOauthIntegration(d, meta)
}

// DeleteExternalOauthIntegration implements schema.DeleteFunc.
func DeleteExternalOauthIntegration(d *schema.ResourceData, meta interface{}) error {
	manager, err := snowflake.NewExternalOauthIntegration3Manager()
	if err != nil {
		return fmt.Errorf("couldn't create external oauth integration builder: %w", err)
	}

	input := &snowflake.ExternalOauthIntegration3DeleteInput{
		TopLevelIdentifier: snowflake.TopLevelIdentifier{
			Name: d.Get("name").(string),
		},
	}

	stmt, err := manager.Delete(input)
	if err != nil {
		return fmt.Errorf("couldn't generate drop statement: %w", err)
	}

	client := meta.(*provider.Context).Client
	db := client.GetConn().DB
	err = snowflake.Exec(db, stmt)
	if err != nil {
		return fmt.Errorf("error executing drop statement: %w", err)
	}

	return nil
}

func ExternalOauthIntegrationID(eoi *snowflake.ExternalOauthIntegration3) string {
	return eoi.QualifiedName()
}

func ExternalOauthIntegrationIdentifier(id string) *snowflake.TopLevelIdentifier {
	return snowflake.TopLevelIdentifierFromQualifiedName(id)
}
