package resources

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	snowflakeValidation "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/validation"
)

const (
	SnowflakeReaderAccountType = "READER"
)

var managedAccountProperties = []string{
	"admin_name",
	"admin_password",
	"type",
	"comment",
}

var managedAccountSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier for the managed account; must be unique for your account.",
		ForceNew:    true,
	},
	"admin_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier, as well as login name, for the initial user in the managed account. This user serves as the account administrator for the account.",
		ForceNew:    true,
	},
	"admin_password": &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Sensitive:    true,
		Description:  "Password for the initial user in the managed account.",
		ValidateFunc: snowflakeValidation.ValidatePassword,
		ForceNew:     true,
	},
	"type": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      SnowflakeReaderAccountType,
		Description:  "Specifies the type of managed account.",
		ValidateFunc: validation.StringInSlice([]string{SnowflakeReaderAccountType}, true),
		ForceNew:     true,
	},
	"comment": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the managed account.",
		ForceNew:    true,
	},
	"cloud": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Cloud in which the managed account is located.",
	},
	"region": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Snowflake Region in which the managed account is located.",
	},
	"locator": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Display name of the managed account.",
	},
	"created_on": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the managed account was created.",
	},
	"url": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "URL for accessing the managed account, particularly through the web interface.",
	},
}

// ManagedAccount returns a pointer to the resource representing a managed account
func ManagedAccount() *schema.Resource {
	return &schema.Resource{
		Create: CreateManagedAccount,
		Read:   ReadManagedAccount,
		Delete: DeleteManagedAccount,
		Exists: ManagedAccountExists,

		Schema: managedAccountSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateManagedAccount implements schema.CreateFunc
func CreateManagedAccount(data *schema.ResourceData, meta interface{}) error {
	return CreateResource(
		"this does not seem to be used",
		managedAccountProperties,
		managedAccountSchema,
		snowflake.ManagedAccount,
		initialReadManagedAccount,
	)(data, meta)
}

// initialReadManagedAccount is used for the first read, since the locator takes
// some time to appear. This is currently implemented as a sleep. @TODO actually
// wait until the locator is generated.
func initialReadManagedAccount(data *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] sleeping to give the locator a chance to be generated")
	time.Sleep(10 * time.Second)
	return ReadManagedAccount(data, meta)
}

// ReadManagedAccount implements schema.ReadFunc
func ReadManagedAccount(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.ManagedAccount(id).Show()
	row := snowflake.QueryRow(db, stmt)
	a, err := snowflake.ScanManagedAccount(row)
	if err != nil {
		return err
	}

	err = data.Set("name", a.Name.String)
	if err != nil {
		return err
	}
	err = data.Set("cloud", a.Cloud.String)
	if err != nil {
		return err
	}

	err = data.Set("region", a.Region.String)
	if err != nil {
		return err
	}

	err = data.Set("locator", a.Locator.String)
	if err != nil {
		return err
	}

	err = data.Set("created_on", a.CreatedOn.String)
	if err != nil {
		return err
	}

	err = data.Set("url", a.Url.String)
	if err != nil {
		return err
	}

	if a.IsReader {
		err = data.Set("type", "READER")
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Unable to determine the account type")
	}

	err = data.Set("comment", a.Comment.String)

	return err
}

// DeleteManagedAccount implements schema.DeleteFunc
func DeleteManagedAccount(data *schema.ResourceData, meta interface{}) error {
	return DeleteResource("this does not seem to be used", snowflake.ManagedAccount)(data, meta)
}

// ManagedAccountExists implements schema.ExistsFunc
func ManagedAccountExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.ManagedAccount(id).Show()
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
