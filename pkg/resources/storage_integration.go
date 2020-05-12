package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Explicitly prohibits external stages that use the integration from referencing one or more storage locations.",
	},
	// This part of the schema is the cloudProviderParams in the Snowflake documentation and differs between vendors
	"storage_provider": &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"S3", "GCS", "AZURE"}, false),
	},
	"storage_aws_external_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The external ID that Snowflake will use when assuming the AWS role.",
	},
	"storage_aws_iam_user_arn": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowflake user that will attempt to assume the AWS role.",
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

	stmt.SetStringList("STORAGE_ALLOWED_LOCATIONS", expandStringList(data.Get("storage_allowed_locations").([]interface{})))

	// Set optional fields
	if v, ok := data.GetOk("comment"); ok {
		stmt.SetString(`COMMENT`, v.(string))
	}

	if _, ok := data.GetOk("storage_blocked_locations"); ok {
		stmt.SetStringList("STORAGE_BLOCKED_LOCATIONS", expandStringList(data.Get("storage_blocked_locations").([]interface{})))
	}

	// Now, set the storage provider
	err := setStorageProviderSettings(data, stmt)
	if err != nil {
		return err
	}

	err = snowflake.Exec(db, stmt.Statement())
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
	row := snowflake.QueryRow(db, stmt)

	// Some properties can come from the SHOW INTEGRATION call

	s, err := snowflake.ScanStorageIntegration(row)
	if err != nil {
		return fmt.Errorf("Could not show storage integration: %w", err)
	}

	// Note: category must be STORAGE or something is broken
	if c := s.Category.String; c != "STORAGE" {
		return fmt.Errorf("Expected %v to be a STORAGE integration, got %v", id, c)
	}

	if err := data.Set("name", s.Name.String); err != nil {
		return err
	}

	if err := data.Set("type", s.IntegrationType.String); err != nil {
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
	stmt = snowflake.StorageIntegration(data.Id()).Describe()
	rows, err := db.Query(stmt)
	if err != nil {
		return fmt.Errorf("Could not describe storage integration: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&k, &pType, &v, &d); err != nil {
			return err
		}
		switch k {
		case "ENABLED":
			// We set this using the SHOW INTEGRATION call so let's ignore it here
		case "STORAGE_PROVIDER":
			if err = data.Set("storage_provider", v.(string)); err != nil {
				return err
			}
		case "STORAGE_ALLOWED_LOCATIONS":
			if err = data.Set("storage_allowed_locations", strings.Split(v.(string), ",")); err != nil {
				return err
			}
		case "STORAGE_BLOCKED_LOCATIONS":
			if val := v.(string); val != "" {
				if err = data.Set("storage_blocked_locations", strings.Split(val, ",")); err != nil {
					return err
				}
			}
		case "STORAGE_AWS_IAM_USER_ARN":
			if err = data.Set("storage_aws_iam_user_arn", v.(string)); err != nil {
				return err
			}
		case "STORAGE_AWS_ROLE_ARN":
			if err = data.Set("storage_aws_role_arn", v.(string)); err != nil {
				return err
			}
		case "STORAGE_AWS_EXTERNAL_ID":
			if err = data.Set("storage_aws_external_id", v.(string)); err != nil {
				return err
			}
		default:
			log.Printf("[WARN] unexpected property %v returned from Snowflake", k)
		}
	}

	return err
}

// UpdateStorageIntegration implements schema.UpdateFunc
func UpdateStorageIntegration(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.StorageIntegration(id).Alter()

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

	if data.HasChange("storage_allowed_locations") {
		runSetStatement = true
		stmt.SetStringList("STORAGE_ALLOWED_LOCATIONS", expandStringList(data.Get("storage_allowed_locations").([]interface{})))
	}

	// We need to UNSET this if we remove all storage blocked locations. I don't think
	// this is documented by Snowflake, but this is how it works.
	//
	// @TODO move the SQL back to the snowflake package
	if data.HasChange("storage_blocked_locations") {
		v := data.Get("storage_blocked_locations").([]interface{})
		if len(v) == 0 {
			err := snowflake.Exec(db, fmt.Sprintf(`ALTER STORAGE INTEGRATION %v UNSET STORAGE_BLOCKED_LOCATIONS`, data.Id()))
			if err != nil {
				return fmt.Errorf("error unsetting storage_blocked_locations: %w", err)
			}
		} else {
			runSetStatement = true
			stmt.SetStringList("STORAGE_BLOCKED_LOCATIONS", expandStringList(v))
		}
	}

	if data.HasChange("storage_provider") {
		runSetStatement = true
		setStorageProviderSettings(data, stmt)
	} else {
		if data.HasChange("storage_aws_role_arn") {
			runSetStatement = true
			stmt.SetString("STORAGE_AWS_ROLE_ARN", data.Get("storage_aws_role_arn").(string))
		}
		if data.HasChange("azure_tenant_id") {
			runSetStatement = true
			stmt.SetString("AZURE_TENANT_ID", data.Get("azure_tenant_id").(string))
		}
	}

	if runSetStatement {
		if err := snowflake.Exec(db, stmt.Statement()); err != nil {
			return fmt.Errorf("error updating storage integration: %w", err)
		}
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
	case "GCS":
		// nothing to set here
	default:
		return fmt.Errorf("Unexpected provider %v", storageProvider)
	}

	return nil
}
