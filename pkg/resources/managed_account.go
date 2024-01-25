package resources

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	snowflakeValidation "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/validation"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	SnowflakeReaderAccountType = "READER"
)

var managedAccountSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier for the managed account; must be unique for your account.",
		ForceNew:    true,
	},
	"admin_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier, as well as login name, for the initial user in the managed account. This user serves as the account administrator for the account.",
		ForceNew:    true,
	},
	"admin_password": {
		Type:         schema.TypeString,
		Required:     true,
		Sensitive:    true,
		Description:  "Password for the initial user in the managed account.",
		ValidateFunc: snowflakeValidation.ValidatePassword,
		ForceNew:     true,
	},
	"type": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      SnowflakeReaderAccountType,
		Description:  "Specifies the type of managed account.",
		ValidateFunc: validation.StringInSlice([]string{SnowflakeReaderAccountType}, true),
		ForceNew:     true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the managed account.",
		ForceNew:    true,
	},
	"cloud": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Cloud in which the managed account is located.",
	},
	"region": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Snowflake Region in which the managed account is located.",
	},
	"locator": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Display name of the managed account.",
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the managed account was created.",
	},
	"url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "URL for accessing the managed account, particularly through the web interface.",
	},
}

// ManagedAccount returns a pointer to the resource representing a managed account.
func ManagedAccount() *schema.Resource {
	return &schema.Resource{
		Create: CreateManagedAccount,
		Read:   ReadManagedAccount,
		Delete: DeleteManagedAccount,

		Schema: managedAccountSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateManagedAccount implements schema.CreateFunc.
func CreateManagedAccount(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)

	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)

	adminName := d.Get("admin_name").(string)
	adminPassword := d.Get("admin_password").(string)
	createParams := sdk.NewCreateManagedAccountParamsRequest(adminName, adminPassword)

	if v, ok := d.GetOk("comment"); ok {
		createParams.WithComment(sdk.String(v.(string)))
	}

	createRequest := sdk.NewCreateManagedAccountRequest(id, *createParams)

	err := client.ManagedAccounts.Create(ctx, createRequest)
	if err != nil {
		return err
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadManagedAccount(d, meta)
}

// initialReadManagedAccount is used for the first read, since the locator takes
// some time to appear. This is currently implemented as a sleep. @TODO actually
// wait until the locator is generated.
func initialReadManagedAccount(d *schema.ResourceData, meta interface{}) error {
	log.Println("[INFO] sleeping to give the locator a chance to be generated")
	// lintignore:R018
	time.Sleep(10 * time.Second)
	return ReadManagedAccount(d, meta)
}

// ReadManagedAccount implements schema.ReadFunc.
func ReadManagedAccount(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	ctx := context.Background()
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	managedAccount, err := client.ManagedAccounts.ShowByID(ctx, objectIdentifier)
	if err != nil {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] managed account (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	if err := d.Set("name", managedAccount.Name); err != nil {
		return err
	}

	if err := d.Set("cloud", managedAccount.Cloud); err != nil {
		return err
	}

	if err := d.Set("region", managedAccount.Region); err != nil {
		return err
	}

	if err := d.Set("locator", managedAccount.Locator); err != nil {
		return err
	}

	if err := d.Set("created_on", managedAccount.CreatedOn); err != nil {
		return err
	}

	if err := d.Set("url", managedAccount.URL); err != nil {
		return err
	}

	if managedAccount.IsReader {
		if err := d.Set("type", "READER"); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unable to determine the account type")
	}

	if err := d.Set("comment", managedAccount.Comment); err != nil {
		return err
	}

	return nil
}

// DeleteManagedAccount implements schema.DeleteFunc.
func DeleteManagedAccount(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	err := client.ManagedAccounts.Drop(ctx, sdk.NewDropManagedAccountRequest(objectIdentifier))
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
