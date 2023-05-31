package resources

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/validation"
)

var shareSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the share; must be unique for the account in which the share is created.",
		ForceNew:    true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the managed account.",
	},
	"accounts": {
		// Changed from Set to List to use DiffSuppressFunc: https://github.com/hashicorp/terraform-plugin-sdk/issues/160
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.ValidateIsNotAccountLocator,
		},
		Optional: true,
		Description: "A list of accounts to be added to the share. Values should not be the account locator, but " +
			"in the form of 'organization_name.account_name",
		DiffSuppressFunc: diffCaseInsensitive,
	},
}

// Share returns a pointer to the resource representing a share.
func Share() *schema.Resource {
	return &schema.Resource{
		Create: CreateShare,
		Read:   ReadShare,
		Update: UpdateShare,
		Delete: DeleteShare,

		Schema: shareSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateShare implements schema.CreateFunc.
func CreateShare(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)
	comment := d.Get("comment").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	var opts sdk.CreateShareOptions
	if comment != "" {
		opts = sdk.CreateShareOptions{
			Comment: sdk.String(comment),
		}
	}
	if err := client.Shares.Create(ctx, id, &opts); err != nil {
		return fmt.Errorf("error creating share err = %w", err)
	}
	d.SetId(name)

	accounts := expandStringList(d.Get("accounts").([]interface{}))
	if len(accounts) > 0 {
		shareID := sdk.NewAccountObjectIdentifier(name)
		accountIdentifiers := make([]sdk.AccountIdentifier, len(accounts))
		for i, account := range accounts {
			parts := strings.Split(account, ".")
			orgName := parts[0]
			accountName := parts[1]
			accountIdentifiers[i] = sdk.NewAccountIdentifier(orgName, accountName)
		}
		err := setShareAccounts(ctx, client, shareID, accountIdentifiers)
		if err != nil {
			return err
		}
	}
	return ReadShare(d, meta)
}

func setShareAccounts(ctx context.Context, client *sdk.Client, shareID sdk.AccountObjectIdentifier, accounts []sdk.AccountIdentifier) error {
	// There is a race condition where error accounts cannot be added to a
	// share until after a database is added to the share. Since a database
	// grant is dependent on the share itself, this is a hack to get the
	// thing working.

	// 1. Create new temporary DB
	tempName := fmt.Sprintf("TEMP_%v_%d", shareID.Name(), time.Now().Unix())
	tempDatabaseID := sdk.NewAccountObjectIdentifier(tempName)
	err := client.Databases.Create(ctx, tempDatabaseID, nil)
	if err != nil {
		return fmt.Errorf("error creating temporary DB %v err = %w", tempName, err)
	}
	defer func() {
		// drop the temporary DB during cleanup
		err = client.Databases.Drop(ctx, tempDatabaseID, nil)
		if err != nil {
			log.Printf("[WARN] error dropping temporary DB %v err = %v", tempName, err)
		}
	}()
	// 2. Create temporary DB grant to the share
	// USAGE can only be granted to one database - granting USAGE on the temp db here
	// conflicts (and errors) with having a database already shared (i.e. when you
	// already have a share and are just adding or removing accounts). Instead, use
	// REFERENCE_USAGE which is intended for multi-database sharing as per Snowflake
	// documentation here:
	// https://docs.snowflake.com/en/sql-reference/sql/grant-privilege-share.html#usage-notes
	// Note however that USAGE will be granted automatically on the temp db for the
	// case where the main db doesn't already exist, so it will need to be revoked
	// before deleting the temp db. Where USAGE hasn't been already granted it is not
	// an error to revoke it, so it's ok to just do the revoke every time.
	err = client.Grants.GrantPrivilegeToShare(ctx, sdk.PrivilegeReferenceUsage, &sdk.GrantPrivilegeToShareOn{
		Database: tempDatabaseID,
	}, shareID)
	if err != nil {
		return fmt.Errorf("error granting privilege to share err = %w", err)
	}
	defer func() {
		// revoke the REFERENCE_USAGE privilege during cleanup
		err = client.Grants.RevokePrivilegeFromShare(ctx, sdk.PrivilegeReferenceUsage, &sdk.RevokePrivilegeFromShareOn{
			Database: tempDatabaseID,
		}, shareID)
		if err != nil {
			log.Printf("[WARN] error revoking privilege from share err = %v", err)
		}
		// revoke the maybe automatically granted USAGE privilege during cleanup
		err = client.Grants.RevokePrivilegeFromShare(ctx, sdk.PrivilegeUsage, &sdk.RevokePrivilegeFromShareOn{
			Database: tempDatabaseID,
		}, shareID)
		if err != nil {
			log.Printf("[WARN] error revoking privilege from share err = %v", err)
		}
	}()
	// 3. Add accounts to the share
	err = client.Shares.Alter(ctx, shareID, &sdk.AlterShareOptions{
		Add: &sdk.ShareAdd{
			Accounts: accounts,
		},
	})
	return err
}

// ReadShare implements schema.ReadFunc.
func ReadShare(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := sdk.NewAccountObjectIdentifier(d.Id())
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	share, err := client.Shares.ShowByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error reading share err = %w", err)
	}
	if err := d.Set("name", share.Name.Name()); err != nil {
		return err
	}
	if err := d.Set("comment", share.Comment); err != nil {
		return err
	}
	accounts := make([]string, len(share.To))
	for i, accountIdentifier := range share.To {
		accounts[i] = accountIdentifier.Name()
	}

	currentAccount := d.Get("accounts")
	if currentAccount != nil {
		currentAccounts := expandStringList(currentAccount.([]interface{}))
		// reorder the accounts so they match the order in the config
		// this is to avoid unnecessary diffs
		accounts = reorderStringList(currentAccounts, accounts)
	}
	if err := d.Set("accounts", accounts); err != nil {
		return err
	}

	return err
}

func accountIdentifiersFromSlice(accounts []string) []sdk.AccountIdentifier {
	accountIdentifiers := make([]sdk.AccountIdentifier, len(accounts))
	for i, account := range accounts {
		parts := strings.Split(account, ".")
		orgName := parts[0]
		accountName := parts[1]
		accountIdentifiers[i] = sdk.NewAccountIdentifier(orgName, accountName)
	}
	return accountIdentifiers
}

// UpdateShare implements schema.UpdateFunc.
func UpdateShare(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	if d.HasChange("accounts") {
		o, n := d.GetChange("accounts")
		oldAccounts := expandStringList(o.([]interface{}))
		newAccounts := expandStringList(n.([]interface{}))
		if len(newAccounts) == 0 {
			accountIdentifiers := accountIdentifiersFromSlice(oldAccounts)
			err := client.Shares.Alter(ctx, sdk.NewAccountObjectIdentifier(d.Id()), &sdk.AlterShareOptions{
				Remove: &sdk.ShareRemove{
					Accounts: accountIdentifiers,
				},
			})
			if err != nil {
				return fmt.Errorf("error removing accounts from share err = %w", err)
			}
		} else {
			accountIdentifiers := accountIdentifiersFromSlice(newAccounts)
			err := setShareAccounts(ctx, client, sdk.NewAccountObjectIdentifier(d.Id()), accountIdentifiers)
			if err != nil {
				return err
			}
		}
	}
	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		err := client.Shares.Alter(ctx, sdk.NewAccountObjectIdentifier(d.Id()), &sdk.AlterShareOptions{
			Set: &sdk.ShareSet{
				Comment: sdk.String(comment),
			},
		})
		if err != nil {
			return fmt.Errorf("error updating share comment err = %w", err)
		}
	}

	return ReadShare(d, meta)
}

// DeleteShare implements schema.DeleteFunc.
func DeleteShare(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	err := client.Shares.Drop(ctx, sdk.NewAccountObjectIdentifier(d.Id()))
	if err != nil {
		return fmt.Errorf("error deleting share err = %w", err)
	}
	return nil
}
