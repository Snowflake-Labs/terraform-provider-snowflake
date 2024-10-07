package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var userSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Name of the user. Note that if you do not supply login_name this will be used as login_name. Check the [docs](https://docs.snowflake.net/manuals/sql-reference/sql/create-user.html#required-parameters)."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"password": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "Password for the user. **WARNING:** this will put the password in the terraform state file. Use carefully.",
	},
	"login_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Sensitive:        true,
		DiffSuppressFunc: SuppressIfAny(ignoreCaseSuppressFunc, IgnoreChangeToCurrentSnowflakeValueInShow("login_name")),
		Description:      "The name users use to log in. If not supplied, snowflake will use name instead. Login names are always case-insensitive.",
		// login_name is case-insensitive
	},
	"display_name": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("display_name"),
		Description:      "Name displayed for the user in the Snowflake web interface.",
	},
	"first_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "First name of the user.",
	},
	"middle_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "Middle name of the user.",
	},
	"last_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "Last name of the user.",
	},
	"email": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "Email address for the user.",
	},
	"must_change_password": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("must_change_password"),
		Description:      booleanStringFieldDescription("Specifies whether the user is forced to change their password on next login (including their first/initial login) into the system."),
		Default:          BooleanDefault,
	},
	"disabled": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("disabled"),
		Description:      booleanStringFieldDescription("Specifies whether the user is disabled, which prevents logging in and aborts all the currently-running queries for the user."),
		Default:          BooleanDefault,
	},
	// TODO [SNOW-1649000]: consider handling external change if there is no config (or zero) for `days_to_expiry` and other similar attributes (what about this the other way around?)
	"days_to_expiry": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: externalChangesNotDetectedFieldDescription("Specifies the number of days after which the user status is set to `Expired` and the user is no longer allowed to log in. This is useful for defining temporary users (i.e. users who should only have access to Snowflake for a limited time period). In general, you should not set this property for [account administrators](https://docs.snowflake.com/en/user-guide/security-access-control-considerations.html#label-accountadmin-users) (i.e. users with the `ACCOUNTADMIN` role) because Snowflake locks them out when they become `Expired`."),
	},
	"mins_to_unlock": {
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(0),
		Default:      IntDefault,
		Description:  externalChangesNotDetectedFieldDescription("Specifies the number of minutes until the temporary lock on the user login is cleared. To protect against unauthorized user login, Snowflake places a temporary lock on a user after five consecutive unsuccessful login attempts. When creating a user, this property can be set to prevent them from logging in until the specified amount of time passes. To remove a lock immediately for a user, specify a value of 0 for this parameter. **Note** because this value changes continuously after setting it, the provider is currently NOT handling the external changes to it."),
	},
	"default_warehouse": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      "Specifies the virtual warehouse that is active by default for the user’s session upon login. Note that the CREATE USER operation does not verify that the warehouse exists.",
	},
	"default_namespace": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInShow("default_namespace")),
		Description:      "Specifies the namespace (database only or database and schema) that is active by default for the user’s session upon login. Note that the CREATE USER operation does not verify that the namespace exists.",
	},
	"default_role": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      "Specifies the role that is active by default for the user’s session upon login. Note that specifying a default role for a user does **not** grant the role to the user. The role must be granted explicitly to the user using the [GRANT ROLE](https://docs.snowflake.com/en/sql-reference/sql/grant-role) command. In addition, the CREATE USER operation does not verify that the role exists.",
	},
	"default_secondary_roles_option": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          sdk.SecondaryRolesOptionDefault,
		ValidateDiagFunc: sdkValidation(sdk.ToSecondaryRolesOption),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToSecondaryRolesOption), IgnoreChangeToCurrentSnowflakeValueInShowWithMapping("default_secondary_roles", func(x any) any {
			return sdk.GetSecondaryRolesOptionFrom(x.(string))
		})),
		Description: fmt.Sprintf("Specifies the secondary roles that are active for the user’s session upon login. Valid values are (case-insensitive): %s. More information can be found in [doc](https://docs.snowflake.com/en/sql-reference/sql/create-user#optional-object-properties-objectproperties).", possibleValuesListed(sdk.ValidSecondaryRolesOptionsString)),
	},
	"mins_to_bypass_mfa": {
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(0),
		Default:      IntDefault,
		Description:  externalChangesNotDetectedFieldDescription("Specifies the number of minutes to temporarily bypass MFA for the user. This property can be used to allow a MFA-enrolled user to temporarily bypass MFA during login in the event that their MFA device is not available."),
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
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the user.",
	},
	"disable_mfa": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		Description:      externalChangesNotDetectedFieldDescription(booleanStringFieldDescription("Allows enabling or disabling [multi-factor authentication](https://docs.snowflake.com/en/user-guide/security-mfa).")),
		Default:          BooleanDefault,
	},
	"user_type": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Specifies a type for the user.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW USER` for the given user.",
		Elem: &schema.Resource{
			Schema: schemas.ShowUserSchema,
		},
	},
	ParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PARAMETERS IN USER` for the given user.",
		Elem: &schema.Resource{
			Schema: schemas.ShowUserParametersSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func User() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: CreateUser,
		UpdateContext: UpdateUser,
		ReadContext:   GetReadUserFunc(true),
		DeleteContext: DeleteUser,
		Description:   "Resource used to manage user objects. For more information, check [user documentation](https://docs.snowflake.com/en/sql-reference/commands-user-role).",

		Schema: helpers.MergeMaps(userSchema, userParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: ImportUser,
		},

		CustomizeDiff: customdiff.All(
			// TODO [SNOW-1629468 - next pr]: test "default_role", "default_secondary_roles"
			ComputedIfAnyAttributeChanged(userSchema, ShowOutputAttributeName, "password", "login_name", "display_name", "first_name", "last_name", "email", "must_change_password", "disabled", "days_to_expiry", "mins_to_unlock", "default_warehouse", "default_namespace", "default_role", "default_secondary_roles_option", "mins_to_bypass_mfa", "rsa_public_key", "rsa_public_key_2", "comment", "disable_mfa"),
			ComputedIfAnyAttributeChanged(userParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllUserParameters), strings.ToLower)...),
			ComputedIfAnyAttributeChanged(userSchema, FullyQualifiedNameAttributeName, "name"),
			userParametersCustomDiff,
			// TODO [SNOW-1645348]: revisit with service user work
			func(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
				if n := diff.Get("user_type"); n != nil {
					logging.DebugLogger.Printf("[DEBUG] new external value for user type %s\n", n.(string))
					if !slices.Contains([]string{"", "PERSON"}, strings.ToUpper(n.(string))) {
						return errors.Join(diff.SetNewComputed("user_type"), diff.ForceNew("user_type"))
					}
				}
				return nil
			},
		),

		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v094UserStateUpgrader,
			},
		},
	}
}

func ServiceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateUser,
		UpdateContext: UpdateUser,
		ReadContext:   GetReadUserFunc(true),
		DeleteContext: DeleteUser,
		Description:   "Resource used to manage service user objects. For more information, check [user documentation](https://docs.snowflake.com/en/sql-reference/commands-user-role).",

		Schema: helpers.MergeMaps(serviceUserSchema, userParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: ImportUser,
		},

		CustomizeDiff: customdiff.All(
			// TODO [SNOW-1645348]: generalize this list
			ComputedIfAnyAttributeChanged(userSchema, ShowOutputAttributeName, "login_name", "display_name", "email", "must_change_password", "disabled", "days_to_expiry", "mins_to_unlock", "default_warehouse", "default_namespace", "default_role", "default_secondary_roles_option", "rsa_public_key", "rsa_public_key_2", "comment"),
			ComputedIfAnyAttributeChanged(userParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllUserParameters), strings.ToLower)...),
			ComputedIfAnyAttributeChanged(userSchema, FullyQualifiedNameAttributeName, "name"),
			userParametersCustomDiff,
			// TODO [SNOW-1645348]: revisit with service user work
			func(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
				if n := diff.Get("user_type"); n != nil {
					logging.DebugLogger.Printf("[DEBUG] new external value for user type %s\n", n.(string))
					if !slices.Contains([]string{"SERVICE"}, strings.ToUpper(n.(string))) {
						return errors.Join(diff.SetNewComputed("user_type"), diff.ForceNew("user_type"))
					}
				}
				return nil
			},
		),
	}
}

func ImportUser(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	logging.DebugLogger.Printf("[DEBUG] Starting user import")
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	userDetails, err := client.Users.Describe(ctx, id)
	if err != nil {
		return nil, err
	}

	u, err := client.Users.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = errors.Join(
		d.Set("name", id.Name()),
		setFromStringPropertyIfNotEmpty(d, "login_name", userDetails.LoginName),
		setFromStringPropertyIfNotEmpty(d, "display_name", userDetails.DisplayName),
		setFromStringPropertyIfNotEmpty(d, "default_namespace", userDetails.DefaultNamespace),
		setBooleanStringFromBoolProperty(d, "must_change_password", userDetails.MustChangePassword),
		setBooleanStringFromBoolProperty(d, "disabled", userDetails.Disabled),
		d.Set("default_secondary_roles_option", u.GetSecondaryRolesOption()),
		// all others are set in read
	)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateUser(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	opts := &sdk.CreateUserOptions{
		ObjectProperties:  &sdk.UserObjectProperties{},
		ObjectParameters:  &sdk.UserObjectParameters{},
		SessionParameters: &sdk.SessionParameters{},
	}
	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)

	errs := errors.Join(
		stringAttributeCreate(d, "password", &opts.ObjectProperties.Password),
		stringAttributeCreate(d, "login_name", &opts.ObjectProperties.LoginName),
		stringAttributeCreate(d, "display_name", &opts.ObjectProperties.DisplayName),
		stringAttributeCreate(d, "first_name", &opts.ObjectProperties.FirstName),
		stringAttributeCreate(d, "middle_name", &opts.ObjectProperties.MiddleName),
		stringAttributeCreate(d, "last_name", &opts.ObjectProperties.LastName),
		stringAttributeCreate(d, "email", &opts.ObjectProperties.Email),
		booleanStringAttributeCreate(d, "must_change_password", &opts.ObjectProperties.MustChangePassword),
		booleanStringAttributeCreate(d, "disabled", &opts.ObjectProperties.Disable),
		intAttributeCreate(d, "days_to_expiry", &opts.ObjectProperties.DaysToExpiry),
		intAttributeWithSpecialDefaultCreate(d, "mins_to_unlock", &opts.ObjectProperties.MinsToUnlock),
		accountObjectIdentifierAttributeCreate(d, "default_warehouse", &opts.ObjectProperties.DefaultWarehouse),
		objectIdentifierAttributeCreate(d, "default_namespace", &opts.ObjectProperties.DefaultNamespace),
		accountObjectIdentifierAttributeCreate(d, "default_role", &opts.ObjectProperties.DefaultRole),
		func() error {
			defaultSecondaryRolesOption, err := sdk.ToSecondaryRolesOption(d.Get("default_secondary_roles_option").(string))
			if err != nil {
				return err
			}
			switch defaultSecondaryRolesOption {
			case sdk.SecondaryRolesOptionDefault:
				return nil
			case sdk.SecondaryRolesOptionNone:
				opts.ObjectProperties.DefaultSecondaryRoles = &sdk.SecondaryRoles{None: sdk.Bool(true)}
			case sdk.SecondaryRolesOptionAll:
				opts.ObjectProperties.DefaultSecondaryRoles = &sdk.SecondaryRoles{All: sdk.Bool(true)}
			}
			return nil
		}(),
		intAttributeWithSpecialDefaultCreate(d, "mins_to_bypass_mfa", &opts.ObjectProperties.MinsToBypassMFA),
		stringAttributeCreate(d, "rsa_public_key", &opts.ObjectProperties.RSAPublicKey),
		stringAttributeCreate(d, "rsa_public_key_2", &opts.ObjectProperties.RSAPublicKey2),
		stringAttributeCreate(d, "comment", &opts.ObjectProperties.Comment),
		// disable mfa cannot be set in create, alter is run after creation
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if parametersCreateDiags := handleUserParametersCreate(d, opts); len(parametersCreateDiags) > 0 {
		return parametersCreateDiags
	}

	err := client.Users.Create(ctx, id, opts)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	// disable mfa cannot be set in create, we need to alter if set in config
	var diags diag.Diagnostics
	if disableMfa := d.Get("disable_mfa").(string); disableMfa != BooleanDefault {
		parsed, err := booleanStringToBool(disableMfa)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Setting disable mfa failed after create for user %s, err: %v", id.FullyQualifiedName(), err),
			})
		}
		alterDisableMfa := sdk.AlterUserOptions{Set: &sdk.UserSet{ObjectProperties: &sdk.UserAlterObjectProperties{DisableMfa: sdk.Bool(parsed)}}}
		err = client.Users.Alter(ctx, id, &alterDisableMfa)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Setting disable mfa failed after create for user %s, err: %v", id.FullyQualifiedName(), err),
			})
		}
	}

	return append(diags, GetReadUserFunc(false)(ctx, d, meta)...)
}

func GetReadUserFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		userDetails, err := client.Users.Describe(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				log.Printf("[DEBUG] user (%s) not found or we are not authorized. Err: %s", d.Id(), err)
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query user. Marking the resource as removed.",
						Detail:   fmt.Sprintf("User: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		u, err := client.Users.ShowByID(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query user. Marking the resource as removed.",
						Detail:   fmt.Sprintf("User: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}
		userParameters, err := client.Users.ShowParameters(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				showMapping{"login_name", "login_name", u.LoginName, u.LoginName, nil},
				showMapping{"display_name", "display_name", u.DisplayName, u.DisplayName, nil},
				showMapping{"must_change_password", "must_change_password", u.MustChangePassword, fmt.Sprintf("%t", u.MustChangePassword), nil},
				showMapping{"disabled", "disabled", u.Disabled, fmt.Sprintf("%t", u.Disabled), nil},
				showMapping{"default_namespace", "default_namespace", u.DefaultNamespace, u.DefaultNamespace, nil},
				showMapping{"default_secondary_roles", "default_secondary_roles_option", u.DefaultSecondaryRoles, u.GetSecondaryRolesOption(), nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, userSchema, []string{
			"login_name",
			"display_name",
			"must_change_password",
			"disabled",
			"default_namespace",
		}); err != nil {
			return diag.FromErr(err)
		}

		errs := errors.Join(
			// not reading name on purpose (we never update the name externally)
			// can't read password
			// not reading login_name on purpose (handled as external change to show output)
			// not reading display_name on purpose (handled as external change to show output)
			setFromStringPropertyIfNotEmpty(d, "first_name", userDetails.FirstName),
			setFromStringPropertyIfNotEmpty(d, "middle_name", userDetails.MiddleName),
			setFromStringPropertyIfNotEmpty(d, "last_name", userDetails.LastName),
			setFromStringPropertyIfNotEmpty(d, "email", userDetails.Email),
			// not reading must_change_password on purpose (handled as external change to show output)
			// not reading disabled on purpose (handled as external change to show output)
			// not reading days_to_expiry on purpose (they always change)
			// not reading mins_to_unlock on purpose (they always change)
			setFromStringPropertyIfNotEmpty(d, "default_warehouse", userDetails.DefaultWarehouse),
			// not reading default_namespace because one-part namespace seems to be capitalized on Snowflake side (handled as external change to show output)
			setFromStringPropertyIfNotEmpty(d, "default_role", userDetails.DefaultRole),
			// not setting default_secondary_role_option (handled as external change to show output)
			// not reading mins_to_bypass_mfa on purpose (they always change)
			setFromStringPropertyIfNotEmpty(d, "rsa_public_key", userDetails.RsaPublicKey),
			setFromStringPropertyIfNotEmpty(d, "rsa_public_key_2", userDetails.RsaPublicKey2),
			setFromStringPropertyIfNotEmpty(d, "comment", userDetails.Comment),
			// can't read disable_mfa
			d.Set("user_type", u.Type),

			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			handleUserParameterRead(d, userParameters),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.UserToSchema(u)}),
			d.Set(ParametersAttributeName, []map[string]any{schemas.UserParametersToSchema(userParameters)}),
		)
		if errs != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func UpdateUser(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newID := sdk.NewAccountObjectIdentifier(d.Get("name").(string))

		err := client.Users.Alter(ctx, id, &sdk.AlterUserOptions{
			NewName: newID,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeResourceIdentifier(newID))
		id = newID
	}

	setObjectProperties := sdk.UserAlterObjectProperties{}
	unsetObjectProperties := sdk.UserObjectPropertiesUnset{}
	errs := errors.Join(
		stringAttributeUpdate(d, "password", &setObjectProperties.Password, &unsetObjectProperties.Password),
		stringAttributeUpdate(d, "login_name", &setObjectProperties.LoginName, &unsetObjectProperties.LoginName),
		stringAttributeUpdate(d, "display_name", &setObjectProperties.DisplayName, &unsetObjectProperties.DisplayName),
		stringAttributeUpdate(d, "first_name", &setObjectProperties.FirstName, &unsetObjectProperties.FirstName),
		stringAttributeUpdate(d, "middle_name", &setObjectProperties.MiddleName, &unsetObjectProperties.MiddleName),
		stringAttributeUpdate(d, "last_name", &setObjectProperties.LastName, &unsetObjectProperties.LastName),
		stringAttributeUpdate(d, "email", &setObjectProperties.Email, &unsetObjectProperties.Email),
		booleanStringAttributeUpdate(d, "must_change_password", &setObjectProperties.MustChangePassword, &unsetObjectProperties.MustChangePassword),
		booleanStringAttributeUpdate(d, "disabled", &setObjectProperties.Disable, &unsetObjectProperties.Disable),
		intAttributeUpdate(d, "days_to_expiry", &setObjectProperties.DaysToExpiry, &unsetObjectProperties.DaysToExpiry),
		intAttributeWithSpecialDefaultUpdate(d, "mins_to_unlock", &setObjectProperties.MinsToUnlock, &unsetObjectProperties.MinsToUnlock),
		accountObjectIdentifierAttributeUpdate(d, "default_warehouse", &setObjectProperties.DefaultWarehouse, &unsetObjectProperties.DefaultWarehouse),
		objectIdentifierAttributeUpdate(d, "default_namespace", &setObjectProperties.DefaultNamespace, &unsetObjectProperties.DefaultNamespace),
		accountObjectIdentifierAttributeUpdate(d, "default_role", &setObjectProperties.DefaultRole, &unsetObjectProperties.DefaultRole),
		func() error {
			if d.HasChange("default_secondary_roles_option") {
				defaultSecondaryRolesOption, err := sdk.ToSecondaryRolesOption(d.Get("default_secondary_roles_option").(string))
				if err != nil {
					return err
				}
				switch defaultSecondaryRolesOption {
				case sdk.SecondaryRolesOptionDefault:
					unsetObjectProperties.DefaultSecondaryRoles = sdk.Bool(true)
				case sdk.SecondaryRolesOptionNone:
					setObjectProperties.DefaultSecondaryRoles = &sdk.SecondaryRoles{None: sdk.Bool(true)}
				case sdk.SecondaryRolesOptionAll:
					setObjectProperties.DefaultSecondaryRoles = &sdk.SecondaryRoles{All: sdk.Bool(true)}
				}
			}
			return nil
		}(),
		intAttributeWithSpecialDefaultUpdate(d, "mins_to_bypass_mfa", &setObjectProperties.MinsToBypassMFA, &unsetObjectProperties.MinsToBypassMFA),
		stringAttributeUpdate(d, "rsa_public_key", &setObjectProperties.RSAPublicKey, &unsetObjectProperties.RSAPublicKey),
		stringAttributeUpdate(d, "rsa_public_key_2", &setObjectProperties.RSAPublicKey2, &unsetObjectProperties.RSAPublicKey2),
		stringAttributeUpdate(d, "comment", &setObjectProperties.Comment, &unsetObjectProperties.Comment),
		booleanStringAttributeUpdate(d, "disable_mfa", &setObjectProperties.DisableMfa, &unsetObjectProperties.DisableMfa),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if (setObjectProperties != sdk.UserAlterObjectProperties{}) {
		err := client.Users.Alter(ctx, id, &sdk.AlterUserOptions{Set: &sdk.UserSet{ObjectProperties: &setObjectProperties}})
		if err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}
	if (unsetObjectProperties != sdk.UserObjectPropertiesUnset{}) {
		err := client.Users.Alter(ctx, id, &sdk.AlterUserOptions{Unset: &sdk.UserUnset{ObjectProperties: &unsetObjectProperties}})
		if err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}

	set := &sdk.UserSet{
		SessionParameters: &sdk.SessionParameters{},
		ObjectParameters:  &sdk.UserObjectParameters{},
	}
	unset := &sdk.UserUnset{
		SessionParameters: &sdk.SessionParametersUnset{},
		ObjectParameters:  &sdk.UserObjectParametersUnset{},
	}
	if updateParamDiags := handleUserParametersUpdate(d, set, unset); len(updateParamDiags) > 0 {
		return updateParamDiags
	}

	if (*set.SessionParameters != sdk.SessionParameters{} || *set.ObjectParameters != sdk.UserObjectParameters{}) {
		err := client.Users.Alter(ctx, id, &sdk.AlterUserOptions{
			Set: set,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if (*unset.SessionParameters != sdk.SessionParametersUnset{}) || (*unset.ObjectParameters != sdk.UserObjectParametersUnset{}) {
		err := client.Users.Alter(ctx, id, &sdk.AlterUserOptions{
			Unset: &sdk.UserUnset{
				SessionParameters: unset.SessionParameters,
				ObjectParameters:  unset.ObjectParameters,
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return GetReadUserFunc(false)(ctx, d, meta)
}

func DeleteUser(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	err := client.Users.Drop(ctx, id, &sdk.DropUserOptions{
		IfExists: sdk.Bool(true),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
