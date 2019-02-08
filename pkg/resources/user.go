package resources

import (
	"database/sql"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

var userProperties = []string{"comment", "password"}

func User() *schema.Resource {
	return &schema.Resource{
		Create: CreateUser,
		Read:   ReadUser,
		Delete: DeleteUser,
		Update: UpdateUser,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the user. Note that if you do not supply login_name this will be used as login_name. [doc](https://docs.snowflake.net/manuals/sql-reference/sql/create-user.html#required-parameters)"},
			// "login_name": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Required: false,
			// 	Computed: true,
			// },
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				// TODO validation
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "**WARNING:** this will put the password in the terraform state file. Use carefully.",
				// TODO validation https://docs.snowflake.net/manuals/sql-reference/sql/create-user.html#optional-parameters
			},
			//    LOGIN_NAME = <string>
			//    DISPLAY_NAME = <string>
			//    FIRST_NAME = <string>
			//    MIDDLE_NAME = <string>
			//    LAST_NAME = <string>
			//    EMAIL = <string>
			//    MUST_CHANGE_PASSWORD = TRUE | FALSE
			//    DISABLED = TRUE | FALSE
			//    SNOWFLAKE_LOCK = TRUE | FALSE
			//    SNOWFLAKE_SUPPORT = TRUE | FALSE
			//    DAYS_TO_EXPIRY = <integer>
			//    MINS_TO_UNLOCK = <integer>
			//    DEFAULT_WAREHOUSE = <string>
			//    DEFAULT_NAMESPACE = <string>
			//    DEFAULT_ROLE = <string>
			//    EXT_AUTHN_DUO = TRUE | FALSE
			//    EXT_AUTHN_UID = <string>
			//    MINS_TO_BYPASS_MFA = <integer>
			//    DISABLE_MFA = TRUE | FALSE
			//    MINS_TO_BYPASS_NETWORK POLICY = <integer>
			//    RSA_PUBLIC_KEY = <string>
			//    RSA_PUBLIC_KEY_2 = <string>
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func CreateUser(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)

	qb := snowflake.User(name).Create()

	for _, field := range userProperties {
		log.Printf("prop %s", field)
		val, ok := data.GetOk(field)
		log.Printf("val, ok %#v, %#v", ok, val)
		if ok {
			valStr := val.(string)
			qb.SetString(field, valStr)
		}
	}
	err := DBExec(db, qb.Statement())

	if err != nil {
		return errors.Wrap(err, "error creating user")
	}

	data.SetId(name)

	return ReadUser(data, meta)
}

func ReadUser(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.User(id).Show()
	row := db.QueryRow(stmt)
	var name, createdOn, loginName, displayName, firstName, lastName, email, minsToUnlock, daysToExpiry, comment, disabled, mustChangePassword, snowflakeLock, defaultWarehouse, defaultNamespace, defaultRole, extAuthnDuo, extAuthnUID, minsToBypassMfa, owner, lastSuccessLogin, expiresAtTime, lockedUntilTime, hasPassword, hasRsaPublicKey sql.NullString
	err := row.Scan(&name, &createdOn, &loginName, &displayName, &firstName, &lastName, &email, &minsToUnlock, &daysToExpiry, &comment, &disabled, &mustChangePassword, &snowflakeLock, &defaultWarehouse, &defaultNamespace, &defaultRole, &extAuthnDuo, &extAuthnUID, &minsToBypassMfa, &owner, &lastSuccessLogin, &expiresAtTime, &lockedUntilTime, &hasPassword, &hasRsaPublicKey)
	if err != nil {
		return err
	}

	err = data.Set("name", name.String)
	if err != nil {
		return err
	}
	err = data.Set("comment", comment.String)
	if err != nil {
		return err
	}

	return err
}

func DeleteUser(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Id()

	stmt := snowflake.User(name).Drop()
	err := DBExec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error dropping user %s", name)
	}

	data.SetId("")
	return nil
}

func UpdateUser(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	if data.HasChange("name") {
		data.Partial(true)
		// I wish this could be done on one line.
		oldNameI, newNameI := data.GetChange("name")
		oldName := oldNameI.(string)
		newName := newNameI.(string)

		stmt := snowflake.User(oldName).Rename(newName)
		err := DBExec(db, stmt)

		if err != nil {
			return errors.Wrapf(err, "error renaming user %s to %s", oldName, newName)
		}
		data.SetId(newName)
		data.SetPartial("name")
		data.Partial(false)
	}

	changes := []string{}

	for _, prop := range userProperties {
		if data.HasChange(prop) {
			changes = append(changes, prop)
		}
	}
	if len(changes) > 0 {
		name := data.Get("name").(string)
		qb := snowflake.User(name).Alter()

		for _, change := range changes {
			val := data.Get(change).(string)
			qb.SetString(change, val)
		}

		err := DBExec(db, qb.Statement())
		if err != nil {
			return errors.Wrap(err, "error altering user")
		}
	}
	return ReadUser(data, meta)
}
