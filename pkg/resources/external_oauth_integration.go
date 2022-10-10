package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pkg/errors"
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
			"OKTA", "AZURE", "PING_FEDERATE", "CUSTOM",
		}, true),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalize := func(s string) string {
				return strings.ToUpper(strings.Replace(s, "-", "", -1))
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
	"snowflake_user_mapping_attribute": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Indicates which Snowflake user record attribute should be used to map the access token to a Snowflake user record.",
		ValidateFunc: validation.StringInSlice([]string{
			"LOGIN_NAME", "EMAIL_ADDRESS",
		}, true),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalize := func(s string) string {
				return strings.ToUpper(strings.Replace(s, "-", "", -1))
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
		Default:     "DISABLE",
		Description: "Specifies whether the OAuth client or user can use a role that is not defined in the OAuth access token.",
		ValidateFunc: validation.StringInSlice([]string{
			"DISABLE", "ENABLE", "ENABLE_FOR_PRIVILEGE",
		}, true),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalize := func(s string) string {
				return strings.ToUpper(strings.Replace(s, "-", "", -1))
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
		Create: CreateExternalOauthIntegration,
		Read:   ReadExternalOauthIntegration,
		Update: UpdateExternalOauthIntegration,
		Delete: DeleteExternalOauthIntegration,

		Schema: oauthExternalIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateExternalOauthIntegration implements schema.CreateFunc.
func CreateExternalOauthIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)

	stmt := snowflake.ExternalOauthIntegration(name).Create()

	// Set required fields
	stmt.SetRaw(`TYPE=EXTERNAL_OAUTH`)
	stmt.SetBool(`ENABLED`, d.Get("enabled").(bool))
	stmt.SetString(`EXTERNAL_OAUTH_TYPE`, d.Get("type").(string))
	stmt.SetString(`EXTERNAL_OAUTH_ISSUER`, d.Get("issuer").(string))
	stmt.SetStringList(`EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM`, expandStringList(d.Get("token_user_mapping_claims").(*schema.Set).List()))
	stmt.SetString(`EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE`, d.Get("snowflake_user_mapping_attribute").(string))

	// Set optional fields
	if _, ok := d.GetOk("jws_keys_urls"); ok {
		stmt.SetStringList(`EXTERNAL_OAUTH_JWS_KEYS_URL`, expandStringList(d.Get("jws_keys_urls").(*schema.Set).List()))
	}
	if _, ok := d.GetOk("rsa_public_key"); ok {
		stmt.SetString(`EXTERNAL_OAUTH_RSA_PUBLIC_KEY`, d.Get("rsa_public_key").(string))
	}
	if _, ok := d.GetOk("rsa_public_key_2"); ok {
		stmt.SetString(`EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2`, d.Get("rsa_public_key_2").(string))
	}
	if _, ok := d.GetOk("blocked_roles"); ok {
		stmt.SetStringList(`EXTERNAL_OAUTH_BLOCKED_ROLES_LIST`, expandStringList(d.Get("blocked_roles").(*schema.Set).List()))
	}
	if _, ok := d.GetOk("allowed_roles"); ok {
		stmt.SetStringList(`EXTERNAL_OAUTH_ALLOWED_ROLES_LIST`, expandStringList(d.Get("allowed_roles").(*schema.Set).List()))
	}
	if _, ok := d.GetOk("audience_urls"); ok {
		stmt.SetStringList(`EXTERNAL_OAUTH_AUDIENCE_LIST`, expandStringList(d.Get("audience_urls").(*schema.Set).List()))
	}
	if _, ok := d.GetOk("any_role_mode"); ok {
		stmt.SetString(`EXTERNAL_OAUTH_ANY_ROLE_MODE`, d.Get("any_role_mode").(string))
	}
	if _, ok := d.GetOk("scope_delimiter"); ok {
		stmt.SetString(`EXTERNAL_OAUTH_SCOPE_DELIMITER`, d.Get("scope_delimiter").(string))
	}
	if _, ok := d.GetOk("comment"); ok {
		stmt.SetString(`COMMENT`, d.Get("comment").(string))
	}

	err := snowflake.Exec(db, stmt.Statement())
	if err != nil {
		return errors.Wrap(err, "error creating security integration"+stmt.Statement())
	}

	d.SetId(name)

	return ReadExternalOauthIntegration(d, meta)
}

// ReadExternalOauthIntegration implements schema.ReadFunc.
func ReadExternalOauthIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	stmt := snowflake.ExternalOauthIntegration(id).Show()
	row := snowflake.QueryRow(db, stmt)

	// Some properties can come from the SHOW INTEGRATION call

	s, err := snowflake.ScanExternalOauthIntegration(row)
	if err != nil {
		return errors.Wrap(err, "could not show security integration")
	}

	// Note: category must be Security or something is broken
	if c := s.Category.String; c != "SECURITY" {
		return fmt.Errorf("expected %v to be an Security integration, got %v", id, c)
	}

	if err := d.Set("type", strings.TrimPrefix(s.IntegrationType.String, "EXTERNAL_OAUTH - ")); err != nil {
		return err
	}

	if err := d.Set("name", s.Name.String); err != nil {
		return err
	}

	if err := d.Set("enabled", s.Enabled.Bool); err != nil {
		return err
	}

	if err := d.Set("comment", s.Comment.String); err != nil {
		return err
	}

	if err := d.Set("created_on", s.CreatedOn.String); err != nil {
		return err
	}

	// Some properties come from the DESCRIBE INTEGRATION call
	// We need to grab them in a loop
	var k, pType string
	var v, unused interface{}
	stmt = snowflake.ExternalOauthIntegration(id).Describe()
	rows, err := db.Query(stmt)
	if err != nil {
		return errors.Wrap(err, "could not describe security integration")
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&k, &pType, &v, &unused); err != nil {
			return errors.Wrap(err, "unable to parse security integration rows")
		}
		switch k {
		case "ENABLED":
			// We set this using the SHOW INTEGRATION call so let's ignore it here
		case "COMMENT":
			// We set this using the SHOW INTEGRATION call so let's ignore it here
		case "EXTERNAL_OAUTH_ISSUER":
			if err = d.Set("issuer", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set issuer for security integration")
			}
		case "EXTERNAL_OAUTH_JWS_KEYS_URL":
			list := []string{}
			list = append(list, strings.Split(v.(string), ",")...)
			if err = d.Set("jws_keys_urls", list); err != nil {
				return errors.Wrap(err, "unable to set jws keys urls for security integration")
			}
		case "EXTERNAL_OAUTH_ANY_ROLE_MODE":
			if err = d.Set("any_role_mode", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set any role mode for security integration")
			}
		case "EXTERNAL_OAUTH_RSA_PUBLIC_KEY":
			if err = d.Set("rsa_public_key", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set rsa public key for security integration")
			}
		case "EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2":
			if err = d.Set("rsa_public_key_2", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set rsa public key 2 for security integration")
			}
		case "EXTERNAL_OAUTH_BLOCKED_ROLES_LIST":
			blockedRolesAll := strings.Split(v.(string), ",")
			// Only roles other than ACCOUNTADMIN, ORGADMIN and SECURITYADMIN can be specified custom,
			// those three are enforced with no option to remove them
			blockedRolesCustom := []string{}
			for _, role := range blockedRolesAll {
				if role != "ACCOUNTADMIN" && role != "ORGADMIN" && role != "SECURITYADMIN" && role != "" {
					blockedRolesCustom = append(blockedRolesCustom, role)
				}
			}

			if err = d.Set("blocked_roles", blockedRolesCustom); err != nil {
				return errors.Wrap(err, "unable to set blocked roles for security integration")
			}
		case "EXTERNAL_OAUTH_ALLOWED_ROLES_LIST":
			list := []string{}
			for _, item := range strings.Split(v.(string), ",") {
				if item != "" {
					list = append(list, item)
				}
			}
			if err = d.Set("allowed_roles", list); err != nil {
				return errors.Wrap(err, "unable to set allowed roles for security integration")
			}
		case "EXTERNAL_OAUTH_AUDIENCE_LIST":
			list := []string{}
			for _, item := range strings.Split(v.(string), ",") {
				if item != "" {
					list = append(list, item)
				}
			}
			if err = d.Set("audience_urls", list); err != nil {
				return errors.Wrap(err, "unable to set audience urls for security integration")
			}
		case "EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM":
			list := []string{}
			for _, item := range strings.Split(strings.Replace(strings.Replace(v.(string), "[", "", 1), "]", "", 1), ",") {
				if item != "" {
					list = append(list, strings.Replace(item, "'", "", 2))
				}
			}
			if err = d.Set("token_user_mapping_claims", list); err != nil {
				return errors.Wrap(err, "unable to set token user mapping claims for security integration")
			}
		case "EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE":
			if err = d.Set("snowflake_user_mapping_attribute", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set snowflake mapping attribute for security integration")
			}
		default:
			log.Printf("[WARN] unexpected security integration property %v returned from Snowflake", k)
		}
	}

	return err
}

// UpdateExternalOauthIntegration implements schema.UpdateFunc.
func UpdateExternalOauthIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	stmt := snowflake.ExternalOauthIntegration(id).Alter()

	var runSetStatement bool

	if d.HasChange("enabled") {
		runSetStatement = true
		stmt.SetBool(`ENABLED`, d.Get("enabled").(bool))
	}
	if d.HasChange("type") {
		runSetStatement = true
		stmt.SetString(`EXTERNAL_OAUTH_TYPE`, d.Get("type").(string))
	}
	if d.HasChange("issuer") {
		runSetStatement = true
		stmt.SetString(`EXTERNAL_OAUTH_ISSUER`, d.Get("issuer").(string))
	}
	if d.HasChange("token_user_mapping_claims") {
		runSetStatement = true
		stmt.SetStringList(`EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM`, expandStringList(d.Get("token_user_mapping_claims").(*schema.Set).List()))
	}
	if d.HasChange("snowflake_user_mapping_attribute") {
		runSetStatement = true
		stmt.SetString(`EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE`, d.Get("snowflake_user_mapping_attribute").(string))
	}
	if d.HasChange("jws_keys_urls") {
		runSetStatement = true
		stmt.SetStringList(`EXTERNAL_OAUTH_JWS_KEYS_URL`, expandStringList(d.Get("jws_keys_urls").(*schema.Set).List()))
	}
	if d.HasChange("rsa_public_key") {
		runSetStatement = true
		stmt.SetString(`EXTERNAL_OAUTH_RSA_PUBLIC_KEY`, d.Get("rsa_public_key").(string))
	}
	if d.HasChange("rsa_public_key_2") {
		runSetStatement = true
		stmt.SetString(`EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2`, d.Get("rsa_public_key_2").(string))
	}
	if d.HasChange("blocked_roles") {
		runSetStatement = true
		stmt.SetStringList(`EXTERNAL_OAUTH_BLOCKED_ROLES_LIST`, expandStringList(d.Get("blocked_roles").(*schema.Set).List()))
	}
	if d.HasChange("allowed_roles") {
		runSetStatement = true
		stmt.SetStringList(`EXTERNAL_OAUTH_ALLOWED_ROLES_LIST`, expandStringList(d.Get("allowed_roles").(*schema.Set).List()))
	}
	if d.HasChange("audience_urls") {
		runSetStatement = true
		stmt.SetStringList(`EXTERNAL_OAUTH_AUDIENCE_LIST`, expandStringList(d.Get("audience_urls").(*schema.Set).List()))
	}
	if d.HasChange("any_role_mode") {
		runSetStatement = true
		stmt.SetString(`EXTERNAL_OAUTH_ANY_ROLE_MODE`, d.Get("any_role_mode").(string))
	}
	if d.HasChange("scope_delimiter") {
		runSetStatement = true
		stmt.SetString(`EXTERNAL_OAUTH_SCOPE_DELIMITER`, d.Get("scope_delimiter").(string))
	}
	if d.HasChange("comment") {
		runSetStatement = true
		stmt.SetString(`COMMENT`, d.Get("comment").(string))
	}

	if runSetStatement {
		if err := snowflake.Exec(db, stmt.Statement()); err != nil {
			return errors.Wrap(err, "error updating security integration")
		}
	}

	return ReadExternalOauthIntegration(d, meta)
}

// DeleteExternalOauthIntegration implements schema.DeleteFunc.
func DeleteExternalOauthIntegration(d *schema.ResourceData, meta interface{}) error {
	return DeleteResource("", snowflake.ExternalOauthIntegration)(d, meta)
}
