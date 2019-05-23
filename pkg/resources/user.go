package resources

import (
	"database/sql"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform/helper/schema"
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
		Update: UpdateUser,
		Delete: DeleteUser,
		Exists: UserExists,

		Schema: userSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// func DeleteResource(t string, builder func(string) *snowflake.Builder) func(*schema.ResourceData, interface{}) error {

func CreateUser(data *schema.ResourceData, meta interface{}) error {
	return CreateResource("user", userProperties, userSchema, snowflake.User, ReadUser)(data, meta)
}

func UserExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.User(id).Show()
	rows, err := db.Query(stmt)
	if err != nil {
		return false, err
	}

	if rows.Next() {
		return true, nil
	}
	return false, nil
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

func UpdateUser(data *schema.ResourceData, meta interface{}) error {
	return UpdateResource("user", userProperties, userSchema, snowflake.User, ReadUser)(data, meta)
}

func DeleteUser(data *schema.ResourceData, meta interface{}) error {
	return DeleteResource("user", snowflake.User)(data, meta)
}
