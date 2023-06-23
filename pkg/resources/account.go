package resources

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	snowflakeValidation "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Note: no test case was created for account since we cannot actually delete them after creation, which is a critical part of the test suite. Instead, this resource
// was manually tested

var accountSchema = map[string]*schema.Schema{
	"name": {
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Specifies the identifier (i.e. name) for the account; must be unique within an organization, regardless of which Snowflake Region the account is in. In addition, the identifier must start with an alphabetic character and cannot contain spaces or special characters except for underscores (_). Note that if the account name includes underscores, features that do not accept account names with underscores (e.g. Okta SSO or SCIM) can reference a version of the account name that substitutes hyphens (-) for the underscores.",
		ValidateFunc: snowflakeValidation.ValidateAccountIdentifier,
		// Name is automatically uppercase by Snowflake
		StateFunc: func(val interface{}) string {
			return strings.ToUpper(val.(string))
		},
	},
	"admin_name": {
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Login name of the initial administrative user of the account. A new user is created in the new account with this name and password and granted the ACCOUNTADMIN role in the account. A login name can be any string consisting of letters, numbers, and underscores. Login names are always case-insensitive.",
		ValidateFunc: snowflakeValidation.ValidateAdminName,
		// We have no way of assuming a role into this account to change the admin user name so this has to be ForceNew even though it's not ideal
		ForceNew:              true,
		DiffSuppressOnRefresh: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// For new resources always show the diff
			if d.Id() == "" {
				return false
			}
			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
			return old == ""
		},
	},
	"admin_password": {
		Type:         schema.TypeString,
		Optional:     true,
		Sensitive:    true,
		Description:  "Password for the initial administrative user of the account. Optional if the `ADMIN_RSA_PUBLIC_KEY` parameter is specified. For more information about passwords in Snowflake, see [Snowflake-provided Password Policy](https://docs.snowflake.com/en/sql-reference/sql/create-account.html#:~:text=Snowflake%2Dprovided%20Password%20Policy).",
		AtLeastOneOf: []string{"admin_password", "admin_rsa_public_key"},
		// We have no way of assuming a role into this account to change the password so this has to be ForceNew even though it's not ideal
		ForceNew:              true,
		DiffSuppressOnRefresh: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// For new resources always show the diff
			if d.Id() == "" {
				return false
			}
			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
			return old == ""
		},
	},
	"admin_rsa_public_key": {
		Type:         schema.TypeString,
		Optional:     true,
		Sensitive:    true,
		Description:  "Assigns a public key to the initial administrative user of the account in order to implement [key pair authentication](https://docs.snowflake.com/en/sql-reference/sql/create-account.html#:~:text=key%20pair%20authentication) for the user. Optional if the `ADMIN_PASSWORD` parameter is specified.",
		AtLeastOneOf: []string{"admin_password", "admin_rsa_public_key"},
		// We have no way of assuming a role into this account to change the admin rsa public key so this has to be ForceNew even though it's not ideal
		ForceNew:              true,
		DiffSuppressOnRefresh: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// For new resources always show the diff
			if d.Id() == "" {
				return false
			}
			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
			return old == ""
		},
	},
	"email": {
		Type:         schema.TypeString,
		Required:     true,
		Sensitive:    true,
		Description:  "Email address of the initial administrative user of the account. This email address is used to send any notifications about the account.",
		ValidateFunc: snowflakeValidation.ValidateEmail,
		// We have no way of assuming a role into this account to change the admin email so this has to be ForceNew even though it's not ideal
		ForceNew:              true,
		DiffSuppressOnRefresh: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// For new resources always show the diff
			if d.Id() == "" {
				return false
			}
			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
			return old == ""
		},
	},
	"edition": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		Description:  "[Snowflake Edition](https://docs.snowflake.com/en/user-guide/intro-editions.html) of the account. Valid values are: STANDARD | ENTERPRISE | BUSINESS_CRITICAL",
		ValidateFunc: validation.StringInSlice([]string{string(sdk.EditionStandard), string(sdk.EditionEnterprise), string(sdk.EditionBusinessCritical)}, false),
	},
	"first_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "First name of the initial administrative user of the account",
		// We have no way of assuming a role into this account to change the admin first name so this has to be ForceNew even though it's not ideal
		ForceNew:              true,
		DiffSuppressOnRefresh: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// For new resources always show the diff
			if d.Id() == "" {
				return false
			}
			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
			return old == ""
		},
	},
	"last_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "Last name of the initial administrative user of the account",
		// We have no way of assuming a role into this account to change the admin last name so this has to be ForceNew even though it's not ideal
		ForceNew:              true,
		DiffSuppressOnRefresh: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// For new resources always show the diff
			if d.Id() == "" {
				return false
			}
			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
			return old == ""
		},
	},
	"must_change_password": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies whether the new user created to administer the account is forced to change their password upon first login into the account.",
		// We have no way of assuming a role into this account to change the admin password policy so this has to be ForceNew even though it's not ideal
		ForceNew:              true,
		DiffSuppressOnRefresh: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// For new resources always show the diff
			if d.Id() == "" {
				return false
			}
			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
			return old == ""
		},
	},
	"region_group": {
		Type:                  schema.TypeString,
		Optional:              true,
		Description:           "ID of the Snowflake Region where the account is created. If no value is provided, Snowflake creates the account in the same Snowflake Region as the current account (i.e. the account in which the CREATE ACCOUNT statement is executed.)",
		ForceNew:              true,
		DiffSuppressOnRefresh: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// For new resources always show the diff
			if d.Id() == "" {
				return false
			}
			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
			return new == ""
		},
	},
	"region": {
		Type:                  schema.TypeString,
		Optional:              true,
		Description:           "ID of the Snowflake Region where the account is created. If no value is provided, Snowflake creates the account in the same Snowflake Region as the current account (i.e. the account in which the CREATE ACCOUNT statement is executed.)",
		ForceNew:              true,
		DiffSuppressOnRefresh: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// For new resources always show the diff
			if d.Id() == "" {
				return false
			}
			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
			return new == ""
		},
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the account.",
		ForceNew:    true,
	},
	"is_org_admin": {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Indicates whether the ORGADMIN role is enabled in an account. If TRUE, the role is enabled.",
	},
}

func Account() *schema.Resource {
	return &schema.Resource{
		Description: "The account resource allows you to create and manage Snowflake accounts.",
		Create:      CreateAccount,
		Read:        ReadAccount,
		Update:      UpdateAccount,
		Delete:      DeleteAccount,

		Schema: accountSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateAccount implements schema.CreateFunc.
func CreateAccount(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	name := d.Get("name").(string)
	objectIdentifier := sdk.NewAccountObjectIdentifier(name)

	createOptions := &sdk.CreateAccountOptions{
		AdminName: d.Get("admin_name").(string),
		Email:     d.Get("email").(string),
		Edition:   sdk.AccountEdition(d.Get("edition").(string)),
	}

	// get optional fields.
	if v, ok := d.GetOk("admin_password"); ok {
		createOptions.AdminPassword = sdk.String(v.(string))
	}
	if v, ok := d.GetOk("admin_rsa_public_key"); ok {
		createOptions.AdminRSAPublicKey = sdk.String(v.(string))
	}
	if v, ok := d.GetOk("first_name"); ok {
		createOptions.FirstName = sdk.String(v.(string))
	}
	if v, ok := d.GetOk("last_name"); ok {
		createOptions.LastName = sdk.String(v.(string))
	}
	if v, ok := d.GetOk("must_change_password"); ok {
		createOptions.MustChangePassword = sdk.Bool(v.(bool))
	}
	if v, ok := d.GetOk("region_group"); ok {
		createOptions.RegionGroup = sdk.String(v.(string))
	} else {
		// For organizations that have accounts in multiple region groups, returns <region_group>.<region> so we need to split on "."
		currentRegion, err := client.ContextFunctions.CurrentRegion(ctx)
		if err != nil {
			return err
		}
		regionParts := strings.Split(currentRegion, ".")
		if len(regionParts) == 2 {
			createOptions.RegionGroup = sdk.String(regionParts[0])
		}
	}
	if v, ok := d.GetOk("region"); ok {
		createOptions.Region = sdk.String(v.(string))
	} else {
		// For organizations that have accounts in multiple region groups, returns <region_group>.<region> so we need to split on "."
		currentRegion, err := client.ContextFunctions.CurrentRegion(ctx)
		if err != nil {
			return err
		}
		regionParts := strings.Split(currentRegion, ".")
		if len(regionParts) == 2 {
			createOptions.Region = sdk.String(regionParts[1])
		} else {
			createOptions.Region = sdk.String(currentRegion)
		}
	}
	if v, ok := d.GetOk("comment"); ok {
		createOptions.Comment = sdk.String(v.(string))
	}

	err := client.Accounts.Create(ctx, objectIdentifier, createOptions)
	if err != nil {
		return err
	}

	account, err := client.Accounts.ShowByID(ctx, objectIdentifier)
	if err != nil {
		return err
	}

	d.SetId(helpers.EncodeSnowflakeID(account.AccountLocator))
	return nil
}

// ReadAccount implements schema.ReadFunc.
func ReadAccount(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	acc, err := client.Accounts.ShowByID(ctx, id)
	if err != nil {
		return err
	}

	if err = d.Set("name", acc.AccountName); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}

	if err = d.Set("edition", acc.Edition); err != nil {
		return fmt.Errorf("error setting edition: %w", err)
	}

	if err = d.Set("region_group", acc.RegionGroup); err != nil {
		return fmt.Errorf("error setting region_group: %w", err)
	}

	if err = d.Set("region", acc.SnowflakeRegion); err != nil {
		return fmt.Errorf("error setting region: %w", err)
	}

	if err = d.Set("comment", acc.Comment); err != nil {
		return fmt.Errorf("error setting comment: %w", err)
	}

	if err = d.Set("is_org_admin", acc.IsOrgAdmin); err != nil {
		return fmt.Errorf("error setting is_org_admin: %w", err)
	}

	return nil
}

// UpdateAccount implements schema.UpdateFunc.
func UpdateAccount(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	// Rename
	if d.HasChange("name") {
		newID := sdk.NewAccountObjectIdentifier(d.Get("name").(string))
		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Rename: &sdk.AccountRename{
				Name:    id,
				NewName: newID,
			},
		})
		if err != nil {
			return err
		}
		d.SetId(helpers.EncodeSnowflakeID(newID))
	}

	// Change comment
	if d.HasChange("comment") {
		err := client.Comments.Set(ctx, &sdk.SetCommentOptions{
			ObjectType: sdk.ObjectTypeAccount,
			ObjectName: id,
			Value:      sdk.String(d.Get("comment").(string)),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteAccount implements schema.DeleteFunc.
func DeleteAccount(_ *schema.ResourceData, _ interface{}) error {
	return fmt.Errorf("cannot delete Snowflake accounts because there is no self service API allowing Terraform to do so. To delete an account, contact Snowflake Support and provide a unique identifier for your account, which can be one of the following:\n  Account name\n  Account locator\nOnce you contact Snowflake Support, it may take up to six weeks for the account to be fully deleted. This delay allows you to recover the account within 30 days of the request. Snowflake usually deducts the account from the number of accounts allowed for your organization within a few days of the initial request")
}
