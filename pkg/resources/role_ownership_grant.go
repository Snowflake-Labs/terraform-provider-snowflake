package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var roleOwnershipGrantSchema = map[string]*schema.Schema{
	"on_role_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the role ownership is granted on.",
		ValidateFunc: func(val interface{}, key string) ([]string, []error) {
			var additionalCharsToIgnoreValidation []string
			return snowflake.ValidateIdentifier(val, additionalCharsToIgnoreValidation)
		},
	},
	"to_role_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the role to grant ownership. Please ensure that the role that terraform is using is granted access.",
		ValidateFunc: func(val interface{}, key string) ([]string, []error) {
			var additionalCharsToIgnoreValidation []string
			return snowflake.ValidateIdentifier(val, additionalCharsToIgnoreValidation)
		},
	},
	"current_grants": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies whether to remove or transfer all existing outbound privileges on the object when ownership is transferred to a new role.",
		Default:     "COPY",
		ValidateFunc: validation.StringInSlice([]string{
			"COPY",
			"REVOKE",
		}, true),
	},
}

func RoleOwnershipGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateRoleOwnershipGrant,
		Read:   ReadRoleOwnershipGrant,
		Delete: DeleteRoleOwnershipGrant,
		Update: UpdateRoleOwnershipGrant,
		Schema: roleOwnershipGrantSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateRoleOwnershipGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	onRoleName := d.Get("on_role_name").(string)
	toRoleName := d.Get("to_role_name").(string)
	currentGrants := d.Get("current_grants").(string)

	g := snowflake.RoleOwnershipGrant(onRoleName, currentGrants)
	err := snowflake.Exec(db, g.Role(toRoleName).Grant())
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf(`%s|%s|%s`, onRoleName, toRoleName, currentGrants))

	return ReadRoleOwnershipGrant(d, meta)
}

func ReadRoleOwnershipGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	log.Println(d.Id())
	onRoleName := strings.Split(d.Id(), "|")[0]
	currentGrants := strings.Split(d.Id(), "|")[2]

	stmt := fmt.Sprintf("SHOW ROLES LIKE '%s'", onRoleName)
	row := snowflake.QueryRow(db, stmt)

	grant, err := snowflake.ScanRoleOwnershipGrant(row)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] role (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	if onRoleName != grant.Name.String {
		return fmt.Errorf("no role found like '%s'", onRoleName)
	}

	grant.Name.String = strings.TrimPrefix(grant.Name.String, `"`)
	grant.Name.String = strings.TrimSuffix(grant.Name.String, `"`)
	err = d.Set("on_role_name", grant.Name.String)
	if err != nil {
		return err
	}

	grant.Owner.String = strings.TrimPrefix(grant.Owner.String, `"`)
	grant.Owner.String = strings.TrimSuffix(grant.Owner.String, `"`)
	err = d.Set("to_role_name", grant.Owner.String)
	if err != nil {
		return err
	}

	err = d.Set("current_grants", currentGrants)
	if err != nil {
		return err
	}

	return nil
}

func UpdateRoleOwnershipGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	onRoleName := d.Get("on_role_name").(string)
	toRoleName := d.Get("to_role_name").(string)
	currentGrants := d.Get("current_grants").(string)

	d.SetId(fmt.Sprintf(`%s|%s|%s`, onRoleName, toRoleName, currentGrants))

	g := snowflake.RoleOwnershipGrant(onRoleName, currentGrants)
	err := snowflake.Exec(db, g.Role(toRoleName).Grant())
	if err != nil {
		return err
	}

	return ReadRoleOwnershipGrant(d, meta)
}

func DeleteRoleOwnershipGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	onRoleName := d.Get("on_role_name").(string)
	currentGrants := d.Get("current_grants").(string)

	g := snowflake.RoleOwnershipGrant(onRoleName, currentGrants)
	err := snowflake.Exec(db, g.Role("ACCOUNTADMIN").Revoke())
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
