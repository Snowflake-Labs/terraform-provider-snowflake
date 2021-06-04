package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pkg/errors"
)

var scimIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the name of the SCIM integration. This name follows the rules for Object Identifiers. The name should be unique among security integrations in your account.",
	},
	"scim_client": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the client type for the scim integration",
		ValidateFunc: validation.StringInSlice([]string{
			"OKTA", "AZURE", "CUSTOM",
		}, true),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalize := func(s string) string {
				return strings.ToUpper(strings.Replace(s, "-", "", -1))
			}
			return normalize(old) == normalize(new)
		},
	},
	"provisioner_role": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specify the SCIM role in Snowflake that owns any users and roles that are imported from the identity provider into Snowflake using SCIM.",
		ValidateFunc: validation.StringInSlice([]string{
			"OKTA_PROVISIONER", "AAD_PROVISIONER", "GENERIC_SCIM_PROVISIONER",
		}, true),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalize := func(s string) string {
				return strings.ToUpper(strings.Replace(s, "-", "", -1))
			}
			return normalize(old) == normalize(new)
		},
	},
	"network_policy": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies an existing network policy active for your account. The network policy restricts the list of user IP addresses when exchanging an authorization code for an access or refresh token and when using a refresh token to obtain a new access token. If this parameter is not set, the network policy for the account (if any) is used instead.",
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the SCIM integration was created.",
	},
}

// SCIMIntegration returns a pointer to the resource representing a network policy
func SCIMIntegration() *schema.Resource {
	return &schema.Resource{
		Create: CreateSCIMIntegration,
		Read:   ReadSCIMIntegration,
		Update: UpdateSCIMIntegration,
		Delete: DeleteSCIMIntegration,

		Schema: scimIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateSCIMIntegration implements schema.CreateFunc
func CreateSCIMIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)

	stmt := snowflake.ScimIntegration(name).Create()

	// Set required fields
	stmt.SetRaw(`TYPE=SCIM`)
	stmt.SetString(`SCIM_CLIENT`, d.Get("scim_client").(string))
	stmt.SetString(`RUN_AS_ROLE`, d.Get("provisioner_role").(string))

	// Set optional fields
	if _, ok := d.GetOk("network_policy"); ok {
		stmt.SetString(`NETWORK_POLICY`, d.Get("network_policy").(string))
	}

	err := snowflake.Exec(db, stmt.Statement())
	if err != nil {
		return errors.Wrap(err, "error creating security integration")
	}

	d.SetId(name)

	return ReadSCIMIntegration(d, meta)
}

// ReadSCIMIntegration implements schema.ReadFunc
func ReadSCIMIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	stmt := snowflake.ScimIntegration(id).Show()
	row := snowflake.QueryRow(db, stmt)

	// Some properties can come from the SHOW INTEGRATION call

	s, err := snowflake.ScanScimIntegration(row)
	if err != nil {
		return errors.Wrap(err, "could not show security integration")
	}

	// Note: category must be Security or something is broken
	if c := s.Category.String; c != "SECURITY" {
		return fmt.Errorf("expected %v to be an Security integration, got %v", id, c)
	}

	if err := d.Set("scim_client", strings.TrimPrefix(s.IntegrationType.String, "SCIM - ")); err != nil {
		return err
	}

	if err := d.Set("name", s.Name.String); err != nil {
		return err
	}

	if err := d.Set("created_on", s.CreatedOn.String); err != nil {
		return err
	}

	// Some properties come from the DESCRIBE INTEGRATION call
	// We need to grab them in a loop
	var k, pType string
	var v, unused interface{}
	stmt = snowflake.ScimIntegration(id).Describe()
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
		case "NETWORK_POLICY":
			if err = d.Set("network_policy", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set network policy for security integration")
			}
		case "RUN_AS_ROLE":
			if err = d.Set("provisioner_role", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set provisioner role for security integration")
			}
		default:
			log.Printf("[WARN] unexpected security integration property %v returned from Snowflake", k)
		}
	}

	return err
}

// UpdateSCIMIntegration implements schema.UpdateFunc
func UpdateSCIMIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	stmt := snowflake.ScimIntegration(id).Alter()

	var runSetStatement bool

	if d.HasChange("scim_client") {
		runSetStatement = true
		stmt.SetString(`SCIM_CLIENT`, d.Get("scim_client").(string))
	}

	if d.HasChange("provisioner_role") {
		runSetStatement = true
		stmt.SetString(`RUN_AS_ROLE`, d.Get("provisioner_role").(string))
	}

	// We need to UNSET this if we remove all api blocked prefixes.
	if d.HasChange("network_policy") {
		v := d.Get("network_policy").(string)
		if len(v) == 0 {
			err := snowflake.Exec(db, fmt.Sprintf(`ALTER SECURITY INTEGRATION %v UNSET NETWORK_POLICY`, id))
			if err != nil {
				return errors.Wrap(err, "error unsetting network_policy")
			}
		} else {
			runSetStatement = true
			stmt.SetString(`NETWORK_POLICY`, d.Get("network_policy").(string))
		}
	}

	if runSetStatement {
		if err := snowflake.Exec(db, stmt.Statement()); err != nil {
			return errors.Wrap(err, "error updating security integration")
		}
	}

	return ReadSCIMIntegration(d, meta)
}

// DeleteSCIMIntegration implements schema.DeleteFunc
func DeleteSCIMIntegration(d *schema.ResourceData, meta interface{}) error {
	return DeleteResource("", snowflake.ScimIntegration)(d, meta)
}
