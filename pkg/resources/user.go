package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var userSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of the user. Note that if you do not supply login_name this will be used as login_name. [doc](https://docs.snowflake.net/manuals/sql-reference/sql/create-user.html#required-parameters)",
	},
	"login_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Sensitive:   true,
		Description: "The name users use to log in. If not supplied, snowflake will use name instead.",
		// login_name is case-insensitive
		DiffSuppressFunc: ignoreCaseSuppressFunc,
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
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      "Specifies the virtual warehouse that is active by default for the user’s session upon login.",
	},
	// TODO [SNOW-1348101 - next PR]: check the exact behavior of default_namespace and default_role because it looks like it is handled in a case-insensitive manner on Snowflake side
	"default_namespace": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      "Specifies the namespace (database only or database and schema) that is active by default for the user’s session upon login.",
	},
	"default_role": {
		Type:             schema.TypeString,
		Optional:         true,
		Computed:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
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
	// TODO [SNOW-1348101 - next PR]: handle #1155 by either forceNew or not reading this value from SF (because it changes constantly after setting; check https://docs.snowflake.com/en/sql-reference/sql/create-user#optional-object-properties-objectproperties)
	//    DAYS_TO_EXPIRY = <integer>
	//    MINS_TO_UNLOCK = <integer>
	//    EXT_AUTHN_DUO = TRUE | FALSE
	//    EXT_AUTHN_UID = <string>
	//    MINS_TO_BYPASS_MFA = <integer>
	//    DISABLE_MFA = TRUE | FALSE
	//    MINS_TO_BYPASS_NETWORK POLICY = <integer>
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

		Schema: helpers.MergeMaps(userSchema, userParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func CreateUser(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	opts := &sdk.CreateUserOptions{
		ObjectProperties:  &sdk.UserObjectProperties{},
		ObjectParameters:  &sdk.UserObjectParameters{},
		SessionParameters: &sdk.SessionParameters{},
	}
	name := d.Get("name").(string)
	objectIdentifier := sdk.NewAccountObjectIdentifier(name)

	if parametersCreateDiags := handleUserParametersCreate(d, opts); len(parametersCreateDiags) > 0 {
		return parametersCreateDiags
	}

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
		opts.ObjectProperties.DefaultWarehouse = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(defaultWarehouse.(string)))
	}
	if defaultNamespace, ok := d.GetOk("default_namespace"); ok {
		defaultNamespaceId, err := helpers.DecodeSnowflakeParameterID(defaultNamespace.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		opts.ObjectProperties.DefaultNamespace = sdk.Pointer(defaultNamespaceId)
	}
	if displayName, ok := d.GetOk("display_name"); ok {
		opts.ObjectProperties.DisplayName = sdk.String(displayName.(string))
	}
	if defaultRole, ok := d.GetOk("default_role"); ok {
		opts.ObjectProperties.DefaultRole = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(defaultRole.(string)))
	}
	if v, ok := d.GetOk("default_secondary_roles"); ok {
		roles := expandStringList(v.(*schema.Set).List())
		secondaryRoles := []sdk.SecondaryRole{}
		for _, role := range roles {
			secondaryRoles = append(secondaryRoles, sdk.SecondaryRole{Value: role})
		}
		opts.ObjectProperties.DefaultSecondaryRoles = &sdk.SecondaryRoles{Roles: secondaryRoles}
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
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeSnowflakeID(objectIdentifier))
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
		if user.RsaPublicKeyFp != nil {
			if err = d.Set("has_rsa_public_key", user.RsaPublicKeyFp.Value != ""); err != nil {
				return diag.FromErr(err)
			}
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
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	if d.HasChange("name") {
		newID := sdk.NewAccountObjectIdentifier(d.Get("name").(string))

		err := client.Users.Alter(ctx, id, &sdk.AlterUserOptions{
			NewName: newID,
		})
		if err != nil {
			return diag.FromErr(err)
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
		if v, ok := d.GetOk("password"); ok {
			runSet = true
			alterOptions.Set.ObjectProperties.Password = sdk.String(v.(string))
		} else {
			// TODO [SNOW-1348101 - next PR]: this is temporary, update logic will be changed with the resource rework
			unsetOptions := &sdk.AlterUserOptions{
				Unset: &sdk.UserUnset{
					ObjectProperties: &sdk.UserObjectPropertiesUnset{
						Password: sdk.Bool(true),
					},
				},
			}
			err := client.Users.Alter(ctx, id, unsetOptions)
			if err != nil {
				d.Partial(true)
				return diag.FromErr(err)
			}
		}
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
		alterOptions.Set.ObjectProperties.DefaultWarehouse = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(n.(string)))
	}
	if d.HasChange("default_namespace") {
		runSet = true
		_, n := d.GetChange("default_namespace")
		defaultNamespaceId, err := helpers.DecodeSnowflakeParameterID(n.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		alterOptions.Set.ObjectProperties.DefaultNamespace = sdk.Pointer(defaultNamespaceId)
	}
	if d.HasChange("default_role") {
		runSet = true
		_, n := d.GetChange("default_role")
		alterOptions.Set.ObjectProperties.DefaultRole = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(n.(string)))
	}
	if d.HasChange("default_secondary_roles") {
		runSet = true
		_, n := d.GetChange("default_secondary_roles")
		roles := expandStringList(n.(*schema.Set).List())
		secondaryRoles := []sdk.SecondaryRole{}
		for _, role := range roles {
			secondaryRoles = append(secondaryRoles, sdk.SecondaryRole{Value: role})
		}
		alterOptions.Set.ObjectProperties.DefaultSecondaryRoles = &sdk.SecondaryRoles{Roles: secondaryRoles}
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
