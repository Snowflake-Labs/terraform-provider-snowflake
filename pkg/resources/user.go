package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// TODO [SNOW-1348101]: handle external type change properly (force new)
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
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "The name users use to log in. If not supplied, snowflake will use name instead. Login names are always case-insensitive.",
		// login_name is case-insensitive
		DiffSuppressFunc: ignoreCaseSuppressFunc,
	},
	// TODO [SNOW-1348101]: handle external changes and the default behavior correctly; same with the login_name
	"display_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Name displayed for the user in the Snowflake web interface.",
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
	// TODO [SNOW-1348101]: handle this properly
	"must_change_password": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		Description:      booleanStringFieldDescription("Specifies whether the user is forced to change their password on next login (including their first/initial login) into the system."),
		Default:          BooleanDefault,
	},
	"disabled": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		Description:      booleanStringFieldDescription("Specifies whether the user is disabled, which prevents logging in and aborts all the currently-running queries for the user."),
		Default:          BooleanDefault,
	},
	// TODO [SNOW-1348101 - next PR]: handle #1155 by either forceNew or not reading this value from SF (because it changes constantly after setting; check https://docs.snowflake.com/en/sql-reference/sql/create-user#optional-object-properties-objectproperties)
	// TODO [SNOW-1348101]: negative values can be set by hand so IntDefault should probably not be used here
	"days_to_expiry": {
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(0),
		Default:      IntDefault,
		Description:  "Specifies the number of days after which the user status is set to `Expired` and the user is no longer allowed to log in. This is useful for defining temporary users (i.e. users who should only have access to Snowflake for a limited time period). In general, you should not set this property for [account administrators](https://docs.snowflake.com/en/user-guide/security-access-control-considerations.html#label-accountadmin-users) (i.e. users with the `ACCOUNTADMIN` role) because Snowflake locks them out when they become `Expired`.",
	},
	"mins_to_unlock": {
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(0),
		Default:      IntDefault,
		Description:  "Specifies the number of minutes until the temporary lock on the user login is cleared. To protect against unauthorized user login, Snowflake places a temporary lock on a user after five consecutive unsuccessful login attempts. When creating a user, this property can be set to prevent them from logging in until the specified amount of time passes. To remove a lock immediately for a user, specify a value of 0 for this parameter.",
	},
	"default_warehouse": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      "Specifies the virtual warehouse that is active by default for the user’s session upon login. Note that the CREATE USER operation does not verify that the warehouse exists.",
	},
	// TODO [SNOW-1348101 - next PR]: check the exact behavior of default_namespace and default_role because it looks like it is handled in a case-insensitive manner on Snowflake side
	"default_namespace": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      "Specifies the namespace (database only or database and schema) that is active by default for the user’s session upon login. Note that the CREATE USER operation does not verify that the namespace exists.",
	},
	"default_role": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      "Specifies the role that is active by default for the user’s session upon login. Note that specifying a default role for a user does **not** grant the role to the user. The role must be granted explicitly to the user using the [GRANT ROLE](https://docs.snowflake.com/en/sql-reference/sql/grant-role) command. In addition, the CREATE USER operation does not verify that the role exists.",
	},
	// TODO [SNOW-1348101]: test (no elems, more elems, duplicated elems, proper setting, update - both ways and external one)
	"default_secondary_roles": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: isValidSecondaryRole(),
		},
		MaxItems:    1,
		MinItems:    1,
		Optional:    true,
		Description: "Specifies the set of secondary roles that are active for the user’s session upon login. Currently only [\"ALL\"] value is supported - more information can be found in [doc](https://docs.snowflake.com/en/sql-reference/sql/create-user#optional-object-properties-objectproperties).",
	},
	// TODO [SNOW-1348101]: note that external changes are not handled
	"mins_to_bypass_mfa": {
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(1),
		Description:  "Specifies the number of minutes to temporarily bypass MFA for the user. This property can be used to allow a MFA-enrolled user to temporarily bypass MFA during login in the event that their MFA device is not available.",
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
	// TODO [SNOW-1348101]: handle properly
	"disable_mfa": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		Description:      booleanStringFieldDescription("Allows enabling or disabling [multi-factor authentication](https://docs.snowflake.com/en/user-guide/security-mfa)."),
		Default:          BooleanDefault,
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
			// TODO [SNOW-1348101]: fill after adding all the attributes
			// ComputedIfAnyAttributeChanged(ShowOutputAttributeName),
			// TODO [SNOW-1348101]: use list from user parameters definition instead listing all here
			ComputedIfAnyAttributeChanged(ParametersAttributeName, strings.ToLower(string(sdk.UserParameterAbortDetachedQuery)), strings.ToLower(string(sdk.UserParameterAutocommit)), strings.ToLower(string(sdk.UserParameterBinaryInputFormat)), strings.ToLower(string(sdk.UserParameterBinaryOutputFormat)), strings.ToLower(string(sdk.UserParameterClientMemoryLimit)), strings.ToLower(string(sdk.UserParameterClientMetadataRequestUseConnectionCtx)), strings.ToLower(string(sdk.UserParameterClientPrefetchThreads)), strings.ToLower(string(sdk.UserParameterClientResultChunkSize)), strings.ToLower(string(sdk.UserParameterClientResultColumnCaseInsensitive)), strings.ToLower(string(sdk.UserParameterClientSessionKeepAlive)), strings.ToLower(string(sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency)), strings.ToLower(string(sdk.UserParameterClientTimestampTypeMapping)), strings.ToLower(string(sdk.UserParameterDateInputFormat)), strings.ToLower(string(sdk.UserParameterDateOutputFormat)), strings.ToLower(string(sdk.UserParameterEnableUnloadPhysicalTypeOptimization)), strings.ToLower(string(sdk.UserParameterErrorOnNondeterministicMerge)), strings.ToLower(string(sdk.UserParameterErrorOnNondeterministicUpdate)), strings.ToLower(string(sdk.UserParameterGeographyOutputFormat)), strings.ToLower(string(sdk.UserParameterGeometryOutputFormat)), strings.ToLower(string(sdk.UserParameterJdbcTreatDecimalAsInt)), strings.ToLower(string(sdk.UserParameterJdbcTreatTimestampNtzAsUtc)), strings.ToLower(string(sdk.UserParameterJdbcUseSessionTimezone)), strings.ToLower(string(sdk.UserParameterJsonIndent)), strings.ToLower(string(sdk.UserParameterLockTimeout)), strings.ToLower(string(sdk.UserParameterLogLevel)), strings.ToLower(string(sdk.UserParameterMultiStatementCount)), strings.ToLower(string(sdk.UserParameterNoorderSequenceAsDefault)), strings.ToLower(string(sdk.UserParameterOdbcTreatDecimalAsInt)), strings.ToLower(string(sdk.UserParameterQueryTag)), strings.ToLower(string(sdk.UserParameterQuotedIdentifiersIgnoreCase)), strings.ToLower(string(sdk.UserParameterRowsPerResultset)), strings.ToLower(string(sdk.UserParameterS3StageVpceDnsName)), strings.ToLower(string(sdk.UserParameterSearchPath)), strings.ToLower(string(sdk.UserParameterSimulatedDataSharingConsumer)), strings.ToLower(string(sdk.UserParameterStatementQueuedTimeoutInSeconds)), strings.ToLower(string(sdk.UserParameterStatementTimeoutInSeconds)), strings.ToLower(string(sdk.UserParameterStrictJsonOutput)), strings.ToLower(string(sdk.UserParameterTimestampDayIsAlways24h)), strings.ToLower(string(sdk.UserParameterTimestampInputFormat)), strings.ToLower(string(sdk.UserParameterTimestampLtzOutputFormat)), strings.ToLower(string(sdk.UserParameterTimestampNtzOutputFormat)), strings.ToLower(string(sdk.UserParameterTimestampOutputFormat)), strings.ToLower(string(sdk.UserParameterTimestampTypeMapping)), strings.ToLower(string(sdk.UserParameterTimestampTzOutputFormat)), strings.ToLower(string(sdk.UserParameterTimezone)), strings.ToLower(string(sdk.UserParameterTimeInputFormat)), strings.ToLower(string(sdk.UserParameterTimeOutputFormat)), strings.ToLower(string(sdk.UserParameterTraceLevel)), strings.ToLower(string(sdk.UserParameterTransactionAbortOnError)), strings.ToLower(string(sdk.UserParameterTransactionDefaultIsolationLevel)), strings.ToLower(string(sdk.UserParameterTwoDigitCenturyStart)), strings.ToLower(string(sdk.UserParameterUnsupportedDdlAction)), strings.ToLower(string(sdk.UserParameterUseCachedResult)), strings.ToLower(string(sdk.UserParameterWeekOfYearPolicy)), strings.ToLower(string(sdk.UserParameterWeekStart)), strings.ToLower(string(sdk.UserParameterEnableUnredactedQuerySyntaxError)), strings.ToLower(string(sdk.UserParameterNetworkPolicy)), strings.ToLower(string(sdk.UserParameterPreventUnloadToInternalStages))),
			ComputedIfAnyAttributeChanged(FullyQualifiedNameAttributeName, "name"),
			userParametersCustomDiff,
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

	err = errors.Join(
		d.Set("name", id.Name()),
		// password can't be set
		setStringProperty(d, "login_name", userDetails.LoginName),
		setStringProperty(d, "display_name", userDetails.DisplayName),
		setStringProperty(d, "first_name", userDetails.FirstName),
		setStringProperty(d, "middle_name", userDetails.MiddleName),
		setStringProperty(d, "last_name", userDetails.LastName),
		setStringProperty(d, "email", userDetails.Email),
		// TODO: must_change_password
		// TODO: disabled
		// TODO: days_to_expiry
		// TODO: mins_to_unlock
		// TODO: default_warehouse
		// TODO: default_namespace
		// TODO: default_role
		// TODO: default_secondary_roles
		// TODO: mins_to_bypass_mfa
		// TODO: rsa_public_key
		// TODO: rsa_public_key_2
		setStringProperty(d, "comment", userDetails.Comment),
		// TODO: disable_mfa
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
		intAttributeCreate(d, "mins_to_unlock", &opts.ObjectProperties.MinsToUnlock),
		accountObjectIdentifierAttributeCreate(d, "default_warehouse", &opts.ObjectProperties.DefaultWarehouse),
		objectIdentifierAttributeCreate(d, "default_namespace", &opts.ObjectProperties.DefaultNamespace),
		accountObjectIdentifierAttributeCreate(d, "default_role", &opts.ObjectProperties.DefaultRole),
		// We do not need value because it is validated on the schema level and ALL is the only supported value currently.
		// Check more in https://docs.snowflake.com/en/sql-reference/sql/create-user#optional-object-properties-objectproperties.
		attributeDirectValueCreate(d, "default_secondary_roles", &opts.ObjectProperties.DefaultSecondaryRoles, &sdk.SecondaryRoles{}),
		intAttributeCreate(d, "mins_to_bypass_mfa", &opts.ObjectProperties.MinsToBypassMFA),
		stringAttributeCreate(d, "rsa_public_key", &opts.ObjectProperties.RSAPublicKey),
		stringAttributeCreate(d, "rsa_public_key_2", &opts.ObjectProperties.RSAPublicKey2),
		stringAttributeCreate(d, "comment", &opts.ObjectProperties.Comment),
		// TODO: handle disable_mfa (not settable in create - check)
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
	return GetReadUserFunc(false)(ctx, d, meta)
}

func GetReadUserFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		// We use User.Describe instead of User.Show because the "SHOW USERS ..." command
		// requires the "MANAGE GRANTS" global privilege
		id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
		user, err := client.Users.Describe(ctx, id)
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

		if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
			return diag.FromErr(err)
		}

		if err := setStringProperty(d, "name", user.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := setStringProperty(d, "comment", user.Comment); err != nil {
			return diag.FromErr(err)
		}
		if err := setStringProperty(d, "login_name", user.LoginName); err != nil {
			return diag.FromErr(err)
		}
		if err := setBoolProperty(d, "disabled", user.Disabled); err != nil {
			return diag.FromErr(err)
		}
		if err := setStringProperty(d, "default_role", user.DefaultRole); err != nil {
			return diag.FromErr(err)
		}

		// TODO [SNOW-1348101]: do we need to read them (probably yes, to handle external changes properly)? Do we need diff suppression for lowercase inside the config?
		var defaultSecondaryRoles []string
		if user.DefaultSecondaryRoles != nil && len(user.DefaultSecondaryRoles.Value) > 0 {
			defaultSecondaryRoles = sdk.ParseCommaSeparatedStringArray(user.DefaultSecondaryRoles.Value, true)
		}
		if err = d.Set("default_secondary_roles", defaultSecondaryRoles); err != nil {
			return diag.FromErr(err)
		}
		if err := setStringProperty(d, "default_namespace", user.DefaultNamespace); err != nil {
			return diag.FromErr(err)
		}
		if err := setStringProperty(d, "default_warehouse", user.DefaultWarehouse); err != nil {
			return diag.FromErr(err)
		}
		if err := setStringProperty(d, "email", user.Email); err != nil {
			return diag.FromErr(err)
		}
		if err := setStringProperty(d, "display_name", user.DisplayName); err != nil {
			return diag.FromErr(err)
		}
		if err := setStringProperty(d, "first_name", user.FirstName); err != nil {
			return diag.FromErr(err)
		}
		if err := setStringProperty(d, "last_name", user.LastName); err != nil {
			return diag.FromErr(err)
		}

		if diags := handleUserParameterRead(d, userParameters); diags != nil {
			return diags
		}

		if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.UserToSchema(u)}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(ParametersAttributeName, []map[string]any{schemas.UserParametersToSchema(userParameters)}); err != nil {
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
		intAttributeUpdate(d, "mins_to_unlock", &setObjectProperties.MinsToUnlock, &unsetObjectProperties.MinsToUnlock),
		accountObjectIdentifierAttributeUpdate(d, "default_warehouse", &setObjectProperties.DefaultWarehouse, &unsetObjectProperties.DefaultWarehouse),
		objectIdentifierAttributeUpdate(d, "default_namespace", &setObjectProperties.DefaultNamespace, &unsetObjectProperties.DefaultNamespace),
		accountObjectIdentifierAttributeUpdate(d, "default_role", &setObjectProperties.DefaultRole, &unsetObjectProperties.DefaultRole),
		// We do not need value because it is validated on the schema level and ALL is the only supported value currently.
		// Check more in https://docs.snowflake.com/en/sql-reference/sql/create-user#optional-object-properties-objectproperties.
		attributeDirectValueUpdate(d, "default_secondary_roles", &setObjectProperties.DefaultSecondaryRoles, &sdk.SecondaryRoles{}, &unsetObjectProperties.DefaultSecondaryRoles),
		intAttributeUpdate(d, "mins_to_bypass_mfa", &setObjectProperties.MinsToBypassMFA, &unsetObjectProperties.MinsToBypassMFA),
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
	// unset is split into two because:
	// 1. this is how it's written in the docs https://docs.snowflake.com/en/sql-reference/sql/alter-user#syntax
	// 2. current implementation of sdk.UserUnset makes distinction between user and session parameters,
	// so adding a comma between them is not trivial in the current SQL builder implementation
	if (*unset.SessionParameters != sdk.SessionParametersUnset{}) {
		err := client.Users.Alter(ctx, id, &sdk.AlterUserOptions{
			Unset: &sdk.UserUnset{
				SessionParameters: unset.SessionParameters,
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if (*unset.ObjectParameters != sdk.UserObjectParametersUnset{}) {
		err := client.Users.Alter(ctx, id, &sdk.AlterUserOptions{
			Unset: &sdk.UserUnset{
				ObjectParameters: unset.ObjectParameters,
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
