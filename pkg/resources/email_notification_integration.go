package resources

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var emailNotificationIntegrationSchema = map[string]*schema.Schema{
	// The first part of the schema is shared between all integration vendors
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"enabled": {
		Type:     schema.TypeBool,
		Required: true,
	},
	"allowed_recipients": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "List of email addresses that should receive notifications.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A comment for the email integration.",
	},
}

// EmailNotificationIntegration returns a pointer to the resource representing a notification integration.
func EmailNotificationIntegration() *schema.Resource {
	return &schema.Resource{
		Create: CreateEmailNotificationIntegration,
		Read:   ReadEmailNotificationIntegration,
		Update: UpdateEmailNotificationIntegration,
		Delete: DeleteEmailNotificationIntegration,

		Schema: emailNotificationIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateEmailNotificationIntegration implements schema.CreateFunc.
func CreateEmailNotificationIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)

	stmt := snowflake.NewNotificationIntegrationBuilder(name).Create()

	stmt.SetString("TYPE", "EMAIL")
	stmt.SetBool(`ENABLED`, d.Get("enabled").(bool))

	if v, ok := d.GetOk("allowed_recipients"); ok {
		stmt.SetStringList(`ALLOWED_RECIPIENTS`, expandStringList(v.(*schema.Set).List()))
	}

	if v, ok := d.GetOk("comment"); ok {
		stmt.SetString(`COMMENT`, v.(string))
	}

	qry := stmt.Statement()
	if err := snowflake.Exec(db, qry); err != nil {
		return fmt.Errorf("error creating notification integration: %w", err)
	}

	d.SetId(name)

	return ReadEmailNotificationIntegration(d, meta)
}

// ReadEmailNotificationIntegration implements schema.ReadFunc.
func ReadEmailNotificationIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	stmt := snowflake.NewEmailNotificationIntegrationBuilder(d.Id()).Show()
	row := snowflake.QueryRow(db, stmt)

	// Some properties can come from the SHOW INTEGRATION call
	s, err := snowflake.ScanEmailNotificationIntegration(row)
	if err != nil {
		return fmt.Errorf("could not show notification integration: %w", err)
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

	// Some properties come from the DESCRIBE INTEGRATION call
	// We need to grab them in a loop
	var k, pType string
	var v, n interface{}
	stmt = snowflake.NewNotificationIntegrationBuilder(d.Id()).Describe()
	rows, err := db.Query(stmt)
	if err != nil {
		return fmt.Errorf("could not describe notification integration: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&k, &pType, &v, &n); err != nil {
			return err
		}
		switch k {
		case "ALLOWED_RECIPIENTS":
			// Empty list returns strange string (it's empty on worksheet level).
			// This is a quick workaround, should be fixed with moving the email integration to SDK.
			r := regexp.MustCompile(`[[:print:]]`)
			if r.MatchString(v.(string)) {
				if err := d.Set("allowed_recipients", strings.Split(v.(string), ",")); err != nil {
					return err
				}
			} else {
				empty := make([]string, 0)
				if err := d.Set("allowed_recipients", empty); err != nil {
					return err
				}
			}
		default:
			log.Printf("[WARN] unexpected property %v returned from Snowflake", k)
		}
	}

	return err
}

// UpdateEmailNotificationIntegration implements schema.UpdateFunc.
func UpdateEmailNotificationIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	stmt := snowflake.NewEmailNotificationIntegrationBuilder(id).Alter()

	if d.HasChange("comment") {
		stmt.SetString("COMMENT", d.Get("comment").(string))
	}

	if d.HasChange("enabled") {
		stmt.SetBool(`ENABLED`, d.Get("enabled").(bool))
	}

	if d.HasChange("allowed_recipients") {
		if v, ok := d.GetOk("allowed_recipients"); ok {
			stmt.SetStringList(`ALLOWED_RECIPIENTS`, expandStringList(v.(*schema.Set).List()))
		} else {
			// raw sql for now; will be updated with SDK rewrite
			// https://docs.snowflake.com/en/sql-reference/sql/alter-notification-integration#syntax
			unset := fmt.Sprintf(`ALTER NOTIFICATION INTEGRATION "%s" UNSET ALLOWED_RECIPIENTS`, id)
			if err := snowflake.Exec(db, unset); err != nil {
				return fmt.Errorf("error unsetting allowed recipients on email notification integration %v err = %w", id, err)
			}
		}
	}

	if err := snowflake.Exec(db, stmt.Statement()); err != nil {
		return fmt.Errorf("error updating notification integration: %w", err)
	}

	return ReadEmailNotificationIntegration(d, meta)
}

// DeleteEmailNotificationIntegration implements schema.DeleteFunc.
func DeleteEmailNotificationIntegration(d *schema.ResourceData, meta interface{}) error {
	return DeleteResource("", snowflake.NewEmailNotificationIntegrationBuilder)(d, meta)
}
