package resources

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
		Type:             schema.TypeString,
		Optional:         true,
		Computed:         true,
		DiffSuppressFunc: diffCaseInsensitive,
		Description:      "Specifies the role that is active by default for the user’s session upon login.",
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
	client := meta.(*provider.Context).Client

	opts := &sdk.CreateUserOptions{
		ObjectProperties:  &sdk.UserObjectProperties{},
		ObjectParameters:  &sdk.UserObjectParameters{},
		SessionParameters: &sdk.SessionParameters{},
	}
	name := d.Get("name").(string)
	ctx := context.Background()
	objectIdentifier := sdk.NewAccountObjectIdentifier(name)

	if loginName, ok := d.GetOk("login_name"); ok {
		opts.ObjectProperties.LoginName = sdk.String(loginName.(string))
	}

	if comment, ok := d.GetOk("comment"); ok {
		opts.ObjectProperties.Comment = sdk.String(comment.(string))
	}
	if password, ok := d.GetOk("password"); ok {
		opts.ObjectProperties.Password = sdk.String(password.(string))
	}
	if v, ok := d.GetOk("disabled"); ok {
		disabled := v.(bool)
		opts.ObjectProperties.Disable = &disabled
	}
	if defaultWarehouse, ok := d.GetOk("default_warehouse"); ok {
		opts.ObjectProperties.DefaultWarehosue = sdk.String(defaultWarehouse.(string))
	}
	if defaultNamespace, ok := d.GetOk("default_namespace"); ok {
		opts.ObjectProperties.DefaultNamespace = sdk.String(defaultNamespace.(string))
	}
	if displayName, ok := d.GetOk("display_name"); ok {
		opts.ObjectProperties.DisplayName = sdk.String(displayName.(string))
	}
	if defaultRole, ok := d.GetOk("default_role"); ok {
		opts.ObjectProperties.DefaultRole = sdk.String(defaultRole.(string))
	}
	if v, ok := d.GetOk("default_secondary_roles"); ok {
		roles := expandStringList(v.(*schema.Set).List())
		secondaryRoles := []sdk.SecondaryRole{}
		for _, role := range roles {
			secondaryRoles = append(secondaryRoles, sdk.SecondaryRole{Value: role})
		}
		opts.ObjectProperties.DefaultSeconaryRoles = &sdk.SecondaryRoles{Roles: secondaryRoles}
	}
	if rsaPublicKey, ok := d.GetOk("rsa_public_key"); ok {
		opts.ObjectProperties.RSAPublicKey = sdk.String(rsaPublicKey.(string))
	}
	if rsaPublicKey2, ok := d.GetOk("rsa_public_key_2"); ok {
		opts.ObjectProperties.RSAPublicKey2 = sdk.String(rsaPublicKey2.(string))
	}
	if v, ok := d.GetOk("must_change_password"); ok {
		mustChangePassword := v.(bool)
		opts.ObjectProperties.MustChangePassword = &mustChangePassword
	}
	if email, ok := d.GetOk("email"); ok {
		opts.ObjectProperties.Email = sdk.String(email.(string))
	}
	if firstName, ok := d.GetOk("first_name"); ok {
		opts.ObjectProperties.FirstName = sdk.String(firstName.(string))
	}
	if lastName, ok := d.GetOk("last_name"); ok {
		opts.ObjectProperties.LastName = sdk.String(lastName.(string))
	}
	err := client.Users.Create(ctx, objectIdentifier, opts)
	if err != nil {
		return err
	}
	d.SetId(helpers.EncodeSnowflakeID(objectIdentifier))
	return ReadUser(d, meta)
}

func ReadUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	// We use User.Describe instead of User.Show because the "SHOW USERS ..." command
	// requires the "MANAGE GRANTS" global privilege
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	ctx := context.Background()
	user, err := client.Users.Describe(ctx, objectIdentifier)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
			log.Printf("[DEBUG] user (%s) not found or we are not authorized. Err: %s", d.Id(), err)
			d.SetId("")
			return nil
		}
		return err
	}

	if err := setStringProperty(d, "name", user.Name); err != nil {
		return err
	}
	if err := setStringProperty(d, "comment", user.Comment); err != nil {
		return err
	}
	if err := setStringProperty(d, "login_name", user.LoginName); err != nil {
		return err
	}
	if err := setBoolProperty(d, "disabled", user.Disabled); err != nil {
		return err
	}
	if err := setStringProperty(d, "default_role", user.DefaultRole); err != nil {
		return err
	}

	var defaultSecondaryRoles []string
	if user.DefaultSecondaryRoles != nil && len(user.DefaultSecondaryRoles.Value) > 0 {
		defaultRoles, _ := strings.CutPrefix(user.DefaultSecondaryRoles.Value, "[\"")
		defaultRoles, _ = strings.CutSuffix(defaultRoles, "\"]")
		defaultSecondaryRoles = strings.Split(defaultRoles, ",")
	}
	if err = d.Set("default_secondary_roles", defaultSecondaryRoles); err != nil {
		return err
	}
	if err := setStringProperty(d, "default_namespace", user.DefaultNamespace); err != nil {
		return err
	}
	if err := setStringProperty(d, "default_warehouse", user.DefaultWarehouse); err != nil {
		return err
	}
	if user.RsaPublicKeyFp != nil {
		if err = d.Set("has_rsa_public_key", user.RsaPublicKeyFp.Value != ""); err != nil {
			return err
		}
	}
	if err := setStringProperty(d, "email", user.Email); err != nil {
		return err
	}
	if err := setStringProperty(d, "display_name", user.DisplayName); err != nil {
		return err
	}
	if err := setStringProperty(d, "first_name", user.FirstName); err != nil {
		return err
	}
	if err := setStringProperty(d, "last_name", user.LastName); err != nil {
		return err
	}
	return nil
}

func UpdateUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client

	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	if d.HasChange("name") {
		_, n := d.GetChange("name")
		newName := n.(string)
		newID := sdk.NewAccountObjectIdentifier(newName)
		alterOptions := &sdk.AlterUserOptions{
			NewName: newID,
		}
		err := client.Users.Alter(ctx, id, alterOptions)
		if err != nil {
			return err
		}
		d.SetId(helpers.EncodeSnowflakeID(newID))
		id = newID
	}
	runSet := false
	alterOptions := &sdk.AlterUserOptions{
		Set: &sdk.UserSet{
			ObjectProperties: &sdk.UserObjectProperties{},
		},
	}
	if d.HasChange("login_name") {
		runSet = true
		_, n := d.GetChange("login_name")
		alterOptions.Set.ObjectProperties.LoginName = sdk.String(n.(string))
	}
	if d.HasChange("comment") {
		runSet = true
		_, n := d.GetChange("comment")
		alterOptions.Set.ObjectProperties.Comment = sdk.String(n.(string))
	}
	if d.HasChange("password") {
		runSet = true
		_, n := d.GetChange("password")
		alterOptions.Set.ObjectProperties.Password = sdk.String(n.(string))
	}

	if d.HasChange("disabled") {
		runSet = true
		_, n := d.GetChange("disabled")
		disabled := n.(bool)
		alterOptions.Set.ObjectProperties.Disable = &disabled
	}
	if d.HasChange("default_warehouse") {
		runSet = true
		_, n := d.GetChange("default_warehouse")
		alterOptions.Set.ObjectProperties.DefaultWarehosue = sdk.String(n.(string))
	}
	if d.HasChange("default_namespace") {
		runSet = true
		_, n := d.GetChange("default_namespace")
		alterOptions.Set.ObjectProperties.DefaultNamespace = sdk.String(n.(string))
	}
	if d.HasChange("default_role") {
		runSet = true
		_, n := d.GetChange("default_role")
		alterOptions.Set.ObjectProperties.DefaultRole = sdk.String(n.(string))
	}
	if d.HasChange("default_secondary_roles") {
		runSet = true
		_, n := d.GetChange("default_secondary_roles")
		roles := expandStringList(n.(*schema.Set).List())
		secondaryRoles := []sdk.SecondaryRole{}
		for _, role := range roles {
			secondaryRoles = append(secondaryRoles, sdk.SecondaryRole{Value: role})
		}
		alterOptions.Set.ObjectProperties.DefaultSeconaryRoles = &sdk.SecondaryRoles{Roles: secondaryRoles}
	}
	if d.HasChange("rsa_public_key") {
		runSet = true
		_, n := d.GetChange("rsa_public_key")
		alterOptions.Set.ObjectProperties.RSAPublicKey = sdk.String(n.(string))
	}
	if d.HasChange("rsa_public_key_2") {
		runSet = true
		_, n := d.GetChange("rsa_public_key_2")
		alterOptions.Set.ObjectProperties.RSAPublicKey2 = sdk.String(n.(string))
	}
	if d.HasChange("must_change_password") {
		runSet = true
		_, n := d.GetChange("must_change_password")
		mustChangePassword := n.(bool)
		alterOptions.Set.ObjectProperties.MustChangePassword = &mustChangePassword
	}
	if d.HasChange("email") {
		runSet = true
		_, n := d.GetChange("email")
		alterOptions.Set.ObjectProperties.Email = sdk.String(n.(string))
	}
	if d.HasChange("display_name") {
		runSet = true
		_, n := d.GetChange("display_name")
		alterOptions.Set.ObjectProperties.DisplayName = sdk.String(n.(string))
	}
	if d.HasChange("first_name") {
		runSet = true
		_, n := d.GetChange("first_name")
		alterOptions.Set.ObjectProperties.FirstName = sdk.String(n.(string))
	}
	if d.HasChange("last_name") {
		runSet = true
		_, n := d.GetChange("last_name")
		alterOptions.Set.ObjectProperties.LastName = sdk.String(n.(string))
	}
	if runSet {
		err := client.Users.Alter(ctx, id, alterOptions)
		if err != nil {
			return err
		}
	}

	return ReadUser(d, meta)
}

func DeleteUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	err := client.Users.Drop(ctx, objectIdentifier)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
