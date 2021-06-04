package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var notificationIntegrationSchema = map[string]*schema.Schema{
	// The first part of the schema is shared between all integration vendors
	"name": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"enabled": &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"type": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "QUEUE",
		ValidateFunc: validation.StringInSlice([]string{"QUEUE"}, true),
		Description:  "A type of integration",
	},
	"notification_provider": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "AZURE_STORAGE_QUEUE",
		ValidateFunc: validation.StringInSlice([]string{"AZURE_STORAGE_QUEUE"}, true),
		Description:  "The third-party cloud message queuing service (e.g. AZURE_STORAGE_QUEUE)",
	},
	"azure_storage_queue_primary_uri": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The queue ID for the Azure Queue Storage queue created for Event Grid notifications",
	},
	"azure_tenant_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The ID of the Azure Active Directory tenant used for identity management",
	},
	"comment": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "A comment for the integration",
	},
	"created_on": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the notification integration was created.",
	},
}

// NotificationIntegration returns a pointer to the resource representing a notification integration
func NotificationIntegration() *schema.Resource {
	return &schema.Resource{
		Create: CreateNotificationIntegration,
		Read:   ReadNotificationIntegration,
		Update: UpdateNotificationIntegration,
		Delete: DeleteNotificationIntegration,
		Exists: NotificationIntegrationExists,

		Schema: notificationIntegrationSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateNotificationIntegration implements schema.CreateFunc
func CreateNotificationIntegration(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)

	stmt := snowflake.NotificationIntegration(name).Create()

	// Set required fields
	stmt.SetString(`TYPE`, data.Get("type").(string))
	stmt.SetBool(`ENABLED`, data.Get("enabled").(bool))

	// Set optional fields
	if v, ok := data.GetOk("comment"); ok {
		stmt.SetString(`COMMENT`, v.(string))
	}
	if v, ok := data.GetOk("notification_provider"); ok {
		stmt.SetString(`NOTIFICATION_PROVIDER`, v.(string))
	}
	if v, ok := data.GetOk("azure_storage_queue_primary_uri"); ok {
		stmt.SetString(`AZURE_STORAGE_QUEUE_PRIMARY_URI`, v.(string))
	}
	if v, ok := data.GetOk("azure_tenant_id"); ok {
		stmt.SetString(`AZURE_TENANT_ID`, v.(string))
	}

	err := snowflake.Exec(db, stmt.Statement())
	if err != nil {
		return fmt.Errorf("error creating notification integration: %w", err)
	}

	data.SetId(name)

	return ReadNotificationIntegration(data, meta)
}

// ReadNotificationIntegration implements schema.ReadFunc
func ReadNotificationIntegration(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.NotificationIntegration(data.Id()).Show()
	row := snowflake.QueryRow(db, stmt)

	// Some properties can come from the SHOW INTEGRATION call

	s, err := snowflake.ScanNotificationIntegration(row)
	if err != nil {
		return fmt.Errorf("Could not show notification integration: %w", err)
	}

	// Note: category must be NOTIFICATION or something is broken
	if c := s.Category.String; c != "NOTIFICATION" {
		return fmt.Errorf("Expected %v to be a NOTIFICATION integration, got %v", id, c)
	}

	if err := data.Set("name", s.Name.String); err != nil {
		return err
	}

	// Snowflake returns "QUEUE - AZURE_STORAGE_QUEUE" instead of simple "QUEUE" as a type
	// so it needs to be parsed in order to not show a diff in Terraform
	typeParts := strings.Split(s.Type.String, "-")
	parsedType := strings.TrimSpace(typeParts[0])
	if err = data.Set("type", parsedType); err != nil {
		return err
	}

	if err := data.Set("created_on", s.CreatedOn.String); err != nil {
		return err
	}

	if err := data.Set("enabled", s.Enabled.Bool); err != nil {
		return err
	}

	// Some properties come from the DESCRIBE INTEGRATION call
	// We need to grab them in a loop
	var k, pType string
	var v, d interface{}
	stmt = snowflake.NotificationIntegration(data.Id()).Describe()
	rows, err := db.Query(stmt)
	if err != nil {
		return fmt.Errorf("Could not describe notification integration: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&k, &pType, &v, &d); err != nil {
			return err
		}
		switch k {
		case "ENABLED":
			// We set this using the SHOW INTEGRATION call so let's ignore it here
		case "NOTIFICATION_PROVIDER":
			if err = data.Set("notification_provider", v.(string)); err != nil {
				return err
			}
		case "AZURE_STORAGE_QUEUE_PRIMARY_URI":
			if err = data.Set("azure_storage_queue_primary_uri", v.(string)); err != nil {
				return err
			}
		case "AZURE_TENANT_ID":
			if err = data.Set("azure_tenant_id", v.(string)); err != nil {
				return err
			}
		default:
			log.Printf("[WARN] unexpected property %v returned from Snowflake", k)
		}
	}

	return err
}

// UpdateNotificationIntegration implements schema.UpdateFunc
func UpdateNotificationIntegration(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.NotificationIntegration(id).Alter()

	// This is required in case the only change is to UNSET STORAGE_ALLOWED_LOCATIONS.
	// Not sure if there is a more elegant way of determining this
	var runSetStatement bool

	if data.HasChange("comment") {
		runSetStatement = true
		stmt.SetString("COMMENT", data.Get("comment").(string))
	}

	if data.HasChange("type") {
		runSetStatement = true
		stmt.SetString("TYPE", data.Get("type").(string))
	}

	if data.HasChange("enabled") {
		runSetStatement = true
		stmt.SetBool(`ENABLED`, data.Get("enabled").(bool))
	}

	if data.HasChange("notification_provider") {
		runSetStatement = true
		stmt.SetString("NOTIFICATION_PROVIDER", data.Get("notification_provider").(string))
	}

	if data.HasChange("azure_storage_queue_primary_uri") {
		runSetStatement = true
		stmt.SetString("AZURE_STORAGE_QUEUE_PRIMARY_URI", data.Get("azure_storage_queue_primary_uri").(string))
	}

	if data.HasChange("azure_tenant_id") {
		runSetStatement = true
		stmt.SetString("AZURE_TENANT_ID", data.Get("azure_tenant_id").(string))
	}

	if runSetStatement {
		if err := snowflake.Exec(db, stmt.Statement()); err != nil {
			return fmt.Errorf("error updating notification integration: %w", err)
		}
	}

	return ReadNotificationIntegration(data, meta)
}

// DeleteNotificationIntegration implements schema.DeleteFunc
func DeleteNotificationIntegration(data *schema.ResourceData, meta interface{}) error {
	return DeleteResource("", snowflake.NotificationIntegration)(data, meta)
}

// NotificationIntegrationExists implements schema.ExistsFunc
func NotificationIntegrationExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.NotificationIntegration(id).Show()
	rows, err := db.Query(stmt)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}
	return false, nil
}
