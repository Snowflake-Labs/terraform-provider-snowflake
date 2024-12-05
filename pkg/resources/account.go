package resources

import (
	"context"
	"errors"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

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
			StateContext: TrackingImportWrapper(resources.Account, ImportAccount),
		},

		// TODO: State upgrader
	}
}

func ImportAccount(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client

	isOrgAdmin, err := client.ContextFunctions.IsRoleInSession(ctx, snowflakeroles.Orgadmin)
	if err != nil {
		return nil, err
	}
	if !isOrgAdmin {
		// TODO:
		return nil, errors.New("current user doesn't have the orgadmin role in session")
	}

	id, err := sdk.ParseAccountIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	account, err := client.Accounts.ShowByID(ctx, id.AccountId())
	if err != nil {
		return nil, err
	}

	if _, err := ImportName[sdk.AccountIdentifier](context.Background(), d, nil); err != nil {
		return nil, err
	}

	if account.RegionGroup != nil {
		if err = d.Set("region_group", *account.RegionGroup); err != nil {
			return nil, err
		}
	}

	if err := errors.Join(
		d.Set("edition", string(*account.Edition)),
		d.Set("region", account.SnowflakeRegion),
		d.Set("comment", *account.Comment),
		d.Set("is_org_admin", booleanStringFromBool(*account.IsOrgAdmin)),
	); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
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

		account, err := client.Accounts.ShowByID(ctx, id.AccountId())
		if err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			var regionGroup string
			if account.RegionGroup != nil {
				regionGroup = *account.RegionGroup
			}
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"edition", "edition", *account.Edition, *account.Edition, nil},
				outputMapping{"is_org_admin", "is_org_admin", *account.IsOrgAdmin, booleanStringFromBool(*account.IsOrgAdmin), nil},
				outputMapping{"region_group", "region_group", regionGroup, regionGroup, nil},
				outputMapping{"snowflake_region", "region", account.SnowflakeRegion, account.SnowflakeRegion, nil},
				outputMapping{"comment", "comment", *account.Comment, *account.Comment, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err = setStateToValuesFromConfig(d, accountSchema, []string{
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

	id, err := sdk.ParseAccountIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewAccountIdentifier(id.OrganizationName(), d.Get("name").(string))

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Rename: &sdk.AccountRename{
				Name:    id.AccountId(),
				NewName: newId.AccountId(),
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	if d.HasChange("is_org_admin") {
		if v := d.Get("is_org_admin").(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
				SetIsOrgAdmin: &sdk.AccountSetIsOrgAdmin{
					Name:     id.AccountId(),
					OrgAdmin: parsed,
				},
			}); err != nil {
				return diag.FromErr(err)
			}
		} else {
			// No unset available for this field (setting Snowflake default)
			if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
				SetIsOrgAdmin: &sdk.AccountSetIsOrgAdmin{
					Name:     id.AccountId(),
					OrgAdmin: false,
				},
			}); err != nil && !strings.Contains(err.Error(), "already has ORGADMIN disabled") { // TODO: What to do about this error?
				return diag.FromErr(err)
			}
		}
	}

	return ReadAccount(false)(ctx, d, meta)
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

	err = client.Accounts.Drop(ctx, id.AccountId(), d.Get("grace_period_in_days").(int), &sdk.DropAccountOptions{
		IfExists: sdk.Bool(true),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
