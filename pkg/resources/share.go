package resources

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/validation"
)

var shareProperties = []string{
	"comment",
}

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

	builder := snowflake.NewShareBuilder(name).Create()
	builder.SetString("COMMENT", d.Get("comment").(string))

	if err := snowflake.Exec(db, builder.Statement()); err != nil {
		return fmt.Errorf("error creating share err = %w", err)
	}
	d.SetId(name)

	// Adding accounts must be done via an ALTER query

	// @TODO flesh out the share type in the snowflake package since it doesn't
	// follow the normal generic rules
	if err := setAccounts(d, meta); err != nil {
		return err
	}

	return ReadShare(d, meta)
}

func setAccounts(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	accs := expandStringList(d.Get("accounts").([]interface{}))

	// There is a race condition where error accounts cannot be added to a
	// share until after a database is added to the share. Since a database
	// grant is dependent on the share itself, this is a hack to get the
	// thing working.
	// 1. Create new temporary DB
	tempName := fmt.Sprintf("TEMP_%v_%d", name, time.Now().Unix())
	tempDB := snowflake.NewDatabaseBuilder(tempName)
	if err := snowflake.Exec(db, tempDB.Create()); err != nil {
		return fmt.Errorf("error creating temporary DB %v err = %w", tempName, err)
	}

	// 2. Create temporary DB grant to the share
	tempDBGrant := snowflake.DatabaseGrant(tempName)

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
	if err := snowflake.Exec(db, tempDBGrant.Share(name).Grant("REFERENCE_USAGE", false)); err != nil {
		return fmt.Errorf("error creating temporary DB REFERENCE_USAGE grant %v err = %w", tempName, err)
	}

	// 3. Add the accounts to the share
	if len(accs) > 0 {
		q := fmt.Sprintf(`ALTER SHARE "%v" SET ACCOUNTS=%v`, name, strings.Join(accs, ","))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error adding accounts to share %v err = %w", name, err)
		}
	} else {
		q := fmt.Sprintf(`ALTER SHARE "%v" UNSET ACCOUNTS`, name)
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error unsetting accounts to share %v err = %w", name, err)
		}
	}

	// 4. Revoke temporary DB grant to the share
	if err := snowflake.ExecMulti(db, tempDBGrant.Share(name).Revoke("REFERENCE_USAGE")); err != nil {
		return fmt.Errorf("error revoking temporary DB REFERENCE_USAGE grant %v err = %w", tempName, err)
	}

	// revoke the maybe automatically granted USAGE privilege.
	if err := snowflake.ExecMulti(db, tempDBGrant.Share(name).Revoke("USAGE")); err != nil {
		return fmt.Errorf("error revoking temporary DB grant %v err = %w", tempName, err)
	}

	// 5. Remove the temporary DB
	if err := snowflake.Exec(db, tempDB.Drop()); err != nil {
		return fmt.Errorf("error dropping temporary DB %v err = %w", tempName, err)
	}

	return nil
}

// ReadShare implements schema.ReadFunc.
func ReadShare(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	stmt := snowflake.NewShareBuilder(id).Show()
	row := snowflake.QueryRow(db, stmt)

	s, err := snowflake.ScanShare(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] share (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	if err := d.Set("name", StripAccountFromName(s.Name.String)); err != nil {
		return err
	}
	if err := d.Set("comment", s.Comment.String); err != nil {
		return err
	}

	accs := strings.FieldsFunc(s.To.String, func(c rune) bool { return c == ',' })
	err = d.Set("accounts", accs)

	return err
}

// UpdateShare implements schema.UpdateFunc.
func UpdateShare(d *schema.ResourceData, meta interface{}) error {
	// Change the accounts first - this is a special case and won't work using the generic method
	if d.HasChange("accounts") {
		if err := setAccounts(d, meta); err != nil {
			return err
		}
	}

	return UpdateResource("this does not seem to be used", shareProperties, shareSchema, snowflake.NewShareBuilder, ReadShare)(d, meta)
}

// DeleteShare implements schema.DeleteFunc.
func DeleteShare(d *schema.ResourceData, meta interface{}) error {
	return DeleteResource("this does not seem to be used", snowflake.NewShareBuilder)(d, meta)
}

// StripAccountFromName removes the account prefix from a resource (e.g. a share)
// that returns it (e.g. yt12345.my_share or org.acc.my_share should just be my_share).
func StripAccountFromName(s string) string {
	return s[strings.LastIndex(s, ".")+1:]
}
