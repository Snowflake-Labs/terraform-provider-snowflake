// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var apiIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the name of the API integration. This name follows the rules for Object Identifiers. The name should be unique among api integrations in your account.",
	},
	"api_provider": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"aws_api_gateway", "aws_private_api_gateway", "azure_api_management", "aws_gov_api_gateway", "aws_gov_private_api_gateway", "google_api_gateway"}, false),
		Description:  "Specifies the HTTPS proxy service type.",
	},
	"api_aws_role_arn": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "ARN of a cloud platform role.",
	},
	// Computed. Info you get by issuing a 'DESCRIBE INTEGRATION <name>' command (API_AWS_IAM_USER_ARN)
	"api_aws_iam_user_arn": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowflake user that will attempt to assume the AWS role.",
	},
	// Computed. Info you get by issuing a 'DESCRIBE INTEGRATION <name>' command (API_AWS_EXTERNAL_ID)
	"api_aws_external_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The external ID that Snowflake will use when assuming the AWS role.",
	},
	"azure_tenant_id": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "Specifies the ID for your Office 365 tenant that all Azure API Management instances belong to.",
	},
	"azure_ad_application_id": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "The 'Application (client) id' of the Azure AD app for your remote service.",
	},
	// Computed. Info you get by issuing a 'DESCRIBE INTEGRATION <name>' command (AZURE_MULTI_TENANT_APP_NAME)
	"azure_multi_tenant_app_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	// Computed. Info you get by issuing a 'DESCRIBE INTEGRATION <name>' command (AZURE_CONSENT_URL)
	"azure_consent_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"google_audience": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "The audience claim when generating the JWT (JSON Web Token) to authenticate to the Google API Gateway.",
	},
	"api_gcp_service_account": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The service account used for communication with the Google API Gateway.",
		Computed:    true,
	},
	"api_allowed_prefixes": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: "Explicitly limits external functions that use the integration to reference one or more HTTPS proxy service endpoints and resources within those proxies.",
		MinItems:    1,
	},
	"api_blocked_prefixes": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Lists the endpoints and resources in the HTTPS proxy service that are not allowed to be called from Snowflake.",
	},
	"api_key": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "The API key (also called a “subscription key”).",
	},
	"enabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Specifies whether this API integration is enabled or disabled. If the API integration is disabled, any external function that relies on it will not work.",
	},
	"comment": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the API integration was created.",
	},
}

// APIIntegration returns a pointer to the resource representing an api integration.
func APIIntegration() *schema.Resource {
	return &schema.Resource{
		Create: CreateAPIIntegration,
		Read:   ReadAPIIntegration,
		Update: UpdateAPIIntegration,
		Delete: DeleteAPIIntegration,

		Schema: apiIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateAPIIntegration implements schema.CreateFunc.
func CreateAPIIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)

	stmt := snowflake.NewAPIIntegrationBuilder(name).Create()

	// Set required fields
	stmt.SetBool(`ENABLED`, d.Get("enabled").(bool))

	stmt.SetStringList("API_ALLOWED_PREFIXES", expandStringList(d.Get("api_allowed_prefixes").([]interface{})))

	// Set optional fields
	if _, ok := d.GetOk("api_blocked_prefixes"); ok {
		stmt.SetStringList("API_BLOCKED_PREFIXES", expandStringList(d.Get("api_blocked_prefixes").([]interface{})))
	}

	if _, ok := d.GetOk("api_key"); ok {
		stmt.SetString("API_KEY", d.Get("api_key").(string))
	}

	if _, ok := d.GetOk("comment"); ok {
		stmt.SetString("COMMENT", d.Get("comment").(string))
	}

	// Now, set the API provider
	if err := setAPIProviderSettings(d, stmt); err != nil {
		return err
	}

	if err := snowflake.Exec(db, stmt.Statement()); err != nil {
		return fmt.Errorf("error creating api integration: %w", err)
	}

	d.SetId(name)

	return ReadAPIIntegration(d, meta)
}

// ReadAPIIntegration implements schema.ReadFunc.
func ReadAPIIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	stmt := snowflake.NewAPIIntegrationBuilder(id).Show()
	row := snowflake.QueryRow(db, stmt)

	// Some properties can come from the SHOW INTEGRATION call

	s, err := snowflake.ScanAPIIntegration(row)
	if err != nil {
		// If no such resource exists, it is not an error but rather not exist
		if err.Error() == snowflake.ErrNoRowInRS {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not show api integration: %w", err)
	}

	// Note: category must be API or something is broken
	if c := s.Category.String; c != "API" {
		return fmt.Errorf("expected %v to be an api integration, got %v", id, c)
	}

	if err := d.Set("name", s.Name.String); err != nil {
		return err
	}

	if err := d.Set("comment", s.Comment.String); err != nil {
		return err
	}

	if err := d.Set("created_on", s.CreatedOn.String); err != nil {
		return err
	}

	if err := d.Set("enabled", s.Enabled.Bool); err != nil {
		return err
	}

	// Some properties come from the DESCRIBE INTEGRATION call
	// We need to grab them in a loop
	var k, pType string
	var v, unused interface{}
	stmt = snowflake.NewAPIIntegrationBuilder(id).Describe()
	rows, err := db.Query(stmt)
	if err != nil {
		return fmt.Errorf("could not describe api integration: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&k, &pType, &v, &unused); err != nil {
			return err
		}
		switch k {
		case "ENABLED":
			// We set this using the SHOW INTEGRATION call so let's ignore it here
		case "API_ALLOWED_PREFIXES":
			if err := d.Set("api_allowed_prefixes", strings.Split(v.(string), ",")); err != nil {
				return err
			}
		case "API_BLOCKED_PREFIXES":
			if val := v.(string); val != "" {
				if err := d.Set("api_blocked_prefixes", strings.Split(val, ",")); err != nil {
					return err
				}
			}
		case "API_AWS_IAM_USER_ARN":
			if err := d.Set("api_aws_iam_user_arn", v.(string)); err != nil {
				return err
			}
		case "API_AWS_ROLE_ARN":
			if err := d.Set("api_aws_role_arn", v.(string)); err != nil {
				return err
			}
		case "API_AWS_EXTERNAL_ID":
			if err := d.Set("api_aws_external_id", v.(string)); err != nil {
				return err
			}
		case "AZURE_CONSENT_URL":
			if err := d.Set("azure_consent_url", v.(string)); err != nil {
				return err
			}
		case "AZURE_MULTI_TENANT_APP_NAME":
			if err := d.Set("azure_multi_tenant_app_name", v.(string)); err != nil {
				return err
			}
		case "GOOGLE_AUDIENCE":
			if err := d.Set("google_audience", v.(string)); err != nil {
				return err
			}
		case "API_GCP_SERVICE_ACCOUNT":
			if err := d.Set("api_gcp_service_account", v.(string)); err != nil {
				return err
			}
		default:
			log.Printf("[WARN] unexpected api integration property %v returned from Snowflake", k)
		}
	}

	return err
}

// UpdateAPIIntegration implements schema.UpdateFunc.
func UpdateAPIIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	stmt := snowflake.NewAPIIntegrationBuilder(id).Alter()

	var runSetStatement bool

	if d.HasChange("enabled") {
		runSetStatement = true
		stmt.SetBool(`ENABLED`, d.Get("enabled").(bool))
	}

	if d.HasChange("api_allowed_prefixes") {
		runSetStatement = true
		stmt.SetStringList("API_ALLOWED_PREFIXES", expandStringList(d.Get("api_allowed_prefixes").([]interface{})))
	}

	if d.HasChange("api_key") {
		runSetStatement = true
		stmt.SetString("API_KEY", d.Get("api_key").(string))
	}

	if d.HasChange("comment") {
		runSetStatement = true
		stmt.SetString("COMMENT", d.Get("comment").(string))
	}

	// We need to UNSET this if we remove all api blocked prefixes.
	if d.HasChange("api_blocked_prefixes") {
		v := d.Get("api_blocked_prefixes").([]interface{})
		if len(v) == 0 {
			if err := snowflake.Exec(db, fmt.Sprintf(`ALTER API INTEGRATION %v UNSET API_BLOCKED_PREFIXES`, id)); err != nil {
				return fmt.Errorf("error unsetting api_blocked_prefixes: %w", err)
			}
		} else {
			runSetStatement = true
			stmt.SetStringList("API_BLOCKED_PREFIXES", expandStringList(v))
		}
	}

	if d.HasChange("api_provider") {
		runSetStatement = true
		err := setAPIProviderSettings(d, stmt)
		if err != nil {
			return err
		}
	} else {
		if d.HasChange("api_aws_role_arn") {
			runSetStatement = true
			stmt.SetString("API_AWS_ROLE_ARN", d.Get("api_aws_role_arn").(string))
		}
		if d.HasChange("azure_tenant_id") {
			runSetStatement = true
			stmt.SetString("AZURE_TENANT_ID", d.Get("azure_tenant_id").(string))
		}
		if d.HasChange("azure_ad_application_id") {
			runSetStatement = true
			stmt.SetString("AZURE_AD_APPLICATION_ID", d.Get("azure_ad_application_id").(string))
		}
		if d.HasChange("google_audience") {
			runSetStatement = true
			stmt.SetString("GOOGLE_AUDIENCE", d.Get("google_audience").(string))
		}
	}

	if runSetStatement {
		if err := snowflake.Exec(db, stmt.Statement()); err != nil {
			return fmt.Errorf("error updating api integration: %w", err)
		}
	}

	return ReadAPIIntegration(d, meta)
}

// DeleteAPIIntegration implements schema.DeleteFunc.
func DeleteAPIIntegration(d *schema.ResourceData, meta interface{}) error {
	return DeleteResource("", snowflake.NewAPIIntegrationBuilder)(d, meta)
}

func setAPIProviderSettings(data *schema.ResourceData, stmt snowflake.SettingBuilder) error {
	apiProvider := data.Get("api_provider").(string)
	stmt.SetRaw("API_PROVIDER=" + apiProvider)

	switch apiProvider {
	case "aws_api_gateway", "aws_private_api_gateway", "aws_gov_api_gateway", "aws_gov_private_api_gateway":
		v, ok := data.GetOk("api_aws_role_arn")
		if !ok {
			return fmt.Errorf("if you use AWS api provider you must specify an api_aws_role_arn")
		}
		stmt.SetString(`API_AWS_ROLE_ARN`, v.(string))
	case "azure_api_management":
		v, ok := data.GetOk("azure_tenant_id")
		if !ok {
			return fmt.Errorf("if you use the Azure api provider you must specify an azure_tenant_id")
		}
		stmt.SetString(`AZURE_TENANT_ID`, v.(string))

		v, ok = data.GetOk("azure_ad_application_id")
		if !ok {
			return fmt.Errorf("if you use the Azure api provider you must specify an azure_ad_application_id")
		}
		stmt.SetString(`AZURE_AD_APPLICATION_ID`, v.(string))
	case "google_api_gateway":
		v, ok := data.GetOk("google_audience")
		if !ok {
			return fmt.Errorf("if you use GCP api provider you must specify a google_audience")
		}
		stmt.SetString(`GOOGLE_AUDIENCE`, v.(string))
	default:
		return fmt.Errorf("unexpected provider %v", apiProvider)
	}

	return nil
}
