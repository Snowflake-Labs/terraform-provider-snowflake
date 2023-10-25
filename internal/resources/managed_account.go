// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/snowflake"
	snowflakeValidation "github.com/Snowflake-Labs/terraform-provider-snowflake/internal/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	return CreateResource(
		"this does not seem to be used",
		managedAccountProperties,
		managedAccountSchema,
		snowflake.NewManagedAccountBuilder,
		initialReadManagedAccount,
	)(d, meta)
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
	id := d.Id()

	stmt := snowflake.NewManagedAccountBuilder(id).Show()
	row := snowflake.QueryRow(db, stmt)
	a, err := snowflake.ScanManagedAccount(row)

	if errors.Is(err, sql.ErrNoRows) {
		// If not found, remove resource from
		log.Printf("[DEBUG] managed account (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	if err := d.Set("name", a.Name.String); err != nil {
		return err
	}

	if err := d.Set("cloud", a.Cloud.String); err != nil {
		return err
	}

	if err := d.Set("region", a.Region.String); err != nil {
		return err
	}

	if err := d.Set("locator", a.Locator.String); err != nil {
		return err
	}

	if err := d.Set("created_on", a.CreatedOn.String); err != nil {
		return err
	}

	if err := d.Set("url", a.URL.String); err != nil {
		return err
	}

	if a.IsReader {
		if err := d.Set("type", "READER"); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unable to determine the account type")
	}

	err = d.Set("comment", a.Comment.String)

	return err
}

// DeleteManagedAccount implements schema.DeleteFunc.
func DeleteManagedAccount(d *schema.ResourceData, meta interface{}) error {
	return DeleteResource("this does not seem to be used", snowflake.NewManagedAccountBuilder)(d, meta)
}
