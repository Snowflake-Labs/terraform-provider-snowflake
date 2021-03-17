package resources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var notificationIntegrationSchema = map[string]*schema.Schema{
	// The first part of the schema is shared between all integration vendors
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"type": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "QUEUE",
		ValidateFunc: validation.StringInSlice([]string{"QUEUE"}, true),
		ForceNew:     true,
	},
	"enabled": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	// This part of the schema is the cloudProviderParams in the Snowflake documentation and differs between vendors
	"notification_provider": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"GCP_PUBSUB", "AZURE_STORAGE_QUEUE"}, false),
	},
	"gcp_pubsub_subscription_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The subscription id that Snowflake will listen to when using the GCP_PUBSUB provider.",
	},
	"azure_tenant_id": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"azure_storage_queue_primary_uri": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The consent URL that is used to create an Azure Snowflake service principle inside your tenant.",
	},
	"gcp_pubsub_service_account": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "This is the name of the Snowflake Google Service Account created for your account.",
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the storage integration was created.",
	},
}

// Notification returns a pointer to the resource representing a notification integration
func NotificationIntegration() *schema.Resource {
	return &schema.Resource{
		Create: CreateNotificationIntegration,
		Read:   ReadNotificationIntegration,
		Update: UpdateNotificationIntegration,
		Delete: DeleteNotificationIntegration,

		Schema: notificationIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateNotificationIntegration implements schema.CreateFunc
func CreateNotificationIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)

	stmt := snowflake.NotificationIntegration(name).Create()

	// Set required fields
	stmt.SetString(`TYPE`, d.Get("type").(string))
	stmt.SetBool(`ENABLED`, d.Get("enabled").(bool))

	// Set optional fields
	if v, ok := d.GetOk("comment"); ok {
		stmt.SetString(`COMMENT`, v.(string))
	}

	// Now, set the notification provider
	err := setNotificationProviderSettings(d, stmt)
	if err != nil {
		return err
	}

	err = snowflake.Exec(db, stmt.Statement())
	if err != nil {
		return fmt.Errorf("error creating notification integration: %w", err)
	}

	d.SetId(name)

	return ReadNotificationIntegration(d, meta)
}

// ReadNotificationIntegration implements schema.ReadFunc
func ReadNotificationIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	stmt := snowflake.NotificationIntegration(d.Id()).Show()
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

	if err := d.Set("name", s.Name.String); err != nil {
		return err
	}

	if err := d.Set("type", s.IntegrationType.String); err != nil {
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
	stmt = snowflake.NotificationIntegration(d.Id()).Describe()
	rows, err := db.Query(stmt)
	if err != nil {
		return fmt.Errorf("Could not describe notification integration: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&k, &pType, &v, &unused); err != nil {
			return err
		}
		switch k {
		case "ENABLED":
			// We set this using the SHOW INTEGRATION call so let's ignore it here

		case "GCP_PUBSUB_SUBSCRIPTION_NAME":
			if err = d.Set("gcp_pubsub_subscription_name", v.(string)); err != nil {
				return err
			}
		case "GCP_PUBSUB_SERVICE_ACCOUNT":
			if err = d.Set("gcp_pubsub_service_account", v.(string)); err != nil {
				return err
			}

		case "AZURE_CONSENT_URL":
			if err = d.Set("azure_consent_url", v.(string)); err != nil {
				return err
			}
		case "AZURE_MULTI_TENANT_APP_NAME":
			if err = d.Set("azure_multi_tenant_app_name", v.(string)); err != nil {
				return err
			}
		default:
			log.Printf("[WARN] unexpected property %v returned from Snowflake", k)
		}
	}

	return err
}

// UpdateNotificationIntegration implements schema.UpdateFunc
func UpdateNotificationIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	stmt := snowflake.NotificationIntegration(id).Alter()

	// Alter is available only for Azure Notification Integrations
	var runSetStatement bool

	if d.HasChange("comment") {
		runSetStatement = true
		stmt.SetString("COMMENT", d.Get("comment").(string))
	}

	if d.HasChange("enabled") {
		runSetStatement = true
		stmt.SetBool(`ENABLED`, d.Get("enabled").(bool))
	}

	if d.HasChange("notification_provider") {
		runSetStatement = true
		err := setNotificationProviderSettings(d, stmt)
		if err != nil {
			return err
		}
	} else {
		if d.HasChange("storage_aws_role_arn") {
			runSetStatement = true
			stmt.SetString("STORAGE_AWS_ROLE_ARN", d.Get("storage_aws_role_arn").(string))
		}
		if d.HasChange("azure_tenant_id") {
			runSetStatement = true
			stmt.SetString("AZURE_TENANT_ID", d.Get("azure_tenant_id").(string))
		}
		if d.HasChange("storage_gcp_service_account") {
			runSetStatement = true
			stmt.SetString("STORAGE_GCP_SERVICE_ACCOUNT", d.Get("storage_gcp_service_account").(string))
		}
	}

	if runSetStatement {
		if err := snowflake.Exec(db, stmt.Statement()); err != nil {
			return fmt.Errorf("error updating storage integration: %w", err)
		}
	}

	return ReadNotificationIntegration(d, meta)
}

// DeleteNotificationIntegration implements schema.DeleteFunc
func DeleteNotificationIntegration(d *schema.ResourceData, meta interface{}) error {
	return DeleteResource("", snowflake.NotificationIntegration)(d, meta)
}

// NotificationIntegrationExists implements schema.ExistsFunc
func NotificationIntegrationExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	id := d.Id()

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

func setNotificationProviderSettings(data *schema.ResourceData, stmt snowflake.SettingBuilder) error {
	notificationProvider := data.Get("notification_provider").(string)
	stmt.SetString("NOTIFICATION_PROVIDER", notificationProvider)

	switch notificationProvider {
	case "AZURE":
		v, ok := data.GetOk("azure_tenant_id")
		if !ok {
			return fmt.Errorf("If you use the Azure storage provider you must specify an azure_tenant_id")
		}
		stmt.SetString(`AZURE_TENANT_ID`, v.(string))
	case "GCP_PUBSUB":
		v, ok := data.GetOk("gcp_pubsub_subscription_name")
		if !ok {
			return fmt.Errorf("If you use the GCP PubSub notification provider you must specify an gcp_pubsub_subscription_name")
		}
		stmt.SetString(`GCP_PUBSUB_SUBSCRIPTION_NAME`, v.(string))
	default:
		return fmt.Errorf("Unexpected provider %v", notificationProvider)
	}

	return nil
}
