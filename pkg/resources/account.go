package resources

import (
	"context"
	"errors"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Note: no test case was created for account since we cannot actually delete them after creation, which is a critical part of the test suite. Instead, this resource
// was manually tested

//var accountSchemaOld = map[string]*schema.Schema{
//	"name": {
//		Type:        schema.TypeString,
//		Required:    true,
//		Description: "Specifies the identifier (i.e. name) for the account; must be unique within an organization, regardless of which Snowflake Region the account is in. In addition, the identifier must start with an alphabetic character and cannot contain spaces or special characters except for underscores (_). Note that if the account name includes underscores, features that do not accept account names with underscores (e.g. Okta SSO or SCIM) can reference a version of the account name that substitutes hyphens (-) for the underscores.",
//		// Name is automatically uppercase by Snowflake
//		StateFunc: func(val interface{}) string {
//			return strings.ToUpper(val.(string))
//		},
//		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
//	},
//	"admin_name": {
//		Type:        schema.TypeString,
//		Required:    true,
//		Description: "Login name of the initial administrative user of the account. A new user is created in the new account with this name and password and granted the ACCOUNTADMIN role in the account. A login name can be any string consisting of letters, numbers, and underscores. Login names are always case-insensitive.",
//		// We have no way of assuming a role into this account to change the admin user name so this has to be ForceNew even though it's not ideal
//		ForceNew:              true,
//		DiffSuppressOnRefresh: true,
//		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
//			// For new resources always show the diff
//			if d.Id() == "" {
//				return false
//			}
//			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
//			return old == ""
//		},
//	},
//	"admin_password": {
//		Type:         schema.TypeString,
//		Optional:     true,
//		Sensitive:    true,
//		Description:  "Password for the initial administrative user of the account. Optional if the `ADMIN_RSA_PUBLIC_KEY` parameter is specified. For more information about passwords in Snowflake, see [Snowflake-provided Password Policy](https://docs.snowflake.com/en/sql-reference/sql/create-account.html#:~:text=Snowflake%2Dprovided%20Password%20Policy).",
//		AtLeastOneOf: []string{"admin_password", "admin_rsa_public_key"},
//		// We have no way of assuming a role into this account to change the password so this has to be ForceNew even though it's not ideal
//		ForceNew:              true,
//		DiffSuppressOnRefresh: true,
//		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
//			// For new resources always show the diff
//			if d.Id() == "" {
//				return false
//			}
//			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
//			return old == ""
//		},
//	},
//	"admin_rsa_public_key": {
//		Type:         schema.TypeString,
//		Optional:     true,
//		Sensitive:    true,
//		Description:  "Assigns a public key to the initial administrative user of the account in order to implement [key pair authentication](https://docs.snowflake.com/en/sql-reference/sql/create-account.html#:~:text=key%20pair%20authentication) for the user. Optional if the `ADMIN_PASSWORD` parameter is specified.",
//		AtLeastOneOf: []string{"admin_password", "admin_rsa_public_key"},
//		// We have no way of assuming a role into this account to change the admin rsa public key so this has to be ForceNew even though it's not ideal
//		ForceNew:              true,
//		DiffSuppressOnRefresh: true,
//		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
//			// For new resources always show the diff
//			if d.Id() == "" {
//				return false
//			}
//			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
//			return old == ""
//		},
//	},
//	"email": {
//		Type:        schema.TypeString,
//		Required:    true,
//		Sensitive:   true,
//		Description: "Email address of the initial administrative user of the account. This email address is used to send any notifications about the account.",
//		// We have no way of assuming a role into this account to change the admin email so this has to be ForceNew even though it's not ideal
//		ForceNew:              true,
//		DiffSuppressOnRefresh: true,
//		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
//			// For new resources always show the diff
//			if d.Id() == "" {
//				return false
//			}
//			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
//			return old == ""
//		},
//	},
//	"edition": {
//		Type:         schema.TypeString,
//		Required:     true,
//		ForceNew:     true,
//		Description:  "[Snowflake Edition](https://docs.snowflake.com/en/user-guide/intro-editions.html) of the account. Valid values are: STANDARD | ENTERPRISE | BUSINESS_CRITICAL",
//		ValidateFunc: validation.StringInSlice([]string{string(sdk.EditionStandard), string(sdk.EditionEnterprise), string(sdk.EditionBusinessCritical)}, false),
//	},
//	"first_name": {
//		Type:        schema.TypeString,
//		Optional:    true,
//		Sensitive:   true,
//		Description: "First name of the initial administrative user of the account",
//		// We have no way of assuming a role into this account to change the admin first name so this has to be ForceNew even though it's not ideal
//		ForceNew:              true,
//		DiffSuppressOnRefresh: true,
//		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
//			// For new resources always show the diff
//			if d.Id() == "" {
//				return false
//			}
//			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
//			return old == ""
//		},
//	},
//	"last_name": {
//		Type:        schema.TypeString,
//		Optional:    true,
//		Sensitive:   true,
//		Description: "Last name of the initial administrative user of the account",
//		// We have no way of assuming a role into this account to change the admin last name so this has to be ForceNew even though it's not ideal
//		ForceNew:              true,
//		DiffSuppressOnRefresh: true,
//		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
//			// For new resources always show the diff
//			if d.Id() == "" {
//				return false
//			}
//			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
//			return old == ""
//		},
//	},
//	"must_change_password": {
//		Type:        schema.TypeBool,
//		Optional:    true,
//		Default:     false,
//		Description: "Specifies whether the new user created to administer the account is forced to change their password upon first login into the account.",
//		// We have no way of assuming a role into this account to change the admin password policy so this has to be ForceNew even though it's not ideal
//		ForceNew:              true,
//		DiffSuppressOnRefresh: true,
//		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
//			// For new resources always show the diff
//			if d.Id() == "" {
//				return false
//			}
//			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
//			return old == ""
//		},
//	},
//	"region_group": {
//		Type:                  schema.TypeString,
//		Optional:              true,
//		Description:           "ID of the Snowflake Region where the account is created. If no value is provided, Snowflake creates the account in the same Snowflake Region as the current account (i.e. the account in which the CREATE ACCOUNT statement is executed.)",
//		ForceNew:              true,
//		DiffSuppressOnRefresh: true,
//		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
//			// For new resources always show the diff
//			if d.Id() == "" {
//				return false
//			}
//			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
//			return new == ""
//		},
//	},
//	"region": {
//		Type:                  schema.TypeString,
//		Optional:              true,
//		Description:           "ID of the Snowflake Region where the account is created. If no value is provided, Snowflake creates the account in the same Snowflake Region as the current account (i.e. the account in which the CREATE ACCOUNT statement is executed.)",
//		ForceNew:              true,
//		DiffSuppressOnRefresh: true,
//		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
//			// For new resources always show the diff
//			if d.Id() == "" {
//				return false
//			}
//			// This suppresses the diff if the old value is empty. This would happen in the event of importing existing accounts since we have no way of reading this value
//			return new == ""
//		},
//	},
//	"comment": {
//		Type:        schema.TypeString,
//		Optional:    true,
//		Description: "Specifies a comment for the account.",
//		ForceNew:    true,
//	},
//	"is_org_admin": {
//		Type:        schema.TypeBool,
//		Computed:    true,
//		Description: "Indicates whether the ORGADMIN role is enabled in an account. If TRUE, the role is enabled.",
//	},
//	"grace_period_in_days": {
//		Type:        schema.TypeInt,
//		Optional:    true,
//		Default:     3,
//		Description: "Specifies the number of days to wait before dropping the account. The default is 3 days.",
//	},
//	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
//}

var accountSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
		// TODO: Sensitive?
		Description:      "TODO",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
	},
	"admin_name": {
		Type:     schema.TypeString,
		Required: true,
		// TODO: Sensitive?
		Description:      externalChangesNotDetectedFieldDescription("TODO"),
		DiffSuppressFunc: IgnoreAfterCreation,
	},
	"admin_password": {
		Type:             schema.TypeString,
		Optional:         true,
		Sensitive:        true,
		Description:      externalChangesNotDetectedFieldDescription("TODO"),
		DiffSuppressFunc: IgnoreAfterCreation,
		AtLeastOneOf:     []string{"admin_password", "admin_rsa_public_key"},
	},
	"admin_rsa_public_key": {
		Type:             schema.TypeString,
		Optional:         true,
		Sensitive:        true,
		Description:      externalChangesNotDetectedFieldDescription("TODO"),
		DiffSuppressFunc: IgnoreAfterCreation,
		AtLeastOneOf:     []string{"admin_password", "admin_rsa_public_key"},
	},
	"admin_user_type": {
		Type:     schema.TypeString,
		Required: true,
		// TODO: Valid options
		Description:      externalChangesNotDetectedFieldDescription("TODO"),
		DiffSuppressFunc: SuppressIfAny(IgnoreAfterCreation, NormalizeAndCompare(sdk.ToUserType)),
		ValidateDiagFunc: sdkValidation(sdk.ToUserType),
	},
	"first_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Sensitive:        true,
		Description:      externalChangesNotDetectedFieldDescription("TODO"),
		DiffSuppressFunc: IgnoreAfterCreation,
	},
	"last_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Sensitive:        true,
		Description:      externalChangesNotDetectedFieldDescription("TODO"),
		DiffSuppressFunc: IgnoreAfterCreation,
	},
	"email": {
		Type:             schema.TypeString,
		Required:         true,
		Sensitive:        true,
		Description:      externalChangesNotDetectedFieldDescription("TODO"),
		DiffSuppressFunc: IgnoreAfterCreation,
	},
	"must_change_password": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		Description:      externalChangesNotDetectedFieldDescription("TODO"),
		DiffSuppressFunc: IgnoreAfterCreation,
		ValidateDiagFunc: validateBooleanString,
	},
	"edition": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "TODO",
		ValidateDiagFunc: sdkValidation(sdk.ToAccountEdition),
	},
	"region_group": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "TODO",
	},
	"region": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "TODO",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies a comment for the account.",
		DiffSuppressFunc: SuppressIfAny(
			IgnoreChangeToCurrentSnowflakeValueInShow("comment"),
			func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				return oldValue == "SNOWFLAKE" && newValue == ""
			},
		),
	},
	"is_org_admin": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("is_org_admin"),
		ValidateDiagFunc: validateBooleanString,
		Description:      "TODO",
	},
	"grace_period_in_days": {
		Type:             schema.TypeInt,
		Required:         true,
		Description:      "TODO",
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(3)),
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW ACCOUNTS` for the given account.",
		Elem: &schema.Resource{
			Schema: schemas.ShowAccountSchema,
		},
	},
}

func Account() *schema.Resource {
	return &schema.Resource{
		// TODO: Desc
		Description:   "The account resource allows you to create and manage Snowflake accounts.",
		CreateContext: TrackingCreateWrapper(resources.Account, CreateAccount),
		ReadContext:   TrackingReadWrapper(resources.Account, ReadAccount(true)),
		UpdateContext: TrackingUpdateWrapper(resources.Account, UpdateAccount),
		DeleteContext: TrackingDeleteWrapper(resources.Account, DeleteAccount),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.Account, customdiff.All(
			ComputedIfAnyAttributeChanged(accountSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(accountSchema, ShowOutputAttributeName, "name", "is_org_admin"),
		)),

		Schema: accountSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		// TODO: State upgrader and import
	}
}

func CreateAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	isOrgAdmin, err := client.ContextFunctions.IsRoleInSession(ctx, snowflakeroles.Orgadmin)
	if err != nil {
		return diag.FromErr(err)
	}
	if !isOrgAdmin {
		return diag.FromErr(errors.New("current user doesn't have the orgadmin role in session"))
	}

	id := sdk.NewAccountObjectIdentifier(d.Get("name").(string))

	opts := &sdk.CreateAccountOptions{
		AdminName: d.Get("admin_name").(string),
		Email:     d.Get("email").(string),
		Edition:   sdk.AccountEdition(d.Get("edition").(string)),
	}

	if v, ok := d.GetOk("admin_password"); ok {
		opts.AdminPassword = sdk.String(v.(string))
	}
	if v, ok := d.GetOk("admin_rsa_public_key"); ok {
		opts.AdminRSAPublicKey = sdk.String(v.(string))
	}
	if v, ok := d.GetOk("admin_user_type"); ok {
		userType, err := sdk.ToUserType(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		opts.AdminUserType = &userType
	}
	if v, ok := d.GetOk("first_name"); ok {
		opts.FirstName = sdk.String(v.(string))
	}
	if v, ok := d.GetOk("last_name"); ok {
		opts.LastName = sdk.String(v.(string))
	}
	if v := d.Get("must_change_password"); v != BooleanDefault {
		parsedBool, err := booleanStringToBool(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		opts.MustChangePassword = &parsedBool
	}
	if v, ok := d.GetOk("region_group"); ok {
		opts.RegionGroup = sdk.String(v.(string))
	}
	if v, ok := d.GetOk("region"); ok {
		opts.Region = sdk.String(v.(string))
	}
	if v, ok := d.GetOk("comment"); ok {
		opts.Comment = sdk.String(v.(string))
	}

	// TODO(TODO): next prs
	//if v := d.Get("polaris"); v != BooleanDefault {
	//	parsedBool, err := booleanStringToBool(v.(string))
	//	if err != nil {
	//		return diag.FromErr(err)
	//	}
	//	opts.Polaris = &parsedBool
	//}

	createResponse, err := client.Accounts.Create(ctx, id, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(sdk.NewAccountIdentifier(createResponse.OrganizationName, createResponse.AccountName)))

	if v, ok := d.GetOk("is_org_admin"); ok && v == BooleanTrue {
		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			SetIsOrgAdmin: &sdk.AccountSetIsOrgAdmin{
				Name:     id,
				OrgAdmin: true,
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadAccount(false)(ctx, d, meta)
}

func ReadAccount(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client

		isOrgAdmin, err := client.ContextFunctions.IsRoleInSession(ctx, snowflakeroles.Orgadmin)
		if err != nil {
			return diag.FromErr(err)
		}
		if !isOrgAdmin {
			return diag.FromErr(errors.New("current user doesn't have the orgadmin role in session"))
		}

		id, err := sdk.ParseAccountIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		account, err := client.Accounts.ShowByID(ctx, sdk.NewAccountObjectIdentifier(id.AccountName()))
		if err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			var regionGroup string
			if account.RegionGroup != nil {
				regionGroup = *account.RegionGroup
			}
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"is_org_admin", "is_org_admin", *account.IsOrgAdmin, booleanStringFromBool(*account.IsOrgAdmin), nil},
				outputMapping{"region_group", "region_group", regionGroup, regionGroup, nil},
				outputMapping{"snowflake_region", "region", account.SnowflakeRegion, account.SnowflakeRegion, nil},
				outputMapping{"comment", "comment", *account.Comment, *account.Comment, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err = setStateToValuesFromConfig(d, taskSchema, []string{
				"name",
				"admin_name",
				"admin_password",
				"admin_rsa_public_key",
				"admin_user_type",
				"first_name",
				"last_name",
				"email",
				"must_change_password",
				"edition",
				"region_group",
				"region",
				"comment",
				"is_org_admin",
				"grace_period_in_days",
			}); err != nil {
				return diag.FromErr(err)
			}
		}

		if errs := errors.Join(
			attributeMappedValueReadOrDefault(d, "edition", account.Edition, func(edition *sdk.AccountEdition) (string, error) {
				if edition != nil {
					return string(*edition), nil
				}
				return "", nil
			}, nil),
			// TODO: Region group is only returned when org is span on multiple region groups, but you can explicitly set it (e.g. PUBLIC)
			//attributeMappedValueReadOrNil(d, "region_group", account.RegionGroup, func(regionGroup *string) (string, error) {
			//	if regionGroup != nil {
			//		return *regionGroup, nil
			//	}
			//	return "", nil
			//}),
			// TODO: Can be left empty and it will be populated with current account's region
			//d.Set("region", account.SnowflakeRegion),
			// TODO: Default comment is "SNOWFLAKE"
			//attributeMappedValueReadOrNil(d, "comment", account.Comment, func(comment *string) (string, error) {
			//	if comment != nil {
			//		return *comment, nil
			//	}
			//	return "", nil
			//}),
			//attributeMappedValueReadOrNil(d, "is_org_admin", account.IsOrgAdmin, func(isOrgAdmin *bool) (string, error) {
			//	if isOrgAdmin != nil {
			//		return booleanStringFromBool(*isOrgAdmin), nil
			//	}
			//	return BooleanDefault, nil
			//}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.AccountToSchema(account)}),
		); errs != nil {
			return diag.FromErr(errs)
		}

		return nil
	}
}

func UpdateAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	isOrgAdmin, err := client.ContextFunctions.IsRoleInSession(ctx, snowflakeroles.Orgadmin)
	if err != nil {
		return diag.FromErr(err)
	}
	if !isOrgAdmin {
		return diag.FromErr(errors.New("current user doesn't have the orgadmin role in session"))
	}

	/*
		todo: comments may eventually work again for accounts, so this can be uncommented when that happens
		client := meta.(*provider.Context).Client
		client := sdk.NewClientFromDB(db)
		ctx := context.Background()

		id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

		// Change comment
		if d.HasChange("comment") {
			// changing comment isn't supported for accounts
			err := client.Comments.Set(ctx, &sdk.SetCommentOptions{
				ObjectType: sdk.ObjectTypeAccount,
				ObjectName: sdk.NewAccountObjectIdentifier(d.Get("name").(string)),
				Value:      sdk.String(d.Get("comment").(string)),
			})
			if err != nil {
				return err
			}
		}
	*/
	return nil
}

func DeleteAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	isOrgAdmin, err := client.ContextFunctions.IsRoleInSession(ctx, snowflakeroles.Orgadmin)
	if err != nil {
		return diag.FromErr(err)
	}
	if !isOrgAdmin {
		return diag.FromErr(errors.New("current user doesn't have the orgadmin role in session"))
	}

	id, err := sdk.ParseAccountIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Accounts.Drop(ctx, sdk.NewAccountObjectIdentifier(id.AccountName()), d.Get("grace_period_in_days").(int), &sdk.DropAccountOptions{
		IfExists: sdk.Bool(true),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
