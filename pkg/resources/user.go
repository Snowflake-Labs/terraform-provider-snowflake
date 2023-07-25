package resources

import (
	"database/sql"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
)

var userProperties = []string{
	"comment",
	"login_name",
	"password",
	"disabled",
	"default_namespace",
	"default_role",
	"default_secondary_roles",
	"default_warehouse",
	"rsa_public_key",
	"rsa_public_key_2",
	"must_change_password",
	"email",
	"display_name",
	"first_name",
	"last_name",
}

var diffCaseInsensitive = func(k, old, new string, d *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}

var userSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Sensitive:   true,
		Description: "Name of the user. Note that if you do not supply login_name this will be used as login_name. [doc](https://docs.snowflake.net/manuals/sql-reference/sql/create-user.html#required-parameters)",
	},
	"login_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Sensitive:   false,
		Description: "The name users use to log in. If not supplied, snowflake will use name instead.",
		// login_name is case-insensitive
		DiffSuppressFunc: diffCaseInsensitive,
	},
	"comment": {
		Type:     schema.TypeString,
		Optional: true,
		// TODO validation
	},
	"password": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "**WARNING:** this will put the password in the terraform state file. Use carefully.",
		// TODO validation https://docs.snowflake.net/manuals/sql-reference/sql/create-user.html#optional-parameters
	},
	"disabled": {
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true,
	},
	"default_warehouse": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the virtual warehouse that is active by default for the user’s session upon login.",
	},
	"default_namespace": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: diffCaseInsensitive,
		Description:      "Specifies the namespace (database only or database and schema) that is active by default for the user’s session upon login.",
	},
	"default_role": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Specifies the role that is active by default for the user’s session upon login.",
	},
	"default_secondary_roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies the set of secondary roles that are active for the user’s session upon login. Currently only [\"ALL\"] value is supported - more information can be found in [doc](https://docs.snowflake.com/en/sql-reference/sql/create-user#optional-object-properties-objectproperties)",
	},
	"rsa_public_key": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the user’s RSA public key; used for key-pair authentication. Must be on 1 line without header and trailer.",
	},
	"rsa_public_key_2": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the user’s second RSA public key; used to rotate the public and private keys for key-pair authentication based on an expiration schedule set by your organization. Must be on 1 line without header and trailer.",
	},
	"has_rsa_public_key": {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Will be true if user as an RSA key set.",
	},
	"must_change_password": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether the user is forced to change their password on next login (including their first/initial login) into the system.",
	},
	"email": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "Email address for the user.",
	},
	"display_name": {
		Type:        schema.TypeString,
		Computed:    true,
		Optional:    true,
		Sensitive:   true,
		Description: "Name displayed for the user in the Snowflake web interface.",
	},
	"first_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "First name of the user.",
	},
	"last_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "Last name of the user.",
	},
	"tag": tagReferenceSchema,

	//    MIDDLE_NAME = <string>
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

		Schema: userSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateUser(d *schema.ResourceData, meta interface{}) error {
	return CreateResource("user", userProperties, userSchema, snowflake.NewUserBuilder, ReadUser)(d, meta)
}

func ReadUser(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	// We use User.Describe instead of User.Show because the "SHOW USERS ..." command
	// requires the "MANAGE GRANTS" global privilege
	stmt := snowflake.NewUserBuilder(d.Id()).Describe()
	rows, err := snowflake.Query(db, stmt)
	if err != nil {
		if snowflake.IsResourceNotExistOrNotAuthorized(err.Error(), "User") {
			// If not found, mark resource to be removed from state file during apply or refresh
			log.Printf("[DEBUG] user (%s) not found or we are not authorized.Err:\n%s", d.Id(), err.Error())
			d.SetId("")
			return nil
		}
		return err
	}

	u, err := snowflake.ScanUserDescription(rows)
	if err != nil {
		return err
	}
	if err = d.Set("name", u.Name.String); err != nil {
		return err
	}
	if err = d.Set("comment", u.Comment.String); err != nil {
		return err
	}
	if err = d.Set("login_name", u.LoginName.String); err != nil {
		return err
	}
	if err = d.Set("disabled", u.Disabled); err != nil {
		return err
	}
	if err = d.Set("default_role", u.DefaultRole.String); err != nil {
		return err
	}

	var defaultSecondaryRoles []string
	if len(u.DefaultSecondaryRoles.String) > 0 {
		defaultSecondaryRoles = strings.Split(u.DefaultSecondaryRoles.String, ",")
	}
	if err = d.Set("default_secondary_roles", defaultSecondaryRoles); err != nil {
		return err
	}
	if err = d.Set("default_namespace", u.DefaultNamespace.String); err != nil {
		return err
	}
	if err = d.Set("default_warehouse", u.DefaultWarehouse.String); err != nil {
		return err
	}
	if err = d.Set("has_rsa_public_key", u.HasRsaPublicKey); err != nil {
		return err
	}
	if err = d.Set("email", u.Email.String); err != nil {
		return err
	}
	if err = d.Set("display_name", u.DisplayName.String); err != nil {
		return err
	}
	if err = d.Set("first_name", u.FirstName.String); err != nil {
		return err
	}
	if err = d.Set("last_name", u.LastName.String); err != nil {
		return err
	}
	return nil
}

func UpdateUser(d *schema.ResourceData, meta interface{}) error {
	return UpdateResource("user", userProperties, userSchema, snowflake.NewUserBuilder, ReadUser)(d, meta)
}

func DeleteUser(d *schema.ResourceData, meta interface{}) error {
	return DeleteResource("user", snowflake.NewUserBuilder)(d, meta)
}
