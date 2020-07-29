package resources

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var shareProperties = []string{
	"comment",
}

var shareSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the share; must be unique for the account in which the share is created.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the managed account.",
	},
	"accounts": {
		// Changed from Set to List to use DiffSuppressFunc: https://github.com/hashicorp/terraform-plugin-sdk/issues/160
		Type:             schema.TypeList,
		Elem:             &schema.Schema{Type: schema.TypeString},
		Optional:         true,
		Description:      "A list of accounts to be added to the share.",
		DiffSuppressFunc: diffCaseInsensitive,
	},
}

// Share returns a pointer to the resource representing a share
func Share() *schema.Resource {
	return &schema.Resource{
		Create: CreateShare,
		Read:   ReadShare,
		Update: UpdateShare,
		Delete: DeleteShare,
		Exists: ShareExists,

		Schema: shareSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateShare implements schema.CreateFunc
func CreateShare(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)

	builder := snowflake.Share(name).Create()
	builder.SetString("COMMENT", data.Get("comment").(string))

	err := snowflake.Exec(db, builder.Statement())
	if err != nil {
		return errors.Wrapf(err, "error creating share")
	}
	data.SetId(name)

	// Adding accounts must be done via an ALTER query

	// @TODO flesh out the share type in the snowflake package since it doesn't
	// follow the normal generic rules
	err = setAccounts(data, meta)
	if err != nil {
		return err
	}

	return ReadShare(data, meta)
}

func setAccounts(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	accs := expandStringList(data.Get("accounts").([]interface{}))

	if len(accs) > 0 {
		// There is a race condition where error accounts cannot be added to a
		// share until after a database is added to the share. Since a database
		// grant is dependent on the share itself, this is a hack to get the
		// thing working.
		// 1. Create new temporary DB
		tempName := fmt.Sprintf("TEMP_%v_%d", name, time.Now().Unix())
		tempDB := snowflake.Database(tempName)
		err := snowflake.Exec(db, tempDB.Create().Statement())
		if err != nil {
			return errors.Wrapf(err, "error creating temporary DB %v", tempName)
		}

		// 2. Create temporary DB grant to the share
		tempDBGrant := snowflake.DatabaseGrant(tempName)
		err = snowflake.Exec(db, tempDBGrant.Share(name).Grant("USAGE"))
		if err != nil {
			return errors.Wrapf(err, "error creating temporary DB grant %v", tempName)
		}
		// 3. Add the accounts to the share
		q := fmt.Sprintf(`ALTER SHARE "%v" SET ACCOUNTS=%v`, name, strings.Join(accs, ","))
		err = snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error adding accounts to share %v", name)
		}
		// 4. Revoke temporary DB grant to the share
		err = snowflake.Exec(db, tempDBGrant.Share(name).Revoke("USAGE"))
		if err != nil {
			return errors.Wrapf(err, "error revoking temporary DB grant %v", tempName)
		}
		// 5. Remove the temporary DB
		err = snowflake.Exec(db, tempDB.Drop())
		if err != nil {
			return errors.Wrapf(err, "error dropping temporary DB %v", tempName)
		}
	}

	return nil
}

// ReadShare implements schema.ReadFunc
func ReadShare(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.Share(id).Show()
	row := snowflake.QueryRow(db, stmt)

	s, err := snowflake.ScanShare(row)
	if err != nil {
		return err
	}

	err = data.Set("name", StripAccountFromName(s.Name.String))
	if err != nil {
		return err
	}
	err = data.Set("comment", s.Comment.String)
	if err != nil {
		return err
	}

	accs := strings.FieldsFunc(s.To.String, func(c rune) bool { return c == ',' })
	err = data.Set("accounts", accs)

	return err
}

// UpdateShare implements schema.UpdateFunc
func UpdateShare(data *schema.ResourceData, meta interface{}) error {
	// Change the accounts first - this is a special case and won't work using the generic method
	if data.HasChange("accounts") {
		err := setAccounts(data, meta)
		if err != nil {
			return err
		}
	}

	return UpdateResource("this does not seem to be used", shareProperties, shareSchema, snowflake.Share, ReadShare)(data, meta)
}

// DeleteShare implements schema.DeleteFunc
func DeleteShare(data *schema.ResourceData, meta interface{}) error {
	return DeleteResource("this does not seem to be used", snowflake.Share)(data, meta)
}

// ShareExists implements schema.ExistsFunc
func ShareExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.Share(id).Show()
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

// StripAccountFromName removes the accout prefix from a resource (e.g. a share)
// that returns it (e.g. yt12345.my_share should just be my_share)
func StripAccountFromName(s string) string {
	return s[strings.Index(s, ".")+1:]
}
