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

var storageIntegrationSchema = map[string]*schema.Schema{
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
		Default:      "EXTERNAL_STAGE",
		ValidateFunc: validation.StringInSlice([]string{"EXTERNAL_STAGE"}, true),
		ForceNew:     true,
	},
	"enabled": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"storage_allowed_locations": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: "Explicitly limits external stages that use the integration to reference one or more storage locations.",
		MinItems:    1,
	},
	"storage_blocked_locations": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Explicitly prohibits external stages that use the integration from referencing one or more storage locations.",
	},
	// This part of the schema is the cloudProviderParams in the Snowflake documentation and differs between vendors
	"storage_provider": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"S3", "GCS", "AZURE"}, false),
	},
	"storage_aws_external_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The external ID that Snowflake will use when assuming the AWS role.",
	},
	"storage_aws_iam_user_arn": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowflake user that will attempt to assume the AWS role.",
	},
	"storage_aws_role_arn": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"azure_tenant_id": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"created_on": {
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

		Schema: storageIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateStorageIntegration implements schema.CreateFunc
func CreateStorageIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)

	stmt := snowflake.StorageIntegration(name).Create()

	// Set required fields
	stmt.SetString(`TYPE`, d.Get("type").(string))
	stmt.SetBool(`ENABLED`, d.Get("enabled").(bool))

	stmt.SetStringList("STORAGE_ALLOWED_LOCATIONS", expandStringList(d.Get("storage_allowed_locations").([]interface{})))

	// Set optional fields
	if v, ok := d.GetOk("comment"); ok {
		stmt.SetString(`COMMENT`, v.(string))
	}

	if _, ok := d.GetOk("storage_blocked_locations"); ok {
		stmt.SetStringList("STORAGE_BLOCKED_LOCATIONS", expandStringList(d.Get("storage_blocked_locations").([]interface{})))
	}

	// Now, set the storage provider
	err := setStorageProviderSettings(d, stmt)
	if err != nil {
		return err
	}

	err = snowflake.Exec(db, stmt.Statement())
	if err != nil {
		return fmt.Errorf("error creating storage integration: %w", err)
	}

	d.SetId(name)

	return ReadStorageIntegration(d, meta)
}

// ReadStorageIntegration implements schema.ReadFunc
func ReadStorageIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	stmt := snowflake.StorageIntegration(d.Id()).Show()
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
	stmt = snowflake.StorageIntegration(d.Id()).Describe()
	rows, err := db.Query(stmt)
	if err != nil {
		return fmt.Errorf("Could not describe storage integration: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&k, &pType, &v, &unused); err != nil {
			return err
		}
		switch k {
		case "ENABLED":
			// We set this using the SHOW INTEGRATION call so let's ignore it here
		case "STORAGE_PROVIDER":
			if err = d.Set("storage_provider", v.(string)); err != nil {
				return err
			}
		case "STORAGE_ALLOWED_LOCATIONS":
			if err = d.Set("storage_allowed_locations", strings.Split(v.(string), ",")); err != nil {
				return err
			}
		case "STORAGE_BLOCKED_LOCATIONS":
			if val := v.(string); val != "" {
				if err = d.Set("storage_blocked_locations", strings.Split(val, ",")); err != nil {
					return err
				}
			}
		case "STORAGE_AWS_IAM_USER_ARN":
			if err = d.Set("storage_aws_iam_user_arn", v.(string)); err != nil {
				return err
			}
		case "STORAGE_AWS_ROLE_ARN":
			if err = d.Set("storage_aws_role_arn", v.(string)); err != nil {
				return err
			}
		case "STORAGE_AWS_EXTERNAL_ID":
			if err = d.Set("storage_aws_external_id", v.(string)); err != nil {
				return err
			}
		default:
			log.Printf("[WARN] unexpected property %v returned from Snowflake", k)
		}
	}

	return err
}

// UpdateStorageIntegration implements schema.UpdateFunc
func UpdateStorageIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	stmt := snowflake.StorageIntegration(id).Alter()

	// This is required in case the only change is to UNSET STORAGE_ALLOWED_LOCATIONS.
	// Not sure if there is a more elegant way of determining this
	var runSetStatement bool

	if d.HasChange("comment") {
		runSetStatement = true
		stmt.SetString("COMMENT", d.Get("comment").(string))
	}

	if d.HasChange("enabled") {
		runSetStatement = true
		stmt.SetBool(`ENABLED`, d.Get("enabled").(bool))
	}

	if d.HasChange("storage_allowed_locations") {
		runSetStatement = true
		stmt.SetStringList("STORAGE_ALLOWED_LOCATIONS", expandStringList(d.Get("storage_allowed_locations").([]interface{})))
	}

	// We need to UNSET this if we remove all storage blocked locations. I don't think
	// this is documented by Snowflake, but this is how it works.
	//
	// @TODO move the SQL back to the snowflake package
	if d.HasChange("storage_blocked_locations") {
		v := d.Get("storage_blocked_locations").([]interface{})
		if len(v) == 0 {
			err := snowflake.Exec(db, fmt.Sprintf(`ALTER STORAGE INTEGRATION %v UNSET STORAGE_BLOCKED_LOCATIONS`, d.Id()))
			if err != nil {
				return fmt.Errorf("error unsetting storage_blocked_locations: %w", err)
			}
		} else {
			runSetStatement = true
			stmt.SetStringList("STORAGE_BLOCKED_LOCATIONS", expandStringList(v))
		}
	}

	if d.HasChange("storage_provider") {
		runSetStatement = true
		err := setStorageProviderSettings(d, stmt)
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
	}

	if runSetStatement {
		if err := snowflake.Exec(db, stmt.Statement()); err != nil {
			return fmt.Errorf("error updating storage integration: %w", err)
		}
	}

	return ReadStorageIntegration(d, meta)
}

// DeleteStorageIntegration implements schema.DeleteFunc
func DeleteStorageIntegration(d *schema.ResourceData, meta interface{}) error {
	return DeleteResource("", snowflake.StorageIntegration)(d, meta)
}

// StorageIntegrationExists implements schema.ExistsFunc
func StorageIntegrationExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	id := d.Id()

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
