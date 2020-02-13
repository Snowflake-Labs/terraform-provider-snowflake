package resources

import (
	"database/sql"
	"fmt"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var storageIntegrationSchema = map[string]*schema.Schema{
	// The first part of the schema is shared between all integration vendors
	"name": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"comment": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"type": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "EXTERNAL_STAGE",
		ValidateFunc: validation.StringInSlice([]string{"EXTERNAL_STAGE"}, true),
	},
	"enabled": &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"storage_allowed_locations": &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: "Explicitly limits external stages that use the integration to reference one or more storage locations.",
	},
	"storage_blocked_locations": &schema.Schema{
		Type:     schema.TypeList,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Optional: true,
		Description: "Explicitly prohibits external stages that use the integration from referencing one or more storage locations	.",
	},
	// This part of the schema is the cloudProviderParams in the Snowflake documentation and differs between vendors
	"storage_provider": &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"S3", "GCS", "AZURE"}, false),
	},
	"storage_aws_role_arn": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"azure_tenant_id": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"created_on": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the storage integration was created.",
	},
}

// StorageIntegration returns a pointer to the resource representing a storage integration
func StorageIntegration() *schema.Resource {
	return &schema.Resource{
		Create: CreateStorageIntegration,
		Read:   ReadStorageIntegration,
		Update: UpdateStorageIntegration,
		Delete: DeleteStorageIntegration,
		Exists: StorageIntegrationExists,

		Schema: storageIntegrationSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateStorageIntegration implements schema.CreateFunc
func CreateStorageIntegration(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)

	stmt := snowflake.StorageIntegration(name).Create()

	// Set required fields
	stmt.SetString(`TYPE`, data.Get("type").(string))
	stmt.SetBool(`ENABLED`, data.Get("enabled").(bool))

	setStorageAllowedLocations(data, stmt)

	// Set optional fields
	if v, ok := data.GetOk("comment"); ok {
		stmt.SetString(`COMMENT`, v.(string))
	}

	if _, ok := data.GetOk("storage_blocked_locations"); ok {
		setStorageBlockedLocations(data, stmt)
	}

	// Now, set the storage provider
	err := setStorageProviderSettings(data, stmt)
	if err != nil {
		return err
	}

	err = DBExec(db, stmt.Statement())
	if err != nil {
		return fmt.Errorf("error creating storage integration: %w", err)
	}

	data.SetId(name)

	return ReadStorageIntegration(data, meta)
}

// ReadStorageIntegration implements schema.ReadFunc
func ReadStorageIntegration(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.StorageIntegration(data.Id()).Show()
	row := db.QueryRow(stmt)

	var name, integrationType, category, comment, createdOn sql.NullString
	var enabled sql.NullBool
	err := row.Scan(&name, &integrationType, &category, &enabled, &comment, &createdOn)
	if err != nil {
		return err
	}

	// Note: category must be STORAGE or something is broken
	if c := category.String; c != "STORAGE" {
		return fmt.Errorf("Expected %v to be a STORAGE integration, got %v", id, c)
	}

	err = data.Set("name", name.String)
	if err != nil {
		return err
	}
	err = data.Set("comment", comment.String)
	if err != nil {
		return err
	}

	err = data.Set("type", integrationType.String)
	if err != nil {
		return err
	}

	err = data.Set("created_on", createdOn.String)
	if err != nil {
		return err
	}

	err = data.Set("enabled", enabled.Bool)

	return err
}

// UpdateStorageIntegration implements schema.UpdateFunc
func UpdateStorageIntegration(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.StorageIntegration(id).Alter()

	if data.HasChange("comment") {
		stmt.SetString("COMMENT", data.Get("comment").(string))
	}

	if data.HasChange("type") {
		stmt.SetString("TYPE", data.Get("type").(string))
	}

	if data.HasChange("enabled") {
		stmt.SetBool(`ENABLED`, data.Get("enabled").(bool))
	}

	if data.HasChange("storage_allowed_locations") {
		setStorageAllowedLocations(data, stmt)
	}

	if data.HasChange("storage_blocked_locations") {
		setStorageBlockedLocations(data, stmt)
	}

	if data.HasChange("storage_provider") {
		setStorageProviderSettings(data, stmt)
	} else {
		if data.HasChange("storage_aws_role_arn") {
			stmt.SetString("STORAGE_AWS_ROLE_ARN", data.Get("storage_aws_role_arn").(string))
		}
		if data.HasChange("azure_tenant_id") {
			stmt.SetString("AZURE_TENANT_ID", data.Get("azure_tenant_id").(string))
		}
	}

	err := DBExec(db, stmt.Statement())
	if err != nil {
		return fmt.Errorf("error updating storage integration: %w", err)
	}

	return ReadStorageIntegration(data, meta)
}

// DeleteStorageIntegration implements schema.DeleteFunc
func DeleteStorageIntegration(data *schema.ResourceData, meta interface{}) error {
	return DeleteResource("", snowflake.StorageIntegration)(data, meta)
}

// StorageIntegrationExists implements schema.ExistsFunc
func StorageIntegrationExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.StorageIntegration(id).Show()
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

func setStorageAllowedLocations(data *schema.ResourceData, stmt snowflake.SettingBuilder) {
	allowed := expandStringListToStorageLocations(data.Get("storage_allowed_locations").([]interface{}))
	stmt.SetString("STORAGE_ALLOWED_LOCATIONS", allowed)
}

func setStorageBlockedLocations(data *schema.ResourceData, stmt snowflake.SettingBuilder) {
	blocked := expandStringListToStorageLocations(data.Get("storage_blocked_locations").([]interface{}))
	stmt.SetString(`STORAGE_BLOCKED_LOCATIONS`, blocked)
}

func setStorageProviderSettings(data *schema.ResourceData, stmt snowflake.SettingBuilder) error {
	storageProvider := data.Get("storage_provider").(string)
	stmt.SetString("STORAGE_PROVIDER", storageProvider)

	switch storageProvider {
	case "S3":
		v, ok := data.GetOk("storage_aws_role_arn")
		if !ok {
			return fmt.Errorf("If you use the S3 storage provider you must specify a storage_aws_role_arn")
		}
		stmt.SetString(`STORAGE_AWS_ROLE_ARN`, v.(string))
	case "AZURE":
		v, ok := data.GetOk("azure_tenant_id")
		if !ok {
			return fmt.Errorf("If you use the Azure storage provider you must specify an azure_tenant_id")
		}
		stmt.SetString(`AZURE_TENANT_ID`, v.(string))
	}

	return nil
}
