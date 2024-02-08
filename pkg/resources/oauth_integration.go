package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var oauthIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the name of the OAuth integration. This name follows the rules for Object Identifiers. The name should be unique among security integrations in your account.",
	},
	"oauth_client": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the OAuth client type.",
		ValidateFunc: validation.StringInSlice([]string{
			"TABLEAU_DESKTOP", "TABLEAU_SERVER", "LOOKER", "CUSTOM",
		}, false),
	},
	"oauth_redirect_uri": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the client URI. After a user is authenticated, the web browser is redirected to this URI.",
	},
	"oauth_client_type": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the type of client being registered. Snowflake supports both confidential and public clients.",
		ValidateFunc: validation.StringInSlice([]string{
			"CONFIDENTIAL", "PUBLIC",
		}, false),
	},
	"oauth_issue_refresh_tokens": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether to allow the client to exchange a refresh token for an access token when the current access token has expired.",
	},
	"oauth_refresh_token_validity": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Specifies how long refresh tokens should be valid (in seconds). OAUTH_ISSUE_REFRESH_TOKENS must be set to TRUE.",
	},
	"oauth_use_secondary_roles": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "NONE",
		Description: "Specifies whether default secondary roles set in the user properties are activated by default in the session being opened.",
		ValidateFunc: validation.StringInSlice([]string{
			"IMPLICIT", "NONE",
		}, false),
	},
	"blocked_roles_list": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "List of roles that a user cannot explicitly consent to using after authenticating. Do not include ACCOUNTADMIN, ORGADMIN or SECURITYADMIN as they are already implicitly enforced and will cause in-place updates.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the OAuth integration.",
	},
	"enabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether this OAuth integration is enabled or disabled.",
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the OAuth integration was created.",
	},
}

// OAuthIntegration returns a pointer to the resource representing an OAuth integration.
func OAuthIntegration() *schema.Resource {
	return &schema.Resource{
		Create: CreateOAuthIntegration,
		Read:   ReadOAuthIntegration,
		Update: UpdateOAuthIntegration,
		Delete: DeleteOAuthIntegration,

		Schema: oauthIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateOAuthIntegration implements schema.CreateFunc.
func CreateOAuthIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)

	stmt := snowflake.NewOAuthIntegrationBuilder(name).Create()

	// Set required fields
	stmt.SetRaw(`TYPE=OAUTH`)
	stmt.SetString(`OAUTH_CLIENT`, d.Get("oauth_client").(string))
	// Set optional fields
	if _, ok := d.GetOk("oauth_redirect_uri"); ok {
		stmt.SetString(`OAUTH_REDIRECT_URI`, d.Get("oauth_redirect_uri").(string))
	}
	if _, ok := d.GetOk("oauth_client_type"); ok {
		stmt.SetString(`OAUTH_CLIENT_TYPE`, d.Get("oauth_client_type").(string))
	}
	if _, ok := d.GetOk("oauth_issue_refresh_tokens"); ok {
		stmt.SetBool(`OAUTH_ISSUE_REFRESH_TOKENS`, d.Get("oauth_issue_refresh_tokens").(bool))
	}
	if _, ok := d.GetOk("oauth_refresh_token_validity"); ok {
		stmt.SetInt(`OAUTH_REFRESH_TOKEN_VALIDITY`, d.Get("oauth_refresh_token_validity").(int))
	}
	if _, ok := d.GetOk("oauth_use_secondary_roles"); ok {
		stmt.SetString(`OAUTH_USE_SECONDARY_ROLES`, d.Get("oauth_use_secondary_roles").(string))
	}
	if _, ok := d.GetOk("blocked_roles_list"); ok {
		stmt.SetStringList(`BLOCKED_ROLES_LIST`, expandStringList(d.Get("blocked_roles_list").(*schema.Set).List()))
	}
	if _, ok := d.GetOk("enabled"); ok {
		stmt.SetBool(`ENABLED`, d.Get("enabled").(bool))
	}
	if _, ok := d.GetOk("comment"); ok {
		stmt.SetString(`COMMENT`, d.Get("comment").(string))
	}

	if err := snowflake.Exec(db, stmt.Statement()); err != nil {
		return fmt.Errorf("error creating security integration err = %w", err)
	}

	d.SetId(name)

	return ReadOAuthIntegration(d, meta)
}

// ReadOAuthIntegration implements schema.ReadFunc.
func ReadOAuthIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	stmt := snowflake.NewOAuthIntegrationBuilder(id).Show()
	row := snowflake.QueryRow(db, stmt)

	// Some properties can come from the SHOW INTEGRATION call

	s, err := snowflake.ScanOAuthIntegration(row)
	if err != nil {
		return fmt.Errorf("could not show security integration err = %w", err)
	}

	// Note: category must be Security or something is broken
	if c := s.Category.String; c != "SECURITY" {
		return fmt.Errorf("expected %v to be an Security integration, got %v err = %w", id, c, err)
	}

	if err := d.Set("oauth_client", strings.TrimPrefix(s.IntegrationType.String, "OAUTH - ")); err != nil {
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
	stmt = snowflake.NewOAuthIntegrationBuilder(id).Describe()
	rows, err := db.Query(stmt)
	if err != nil {
		return fmt.Errorf("could not describe security integration err = %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&k, &pType, &v, &unused); err != nil {
			return fmt.Errorf("unable to parse security integration rows err = %w", err)
		}
		switch k {
		case "ENABLED":
			// We set this using the SHOW INTEGRATION call so let's ignore it here
		case "COMMENT":
			// We set this using the SHOW INTEGRATION call so let's ignore it here
		case "OAUTH_ISSUE_REFRESH_TOKENS":
			b, err := strconv.ParseBool(v.(string))
			if err != nil {
				return fmt.Errorf("returned OAuth issue refresh tokens that is not boolean err = %w", err)
			}
			if err := d.Set("oauth_issue_refresh_tokens", b); err != nil {
				return fmt.Errorf("unable to set OAuth issue refresh tokens for security integration err = %w", err)
			}
		case "OAUTH_REFRESH_TOKEN_VALIDITY":
			i, err := strconv.Atoi(v.(string))
			if err != nil {
				return fmt.Errorf("returned OAuth refresh token validity that is not integer err = %w", err)
			}
			if err := d.Set("oauth_refresh_token_validity", i); err != nil {
				return fmt.Errorf("unable to set OAuth refresh token validity for security integration err = %w", err)
			}
		case "OAUTH_USE_SECONDARY_ROLES":
			if err := d.Set("oauth_use_secondary_roles", v.(string)); err != nil {
				return fmt.Errorf("unable to set OAuth use secondary roles for security integration err = %w", err)
			}
		case "BLOCKED_ROLES_LIST":
			blockedRolesAll := strings.Split(v.(string), ",")

			// Only roles other than ACCOUNTADMIN, ORGADMIN and SECURITYADMIN can be specified custom,
			// those three are enforced with no option to remove them
			blockedRolesCustom := []string{}
			for _, role := range blockedRolesAll {
				if role != "ACCOUNTADMIN" && role != "ORGADMIN" && role != "SECURITYADMIN" {
					blockedRolesCustom = append(blockedRolesCustom, role)
				}
			}

			if err := d.Set("blocked_roles_list", blockedRolesCustom); err != nil {
				return fmt.Errorf("unable to set blocked roles list for security integration err = %w", err)
			}
		case "OAUTH_REDIRECT_URI":
			if err := d.Set("oauth_redirect_uri", v.(string)); err != nil {
				return fmt.Errorf("unable to set OAuth redirect URI for security integration err = %w", err)
			}
		case "OAUTH_CLIENT_TYPE":
			isTableau := strings.HasSuffix(s.IntegrationType.String, "TABLEAU_DESKTOP") ||
				strings.HasSuffix(s.IntegrationType.String, "TABLEAU_SERVER")
			if !isTableau {
				if err = d.Set("oauth_client_type", v.(string)); err != nil {
					return fmt.Errorf("unable to set OAuth client type for security integration err = %w", err)
				}
			}
		case "OAUTH_ENFORCE_PKCE":
			// Only used for custom OAuth clients (not supported yet)
		case "OAUTH_AUTHORIZATION_ENDPOINT":
			// Only used for custom OAuth clients (not supported yet)
		case "OAUTH_TOKEN_ENDPOINT":
			// Only used for custom OAuth clients (not supported yet)
		case "OAUTH_ALLOWED_AUTHORIZATION_ENDPOINTS":
			// Only used for custom OAuth clients (not supported yet)
		case "OAUTH_ALLOWED_TOKEN_ENDPOINTS":
			// Only used for custom OAuth clients (not supported yet)
		case "PRE_AUTHORIZED_ROLES_LIST":
			// Only used for custom OAuth clients (not supported yet)

		default:
			log.Printf("[WARN] unexpected security integration property %v returned from Snowflake", k)
		}
	}

	return err
}

// UpdateOAuthIntegration implements schema.UpdateFunc.
func UpdateOAuthIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	stmt := snowflake.NewOAuthIntegrationBuilder(id).Alter()

	var runSetStatement bool

	if d.HasChange("oauth_client") {
		runSetStatement = true
		stmt.SetString(`OAUTH_CLIENT`, d.Get("oauth_client").(string))
	}

	if d.HasChange("oauth_redirect_uri") {
		runSetStatement = true
		stmt.SetString(`OAUTH_REDIRECT_URI`, d.Get("oauth_redirect_uri").(string))
	}

	if d.HasChange("oauth_client_type") {
		runSetStatement = true
		stmt.SetString(`OAUTH_CLIENT_TYPE`, d.Get("oauth_client_type").(string))
	}

	if d.HasChange("oauth_issue_refresh_tokens") {
		runSetStatement = true
		stmt.SetBool(`OAUTH_ISSUE_REFRESH_TOKENS`, d.Get("oauth_issue_refresh_tokens").(bool))
	}

	if d.HasChange("oauth_refresh_token_validity") {
		runSetStatement = true
		stmt.SetInt(`OAUTH_REFRESH_TOKEN_VALIDITY`, d.Get("oauth_refresh_token_validity").(int))
	}

	if d.HasChange("oauth_use_secondary_roles") {
		runSetStatement = true
		stmt.SetString(`OAUTH_USE_SECONDARY_ROLES`, d.Get("oauth_use_secondary_roles").(string))
	}

	if d.HasChange("blocked_roles_list") {
		runSetStatement = true
		stmt.SetStringList(`BLOCKED_ROLES_LIST`, expandStringList(d.Get("blocked_roles_list").(*schema.Set).List()))
	}

	if d.HasChange("enabled") {
		runSetStatement = true
		stmt.SetBool(`ENABLED`, d.Get("enabled").(bool))
	}

	if d.HasChange("comment") {
		runSetStatement = true
		stmt.SetString(`COMMENT`, d.Get("comment").(string))
	}

	if runSetStatement {
		if err := snowflake.Exec(db, stmt.Statement()); err != nil {
			return fmt.Errorf("error updating security integration err = %w", err)
		}
	}

	return ReadOAuthIntegration(d, meta)
}

// DeleteOAuthIntegration implements schema.DeleteFunc.
func DeleteOAuthIntegration(d *schema.ResourceData, meta interface{}) error {
	return DeleteResource("", snowflake.NewOAuthIntegrationBuilder)(d, meta)
}
