package resources

import (
	"database/sql"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

var userProperties = []string{
	"comment",
	"login_name",
	"password",
	"disabled",
	"default_namespace",
	"default_role",
	"default_warehouse",
	"rsa_public_key",
	"rsa_public_key_2",
}

var diffCaseInsensitive = func(k, old, new string, d *schema.ResourceData) bool {
	return strings.ToUpper(old) == strings.ToUpper(new)
}

var userSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of the user. Note that if you do not supply login_name this will be used as login_name. [doc](https://docs.snowflake.net/manuals/sql-reference/sql/create-user.html#required-parameters)"},
	"login_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "The name users use to log in. If not supplied, snowflake will use name instead.",
		// login_name is case-insensitive
		DiffSuppressFunc: diffCaseInsensitive,
	},
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
	"disabled": &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true,
	},
	"default_warehouse": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the virtual warehouse that is active by default for the user’s session upon login.",
	},
	"default_namespace": &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: diffCaseInsensitive,
		Description:      "Specifies the namespace (database only or database and schema) that is active by default for the user’s session upon login.",
	},
	"default_role": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Specifies the role that is active by default for the user’s session upon login.",
	},
	"rsa_public_key": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the user’s RSA public key; used for key-pair authentication. Must be on 1 line without header and trailer.",
	},
	"rsa_public_key_2": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the user’s second RSA public key; used to rotate the public and private keys for key-pair authentication based on an expiration schedule set by your organization. Must be on 1 line without header and trailer.",
	},
	"has_rsa_public_key": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Will be true if user as an RSA key set.",
	},

	//    DISPLAY_NAME = <string>
	//    FIRST_NAME = <string>
	//    MIDDLE_NAME = <string>
	//    LAST_NAME = <string>
	//    EMAIL = <string>
	//    MUST_CHANGE_PASSWORD = TRUE | FALSE
	//    SNOWFLAKE_LOCK = TRUE | FALSE
	//    SNOWFLAKE_SUPPORT = TRUE | FALSE
	//    DAYS_TO_EXPIRY = <integer>
	//    MINS_TO_UNLOCK = <integer>
	//    EXT_AUTHN_DUO = TRUE | FALSE
	//    EXT_AUTHN_UID = <string>
	//    MINS_TO_BYPASS_MFA = <integer>
	//    DISABLE_MFA = TRUE | FALSE
	//    MINS_TO_BYPASS_NETWORK POLICY = <integer>
}

func User() *schema.Resource {
	return &schema.Resource{
		Create: CreateUser,
		Read:   ReadUser,
		Delete: DeleteResource("user", snowflake.User),
		Update: UpdateResource("user", userProperties, userSchema, snowflake.User, ReadUser),

		Schema: userSchema,
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
		val, ok := data.GetOk(field)
		if ok {
			switch userSchema[field].Type {
			case schema.TypeString:
				valStr := val.(string)
				qb.SetString(field, valStr)
			case schema.TypeBool:
				valBool := val.(bool)
				qb.SetBool(field, valBool)
			}
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
	var name, createdOn, loginName, displayName, firstName, lastName, email, minsToUnlock, daysToExpiry, comment, mustChangePassword, snowflakeLock, defaultWarehouse, defaultNamespace, defaultRole, extAuthnDuo, extAuthnUID, minsToBypassMfa, owner, lastSuccessLogin, expiresAtTime, lockedUntilTime, hasPassword sql.NullString
	var disabled, hasRsaPublicKey bool
	err := row.Scan(&name, &createdOn, &loginName, &displayName, &firstName, &lastName, &email, &minsToUnlock, &daysToExpiry, &comment, &disabled, &mustChangePassword, &snowflakeLock, &defaultWarehouse, &defaultNamespace, &defaultRole, &extAuthnDuo, &extAuthnUID, &minsToBypassMfa, &owner, &lastSuccessLogin, &expiresAtTime, &lockedUntilTime, &hasPassword, &hasRsaPublicKey)
	if err != nil {
		return err
	}

	// TODO turn this into a loop after we switch to scaning in a struct
	err = data.Set("name", name.String)
	if err != nil {
		return err
	}
	err = data.Set("comment", comment.String)
	if err != nil {
		return err
	}

	err = data.Set("login_name", loginName.String)
	if err != nil {
		return err
	}

	err = data.Set("disabled", disabled)
	if err != nil {
		return err
	}

	err = data.Set("default_role", defaultRole.String)
	if err != nil {
		return err
	}

	err = data.Set("default_namespace", defaultNamespace.String)
	if err != nil {
		return err
	}

	err = data.Set("default_warehouse", defaultWarehouse.String)
	if err != nil {
		return err
	}

	err = data.Set("has_rsa_public_key", hasRsaPublicKey)

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

		for _, field := range changes {
			val := data.Get(field)
			switch userSchema[field].Type {
			case schema.TypeString:
				valStr := val.(string)
				qb.SetString(field, valStr)
			case schema.TypeBool:
				valBool := val.(bool)
				qb.SetBool(field, valBool)
			}
		}

		err := DBExec(db, qb.Statement())
		if err != nil {
			return errors.Wrap(err, "error altering user")
		}
	}
	return ReadUser(data, meta)
}
